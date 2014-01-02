package beacon

import (
	"fmt"
	"hash/fnv"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"appengine"
	"appengine/urlfetch"
)

const beaconURL = "http://www.google-analytics.com/collect"

var (
	pixel        = mustReadFile("static/pixel.gif")
	badge        = mustReadFile("static/badge.gif")
	pageTemplate = template.Must(template.New("page").ParseFiles("ga-beacon/page.html"))
)

func init() {
	http.HandleFunc("/", handler)
}

func mustReadFile(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return b
}

func handler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	params := strings.SplitN(strings.Trim(r.URL.Path, "/"), "/", 2)

	// / -> redirect
	if len(params[0]) == 0 {
		http.Redirect(w, r, "https://github.com/igrigorik/ga-beacon", http.StatusFound)
		return
	}

	// /account -> account template
	if len(params) == 1 {
		templateParams := struct {
			Account string
		}{
			Account: params[0],
		}
		if err := pageTemplate.ExecuteTemplate(w, "page.html", templateParams); err != nil {
			http.Error(w, "could not show account page", 500)
			c.Errorf("Cannot execute template: %v", err)
		}
		return
	}

	// /account/page -> GIF + log pageview to GA collector

	hash := fnv.New32a()
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err != nil {
		hash.Write([]byte(host))
	} else {
		c.Errorf("could not parse addr: %q", r.RemoteAddr)
		hash.Write([]byte(r.RemoteAddr))
	}
	hash.Write([]byte(r.Header.Get("User-Agent")))
	cid := fmt.Sprintf("%d", hash.Sum32())

	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("CID", cid)

	// https://developers.google.com/analytics/devguides/collection/protocol/v1/reference
	payload := url.Values{
		"v":   {"1"},        // protocol version = 1
		"t":   {"pageview"}, // hit type
		"tid": {params[0]},  // tracking / property ID
		"cid": {cid},        // unique client ID (IP + UA hash)
		"dp":  {params[1]},  // page path
	}

	req, _ := http.NewRequest("POST", beaconURL, strings.NewReader(payload.Encode()))
	req.Header.Add("User-Agent", r.Header.Get("User-Agent"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if resp, err := urlfetch.Client(c).Do(req); err != nil {
		c.Errorf("GA collector POST error: %s", err.Error())
	} else {
		c.Debugf("GA collector status: %v, cid: %v", resp.Status, cid)
	}

	// Write out GIF pixel or badge, based on presence of "pixel" param.
	query, _ := url.ParseQuery(r.URL.RawQuery)
	if _, ok := query["pixel"]; ok {
		w.Write(pixel)
	} else {
		w.Write(badge)
	}
}
