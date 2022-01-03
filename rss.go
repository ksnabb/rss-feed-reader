package main

import (
	"encoding/xml"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
)

type Item struct {
	Title       string `xml:"title" json:"title"`
	Description string `xml:"description" json:"desciption"`
	Link        string `xml:"link" json:"link"`
	PubDate     string `xml:"pubDate" json:"pubDate"`
}

type Feed struct {
	Title string `xml:"channel>title" json:"title"`
	Items []Item `xml:"channel>item" json:"items"`
}

func parseFeedFromXML(rssXml []byte) (Feed, error) {
	var feed Feed
	if err := xml.Unmarshal(rssXml, &feed); err != nil {
		return feed, errors.New("could not parse xml")
	}
	return feed, nil
}

func parseFeedFromUrl(rssFeedUrl *url.URL) (Feed, error) {
	resp, err := http.Get(rssFeedUrl.String())
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 200 {
		return Feed{}, errors.New("url did not return statusCode 200")
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return parseFeedFromXML(content)
}
