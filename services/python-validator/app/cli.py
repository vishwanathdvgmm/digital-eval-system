"""
CLI wrapper to call extractor in batch or single mode.
Usage:
  python -m services.python-validator.app.cli <file-or-dir>
"""

import argparse
import os
import sys
from pathlib import Path
import json

from extractor import process_input_file

# Resolve project root: app/cli.py -> python-validator -> services -> digital-eval-system
_PROJECT_ROOT = Path(__file__).resolve().parents[3]

def main(argv=None):
    parser = argparse.ArgumentParser()
    parser.add_argument("input", help="file path or directory of files to process")
    parser.add_argument("--pdf-out", default=str(_PROJECT_ROOT / "documents" / "PDFs"))
    parser.add_argument("--meta-out", default=str(_PROJECT_ROOT / "documents" / "Metadata"))
    parser.add_argument("--no-ipfs", action="store_true", help="do not upload to IPFS")
    args = parser.parse_args(argv)

    inp = Path(args.input)
    if not inp.exists():
        print(json.dumps({"status":"error","error":"input not found"}))
        sys.exit(2)

    files = []
    if inp.is_dir():
        for ext in ("*.pdf","*.png","*.jpg","*.jpeg"):
            files += list(inp.glob(ext))
    else:
        files = [inp]

    results = []
    for f in files:
        try:
            out = process_input_file(str(f), pdf_out_dir=args.pdf_out, meta_out_dir=args.meta_out, upload_to_ipfs=not args.no_ipfs)
            results.append(out)
            print(json.dumps({"status":"ok","file": str(f), "pdf_cid": out.get("pdf_cid")}))
        except Exception as e:
            results.append({"status":"error","file": str(f), "error": str(e)})
            print(json.dumps({"status":"error","file": str(f), "error": str(e)}))

    # write summary
    summary_path = os.path.join(args.meta_out, "batch_summary.json")
    os.makedirs(args.meta_out, exist_ok=True)
    with open(summary_path, "w", encoding="utf-8") as fh:
        fh.write(json.dumps(results, indent=2))

    return 0

if __name__ == "__main__":
    sys.exit(main())
