#!/bin/bash
# build.sh — Build the single-binary release of digital-eval-system.
# Run this from the project root (digital-eval-system/).
set -e

echo "=== Building frontend ==="
cd digital-eval-ui
npm run build
cd ..

echo "=== Copying dist into Go project ==="
rm -rf services/go-node/ui/dist
cp -r digital-eval-ui/dist services/go-node/ui/dist

echo "=== Building Go binary ==="
cd services/go-node
go build -o node.exe ./cmd/node

echo ""
echo "✅  Build complete: services/go-node/node.exe"
echo "    Run with: cd services/go-node && ./node.exe -config configs/config.yaml"
echo "    Then open: http://127.0.0.1:8443"
