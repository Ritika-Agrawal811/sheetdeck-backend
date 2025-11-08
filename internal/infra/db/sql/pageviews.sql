-- name: StorePageview :exec
INSERT INTO pageviews (pathname, hashed_ip, country, browser, os, device, user_agent, referrer)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetTotalViewsAndVisitors :one
SELECT COUNT(id) as total_views, COUNT(DISTINCT hashed_ip) as total_visitors
FROM pageviews;

-- name: GetPageviewTimeseriesByDay :many
SELECT 
    d::date AS date,
    COALESCE(COUNT(p.viewed_at), 0)::bigint AS views,
    COALESCE(COUNT(DISTINCT p.hashed_ip), 0)::bigint AS unique_visitors
FROM generate_series(
    (NOW() - make_interval(days => sqlc.arg(days)::int))::date,
    NOW()::date,
    '1 day'
) AS d
LEFT JOIN pageviews p
    ON DATE_TRUNC('day', p.viewed_at)::date = d
    AND p.browser != 'Headless Chrome'
GROUP BY d
ORDER BY d;

-- name: GetPageviewTimeseriesForLast24Hours :many
SELECT 
    h::timestamp AS hour,
    COALESCE(COUNT(p.viewed_at), 0)::bigint AS views,
    COALESCE(COUNT(DISTINCT p.hashed_ip), 0)::bigint AS unique_visitors
FROM generate_series(
    date_trunc('hour', NOW() - INTERVAL '23 hours'),
    date_trunc('hour', NOW()),
    '1 hour'
) AS h
LEFT JOIN pageviews p
    ON date_trunc('hour', p.viewed_at) = h
    AND p.browser != 'Headless Chrome'
GROUP BY h
ORDER BY h;

-- name: GetDevicesSummaryByDay :many
SELECT 
  d::date AS date,
  -- Mobile
  COALESCE(COUNT(p.viewed_at) FILTER (WHERE p.device = 'mobile'), 0)::bigint AS mobile_views,
  COALESCE(COUNT(DISTINCT p.hashed_ip) FILTER (WHERE p.device = 'mobile'), 0)::bigint AS mobile_visitors,
  -- Desktop
  COALESCE(COUNT(p.viewed_at) FILTER (WHERE p.device = 'desktop'), 0)::bigint AS desktop_views,
  COALESCE(COUNT(DISTINCT p.hashed_ip) FILTER (WHERE p.device = 'desktop'), 0)::bigint AS desktop_visitors
FROM generate_series(
    (NOW() - make_interval(days => sqlc.arg(days)::int))::date,
    NOW()::date,
    '1 day'
) AS d
LEFT JOIN pageviews p
  ON DATE_TRUNC('day', p.viewed_at)::date = d
  AND p.browser != 'Headless Chrome'
GROUP BY d
ORDER BY d;

-- name: GetDevicesSummaryForLast24Hours :many
SELECT 
  h::timestamp AS hour,
  -- Mobile
  COALESCE(COUNT(p.viewed_at) FILTER (WHERE p.device = 'mobile'), 0)::bigint AS mobile_views,
  COALESCE(COUNT(DISTINCT p.hashed_ip) FILTER (WHERE p.device = 'mobile'), 0)::bigint AS mobile_visitors,
  -- Desktop
  COALESCE(COUNT(p.viewed_at) FILTER (WHERE p.device = 'desktop'), 0)::bigint AS desktop_views,
  COALESCE(COUNT(DISTINCT p.hashed_ip) FILTER (WHERE p.device = 'desktop'), 0)::bigint AS desktop_visitors
FROM generate_series(
     date_trunc('hour', NOW() - INTERVAL '23 hours'),
    date_trunc('hour', NOW()),
    '1 hour'
) AS h
LEFT JOIN pageviews p
  ON DATE_TRUNC('hour', p.viewed_at)::timestamp = h
  AND p.browser != 'Headless Chrome'
GROUP BY h
ORDER BY h;

-- name: GetBrowsersSummaryByDay :many
SELECT 
   DISTINCT(browser), 
   COALESCE(COUNT(viewed_at), 0)::bigint AS views,
   COALESCE(COUNT(DISTINCT hashed_ip), 0)::bigint AS unique_visitors
FROM pageviews 
WHERE browser != 'Headless Chrome'
  AND DATE_TRUNC('day', viewed_at)::date >= (NOW() - make_interval(days => sqlc.arg(days)::int))::date
  AND DATE_TRUNC('day', viewed_at)::date <= NOW()::date
GROUP BY browser;

-- name: GetBrowsersSummaryForLast24Hours :many
SELECT 
   DISTINCT(browser), 
   COALESCE(COUNT(viewed_at), 0)::bigint AS views,
   COALESCE(COUNT(DISTINCT hashed_ip), 0)::bigint AS unique_visitors
FROM pageviews 
WHERE browser != 'Headless Chrome'
  AND viewed_at >= NOW() - INTERVAL '23 hours'
GROUP BY browser;
