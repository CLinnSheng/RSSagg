package main

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/CLinnSheng/RSSagg/internal/database"
	"github.com/google/uuid"
)

// it will run at the background concurrently
func startScrapping(db *database.Queries, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Scrapping on %v go routines every %s duration", concurrency, timeBetweenRequest)

	ticker := time.NewTicker(timeBetweenRequest)
	// .C stands for channel
	// every timeBetweenRequest a value will be send across the channel
	// empty initialization is because we want it to straight fire when the server run
	// so every interval will get the next feed to fetch
	for ; ; <-ticker.C {
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrency))

		if err != nil {
			log.Println("Error fetching feeds: ", err)
			// because this function will be keep running as long as the server is running
			continue
		}

		// waitgroup use to synchronize the execution of multiple go routines or "thread" to
		// sneusre the main thread waits for all the other thread to complete before exiting
		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)

			// create new thread
			go scrapeFeed(db, wg, feed)
		}

		// make sure all the goroutines/thread has return or done
		wg.Wait()
	}
}

func scrapeFeed(db *database.Queries, wg *sync.WaitGroup, feed database.Feed) {
	// wg.Done when the thread finish it works
	// it actually decrement the wg by 1
	defer wg.Done()

	_, err := db.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Println("Error marking feed as fetched: ", err)
		return
	}

	rssFeed, err := urlToFeed(feed.Url)
	if err != nil {
		log.Println("Error fetching feed: ", err)
		return
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}

		t, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Couldnt parse date %v, with err %v", item.PubDate, err)
			continue
		}

		_, err = db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: t,
			Url:         item.Link,
			FeedID:      feed.ID,
		})

		if err != nil {
			// prevent logging for duplicate post
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Println("Failed to Create Post")
		}
	}

	log.Printf("Feed %s collected, %v posts found", feed.Name, len(rssFeed.Channel.Item))
}
