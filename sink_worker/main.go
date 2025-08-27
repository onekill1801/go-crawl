package main

import (
	"context"
	"fmt"
	"log"

	dbgen "github.com/chungtv/sink_worker/internal/db" // sqlc generated

	"github.com/chungtv/sink_worker/internal/redis"
	"github.com/chungtv/sink_worker/internal/service"
)

func main() {
	ctx := context.Background()

	// Kết nối MySQL
	database, err := dbgen.NewDB("root:your_root_password@tcp(192.168.1.6:5306)/test?parseTime=true")
	if err != nil {
		log.Fatal("cannot connect db:", err)
	}
	q := dbgen.New(database)

	// Kết nối Redis
	rdb := redis.NewClient("localhost:6379")

	processor := service.NewProcessor(q, rdb)

	stream := "events"
	group := "worker-group"
	consumer := "worker-1"

	for {
		msgs, err := redis.ReadStream(ctx, rdb, stream, group, consumer)
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		for _, m := range msgs {
			if err := processor.HandleMessage(ctx, m); err != nil {
				log.Println("process error:", err)
			} else {
				fmt.Println("Inserted:", m.ID)
			}
		}
	}
}
