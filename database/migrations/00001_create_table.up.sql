BEGIN;

CREATE TABLE IF NOT EXISTS users
(
    id          UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL,
    email       TEXT NOT NULL,
    password    TEXT NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS unique_user ON users (email) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS todos
(
    id           UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    user_id      UUID REFERENCES users (id) NOT NULL,
    title         TEXT                       NOT NULL,
    description  TEXT                       NOT NULL,
    is_completed BOOLEAN                  DEFAULT FALSE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at  TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS unique_todo ON todos (user_id,title) WHERE archived_at IS NULL;

CREATE TABLE IF NOT EXISTS user_session
(
    id          UUID PRIMARY KEY         DEFAULT gen_random_uuid(),
    user_id     UUID REFERENCES users (id) NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

COMMIT;