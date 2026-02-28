#!/usr/bin/env python3
# FastAPI service that validates evaluation payloads.
# Checks consistency (lengths, totals) and returns {"valid": true/false, "errors": [...]}

from fastapi import FastAPI
from pydantic import BaseModel, conlist
from typing import Optional
import uvicorn

app = FastAPI(title="Evaluation Validator")

class EvalPayload(BaseModel):
    script_id: str
    evaluator_id: str
    total_questions: int
    marks_per_question: int
    total_marks: int
    questions_answered: int
    marks_allotted_per_question: conlist(int, min_length=0)
    marks_scored: conlist(int, min_length=0)
    additional_metadata: Optional[dict] = None

@app.post("/validate-evaluation")
def validate_evaluation(p: EvalPayload):
    errors = []

    if p.total_questions <= 0:
        errors.append("total_questions must be > 0")
    if p.marks_per_question <= 0:
        errors.append("marks_per_question must be > 0")
    if p.total_marks > p.total_questions * p.marks_per_question:
        errors.append("total_marks exceeds theoretical maximum (total_questions * marks_per_question)")
    if p.questions_answered < 0 or p.questions_answered > p.total_questions:
        errors.append("questions_answered out of range")
    if len(p.marks_allotted_per_question) not in (0, p.questions_answered, p.total_questions):
        # allow 0 (not provided) or one entry per answered question or per total question
        errors.append("marks_allotted_per_question length invalid")
    if len(p.marks_scored) not in (0, p.questions_answered, p.total_questions):
        errors.append("marks_scored length invalid")

    # ensure scores do not exceed allotted/marks_per_question
    for s in p.marks_scored:
        if s < 0 or s > p.marks_per_question:
            errors.append("marks_scored contains invalid value")

    valid = len(errors) == 0
    return {"valid": valid, "errors": errors}

if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=8082)
