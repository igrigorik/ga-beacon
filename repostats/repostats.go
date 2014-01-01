package stats

import (
  "appengine"
  // "appengine/urlfetch"
  "html/template"
  "encoding/base64"
  // "crypto/md5"
  "net/http"
  "strings"
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
  }

  c.Infof("Params size: %v, %v", len(params), params)
  return
}
