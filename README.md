# choochoo
Hop on the train

A simple Go server for accepting GitHub webhooks.

## Features

- ✅ HTTP server with webhook endpoint
- ✅ GitHub webhook signature validation
- ✅ Request logging and error handling
- ✅ Health check endpoint
- ✅ Configurable port and webhook secret

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
   
   # With both
   PORT=3000 GITHUB_WEBHOOK_SECRET="your-secret" ./choochoo
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

## Development

Run the server in development mode:
```bash
go run main.go
```

Test the webhook endpoint:
```bash
curl -X POST http://localhost:8080/webhook \
  -H "Content-Type: application/json" \
  -H "X-GitHub-Event: push" \
  -H "X-GitHub-Delivery: test-delivery-id" \
  -d '{"action":"push","repository":{"full_name":"user/repo"},"sender":{"login":"username"}}'
```
