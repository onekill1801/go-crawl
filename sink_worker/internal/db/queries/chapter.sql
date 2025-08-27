-- name: CreateChapter :exec
INSERT INTO chapter (story_id, title, content, order_stt, image_url)
VALUES (?, ?, ?, ?, ?);

-- name: GetChapter :one
SELECT * FROM chapter WHERE id = ? LIMIT 1;

-- name: ListChaptersByStory :many
SELECT * FROM chapter
WHERE story_id = ?
ORDER BY order_stt ASC, created_at ASC
LIMIT ? OFFSET ?;

-- name: UpdateChapter :exec
UPDATE chapter
SET title = ?, content = ?, order_stt = ?, image_url = ?
WHERE id = ?;

-- name: DeleteChapter :exec
DELETE FROM chapter WHERE id = ?;
