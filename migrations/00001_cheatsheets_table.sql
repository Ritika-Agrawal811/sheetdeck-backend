-- +goose Up

-- First create UUID extension if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- First create type if not exists
CREATE TYPE category AS ENUM ('html', 'css', 'javascript', 'react');
CREATE TYPE subcategory AS ENUM (
  'concepts',
  'attributes',
  'elements',
  'properties',
  'pseudo_classes',
  'methods',
  'selectors',
  'advanced_syntax',
  'dom_manipulation',
  'operators'
);

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS cheatsheets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slug VARCHAR(255) UNIQUE NOT NULL,
    title TEXT NOT NULL,
    category category NOT NULL,
    subcategory subcategory NOT NULL,
    image_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Add indexes for faster filtering
CREATE INDEX IF NOT EXISTS idx_cheatsheets_category ON cheatsheets (category);
CREATE INDEX IF NOT EXISTS idx_cheatsheets_subcategory ON cheatsheets (subcategory);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_cheatsheets_category;
DROP INDEX IF EXISTS idx_cheatsheets_subcategory;
DROP TYPE IF EXISTS category;
DROP TYPE IF EXISTS subcategory;
DROP TABLE IF EXISTS cheatsheets;
-- +goose StatementEnd
