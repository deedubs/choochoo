# Choochoo Features

This document provides a comprehensive overview of all features and capabilities in the choochoo GitHub webhook server.

## Core Features

### 🚀 HTTP Server
- **Fast startup**: Starts immediately with minimal dependencies
- **Configurable port**: Default port 8080, customizable via `PORT` environment variable
- **Graceful shutdown**: Proper HTTP server lifecycle management
- **Zero external runtime dependencies**: Single binary deployment

### 🔐 GitHub Webhook Processing
- **Webhook endpoint**: `POST /webhook` for receiving GitHub webhook events
- **Event logging**: Comprehensive logging of all received webhook events
- **Event validation**: Validates required GitHub headers (`X-GitHub-Event`, `X-GitHub-Delivery`)
- **JSON payload parsing**: Robust parsing with error handling for malformed payloads
- **Repository and sender tracking**: Extracts and logs repository name and sender information

### 🛡️ Security Features
- **Webhook signature validation**: HMAC-SHA256 signature verification using `X-Hub-Signature-256` header
- **Configurable security**: Enable/disable signature validation via `GITHUB_WEBHOOK_SECRET` environment variable
- **Constant-time comparison**: Secure signature validation to prevent timing attacks
- **Request method validation**: Only accepts POST requests to webhook endpoint
- **Input validation**: Validates all incoming data before processing

### 💾 Database Integration
- **PostgreSQL support**: Optional PostgreSQL database integration for webhook storage
- **Type-safe SQL operations**: Uses [sqlc](https://sqlc.dev/) for generated, type-safe database code
- **Selective event storage**: Only stores supported event types (push, issue_comment, pull_request)
- **Comprehensive database schema**: Includes indexes for efficient querying
- **Database connection management**: Automatic connection handling with error recovery

#### Supported Event Types for Database Storage
- **`push`**: Git push events (commits, branch updates)
- **`issue_comment`**: Comments on issues and pull requests  
- **`pull_request`**: Pull request creation, updates, and state changes

All other webhook events are logged but not stored in the database.

#### Database Schema
```sql
CREATE TABLE webhook_events (
    id SERIAL PRIMARY KEY,
    delivery_id VARCHAR(255) NOT NULL UNIQUE,
    event_type VARCHAR(50) NOT NULL,
    repository_name VARCHAR(255),
    sender_login VARCHAR(255),
    action VARCHAR(100),
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

#### Database Operations
- **Event creation**: Store webhook events with full JSON payload
- **Event retrieval**: Query by delivery ID, event type, or repository
- **Event listing**: Paginated listing with filtering and sorting
- **Event counting**: Count events by type for analytics
- **Event cleanup**: Delete old events for maintenance

## API Endpoints

### `POST /webhook`
**Purpose**: Receive and process GitHub webhook events

**Headers**:
- `Content-Type: application/json` (required)
- `X-GitHub-Event`: Event type (e.g., "push", "pull_request")
- `X-GitHub-Delivery`: Unique delivery identifier
- `X-Hub-Signature-256`: HMAC-SHA256 signature (when webhook secret is configured)

**Request Body**: JSON payload from GitHub webhook

**Responses**:
- `200 OK`: Webhook processed successfully
- `400 Bad Request`: Invalid request body or missing headers
- `401 Unauthorized`: Invalid webhook signature
- `405 Method Not Allowed`: Non-POST request

**Example Response**:
```json
{
  "message": "Webhook received and processed",
  "status": "success"
}
```

### `GET /health`
**Purpose**: Health check endpoint for monitoring and load balancers

**Response**:
```json
{
  "service": "choochoo-webhook-server",
  "status": "healthy"
}
```

### `GET /`
**Purpose**: Server information and endpoint listing

**Response**: HTML page with server information and available endpoints

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | HTTP server port | `8080` | No |
| `GITHUB_WEBHOOK_SECRET` | Secret for webhook signature validation | (none) | No |
| `DATABASE_URL` | PostgreSQL connection string | (none) | No |

### Configuration Modes

#### Basic Mode (No Configuration)
- Server runs on port 8080
- Webhook signature validation is skipped
- Events are logged but not stored in database
- Suitable for development and testing

#### Secure Mode (With Webhook Secret)
```bash
GITHUB_WEBHOOK_SECRET="your-secret-key"
```
- Enables webhook signature validation
- Rejects webhooks with invalid signatures
- Recommended for production deployments

#### Database Mode (With Database URL)
```bash
DATABASE_URL="postgres://user:pass@localhost:5432/choochoo?sslmode=disable"
```
- Enables webhook event storage
- Stores supported events in PostgreSQL
- Provides data persistence and analytics capabilities

#### Full Production Mode
```bash
PORT=3000
GITHUB_WEBHOOK_SECRET="your-secret-key"
DATABASE_URL="postgres://user:pass@localhost:5432/choochoo?sslmode=disable"
```

## Architecture

### Package Structure
- **`main.go`**: Application entry point
- **`internal/server`**: HTTP server setup and configuration
- **`internal/handlers`**: HTTP request handlers for all endpoints
- **`internal/webhook`**: Webhook event types and processing logic
- **`internal/database`**: Database connection management
- **`internal/db`**: Generated sqlc database code (do not edit manually)

### Request Flow
1. **HTTP Request**: Incoming webhook request to `/webhook`
2. **Method Validation**: Verify POST method
3. **Header Extraction**: Extract GitHub headers (event type, delivery ID, signature)
4. **Signature Validation**: Verify HMAC-SHA256 signature (if secret configured)
5. **Payload Parsing**: Parse JSON payload and extract event data
6. **Event Processing**: Log event details and extract metadata
7. **Database Storage**: Store event in database (if configured and supported event type)
8. **Response**: Return success/error response

## Development Features

### Testing
- **Comprehensive test suite**: 100% coverage of core functionality
- **Unit tests**: Individual function and method testing
- **Integration tests**: End-to-end request/response testing
- **Security tests**: Signature validation and authentication testing
- **Edge case testing**: Error conditions and malformed input handling

### Build System
- **Makefile**: Standardized build targets
- **Fast builds**: Single binary compilation in ~0.4 seconds
- **Test execution**: Full test suite runs in ~18 seconds
- **Coverage reporting**: HTML coverage reports generated
- **Code generation**: Automatic sqlc database code generation

### Development Workflow
- **Live reloading**: Use `go run main.go` for development
- **Database migrations**: SQL migrations for schema management
- **Code formatting**: Enforced Go formatting standards
- **Static analysis**: `go vet` integration for code quality

## Security Considerations

### Webhook Security
- **Signature verification**: Prevents unauthorized webhook submissions
- **Timing attack protection**: Constant-time signature comparison
- **Secret management**: Environment variable-based secret configuration
- **HTTPS requirement**: Recommended for production deployments

### Input Validation
- **JSON validation**: Robust parsing with error handling
- **Header validation**: Required GitHub headers verification
- **Method validation**: Only POST requests accepted
- **Content-type validation**: Requires application/json

### Database Security
- **Parameterized queries**: Protection against SQL injection
- **Connection string security**: Secure credential management
- **Data sanitization**: Safe handling of user-provided data

## Performance Characteristics

### Server Performance
- **Startup time**: Immediate (milliseconds)
- **Build time**: ~0.4 seconds
- **Test time**: ~18 seconds for full suite
- **Memory usage**: Minimal footprint with no external dependencies
- **Concurrent requests**: Standard Go HTTP server concurrency

### Database Performance
- **Indexed queries**: Optimized database indexes for common queries
- **JSON storage**: Efficient JSONB storage for webhook payloads
- **Pagination support**: Efficient pagination for large datasets
- **Connection pooling**: Managed database connection pooling

## Monitoring and Observability

### Logging
- **Structured logging**: Consistent log format across components
- **Event logging**: Detailed webhook event information
- **Error logging**: Comprehensive error tracking and reporting
- **Security logging**: Authentication and validation events

### Health Monitoring
- **Health endpoint**: `/health` for load balancer checks
- **Database health**: Connection status monitoring
- **Service status**: Overall service health reporting

### Metrics and Analytics
- **Event counting**: Database queries for event analytics
- **Repository tracking**: Events grouped by repository
- **Sender tracking**: Events grouped by GitHub user
- **Time-based analysis**: Events with timestamp tracking

## Deployment Options

### Single Binary Deployment
- **No dependencies**: Single executable file
- **Cross-platform**: Builds for multiple operating systems
- **Container ready**: Suitable for Docker/Kubernetes deployment
- **Minimal resources**: Low memory and CPU requirements

### Environment-based Configuration
- **12-factor app**: Configuration via environment variables
- **No config files**: Zero configuration file dependencies
- **Runtime configuration**: Dynamic configuration without rebuilds
- **Secret management**: Secure environment variable handling

## Extensibility

### Adding New Event Types
1. Update `internal/webhook/types.go` `SupportedEventTypes` map
2. Add tests in `internal/webhook/types_test.go`
3. Update this documentation

### Adding New Endpoints
1. Create handler in `internal/handlers/`
2. Register route in `internal/server/server.go`
3. Add tests for new functionality
4. Update documentation

### Database Schema Changes
1. Create new migration in `sql/migrations/`
2. Update queries in `sql/queries/`
3. Regenerate sqlc code with `make sqlc-generate`
4. Update tests and documentation

## Version Information

This documentation reflects the current state of the choochoo codebase. Features and capabilities are actively maintained and updated as the codebase evolves.

For the most up-to-date information about specific implementation details, refer to:
- **README.md**: Quick start and basic usage
- **CONTRIBUTING.md**: Development guidelines and testing requirements
- **Source code**: Authoritative implementation details