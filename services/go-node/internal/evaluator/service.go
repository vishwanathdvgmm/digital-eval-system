package evaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"digital-eval-system/services/go-node/internal/block"
	"digital-eval-system/services/go-node/internal/db"
	"digital-eval-system/services/go-node/internal/pybridge"
	"digital-eval-system/services/go-node/internal/storage"
)

// Service orchestrates evaluator actions.
type Service struct {
	pg    *db.PostgresDB
	store storage.Storage
	py    *pybridge.Client
	chain interface {
		AppendBlock(*block.Block) (string, error)
	}
}

// NewService creates evaluator service
func NewService(pg *db.PostgresDB, st storage.Storage, py *pybridge.Client, chain interface {
	AppendBlock(*block.Block) (string, error)
}) *Service {
	return &Service{pg: pg, store: st, py: py, chain: chain}
}

func (s *Service) CreateRequest(ctx context.Context, evaluatorID, courseID, semester, academicYear, desc string) (int64, error) {
	return s.pg.InsertEvaluationRequest(ctx, evaluatorID, courseID, semester, academicYear, desc)
}

func (s *Service) ListRequests(ctx context.Context, evaluatorID string) ([]db.EvalRequestHistoryRow, error) {
	return s.pg.ListRequestsByEvaluator(ctx, evaluatorID)
}

// ListAssigned returns assigned scripts for evaluator
func (s *Service) ListAssigned(ctx context.Context, evaluatorID string) ([]db.AssignedScriptRow, error) {
	return s.pg.ListAssignedByEvaluator(ctx, evaluatorID)
}

// GetScriptMetadata fetches metadata for a script from BoltDB
// returns map[string]string or error if not found
func (s *Service) GetScriptMetadata(ctx context.Context, scriptID string) (map[string]string, error) {
	var result map[string]string

	err := s.store.ForEachBlock(func(blk *block.Block) {
		for _, tx := range blk.Transactions {
			// match script id (case-insensitive trim)
			if strings.EqualFold(strings.TrimSpace(tx.ScriptID), strings.TrimSpace(scriptID)) {
				// prefer upload records (uploaded blocks store metadata and set _upload_record)
				if _, isUpload := tx.Meta["_upload_record"]; isUpload {
					// make a copy to avoid mutating the block's map if it's shared
					result = make(map[string]string)
					for k, v := range tx.Meta {
						result[k] = v
					}
					result["pdf_cid"] = tx.CID
					return
				}
				// if no explicit upload marker but meta present, accept it as metadata
				if result == nil {
					result = make(map[string]string)
					for k, v := range tx.Meta {
						result[k] = v
					}
					result["pdf_cid"] = tx.CID
				}
			}
		}
	})
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("not found")
	}
	return result, nil
}

// SubmitEvaluation validates marks via Python and stores record on blockchain + Postgres
func (s *Service) SubmitEvaluation(ctx context.Context, payload SubmitPayload) (string, error) {
	// basic validation
	if payload.ScriptID == "" || payload.EvaluatorID == "" {
		return "", fmt.Errorf("missing fields")
	}

	// local strict validation
	if payload.TotalMarks <= 0 {
		return "", fmt.Errorf("total_marks must be > 0")
	}
	if payload.TotalQuestions <= 0 {
		return "", fmt.Errorf("total_questions must be > 0")
	}
	if len(payload.MarksScored) == 0 {
		return "", fmt.Errorf("marks_scored cannot be empty")
	}
	if payload.QuestionsAnswered != len(payload.MarksScored) {
		return "", fmt.Errorf("questions_answered mismatch: expected %d marks, got %d", payload.QuestionsAnswered, len(payload.MarksScored))
	}
	if len(payload.MarksAllotted) != 0 && len(payload.MarksAllotted) != len(payload.MarksScored) {
		return "", fmt.Errorf("marks_allotted_per_question length mismatch")
	}
	sumScored := 0
	for i, sc := range payload.MarksScored {
		if sc < 0 {
			return "", fmt.Errorf("marks_scored[%d] negative", i)
		}
		if len(payload.MarksAllotted) == len(payload.MarksScored) {
			if sc > payload.MarksAllotted[i] {
				return "", fmt.Errorf("marks_scored[%d] out of range: got %d, max %d", i, sc, payload.MarksAllotted[i])
			}
		}
		sumScored += sc
	}
	if sumScored > payload.TotalMarks {
		return "", fmt.Errorf("sum of marks_scored (%d) exceeds total_marks (%d)", sumScored, payload.TotalMarks)
	}

	// fetch original metadata to get USN/course/semester
	meta, err := s.GetScriptMetadata(ctx, payload.ScriptID)
	if err != nil {
		return "", fmt.Errorf("cannot fetch script metadata: %w", err)
	}
	usn := strings.TrimSpace(meta["USN"])
	course := payload.CourseID
	if course == "" {
		course = strings.TrimSpace(meta["CourseID"])
	}
	semester := payload.Semester
	if semester == "" {
		semester = strings.TrimSpace(meta["Semester"])
	}

	academicYear := payload.AcademicYear
	if academicYear == "" {
		academicYear = strings.TrimSpace(meta["AcademicYear"])
	}

	// ensure script is assigned to evaluator and not already evaluated
	assigned, err := s.pg.GetAssignedScriptByScriptID(ctx, payload.ScriptID)
	if err != nil {
		return "", fmt.Errorf("failed to fetch assignment: %w", err)
	}
	if assigned == nil || assigned.Evaluator != payload.EvaluatorID || assigned.Status != "assigned" {
		return "", fmt.Errorf("script not assigned to evaluator or not in assigned state")
	}
	// assignedScripts retrieved via GetAssignedScript must include academic_year
	if assigned.AcademicYear != payload.AcademicYear {
		return "", fmt.Errorf("academic_year mismatch: assigned=%s submitted=%s",
			assigned.AcademicYear, payload.AcademicYear)
	}

	// ensure no existing evaluation for this script
	exists, err := s.pg.EvaluationExistsForScript(ctx, payload.ScriptID)
	if err != nil {
		return "", fmt.Errorf("db check failed: %w", err)
	}
	if exists {
		return "", fmt.Errorf("evaluation already submitted for this script")
	}

	// ensure no existing evaluation for this student+semester+course (enforce single row per subject)
	dup, err := s.pg.EvaluationExistsForStudentSemesterCourse(ctx, usn, semester, academicYear, course)
	if err != nil {
		return "", fmt.Errorf("db check failed: %w", err)
	}
	if dup {
		return "", fmt.Errorf("evaluation already exists for student %s semester %s course %s", usn, semester, course)
	}

	// create marks JSON
	marksStruct := map[string]interface{}{
		"total_questions":     payload.TotalQuestions,
		"marks_per_question":  payload.MarksPerQuestion,
		"total_marks":         payload.TotalMarks,
		"course_id":           course,
		"semester":            semester,
		"academic_year":       payload.AcademicYear,
		"questions_answered":  payload.QuestionsAnswered,
		"marks_allotted":      payload.MarksAllotted,
		"marks_scored":        payload.MarksScored,
		"additional_metadata": payload.AdditionalMetadata,
		"course_credits":      payload.CourseCredits,
	}
	marksJSON, _ := json.Marshal(marksStruct)

	// create evaluation transaction and append to chain
	tx := block.Transaction{
		ScriptID:     payload.ScriptID,
		USN:          usn,
		CourseID:     course,
		Semester:     semester,
		AcademicYear: payload.AcademicYear,
		CID:          "",
		Meta:         map[string]string{"_evaluation": string(marksJSON)},
		CreatedAt:    int64(time.Now().Unix()),
		SignerID:     payload.EvaluatorID,
	}
	prevHash := ""
	if h, herr := s.store.Head(); herr == nil {
		prevHash = h
	}
	newBlock := block.NewBlock(prevHash, []block.Transaction{tx}, payload.EvaluatorID)
	blockHash, err := s.chain.AppendBlock(newBlock)
	if err != nil {
		return "", fmt.Errorf("append block failed: %w", err)
	}

	// persist evaluation: use corrected InsertEvaluationResult signature
	if err := s.pg.InsertEvaluationResult(ctx, payload.ScriptID, usn, course, semester, payload.AcademicYear, payload.CourseCredits, payload.EvaluatorID, marksJSON, payload.TotalMarks, "PASS", blockHash); err != nil {
		logrus.Warnf("failed to insert evaluation to pg: %v", err)
		// block appended; return error to caller
		return "", fmt.Errorf("insert evaluation failed: %w", err)
	}

	return blockHash, nil
}
