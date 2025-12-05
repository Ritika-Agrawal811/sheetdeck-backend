-- name: CreateCheatsheet :exec
INSERT INTO cheatsheets (slug, title, category, subcategory, image_url, image_size)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: GetCheatsheetByID :one
SELECT id, slug, title, category, subcategory, image_url, created_at, updated_at, image_size
FROM cheatsheets
WHERE id = $1;

-- name: GetCheatsheetBySlug :one
SELECT id, slug, title, category, subcategory, image_url, created_at, updated_at, image_size
FROM cheatsheets
WHERE slug = $1;    

-- name: ListCheatsheets :many
SELECT 
    c.id, 
    c.slug, 
    c.title, 
    c.category::varchar AS category, 
    c.subcategory::varchar AS subcategory, 
    c.image_url, 
    c.image_size, 
    c.created_at, 
    c.updated_at,
    COUNT(e.hashed_ip) FILTER (WHERE e.event_type = 'download') AS downloads,
    COUNT(e.hashed_ip) FILTER (WHERE e.event_type = 'click') AS views
FROM cheatsheets c
LEFT JOIN events e ON e.cheatsheet_id = c.id
WHERE (sqlc.narg(category)::category IS NULL OR c.category = sqlc.narg(category))
  AND (sqlc.narg(subcategory)::subcategory IS NULL OR c.subcategory = sqlc.narg(subcategory))
GROUP BY c.id
ORDER BY 
  CASE sqlc.arg(sort_by)::text
    WHEN 'recent' THEN c.created_at
  END DESC,
  CASE sqlc.arg(sort_by)::text
    WHEN 'oldest' THEN c.created_at
  END ASC,
  CASE sqlc.arg(sort_by)::text
    WHEN 'most_downloaded' THEN COUNT(e.hashed_ip) FILTER (WHERE e.event_type = 'download')
    WHEN 'most_viewed' THEN COUNT(e.hashed_ip) FILTER (WHERE e.event_type = 'click')
  END DESC,
  CASE sqlc.arg(sort_by)::text
    WHEN 'least_downloaded' THEN COUNT(e.hashed_ip) FILTER (WHERE e.event_type = 'download')
    WHEN 'least_viewed' THEN COUNT(e.hashed_ip) FILTER (WHERE e.event_type = 'click')
  END ASC
LIMIT $1 OFFSET $2;


-- name: UpdateCheatsheet :exec
UPDATE cheatsheets
SET slug = COALESCE(NULLIF(sqlc.arg(slug)::varchar, ''), slug),
    title = COALESCE(NULLIF(sqlc.arg(title)::text, ''), title),
    category = COALESCE(NULLIF(sqlc.narg(category), '')::category, category),
    subcategory = COALESCE(NULLIF(sqlc.narg(subcategory), '')::subcategory, subcategory),
    image_url = COALESCE(NULLIF(sqlc.arg(image_url)::text, ''), image_url),
    updated_at = NOW()
WHERE id = sqlc.arg(id);

-- name: GetTotalCheasheetsCount :one
Select COUNT(id) from cheatsheets;

-- name: CountCheatsheetsByCategoryAndSubcategory :many
SELECT category, subcategory, COUNT(DISTINCT id) AS cheatsheet_count
FROM cheatsheets
GROUP BY category, subcategory;

-- name: GetCategoryDetails :many
SELECT category, COUNT(DISTINCT id) AS cheatsheet_count, ARRAY_AGG(DISTINCT subcategory)::varchar[] AS subcategories
FROM cheatsheets
GROUP BY category;

-- name: GetCategories :many
SELECT unnest(enum_range(NULL::category))::varchar as categories;

-- name: GetSubcategories :many
SELECT unnest(enum_range(NULL::subcategory))::varchar as subcategories;

-- name: GetTotalImageSize :one
SELECT 
  COALESCE(SUM(image_size), 0)::bigint as total_size,
  pg_size_pretty(COALESCE(SUM(image_size), 0)) as total_size_pretty
FROM cheatsheets;

-- name: GetLargestCheatsheets :many
SELECT
  title,
  category::varchar,
  pg_size_pretty(image_size) AS size
FROM cheatsheets
ORDER BY image_size DESC
LIMIT 2;

