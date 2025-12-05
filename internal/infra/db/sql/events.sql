-- name: StoreEvent :exec
INSERT INTO events (cheatsheet_id, event_type, pathname, hashed_ip)
VALUES ($1, $2, $3, $4);

-- name: GetTotalClicksAndDownloads :one
SELECT
COALESCE(COUNT(hashed_ip) FILTER (WHERE event_type = 'click'), 0)::bigint AS total_clicks,
COALESCE(COUNT(hashed_ip) FILTER (WHERE event_type = 'download'), 0)::bigint AS total_downloads
FROM events;