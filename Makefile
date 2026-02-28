SHELL := /bin/bash
.PHONY: help keys setup start test lint migrate clean

help:
	@echo "Available commands:"
	@echo "  help     - show this help"
	@echo "  keys     - generate TLS certs and JWT signing keys (local)"
	@echo "  setup    - run local setup script (installs CLI deps, checks services)"
	@echo "  start    - placeholder to start all services (phase 0 no-op)"
	@echo "  test     - run unit test placeholders"
	@echo "  migrate  - run DB migrations (placeholder - apply SQL in infra/migrations)"
	@echo "  lint     - run lint checks (placeholders)"
	@echo "  clean    - remove generated certs and keys in infra/certs"

keys:
	@echo "Generating TLS certs and JWT keys..."
	@chmod +x tools/keygen/gen_keys.sh
	@./tools/keygen/gen_keys.sh infra/certs

setup:
	@echo "Running local setup checks/installers..."
	@chmod +x tools/local_setup.sh
	@./tools/local_setup.sh

start:
	@echo "Phase 0: no services to start."

test:
	@echo "Phase 0: no tests yet."

migrate:
	@echo "Phase 0: migrations placeholder."

lint:
	@echo "Phase 0: no linters configured."

clean:
	@echo "Cleaning generated certs in infra/certs"
	@rm -rf infra/certs/* || true
	@mkdir -p infra/certs
	@touch infra/certs/.gitkeep
	@echo "Done."
