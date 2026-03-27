# Complete Developer Onboarding Guide: QuickBooks Online MCP Server

Welcome to the **QuickBooks Online MCP Server**! This document is the **only** guide you need to understand, build, set up, and extend this project.

## 1. What the Project is About

This application is a **Model Context Protocol (MCP)** server built in Go. It acts as a bridge between an AI Agent (like Claude, ChatGPT with MCP support, or custom orchestration platforms) and **QuickBooks Online**. 

Instead of an AI hallucinating API endpoints, this server exposes over **40 specific QuickBooks operations** as structured "Tools" that the AI can call naturally. It handles the difficult parts:
- **Authentication**: Managing QuickBooks OAuth 2.0 flows and token refreshes.
- **Query Building**: Translating simple JSON search requests into QuickBooks SQL (`QBO`).
- **Transport**: Running over `stdio` (for local desktop AI agents) or `HTTP` (for distributed cloud deployments).

**Domains Supported:**
The server enables AI read/write functionality for Customers, Invoices, Bills, Bill Payments, Estimates, Vendors, Chart of Accounts, Items, Employees, Journal Entries, and Purchases.

---

## 2. Core Architecture

The architecture relies on a few core components operating together:

```
┌─────────────────────────┐         ┌────────────────────────┐
│     AI Agent Console    │ ──────► │   cmd/qbo-mcp/main.go  │
│ (Claude Desktop/Stream) │  stdio  │  (CLI Wizard & Server) │
└─────────────────────────┘  HTTP   └───────────┬────────────┘
                                                │
┌─────────────────────────┐         ┌───────────▼────────────┐
│      intuit-oauth       │ ◄────── │ internal/tools/        │
│    (OAuth 2.0 flow)     │         │ (Tool Handlers & QBO)  │
└───────────┬─────────────┘         └───────────┬────────────┘
            │                                   │
┌───────────▼─────────────┐         ┌───────────▼────────────┐
│ Platform API (Optional) │         │  QuickBooks REST API   │
│ (Per-user tenant creds) │         │  (Invoices, Bills...)  │
└─────────────────────────┘         └────────────────────────┘
```

### 2.1 The Two Modes: Single-Tenant vs Multi-Tenant

**A. Single-Tenant (Local Development Mode)**
- Ideal for personal use or local testing.
- Uses hardcoded QuickBooks Developer credentials stored in a `.env` file.
- `ENFORCE_AUTH=false`.
- You authenticate one specific QuickBooks Company.

**B. Multi-Tenant (Platform Integration Mode)**
- Built for enterprise deployment and cloud hosting.
- `ENFORCE_AUTH=true`.
- The server ignores local QuickBooks credentials. Instead, every incoming HTTP request must include an `Authorization: Bearer <token>` header.
- The server calls out to a `PLATFORM_INT_URL` to dynamically fetch the exact QuickBooks tokens for the user making the request.

### 2.2 The Interactive CLI Wizard
When you run the bare executable (`./qbo-mcp`), an interactive command-line wizard starts up before the MCP server initializes. 
- It prompts for a **License Key**.
- If you enter `1234`, it enters **Single-Tenant Mode** and asks you directly for your QuickBooks Client ID, Secret, and Tokens.
- If you enter `6789`, it enters **Multi-Tenant Mode** (Platform Integration).
- It finally prompts for a Port. If provided, the server spins up a Streamable HTTP MCP server; otherwise, it falls back to `stdio`.

---

## 3. How to Set Things Up

### Prerequisites
1. **Go 1.25+** installed on your machine.
2. A QuickBooks Developer Account.

### Step 1: Get QuickBooks Credentials (Single-Tenant)
1. Go to the [Intuit Developer Portal](https://developer.intuit.com/).
2. Create or select an app.
3. Grab your **Client ID** and **Client Secret**.
4. Use the OAuth Playground in the developer dashboard to generate a **Refresh Token** and **Realm ID** (Company ID) for your sandbox or production environment.

### Step 2: Configure the Environment
Create a `.env` file in the root of the project (or interactively enter these when the CLI wizard prompts you):

```env
# 🏢 SINGLE-TENANT (LOCAL) CONFIGURATION
QUICKBOOKS_CLIENT_ID=your_client_id_here
QUICKBOOKS_CLIENT_SECRET=your_client_secret_here
QUICKBOOKS_REALM_ID=your_realm_id_here
QUICKBOOKS_ENVIRONMENT=sandbox
QUICKBOOKS_REFRESH_TOKEN=your_refresh_token_here

# 🏢 MULTI-TENANT CONFIGURATION (Leave blank for local setup)
ENFORCE_AUTH=false
PLATFORM_INT_URL=

# 🌐 HTTP SERVER SETTINGS
# PORT=3000
MCP_HTTP_PATH=/mcp
```

### Step 3: Build the Source Code (If altering the code)
```bash
# Build the standard binary
go build -o qbo-mcp ./cmd/qbo-mcp

# Build a Linux-specific binary (e.g., for Docker)
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o qbo-mcp ./cmd/qbo-mcp
```

### Step 4: Run the Server
```bash
./qbo-mcp
```
Follow the CLI Wizard prompts. If `PORT` is omitted, the server runs in `stdio` mode (listening for raw JSON-RPC on standard input). If you input a port, it will start an HTTP server.

---

## 4. Understanding the Windows Binary

For Windows users or developers testing on Windows systems, a pre-compiled executable is included in the project root: **`qbo-mcp-windows.exe`**.

**Details:**
- **Architecture:** 64-bit (`x86-64`)
- **Format:** Native PE32+ console executable for Microsoft Windows.
- **Portability:** This binary was statically compiled and requires no additional dependencies (like a Go runtime) to execute on a Windows machine.
- **Linux Execution:** If you are developing on a Linux host (like Ubuntu) and need to test it, you can execute it natively using the `wine` compatibility layer (`wine ./qbo-mcp-windows.exe`).

---

## 5. Adding New Features (Developer Guide)

If you need to add a new QuickBooks endpoint to the MCP server:

1. **Create the Entity Handler:** Inside `internal/qbo/`, create a new file (e.g., `item_read.go`).
2. **Define the arguments schema:** Map the QuickBooks JSON payload structure to the expected MCP Tool input.
3. **Write the API Request:** Use the internal authenticated `client` to make a standard Go `http.NewRequest` to the Intuit API.
4. **Register the Tool:** Open `internal/tools/register.go` and add your new function to the master `RegisterAll()` list so the AI Agent becomes aware of it immediately upon connection.

*You are now equipped to run, test, and expand the QuickBooks Online MCP integration!*
