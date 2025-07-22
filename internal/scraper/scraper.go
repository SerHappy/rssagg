package scraper

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/serhappy/rssagg/internal/app"
	"github.com/serhappy/rssagg/internal/db"
)

func StartScraping(app *app.App, concurrency int, timeBetweenRequest time.Duration) {
	log.Printf("Collecting feeds every %s on %v goroutines...", timeBetweenRequest, concurrency)
	ticker := time.NewTicker(timeBetweenRequest)

	for ; ; <-ticker.C {
		feeds, err := app.DB.GetNextFeedsToFetch(context.Background(), int32(concurrency))
		if err != nil {
			log.Println("Couldn't get next feeds to fetch", err)
			continue
		}
		log.Printf("Found %v feeds to fetch!", len(feeds))

		wg := &sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go scrapeFeed(app, wg, feed)
		}
		wg.Wait()
	}
}

func scrapeFeed(app *app.App, wg *sync.WaitGroup, feed db.Feed) {
	defer wg.Done()
	_, err := app.DB.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Couldn't mark feed %s fetched: %v", feed.Name, err)
		return
	}

	feedData, err := fetchFeed(feed.Url)
	if err != nil {
		log.Printf("Couldn't collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		description := sql.NullString{}
		if item.Description != "" {
			description = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		}
		pubAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Couldn't parse date %s: %v", item.PubDate, err)
			continue
		}
		_, err = app.DB.CreatePost(context.Background(), db.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Description: description,
			PublishedAt: pubAt,
			Url:         item.Link,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
			log.Printf("Couldn't create post: %v", err)
			continue
		}
	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}
