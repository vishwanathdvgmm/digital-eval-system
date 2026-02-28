BEGIN;

-- Legacy results table (not used by final release logic, but safe to keep)
CREATE TABLE IF NOT EXISTS results (
    id serial PRIMARY KEY,
    student_usn text NOT NULL,
    course_id text NOT NULL,
    semester text NOT NULL,
    marks int NOT NULL,
    total_marks int NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT uq_results_student_course_semester UNIQUE (student_usn, course_id, semester),
    CONSTRAINT chk_marks_valid CHECK (
        marks >= 0 AND
        total_marks > 0 AND
        marks <= total_marks
    )
);

CREATE INDEX IF NOT EXISTS idx_results_usn ON results(student_usn);
CREATE INDEX IF NOT EXISTS idx_results_semester ON results(semester);


COMMIT;
