-- name: CreateStory :exec
INSERT INTO stories (id, title, author, cover_url, domain_link, image_url)
VALUES (?, ?, ?, ?, ?, ?);

-- name: GetStory :one
SELECT * FROM stories WHERE id = ? LIMIT 1;

-- name: ListStories :many
SELECT * FROM stories ORDER BY created_at DESC LIMIT ? OFFSET ?;

-- name: UpdateStory :exec
UPDATE stories
SET title = ?, author = ?, cover_url = ?, domain_link = ?, image_url = ?
WHERE id = ?;

-- name: DeleteStory :exec
DELETE FROM stories WHERE id = ?;

