BEGIN;

CREATE TABLE IF NOT EXISTS audit_logs (
    id serial PRIMARY KEY,
    user_id text NOT NULL,
    action text NOT NULL,
    detail jsonb,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT chk_action_not_empty CHECK (length(trim(action)) > 0)
);

CREATE INDEX IF NOT EXISTS idx_audit_user_id ON audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_created_at ON audit_logs(created_at);


COMMIT;
