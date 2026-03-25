#!/usr/bin/env python3
"""
QuickBooks MCP Server Test Client
Proper MCP protocol client that maintains session state
"""

import subprocess
import json
import sys
import os

class MCPClient:
    def __init__(self, server_cmd):
        self.proc = subprocess.Popen(
            server_cmd,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=1
        )
        self.request_id = 0
        self._initialize()
    
    def _initialize(self):
        """Initialize the MCP session"""
        self.request_id += 1
        req = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": "initialize",
            "params": {
                "protocolVersion": "2024-11-05",
                "capabilities": {},
                "clientInfo": {
                    "name": "test-client",
                    "version": "1.0"
                }
            }
        }
        response = self._send(req)
        print("✓ MCP Session initialized\n")
        return response
    
    def _send(self, request):
        """Send request and get response"""
        json_str = json.dumps(request)
        self.proc.stdin.write(json_str + '\n')
        self.proc.stdin.flush()
        
        response_str = self.proc.stdout.readline()
        try:
            return json.loads(response_str)
        except json.JSONDecodeError:
            print(f"Error parsing response: {response_str}")
            return None
    
    def list_tools(self):
        """List all available tools"""
        self.request_id += 1
        req = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": "tools/list",
            "params": {}
        }
        return self._send(req)
    
    def call_tool(self, tool_name, arguments):
        """Call a specific tool"""
        self.request_id += 1
        req = {
            "jsonrpc": "2.0",
            "id": self.request_id,
            "method": "tools/call",
            "params": {
                "name": tool_name,
                "arguments": arguments
            }
        }
        return self._send(req)
    
    def close(self):
        """Close the connection"""
        self.proc.terminate()
        self.proc.wait()


def pretty_print(obj, indent=2):
    """Pretty print JSON"""
    print(json.dumps(obj, indent=indent))


def run_tests():
    """Run all tests"""
    
    # Start the server
    server_cmd = ['/home/emma/Downloads/quickbooks-online-mcp/qbo-mcp']
    
    print("\n" + "="*50)
    print("QuickBooks MCP Server Test Suite")
    print("="*50 + "\n")
    
    try:
        client = MCPClient(server_cmd)
        
        # Test 1: List all tools
        print("📋 [TEST 1] List All Available Tools")
        print("-" * 50)
        response = client.list_tools()
        if response and 'result' in response:
            tools = response['result'].get('tools', [])
            print(f"✓ Found {len(tools)} tools:\n")
            for tool in tools:
                print(f"  • {tool['name']}")
                print(f"    {tool['description']}\n")
        else:
            print("✗ Error listing tools")
            pretty_print(response)
        
        # Test 2: Search customers
        print("\n👥 [TEST 2] Search Customers")
        print("-" * 50)
        response = client.call_tool('quickbook_search_customers', {
            'params': {
                'query': 'SELECT * FROM Customer MAXRESULTS 5'
            }
        })
        if response and 'result' in response:
            print("✓ Customer search successful")
            pretty_print(response['result'])
        else:
            print("✗ Error searching customers")
            pretty_print(response)
        
        # Test 3: Search invoices
        print("\n📄 [TEST 3] Search Invoices")
        print("-" * 50)
        response = client.call_tool('quickbook_search_invoices', {
            'params': {
                'query': 'SELECT * FROM Invoice MAXRESULTS 5'
            }
        })
        if response and 'result' in response:
            print("✓ Invoice search successful")
            pretty_print(response['result'])
        else:
            print("✗ Error searching invoices")
            pretty_print(response)
        
        # Test 4: Search bills
        print("\n💰 [TEST 4] Search Bills")
        print("-" * 50)
        response = client.call_tool('quickbook_search_bills', {
            'params': {
                'query': 'SELECT * FROM Bill MAXRESULTS 5'
            }
        })
        if response and 'result' in response:
            print("✓ Bill search successful")
            pretty_print(response['result'])
        else:
            print("✗ Error searching bills")
            pretty_print(response)
        
        # Test 5: Search estimates
        print("\n📊 [TEST 5] Search Estimates")
        print("-" * 50)
        response = client.call_tool('quickbook_search_estimates', {
            'params': {
                'query': 'SELECT * FROM Estimate MAXRESULTS 5'
            }
        })
        if response and 'result' in response:
            print("✓ Estimate search successful")
            pretty_print(response['result'])
        else:
            print("✗ Error searching estimates")
            pretty_print(response)
        
        client.close()
        
        print("\n" + "="*50)
        print("✅ Test Suite Complete!")
        print("="*50 + "\n")
        
    except Exception as e:
        print(f"❌ Error: {e}")
        sys.exit(1)


if __name__ == '__main__':
    run_tests()
