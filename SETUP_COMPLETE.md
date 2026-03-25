# 🚀 Complete Setup Guide - QuickBooks MCP Server

## Overview
You now have three complete setups ready:

1. ✅ **Claude Desktop Integration** - Use tools directly in Claude
2. 📝 **QB Credentials** - Guide to get real API credentials
3. 🐳 **Docker Deployment** - Production-ready containerization

---

## 1️⃣ Claude Desktop Integration ✅ DONE

### What Was Set Up
- **Location:** `~/.config/claude/claude.json`
- **Status:** Ready to use
- **Content:** MCP server configuration pointing to `/home/emma/Downloads/quickbooks-online-mcp/qbo-mcp`

### Steps to Activate

#### Option A: Start Using (Quickest)
1. **Restart Claude Desktop** (if running)
2. **Wait 5 seconds** for it to load tools
3. **Open a new conversation**
4. **Ask Claude to help with QuickBooks:**
   ```
   "Create a customer named 'Acme Corp' in QuickBooks"
   "Search for all invoices over $1000"
   "Update customer John Smith's email"
   ```

#### Option B: Verify Configuration
```bash
# Check config exists
cat ~/.config/claude/claude.json

# Expected output:
# {
#   "version": "1",
#   "tools": [
#     {
#       "name": "qbo-mcp",
#       ...
#     }
#   ]
# }
```

#### Option C: Run Server Manually First (For Testing)
```bash
cd /home/emma/Downloads/quickbooks-online-mcp
./qbo-mcp
```

### Accessing Tools in Claude
Once configured, you can ask Claude:
- "What customers do we have in QuickBooks?"
- "Create a new invoice for $5000"
- "List all bills due this month"
- "Update the address for customer XYZ"

Claude will automatically use the MCP tools to handle these requests.

---

## 2️⃣ QuickBooks Credentials 📝 

### Current Status
- **File Created:** `QUICKBOOKS_CREDENTIALS.md` (in repo root)
- **Test Credentials:** Not yet added (optional for testing)
- **Real Credentials:** Need to get from Intuit Developer Portal

### Getting Real Credentials (Required for Production)

**Follow these steps:**

1. **Go to Intuit Developer Portal**
   ```
   https://developer.intuit.com/
   ```

2. **Create/Select App**
   - Sign in
   - Go to "My Apps"
   - Click "Create an app"
   - Choose "QuickBooks Online + Payments"
   - Give it a name

3. **Copy These 5 Values:**
   ```
   Client ID:      (from Settings → Keys & Credentials)
   Client Secret:  (from Settings → Keys & Credentials)
   Realm ID:       (your QB company ID)
   Refresh Token:  (get via OAuth flow)
   Environment:    sandbox (or production)
   ```

4. **Create `.env` File**
   ```bash
   cd /home/emma/Downloads/quickbooks-online-mcp
   cat > .env << 'EOF'
   QUICKBOOKS_CLIENT_ID=your_id_here
   QUICKBOOKS_CLIENT_SECRET=your_secret_here
   QUICKBOOKS_REALM_ID=your_realm_here
   QUICKBOOKS_REFRESH_TOKEN=your_token_here
   QUICKBOOKS_ENVIRONMENT=sandbox
   EOF
   ```

5. **Test Connection**
   ```bash
   export $(cat .env | xargs)
   ./qbo-mcp
   ```

### For Detailed Instructions
See: `QUICKBOOKS_CREDENTIALS.md`
- OAuth flow details
- Token refresh explanation
- Security best practices
- Troubleshooting guide

---

## 3️⃣ Docker Deployment 🐳

### Current Status
- **Dockerfile:** ✅ Production-ready
- **Docker Compose:** ✅ Configured
- **Image:** Ready to build

### Quick Start: Build & Deploy Locally

#### Step 1: Build Docker Image
```bash
cd /home/emma/Downloads/quickbooks-online-mcp

# Build
docker build -t qbo-mcp-server:latest .

# Verify
docker images | grep qbo-mcp
# Output: qbo-mcp-server    latest    <image-id>    <size>
```

#### Step 2: Create .env for Docker
```bash
cd /home/emma/Downloads/quickbooks-online-mcp
cat > .env << 'EOF'
QUICKBOOKS_CLIENT_ID=test_client_id
QUICKBOOKS_CLIENT_SECRET=test_client_secret
QUICKBOOKS_REALM_ID=test_realm_id
QUICKBOOKS_REFRESH_TOKEN=test_refresh_token
QUICKBOOKS_ENVIRONMENT=sandbox
PORT=3000
MCP_TRANSPORT=http
EOF
```

#### Step 3: Run with Docker Compose
```bash
# Start
docker compose up -d

# Check logs
docker compose logs -f

# Test
curl http://localhost:3000/mcp

# Stop
docker compose down
```

#### Step 4: Verify It Works
```bash
# Check running
docker ps | grep qbo-mcp

# Expected: Container running on port 3000

# Test endpoint
curl -X POST http://localhost:3000/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}'

# Should return: {"jsonrpc":"2.0"...}
```

### Deployment Options

| Option | When to Use | Effort | Cost |
|--------|------------|--------|------|
| **Local Docker** | Development/Testing | ⭐ Easiest | Free |
| **Docker Compose** | Small Teams | ⭐⭐ Easy | Free |
| **Docker Hub** | Share publicly | ⭐⭐ Medium | Free/Paid |
| **Azure Container Registry** | Enterprise Azure | ⭐⭐⭐ Hard | Paid |
| **Kubernetes** | Large Scale | ⭐⭐⭐⭐ Complex | Paid |

### Full Deployment Guide
See: `DOCKER_DEPLOYMENT.md`
- Local testing
- Docker Compose setup
- Production configurations
- Azure Container Registry
- Kubernetes deployment
- Monitoring & logging

---

## 📊 Testing Checklist

- [ ] **Claude Desktop**
  - [ ] Config file created at `~/.config/claude/claude.json`
  - [ ] Verified by restarting Claude
  - [ ] Can ask Claude about QB operations

- [ ] **QB Credentials**
  - [ ] Got credentials from Intuit Developer Portal
  - [ ] Created `.env` file
  - [ ] Tested with: `export $(cat .env | xargs) && ./qbo-mcp`

- [ ] **Docker**
  - [ ] Built image: `docker build -t qbo-mcp-server:latest .`
  - [ ] Ran with compose: `docker compose up -d`
  - [ ] Verified endpoint: `curl http://localhost:3000/mcp`
  - [ ] Stopped cleanly: `docker compose down`

---

## 🚀 Next Steps

### For Development
```bash
# 1. Get QB credentials from Intuit
# 2. Create .env file in repo
# 3. Run server
./qbo-mcp

# 4. Use in Claude Desktop immediately
# (auto-loaded from ~/.config/claude/claude.json)
```

### For Deployment
```bash
# 1. Build Docker image
docker build -t qbo-mcp-server:1.0.0 .

# 2. Push to registry (if needed)
docker tag qbo-mcp-server:1.0.0 myregistry/qbo-mcp-server:1.0.0
docker push myregistry/qbo-mcp-server:1.0.0

# 3. Deploy with compose or orchestrator
docker compose -f docker-compose.yaml up -d
```

### For Production
```bash
# 1. Set up secrets management
# 2. Use environment variables for credentials
# 3. Enable HTTPS with reverse proxy
# 4. Set up monitoring & logging
# 5. Deploy to Kubernetes or cloud
```

---

## 📚 Quick Reference

| Component | File | Status | Next Action |
|-----------|------|--------|-------------|
| Claude Config | `~/.config/claude/claude.json` | ✅ Done | Restart Claude |
| QB Credentials | `QUICKBOOKS_CREDENTIALS.md` | Created | Follow guide to get credentials |
| Docker Build | `Dockerfile` | ✅ Ready | `docker build -t qbo-mcp-server:latest .` |
| Docker Compose | `compose.yaml` | ✅ Ready | Add `.env` then `docker compose up -d` |
| Deployment Guide | `DOCKER_DEPLOYMENT.md` | Created | Choose deployment option |

---

## ❓ FAQ

**Q: Do I need all three?**
A: No. Use Claude Desktop for immediate access, skip Docker if not deploying.

**Q: Can I use the binary directly?**
A: Yes! `./qbo-mcp` or `PORT=3000 ./qbo-mcp` for HTTP mode.

**Q: Where do I get QB credentials?**
A: https://developer.intuit.com/ - Create an app and follow OAuth flow.

**Q: Is Docker required for Claude Desktop?**
A: No - Claude Desktop uses the binary directly from the filesystem.

**Q: How often do I need to refresh the token?**
A: Automatically! The server handles refresh_token rotation.

**Q: Can I use production QB credentials?**
A: Yes, just set `QUICKBOOKS_ENVIRONMENT=production` instead of `sandbox`.

---

## 🎉 You're All Set!

Your QuickBooks MCP server is ready to:
- ✅ Integrate with Claude Desktop immediately
- ✅ Connect to real QB data (once you get credentials)
- ✅ Deploy anywhere with Docker

**Next recommended action:** Get QB credentials and start using with Claude!
