package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// ⚙️ Redis config
var rdb = redis.NewClient(&redis.Options{
	Addr: "localhost:6379", // hoặc Redis container nếu dùng Docker
})

// Prometheus metrics
var (
	pendingGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redis_stream_pending_messages",
			Help: "Pending messages in Redis stream consumer group",
		},
		[]string{"stream", "group"},
	)

	consumerPending = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redis_stream_consumer_pending",
			Help: "Pending messages per consumer",
		},
		[]string{"stream", "group", "consumer"},
	)

	consumerIdle = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "redis_stream_consumer_idle_seconds",
			Help: "Idle time per consumer in seconds",
		},
		[]string{"stream", "group", "consumer"},
	)
)

func collectMetrics() {
	stream := "domain_queue"
	group := "mygroup"

	for {
		// XINFO GROUPS
		groups, err := rdb.XInfoGroups(ctx, stream).Result()
		if err == nil {
			for _, g := range groups {
				pendingGauge.WithLabelValues(stream, g.Name).Set(float64(g.Pending))
			}
		} else {
			log.Println("XInfoGroups error:", err)
		}

		// XINFO CONSUMERS
		consumers, err := rdb.XInfoConsumers(ctx, stream, group).Result()
		if err == nil {
			for _, c := range consumers {
				consumerPending.WithLabelValues(stream, group, c.Name).Set(float64(c.Pending))
				consumerIdle.WithLabelValues(stream, group, c.Name).Set(float64(c.Idle) / 1000.0) // ms → sec
			}
		} else {
			log.Println("XInfoConsumers error:", err)
		}

		time.Sleep(10 * time.Second)
	}
}

func main() {
	// Register metrics
	prometheus.MustRegister(pendingGauge)
	prometheus.MustRegister(consumerPending)
	prometheus.MustRegister(consumerIdle)

	// Start collection loop
	go collectMetrics()

	// Expose metrics
	http.Handle("/metrics", promhttp.Handler())
	log.Println("Serving Prometheus metrics at :2112/metrics")
	log.Fatal(http.ListenAndServe(":2112", nil))
}
