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

var ctx = context.Background()

func main() {
	q := queue.NewRedisQueue("localhost:6379")

	var wg sync.WaitGroup
	streams := []struct {
		Name     string
		Group    string
		Consumer string
		Handler  func(ctx context.Context, q *queue.RedisQueue, msg redis.XMessage)
	}{
		{"domain_queue", "worker-group", "worker-1", storyRunning},
		{"series_queue", "worker-group", "worker-2", chapterRunning},
		{"chapter_queue", "worker-group", "worker-3", imageRunning},
	}

	for _, s := range streams {
		if err := q.EnsureConsumerGroup(s.Name, s.Group, true); err != nil {
			log.Fatalf("cannot create group for stream %s: %v", s.Name, err)
		} else {
			fmt.Printf("Consumer group ensured for stream %s and group %s\n", s.Name, s.Group)
		}
		wg.Add(1)
		go func(streamName, group, consumer string, handler func(ctx context.Context, q *queue.RedisQueue, msg redis.XMessage)) {
			defer wg.Done()
			for {
				msgs, err := q.ReadStream(ctx, streamName, group, consumer)
				if err != nil {
					log.Println("read error:", err)
					time.Sleep(time.Second)
					continue
				} else {
					fmt.Printf("Read %d messages from %s\n", len(msgs), streamName)
				}
				for _, msg := range msgs {
					handler(ctx, q, msg)
				}
			}
		}(s.Name, s.Group, s.Consumer, s.Handler)
	}

	wg.Wait()
	log.Println("All streams stopped gracefully")
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
