-- name: StorePageview :exec
INSERT INTO pageviews (pathname, hashed_ip, country, browser, os, device, user_agent, referrer)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetTotalViewsAndVisitors :one
SELECT COUNT(id) as total_views, COUNT(DISTINCT hashed_ip) as total_visitors
FROM pageviews;