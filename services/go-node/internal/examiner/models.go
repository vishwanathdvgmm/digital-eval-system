package examiner

import "time"

// ScriptUpload is a compact record stored as a block (and optionally indexed later).
type ScriptUpload struct {
	ScriptID   string            `json:"script_id"`
	USN        string            `json:"usn"`
	CourseID   string            `json:"course_id"`
	Semester   string            `json:"semester"`
	CourseName string            `json:"course_name,omitempty"`
	PDFCid     string            `json:"pdf_cid"`
	PDFPath    string            `json:"pdf_path"`
	Metadata   map[string]string `json:"metadata"`
	Status     string            `json:"status"` // uploaded | validated | errored
	CreatedAt  time.Time         `json:"created_at"`
}
