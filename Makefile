# ==============================================================================
# Arguments passing to Makefile commands
GO_INSTALLED := $(shell which go)
SS_INSTALLED := $(shell which staticcheck 2> /dev/null)

GITHUB=gotrackery
PROJECT_NAME=$(notdir $(shell pwd))

# ==============================================================================
# Install commands
install-tools:
	@echo Checking tools are installed...
ifndef SS_INSTALLED
	@echo Installing staticcheck...
	@go install honnef.co/go/tools/cmd/staticcheck@latest
endif

# ==============================================================================
# Modules support
tidy:
	@echo Running go mod tidy...
	@go mod tidy

# ==============================================================================
# Test commands
lint: install-tools
	@echo Running lints...
	@go vet ./...
	@staticcheck ./...
	@golangci-lint run

tests:
	@echo Running tests...
	@go test -v -race -vet=off $$(go list ./... | grep -v /test/)

cover:
	@echo Running coverage tests...
	@go test -vet=off -coverprofile ./cover.out $$(go list ./... | grep -v /test/)
	@go tool cover -html=./cover.out