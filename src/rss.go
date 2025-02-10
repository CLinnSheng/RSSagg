package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"time"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

// using xml tag so that when it parse from struct to xml it will follow all the xml tag defined below or else
// by default it just use the field name
type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func urlToFeed(url string) (RSSFeed, error) {
	httpClient := http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := httpClient.Get(url)
	if err != nil {
		return RSSFeed{}, err
	}

	// defer will run at last after exiting the function
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return RSSFeed{}, err
	}

	rssFeed := RSSFeed{}
	err = xml.Unmarshal(data, &rssFeed)
	if err != nil {
		return RSSFeed{}, err
	}

	return rssFeed, nil
}
