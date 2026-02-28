#!/usr/bin/env bash
# local_setup.sh - idempotent local setup checks for Phase 0
# Installs or checks for presence of: ipfs (go-ipfs CLI), psql client
set -euo pipefail

echo "Local setup check: verifying required CLI tools..."

check_cmd() {
  cmd=$1
  install_hint=$2
  if command -v "${cmd}" >/dev/null 2>&1; then
    echo " - ${cmd} found: $(command -v ${cmd})"
  else
    echo " - ${cmd} NOT found. Hint: ${install_hint}"
  fi
}

check_cmd "ipfs" "Install go-ipfs from https://docs.ipfs.io/install/ or use your package manager."
check_cmd "psql" "Install PostgreSQL client tools (psql). On Debian/Ubuntu: apt install postgresql-client."

echo ""
echo "Optional (recommended) checks:"
check_cmd "openssl" "Install openssl for key generation."
check_cmd "jq" "Install jq for JSON handling in scripts (optional)."

echo ""
echo "Local setup script complete. This script does not install missing packages automatically."
echo "If you want to install packages, run the appropriate OS package manager commands yourself with root privileges."
