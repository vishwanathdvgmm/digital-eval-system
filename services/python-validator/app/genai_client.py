"""
genai_client.py
Wraps google.genai client usage with retries and content assembly.
Uses GENAI_API_KEY and GENAI_MODEL from env by default.
"""

import os
import time
from typing import List, Tuple

from google import genai
from google.genai import types

# environment-driven configuration
API_KEY = "Your API Key."
MODEL = "gemini-2.5-flash"
RETRIES = 2
BACKOFF_BASE = 1.0

if not API_KEY:
    raise RuntimeError("GENAI_API_KEY environment variable is required for genai_client")

_client = genai.Client(api_key=API_KEY)

def generate_with_retries(parts: List[Tuple[str, bytes]], prompt: str, retries: int = RETRIES) -> str:
    """
    parts: list of (mime_type, bytes)
    prompt: text prompt appended after parts
    returns: response text (string)
    """
    last_err = None
    for attempt in range(1, retries + 1):
        try:
            contents = []
            for mime, b in parts:
                contents.append(types.Part.from_bytes(data=b, mime_type=mime))
            contents.append(prompt)
            resp = _client.models.generate_content(model=MODEL, contents=contents)
            # resp.text is the textual content; resp may be None depending on API.
            return resp.text or ""
        except Exception as e:
            last_err = e
            # exponential-ish backoff
            time.sleep(BACKOFF_BASE * attempt)
    # final fallback: raise the last exception string for caller to handle
    raise RuntimeError(f"genai generation failed after {retries} attempts: {last_err}")
