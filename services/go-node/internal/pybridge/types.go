package pybridge

// Types that mirror Python validator service JSON payloads.
// Keep these structs simple and tolerant to extra fields.

type ExtractRequest struct {
	FilePath string `json:"file_path"`
}

type ExtractResponse struct {
	Status    string            `json:"status"`
	Metadata  map[string]string `json:"metadata"`
	PDFCid    string            `json:"pdf_cid"`
	PDFPath   string            `json:"pdf_path"`
	Timestamp string            `json:"timestamp"`
	Error     string            `json:"error,omitempty"`
}

type ValidateRequest struct {
	Metadata       map[string]string `json:"metadata"`
	ExpectedUSN    string            `json:"expected_usn,omitempty"`
	ExpectedCourse string            `json:"expected_courseid,omitempty"`
}

type ValidateResponse struct {
	Status string   `json:"status"`
	Valid  bool     `json:"valid"`
	Errors []string `json:"errors,omitempty"`
	Error  string   `json:"error,omitempty"`
}
