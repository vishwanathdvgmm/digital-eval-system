package authority

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"digital-eval-system/services/go-node/internal/block"
	"digital-eval-system/services/go-node/internal/db"
)

// ReleaseService performs semester result release
type ReleaseService struct {
	pg    *db.PostgresDB
	chain interface {
		AppendBlock(*block.Block) (string, error)
	}
}

func NewReleaseService(pg *db.PostgresDB, chain interface {
	AppendBlock(*block.Block) (string, error)
}) *ReleaseService {
	return &ReleaseService{pg: pg, chain: chain}
}

// ReleaseResults aggregates evaluations for semester, writes a release block, and records release.
func (s *ReleaseService) ReleaseResults(ctx context.Context, semester, academicYear, releasedBy string) (string, error) {
	// 1. fetch all evaluations for semester
	rows, err := s.pg.FetchResultsBySemester(ctx, semester, academicYear)
	if err != nil {
		return "", err
	}

	if len(rows) == 0 {
		return "", fmt.Errorf("no evaluations for semester %s", semester)
	}

	// 2. prepare aggregated JSON: list of records
	records := []map[string]interface{}{}
	for _, r := range rows {
		rec := map[string]interface{}{
			"script_id":     r.ScriptID,
			"student_usn":   nil,
			"course_id":     r.CourseID,
			"semester":      r.Semester,
			"academic_year": r.AcademicYear,
			"evaluator_id":  r.Evaluator,
			"marks":         json.RawMessage(r.Marks),
			"total_marks":   r.TotalMarks,
			"result":        r.Result,
			"created_at":    r.CreatedAt,
		}
		if r.StudentUSN.Valid {
			rec["student_usn"] = r.StudentUSN.String
		}
		records = append(records, rec)
	}

	data, _ := json.Marshal(records)

	// 3. create a ResultRelease block transaction
	tx := block.Transaction{
		ScriptID:     "", // not per-script, it's an aggregate release
		USN:          "",
		CourseID:     "",
		Semester:     semester,
		AcademicYear: academicYear,
		CID:          "",
		Meta:         map[string]string{"_result_release": string(data)},
		CreatedAt:    time.Now().Unix(),
		SignerID:     releasedBy,
	}
	newBlock := block.NewBlock("", []block.Transaction{tx}, releasedBy)
	blockHash, err := s.chain.AppendBlock(newBlock)
	if err != nil {
		return "", err
	}

	// 4. insert a record in postgres
	_, err = s.pg.RecordRelease(ctx, semester, academicYear, releasedBy, blockHash)
	if err != nil {
		logrus.Warnf("failed to record release in pg: %v", err)
	}

	return blockHash, nil
}
