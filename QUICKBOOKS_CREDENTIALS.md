# QuickBooks Online API Credentials - Setup Guide

## Where to Get Your Credentials

curl -X POST http://localhost:3000/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}'### Step 1: Create/Access Intuit Developer Portal
1. Go to: https://developer.intuit.com/
2. Sign in with your Intuit account (create one if needed)
3. Go to **My Apps** → **Create an app**

### Step 2: Create a QuickBooks Online App
1. Choose: **QuickBooks Online + Payments**
2. App Type: **Web app**
3. Give it a name (e.g., "QuickBooks MCP Server")
4. Click **Create**

### Step 3: Get Your Credentials
After creating the app, you'll see these in your app settings:

**Copy these 5 values:**

```
QUICKBOOKS_CLIENT_ID=<Your Client ID>
QUICKBOOKS_CLIENT_SECRET=<Your Client Secret>
QUICKBOOKS_ENVIRONMENT=sandbox  (or 'production' for live)
QUICKBOOKS_REALM_ID=<Your Realm ID>
QUICKBOOKS_REFRESH_TOKEN=<Your Refresh Token>
```

### Step 4: OAuth Flow to Get Refresh Token & Realm ID

If you don't have these yet, you need to authenticate:

1. **Get Authorization Code:**
   - Visit this URL (replace YOUR_CLIENT_ID and REDIRECT_URI):
   ```
   https://appcenter.intuit.com/connect/oauth2?client_id=YOUR_CLIENT_ID&response_type=code&scope=com.intuit.quickbooks.accounting&redirect_uri=http://localhost:3000/callback&state=1
   ```

2. **Exchange for Tokens:**
   ```bash
   curl -X POST https://oauth.platform.intuit.com/oauth2/v1/tokens/bearer \
     -H "Content-Type: application/x-www-form-urlencoded" \
     -d "grant_type=authorization_code" \
     -d "code=AUTHORIZATION_CODE" \
     -d "redirect_uri=http://localhost:3000/callback" \
     -d "client_id=YOUR_CLIENT_ID" \
     -d "client_secret=YOUR_CLIENT_SECRET"
   ```

3. **Response will contain:**
   - `access_token` (temporary, refreshed automatically)
   - `refresh_token` (long-lived, save this!)
   - `realm_id` (your QuickBooks company ID)

## Where to Store Credentials

### Option A: Environment File (Recommended for Testing)
Create `.env` file in `/home/emma/Downloads/quickbooks-online-mcp/`:

```env
QUICKBOOKS_CLIENT_ID=your_client_id_here
QUICKBOOKS_CLIENT_SECRET=your_client_secret_here
QUICKBOOKS_REALM_ID=your_realm_id_here
QUICKBOOKS_REFRESH_TOKEN=your_refresh_token_here
QUICKBOOKS_ENVIRONMENT=sandbox
```

Then run:
```bash
cd /home/emma/Downloads/quickbooks-online-mcp
export $(cat .env | xargs)
./qbo-mcp
```

### Option B: System Environment Variables (Better for Production)
```bash
export QUICKBOOKS_CLIENT_ID=your_client_id_here
export QUICKBOOKS_CLIENT_SECRET=your_client_secret_here
export QUICKBOOKS_REALM_ID=your_realm_id_here
export QUICKBOOKS_REFRESH_TOKEN=your_refresh_token_here
export QUICKBOOKS_ENVIRONMENT=sandbox

./qbo-mcp
```

### Option C: Docker Environment (Production)
Use in docker-compose.yaml or `docker run -e`:

```bash
docker run \
  -e QUICKBOOKS_CLIENT_ID=your_client_id \
  -e QUICKBOOKS_CLIENT_SECRET=your_client_secret \
  -e QUICKBOOKS_REALM_ID=your_realm_id \
  -e QUICKBOOKS_REFRESH_TOKEN=your_refresh_token \
  -e QUICKBOOKS_ENVIRONMENT=sandbox \
  -e PORT=3000 \
  -e MCP_TRANSPORT=http \
  -p 3000:3000 \
  qbo-mcp-server:latest
```

## Security Best Practices

⚠️ **NEVER commit credentials to git!**

1. Add `.env` to `.gitignore` (already done ✓)
2. Use environment variables in CI/CD systems
3. Use secrets management (HashiCorp Vault, AWS Secrets Manager, Azure Key Vault)
4. Rotate tokens regularly
5. Use sandbox environment for development

## Testing Your Credentials

Once you have credentials, test with:

```bash
# Run the server with your credentials
./qbo-mcp

# In another terminal, test a tool:
curl -X POST http://localhost:3000/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "tools/call",
    "params": {
      "name": "quickbook_search_customers",
      "arguments": {
        "params": {
          "criteria": { "limit": 5 }
        }
      }
    }
  }'
```

## Common Issues

| Issue | Solution |
|-------|----------|
| `401 Unauthorized` | Token expired - get a new refresh token |
| `invalid_client` | Check CLIENT_ID and CLIENT_SECRET |
| `invalid_realm` | Check REALM_ID matches your QB company |
| `Invalid oauth state` | Realm ID doesn't match the session |

## Resources

- [Intuit Developer Portal](https://developer.intuit.com/)
- [QuickBooks Online API Docs](https://developer.intuit.com/app/developer/qbo/docs/api)
- [OAuth 2.0 Guide](https://developer.intuit.com/app/developer/qbo/docs/develop/authentication-and-authorization/oauth-2)
- [Refresh Token Expiry Info](https://developer.intuit.com/app/developer/qbo/docs/develop/authentication-and-authorization/oauth-2#oauth-20-refresh-tokens)
