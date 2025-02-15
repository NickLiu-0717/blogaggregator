package main

import (
	"context"
	"database/sql"
	"html"
	"log"

	"time"

	"github.com/NickLiu-0717/blogaggregator/internal/database"
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
		var pubtime time.Time
		if item.PubDate != "" {
			pubtime, err = time.Parse(time.RFC1123Z, item.PubDate)
			if err != nil {
				log.Printf("Error fetching feed: %s", err)
				return
			}
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			Title: sql.NullString{
				String: html.UnescapeString(item.Title),
				Valid:  true,
			},
			Url:         item.Link,
			Description: html.UnescapeString(item.Description),
			PublishedAt: pubtime,
			FeedID:      feed.ID,
		})
		if err != nil {
			log.Printf("Error creating post: %s", err)
			return
		}
	}
}
