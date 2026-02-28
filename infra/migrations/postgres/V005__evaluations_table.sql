BEGIN;

DROP TABLE IF EXISTS evaluations;

CREATE TABLE evaluations (
    id SERIAL PRIMARY KEY,

    script_id TEXT NOT NULL,
    student_usn TEXT NOT NULL,

    course_id TEXT NOT NULL,
    semester TEXT NOT NULL,
    academic_year TEXT NOT NULL DEFAULT '2025-2026',

    course_credits INT NOT NULL DEFAULT 0,

    evaluator_id TEXT NOT NULL,

    marks JSONB NOT NULL,

    total_marks INT NOT NULL,
    result TEXT NOT NULL,

    block_hash TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_eval_once_per_subject UNIQUE (student_usn, course_id, semester)
);

CREATE INDEX idx_eval_script      ON evaluations(script_id);
CREATE INDEX idx_eval_student     ON evaluations(student_usn);
CREATE INDEX idx_eval_course_sem  ON evaluations(course_id, semester);

COMMIT;
