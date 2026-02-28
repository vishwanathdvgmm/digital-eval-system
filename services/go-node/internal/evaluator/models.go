package evaluator

type SubmitPayload struct {
	ScriptID           string                 `json:"script_id"`
	EvaluatorID        string                 `json:"evaluator_id"`
	TotalQuestions     int                    `json:"total_questions"`
	MarksPerQuestion   int                    `json:"marks_per_question"`
	TotalMarks         int                    `json:"total_marks"`
	QuestionsAnswered  int                    `json:"questions_answered"`
	MarksAllotted      []int                  `json:"marks_allotted_per_question"`
	MarksScored        []int                  `json:"marks_scored"`
	CourseID           string                 `json:"course_id"`
	Semester           string                 `json:"semester,omitempty"`
	AcademicYear       string                 `json:"academic_year,omitempty"`
	CourseCredits      int                    `json:"course_credits"`
	AdditionalMetadata map[string]interface{} `json:"additional_metadata"`
}

// you may use []byte or json.RawMessage depending on your code
type RequestCreate struct {
	EvaluatorID  string `json:"evaluator_id"`
	CourseID     string `json:"course_id"`
	Semester     string `json:"semester"`
	AcademicYear string `json:"academic_year"`
	Description  string `json:"description,omitempty"`
}
