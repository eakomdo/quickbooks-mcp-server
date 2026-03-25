# QuickBooks Online MCP Server — Documentation

---

## 1. Summary

The **QuickBooks Online MCP Server** is a [Model Context Protocol (MCP)](https://modelcontextprotocol.io) server that exposes QuickBooks Online accounting operations as structured tools that AI agents can call. It acts as a bridge between an AI agent (like Claude) and the QuickBooks Online API, translating natural-language-driven tool calls into authenticated QuickBooks API requests.

The server supports two transport modes:

- **Stdio** — for local AI clients that speak MCP over stdin/stdout (e.g. Claude Desktop).
- **HTTP** — for platform-based deployments where the server listens on a port and handles requests from AI agent orchestration services.

It ships with **40+ tools** across 10 accounting domains: Customers, Invoices, Bills, Bill Payments, Estimates, Vendors, Chart of Accounts, Items, Employees, Journal Entries, and Purchases.

---

## 2. Language & Stack

| Concern | Technology |
|---|---|
| Runtime | Go 1.25+ |
| Language | Go |
| MCP Framework | `github.com/modelcontextprotocol/go-sdk` |
| QuickBooks API | Native HTTP client |
| OAuth 2.0 | Native Go implementation |
| HTTP Server | Go's `net/http` |
| Configuration | Environment variables with `godotenv` |

**Build:**
```bash
go build -o qbo-mcp ./cmd/qbo-mcp
```

**Run:**
```bash
# Stdio mode (default — used by Claude Desktop / MCP clients)
./qbo-mcp

# HTTP mode
PORT=3000 ./qbo-mcp
```

---

## 3. Tools Reference

Tools are grouped by domain. Every tool follows the same response envelope:

```json
{
  "content": [
    { "type": "text", "text": "..." }
  ]
}
```

Errors are returned as text with the prefix `Error <action>: <message>`.

---

### Search Criteria Format

Many tools accept a `criteria` parameter. Two forms are supported:

**Simple key-value object:**
```json
{ "criteria": { "Active": true, "DisplayName": "Acme Corp" } }
```

**Advanced filter array (with operators):**
```json
{
  "criteria": {
    "filters": [
      { "field": "TotalAmt", "value": 500, "operator": ">=" },
      { "field": "TxnDate",  "value": "2025-01-01", "operator": ">=" }
    ],
    "asc": "TxnDate",
    "limit": 20,
    "offset": 0
  }
}
```

Valid operators: `=`, `IN`, `<`, `>`, `<=`, `>=`, `LIKE`

Pagination fields available on all search tools: `limit`, `offset`, `asc`, `desc`, `count`, `fetchAll`.

---

### Customer Tools

| Tool Name | Description |
|---|---|
| `create_customer` | Create a new customer in QuickBooks Online |
| `get_customer` | Retrieve a customer by their QuickBooks ID |
| `update_customer` | Update fields on an existing customer |
| `delete_customer` | Deactivate / delete a customer |
| `search_customers` | Search customers using criteria and pagination |

#### `create_customer`
```go
// Input schema
{
  customer: any   // Raw QuickBooks Customer object
}
```

#### `get_customer`
```go
{ id: string }   // QuickBooks entity ID
```

#### `update_customer`
```go
{
  id: string,
  patch: Record<string, any>   // Fields to update
}
```

#### `delete_customer`
```go
{ idOrEntity: any }   // ID string or full Customer object
```

#### `search_customers`
```go
{
  criteria?: Array<any>,   // Filter criteria
  limit?:    number,
  offset?:   number,
  asc?:      string,       // Sort ascending by field
  desc?:     string,       // Sort descending by field
  count?:    boolean,      // Return count only
  fetchAll?: boolean       // Fetch all pages
}
```

---

### Invoice Tools

| Tool Name | Description |
|---|---|
| `read_invoice` | Get a single invoice by ID |
| `create_invoice` | Create a new invoice with line items |
| `update_invoice` | Update an existing invoice |
| `search_invoices` | Search invoices with field-validated criteria |

#### `create_invoice`
```go
{
  customer_ref: string,          // Customer QuickBooks ID (required)
  line_items: Array<{
    item_ref:    string,         // Item QuickBooks ID (required)
    qty:         number,         // Positive quantity (required)
    unit_price:  number,         // Non-negative price (required)
    description: string          // Optional line description
  }>,                            // At least one line item required
  doc_number?: string,           // Invoice number (optional)
  txn_date?:   string            // Date string YYYY-MM-DD (optional)
}
```

**Example:**
```json
{
  "customer_ref": "1",
  "line_items": [
    { "item_ref": "5", "qty": 2, "unit_price": 150.00, "description": "Consulting" }
  ],
  "txn_date": "2025-06-01"
}
```

#### `search_invoices`

Filterable fields: `Id`, `MetaData.CreateTime`, `MetaData.LastUpdatedTime`, `DocNumber`, `TxnDate`, `DueDate`, `CustomerRef`, `ClassRef`, `DepartmentRef`, `Balance`, `TotalAmt`

Sortable fields: `Id`, `MetaData.CreateTime`, `MetaData.LastUpdatedTime`, `DocNumber`, `TxnDate`, `Balance`, `TotalAmt`

```go
{
  criteria: {
    filters?: Array<{ field: string, value: any, operator?: string }>,
    asc?:     string,
    desc?:    string,
    limit?:   number,
    offset?:  number
  }
}
```

---

### Bill Tools

| Tool Name | Description |
|---|---|
| `create-bill` | Create a vendor bill |
| `get_bill` | Retrieve a bill by ID |
| `update_bill` | Update a bill |
| `delete_bill` | Delete a bill |
| `search_bills` | Search bills with criteria |

#### `create-bill`
```go
{
  bill: {
    Line:      Array<any>,          // Line items
    VendorRef: { value: string },   // Vendor ID
    DueDate?:  string,
    Balance?:  number,
    TotalAmt?: number
  }
}
```

---

### Bill Payment Tools

| Tool Name | Description |
|---|---|
| `create_bill_payment` | Record a payment against a bill |
| `get_bill_payment` | Get a bill payment by ID |
| `update_bill_payment` | Update a bill payment |
| `delete_bill_payment` | Delete a bill payment |
| `search_bill_payments` | Search bill payments |

#### `create_bill_payment`
```go
{ billPayment: any }   // Raw QuickBooks BillPayment object
```

---

### Estimate Tools

| Tool Name | Description |
|---|---|
| `create_estimate` | Create a sales estimate / quote |
| `get_estimate` | Get an estimate by ID |
| `update_estimate` | Update an estimate |
| `delete_estimate` | Delete an estimate |
| `search_estimates` | Search estimates with criteria |

#### `create_estimate`
```go
{ estimate: any }   // Raw QuickBooks Estimate object
```

---

### Vendor Tools

| Tool Name | Description |
|---|---|
| `create-vendor` | Create a new vendor |
| `get_vendor` | Get a vendor by ID |
| `update_vendor` | Update a vendor |
| `delete_vendor` | Delete / deactivate a vendor |
| `search_vendors` | Search vendors with criteria |

#### `create-vendor`
```go
{
  vendor: {
    DisplayName:      string,                   // Required
    GivenName?:       string,
    FamilyName?:      string,
    CompanyName?:     string,
    PrimaryEmailAddr?: { Address?: string },
    PrimaryPhone?:    { FreeFormNumber?: string },
    BillAddr?: {
      Line1?:                  string,
      City?:                   string,
      Country?:                string,
      CountrySubDivisionCode?: string,
      PostalCode?:             string
    }
  }
}
```

---

### Chart of Accounts Tools

| Tool Name | Description |
|---|---|
| `create_account` | Create a new account in the Chart of Accounts |
| `update_account` | Update an existing account |
| `search_accounts` | Search accounts with criteria |

Filterable fields: `Id`, `MetaData.CreateTime`, `MetaData.LastUpdatedTime`, `Name`, `SubAccount`, `ParentRef`, `Description`, `Active`, `Classification`, `AccountType`, `CurrentBalance`

#### `create_account`
```go
{
  name:         string,
  type:         string,    // QuickBooks AccountType (e.g. "Income", "Expense")
  sub_type?:    string,
  description?: string
}
```

---

### Item Tools

| Tool Name | Description |
|---|---|
| `read_item` | Get an item/product by ID |
| `create_item` | Create a new item |
| `update_item` | Update an item |
| `search_items` | Search items with criteria |

#### `create_item`
```go
{ item: any }   // Raw QuickBooks Item object
```

---

### Employee Tools

| Tool Name | Description |
|---|---|
| `create_employee` | Create a new employee record in QuickBooks Online |
| `get_employee` | Retrieve an employee by their QuickBooks ID |
| `update_employee` | Update fields on an existing employee |
| `search_employees` | Search employees using criteria and pagination options |

#### `create_employee`
```go
{ employee: any }   // Raw QuickBooks Employee object
```

**Example payload:**
```json
{
  "employee": {
    "GivenName":    "Jane",
    "FamilyName":   "Smith",
    "DisplayName":  "Jane Smith",
    "PrimaryAddr": {
      "Line1":                  "123 Main St",
      "City":                   "Mountain View",
      "CountrySubDivisionCode": "CA",
      "PostalCode":             "94043"
    },
    "PrimaryPhone": { "FreeFormNumber": "555-1234" },
    "SSN":          "XXX-XX-XXXX"
  }
}
```

#### `get_employee`
```go
{ id: string }   // QuickBooks employee ID
```

#### `update_employee`
```go
{
  id:    string,
  patch: Record<string, any>   // Fields to update
}
```

#### `search_employees`
```go
{
  criteria?: Array<any>,   // Filter criteria
  asc?:      string,       // Sort ascending by field name
  desc?:     string,       // Sort descending by field name
  limit?:    number,       // Max results to return
  offset?:   number,       // Skip N results
  count?:    boolean,      // Return only the count
  fetchAll?: boolean       // Fetch all pages automatically
}
```

**Example — find all active employees:**
```json
{
  "criteria": [{ "field": "Active", "value": true, "operator": "=" }],
  "limit": 50
}
```

---

### Journal Entry Tools

| Tool Name | Description |
|---|---|
| `create_journal_entry` | Create a double-entry journal entry in QuickBooks Online |
| `get_journal_entry` | Retrieve a journal entry by its QuickBooks ID |
| `update_journal_entry` | Update an existing journal entry |
| `delete_journal_entry` | Delete a journal entry |
| `search_journal_entries` | Search journal entries using filter criteria |

#### `create_journal_entry`
```go
{ journalEntry: any }   // Raw QuickBooks JournalEntry object
```

**Example — simple debit/credit entry:**
```json
{
  "journalEntry": {
    "TxnDate": "2025-06-01",
    "PrivateNote": "Monthly accrual",
    "Line": [
      {
        "Description":     "Accrued expense",
        "Amount":          500.00,
        "DetailType":      "JournalEntryLineDetail",
        "JournalEntryLineDetail": {
          "PostingType": "Debit",
          "AccountRef":  { "value": "7" }
        }
      },
      {
        "Description":     "Accrued liability",
        "Amount":          500.00,
        "DetailType":      "JournalEntryLineDetail",
        "JournalEntryLineDetail": {
          "PostingType": "Credit",
          "AccountRef":  { "value": "33" }
        }
      }
    ]
  }
}
```

#### `get_journal_entry`
```go
{ id: string }   // QuickBooks journal entry ID
```

#### `update_journal_entry`
```go
{
  id:    string,
  patch: Record<string, any>
}
```

#### `delete_journal_entry`
```go
{ idOrEntity: any }   // ID string or full JournalEntry object
```

#### `search_journal_entries`
```go
{
  criteria?: Array<any>,
  asc?:      string,
  desc?:     string,
  limit?:    number,
  offset?:   number,
  count?:    boolean,
  fetchAll?: boolean
}
```

---

### Purchase Tools

| Tool Name | Description |
|---|---|
| `create_purchase` | Create a purchase / expense transaction |
| `get_purchase` | Get a purchase by ID |
| `update_purchase` | Update a purchase |
| `delete_purchase` | Delete a purchase |
| `search_purchases` | Search purchases with criteria |

#### `create_purchase`
```go
{ purchase: any }   // Raw QuickBooks Purchase object
```

---

## 4. Authentication & Platform Integration

The server has two authentication layers: one for verifying incoming requests (who is calling the server) and one for authenticating with QuickBooks (OAuth 2.0).

---

### 4.1 Incoming Request Authentication — JWT Verification

When running in HTTP mode with `ENFORCE_AUTH=true`, every request must include a Bearer token:

```
Authorization: Bearer <jwt-token>
```

The server validates the token in [src/auth.ts](src/auth.ts) using `jose`. Two JWT versions are supported:

**Version 1 — Symmetric key (HS256):**
```go
// Token payload contains version=1 and client_id
// Key comes from SECRET_KEY env var keyed by client_id
const result = await jwtVerify(token, encoder.encode(KEYS[client_id]));
```

**Version 2 — Asymmetric key (RS256/ES256 via JWKS):**
```go
// Token payload contains version=2 and a kid header
// Public key fetched from ACCOUNT_SERVICE_URL + ACCOUNT_SERVICE_JWKS_ENDPOINT
const publicKey = await getPublicKey(kid);
const result = await jwtVerify(token, publicKey, { audience: SERVICE_ID });
```

JWKS responses are cached in-memory with a configurable TTL (`ACCOUNT_SERVICE_JWKS_CACHE_TTL`, default 600s) to avoid hitting the key server on every request.

On success, the token payload is decoded into a `TokenData` object:

```go
interface TokenData {
  id:           string;
  email:        string;
  username:     string | null;
  role:         string;
  type:         string;
  client_id:    string;
  access_token: string;
}
```

If verification fails, the server returns `401 Unauthorized` with a JSON error body.

---

### 4.2 QuickBooks OAuth 2.0 Authentication

The [QuickbooksClient](src/clients/quickbooks-client.ts) handles OAuth 2.0 against Intuit's identity platform. Every tool call goes through `authenticate()` before touching the API.

```
┌─────────────┐        authenticate()        ┌──────────────────┐
│   MCP Tool  │ ─────────────────────────►  │ QuickbooksClient │
└─────────────┘                              └────────┬─────────┘
                                                      │
                        ┌─────────────────────────────┤
                        │  Platform integration?       │
                        │  (PLATFORM_INT_URL set?)     │
                        └──────────┬──────────────────┘
                    Yes ◄──────────┤──────────► No
                                   │
             ┌─────────────────────┤              ┌──────────────────────┐
             │ loadFromPlatform()  │              │ Use env var creds     │
             │ fetch per-user creds│              │ QUICKBOOKS_CLIENT_ID  │
             │ from platform API   │              │ QUICKBOOKS_CLIENT_SECRET│
             └─────────┬───────────┘              └──────────────────────┘
                       │
            ┌──────────▼──────────────────────────────────────────┐
            │  refreshToken available?                             │
            │   Yes → refreshAccessToken() (auto-refresh)         │
            │   No  → startOAuthFlow() (open browser for consent) │
            └──────────────────────────────────────────────────────┘
```

**Token Refresh:**

```go
// Access token is refreshed automatically when expired
const now = new Date();
if (!this.accessToken || this.accessTokenExpiry <= now) {
  await this.refreshAccessToken();
}
```

If the refresh token itself has expired (HTTP 400 from Intuit), the client automatically re-triggers the browser OAuth flow to obtain a fresh token pair.

**Saving tokens:** After a successful OAuth flow, tokens are written back to `.env` and (if platform integration is active) synced to the platform service via `PATCH /api/v1/quickbooks_projects/{id}`.

---

### 4.3 Platform Integration — Per-User Credentials

When `PLATFORM_INT_URL` is set, the server switches to **per-user credential mode**. Instead of using a single set of QuickBooks credentials from `.env`, it fetches the calling user's credentials from the platform integration service on every request.

This enables multi-tenant deployments where each user has their own connected QuickBooks company.

**Flow:**

```
HTTP Request (Authorization: Bearer <user-jwt>)
         │
         ▼
  Extract & verify JWT
         │
         ▼
  Store token in AsyncLocalStorage (requestContext)
         │
         ▼
  Tool call → quickbooksClient.authenticate()
         │
         ▼
  PlatformIntegrationClient.getQuickbooksProject(userToken)
  GET {PLATFORM_INT_URL}/api/v1/quickbooks_projects?limit=1
         │
         ▼
  Returns: { client_id, client_secret, refresh_token, realmid, ... }
         │
         ▼
  Re-initialize OAuthClient with per-user credentials
         │
         ▼
  Execute QuickBooks API call
```

**Platform client** ([src/clients/platform-integration-client.ts](src/clients/platform-integration-client.ts)):

```go
// Fetch user's QB project
const project = await platformClient.getQuickbooksProject(userToken);
// project: { id, client_id, client_secret, refresh_token, realmid, ... }

// Sync updated refresh token back to platform
await platformClient.updateRefreshToken(project.id, newRefreshToken, userToken);
```

The `userToken` is threaded through requests using Node.js `AsyncLocalStorage` ([src/clients/request-context.ts](src/clients/request-context.ts)), so it is available inside tool handlers without passing it as a parameter:

```go
// Set on incoming request
requestContext.run({ userToken: bearerToken }, () => {
  // handle MCP request
});

// Read inside QuickbooksClient
const userToken = requestContext.getStore()?.userToken;
```

---

### 4.4 Environment Variables (Saved in Platform Integration Service)

| Variable | Required | Description |
|---|---|---|
| `QUICKBOOKS_CLIENT_ID` | Yes* | Intuit app client ID |
| `QUICKBOOKS_CLIENT_SECRET` | Yes* | Intuit app client secret |
| `QUICKBOOKS_ENVIRONMENT` | No | `sandbox` (default) or `production` |
| `QUICKBOOKS_REFRESH_TOKEN` | No | Pre-obtained refresh token |
| `QUICKBOOKS_REALM_ID` | No | QuickBooks company realm ID |
| `PLATFORM_INT_URL` | No | Enables per-user credential mode |
| `ACCOUNT_SERVICE_URL` | No | Base URL of JWT issuer (for v2 tokens) |
| `ACCOUNT_SERVICE_JWKS_ENDPOINT` | No | JWKS path (default `/.well-known/jwks.json`) |
| `ACCOUNT_SERVICE_JWKS_CACHE_TTL` | No | JWKS cache seconds (default `600`) |
| `SECRET_KEY` | No | Symmetric key for v1 JWT verification |
| `ENFORCE_AUTH` | No | Set to `true` to require Bearer token on HTTP |
| `PORT` | No | HTTP server port (enables HTTP mode) |
| `MCP_TRANSPORT` | No | Set to `http` to force HTTP mode |
| `MCP_HTTP_PATH` | No | HTTP endpoint path (default `/mcp`) |
| `SERVICE_ID` | No | JWT audience claim (default `quickboos-mcp`) |
| `LICENSE_KEY` | No | License key for the license watcher |
| `USAGE_REPORT_ENDPOINT` | No | URL to POST usage telemetry |

*Required unless `PLATFORM_INT_URL` is set.

---

## 5. Architecture Overview

```
┌──────────────────────────────────────────────────────────────────┐
│                        AI Agent / Twynity                          │
└────────────────────────┬─────────────────────────────────────────┘
                         │  MCP protocol (stdio or HTTP)
                         ▼
┌──────────────────────────────────────────────────────────────────┐
│                      src/index.ts (Entry Point)                   │
│                                                                    │
│   ┌────────────────────────────────────────────────────────────┐  │
│   │                  registerAllTools()                         │  │
│   │   49 tools registered via RegisterTool() helper            │  │
│   └────────────────────────────────────────────────────────────┘  │
│                                                                    │
│   ┌──────────────────┐        ┌───────────────────────────────┐   │
│   │   Stdio mode      │        │         HTTP mode              │   │
│   │  StdioTransport   │        │  StreamableHTTPTransport       │   │
│   │  (one server)     │        │  (new server per request)      │   │
│   └──────────────────┘        └───────────────────────────────┘   │
└──────────────────────────────────────────────────────────────────┘
                         │
              ┌──────────▼──────────┐
              │  QuickbooksMCPServer │  (src/server/qbo-mcp-server.ts)
              │   McpServer factory  │
              └──────────┬──────────┘
                         │
              ┌──────────▼──────────┐
              │     Tool Handler     │  (src/handlers/*.handler.ts)
              └──────────┬──────────┘
                         │
              ┌──────────▼──────────┐
              │  QuickbooksClient    │  (src/clients/quickbooks-client.ts)
              │  authenticate()      │
              │  getQuickbooks()     │
              └──────────┬──────────┘
                         │
         ┌───────────────┴───────────────┐
         │                               │
┌────────▼────────┐           ┌──────────▼──────────┐
│  intuit-oauth   │           │  PlatformClient      │
│  OAuth 2.0 flow │           │  (per-user creds)    │
└────────┬────────┘           └──────────────────────┘
         │
┌────────▼────────┐
│  node-quickbooks│
│  QBO REST API   │
└─────────────────┘
```

### Tool Registration

Every tool is defined as a `ToolDefinition` object with a name, description, Zod schema, and handler:

```go
// src/helpers/register-tool.ts
function RegisterTool<T extends z.ZodType>(server: McpServer, tool: ToolDefinition<T>) {
  server.tool(
    tool.name,
    tool.description,
    { params: tool.schema },
    tool.handler
  );
}
```

All 49 tools are registered at startup in `registerAllTools()` before the transport connects.


## 6. Sample MCP Requests — Create Endpoints

All `customer`, `estimate`, `bill`, `billPayment`, `vendor`, `item`, `employee`, `journalEntry`, and `purchase` fields must be passed as **plain JSON objects** (not stringified). Fields marked as read-only by QuickBooks (`FullyQualifiedName`, `Balance`, `TotalAmt` where auto-computed, address `Id`/`Lat`/`Long`, and line `Id`) are omitted.

---

### `create_customer`

```json
{
  "method": "tools/call",
  "params": {
    "name": "create_customer",
    "arguments": {
      "params": {
        "customer": {
          "DisplayName": "King's Groceries",
          "CompanyName": "King Groceries",
          "Title": "Mr",
          "GivenName": "James",
          "MiddleName": "B",
          "FamilyName": "King",
          "Suffix": "Jr",
          "Notes": "Here are other details.",
          "PrimaryEmailAddr": {
            "Address": "jdrew@myemail.com"
          },
          "PrimaryPhone": {
            "FreeFormNumber": "(555) 555-5555"
          },
          "BillAddr": {
            "Line1": "123 Main Street",
            "City": "Mountain View",
            "CountrySubDivisionCode": "CA",
            "PostalCode": "94042",
            "Country": "USA"
          }
        }
      }
    }
  }
}
```

---

### `create_estimate`

```json
{
  "method": "tools/call",
  "params": {
    "name": "create_estimate",
    "arguments": {
      "params": {
        "estimate": {
          "CustomerRef": {
            "value": "3",
            "name": "Cool Cars"
          },
          "BillEmail": {
            "Address": "Cool_Cars@intuit.com"
          },
          "CustomerMemo": {
            "value": "Thank you for your business and have a great day!"
          },
          "PrintStatus": "NeedToPrint",
          "EmailStatus": "NotSet",
          "ApplyTaxAfterDiscount": false,
          "BillAddr": {
            "Line1": "65 Ocean Dr.",
            "City": "Half Moon Bay",
            "CountrySubDivisionCode": "CA",
            "PostalCode": "94213"
          },
          "ShipAddr": {
            "Line1": "65 Ocean Dr.",
            "City": "Half Moon Bay",
            "CountrySubDivisionCode": "CA",
            "PostalCode": "94213"
          },
          "TxnTaxDetail": {
            "TotalTax": 0
          },
          "Line": [
            {
              "LineNum": 1,
              "Description": "Pest Control Services",
              "Amount": 35.0,
              "DetailType": "SalesItemLineDetail",
              "SalesItemLineDetail": {
                "ItemRef": {
                  "value": "10",
                  "name": "Pest Control"
                },
                "Qty": 1,
                "UnitPrice": 35,
                "TaxCodeRef": {
                  "value": "NON"
                }
              }
            },
            {
              "Amount": 3.5,
              "DetailType": "DiscountLineDetail",
              "DiscountLineDetail": {
                "PercentBased": true,
                "DiscountPercent": 10,
                "DiscountAccountRef": {
                  "value": "86",
                  "name": "Discounts given"
                }
              }
            }
          ]
        }
      }
    }
  }
}
```

---

### `create_invoice`

> `create_invoice` uses a structured schema instead of a raw QuickBooks object. `customer_ref` and `item_ref` are QuickBooks entity IDs.

```json
{
  "method": "tools/call",
  "params": {
    "name": "create_invoice",
    "arguments": {
      "params": {
        "customer_ref": "3",
        "txn_date": "2026-03-04",
        "doc_number": "INV-001",
        "line_items": [
          {
            "item_ref": "10",
            "qty": 2,
            "unit_price": 150.00,
            "description": "Consulting Services"
          },
          {
            "item_ref": "5",
            "qty": 1,
            "unit_price": 75.00,
            "description": "Setup Fee"
          }
        ]
      }
    }
  }
}
```

---

### `create-bill`

> `VendorRef.value` is the QuickBooks vendor ID. `AccountRef.value` is the expense account ID. `DetailType` must be `"AccountBasedExpenseLineDetail"`.

```json
{
  "method": "tools/call",
  "params": {
    "name": "create-bill",
    "arguments": {
      "params": {
        "bill": {
          "VendorRef": {
            "value": "56",
            "name": "Bob's Hardware"
          },
          "DueDate": "2026-04-03",
          "Balance": 200.00,
          "TotalAmt": 200.00,
          "Line": [
            {
              "Amount": 200.00,
              "DetailType": "AccountBasedExpenseLineDetail",
              "Description": "Office supplies",
              "AccountRef": {
                "value": "7",
                "name": "Office Expenses"
              }
            }
          ]
        }
      }
    }
  }
}
```

---

### `create_bill_payment`

> `PayType` must be `"Check"` or `"CreditCard"`. `LinkedTxn.TxnId` is the ID of the bill being paid.

```json
{
  "method": "tools/call",
  "params": {
    "name": "create_bill_payment",
    "arguments": {
      "params": {
        "billPayment": {
          "VendorRef": {
            "value": "56",
            "name": "Bob's Hardware"
          },
          "PayType": "Check",
          "CheckPayment": {
            "BankAccountRef": {
              "value": "35",
              "name": "Checking"
            }
          },
          "TxnDate": "2026-03-04",
          "Line": [
            {
              "Amount": 200.00,
              "LinkedTxn": [
                {
                  "TxnId": "25",
                  "TxnType": "Bill"
                }
              ]
            }
          ]
        }
      }
    }
  }
}
```

---

### `create-vendor`

```json
{
  "method": "tools/call",
  "params": {
    "name": "create-vendor",
    "arguments": {
      "params": {
        "vendor": {
          "DisplayName": "Bob's Hardware",
          "CompanyName": "Bob's Hardware",
          "GivenName": "Bob",
          "FamilyName": "Johnson",
          "PrimaryEmailAddr": {
            "Address": "bob@hardware.com"
          },
          "PrimaryPhone": {
            "FreeFormNumber": "(555) 444-3333"
          },
          "BillAddr": {
            "Line1": "456 Supply Lane",
            "City": "Oakland",
            "CountrySubDivisionCode": "CA",
            "PostalCode": "94601",
            "Country": "US"
          }
        }
      }
    }
  }
}
```

---

### `create_account`

> `type` must be a valid QuickBooks AccountType (e.g. `"Income"`, `"Expense"`, `"Asset"`, `"Liability"`, `"Equity"`). `sub_type` narrows the classification.

```json
{
  "method": "tools/call",
  "params": {
    "name": "create_account",
    "arguments": {
      "params": {
        "name": "Office Supplies",
        "type": "Expense",
        "sub_type": "OfficeGeneralAdministrativeExpenses",
        "description": "General office supply purchases"
      }
    }
  }
}
```

---

### `create_item`

> `type` must be `"Service"`, `"Inventory"`, or `"NonInventory"`. `income_account_ref` is the QuickBooks account ID for tracking income from this item.

```json
{
  "method": "tools/call",
  "params": {
    "name": "create_item",
    "arguments": {
      "params": {
        "name": "Consulting Hour",
        "type": "Service",
        "income_account_ref": "79",
        "expense_account_ref": "80",
        "unit_price": 150.00,
        "description": "Professional consulting services per hour"
      }
    }
  }
}
```

---

### `create_employee`

```json
{
  "method": "tools/call",
  "params": {
    "name": "create_employee",
    "arguments": {
      "params": {
        "employee": {
          "GivenName": "Jane",
          "FamilyName": "Smith",
          "DisplayName": "Jane Smith",
          "PrimaryPhone": {
            "FreeFormNumber": "(555) 123-4567"
          },
          "PrimaryAddr": {
            "Line1": "123 Main St",
            "City": "Mountain View",
            "CountrySubDivisionCode": "CA",
            "PostalCode": "94043"
          },
          "BillableTime": false,
          "HiredDate": "2026-03-04"
        }
      }
    }
  }
}
```

---

### `create_journal_entry`

> Debits and credits must balance (total debit amount = total credit amount). `AccountRef.value` is the QuickBooks account ID.

```json
{
  "method": "tools/call",
  "params": {
    "name": "create_journal_entry",
    "arguments": {
      "params": {
        "journalEntry": {
          "TxnDate": "2026-03-04",
          "PrivateNote": "Monthly accrual",
          "Line": [
            {
              "Description": "Accrued expense",
              "Amount": 500.00,
              "DetailType": "JournalEntryLineDetail",
              "JournalEntryLineDetail": {
                "PostingType": "Debit",
                "AccountRef": {
                  "value": "7"
                }
              }
            },
            {
              "Description": "Accrued liability",
              "Amount": 500.00,
              "DetailType": "JournalEntryLineDetail",
              "JournalEntryLineDetail": {
                "PostingType": "Credit",
                "AccountRef": {
                  "value": "33"
                }
              }
            }
          ]
        }
      }
    }
  }
}
```

---

### `create_purchase`

> `PaymentType` must be `"Cash"`, `"Check"`, or `"CreditCard"`. `AccountRef.value` is the payment account (bank/credit card). `EntityRef` links the purchase to a vendor or customer.

```json
{
  "method": "tools/call",
  "params": {
    "name": "create_purchase",
    "arguments": {
      "params": {
        "purchase": {
          "PaymentType": "Check",
          "TxnDate": "2026-03-04",
          "AccountRef": {
            "value": "35",
            "name": "Checking"
          },
          "EntityRef": {
            "value": "56",
            "type": "Vendor",
            "name": "Bob's Hardware"
          },
          "Line": [
            {
              "DetailType": "AccountBasedExpenseLineDetail",
              "Amount": 50.00,
              "Description": "Office supplies purchase",
              "AccountBasedExpenseLineDetail": {
                "AccountRef": {
                  "value": "7",
                  "name": "Office Expenses"
                }
              }
            }
          ]
        }
      }
    }
  }
}
```