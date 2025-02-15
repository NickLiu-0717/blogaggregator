package main

import (
	"context"
	"fmt"
	"html"
	"log"

	"github.com/google/uuid"
)

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Printf("Error getting next feed: %s", err)
		return
	}
	if feed.ID == uuid.Nil && feed.Url == "" {
		log.Println("No feed to fetch")
		return
	}
	if err = s.db.MarkFeedFetched(context.Background(), feed.ID); err != nil {
		log.Printf("Error marking feed fetch: %s", err)
		return
	}
	rss, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("Error fetching feed: %s", err)
		return
	}
	for _, item := range rss.Channel.Item {
		if html.UnescapeString(item.Title) == "" {
			fmt.Printf("Title: emtpy string with feed id: %s\n", feed.ID)
			continue
		}
		fmt.Printf("Title: %+v\n", html.UnescapeString(item.Title))
	}
}
