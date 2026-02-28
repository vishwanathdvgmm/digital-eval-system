"""
schema.py
Pydantic models and JSON schema for metadata and API payloads/responses.
"""

from typing import Dict, Optional
from pydantic import BaseModel, Field

class Metadata(BaseModel):
    USN: str = Field(..., description="University Student Number")
    CourseID: str = Field(..., description="Course identifier")
    Semester: str = Field(..., description="Semester number or code")
    CourseName: Optional[str] = Field(None, description="Course name (optional)")
    Date: Optional[str] = Field(None, description="Exam or script date")
    Institute: Optional[str] = Field(None, description="Institute name")
    # allow extra fields as meta
    extra: Optional[Dict[str, str]] = None

class ExtractRequest(BaseModel):
    # either local file path or base64 bytes can be supported later; for now local path
    file_path: str

class ExtractResponse(BaseModel):
    status: str
    metadata: Metadata
    pdf_cid: str
    pdf_path: str
    timestamp: str

class ValidateRequest(BaseModel):
    metadata: Metadata
    expected_usn: Optional[str] = None
    expected_courseid: Optional[str] = None

class ValidateResponse(BaseModel):
    status: str
    valid: bool
    errors: Optional[Dict[str, str]] = None
