package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	defaultTID = "UA-XXXXX-X"
	defaultDP  = "homepage"
	defaultURL = "/" + defaultTID + "/" + defaultDP
)

// Record the response from the handler function
func record(req *http.Request) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	handler(rec, req)
	return rec
}

// Create a request for a URL
func newRequest(t *testing.T, url string) *http.Request {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

// Record the response for a URL
func recordURL(t *testing.T, url string) *httptest.ResponseRecorder {
	return record(newRequest(t, url))
}

// Read a file
func readFile(t *testing.T, path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Read an input stream
func readAll(t *testing.T, r io.Reader) []byte {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// Test the correct data is sent to Google Analytics
func testBeaconRequest(t *testing.T, r *http.Request, tid string, dp string) {
	assert.Equal(t, "POST", r.Method)
	assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
	err := r.ParseForm()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "1", r.PostForm.Get("v"))
	assert.Equal(t, "pageview", r.PostForm.Get("t"))
	assert.Equal(t, tid, r.PostForm.Get("tid"))
	assert.Equal(t, dp, r.PostForm.Get("dp"))
	// The requests are not actually sent, so should not have an IP address
	assert.Empty(t, r.PostForm.Get("uip"))
}

// Test the tracking request
func testTrackRequest(t *testing.T, tid string, dp string, cid string, req *http.Request) *http.Response {
	if req == nil {
		req = newRequest(t, "/"+tid+"/"+dp)
	}
	var cidf string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		testBeaconRequest(t, r, tid, dp)
		cidf = r.PostForm.Get("cid")
		if cid != "" {
			assert.Equal(t, cid, cidf)
		}
	}))
	defer server.Close()
	beaconURL = server.URL
	res := record(req).Result()
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, "no-cache, no-store, must-revalidate, private", res.Header.Get("Cache-Control"))
	expires, err := time.Parse(http.TimeFormat, res.Header.Get("Expires"))
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, expires.Before(time.Now()))
	if cid == "" {
		// ensure cid cookie has been set
		var cidc string
		for _, c := range res.Cookies() {
			if c.Name == "cid" {
				cidc = c.Value
				break
			}
		}
		assert.Equal(t, cidf, cidc)
	}
	return res
}

func TestBeacon(t *testing.T) {
	t.Run("redirect on no params", func(t *testing.T) {
		rec := recordURL(t, "/")
		assert.Equal(t, http.StatusFound, rec.Code)
		assert.Equal(t, "https://github.com/igrigorik/ga-beacon", rec.Header().Get("Location"))
	})
	t.Run("account page", func(t *testing.T) {
		req := newRequest(t, "/UA-XXXXX-X")
		req.Header.Set("Referer", "https://example.com")
		rec := record(req)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, readFile(t, "page_test.html"), rec.Body.Bytes())
	})
	t.Run("badge", func(t *testing.T) {
		res := testTrackRequest(t, defaultTID, defaultDP, "", nil)
		assert.Equal(t, "image/svg+xml", res.Header.Get("Content-Type"))
		assert.Equal(t, readFile(t, "static/badge.svg"), readAll(t, res.Body))
	})
	t.Run("pixel", func(t *testing.T) {
		res := testTrackRequest(t, defaultTID, defaultDP, "", newRequest(t, defaultURL+"?pixel"))
		assert.Equal(t, "image/gif", res.Header.Get("Content-Type"))
		assert.Equal(t, readFile(t, "static/pixel.gif"), readAll(t, res.Body))
	})
	t.Run("badge gif", func(t *testing.T) {
		res := testTrackRequest(t, defaultTID, defaultDP, "", newRequest(t, defaultURL+"?gif"))
		assert.Equal(t, "image/gif", res.Header.Get("Content-Type"))
		assert.Equal(t, readFile(t, "static/badge.gif"), readAll(t, res.Body))
	})
	t.Run("badge flat", func(t *testing.T) {
		res := testTrackRequest(t, defaultTID, defaultDP, "", newRequest(t, defaultURL+"?flat"))
		assert.Equal(t, "image/svg+xml", res.Header.Get("Content-Type"))
		assert.Equal(t, readFile(t, "static/badge-flat.svg"), readAll(t, res.Body))
	})
	t.Run("badge flat gif", func(t *testing.T) {
		res := testTrackRequest(t, defaultTID, defaultDP, "", newRequest(t, defaultURL+"?flat-gif"))
		assert.Equal(t, "image/gif", res.Header.Get("Content-Type"))
		assert.Equal(t, readFile(t, "static/badge-flat.gif"), readAll(t, res.Body))
	})
	t.Run("referer as path", func(t *testing.T) {
		req := newRequest(t, "/"+defaultTID+"?useReferer")
		t.Run("warn on missing referer", func(t *testing.T) {
			rec := record(req)
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), "could not extract referer from headers")
		})
		dp := "example.com"
		req.Header.Set("Referer", "https://"+dp)
		testTrackRequest(t, defaultTID, dp, "", req)
	})
	t.Run("existing cid", func(t *testing.T) {
		req := newRequest(t, defaultURL)
		cid := "5d7b632fef264b76a7938362e5aba2c8"
		req.AddCookie(&http.Cookie{
			Name:  "cid",
			Value: cid,
			Path:  "/" + defaultTID,
		})
		testTrackRequest(t, defaultTID, defaultDP, cid, req)
	})
}
