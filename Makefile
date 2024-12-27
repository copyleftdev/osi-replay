# -------------------------------------------------------------
# Makefile for osi-replay
# A multi-command Go project with subcommands for capture,
# replay, transform, and rewriter.
#
# Targets:
#   make build         - Compiles all commands into ./bin/ folder
#   make test          - Runs go tests on the entire project
#   make tidy          - Cleans up and updates go.mod / go.sum
#   make run-capture   - Builds and runs the capture command (example usage)
#   make run-replay    - Builds and runs the replay command (example usage)
#   make run-transform - Builds and runs the transform command (example usage)
#   make run-rewriter  - Builds and runs the rewriter command (example usage)
#   make clean         - Removes build artifacts
#   make help          - Displays these instructions
# -------------------------------------------------------------

# Default target, for convenience
.DEFAULT_GOAL := build

# List of subcommands in cmd/
SUBCOMMANDS = capture replay transform rewriter

# The binary directory
BIN_DIR = bin

# -------------------------------------------------------------
# PHONY declarations (targets that are not actual files)
# -------------------------------------------------------------
.PHONY: build test tidy clean help \
        run-capture run-replay run-transform run-rewriter

# -------------------------------------------------------------
# BUILD: Compile all commands
# -------------------------------------------------------------
build:
	@echo ">> Building all subcommands into $(BIN_DIR)/ ..."
	@mkdir -p $(BIN_DIR)
	@for cmd in $(SUBCOMMANDS); do \
		echo "   -> Building $$cmd"; \
		go build -o $(BIN_DIR)/$$cmd ./cmd/$$cmd; \
	done
	@echo ">> Build complete!"

# -------------------------------------------------------------
# TEST: Run go test on all packages
# -------------------------------------------------------------
test:
	@echo ">> Running tests..."
	@go test ./... -cover
	@echo ">> Tests finished"

# -------------------------------------------------------------
# TIDY: Cleans up dependencies
# -------------------------------------------------------------
tidy:
	@echo ">> Tidying go.mod and go.sum..."
	@go mod tidy
	@echo ">> Tidy complete"

# -------------------------------------------------------------
# RUN-SUBCOMMAND: Example usage for each command
# -------------------------------------------------------------
run-capture: build
	@echo ">> Running capture command..."
	@$(BIN_DIR)/capture -i eth0 -o capture.pcap

run-replay: build
	@echo ">> Running replay command..."
	@$(BIN_DIR)/replay -i eth0 -f capture.pcap

run-transform: build
	@echo ">> Running transform command..."
	@$(BIN_DIR)/transform -in capture.pcap -out sanitized_capture.pcap

run-rewriter: build
	@echo ">> Running rewriter command..."
	@$(BIN_DIR)/rewriter -in capture.pcap -out rewritten_capture.pcap

# -------------------------------------------------------------
# CLEAN: Removes build artifacts
# -------------------------------------------------------------
clean:
	@echo ">> Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)
	@echo ">> Clean complete"

# -------------------------------------------------------------
# HELP: Display usage
# -------------------------------------------------------------
help:
	@echo "OSI-REPLAY MAKEFILE - Available targets:"
	@echo ""
	@echo "  build           Compile all commands into 'bin/' folder."
	@echo "  test            Run 'go test' for all packages."
	@echo "  tidy            Run 'go mod tidy' to clean up modules."
	@echo "  run-capture     Build & run the capture command (example usage)."
	@echo "  run-replay      Build & run the replay command (example usage)."
	@echo "  run-transform   Build & run the transform command (example usage)."
	@echo "  run-rewriter    Build & run the rewriter command (example usage)."
	@echo "  clean           Remove build artifacts."
	@echo "  help            Show this help text."
	@echo ""
	@echo "Usage examples:"
	@echo "  make build          - builds all commands"
	@echo "  make test           - runs go tests"
	@echo "  make run-capture    - builds then runs './bin/capture'"
	@echo "  make clean          - cleans up"
