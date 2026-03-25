#!/bin/bash

# QuickBooks Online MCP Server Test Script
# Tests all available tools via curl
# Prerequisites: Server running on http://localhost:3000

BASE_URL="http://localhost:3000"
SERVER_ENDPOINT="$BASE_URL/mcp"

# Color codes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}QuickBooks MCP Server Test Suite${NC}"
echo -e "${BLUE}================================${NC}\n"

# Test 1: Initialize connection
test_initialize() {
    echo -e "${YELLOW}[TEST 1] Initialize MCP Connection${NC}"
    echo "Testing basic MCP protocol initialization..."
    
    curl -s "$SERVER_ENDPOINT" -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {
                "protocolVersion": "2024-11-05",
                "capabilities": {},
                "clientInfo": {
                    "name": "test-client",
                    "version": "1.0"
                }
            }
        }' | grep -oP 'data: \K.*' | jq .
    
    echo -e "${GREEN}✓ Initialize test complete\n${NC}"
}

# Test 2: List all available tools
test_list_tools() {
    echo -e "${YELLOW}[TEST 2] List All Available Tools${NC}"
    echo "Fetching all registered QuickBooks tools..."
    
    curl -s "$SERVER_ENDPOINT" -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "jsonrpc": "2.0",
            "id": 2,
            "method": "tools/list",
            "params": {}
        }' | grep -oP 'data: \K.*' | jq .
    
    echo -e "${GREEN}✓ Tools list test complete\n${NC}"
}

# Test 3: Search customers
test_search_customers() {
    echo -e "${YELLOW}[TEST 3] Search Customers${NC}"
    echo "Searching for customers (first 10)..."
    
    curl -s "$SERVER_ENDPOINT" -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "jsonrpc": "2.0",
            "id": 3,
            "method": "tools/call",
            "params": {
                "name": "quickbook_search_customers",
                "arguments": {
                    "params": {
                        "query": "SELECT * FROM Customer MAXRESULTS 10"
                    }
                }
            }
        }' | grep -oP 'data: \K.*' | jq .
    
    echo -e "${GREEN}✓ Customer search test complete\n${NC}"
}

# Test 4: Search invoices
test_search_invoices() {
    echo -e "${YELLOW}[TEST 4] Search Invoices${NC}"
    echo "Searching for invoices (first 10)..."
    
    curl -s "$SERVER_ENDPOINT" -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "jsonrpc": "2.0",
            "id": 4,
            "method": "tools/call",
            "params": {
                "name": "quickbook_search_invoices",
                "arguments": {
                    "params": {
                        "query": "SELECT * FROM Invoice MAXRESULTS 10"
                    }
                }
            }
        }' | grep -oP 'data: \K.*' | jq .
    
    echo -e "${GREEN}✓ Invoice search test complete\n${NC}"
}

# Test 5: Search bills
test_search_bills() {
    echo -e "${YELLOW}[TEST 5] Search Bills${NC}"
    echo "Searching for bills (first 10)..."
    
    curl -s "$SERVER_ENDPOINT" -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "jsonrpc": "2.0",
            "id": 5,
            "method": "tools/call",
            "params": {
                "name": "quickbook_search_bills",
                "arguments": {
                    "params": {
                        "query": "SELECT * FROM Bill MAXRESULTS 10"
                    }
                }
            }
        }' | grep -oP 'data: \K.*' | jq .
    
    echo -e "${GREEN}✓ Bill search test complete\n${NC}"
}

# Test 6: Search estimates
test_search_estimates() {
    echo -e "${YELLOW}[TEST 6] Search Estimates${NC}"
    echo "Searching for estimates (first 10)..."
    
    curl -s "$SERVER_ENDPOINT" -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "jsonrpc": "2.0",
            "id": 6,
            "method": "tools/call",
            "params": {
                "name": "quickbook_search_estimates",
                "arguments": {
                    "params": {
                        "query": "SELECT * FROM Estimate MAXRESULTS 10"
                    }
                }
            }
        }' | grep -oP 'data: \K.*' | jq .
    
    echo -e "${GREEN}✓ Estimate search test complete\n${NC}"
}

# Test 7: Create test customer (sandbox only!)
test_create_customer() {
    echo -e "${YELLOW}[TEST 7] Create Test Customer (Sandbox Only)${NC}"
    echo "Creating a test customer named 'Test Corp'..."
    
    curl -s "$SERVER_ENDPOINT" -X POST \
        -H "Content-Type: application/json" \
        -d '{
            "jsonrpc": "2.0",
            "id": 7,
            "method": "tools/call",
            "params": {
                "name": "quickbook_create_customer",
                "arguments": {
                    "params": {
                        "customer": {
                            "DisplayName": "Test Corp - '$(date +%s)'",
                            "FullyQualifiedName": "Test Corp - '$(date +%s)'",
                            "CustomerTypeRef": {
                                "value": "1"
                            }
                        }
                    }
                }
            }
        }' | grep -oP 'data: \K.*' | jq .
    
    echo -e "${GREEN}✓ Create customer test complete\n${NC}"
}

# Display usage info
usage() {
    echo "Usage: ./test-mcp-server.sh [test-number|all]"
    echo ""
    echo "Available tests:"
    echo "  1 - Initialize MCP Connection"
    echo "  2 - List All Available Tools"
    echo "  3 - Search Customers"
    echo "  4 - Search Invoices"
    echo "  5 - Search Bills"
    echo "  6 - Search Estimates"
    echo "  7 - Create Test Customer (Sandbox Only)"
    echo "  all - Run all tests"
    echo ""
    echo "Example:"
    echo "  ./test-mcp-server.sh 1        # Run test 1"
    echo "  ./test-mcp-server.sh all      # Run all tests"
}

# Check if jq is installed
check_jq() {
    if ! command -v jq &> /dev/null; then
        echo -e "${YELLOW}Warning: jq not found. Install it for pretty JSON output:${NC}"
        echo "  sudo apt install jq"
        echo ""
    fi
}

# Check if server is running
check_server() {
    echo -e "${BLUE}Checking if MCP server is running...${NC}"
    if ! curl -s "$SERVER_ENDPOINT" -X POST \
        -H "Content-Type: application/json" \
        -d '{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"check","version":"1.0"}}}' | grep -q "serverInfo"; then
        echo -e "${YELLOW}Error: Cannot connect to MCP server at $SERVER_ENDPOINT${NC}"
        echo "Make sure the server is running:"
        echo "  cd ~/Downloads/quickbooks-online-mcp"
        echo "  PORT=3000 ./qbo-mcp"
        exit 1
    fi
    echo -e "${GREEN}✓ Server is running\n${NC}"
}

# Main execution
check_jq
check_server

if [ -z "$1" ]; then
    usage
    exit 0
fi

case "$1" in
    1)
        test_initialize
        ;;
    2)
        test_list_tools
        ;;
    3)
        test_search_customers
        ;;
    4)
        test_search_invoices
        ;;
    5)
        test_search_bills
        ;;
    6)
        test_search_estimates
        ;;
    7)
        test_create_customer
        ;;
    all)
        test_initialize
        test_list_tools
        test_search_customers
        test_search_invoices
        test_search_bills
        test_search_estimates
        test_create_customer
        ;;
    *)
        echo "Unknown test: $1"
        usage
        exit 1
        ;;
esac

echo -e "${BLUE}================================${NC}"
echo -e "${BLUE}Test suite execution complete${NC}"
echo -e "${BLUE}================================${NC}"
