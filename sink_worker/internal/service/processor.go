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

	_, _ = p.q.InsertEvent(ctx, db.InsertEventParams{
		StreamID: msg.ID,
		Payload:  payload,
	})

	// dataStr, _ := msg.Values["data"].(string)
	// if err := json.Unmarshal([]byte(dataStr), &parsed); err == nil {
	// 	items = append(items, parsed)
	// }

	err = p.q.CreateChapter(ctx, db.CreateChapterParams{
		StoryID: msg.ID,
		Title:   msg.ID,
	})
	return err
}

func (p *Processor) HandleMessageStories(ctx context.Context, msg redis.XMessage) error {

	_, _ = msg.Values["data"].(string)
	// if err := json.Unmarshal([]byte(dataStr), &parsed); err == nil {
	// 	items = append(items, parsed)
	// }

	// err = p.q.CreateChapter(ctx, db.CreateChapterParams{
	// 	StoryID: msg.ID,
	// 	Title:   msg.ID,
	// })
	return nil
}

func (p *Processor) HandleMessageChapter(ctx context.Context, msg redis.XMessage) error {

	_, _ = msg.Values["data"].(string)
	// if err := json.Unmarshal([]byte(dataStr), &parsed); err == nil {
	// 	items = append(items, parsed)
	// }

	// err = p.q.CreateChapter(ctx, db.CreateChapterParams{
	// 	StoryID: msg.ID,
	// 	Title:   msg.ID,
	// })
	return nil
}

func (p *Processor) HandleMessageImages(ctx context.Context, msg redis.XMessage) error {

	_, _ = msg.Values["data"].(string)
	// if err := json.Unmarshal([]byte(dataStr), &parsed); err == nil {
	// 	items = append(items, parsed)
	// }

	// err = p.q.CreateChapter(ctx, db.CreateChapterParams{
	// 	StoryID: msg.ID,
	// 	Title:   msg.ID,
	// })
	return nil
}
