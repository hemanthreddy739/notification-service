.PHONY: test coverage coverage-html coverage-func bench build clean run help

help:
	@echo "Notification Service - Available Commands"
	@echo "=========================================="
	@echo "make test          - Run all tests"
	@echo "make coverage      - Run tests with coverage"
	@echo "make coverage-html - Generate HTML coverage report"
	@echo "make coverage-func - Display function-level coverage"
	@echo "make bench         - Run benchmark tests"
	@echo "make build         - Build the application"
	@echo "make run           - Run the application"
	@echo "make clean         - Clean build artifacts and coverage files"
	@echo ""

test:
	@echo "Running tests..."
	@go test -v

coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out

coverage-html: coverage
	@echo "Generating HTML coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

coverage-func: coverage
	@echo "Function-level coverage:"
	@go tool cover -func=coverage.out
	@echo ""
	@echo "Total coverage:"
	@go tool cover -func=coverage.out | grep total

bench:
	@echo "Running benchmark tests..."
	@go test -bench=. -benchmem

build:
	@echo "Building notification service..."
	@go build -o notification-service

run: build
	@echo "Starting notification service on :8080..."
	@./notification-service

clean:
	@echo "Cleaning up..."
	@rm -f coverage.out coverage.html notification-service
	@go clean

all: clean test coverage-html bench
	@echo "All tasks completed!"
