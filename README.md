# choochoo
Hop on the train

A simple Go server for accepting GitHub webhooks.

## Features

- ✅ HTTP server with webhook endpoint
- ✅ GitHub webhook signature validation
- ✅ Request logging and error handling
- ✅ Health check endpoint
- ✅ Configurable port and webhook secret
- ✅ PostgreSQL database integration with sqlc
- ✅ Selective webhook storage (push, issue_comment, pull_request events)

## Quick Start

1. **Build the application:**
   ```bash
   go build -o choochoo .
   ```

2. **Set up your environment** (optional):
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Run the server:**
   ```bash
   # With default settings (port 8080, no signature validation)
   ./choochoo
   
   # With custom port
   PORT=3000 ./choochoo
   
   # With webhook secret for signature validation
   GITHUB_WEBHOOK_SECRET="your-secret" ./choochoo
   
   # With database for storing webhook events
   DATABASE_URL="postgres://user:pass@localhost:5432/choochoo?sslmode=disable" ./choochoo
   
   # With both secret and database
   GITHUB_WEBHOOK_SECRET="your-secret" DATABASE_URL="postgres://user:pass@localhost:5432/choochoo?sslmode=disable" ./choochoo
   ```

## Endpoints

- `POST /webhook` - GitHub webhook endpoint
- `GET /health` - Health check endpoint
- `GET /` - Server information

## Configuration

The server can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Port to run the server on | `8080` |
| `GITHUB_WEBHOOK_SECRET` | Secret for webhook signature validation | (none) |
| `DATABASE_URL` | PostgreSQL connection string for storing webhook events | (none) |

### Database Configuration

When `DATABASE_URL` is set, the server will store supported webhook events in a PostgreSQL database. The following event types are stored:

- `push` - Git push events
- `issue_comment` - Issue comment events  
- `pull_request` - Pull request events

All other webhook events are logged but not stored in the database.

**Database URL Format:**
```
postgres://username:password@hostname:port/database_name?sslmode=disable
```

**Example:**
```bash
DATABASE_URL="postgres://postgres:password@localhost:5432/choochoo?sslmode=disable"
```

## Database Setup

If you want to store webhook events in a PostgreSQL database:

1. **Set up PostgreSQL** (if not already installed):
   ```bash
   # On macOS
   brew install postgresql
   brew services start postgresql
   
   # On Ubuntu/Debian
   sudo apt-get install postgresql postgresql-contrib
   sudo systemctl start postgresql
   ```

2. **Create the database:**
   ```bash
   # Connect to PostgreSQL
   psql -U postgres
   
   # Create database
   CREATE DATABASE choochoo;
   
   # Exit psql
   \q
   ```

3. **Run database migrations:**
   ```bash
   # Apply the schema
   psql -U postgres -d choochoo -f sql/migrations/001_webhook_events.sql
   ```

4. **Set the DATABASE_URL environment variable:**
   ```bash
   export DATABASE_URL="postgres://postgres:password@localhost:5432/choochoo?sslmode=disable"
   ```

The server will automatically connect to the database on startup and store supported webhook events.

## GitHub Webhook Setup

1. Go to your GitHub repository settings
2. Navigate to "Webhooks" → "Add webhook"
3. Set the payload URL to: `http://your-server:port/webhook`
4. Set content type to `application/json`
5. Set a secret (optional but recommended)
6. Select the events you want to receive
7. Save the webhook

## Security

- The server validates GitHub webhook signatures when `GITHUB_WEBHOOK_SECRET` is set
- Always use HTTPS in production environments
- Keep your webhook secret secure and rotate it regularly

## Continuous Integration

This repository uses GitHub Actions for automated testing and quality assurance:

- **Automated Testing**: All pull requests and pushes to main automatically run the full test suite
- **Build Verification**: The CI pipeline ensures the application builds successfully on every change
- **Coverage Reports**: Test coverage is automatically generated and tracked
- **Branch Protection**: The main branch requires passing CI checks before merging pull requests

### Setting Up Branch Protection

To require passing tests for pull request merges, configure branch protection rules. For detailed step-by-step instructions, see [docs/branch-protection.md](docs/branch-protection.md).

**Quick Setup:**
1. Go to repository **Settings** → **Branches**
2. Click **Add rule** for `main` branch
3. Enable:
   - ☑️ **Require a pull request before merging**
   - ☑️ **Require status checks to pass before merging**
   - ☑️ **Require branches to be up to date before merging**
   - Search and select the **CI** status check
   - ☑️ **Include administrators** (recommended)

This ensures all code changes go through proper review and testing before being merged.

## Development

### Running the Server

Run the server in development mode:
```bash
go run main.go
```

### Testing

This project includes comprehensive tests to ensure functionality works as expected. **All contributions must include tests.**

Run tests:
```bash
# Run all tests
make test

# Run tests with coverage report
make coverage

# Or using go directly
go test -v
go test -v -cover
```

### Database Development

The project uses [sqlc](https://sqlc.dev/) for type-safe SQL operations. After modifying SQL queries or schema:

```bash
# Regenerate database code
make sqlc-generate

# Or directly:
~/go/bin/sqlc generate
```

Database files:
- `sql/migrations/` - Database schema migrations
- `sql/queries/` - SQL queries for sqlc
- `internal/db/` - Generated sqlc code (do not edit manually)

### Manual Testing

Test the webhook endpoint manually:
```bash
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: push" \
  -H "X-GitHub-Delivery: test-delivery-id" \
  -d '{"action":"push","repository":{"full_name":"user/repo"},"sender":{"login":"username"}}'
```

### Contributing

Before contributing, please read [CONTRIBUTING.md](CONTRIBUTING.md) for development guidelines and testing requirements.

### Available Make Targets

```bash
make test      # Run all tests
make coverage  # Run tests with coverage report
make build     # Build the application
make run       # Run the application locally
make clean     # Clean build artifacts
make help      # Show available targets
```
