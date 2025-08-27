package main

import (
	"catalog/queue"
	"context"
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
		{"stories_queue", "worker-group", "worker-2", storyRunning},
		{"chapters_queue", "worker-group", "worker-2", chapterRunning},
		{"images_queue", "worker-group", "worker-2", imageRunning},
	}
	for _, s := range streams {
		wg.Add(1)
		go func(streamName string, handler func(ctx context.Context, q *queue.RedisQueue, msg redis.XMessage)) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					log.Printf("[%s] stopping goroutine\n", streamName)
					return
				default:
					msgs, err := queue.ReadStream(ctx, s.Name, s.Group, s.Consumer)
					if err != nil {
						log.Println("read error:", err)
						time.Sleep(time.Second) // retry
						continue
					}

					for _, msg := range msgs {
						s.Handler(ctx, q, msg) // gọi trực tiếp
					}
				}
			}
		}(s.Name, s.Handler)
	}

	// Chờ tất cả goroutine kết thúc khi cancel được gọi
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
