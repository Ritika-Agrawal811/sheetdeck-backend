-- +goose Up

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pageviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pathname TEXT NOT NULL,
    hashed_ip TEXT NOT NULL,        -- hashed for uniqueness tracking
    country TEXT,                   -- country name
    browser TEXT,                   -- e.g. Chrome, Safari
    os TEXT,                        -- e.g. macOS, Windows
    device TEXT,                    -- e.g. desktop, mobile
    user_agent TEXT NOT NULL,       -- raw UA string (keep for reference)
    referrer TEXT,
    viewed_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX idx_pageviews_pathname ON pageviews(pathname);
CREATE INDEX idx_pageviews_hashed_ip ON pageviews(hashed_ip);
CREATE INDEX idx_pageviews_viewed_at ON pageviews(viewed_at);
CREATE INDEX idx_pageviews_referrer ON pageviews(referrer);
CREATE INDEX idx_pageviews_route_date ON pageviews(pathname, viewed_at);
CREATE INDEX idx_pageviews_country ON pageviews(country);
CREATE INDEX idx_pageviews_browser ON pageviews(browser);
CREATE INDEX idx_pageviews_os ON pageviews(os);
CREATE INDEX idx_pageviews_device ON pageviews(device);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_pageviews_device;
DROP INDEX IF EXISTS idx_pageviews_os;
DROP INDEX IF EXISTS idx_pageviews_browser;
DROP INDEX IF EXISTS idx_pageviews_country;
DROP INDEX IF EXISTS idx_pageviews_referrer;
DROP INDEX IF EXISTS idx_pageviews_viewed_at;
DROP INDEX IF EXISTS idx_pageviews_hashed_ip;
DROP INDEX IF EXISTS idx_pageviews_pathname;
DROP INDEX IF EXISTS idx_pageviews_route_date;

DROP TABLE IF EXISTS pageviews;
-- +goose StatementEnd
