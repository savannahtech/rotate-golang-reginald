# Define variables
OSQUERY_SOCKET_PATH := /Users/$(USER)/.osquery/shell.em
CMD_DIR := cmd/api

# Default target
.PHONY: all
all: help

# Help target
.PHONY: help
help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Available targets:"
	@echo "  run-mac-dev        - Run the development build on macOS"
	@echo "  run-windows-dev    - Run the development build on Windows"
	@echo "  build              - Build the project"
	@echo ""

# macOS dev target
.PHONY: run-mac-dev
run-mac-dev:
	@echo "Starting osquery in a separate terminal"
	@echo "Run the following command in a separate terminal to start osquery:"
	@echo "osqueryi --nodisable_extensions"
	@echo "Retrieve the extensions socket using the following query:"
	@echo "osquery> select value from osquery_flags where name = 'extensions_socket';"
	@echo "Ensure the socket path matches: $(OSQUERY_SOCKET_PATH)"
	@echo ""
	@echo "Starting Wails dev server"
	cd $(CMD_DIR) && wails dev

# Windows dev target
.PHONY: run-windows-dev
run-windows-dev:
	@echo "Starting osqueryd"
	@echo "Make sure osqueryd is running before starting the dev server."
	@echo "Run the following command in a separate terminal to start osquery:"
	@echo "osqueryi --nodisable_extensions"
	@echo "osquery> select value from osquery_flags where name = 'extensions_socket';"
	@echo "Ensure the socket path matches: | \\.\pipe\shell.em |"
	cd $(CMD_DIR) && wails dev

# Build target
.PHONY: build
build:
	@echo "Building the Wails project"
	cd $(CMD_DIR) && wails build


.PHONY: build-nsis
build-nsis:
	@echo "Building the Wails project"
	cd $(CMD_DIR) && wails build -nsis
