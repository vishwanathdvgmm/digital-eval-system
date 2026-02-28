package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

const insertEvaluationResultSQL = `
INSERT INTO evaluations (script_id, student_usn, course_id, semester, academic_year, course_credits, evaluator_id, marks, total_marks, result, block_hash) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11) RETURNING id;`

// InsertEvaluationResult inserts evaluated data into evaluations table.
func (p *PostgresDB) InsertEvaluationResult(ctx context.Context, scriptID string, studentUSN string, courseID string, semester string, academicYear string, courseCredits int, evaluatorID string, marksJSON []byte, totalMarks int, result string, blockHash string) error {

	var id int
	err := p.DB.QueryRowContext(
		ctx,
		insertEvaluationResultSQL,
		scriptID,
		studentUSN,
		courseID,
		semester,
		academicYear,
		courseCredits,
		evaluatorID,
		marksJSON,
		totalMarks,
		result,
		blockHash,
	).Scan(&id)

	return err
}

type EvaluationRow struct {
	ID            int64
	ScriptID      string
	StudentUSN    sql.NullString
	CourseID      string
	Semester      string
	AcademicYear  string
	CourseCredits sql.NullInt32
	Evaluator     string
	Marks         json.RawMessage
	TotalMarks    int
	Result        string
	CreatedAt     time.Time
}

func (p *PostgresDB) FetchResultsByUSN(ctx context.Context, usn string, academicYear string) ([]EvaluationRow, error) {
	rows, err := p.DB.QueryContext(ctx, `SELECT id, script_id, student_usn, course_id, semester, academic_year, course_credits, evaluator_id, marks, total_marks, result, created_at FROM evaluations WHERE student_usn=$1 AND academic_year=$2 ORDER BY created_at DESC`, usn, academicYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []EvaluationRow
	for rows.Next() {
		var r EvaluationRow
		if err := rows.Scan(&r.ID, &r.ScriptID, &r.StudentUSN, &r.CourseID, &r.Semester, &r.AcademicYear, &r.CourseCredits, &r.Evaluator, &r.Marks, &r.TotalMarks, &r.Result, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (p *PostgresDB) FetchResultsBySemester(ctx context.Context, semester string, academicYear string) ([]EvaluationRow, error) {
	rows, err := p.DB.QueryContext(ctx, `SELECT id, script_id, student_usn, course_id, semester, academic_year, course_credits, evaluator_id, marks, total_marks, result, created_at FROM evaluations WHERE semester = $1 AND academic_year = $2 ORDER BY created_at ASC`, semester, academicYear)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []EvaluationRow
	for rows.Next() {
		var r EvaluationRow
		if err := rows.Scan(&r.ID, &r.ScriptID, &r.StudentUSN, &r.CourseID, &r.Semester, &r.AcademicYear, &r.CourseCredits, &r.Evaluator, &r.Marks, &r.TotalMarks, &r.Result, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// RecordRelease inserts release record
func (p *PostgresDB) RecordRelease(ctx context.Context, semester, academicYear, releasedBy, blockHash string) (int64, error) {
	var id int64
	err := p.DB.QueryRowContext(ctx, `INSERT INTO result_releases (semester, academic_year, released_by, block_hash, released_at) VALUES ($1,$2,$3,$4, now()) RETURNING id`, semester, academicYear, releasedBy, blockHash).Scan(&id)
	return id, err
}
