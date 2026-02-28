"""
FastAPI service exposing:
- POST /extract  : accepts JSON {file_path: "..."} returns extracted metadata + CID
- POST /validate : accepts metadata to validate against schema
"""

import os
from pathlib import Path
from fastapi import FastAPI, HTTPException
from fastapi.responses import JSONResponse
from pydantic import ValidationError
from datetime import datetime, timezone

from schema import ExtractRequest, ExtractResponse, ValidateRequest, ValidateResponse, Metadata
from extractor import process_input_file
from validator import normalize_meta, validate_meta
import uvicorn

# Resolve project root: app/main.py -> python-validator -> services -> digital-eval-system
_PROJECT_ROOT = Path(__file__).resolve().parents[3]
PDF_OUT = str(_PROJECT_ROOT / "documents" / "PDFs")
META_OUT = str(_PROJECT_ROOT / "documents" / "Metadata")

app = FastAPI(title="Python Validator Service", version="1.0.0")

@app.post("/extract", response_model=ExtractResponse)
async def extract(req: ExtractRequest):
    file_path = req.file_path
    if not os.path.exists(file_path):
        raise HTTPException(status_code=400, detail="file not found")
    try:
        out = process_input_file(file_path, pdf_out_dir=PDF_OUT, meta_out_dir=META_OUT, upload_to_ipfs=True)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
    # pydantic response
    meta = normalize_meta(out.get("metadata", {}))
    timestamp = out.get("timestamp", datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"))
    resp = {
        "status": "success",
        "metadata": meta,
        "pdf_cid": out.get("pdf_cid", ""),
        "pdf_path": out.get("pdf_path", ""),
        "timestamp": timestamp
    }
    return JSONResponse(status_code=200, content=resp)

@app.post("/validate", response_model=ValidateResponse)
async def validate(req: ValidateRequest):
    try:
        m = req.metadata.model_dump()
    except ValidationError as e:
        raise HTTPException(status_code=400, detail=str(e))
    norm = normalize_meta(m)
    valid, errors = validate_meta(norm)
    return {"status": "ok" if valid else "error", "valid": valid, "errors": errors}

if __name__ == "__main__":
    uvicorn.run(app, host="127.0.0.1", port=8081)