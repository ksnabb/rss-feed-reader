package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

func requestHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	u, err := url.Parse(string(body))
	if err != nil {
		w.WriteHeader(400)
		if _, err := w.Write([]byte("error in request, url passed in body was not valid")); err != nil {
			log.Fatal(err)
		}
		return
	}
	feed, err := parseFeedFromUrl(u)
	if err != nil {
		w.WriteHeader(400)
		if _, err := w.Write([]byte("error in request, could not parse rss from url")); err != nil {
			log.Fatal(err)
		}
		return
	}
	b, err := json.Marshal(feed)
	if _, err := w.Write(b); err != nil {
		w.WriteHeader(400)
		if _, err := w.Write([]byte("error in request url, could not parse the feed,")); err != nil {
			log.Fatal(err)
		}
		return
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func main() {
	http.HandleFunc("/", requestHandler)
	http.HandleFunc("/ping", pingHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
