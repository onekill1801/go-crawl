package main

import (
	"catalog/queue"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

func storyRunning(ctx context.Context, q *queue.RedisQueue, msg redis.XMessage) {
	var job DomainJob
	job.Domain = "www.webtoons.com/en/"
	if err := ProcessDomainJob(job, q); err != nil {
		log.Println("DomainJob error:", err)
	}
}

func chapterRunning(ctx context.Context, q *queue.RedisQueue, msg redis.XMessage) {
	const totalPages = 56
	const batchSize = 2

	var wg sync.WaitGroup

	for i := 1; i <= totalPages; i += batchSize {
		// Xử lý batch 4 trang một lúc
		for j := 0; j < batchSize && (i+j) <= totalPages; j++ {
			wg.Add(1)
			go func(page int) {
				defer wg.Done()

				var seriesJob SeriesJob
				seriesJob.SeriesURL = fmt.Sprintf("https://www.webtoons.com/en/drama/lookism/list?title_no=1049&page=%d", page)

				fmt.Println("seriesJob.SeriesURL:", seriesJob.SeriesURL)
				// seriesJob.SeriesURL = "https://www.webtoons.com/en/drama/lookism/list?title_no=1049&page=3"
				err := ProcessSeriesJob(seriesJob, q)
				if err != nil {
					log.Println("DomainJob error on page", page, ":", err)
				}
			}(i + j)
		}

		// Chờ 4 request chạy xong
		wg.Wait()

		// Nghỉ 300ms trước khi xử lý batch tiếp theo
		time.Sleep(2000 * time.Millisecond)
	}
}

func imageRunning(ctx context.Context, q *queue.RedisQueue, msg redis.XMessage) {
	var seriesJob SeriesJob
	seriesJob.SeriesURL = "https://www.webtoons.com/en/drama/lookism/ep-554-cheonmyeong-3/viewer?title_no=1049&episode_no=554"
	if err := ProcessImagesJob(seriesJob, q); err != nil {
		log.Println("DomainJob error on page", ":", err)
	}
}
