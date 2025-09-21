.PHONY: test build run clean coverage help sqlc-generate

# Default target
help:
	@echo "Available targets:"
	@echo "  test            - Run all tests"
	@echo "  coverage        - Run tests with coverage report"
	@echo "  build           - Build the application"
	@echo "  run             - Run the application locally"
	@echo "  clean           - Clean build artifacts"
	@echo "  sqlc-generate   - Generate sqlc database code"
	@echo "  help            - Show this help message"

# Run tests
test:
	go test -v ./...

# Run tests with coverage
coverage:
	go test -v -cover ./...
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build the application
build:
	go build -o choochoo .

# Run the application locally
run: build
	./choochoo

# Clean build artifacts
clean:
	rm -f choochoo coverage.out coverage.html

# Install dependencies (if any)
deps:
	go mod tidy
	go mod download

# Verify everything is working
verify: deps test build
	@echo "All checks passed!"

# Generate sqlc database code
sqlc-generate:
	~/go/bin/sqlc generate