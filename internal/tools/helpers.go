package tools

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func okResult(parts ...string) *mcp.CallToolResult {
	var c []mcp.Content
	for _, p := range parts {
		c = append(c, &mcp.TextContent{Text: p})
	}
	return &mcp.CallToolResult{Content: c}
}

func errResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}
}

func jsonPretty(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprint(v)
	}
	return string(b)
}

func toMap(v any) (map[string]any, bool) {
	m, ok := v.(map[string]any)
	return m, ok
}

func num(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	case int64:
		return float64(t)
	case json.Number:
		f, _ := t.Float64()
		return f
	default:
		f, err := strconv.ParseFloat(fmt.Sprint(v), 64)
		if err != nil {
			return 0
		}
		return f
	}
}
