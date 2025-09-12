-- name: StorePageview :exec
INSERT INTO pageviews (pathname, ip_address, user_agent, referrer)
VALUES ($1, $2, $3, $4);