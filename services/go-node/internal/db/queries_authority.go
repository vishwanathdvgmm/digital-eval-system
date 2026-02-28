package db

import (
	"context"
	"fmt"
	"time"
)

// Authority-related DB helper functions

// InsertEvaluationRequest inserts a new evaluator request and returns id.
func (p *PostgresDB) InsertEvaluationRequest(ctx context.Context, evaluatorID, courseID, semester, academicYear, desc string) (int64, error) {

	// 1. Check for existing pending request
	var existingID int64
	err := p.DB.QueryRowContext(ctx,
		`SELECT id FROM evaluation_requests 
         WHERE evaluator_id=$1 AND course_id=$2 AND semester=$3 AND academic_year=$4 AND status='pending'`,
		evaluatorID, courseID, semester, academicYear,
	).Scan(&existingID)

	if err == nil {
		// There is already a pending request
		return 0, fmt.Errorf("pending request exists")
	}

	// 2. Delete old non-pending requests (approved/rejected)
	_, _ = p.DB.ExecContext(ctx,
		`DELETE FROM evaluation_requests 
         WHERE evaluator_id=$1 AND course_id=$2 AND semester=$3 AND academic_year=$4 AND status!='pending'`,
		evaluatorID, courseID, semester, academicYear,
	)

	// 3. Insert new request
	var id int64
	err = p.DB.QueryRowContext(ctx,
		`INSERT INTO evaluation_requests (evaluator_id, course_id, semester, academic_year, description, status, created_at, updated_at)
         VALUES ($1,$2,$3,$4,$5,'pending', now(), now()) RETURNING id`,
		evaluatorID, courseID, semester, academicYear, desc,
	).Scan(&id)

	return id, err
}

// ListPendingRequests returns rows of pending requests.
type EvalRequestRow struct {
	ID           int64     `json:"id"`
	EvaluatorID  string    `json:"evaluator_id"`
	CourseID     string    `json:"course_id"`
	Semester     string    `json:"semester"`
	AcademicYear string    `json:"academic_year"`
	Description  string    `json:"description"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

func (p *PostgresDB) ListPendingRequests(ctx context.Context) ([]EvalRequestRow, error) {
	rows, err := p.DB.QueryContext(ctx, `SELECT id,evaluator_id,course_id,semester,academic_year,description,created_at FROM evaluation_requests WHERE status='pending' ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []EvalRequestRow
	for rows.Next() {
		var r EvalRequestRow
		if err := rows.Scan(&r.ID, &r.EvaluatorID, &r.CourseID, &r.Semester, &r.AcademicYear, &r.Description, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// UpdateRequestStatus sets status to approved/rejected and updates updated_at.
func (p *PostgresDB) UpdateRequestStatus(ctx context.Context, requestID int64, status string) error {
	_, err := p.DB.ExecContext(ctx, `UPDATE evaluation_requests SET status=$1, updated_at=now() WHERE id=$2`, status, requestID)
	return err
}

// ListRequestHistory returns rows of non-pending requests (approved/rejected).
func (p *PostgresDB) ListRequestHistory(ctx context.Context) ([]EvalRequestRow, error) {
	rows, err := p.DB.QueryContext(ctx, `SELECT id,evaluator_id,course_id,semester,academic_year,description,status,created_at FROM evaluation_requests WHERE status!='pending' ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []EvalRequestRow
	for rows.Next() {
		var r EvalRequestRow
		if err := rows.Scan(&r.ID, &r.EvaluatorID, &r.CourseID, &r.Semester, &r.AcademicYear, &r.Description, &r.Status, &r.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}
