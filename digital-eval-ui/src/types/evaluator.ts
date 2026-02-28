export interface EvaluationRequestCreate {
    evaluator_id: string;
    course_id: string;
    semester: string;
    academic_year: string;
    description?: string;
}

export interface EvaluationRequest {
    id: number;
    course_id: string;
    semester: string;
    academic_year: string;
    description: string;
    status: string;
    created_at: string;
}

export interface AssignedScript {
    ID: number;
    ScriptID: string;
    Evaluator: string;
    CourseID: string;
    Semester: string;
    AcademicYear: string;
    CourseCredits: number;
    AssignedAt: string;
    Status: string;
}

export interface ScriptMetadata {
    [key: string]: string;
}

export interface EvaluationSubmit {
    script_id: string;
    evaluator_id: string;
    total_questions: number;
    marks_per_question: number;
    total_marks: number;
    questions_answered: number;
    marks_allotted_per_question: number[];
    marks_scored: number[];
    course_id: string;
    semester: string;
    academic_year: string;
    course_credits: number;
    additional_metadata: Record<string, any>;
}

export interface SubmitResponse {
    block_hash: string;
}
