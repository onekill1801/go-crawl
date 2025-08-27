-- name: CreateImage :exec
INSERT INTO images (chapter_id, url, referer, title, order_stt)
VALUES (?, ?, ?, ?, ?);

-- name: GetImage :one
SELECT * FROM images WHERE id = ? LIMIT 1;

-- name: ListImagesByChapter :many
SELECT * FROM images
WHERE chapter_id = ?
ORDER BY order_stt ASC, created_at ASC
LIMIT ? OFFSET ?;

-- name: UpdateImage :exec
UPDATE images
SET url = ?, referer = ?, title = ?, order_stt = ?
WHERE id = ?;

-- name: DeleteImage :exec
DELETE FROM images WHERE id = ?;
