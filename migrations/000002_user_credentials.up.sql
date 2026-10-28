CREATE TABLE todoapp.user_credentials (
    user_id INTEGER PRIMARY KEY REFERENCES todoapp.users(id) ON DELETE CASCADE,
    login VARCHAR(64) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    CHECK (char_length(login) BETWEEN 3 AND 64),
    CHECK (login = lower(login)),
    CHECK (login ~ '^[a-z0-9_]+$')
)