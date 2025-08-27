package main

import (
	"catalog/queue"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func storyRunning(ctx context.Context, q *queue.RedisQueue, msg redis.XMessage) {
	if storyURL, ok := msg.Values["domain_url"].(string); ok {
		var job DomainJob
		job.Domain = storyURL
		if err := ProcessDomainJob(job, q); err != nil {
			log.Println("DomainJob error:", err)
		}
	}

}

func chapterRunning(ctx context.Context, q *queue.RedisQueue, msg redis.XMessage) {
	series := msg.Values["series_url"].(string)

	page := 1
	for {
		var seriesJob SeriesJob
		seriesJob.SeriesURL = fmt.Sprintf("%s&page=%d", series, page)

		fmt.Println("seriesJob.SeriesURL:", seriesJob.SeriesURL)

		maxPage, err := ProcessSeriesJob(seriesJob, q)
		if err != nil {
			log.Println("DomainJob error on page", page, ":", err)
			break
		}

		// nếu maxPage >= page thì dừng luôn
		if maxPage >= page {
			fmt.Println("Max page found:", maxPage, " -> stop processing")
			break
		}

		page++

		// nghỉ 2s trước khi request tiếp
		time.Sleep(1 * time.Second)
	}
}

func imageRunning(ctx context.Context, q *queue.RedisQueue, msg redis.XMessage) {
	chapterUrl := msg.Values["chapter_url"].(string)
	var seriesJob SeriesJob
	seriesJob.SeriesURL = chapterUrl
	if err := ProcessImagesJob(seriesJob, q); err != nil {
		log.Println("DomainJob error on page", ":", err)
	}
}
