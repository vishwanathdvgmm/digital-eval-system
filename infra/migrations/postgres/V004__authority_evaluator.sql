-- V004__phase5_authority_evaluator.sql
-- Tables for authority / evaluator workflows

BEGIN;

CREATE TABLE IF NOT EXISTS evaluation_requests (
    id serial PRIMARY KEY,
    evaluator_id text NOT NULL,
    course_id text NOT NULL,
    semester text NOT NULL,
    academic_year text NOT NULL DEFAULT '2025-2026',
    description text,
    status text NOT NULL DEFAULT 'pending', -- pending / approved / rejected
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT chk_eval_status_not_empty CHECK (length(trim(status)) > 0)
);

CREATE INDEX idx_eval_req_pending 
ON evaluation_requests (evaluator_id, course_id, semester)
WHERE status = 'pending';
CREATE INDEX IF NOT EXISTS idx_eval_requests_evaluator_id 
    ON evaluation_requests(evaluator_id);

CREATE INDEX IF NOT EXISTS idx_eval_requests_status 
    ON evaluation_requests(status);

CREATE TABLE IF NOT EXISTS assigned_scripts (
    id serial PRIMARY KEY,
    script_id text NOT NULL,
    evaluator_id text NOT NULL,
    course_id text NOT NULL,
    semester text NOT NULL,
    academic_year TEXT NOT NULL DEFAULT '2025-2026',
    course_credits integer NOT NULL DEFAULT 0,
    assigned_at timestamptz NOT NULL DEFAULT now(),
    status text NOT NULL DEFAULT 'assigned', -- assigned / in_progress / evaluated / revoked
    CONSTRAINT uq_assigned_script UNIQUE (script_id),
    CONSTRAINT chk_assigned_status_not_empty CHECK (length(trim(status)) > 0)
);

CREATE INDEX IF NOT EXISTS idx_assigned_scripts_evaluator_id 
    ON assigned_scripts(evaluator_id);

CREATE INDEX IF NOT EXISTS idx_assigned_scripts_status 
    ON assigned_scripts(status);

COMMIT;
