BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id serial PRIMARY KEY,
    user_id text UNIQUE NOT NULL,         -- login username
    role text NOT NULL,                    -- admin / authority / examiner / evaluator / student
    name text,
    email text UNIQUE,
    password_hash text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_users_user_id ON users(user_id);
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

COMMIT;
