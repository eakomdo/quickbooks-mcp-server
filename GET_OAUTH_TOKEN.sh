#!/bin/bash

# Get OAuth tokens from Intuit
# Usage: ./GET_OAUTH_TOKEN.sh YOUR_AUTH_CODE

if [ -z "$1" ]; then
  echo "❌ You need to provide the Authorization Code"
  echo ""
  echo "Steps to get the code:"
  echo "1. Go to: https://developer.intuit.com/app/developer/qbo/playground"
  echo "2. Click 'Start OAuth Dance'"
  echo "3. Sign in with your Intuit sandbox account"
  echo "4. Click 'Allow'"
  echo "5. Copy the Authorization Code from the next page"
  echo ""
  echo "Usage: $0 YOUR_AUTH_CODE"
  echo "Example: $0 ABCDEFGHijklmnop"
  exit 1
fi

AUTH_CODE="$1"
CLIENT_ID="ABHhKijre3B7l4kCnyg5LBlPGAwzwQvPpBU8TYSom0L5Hg7EEI"
CLIENT_SECRET="zquPRbd362QOFu6QynbWIiyaJjlkOATTchj1TmjK"
REDIRECT_URI="https://developer.intuit.com/v2/OAuth2Playground/RedirectUrl"

echo "🔄 Getting OAuth tokens..."
echo ""

RESPONSE=$(curl -s -X POST https://quickbooks.api.intuit.com/oauth2/tokens/bearer \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code" \
  -d "code=$AUTH_CODE" \
  -d "redirect_uri=$REDIRECT_URI" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET")

echo "$RESPONSE" | grep -q "refresh_token"

if [ $? -eq 0 ]; then
  echo "✅ SUCCESS!"
  echo ""
  echo "Your tokens:"
  REFRESH_TOKEN=$(echo "$RESPONSE" | grep -o '"refresh_token":"[^"]*"' | cut -d'"' -f4)
  REALM_ID=$(echo "$RESPONSE" | grep -o '"realmId":"[^"]*"' | cut -d'"' -f4)
  echo "  Refresh Token: $REFRESH_TOKEN"
  echo "  Realm ID: $REALM_ID"
  echo ""
  echo "Add these to your .env file"
else
  echo "❌ ERROR - Token request failed"
  echo ""
  echo "Response: $RESPONSE"
fi
