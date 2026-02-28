"""
validator.py
Validation rules, canonicalization and metadata signing placeholder.
"""

import re
from typing import Tuple

USN_REGEX = re.compile(r"[0-9][A-Z]{2}[0-9]{2}[A-Z]{2}[0-9]{3}")
COURSEID_REGEX = re.compile(r"[A-Z]{3}[0-9]{3}[A-Z]?")
SEM_REGEX = re.compile(r"^[1-8]$")

def normalize_meta(raw: dict) -> dict:
    """
    Convert keys to canonical shapes, uppercase some fields, trim.
    Does not mutate input.
    """
    out = {}
    for k, v in raw.items():
        if v is None:
            continue
        key = k.strip()
        val = str(v).strip()
        if key.lower() in ("usn", "student_id"):
            out["USN"] = val.upper()
        elif key.lower() in ("courseid", "course_id", "course"):
            out["CourseID"] = val.upper()
        elif key.lower() in ("semester", "sem"):
            out["Semester"] = val.upper()
        elif key.lower() in ("coursename", "course_name"):
            out["CourseName"] = val.strip()
        elif key.lower() in ("date",):
            out["Date"] = val
        elif key.lower() in ("institute", "college"):
            out["Institute"] = val
        else:
            # keep other keys under meta.extra
            out.setdefault("extra", {})[key] = val
    return out

def validate_meta(meta: dict) -> Tuple[bool, dict]:
    """
    Validate presence of required fields and format.
    Returns (valid, errors)
    """
    errors = {}
    usn = meta.get("USN", "")
    course = meta.get("CourseID", "")
    sem = meta.get("Semester", "")

    if not usn:
        errors["USN"] = "missing"
    elif not USN_REGEX.search(usn):
        errors["USN"] = "invalid_format"

    if not course:
        errors["CourseID"] = "missing"
    elif not COURSEID_REGEX.search(course):
        errors["CourseID"] = "invalid_format"

    if not sem:
        errors["Semester"] = "missing"
    elif not SEM_REGEX.search(sem):
        errors["Semester"] = "invalid_format"

    return (len(errors) == 0), errors
