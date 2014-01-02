package beacon

import (
	"appengine"
	"appengine/urlfetch"
	"fmt"
	"hash/fnv"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const beaconURL = "http://www.google-analytics.com/collect"

var pixel, _ = ioutil.ReadFile("static/pixel.gif")
var badge, _ = ioutil.ReadFile("static/badge.gif")
var pageTemplate, _ = template.New("page").ParseFiles("ga-beacon/page.html")

func init() {
	http.HandleFunc("/", handler)
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
	// /account/page -> GIF + log pageview to GA collector
	if len(params) == 1 {
		account := map[string]interface{}{"Account": params[0]}
		if err := pageTemplate.ExecuteTemplate(w, "page", account); err != nil {
			panic("Cannot execute template")
		}

	} else {
		hash := fnv.New32a()
		hash.Write([]byte(strings.Split(r.RemoteAddr, ":")[0]))
		hash.Write([]byte(r.Header.Get("User-Agent")))
		cid := fmt.Sprintf("%d", hash.Sum32())

		w.Header().Set("Content-Type", "image/gif")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("CID", cid)

		payload := url.Values{
			"v":   {"1"},        // protocol version = 1
			"t":   {"pageview"}, // hit type
			"tid": {params[0]},  // tracking / property ID
			"cid": {cid},        // unique client ID (IP + UA hash)
			"dp":  {params[1]},  // page path
		}

		req, _ := http.NewRequest("POST", beaconURL,
			strings.NewReader(payload.Encode()))
		req.Header.Add("User-Agent", r.Header.Get("User-Agent"))

		client := urlfetch.Client(c)
		resp, err := client.Do(req)
		if err != nil {
			c.Errorf("GA collector POST error: %s", err.Error())
		}
		c.Infof("GA collector status: %v, cid: %v", resp.Status, cid)

		// write out GIF pixel
		query, _ := url.ParseQuery(r.URL.RawQuery)
		_, ok_pixel := query["pixel"]

		if ok_pixel {
			io.WriteString(w, string(pixel))
		} else {
			io.WriteString(w, string(badge))
		}
	}

	return
}
