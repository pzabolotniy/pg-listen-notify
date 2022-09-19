-- +migrate Up
-- +migrate StatementBegin
ALTER TABLE events ADD COLUMN IF NOT EXISTS locked bool NOT NULL DEFAULT false;
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
ALTER TABLE events DROP COLUMN IF EXISTS locked;
-- +migrate StatementEnd
