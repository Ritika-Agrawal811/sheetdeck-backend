-- +goose Up

CREATE TYPE event_type AS ENUM ('view', 'click', 'download', 'search', 'copy_code');

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    cheatsheet_id UUID NOT NULL REFERENCES cheatsheets(id),    
    event_type event_type NOT NULL,      
    pathname TEXT NOT NULL,         
    hashed_ip TEXT NOT NULL,       
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_events_type ON events(event_type);
CREATE INDEX idx_events_cheatsheet ON events(cheatsheet_id);
CREATE INDEX idx_events_created_at ON events(created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_events_type ON events(event_type);
DROP INDEX IF EXISTS idx_events_cheatsheet ON events(cheatsheet_id);
DROP INDEX IF EXISTS idx_events_created_at ON events(created_at);
DROP TABLE IF EXISTS events;
-- +goose StatementEnd
