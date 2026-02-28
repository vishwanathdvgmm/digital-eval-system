package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	DB *sqlx.DB
}

type PostgresConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func NewPostgres(cfg PostgresConfig) (*PostgresDB, error) {
	db, err := sqlx.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}
	if cfg.MaxOpenConns > 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns > 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConns)
	}
	if cfg.ConnMaxLifetime > 0 {
		db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	}
	// ping
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresDB{DB: db}, nil
}

func (p *PostgresDB) Close() error {
	if p.DB == nil {
		return nil
	}
	return p.DB.Close()
}

// EvaluationExistsForScript returns true if an evaluation exists for the given script_id.
func (p *PostgresDB) EvaluationExistsForScript(ctx context.Context, scriptID string) (bool, error) {
	var cnt int
	err := p.DB.QueryRowContext(ctx, `SELECT count(1) FROM evaluations WHERE script_id = $1`, scriptID).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// EvaluationExistsForStudentSemesterCourse returns true if an evaluation exists for (student_usn, semester, academic_year, course_id)
func (p *PostgresDB) EvaluationExistsForStudentSemesterCourse(ctx context.Context, studentUSN, semester, academicYear, courseID string) (bool, error) {
	var cnt int
	err := p.DB.QueryRowContext(ctx, `SELECT count(1) FROM evaluations WHERE student_usn = $1 AND semester = $2 AND academic_year = $3 AND course_id = $4`, studentUSN, semester, academicYear, courseID).Scan(&cnt)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// GetAssignedScriptByScriptID fetches an assigned_scripts row by script_id
// returns nil, nil if not found.
func (p *PostgresDB) GetAssignedScriptByScriptID(ctx context.Context, scriptID string) (*AssignedScriptRow, error) {
	// Assumes there is a struct AssignedScriptRow in your queries or models package.
	var row AssignedScriptRow
	err := p.DB.QueryRowContext(ctx, `
        SELECT id, script_id, evaluator_id, course_id, semester, academic_year, course_credits, assigned_at, status
        FROM assigned_scripts
        WHERE script_id = $1
        LIMIT 1
    `, scriptID).Scan(&row.ID, &row.ScriptID, &row.Evaluator, &row.CourseID, &row.Semester, &row.AcademicYear, &row.CourseCredits, &row.AssignedAt, &row.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}
