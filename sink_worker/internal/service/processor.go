package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

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
	var params db.CreateStoryParams

	// Parse StoryID
	if storyURL, ok := msg.Values["series_url"].(string); ok {
		id, err := extractTitleNo(storyURL)
		if err != nil {
			return fmt.Errorf("cannot extract title_no: %w", err)
		}
		params.ID = strconv.FormatInt(id, 10) // chuyển int64 -> string
	}

	// Parse Url
	if url, ok := msg.Values["series_url"].(string); ok {
		params.CoverUrl = sql.NullString{String: url, Valid: url != ""}
	}

	// Insert vào DB (sqlc generated)
	if err := p.q.CreateStory(ctx, params); err != nil {
		return fmt.Errorf("insert story error: %w", err)
	}

	return nil
}

func extractTitleNo(url string) (int64, error) {
	re := regexp.MustCompile(`title_no=(\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return 0, fmt.Errorf("title_no not found in url: %s", url)
	}
	var num int64
	fmt.Sscanf(matches[1], "%d", &num)
	return num, nil
}

func (p *Processor) HandleMessageChapter(ctx context.Context, msg redis.XMessage) error {
	var params db.CreateChapterParams

	// Parse StoryID
	if storyID, ok := msg.Values["story_id"].(string); ok {
		params.StoryID = storyID
	}

	// Parse Url
	if url, ok := msg.Values["chapter_url"].(string); ok {
		params.Content = sql.NullString{String: url, Valid: url != ""}
	}

	// Parse Title
	if title, ok := msg.Values["title"].(string); ok {
		params.Title = title
	}

	if title, ok := msg.Values["title"].(string); ok {
		if num, err := extractChapterNumber(title); err == nil {
			params.OrderStt = sql.NullInt32{
				Int32: int32(num),
				Valid: true,
			}
		} else {
			params.OrderStt = sql.NullInt32{Valid: false} // NULL
		}
	}

	// Insert vào DB (sqlc generated)
	if err := p.q.CreateChapter(ctx, params); err != nil {
		return fmt.Errorf("insert chapter error: %w", err)
	}

	return nil
}

func extractChapterNumber(title string) (int32, error) {
	re := regexp.MustCompile(`(?i)^Ep\.?\s*(\d+)`)
	matches := re.FindStringSubmatch(title)
	if len(matches) < 2 {
		return 0, fmt.Errorf("no chapter number found in: %s", title)
	}

	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}
	return int32(num), nil
}

func (p *Processor) HandleMessageImages(ctx context.Context, msg redis.XMessage) error {
	var params db.CreateImageParams

	// Parse ChapterID
	if chapStr, ok := msg.Values["chapter_id"].(string); ok {
		if chapID, err := strconv.ParseInt(chapStr, 10, 64); err == nil {
			params.ChapterID = chapID
		} else {
			return fmt.Errorf("invalid chapter_id: %v", chapStr)
		}
	}

	// Parse StoryID
	if storyID, ok := msg.Values["story_id"].(string); ok {
		params.StoryID = storyID
	}

	// Parse Url
	if url, ok := msg.Values["image_url"].(string); ok {
		params.Url = url
	}

	// Parse Referer
	if referer, ok := msg.Values["referer"].(string); ok {
		params.Referer = sql.NullString{String: referer, Valid: referer != ""}
	}

	// Parse Title
	if title, ok := msg.Values["title"].(string); ok {
		params.Title = sql.NullString{String: title, Valid: title != ""}
	}

	// Parse OrderStt
	if orderStr, ok := msg.Values["order"].(string); ok {
		if orderInt, err := strconv.ParseUint(orderStr, 10, 32); err == nil {
			params.OrderStt = uint32(orderInt)
		}
	}

	// Insert vào DB (sqlc generated)
	if err := p.q.CreateImage(ctx, params); err != nil {
		return fmt.Errorf("insert image error: %w", err)
	}

	return nil
}
