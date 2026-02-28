package db

import (
	"context"
	"time"
)

// Assigned script helpers

func (p *PostgresDB) CreateAssignment(ctx context.Context, scriptID, evaluatorID, courseID, semester string, academicYear string, courseCredits int) (int64, error) {
	var id int64
	err := p.DB.QueryRowContext(ctx,
		`INSERT INTO assigned_scripts (script_id, evaluator_id, course_id, semester, academic_year, course_credits, assigned_at, status)
		 VALUES ($1,$2,$3,$4,$5,$6, now(), 'assigned') RETURNING id`,
		scriptID, evaluatorID, courseID, semester, academicYear, courseCredits).Scan(&id)
	return id, err
}

type AssignedScriptRow struct {
	ID            int64
	ScriptID      string
	Evaluator     string
	CourseID      string
	Semester      string
	AcademicYear  string
	CourseCredits int
	AssignedAt    time.Time
	Status        string
}

func (p *PostgresDB) ListAssignedByEvaluator(ctx context.Context, evaluatorID string) ([]AssignedScriptRow, error) {
	rows, err := p.DB.QueryContext(ctx, `SELECT id,script_id,evaluator_id,course_id,semester,academic_year,course_credits,assigned_at,status FROM assigned_scripts WHERE evaluator_id=$1 ORDER BY assigned_at DESC`, evaluatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AssignedScriptRow
	for rows.Next() {
		var r AssignedScriptRow
		if err := rows.Scan(&r.ID, &r.ScriptID, &r.Evaluator, &r.CourseID, &r.Semester, &r.AcademicYear, &r.CourseCredits, &r.AssignedAt, &r.Status); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (p *PostgresDB) UpdateAssignmentStatus(ctx context.Context, id int64, status string) error {
	_, err := p.DB.ExecContext(ctx, `UPDATE assigned_scripts SET status=$1 WHERE id=$2`, status, id)
	return err
}

type EvalRequestHistoryRow struct {
	ID           int64     `json:"id"`
	CourseID     string    `json:"course_id"`
	Semester     string    `json:"semester"`
	AcademicYear string    `json:"academic_year"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func (p *PostgresDB) ListRequestsByEvaluator(ctx context.Context, evaluatorID string) ([]EvalRequestHistoryRow, error) {
	rows, err := p.DB.QueryContext(ctx, `SELECT id,course_id,semester,academic_year,description,status,created_at FROM evaluation_requests WHERE evaluator_id=$1 ORDER BY created_at DESC`, evaluatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []EvalRequestHistoryRow
	for rows.Next() {
		var r EvalRequestHistoryRow
		if err := rows.Scan(&r.ID, &r.CourseID, &r.Semester, &r.AcademicYear, &r.Description, &r.Status, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
