package stats

import (
  "appengine"
  "appengine/urlfetch"
  "encoding/base64"
  "html/template"
  "hash/fnv"
  "net/http"
  "strings"
  "net/url"
  "io"
)

const base64GifPixel = "R0lGODlhAQABAIAAAP///wAAACwAAAAAAQABAAACAkQBADs="
var pageTemplate, _ = template.New("page").ParseFiles("repostats/page.html")

func init() {
  http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  c.Infof("Requested URL: %v", r.URL)

  // / -> redirect
  // /account -> account template
  // /account/page -> GIF hit
  params := strings.SplitN(strings.Trim(r.URL.Path,"/"), "/", 2)

  if len(params) == 0 {
    http.Redirect(w, r, "https://github.com/igrigorik/repostats", http.StatusFound)
    return
  }

  if len(params) == 1 {
    account := map[string]interface{}{"Account": params[0]}
    if err := pageTemplate.ExecuteTemplate(w, "page", account); err != nil {
      panic("Cannot execute template")
    }

  } else {
    w.Header().Set("Content-Type", "image/gif")
    w.Header().Set("Cache-Control", "no-cache")

    output, _ := base64.StdEncoding.DecodeString(base64GifPixel)
    io.WriteString(w, string(output))

    go func(){
      h := fnv.New32a()
      h.Write([]byte(strings.Split(r.RemoteAddr,":")[0])) // IP address
      h.Write([]byte(r.Header.Get("User-Agent")))
      cid := h.Sum32()

      client := urlfetch.Client(c)
      resp, err := client.PostForm("http://www.google-analytics.com/collect",
          url.Values{
            "v": {"1"},           // protocol version = 1
            "t": {"pageview"},    // hit type
            "tid": {params[0]},   // tracking / property ID
            "cid": {string(cid)}, // unique client ID (IP + UA hash)
        })

      if err != nil {
        c.Errorf("GA collector POST error: %s", err.Error())
      }
      c.Infof("GA collector status: %v, cid: %v", resp.Status, cid)
    }()
  }

  c.Infof("Params size: %v, %v", len(params), params)
  return
}
