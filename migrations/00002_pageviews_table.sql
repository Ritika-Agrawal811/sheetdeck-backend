-- +goose Up

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pageviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pathname TEXT NOT NULL,
    ip_address TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    referrer TEXT,
    viewed_at TIMESTAMP DEFAULT NOW()
);

-- Add indexes for faster filtering
CREATE INDEX idx_pageviews_pathname ON pageviews(pathname);
CREATE INDEX idx_pageviews_ip ON pageviews(ip_address);
CREATE INDEX idx_pageviews_viewed_at ON pageviews(viewed_at);
CREATE INDEX idx_pageviews_referrer ON pageviews(referrer);
CREATE INDEX idx_pageviews_route_date ON pageviews(pathname, viewed_at);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_pageviews_referrer;
DROP INDEX IF EXISTS idx_pageviews_viewed_at;
DROP INDEX IF EXISTS idx_pageviews_ip;
DROP INDEX IF EXISTS idx_pageviews_pathname;
DROP INDEX IF EXISTS idx_pageviews_route_date;

DROP TABLE IF EXISTS pageviews;
-- +goose StatementEnd
