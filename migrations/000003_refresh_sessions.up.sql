CREATE TABLE todoapp.auth_sessions (
    id UUID PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES todoapp.users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,

    CHECK (expires_at > created_at)
);

CREATE INDEX auth_sessions_user_id_idx ON todoapp.auth_sessions(user_id);

CREATE TABLE todoapp.refresh_tokens (
    id UUID PRIMARY KEY,
    session_id UUID NOT NULL REFERENCES todoapp.auth_sessions(id) ON DELETE CASCADE,
    token_hash BYTEA NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,

    CHECK (octet_length(token_hash) = 32),
    CHECK (expires_at > created_at)
);

CREATE INDEX refresh_tokens_session_id_idx ON todoapp.refresh_tokens(session_id);