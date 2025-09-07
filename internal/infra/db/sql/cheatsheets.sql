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
SELECT id, slug, title, category, subcategory, image_url, created_at, updated_at
FROM cheatsheets
WHERE ($1::category IS NULL OR category = $1)
  AND ($2::subcategory IS NULL OR subcategory = $2)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4; 