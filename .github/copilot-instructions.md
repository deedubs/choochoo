# Choochoo GitHub Webhook Server

**ALWAYS follow these instructions first and only fallback to additional search and context gathering if the information here is incomplete or found to be in error.**

Choochoo is a simple Go-based HTTP server that receives and processes GitHub webhooks. It validates webhook signatures, logs events, and provides health monitoring endpoints. The application has no external dependencies and builds/runs quickly.

## Quick Start - Build, Test, and Run

Bootstrap and validate the repository:
```bash
# Install dependencies (completes in ~4 seconds)
make deps

# Run all tests (completes in ~18 seconds, timeout: 60 seconds)
make test

# Build application (completes in ~0.4 seconds, timeout: 30 seconds)  
make build

# Run complete verification (completes in ~0.2 seconds, timeout: 60 seconds)
make verify
```

**NEVER CANCEL** any command. All operations complete quickly but always set appropriate timeouts.

## Development Workflow

### Building and Testing

Run the complete development cycle:
```bash
# Dependencies and build
make deps     # ~4 seconds - downloads Go modules (none currently)
make build    # ~0.4 seconds - creates ./choochoo binary

# Testing (comprehensive test suite with 75.3% coverage)
make test     # ~18 seconds - runs 17 tests, all should pass
make coverage # ~1.5 seconds - generates coverage.html report

# Alternative: use Go directly
go test -v           # ~0.3 seconds - run tests with verbose output
go test -v -cover    # ~0.3 seconds - run tests with coverage percentage
go vet ./...         # ~0.1 seconds - static analysis, should show no issues
```

### Running the Application

Start the webhook server:
```bash
# Option 1: Using make (builds first)
make run              # Starts on port 8080 (default)

# Option 2: Using go run directly
go run main.go        # Starts immediately

# Option 3: Using built binary
./choochoo           # Must run `make build` first

# With custom configuration
PORT=3000 ./choochoo                                    # Custom port
GITHUB_WEBHOOK_SECRET="your-secret" ./choochoo        # Enable signature validation
PORT=3000 GITHUB_WEBHOOK_SECRET="secret" ./choochoo   # Both custom port and secret
```

The server provides these endpoints:
- `GET /health` - Health check endpoint (returns JSON status)
- `GET /` - Server information page
- `POST /webhook` - GitHub webhook endpoint with signature validation

## Validation Scenarios

**ALWAYS test these scenarios after making changes:**

### 1. Basic Server Functionality
```bash
# Start server in background
./choochoo &
sleep 2

# Test health endpoint
curl -s http://localhost:8080/health
# Expected: {"service":"choochoo-webhook-server","status":"healthy"}

# Test info endpoint  
curl -s http://localhost:8080/
# Expected: "Choochoo GitHub Webhook Server" with endpoint listing

# Stop server
pkill choochoo
```

### 2. Webhook Processing (No Secret)
```bash
# Start server without secret
./choochoo &
sleep 2

# Send test webhook
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: push" \
  -H "X-GitHub-Delivery: test-delivery-id" \
  -d '{"action":"push","repository":{"full_name":"user/repo"},"sender":{"login":"username"}}'

# Expected response: {"message":"Webhook received and processed","status":"success"}
# Expected server log: "Received push event from user/repo (delivery: test-delivery-id, sender: username)"

pkill choochoo
```

### 3. Webhook Signature Validation  
```bash
# Start server with secret
GITHUB_WEBHOOK_SECRET="test-secret" ./choochoo &
sleep 2

# Generate valid signature
PAYLOAD='{"action":"push","repository":{"full_name":"user/repo"},"sender":{"login":"username"}}'
SIGNATURE=$(echo -n "$PAYLOAD" | openssl dgst -sha256 -hmac "test-secret" | sed 's/.* /sha256=/')

# Test with valid signature
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: push" \
  -H "X-GitHub-Delivery: test-delivery-id" \
  -H "X-Hub-Signature-256: $SIGNATURE" \
  -d "$PAYLOAD"
# Expected: Success response and logged event

# Test with invalid signature  
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: push" \
  -H "X-GitHub-Delivery: test-delivery-id" \
  -H "X-Hub-Signature-256: sha256=invalid" \
  -d "$PAYLOAD"
# Expected: "Invalid signature" response and error log

pkill choochoo
```

**ALWAYS run at least scenario 1 after making any changes.**

## Code Quality Standards

### Before Committing Changes

Run these commands before any commit:
```bash
# Format code (fixes whitespace and formatting)
gofmt -w .

# Verify code quality
go vet ./...          # Should show no issues
make test            # All 17 tests must pass
make build           # Must build successfully
```

**Required:** Code must be properly formatted with `gofmt -w .` - the repository currently has formatting issues that should be fixed.

### Test Coverage Requirements

- Current coverage: 75.3% of statements
- All new code must include tests
- Use existing test patterns in `main_test.go`
- Test both success and error scenarios
- Include signature validation tests for webhook changes

## Repository Structure

### Key Files
```
/
├── main.go              # Main server implementation (~150 lines)
├── main_test.go         # Comprehensive test suite (17 tests)  
├── Makefile            # Build targets: deps, test, coverage, build, run, clean, verify
├── go.mod              # Go module definition (no external dependencies)
├── .env.example        # Environment variable examples
├── README.md           # Project documentation
├── CONTRIBUTING.md     # Development guidelines and CI requirements
└── .github/
    └── workflows/
        └── ci.yml      # GitHub Actions CI pipeline
```

### Important Directories
- `.github/workflows/` - CI configuration (runs make verify on PRs)
- `docs/` - Additional documentation (branch protection setup)

## Environment Configuration

The server uses these environment variables:
- `PORT` (default: 8080) - Server port
- `GITHUB_WEBHOOK_SECRET` (optional) - Enable webhook signature validation

Example setup:
```bash
cp .env.example .env
# Edit .env with your values
source .env
./choochoo
```

## CI/CD Integration

This repository uses GitHub Actions CI that runs on every PR and push to main:
- Installs Go 1.24.7
- Runs `make deps` (install dependencies)
- Runs `make test` (test suite)
- Runs `make coverage` (coverage report)  
- Runs `make build` (compilation)
- Runs `make verify` (complete validation)

**Always ensure your changes pass `make verify` locally before pushing.**

## Performance Characteristics

All operations are fast:
- Dependencies: ~4 seconds (no external deps to download)
- Tests: ~18 seconds (17 comprehensive tests)
- Build: ~0.4 seconds (single binary)
- Coverage: ~1.5 seconds (with HTML report)
- Startup: Immediate (server starts in milliseconds)

## Common Tasks Reference

### Frequently Used Commands
```bash
make help           # Show all available targets
make clean          # Remove build artifacts (choochoo binary, coverage files)
go mod tidy         # Clean up go.mod (part of make deps)
go run main.go      # Quick development server start
```

### Debugging Issues
- Check server logs for webhook processing details
- Use `go vet ./...` to find code issues
- Run tests with `-v` flag for verbose output
- Generate coverage report with `make coverage` (creates coverage.html)

### File Locations for Common Changes
- **Webhook handling logic**: `main.go`, `handleWebhook()` function
- **Server configuration**: `main.go`, `NewWebhookServer()` function  
- **Signature validation**: `main.go`, `validateSignature()` method
- **Test scenarios**: `main_test.go` (organized by function being tested)
- **CI configuration**: `.github/workflows/ci.yml`

## Security Notes

- Always validate webhook signatures in production (set GITHUB_WEBHOOK_SECRET)
- The server logs webhook events but not sensitive payload details
- Use HTTPS in production deployments
- Keep webhook secrets secure and rotate regularly

---

**Remember: Follow these instructions first, validate all changes with the test scenarios, and ensure proper formatting before committing.**