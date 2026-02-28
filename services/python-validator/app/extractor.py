"""
extractor.py
Refactor of user's extraction code:
- image deskew, enhance, crops
- integrate genai_client.generate_with_retries
- PDF first-page rasterization using fitz (PyMuPDF)
- save PDF fullpage
- returns canonical metadata dict
"""

import os
import re
import tempfile
import shutil
from io import BytesIO
from datetime import datetime, timezone
from pathlib import Path
from typing import Dict

import cv2
import numpy as np
from PIL import Image
from reportlab.pdfgen import canvas
from reportlab.lib.pagesizes import A4

from genai_client import generate_with_retries
from ipfs_client import add_file

# Regex patterns
USN_REGEX = r"[0-9][A-Z]{2}[0-9]{2}[A-Z]{2}[0-9]{3}"
CID_REGEX = r"[A-Z]{3}[0-9]{3}"
SEM_REGEX = r"[1-8]"
DATE_REGEX = r"\d{1,2}[-/]\d{1,2}[-/]\d{4}"

# prompts (kept from original)
META_PROMPT = """
You are an OMR metadata extractor.

Return JSON ONLY in this exact structure:
{
  "USN": "",
  "CourseID": "",
  "Semester": "",
  "CourseName": "",
  "Date": "",
  "Institute": ""
}

Rules:
- USN pattern: 1 digit + 2 letters + 2 digits + 2 letters + 3 digits
- CourseID: 3 letters + 3 digits + maybe 1 letter
- Semester: 1 digit
- Only output valid JSON.
"""

COURSE_RECOVER_PROMPT = """
Extract ONLY the course name. Return JSON:
{"CourseName": "..."}
"""

# helper functions
def clean_json_string(s: str) -> str:
    if not isinstance(s, str):
        return ""
    s = s.replace("\n", " ").replace("\r", " ")
    return "".join(ch for ch in s if ord(ch) >= 32).strip()

def safe_filename(name: str) -> str:
    s = re.sub(r'[^A-Za-z0-9._-]', '_', str(name))
    return s[:120]

# image processing
def deskew_image(img):
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    edges = cv2.Canny(gray, 40, 140)
    lines = cv2.HoughLines(edges, 1, np.pi/180, 150)
    if lines is None:
        return img
    angles = []
    for l in lines:
        rho, theta = l[0]
        ang = (theta * 180 / np.pi) - 90
        if -45 < ang < 45:
            angles.append(ang)
    if not angles:
        return img
    angle = float(np.median(angles))
    h, w = img.shape[:2]
    M = cv2.getRotationMatrix2D((w/2, h/2), angle, 1.0)
    return cv2.warpAffine(img, M, (w, h), flags=cv2.INTER_LINEAR, borderMode=cv2.BORDER_REPLICATE)

def enhance_image(img):
    gray = cv2.cvtColor(img, cv2.COLOR_BGR2GRAY)
    clahe = cv2.createCLAHE(clipLimit=2.0, tileGridSize=(8,8))
    enhanced = clahe.apply(gray)
    kernel = np.array([[0,-1,0],[-1,5,-1],[0,-1,0]])
    sharp = cv2.filter2D(enhanced, -1, kernel)
    return cv2.cvtColor(sharp, cv2.COLOR_GRAY2BGR)

def crop_header_regions(img):
    h, w = img.shape[:2]
    r1 = img[0:int(h*0.18), 0:w]
    r2 = img[int(h*0.08):int(h*0.28), 0:w]
    r3 = img[0:int(h*0.22), 0:int(w*0.60)]
    return [r1, r2, r3]

def image_to_bytes(img, fmt="PNG"):
    rgb = cv2.cvtColor(img, cv2.COLOR_BGR2RGB)
    pil = Image.fromarray(rgb)
    bio = BytesIO()
    pil.save(bio, format=fmt)
    return bio.getvalue()

# PDF first page to image
def pdf_first_page_to_image(pdf_path: str, out_png: str) -> str:
    import fitz  # PyMuPDF
    doc = fitz.open(pdf_path)
    page = doc.load_page(0)
    pix = page.get_pixmap(dpi=200)
    pix.save(out_png)
    return out_png

def save_pdf_fullpage(image_path: str, usn: str, course: str, dest_dir: str) -> str:
    course = safe_filename(course or "UNKNOWN")
    stamp = datetime.now(timezone.utc).strftime("%Y%m%d_%H%M%S")
    pdf_name = f"{usn}_{course}_{stamp}.pdf"
    pdf_path = os.path.join(dest_dir, pdf_name)

    img = Image.open(image_path)
    iw, ih = img.size
    pw, ph = A4

    img_aspect = iw / ih
    page_aspect = pw / ph

    if img_aspect > page_aspect:
        new_w = pw
        new_h = pw / img_aspect
        x = 0
        y = (ph - new_h) / 2
    else:
        new_h = ph
        new_w = ph * img_aspect
        x = (pw - new_w) / 2
        y = 0

    c = canvas.Canvas(pdf_path, pagesize=A4)
    c.drawImage(image_path, x, y, new_w, new_h)
    c.save()
    return pdf_path

# high-level extraction
def extract_metadata_from_image(image_path: str) -> Dict:
    orig = cv2.imread(image_path)
    if orig is None:
        raise ValueError(f"cannot read image {image_path}")

    desk = deskew_image(orig)
    enhanced = enhance_image(desk)
    crops = crop_header_regions(enhanced)
    parts = [("image/png", image_to_bytes(c)) for c in crops]

    # full page JPEG for robust OCR/LLM
    pil_full = Image.fromarray(cv2.cvtColor(enhanced, cv2.COLOR_BGR2RGB))
    buf = BytesIO()
    pil_full.save(buf, "JPEG", quality=80)
    full_bytes = buf.getvalue()
    parts.insert(0, ("image/jpeg", full_bytes))

    raw = generate_with_retries(parts, META_PROMPT)
    raw = (raw or "").strip()

    if raw.startswith("```"):
        # strip code fences
        raw = raw.split("```")[1].strip()
        if raw.lower().startswith("json"):
            raw = raw[4:].strip()

    # safe json extraction: find first {...}
    import json, re as _re
    def safe_json_extract(txt: str):
        try:
            return json.loads(txt)
        except Exception:
            m = _re.search(r"\{(?:[^{}]|\n|(?R))*\}", txt, flags=_re.S)
            if m:
                try:
                    return json.loads(m.group(0))
                except Exception:
                    return {}
            return {}

    meta = safe_json_extract(raw)
    if not isinstance(meta, dict):
        meta = {}

    # normalize string values
    meta = {k: clean_json_string(str(v).upper()) for k, v in meta.items()}

    raw_up = clean_json_string(raw.upper())

    if not re.match(USN_REGEX, meta.get("USN", "")):
        found = re.findall(USN_REGEX, raw_up)
        meta["USN"] = found[0] if found else "UNKNOWN_USN"

    if not re.match(CID_REGEX, meta.get("CourseID", "")):
        found = re.findall(CID_REGEX, raw_up)
        meta["CourseID"] = found[0] if found else "UNKNOWN_COURSE"

    if not re.match(SEM_REGEX, meta.get("Semester", "")):
        found = re.findall(SEM_REGEX, raw_up)
        meta["Semester"] = found[0] if found else "N/A"

    if not meta.get("CourseName") or meta["CourseName"] in ["", "N/A"]:
        # recover course name with dedicated prompt
        try:
            course_raw = generate_with_retries([("image/jpeg", full_bytes)], COURSE_RECOVER_PROMPT)
            import json as _json
            parsed = {}
            try:
                parsed = _json.loads(course_raw)
            except Exception:
                # try to find JSON in text
                m = _re.search(r"\{(?:[^{}]|\n|(?R))*\}", course_raw, flags=_re.S)
                if m:
                    try:
                        parsed = _json.loads(m.group(0))
                    except Exception:
                        parsed = {}
            meta["CourseName"] = clean_json_string(parsed.get("CourseName", "N/A")).upper()
        except Exception:
            meta["CourseName"] = "N/A"

    if not re.match(DATE_REGEX, meta.get("Date", "")):
        found = re.findall(DATE_REGEX, raw_up)
        meta["Date"] = found[0] if found else ""

    meta.setdefault("Institute", "")

    return meta

def process_input_file(input_path: str, pdf_out_dir: str, meta_out_dir: str, upload_to_ipfs: bool = True) -> dict:
    """
    - If input is PDF, rasterize first page to image and extract.
    - Save a full-page PDF copy (A4-fit) to pdf_out_dir.
    - Save metadata JSON to meta_out_dir.
    - Optionally upload PDF to IPFS and return CID.
    Returns dict with metadata, pdf_path, pdf_cid, timestamp.
    """
    tmpdir = tempfile.mkdtemp(prefix="val_")
    try:
        p = Path(input_path)
        if not p.exists():
            raise FileNotFoundError(input_path)

        if p.suffix.lower() == ".pdf":
            img_path = os.path.join(tmpdir, "first.png")
            pdf_first_page_to_image(str(p), img_path)
            meta = extract_metadata_from_image(img_path)
            usn = safe_filename(meta.get("USN", "UNKNOWN_USN"))
            course = meta.get("CourseID", "UNKNOWN")
            # create saved pdf copy (copy original but ensure naming)
            pdf_path = save_pdf_fullpage(img_path, usn, course, pdf_out_dir)
            # copy original content into the new pdf (preserve original multipage)
            shutil.copy2(str(p), pdf_path)
        else:
            # assume image
            meta = extract_metadata_from_image(str(p))
            usn = safe_filename(meta.get("USN", "UNKNOWN_USN"))
            course = meta.get("CourseID", "UNKNOWN")
            img_path = str(p)
            pdf_path = save_pdf_fullpage(img_path, usn, course, pdf_out_dir)

        # save metadata JSON
        import json
        stamp = datetime.now(timezone.utc).strftime("%Y%m%d_%H%M%S")
        json_name = f"{usn}_{safe_filename(course)}_{stamp}.json"
        json_path = os.path.join(meta_out_dir, json_name)
        os.makedirs(meta_out_dir, exist_ok=True)
        with open(json_path, "w", encoding="utf-8") as fh:
            fh.write(json.dumps(meta, indent=2))

        cid = ""
        if upload_to_ipfs:
            # upload the saved pdf copy to IPFS
            cid = add_file(pdf_path)

        out = {
            "student_id": meta.get("USN", ""),
            "exam_id": meta.get("CourseID", ""),
            "course_name": meta.get("CourseName", ""),
            "semester": meta.get("Semester", ""),
            "timestamp": datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ"),
            "pdf_path": str(Path(pdf_path).resolve()),
            "pdf_cid": cid,
            "metadata": meta,
            "status": "success"
        }
        return out
    finally:
        try:
            shutil.rmtree(tmpdir)
        except Exception:
            pass
