#!/bin/bash

# After doing OAuth Dance and getting a CODE, use this script
# Usage: ./GET_REFRESH_TOKEN.sh YOUR_CODE YOUR_CLIENT_ID YOUR_CLIENT_SECRET

if [ $# -lt 3 ]; then
  echo "❌ Usage: ./GET_REFRESH_TOKEN.sh <CODE> <CLIENT_ID> <CLIENT_SECRET>"
  echo ""
  echo "Example:"
  echo "  ./GET_REFRESH_TOKEN.sh AUTH_CODE_FROM_OAUTH_DANCE ABCDEfghij aBcDeFgHiJkL"
  exit 1
fi

CODE="$1"
CLIENT_ID="$2"
CLIENT_SECRET="$3"

echo "🔄 Getting refresh token from Intuit..."
echo ""

RESPONSE=$(curl -s -X POST https://quickbooks.api.intuit.com/oauth2/tokens/bearer \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=authorization_code" \
  -d "code=$CODE" \
  -d "redirect_uri=https://developer.intuit.com/v2/oauth/authorize" \
  -d "client_id=$CLIENT_ID" \
  -d "client_secret=$CLIENT_SECRET")

echo "$RESPONSE" | grep -q "refresh_token"

if [ $? -eq 0 ]; then
  echo "✅ SUCCESS! Here's your refresh token:"
  echo ""
  echo "$RESPONSE" | grep -o '"refresh_token":"[^"]*"' | cut -d'"' -f4
  echo ""
  echo "Keep this safe! You need it for .env"
else
  echo "❌ ERROR - Token request failed"
  echo "Response: $RESPONSE"
fi
