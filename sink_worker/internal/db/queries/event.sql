-- name: InsertEvent :execresult
INSERT INTO events (stream_id, payload) 
VALUES (?, ?);
