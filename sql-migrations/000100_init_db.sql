-- +migrate Up
-- +migrate StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS events (
    id uuid PRIMARY KEY,
    payload jsonb NOT NULL,
    received_at TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
-- +migrate StatementEnd

-- +migrate Down
-- +migrate StatementBegin
DROP TABLE IF EXISTS events;
DROP EXTENSION IF EXISTS "uuid-ossp";
-- +migrate StatementEnd