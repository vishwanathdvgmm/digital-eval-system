package authority

import "time"

// Authority-side request / assignment models

type RequestRow struct {
	ID           int64     `json:"id"`
	EvaluatorID  string    `json:"evaluator_id"`
	CourseID     string    `json:"course_id"`
	Semester     string    `json:"semester"`
	AcademicYear string    `json:"academic_year"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

type ApprovePayload struct {
	RequestID int64 `json:"request_id"`
	AssignNum int   `json:"assign_num"` // number of scripts to assign (default 5)
}
