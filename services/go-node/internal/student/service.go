package student

import (
	"context"
	"encoding/json"
	"fmt"
	"math"

	"digital-eval-system/services/go-node/internal/db"
	"digital-eval-system/services/go-node/internal/rootdir"
)

// Service provides student result access
type Service struct {
	pg *db.PostgresDB
}

func NewService(pg *db.PostgresDB) *Service {
	return &Service{pg: pg}
}

func (s *Service) FetchResults(ctx context.Context, usn, semester string, academicYear string) ([]db.EvaluationRow, error) {
	// Your DB helper already filters by USN, so we fetch all rows
	rows, err := s.pg.FetchResultsByUSN(ctx, usn, academicYear)
	if err != nil {
		return nil, err
	}

	// Filter by semester here for safety
	// If semester is empty, return all rows (useful for CGPA calculation)
	var filtered []db.EvaluationRow
	for _, r := range rows {
		if semester == "" || r.Semester == semester {
			filtered = append(filtered, r)
		}
	}

	return filtered, nil
}

// ---------- PHASE C: FetchResultsWithGPA ----------
func (s *Service) FetchResultsWithGPA(ctx context.Context, usn, semester string, academicYear string) (map[string]interface{}, error) {

	rows, err := s.FetchResults(ctx, usn, semester, academicYear)
	if err != nil {
		return nil, err
	}

	sgpa := CalculateSGPA(rows)

	return map[string]interface{}{
		"usn":           usn,
		"semester":      semester,
		"academic_year": academicYear,
		"rows":          rows,
		"sgpa":          sgpa,
	}, nil
}

// GenerateResultPDF fetches results and produces the PDF bytes using pdf_generator.go
func (s *Service) GenerateResultPDF(ctx context.Context, usn, semester string, academicYear string) ([]byte, error) {
	// Fix: Use FetchResults directly to get []db.EvaluationRow, not the map from FetchResultsWithGPA
	rows, err := s.FetchResults(ctx, usn, semester, academicYear)
	if err != nil {
		return nil, fmt.Errorf("fetch results: %w", err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("no results found for %s semester %s", usn, semester)
	}

	pdfBytes, err := GenerateResultPDF(ctx, usn, semester, academicYear, rows, PDFOptions{
		LogoPath: rootdir.Resolve("services/go-node/internal/student/assets/biet_logo.jpg"),
		FontDir:  rootdir.Resolve("services/go-node/internal/student/assets/fonts"),
	})
	if err != nil {
		return nil, fmt.Errorf("generate PDF: %w", err)
	}

	return pdfBytes, nil
}

func CalculateSGPA(rows []db.EvaluationRow) float64 {
	if len(rows) == 0 {
		return 0.0
	}
	var totalWeightedPoints float64
	var totalCredits float64

	for _, r := range rows {
		// r.TotalMarks is int, r.Result string. We need marks scored from r.Marks JSON.
		// Expect r.Marks JSON to contain "marks_scored" array or "marks_scored_total" numeric.
		// We'll try to extract marks_scored array sum; fallback to 0.
		var marksMap map[string]interface{}
		_ = json.Unmarshal(r.Marks, &marksMap)

		// sum marks_scored
		// sum marks_scored using module logic
		marksScored := calculateModuleScore(marksMap["marks_scored"])

		// If total marks available, compute percentage for grade points.
		credit := float64(0)
		if r.CourseCredits.Valid {
			credit = float64(r.CourseCredits.Int32) // adjust type if CourseCredits is different
		} else {
			// fallback: try CourseCredits column as int field if present
			if r.CourseCredits.Int32 == 0 {
				credit = 0
			}
		}

		if credit <= 0 {
			continue
		}

		totalCredits += credit

		// compute percentage and map to grade point (simple 10-point scale).
		var perc float64
		if r.TotalMarks > 0 {
			perc = (float64(marksScored) / float64(r.TotalMarks)) * 100.0
		} else {
			perc = 0
		}

		// Convert percentage to grade point (0..10). Adjust thresholds as required.
		var gp float64
		switch {
		case perc >= 90:
			gp = 10.0
		case perc >= 80:
			gp = 9.0
		case perc >= 70:
			gp = 8.0
		case perc >= 60:
			gp = 7.0
		case perc >= 50:
			gp = 6.0
		case perc >= 40:
			gp = 5.0
		default:
			gp = 0.0
		}

		totalWeightedPoints += gp * credit
	}

	if totalCredits == 0 {
		return 0.0
	}
	sgpa := totalWeightedPoints / totalCredits
	// round to 2 decimals
	return math.Round(sgpa*100) / 100
}
