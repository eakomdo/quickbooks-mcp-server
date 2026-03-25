# ⚡ SPEED GUIDE: Get 5 Credentials in 10 Minutes

You have: ✅ Client ID, ✅ Client Secret  
You need: ⏳ Realm ID, ⏳ Refresh Token, ⏳ Environment

---

## Step 1: Get Realm ID (2 minutes)

**In Intuit Developer Portal:**
1. Go to https://developer.intuit.com
2. Click "My Apps" (top right)
3. Click your QB app
4. Go to "Keys & OAuth" tab
5. **Scroll down** → Find "Realm IDs" section
6. Copy your Realm ID (looks like: `1234567890`)

✅ **Realm ID = copied**

---

## Step 2: Get Refresh Token (5 minutes) - FASTEST METHOD

**Option A: Use OAuth Playground (Easiest)**

1. Go to: https://developer.intuit.com/app/developer/qbo/playground
2. Select "QuickBooks Online" 
3. Click "Start OAuth Dance"
4. Sign in with your Intuit account
5. After redirect, **Check browser DevTools (F12) → Network tab**
6. Look for request with `access_token` and `refresh_token` in response
7. Copy `refresh_token` value

✅ **Refresh Token = copied**

---

**Option B: cURL (If OAuth Playground doesn't work)**

Replace with your actual values:

```bash
curl -X POST https://quickbooks.api.intuit.com/oauth2/tokens/bearer \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code&code=YOUR_AUTH_CODE&redirect_uri=https://developer.intuit.com&client_id=YOUR_CLIENT_ID&client_secret=YOUR_CLIENT_SECRET"
```

The response will have: `"refresh_token": "your_refresh_token_value"`

---

## Step 3: Set Environment (30 seconds)

```
QUICKBOOKS_ENVIRONMENT=sandbox
```

Or if using production:
```
QUICKBOOKS_ENVIRONMENT=production
```

**Default: Use `sandbox` for testing**

---

## Step 4: Create .env File (1 minute)

```bash
cd /home/emma/Downloads/quickbooks-online-mcp

cat > .env << 'EOF'
QUICKBOOKS_CLIENT_ID=your_client_id_here
QUICKBOOKS_CLIENT_SECRET=your_client_secret_here
QUICKBOOKS_REALM_ID=your_realm_id_here
QUICKBOOKS_REFRESH_TOKEN=your_refresh_token_here
QUICKBOOKS_ENVIRONMENT=sandbox
EOF
```

---

## Step 5: Verify & Run (2 minutes)

```bash
# Verify credentials loaded
cat .env

# Run server
./qbo-mcp
```

You should see:
```
INFO: MCP Server initialized
INFO: Tools registered: 40+
```

---

## Troubleshooting Fast

| Problem | Fix |
|---------|-----|
| "Invalid client" | Check Client ID/Secret spelling |
| "Invalid refresh token" | Re-run OAuth flow at playground |
| "Realm ID not found" | Check Intuit portal → Keys & OAuth tab |
| Server won't start | Check .env file has no extra spaces |

---

## Direct Links

- **Intuit Developer Portal:** https://developer.intuit.com
- **My Apps:** https://developer.intuit.com/app/developer/myapps
- **OAuth Playground:** https://developer.intuit.com/app/developer/qbo/playground
- **Keys & OAuth:** https://developer.intuit.com/app/developer/qbo/playground (scroll to Keys section)

---

## Total Time: ~10 minutes ✨

1. Get Realm ID (2 min)
2. Get Refresh Token (5 min)
3. Create .env (1 min)
4. Run & test (2 min)

**Go! Timer starts now.** 🚀
