# QuickBooks Online MCP Server

A native Go implementation of a Model Context Protocol (MCP) server for QuickBooks Online integration.

## Quick Start

### Prerequisites

- Go 1.25+
- QuickBooks Online credentials

### Build & Run

```bash
# Build the binary
go build -o qbo-mcp ./cmd/qbo-mcp

# Run with stdio (default for local MCP clients)
./qbo-mcp

# Or run with HTTP
PORT=3000 ./qbo-mcp
```

### Transport Modes

- **Stdio (default)** — For Claude Desktop and local MCP clients. No `PORT` environment variable needed.
- **HTTP (Streamable MCP)** — Set `PORT` (e.g., `PORT=3000`) or `MCP_TRANSPORT=http`. HTTP path defaults to `/mcp` (customizable with `MCP_HTTP_PATH`).

## Configuration

Set environment variables in a `.env` file or export them. In single-tenant mode, the server will actively prompt you for missing QuickBooks credentials at startup if they are not provided:

```env
# QuickBooks API credentials
QUICKBOOKS_CLIENT_ID=your_client_id
QUICKBOOKS_CLIENT_SECRET=your_client_secret
QUICKBOOKS_REALM_ID=your_realm_id
QUICKBOOKS_REFRESH_TOKEN=your_refresh_token
QUICKBOOKS_ENVIRONMENT=sandbox  # or 'production'

# Optional: Platform integration (multi-tenant mode)
PLATFORM_INT_URL=https://platform-api.example.com
ENFORCE_AUTH=true

# Optional: Usage tracking
USAGE_REPORT_ENDPOINT=https://usage-tracker.example.com

# Optional: License server integration
LICENSE_SERVER_BASE_URL=https://license-server.example.com
LICENSE_SERVER_JWKS_ENDPOINT=/.well-known/jwks.json
LICENSE_SERVER_ACTIVATION_ENDPOINT=/activate
LICENSE_KEY=your_license_key

# Optional: HTTP server
PORT=3000
MCP_HTTP_PATH=/mcp
```

### Getting Credentials

1. Go to the [Intuit Developer Portal](https://developer.intuit.com/)
2. Create or select your app
3. Copy your Client ID and Client Secret
4. Add `http://localhost:3000/callback` to your app's Redirect URIs
5. Obtain your Refresh Token and Realm ID (see OAuth Flow below)

## Authentication

### Option 1: Pre-configured Tokens (Recommended for local development)

If you already have a refresh token and realm ID:

```env
QUICKBOOKS_REFRESH_TOKEN=your_refresh_token
QUICKBOOKS_REALM_ID=your_realm_id
```

### Option 2: OAuth Flow Setup

The server automatically handles OAuth token refreshes using your `QUICKBOOKS_REFRESH_TOKEN`. To obtain initial tokens:

1. Use the Intuit OAuth flow in the [Developer Portal](https://developer.intuit.com/)
2. Save the returned `refresh_token` and `realmId`
3. Add them to your `.env` file

## Available Tools

The server provides the following tool operations for QuickBooks Online entities:

**Supported Entities:**
- Account
- Bill
- Bill Payment
- Customer
- Employee
- Estimate
- Invoice
- Item
- Journal Entry
- Purchase
- Vendor

**Operations:**
- `quickbook_create_*` — Create a new entity
- `quickbook_get_*` — Retrieve entity by ID
- `quickbook_update_*` — Update an existing entity
- `quickbook_delete_*` — Delete/deactivate an entity
- `quickbook_search_*` — Search entities with criteria

Example: `quickbook_search_customers`, `quickbook_create_invoice`, `quickbook_update_bill`

## Platform Integration Mode

For multi-tenant deployments, set `PLATFORM_INT_URL` to enable platform integration:

```env
PLATFORM_INT_URL=https://platform-api.example.com
ENFORCE_AUTH=true
```

In this mode:
- The server fetches QuickBooks credentials from the platform API using the request's Bearer token
- Each HTTP request must include a valid `Authorization: Bearer <token>` header
- The Bearer token is used to retrieve the authenticated user's QB project credentials

## Docker Deployment

### Build Docker Image

```bash
docker build -t qbo-mcp-server:latest .
```

### Run with compose.yaml (HTTP mode)

```bash
docker compose up
```

This runs the server on `http://localhost:3000/mcp`

### Run with docker-compose.yaml (Production)

```bash
SHA=latest docker compose -f docker-compose.yaml up
```

This assumes credentials are in `.env` and maps port 8011 to the container's port 3000.

## Error Handling

**"QuickBooks not connected"**
- Verify `.env` contains all required variables
- Check that tokens are valid (refresh tokens expire periodically)
- For platform mode, ensure Bearer token is valid and has access to QB credentials

**"Missing authorization token"**
- In platform integration mode, include `Authorization: Bearer <token>` in HTTP requests
- Or disable with `ENFORCE_AUTH=false` for local development

**"Invalid OAuth token"**
- Your refresh token may have expired; re-authenticate via the OAuth flow
- Verify `QUICKBOOKS_ENVIRONMENT` matches where your credentials were issued

## Development

### Project Structure

```
.
├── cmd/
│   └── qbo-mcp/
│       └── main.go           # Server entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration & environment loading
│   ├── qbo/
│   │   ├── client.go         # QB API client
│   │   ├── query.go          # Query building
│   │   ├── unwrap.go         # Response parsing
│   │   └── customer_delete.go # Entity-specific operations
│   └── tools/
│       ├── register.go       # Tool registration
│       └── helpers.go        # Utility functions
└── go.mod / go.sum         # Go module management
```

### Running Tests

```bash
go test ./...
```

### Build Flags

```bash
# Disable CGO for better portability
CGO_ENABLED=0 go build -o qbo-mcp ./cmd/qbo-mcp

# Linux binary for Docker
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o qbo-mcp ./cmd/qbo-mcp
```

## License

MIT

## Support

For issues, questions, or contributions, please refer to [CONTRIBUTING.md](CONTRIBUTING.md) and [DOCUMENTATION.md](DOCUMENTATION.md).

