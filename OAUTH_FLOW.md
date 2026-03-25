# Get Refresh Token (Takes 3 minutes)

## Step 1: Go to OAuth Playground
https://developer.intuit.com/app/developer/qbo/playground

## Step 2: Follow the OAuth "Dance"
1. Select "QuickBooks Online" from dropdown
2. Click the big blue "Start OAuth Dance" button
3. Sign in with your Intuit account (same as developer portal)
4. Click "Allow" for permissions
5. You'll be redirected to a success page with an **Authorization Code**

## Step 3: Copy the Authorization Code
Save this code from the redirect page - looks like: `ABCDEFGHijklmnop...`

## Step 4: Get Refresh Token

Replace these values with YOUR actual values and run in your terminal:

```bash
curl -X POST https://quickbooks.api.intuit.com/oauth2/tokens/bearer \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code" \
  -d "code=YOUR_AUTHORIZATION_CODE" \
  -d "redirect_uri=https://developer.intuit.com" \
  -d "client_id=ABHhKijre3B7l4kCnyg5LBlPGAwzwQvPpBU8TYSom0L5Hg7EEI" \
  -d "client_secret=zquPRbd362QOFu6QynbWIiyaJjlkOATTchj1TmjK"
```

## Step 5: Copy from Response
The response will look like:
```json
{
  "access_token": "...",
  "refresh_token": "YOUR_REFRESH_TOKEN_HERE",
  "token_type": "Bearer",
  "x_refresh_token_expires_in": 8726400,
  "expires_in": 3600
}
```

**COPY the `refresh_token` value**

## Step 6: Find Realm ID
Also from the OAuth Dance response, look for:
- "realmId" field (your Realm ID)

If not in response, check Intuit portal under your app settings.

