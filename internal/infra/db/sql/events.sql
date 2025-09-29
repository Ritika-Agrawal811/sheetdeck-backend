-- name: StoreEvent :exec
INSERT INTO events (cheatsheet_id, event_type, pathname, hashed_ip)
VALUES ($1, $2, $3, $4);