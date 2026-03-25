#!/bin/bash

# SPEED SCRIPT: Fill .env with your credentials
# Usage: ./FILL_CREDENTIALS.sh YOUR_CLIENT_ID YOUR_CLIENT_SECRET YOUR_REALM_ID YOUR_REFRESH_TOKEN

if [ $# -lt 4 ]; then
  echo "❌ Usage: ./FILL_CREDENTIALS.sh <CLIENT_ID> <CLIENT_SECRET> <REALM_ID> <REFRESH_TOKEN>"
  echo ""
  echo "Example:"
  echo "  ./FILL_CREDENTIALS.sh ABCDEfghijKL1234567 aBcDeFgHiJkLmNoPqRsT 1234567890 refresh123456789"
  exit 1
fi

CLIENT_ID="$1"
CLIENT_SECRET="$2"
REALM_ID="$3"
REFRESH_TOKEN="$4"

cat > .env << ENVEOF
QUICKBOOKS_CLIENT_ID=$CLIENT_ID
QUICKBOOKS_CLIENT_SECRET=$CLIENT_SECRET
QUICKBOOKS_REALM_ID=$REALM_ID
QUICKBOOKS_REFRESH_TOKEN=$REFRESH_TOKEN
QUICKBOOKS_ENVIRONMENT=sandbox
ENVEOF

echo "✅ .env file created!"
echo ""
echo "Credentials set:"
echo "  Client ID: ${CLIENT_ID:0:10}..."
echo "  Client Secret: ${CLIENT_SECRET:0:10}..."
echo "  Realm ID: $REALM_ID"
echo "  Refresh Token: ${REFRESH_TOKEN:0:15}..."
echo "  Environment: sandbox"
echo ""
echo "🚀 Ready to run: ./qbo-mcp"
