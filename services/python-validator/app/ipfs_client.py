"""
ipfs_client.py
Simple IPFS HTTP add helper. Uses IPFS_API env var or default http://127.0.0.1:5001/api/v0
"""

import os
from pathlib import Path
from typing import Optional

import requests

IPFS_API = os.getenv("IPFS_API", "http://127.0.0.1:5001/api/v0")
IPFS_TIMEOUT = int(os.getenv("IPFS_TIMEOUT_SECONDS", "120"))

def add_file(filepath: str, timeout: Optional[int] = None) -> str:
    """
    Upload file at filepath to IPFS via HTTP API `/add`.
    Returns CID string on success. Raises RuntimeError on failure.
    """
    url = IPFS_API.rstrip("/") + "/add"
    to = timeout or IPFS_TIMEOUT
    with open(filepath, "rb") as f:
        files = {"file": (Path(filepath).name, f)}
        resp = requests.post(url, files=files, timeout=to)
    resp.raise_for_status()
    # ipfs returns streaming JSON lines in some setups; parse robustly
    text = resp.text.strip()
    # try direct json
    try:
        j = resp.json()
        cid = j.get("Hash") or j.get("hash")
        if cid:
            return cid
    except Exception:
        pass
    # fallback: scan lines for JSON with Hash
    for line in text.splitlines()[::-1]:
        line = line.strip()
        if not line:
            continue
        try:
            obj = __import__("json").loads(line)
            if isinstance(obj, dict):
                cid = obj.get("Hash") or obj.get("hash")
                if isinstance(cid, str) and cid.strip():
                    return cid

        except Exception:
            continue
    raise RuntimeError("IPFS add returned no CID")
