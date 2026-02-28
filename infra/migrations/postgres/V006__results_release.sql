BEGIN;

CREATE TABLE IF NOT EXISTS result_releases (
    id serial PRIMARY KEY,
    semester text NOT NULL,
    academic_year text NOT NULL DEFAULT '2025-2026',
    released_at timestamptz NOT NULL DEFAULT now(),
    released_by text NOT NULL,
    block_hash text,
    CONSTRAINT uq_result_releases_semester_year UNIQUE (semester, academic_year),
    CONSTRAINT chk_released_by_not_empty CHECK (length(trim(released_by)) > 0)
);

COMMIT;
