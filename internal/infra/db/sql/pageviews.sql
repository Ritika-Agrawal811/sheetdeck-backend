-- name: StorePageview :exec
INSERT INTO pageviews (pathname, hashed_ip, country, browser, os, device, user_agent, referrer)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetTotalViewsAndVisitors :one
SELECT COUNT(id) as total_views, COUNT(DISTINCT hashed_ip) as total_visitors
FROM pageviews
WHERE browser NOT IN ('Headless Chrome', 'Google-Read-Aloud');

-- name: GetMetricsTimeseriesByDay :many
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
    AND p.browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
GROUP BY d
ORDER BY d;

-- name: GetMetricsTimeseriesForLast24Hours :many
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
    AND p.browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
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
  AND p.browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
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
  AND p.browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
GROUP BY h
ORDER BY h;

-- name: GetBrowsersSummaryByDay :many
SELECT 
   DISTINCT(browser), 
   COALESCE(COUNT(viewed_at), 0)::bigint AS views,
   COALESCE(COUNT(DISTINCT hashed_ip), 0)::bigint AS unique_visitors
FROM pageviews 
WHERE browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
  AND DATE_TRUNC('day', viewed_at)::date >= (NOW() - make_interval(days => sqlc.arg(days)::int))::date
  AND DATE_TRUNC('day', viewed_at)::date <= NOW()::date
GROUP BY browser
ORDER BY views DESC;

-- name: GetBrowsersSummaryForLast24Hours :many
SELECT 
   DISTINCT(browser), 
   COALESCE(COUNT(viewed_at), 0)::bigint AS views,
   COALESCE(COUNT(DISTINCT hashed_ip), 0)::bigint AS unique_visitors
FROM pageviews 
WHERE browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
  AND viewed_at >= NOW() - INTERVAL '23 hours'
GROUP BY browser
ORDER BY views DESC;

-- name: GetOSSummaryByDay :many
SELECT 
   CASE 
      -- iOS variants
      WHEN os ILIKE '%iPhone OS%' OR os ILIKE '%iOS%' THEN 'iOS'::varchar
      -- Android variants
      WHEN os ILIKE '%Android%' THEN 'Android'::varchar
      -- Windows variants
      WHEN os ILIKE '%Windows%' THEN 'Windows'::varchar
      -- Mac OS variants
      WHEN os ILIKE '%Mac OS%' AND os NOT ILIKE '%iPhone%' THEN 'Mac OS'::varchar
      -- Linux variants
      WHEN os ILIKE '%Linux%' AND os NOT ILIKE '%Android%' THEN 'Linux'::varchar
      -- Chrome OS
      WHEN os ILIKE '%Chrome OS%' THEN 'Chrome OS'::varchar
      -- Keep other OS as-is or mark as Other
      ELSE COALESCE(os, 'Other')::varchar
   END AS os_group,
   COALESCE(COUNT(viewed_at), 0)::bigint AS views,
   COALESCE(COUNT(DISTINCT hashed_ip), 0)::bigint AS unique_visitors
FROM pageviews 
WHERE browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
  AND DATE_TRUNC('day', viewed_at)::date >= (NOW() - make_interval(days => sqlc.arg(days)::int))::date
  AND DATE_TRUNC('day', viewed_at)::date <= NOW()::date
  AND os IS NOT NULL
GROUP BY os_group
ORDER BY views DESC;

-- name: GetOSSummaryForLast24Hours :many
SELECT 
   CASE 
      -- iOS variants
      WHEN os ILIKE '%iPhone OS%' OR os ILIKE '%iOS%' THEN 'iOS'::varchar
      -- Android variants
      WHEN os ILIKE '%Android%' THEN 'Android'::varchar
      -- Windows variants
      WHEN os ILIKE '%Windows%' THEN 'Windows'::varchar
      -- Mac OS variants
      WHEN os ILIKE '%Mac OS%' AND os NOT ILIKE '%iPhone%' THEN 'Mac OS'::varchar
      -- Linux variants
      WHEN os ILIKE '%Linux%' AND os NOT ILIKE '%Android%' THEN 'Linux'::varchar
      -- Chrome OS
      WHEN os ILIKE '%Chrome OS%' THEN 'Chrome OS'::varchar
      -- Keep other OS as-is or mark as Other
      ELSE COALESCE(os, 'Other')::varchar
   END AS os_group,
   COALESCE(COUNT(viewed_at), 0)::bigint AS views,
   COALESCE(COUNT(DISTINCT hashed_ip), 0)::bigint AS unique_visitors
FROM pageviews 
WHERE browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
  AND viewed_at >= NOW() - INTERVAL '23 hours'
  AND os IS NOT NULL
GROUP BY os_group
ORDER BY views DESC;

-- name: GetReferrerSummaryByDay :many
SELECT 
   DISTINCT(referrer), 
   COALESCE(COUNT(viewed_at), 0)::bigint AS views,
   COALESCE(COUNT(DISTINCT hashed_ip), 0)::bigint AS unique_visitors
FROM pageviews 
WHERE browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
  AND DATE_TRUNC('day', viewed_at)::date >= (NOW() - make_interval(days => sqlc.arg(days)::int))::date
  AND DATE_TRUNC('day', viewed_at)::date <= NOW()::date
  AND referrer IS NOT NULL
GROUP BY referrer
ORDER BY views DESC;

-- name: GetReferrerSummaryForLast24Hours :many
SELECT 
   DISTINCT(referrer), 
   COALESCE(COUNT(viewed_at), 0)::bigint AS views,
   COALESCE(COUNT(DISTINCT hashed_ip), 0)::bigint AS unique_visitors
FROM pageviews 
WHERE browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
 AND viewed_at >= NOW() - INTERVAL '23 hours'
 AND referrer IS NOT NULL
GROUP BY referrer
ORDER BY views DESC;

-- name: GetRoutesSummaryByDay :many
SELECT 
   DISTINCT(pathname), 
   COALESCE(COUNT(viewed_at), 0)::bigint AS views,
   COALESCE(COUNT(DISTINCT hashed_ip), 0)::bigint AS unique_visitors
FROM pageviews 
WHERE browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
 AND DATE_TRUNC('day', viewed_at)::date >= (NOW() - make_interval(days => sqlc.arg(days)::int))::date
 AND DATE_TRUNC('day', viewed_at)::date <= NOW()::date
GROUP BY pathname
ORDER BY views DESC;

-- name: GetRoutesSummaryForLast24Hours :many
SELECT 
   DISTINCT(pathname), 
   COALESCE(COUNT(viewed_at), 0)::bigint AS views,
   COALESCE(COUNT(DISTINCT hashed_ip), 0)::bigint AS unique_visitors
FROM pageviews 
WHERE browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
 AND viewed_at >= NOW() - INTERVAL '23 hours'
GROUP BY pathname
ORDER BY views DESC;

-- name: GetCountriesSummaryByDay :many
SELECT 
   DISTINCT(country), 
   COALESCE(COUNT(viewed_at), 0)::bigint AS views,
   COALESCE(COUNT(DISTINCT hashed_ip), 0)::bigint AS unique_visitors
FROM pageviews 
WHERE browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
 AND DATE_TRUNC('day', viewed_at)::date >= (NOW() - make_interval(days => sqlc.arg(days)::int))::date
 AND DATE_TRUNC('day', viewed_at)::date <= NOW()::date
GROUP BY country
ORDER BY views DESC;

-- name: GetCountriesSummaryForLast24Hours :many
SELECT 
   DISTINCT(country), 
   COALESCE(COUNT(viewed_at), 0)::bigint AS views,
   COALESCE(COUNT(DISTINCT hashed_ip), 0)::bigint AS unique_visitors
FROM pageviews 
WHERE browser NOT IN ('Headless Chrome', 'Google-Read-Aloud')
 AND viewed_at >= NOW() - INTERVAL '23 hours'
GROUP BY country
ORDER BY views DESC;



