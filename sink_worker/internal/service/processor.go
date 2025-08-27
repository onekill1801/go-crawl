package service

import (
	"context"
	"encoding/json"

	"github.com/chungtv/sink_worker/internal/db"

	"github.com/redis/go-redis/v9"
)

// Processor xử lý message từ Redis và insert vào DB
type Processor struct {
	q   *db.Queries
	rdb *redis.Client
}

func NewProcessor(q *db.Queries, rdb *redis.Client) *Processor {
	return &Processor{q: q, rdb: rdb}
}

func (p *Processor) HandleMessage(ctx context.Context, msg redis.XMessage) error {
	payload, err := json.Marshal(msg.Values)
	if err != nil {
		return err
	}

	_, err = p.q.InsertEvent(ctx, db.InsertEventParams{
		StreamID: msg.ID,
		Payload:  payload,
	})
	return err
}
