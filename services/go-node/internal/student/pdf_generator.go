package student

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"

	"digital-eval-system/services/go-node/internal/db"
	"digital-eval-system/services/go-node/internal/rootdir"
)

type PDFOptions struct {
	LogoPath   string // e.g. services/go-node/internal/student/assets/biet_logo.jpg
	FontDir    string // e.g. services/go-node/internal/student/assets/fonts
	IncludeSig bool
}

// GenerateResultPDF builds the BIET-style result PDF.
// It will fetch CourseName from (A) evaluation marks JSON, or (B) by scanning
// extractor metadata JSON files on disk (common places), preferring metadata.
func GenerateResultPDF(ctx context.Context, usn string, semester string, academicYear string, rows []db.EvaluationRow, opts PDFOptions) ([]byte, error) {
	if len(rows) == 0 {
		return nil, fmt.Errorf("no results for %s semester %s", usn, semester)
	}

	// ---------------------------------------------------------------------
	// Font / resource paths
	// ---------------------------------------------------------------------
	roboto := filepath.Join(opts.FontDir, "Roboto-Regular.ttf")
	robotoB := filepath.Join(opts.FontDir, "Roboto-Bold.ttf")
	kannada := filepath.Join(opts.FontDir, "NotoSansKannada-Regular.ttf")

	// ---------------------------------------------------------------------
	// Setup PDF
	// ---------------------------------------------------------------------
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 18)
	pdf.AddPage()

	// Add fonts (safe adds; if files missing gofpdf will fallback)
	pdf.AddUTF8Font("Rob", "", roboto)
	pdf.AddUTF8Font("RobB", "", robotoB)
	pdf.AddUTF8Font("Kan", "", kannada)

	// ---------------------------------------------------------------------
	// Header
	// ---------------------------------------------------------------------
	if opts.LogoPath != "" {
		// place logo left, keep height ~26mm
		pdf.Image(opts.LogoPath, 15, 12, 26, 0, false, "", 0, "")
	}

	pdf.SetFont("RobB", "", 14)
	pdf.SetXY(15, 12)
	pdf.CellFormat(180, 10, "BAPUJI INSTITUTE OF ENGINEERING & TECHNOLOGY", "", 1, "C", false, 0, "")

	pdf.SetFont("Rob", "", 12)
	pdf.CellFormat(180, 6, "DAVANAGERE - 577004", "", 1, "C", false, 0, "")

	// Kannada motto (if font available)
	pdf.SetFont("Kan", "", 11)
	pdf.CellFormat(180, 7, "ಕರ್ಮಣೇಯೇವಾಧಿಕಾರಸ್ತೇ ಮಾಫಲೇಷು ಕದಾಚನ", "", 1, "C", false, 0, "")
	pdf.Ln(6)

	// ---------------------------------------------------------------------
	// Title & student info box
	// ---------------------------------------------------------------------
	pdf.SetFont("RobB", "", 12)
	pdf.CellFormat(0, 9, "STUDENT RESULT REPORT", "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Student info - use first row's StudentUSN (DB stores it properly)
	infoRow := rows[0]
	studentUSN := infoRow.StudentUSN.String

	pdf.SetFont("Rob", "", 11)
	pdf.SetFillColor(248, 248, 248)

	pdf.CellFormat(95, 8, fmt.Sprintf("USN: %s", studentUSN), "1", 0, "L", true, 0, "")
	pdf.CellFormat(95, 8, fmt.Sprintf("Semester: %s", semester), "1", 1, "L", true, 0, "")

	pdf.CellFormat(95, 8, "Institute: BIET Davangere", "1", 0, "L", true, 0, "")
	pdf.CellFormat(95, 8, fmt.Sprintf("Exam Date: %s", time.Now().Format("02-01-2006")), "1", 1, "L", true, 0, "")
	pdf.CellFormat(190, 8, fmt.Sprintf("Academic Year: %s", academicYear), "1", 1, "L", true, 0, "")

	pdf.Ln(12)

	// ---------------------------------------------------------------------
	// Table header
	// ---------------------------------------------------------------------
	pdf.SetFont("RobB", "", 12)
	pdf.SetFillColor(230, 230, 230)

	pdf.CellFormat(30, 8, "Course ID", "1", 0, "C", true, 0, "")
	pdf.CellFormat(85, 8, "Course Name", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, "Marks", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, "Total", "1", 0, "C", true, 0, "")
	pdf.CellFormat(25, 8, "Result", "1", 1, "C", true, 0, "")

	// ---------------------------------------------------------------------
	// Table rows
	// ---------------------------------------------------------------------
	pdf.SetFont("Rob", "", 11)

	totalScoredAll := 0
	totalMarksAll := 0

	for _, r := range rows {
		// parse marks JSON (safely)
		var marksMap map[string]interface{}
		_ = json.Unmarshal(r.Marks, &marksMap)

		// marks_scored array -> sum
		scored := calculateModuleScore(marksMap["marks_scored"])
		totalScoredAll += scored
		totalMarksAll += r.TotalMarks

		// obtain course name:
		// 1) try marks JSON course_name
		// 2) otherwise try metadata files on disk (extractor output)
		courseName := ""
		if cn, ok := marksMap["course_name"].(string); ok && strings.TrimSpace(cn) != "" {
			courseName = cn
		} else {
			// attempt to locate metadata JSON that matches USN and CourseID
			courseName = findCourseNameFromMetadata(studentUSN, r.CourseID)
		}

		// safe formatting
		if courseName == "" {
			courseName = "-"
		}

		pdf.CellFormat(30, 8, r.CourseID, "1", 0, "C", false, 0, "")
		pdf.CellFormat(85, 8, courseName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 8, fmt.Sprintf("%d", scored), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 8, fmt.Sprintf("%d", r.TotalMarks), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 8, r.Result, "1", 1, "C", false, 0, "")
	}

	pdf.Ln(8)

	// ---------------------------------------------------------------------
	// Summary block
	// ---------------------------------------------------------------------
	pdf.SetFont("RobB", "", 12)
	label := "Total Marks Scored:"
	labelWidth := pdf.GetStringWidth(label) + 2

	pdf.CellFormat(labelWidth, 8, label, "", 0, "L", false, 0, "")

	pdf.SetFont("Rob", "", 12)
	pdf.CellFormat(0, 8, fmt.Sprintf("%d/%d", totalScoredAll, totalMarksAll), "", 1, "L", false, 0, "")

	pdf.Ln(6)

	sgpa := CalculateSGPA(rows)
	renderTightLine(pdf, "SGPA (this semester):", fmt.Sprintf("%.2f", sgpa))

	pdf.Ln(6)

	// ---------------------------------------------------------------------
	// Output bytes
	// ---------------------------------------------------------------------
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// calculateModuleScore calculates the total score based on Best of 2 per module logic.
// It expects an array of 10 marks (5 modules * 2 questions).
func calculateModuleScore(v interface{}) int {
	if v == nil {
		return 0
	}

	// Convert interface to []int
	var marks []int
	switch arr := v.(type) {
	case []interface{}:
		for _, it := range arr {
			switch n := it.(type) {
			case float64:
				marks = append(marks, int(n))
			case int:
				marks = append(marks, n)
			case int64:
				marks = append(marks, int(n))
			}
		}
	}

	// Calculate score
	sum := 0
	// Iterate in pairs
	for i := 0; i < len(marks); i += 2 {
		m1 := marks[i]
		m2 := 0
		if i+1 < len(marks) {
			m2 = marks[i+1]
		}

		// Max of the pair
		if m1 > m2 {
			sum += m1
		} else {
			sum += m2
		}
	}
	return sum
}

func renderTightLine(pdf *gofpdf.Fpdf, label, value string) {
	pdf.SetFont("RobB", "", 12)
	labelWidth := pdf.GetStringWidth(label) + 2
	pdf.CellFormat(labelWidth, 8, label, "", 0, "L", false, 0, "")

	pdf.SetFont("Rob", "", 12)
	pdf.CellFormat(0, 8, value, "", 1, "L", false, 0, "")
}

// safe helper to find course name from metadata JSON files written by extractor.
// It searches a list of likely directories and matches JSON files containing both USN and CourseID.
// Returns empty string if not found.
func findCourseNameFromMetadata(usn, courseID string) string {
	// Candidate metadata directories (ordered)
	candidates := []string{
		rootdir.Resolve("documents/Metadata"),
	}

	usn = strings.TrimSpace(strings.ToUpper(usn))
	courseID = strings.TrimSpace(strings.ToUpper(courseID))

	for _, dir := range candidates {
		_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			// ignore errors traversing (we just return nil)
			if err != nil || d == nil || d.IsDir() {
				return nil
			}
			// only JSON files
			if !strings.HasSuffix(strings.ToLower(d.Name()), ".json") {
				return nil
			}
			// quick filename filter: must contain usn or courseID or both
			filename := strings.ToUpper(d.Name())
			if !strings.Contains(filename, usn) && !strings.Contains(filename, courseID) {
				return nil
			}
			// read file
			bs, rerr := os.ReadFile(path)
			if rerr != nil {
				return nil
			}
			var jm map[string]interface{}
			if jerr := json.Unmarshal(bs, &jm); jerr != nil {
				return nil
			}
			// check USN and CourseID match if present
			jusn := strings.ToUpper(strings.TrimSpace(interfaceToString(jm["USN"])))
			jcid := strings.ToUpper(strings.TrimSpace(interfaceToString(jm["CourseID"])))
			if jusn != "" && jcid != "" {
				if jusn == usn && (courseID == "" || jcid == courseID) {
					// return CourseName if present
					cn := strings.TrimSpace(interfaceToString(jm["CourseName"]))
					if cn != "" {
						// found — set result by returning a sentinel error to stop WalkDir
						courseNameFound = cn
						// Use panic/recover free method: set global var then stop by returning fs.SkipDir
						return fs.SkipDir
					}
				}
			} else {
				// more lenient: if CourseID matches and CourseName exists
				if jcid == courseID {
					cn := strings.TrimSpace(interfaceToString(jm["CourseName"]))
					if cn != "" {
						courseNameFound = cn
						return fs.SkipDir
					}
				}
			}
			return nil
		})
		if courseNameFound != "" {
			// copy and reset sentinel
			res := courseNameFound
			courseNameFound = ""
			return res
		}
	}
	return ""
}

// small sentinel for findCourseNameFromMetadata to capture match within filepath.WalkDir closure
var courseNameFound string

func interfaceToString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	case fmt.Stringer:
		return t.String()
	default:
		bs, _ := json.Marshal(t)
		return string(bs)
	}
}
