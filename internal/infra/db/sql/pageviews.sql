-- name: StorePageview :exec
INSERT INTO pageviews (pathname, hashed_ip, country, browser, os, device, user_agent, referrer)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);