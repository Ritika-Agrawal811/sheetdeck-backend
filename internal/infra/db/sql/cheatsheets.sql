-- name: CreateCheatsheet :exec
INSERT INTO cheatsheets (slug, title, category, subcategory, image_url)
VALUES ($1, $2, $3, $4, $5);


-- name: GetCheatsheetByID :one
SELECT id, slug, title, category, subcategory, image_url, created_at, updated_at
FROM cheatsheets
WHERE id = $1;

-- name: GetCheatsheetBySlug :one
SELECT id, slug, title, category, subcategory, image_url, created_at, updated_at
FROM cheatsheets
WHERE slug = $1;    

-- name: ListCheatsheets :many
SELECT id, slug, title, category::varchar, subcategory::varchar, image_url, created_at, updated_at
FROM cheatsheets
WHERE (sqlc.narg(category)::category IS NULL OR category = sqlc.narg(category))
  AND (sqlc.narg(subcategory)::subcategory IS NULL OR subcategory = sqlc.narg(subcategory))
ORDER BY created_at DESC
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



