package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	dbgen "github.com/chungtv/sink_worker/internal/db" // sqlc generated

	"github.com/chungtv/sink_worker/internal/redis"
	"github.com/chungtv/sink_worker/internal/service"
	redislib "github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	dns := "root:your_root_password@tcp(192.168.1.6:5306)"
	migrate1(dns)

	// Kết nối MySQL
	database, err := dbgen.NewDB(dns + "?parseTime=true")
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}
	q := dbgen.New(database)

	// Kết nối Redis
	rdb := redis.NewClient("localhost:6379")

	processor := service.NewProcessor(q, rdb)

	streams := []struct {
		Name     string
		Group    string
		Consumer string
		Handler  func(ctx context.Context, msg redislib.XMessage) error
	}{
		// {"events", "worker-group", "worker-1", processor.HandleMessage},
		{"series_queue", "worker-sink", "worker-2", processor.HandleMessageStories},
		{"chapter_queue", "worker-sink", "worker-3", processor.HandleMessageChapter},
		{"images_queue", "worker-sink", "worker-4", processor.HandleMessageImages},
	}

	var wg sync.WaitGroup

	// Tạo 4 goroutine để xử lý song song 4 queue
	for _, s := range streams {
		if err := redis.EnsureConsumerGroup(rdb, s.Name, s.Group, false); err != nil {
			log.Fatalf("cannot create group for stream %s: %v", s.Name, err)
		} else {
			fmt.Printf("Consumer group ensured for stream %s and group %s\n", s.Name, s.Group)
		}
		wg.Add(1)
		go func(s struct {
			Name     string
			Group    string
			Consumer string
			Handler  func(ctx context.Context, msg redislib.XMessage) error
		}) {
			defer wg.Done()
			for {
				msgs, err := redis.ReadStream(ctx, rdb, s.Name, s.Group, s.Consumer)
				if err != nil {
					log.Println("read error:", err)
					continue
				}

				for _, m := range msgs {
					if err := s.Handler(ctx, m); err != nil {
						fmt.Printf("Inserted error: %s", s.Name)
						log.Println("process error:", err)
					} else {
						fmt.Printf("[%s] Inserted: %s\n", s.Name, m.ID)
					}
				}
			}
		}(s)
	}

	wg.Wait()
}
