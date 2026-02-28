#!/usr/bin/env bash
# gen_keys.sh - generate local TLS cert and RSA keypair for JWT (development only)
# Usage: ./gen_keys.sh <output-dir>
set -euo pipefail

OUT_DIR="${1:-infra/certs}"
mkdir -p "${OUT_DIR}"

echo "Output directory: ${OUT_DIR}"

# Check for openssl
if ! command -v openssl >/dev/null 2>&1; then
  echo "Error: openssl not found. Install openssl and retry."
  exit 2
fi

# TLS cert (self-signed)
TLS_KEY="${OUT_DIR}/server.key"
TLS_CRT="${OUT_DIR}/server.crt"

if [ -f "${TLS_KEY}" ] || [ -f "${TLS_CRT}" ]; then
  echo "TLS cert/key already exist, skipping generation."
else
  echo "Generating self-signed TLS certificate..."
  openssl req -x509 -nodes -days 3650 -newkey rsa:4096 \
    -keyout "${TLS_KEY}" -out "${TLS_CRT}" -subj "/C=IN/ST=State/L=City/O=DigitalEval/OU=Dev/CN=localhost"
  chmod 600 "${TLS_KEY}"
  echo "Generated TLS: ${TLS_CRT}, ${TLS_KEY}"
fi

# RSA keypair for JWT (RS256)
JWT_PRIV="${OUT_DIR}/jwt_private.pem"
JWT_PUB="${OUT_DIR}/jwt_public.pem"

if [ -f "${JWT_PRIV}" ] || [ -f "${JWT_PUB}" ]; then
  echo "JWT keypair already exist, skipping generation."
else
  echo "Generating RSA keypair for JWT (RS256)..."
  openssl genpkey -algorithm RSA -out "${JWT_PRIV}" -pkeyopt rsa_keygen_bits:4096
  openssl rsa -pubout -in "${JWT_PRIV}" -out "${JWT_PUB}"
  chmod 600 "${JWT_PRIV}"
  chmod 644 "${JWT_PUB}"
  echo "Generated JWT keys: ${JWT_PRIV}, ${JWT_PUB}"
fi

# Create .gitignore reminder
GITIGNORE="${OUT_DIR}/.gitignore"
if [ ! -f "${GITIGNORE}" ]; then
  cat > "${GITIGNORE}" <<EOF
# Do not commit generated keys
*
!.gitkeep
EOF
  echo "Wrote ${GITIGNORE} (reminder to ignore certs in git)."
fi

echo "Key generation completed."
