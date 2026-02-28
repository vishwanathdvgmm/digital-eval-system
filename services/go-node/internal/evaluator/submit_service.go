package evaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"digital-eval-system/services/go-node/internal/block"
	"digital-eval-system/services/go-node/internal/chain"
	"digital-eval-system/services/go-node/internal/db"
	"digital-eval-system/services/go-node/internal/pybridge"
	"digital-eval-system/services/go-node/internal/storage"
)

// SubmitService handles evaluation submission flow.
type SubmitService struct {
	pg    *db.PostgresDB
	store storage.Storage
	pyURL *pybridge.Client
	chain interface {
		AppendBlock(*block.Block) (string, error)
	}
	client *http.Client
}

func NewSubmitService(pg *db.PostgresDB, store storage.Storage, pyValidator *pybridge.Client, chain *chain.Chain) *SubmitService {

	if pyValidator == nil {
		pyValidator = pybridge.NewClient("http://127.0.0.1:8082", 120*time.Second)
	}

	return &SubmitService{
		pg:     pg,
		store:  store, // assign interface
		chain:  chain,
		pyURL:  pyValidator,
		client: &http.Client{Timeout: 120 * time.Second},
	}
}

// ValidateAgainstPython calls python validator service.
// ValidateAgainstPython calls python validator service using pybridge client.
func (s *SubmitService) ValidateAgainstPython(ctx context.Context, payload SubmitPayload) (bool, []string, error) {
	if s.pyURL == nil {
		s.pyURL = pybridge.NewClient("http://127.0.0.1:8082", 120*time.Second)
	}

	vIn := pybridge.EvalValidationInput{
		ScriptID:          payload.ScriptID,
		EvaluatorID:       payload.EvaluatorID,
		TotalQuestions:    payload.TotalQuestions,
		MarksPerQuestion:  payload.MarksPerQuestion,
		TotalMarks:        payload.TotalMarks,
		QuestionsAnswered: payload.QuestionsAnswered,
		MarksAllotted:     payload.MarksAllotted,
		MarksScored:       payload.MarksScored,
		CourseID:          payload.CourseID,
		Semester:          payload.Semester,
		CourseCredits:     payload.CourseCredits,
	}

	vResp, err := s.pyURL.ValidateEvaluation(ctx, vIn)
	if err != nil {
		return false, nil, err
	}
	return vResp.Valid, vResp.Errors, nil
}

// SubmitEvaluation runs the full flow: validate, create evaluation block, store in pg.
func (s *SubmitService) SubmitEvaluation(ctx context.Context, payload SubmitPayload) (string, error) {
	// validate basic
	if payload.ScriptID == "" || payload.EvaluatorID == "" {
		return "", fmt.Errorf("missing fields")
	}

	// 1. call python validator (best-effort)
	valid, errors, err := s.ValidateAgainstPython(ctx, payload)
	if err != nil {
		logrus.Warnf("python validator call failed: %v", err)
		// allow proceed but warn
	}
	if !valid {
		return "", fmt.Errorf("validation failed: %v", errors)
	}

	// 2. create marks JSON
	marksStruct := map[string]interface{}{
		"total_questions":    payload.TotalQuestions,
		"marks_per_question": payload.MarksPerQuestion,
		"total_marks":        payload.TotalMarks,
		"course_id":          payload.CourseID,
		"semester":           payload.Semester,
		"academic_year":      payload.AcademicYear,
		"course_credits":     payload.CourseCredits,
		"questions_answered": payload.QuestionsAnswered,
		"marks_allotted":     payload.MarksAllotted,
		"marks_scored":       payload.MarksScored,
		"additional":         payload.AdditionalMetadata,
	}
	marksJSON, _ := json.Marshal(marksStruct)

	// 3. create evaluation transaction and block
	tx := block.Transaction{
		ScriptID:     payload.ScriptID,
		USN:          "", // hidden (we fetch from storage below)
		CourseID:     payload.CourseID,
		Semester:     payload.Semester,
		AcademicYear: payload.AcademicYear,
		CID:          "",
		Meta:         map[string]string{"_evaluation": string(marksJSON)},
		CreatedAt:    time.Now().Unix(),
		SignerID:     payload.EvaluatorID,
	}
	newBlock := block.NewBlock("", []block.Transaction{tx}, payload.EvaluatorID)
	blockHash, err := s.chain.AppendBlock(newBlock)
	if err != nil {
		return "", fmt.Errorf("append block failed: %w", err)
	}

	// 4. compute PASS/FAIL
	// 4. compute PASS/FAIL with Module-based logic (Best of 2)
	// Expecting 10 questions (5 modules * 2 questions)
	if len(payload.MarksScored) != 10 {
		return "", fmt.Errorf("expected 10 questions for module-based evaluation, got %d", len(payload.MarksScored))
	}

	sum := 0
	attemptedCount := 0
	
	// Iterate in pairs (Module 1: Q1,Q2; Module 2: Q3,Q4; etc.)
	for i := 0; i < 10; i += 2 {
		m1 := payload.MarksScored[i]
		m2 := payload.MarksScored[i+1]
		
		// Count attempted (non-zero marks)
		if m1 > 0 {
			attemptedCount++
		}
		if m2 > 0 {
			attemptedCount++
		}

		// Take max of the two for the module score
		moduleScore := m1
		if m2 > m1 {
			moduleScore = m2
		}
		sum += moduleScore
	}

	// Validation: Min questions to be attempted >= 5
	if attemptedCount < 5 {
		return "", fmt.Errorf("minimum 5 questions must be attempted (got %d)", attemptedCount)
	}

	// Validation: Max Marks <= 100
	if sum > 100 {
		return "", fmt.Errorf("total calculated score %d exceeds maximum 100", sum)
	}

	result := "FAIL"
	if payload.TotalMarks > 0 {
		perc := (float64(sum) / float64(payload.TotalMarks)) * 100.0
		if perc >= 36.0 {
			result = "PASS"
		}
	}

	// 5. attempt to find student USN from storage blocks (best-effort)
	studentUSN := ""
	_ = s.store.ForEachBlock(func(blk *block.Block) {
		for _, t := range blk.Transactions {
			if strings.EqualFold(strings.TrimSpace(t.ScriptID), strings.TrimSpace(payload.ScriptID)) {
				if usn, ok := t.Meta["USN"]; ok && usn != "" {
					studentUSN = usn
					return
				}
			}
		}
	})

	// 6. persist to Postgres
	// InsertEvaluationResult signature (Option A) expects:
	// ctx, script_id, student_usn, course_id, semester, course_credits (string), evaluator_id, marks_json (string), total_marks (int), result, block_hash
	err = s.pg.InsertEvaluationResult(ctx, payload.ScriptID, studentUSN, payload.CourseID, payload.Semester, payload.AcademicYear, payload.CourseCredits, payload.EvaluatorID, marksJSON, payload.TotalMarks, result, blockHash)
	if err != nil {
		logrus.Warnf("failed to insert evaluation result: %v", err)
		// continue
	}

	// 7. Update assigned_scripts status to 'evaluated'
	assignedRow, err := s.pg.GetAssignedScriptByScriptID(ctx, payload.ScriptID)
	if err != nil {
		logrus.Warnf("failed to fetch assigned script for status update: %v", err)
	} else if assignedRow != nil {
		if err := s.pg.UpdateAssignmentStatus(ctx, assignedRow.ID, "evaluated"); err != nil {
			logrus.Warnf("failed to update assignment status: %v", err)
		} else {
			logrus.Infof("updated assignment status to evaluated for script %s", payload.ScriptID)
		}
	}

	return blockHash, nil
}
