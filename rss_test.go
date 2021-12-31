package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseRSS(t *testing.T) {

	rssThirdPartyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/basic" {
			rssXml, err := ioutil.ReadFile("fixtures/basic-rss.xml")
			if err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(rssXml)
		} else if r.URL.Path == "/invalidformat" {
			rssXml, err := ioutil.ReadFile("fixtures/invalid.xml")
			if err != nil {
				t.Fatal(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(rssXml)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(``))
		}
	}))
	defer rssThirdPartyServer.Close()

	t.Run("test basic RSS feed (simple and correct)", func(t *testing.T) {
		r := httptest.NewRequest(
			http.MethodPost,
			"/basic",
			strings.NewReader(fmt.Sprintf("%s%s", rssThirdPartyServer.URL, "/basic")))
		w := httptest.NewRecorder()
		requestHandler(w, r)
		feed, _ := io.ReadAll(w.Result().Body)
		expectedFeed := `{"title":"test RSS feed title","items":[{"title":"article 1","desciption":"one","link":"https://example.com/1","pubDate":""},{"title":"article 2","desciption":"two","link":"https://example.com/2","pubDate":""}]}`
		if string(feed) != expectedFeed {
			t.Fatalf("expected feed %v does not equal parsed feed %s", expectedFeed, feed)
		}
	})

	t.Run("test not found rss url", func(t *testing.T) {
		r := httptest.NewRequest(
			http.MethodPost,
			"/",
			strings.NewReader(fmt.Sprintf("%s%s", rssThirdPartyServer.URL, "/doesnotexist")))
		w := httptest.NewRecorder()
		requestHandler(w, r)
		if w.Result().StatusCode != 400 {
			t.Fatal("Status code returned was not 400")
		}
		content, _ := io.ReadAll(w.Result().Body)
		if string(content) != "error in request, could not parse rss from url" {
			t.Fatal("Message returned was incorrect")
		}
	})

	t.Run("test invalid rss url", func(t *testing.T) {
		r := httptest.NewRequest(
			http.MethodPost,
			"/",
			strings.NewReader("this :// is not a url"))
		w := httptest.NewRecorder()
		requestHandler(w, r)
		if w.Result().StatusCode != 400 {
			t.Fatal("Status code returned was not 400")
		}
		content, _ := io.ReadAll(w.Result().Body)
		if string(content) != "error in request, url passed in body was not valid" {
			t.Fatal("Message returned was incorrect")
		}
	})

	t.Run("test invalid RSS feed", func(t *testing.T) {
		r := httptest.NewRequest(
			http.MethodPost,
			"/invalidformat",
			strings.NewReader(fmt.Sprintf("%s%s", rssThirdPartyServer.URL, "/invalidformat")))
		w := httptest.NewRecorder()
		requestHandler(w, r)
		if w.Result().StatusCode != 400 {
			t.Fatal("Status code returned was not 400")
		}
		content, _ := io.ReadAll(w.Result().Body)
		if string(content) != "error in request, could not parse rss from url" {
			t.Fatal("Message returned was incorrect")
		}
	})
}
