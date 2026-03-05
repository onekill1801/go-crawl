// Seed domain_queue để bắt đầu crawl.
// Chạy: cd catalog_worker && go run ./cmd/seed [domain_url]
// Mặc định: https://www.webtoons.com/en/
package main

import (
	"fmt"
	"log"
	"os"

	"catalog/queue"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	domainURL := "https://www.webtoons.com/en/"
	if len(os.Args) > 1 {
		domainURL = os.Args[1]
	}

	q := queue.NewRedisQueue(redisAddr)
	err := q.Push("domain_queue", struct {
		DomainURL string `json:"domain_url"`
	}{DomainURL: domainURL})
	if err != nil {
		log.Fatalf("Push domain_queue: %v", err)
	}
	fmt.Printf("✅ Đã gửi domain_url=%s vào domain_queue (Redis %s)\n", domainURL, redisAddr)
}
