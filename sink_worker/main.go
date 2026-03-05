package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	dbgen "github.com/chungtv/sink_worker/internal/db" // sqlc generated

	"github.com/chungtv/sink_worker/internal/redis"
	"github.com/chungtv/sink_worker/internal/service"
	redislib "github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		dsn = "root:crawl_secret@tcp(localhost:3306)/story?parseTime=true"
	}
	// Migrate cần multiStatements=true
	migrateDSN := dsn
	if strings.Contains(dsn, "?") {
		migrateDSN = dsn + "&multiStatements=true"
	} else {
		migrateDSN = dsn + "?multiStatements=true"
	}
	migrate1(migrateDSN)

	// Kết nối MySQL (đảm bảo có parseTime)
	if !strings.Contains(dsn, "parseTime") {
		if strings.Contains(dsn, "?") {
			dsn += "&parseTime=true"
		} else {
			dsn += "?parseTime=true"
		}
	}
	database, err := dbgen.NewDB(dsn)
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}
	q := dbgen.New(database)

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	rdb := redis.NewClient(redisAddr)

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
