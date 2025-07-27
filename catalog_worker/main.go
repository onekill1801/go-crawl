package main

import (
	"catalog/queue"
	"log"
)

// var ctx = context.Background()

// var rdb = redis.NewClient(&redis.Options{
// 	Addr: "localhost:6379",
// })

func main() {
	q := queue.NewRedisQueue("localhost:6379")

	// rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	// ctx := context.Background()

	// for {
	// 	streams, err := rdb.XRead(ctx, &redis.XReadArgs{
	// 		Streams: []string{"domain_queue", "0"},
	// 		Block:   5 * time.Second,
	// 		Count:   10,
	// 	}).Result()

	// stream := "domain_queue"
	// group := "mygroup"
	// newConsumer := "worker-1"

	// err := rdb.XGroupCreateMkStream(ctx, "domain_queue", "mygroup", "0").Err()
	// if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
	// 	log.Fatalf("Failed to create group: %v", err)
	// }

	// checkPendingMessages(stream, group)
	// claimPendingMessages(stream, group, newConsumer)

	// streams, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
	// 	Group:    "mygroup",
	// 	Consumer: "worker-1",
	// 	Streams:  []string{"domain_queue", ">"},
	// 	Block:    5 * time.Second,
	// }).Result()

	// if err != nil && err != redis.Nil {
	// 	log.Println("Redis read error:", err)
	// 	continue
	// }

	// for _, stream := range streams {
	// 	fmt.Println(stream)
	// for _, msg := range stream.Messages {

	// fmt.Println(msg)
	// var job DomainJob
	// json.Unmarshal([]byte(msg.Values["data"].(string)), &job)
	// if err := ProcessDomainJob(job, q); err != nil {
	// 	log.Println("DomainJob error:", err)
	// }
	// }
	// }

	// Lặp tương tự cho series_queue (gộp vào vòng sau)
	// }
	// TODO get list novel
	// var job DomainJob
	// job.Domain = "www.webtoons.com/en/"
	// if err := ProcessDomainJob(job, q); err != nil {
	// 	log.Println("DomainJob error:", err)
	// }

	// var wg sync.WaitGroup
	// sem := make(chan struct{}, 2) // Giới hạn 10 goroutine đồng thời

	// for i := 1; i <= 10; i++ {
	// 	i := i
	// 	wg.Add(1)
	// 	sem <- struct{}{} // chặn nếu vượt quá 10

	// 	go func() {
	// 		defer wg.Done()
	// 		defer func() { <-sem }()

	// 		var seriesJob SeriesJob
	// 		seriesJob.SeriesURL = fmt.Sprintf("https://www.webtoons.com/en/drama/lookism/list?title_no=1049&page=%d", i)

	// 		if err := ProcessSeriesJob(seriesJob, q); err != nil {
	// 			log.Println("DomainJob error on page", i, ":", err)
	// 		}
	// 	}()
	// }

	// wg.Wait()

	// TODO get list chappter
	// const totalPages = 56
	// const batchSize = 2

	// var wg sync.WaitGroup

	// for i := 1; i <= totalPages; i += batchSize {
	// 	// Xử lý batch 4 trang một lúc
	// 	for j := 0; j < batchSize && (i+j) <= totalPages; j++ {
	// 		wg.Add(1)
	// 		go func(page int) {
	// 			defer wg.Done()

	// 			var seriesJob SeriesJob
	// 			seriesJob.SeriesURL = fmt.Sprintf("https://www.webtoons.com/en/drama/lookism/list?title_no=1049&page=%d", page)

	// 			fmt.Println("seriesJob.SeriesURL:", seriesJob.SeriesURL)
	// 			// seriesJob.SeriesURL = "https://www.webtoons.com/en/drama/lookism/list?title_no=1049&page=3"
	// 			err := ProcessSeriesJob(seriesJob, q)
	// 			if err != nil {
	// 				log.Println("DomainJob error on page", page, ":", err)
	// 			}
	// 		}(i + j)
	// 	}

	// 	// Chờ 4 request chạy xong
	// 	wg.Wait()

	// 	// Nghỉ 300ms trước khi xử lý batch tiếp theo
	// 	time.Sleep(2000 * time.Millisecond)
	// }

	// TODO get list images in chapter
	var seriesJob SeriesJob
	seriesJob.SeriesURL = "https://www.webtoons.com/en/drama/lookism/ep-554-cheonmyeong-3/viewer?title_no=1049&episode_no=554"
	if err := ProcessImagesJob(seriesJob, q); err != nil {
		log.Println("DomainJob error on page", ":", err)
	}

}

// func checkPendingMessages(stream, group string) {
// 	result, err := rdb.XPending(ctx, stream, group).Result()
// 	if err != nil {
// 		log.Fatalf("XPENDING error: %v", err)
// 	}

// 	fmt.Printf("Total pending: %d\n", result.Count)
// 	fmt.Printf("First ID: %s\n", result.Lower)
// 	fmt.Printf("Last ID: %s\n", result.Higher)
// 	fmt.Printf("Consumers: %+v\n", result.Consumers)
// }

// func claimPendingMessages(stream, group, newConsumer string) {
// 	// Lấy các message ID đang pending
// 	entries, err := rdb.XPendingExt(ctx, &redis.XPendingExtArgs{
// 		Stream: stream,
// 		Group:  group,
// 		Start:  "-",
// 		End:    "+",
// 		Count:  10,
// 	}).Result()
// 	if err != nil {
// 		log.Fatalf("XPENDING EXT error: %v", err)
// 	}

// 	var idsToClaim []string
// 	for _, entry := range entries {
// 		if entry.Idle >= 30*time.Second {
// 			idsToClaim = append(idsToClaim, entry.ID)
// 		}
// 	}

// 	if len(idsToClaim) == 0 {
// 		fmt.Println("⏳ No messages need claiming")
// 		return
// 	}

// 	// Claim lại các message
// 	claimed, err := rdb.XClaim(ctx, &redis.XClaimArgs{
// 		Stream:   stream,
// 		Group:    group,
// 		Consumer: newConsumer,
// 		MinIdle:  30 * time.Second,
// 		Messages: idsToClaim,
// 	}).Result()
// 	if err != nil {
// 		log.Fatalf("XCLAIM error: %v", err)
// 	}

// 	for _, msg := range claimed {
// 		fmt.Printf("✅ Claimed message: ID=%s, Values=%v\n", msg.ID, msg.Values)
// 	}
// }
