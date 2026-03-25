package qbo

import (
	"context"
	"fmt"
	"strings"
)

// idFromFlexible extracts an Id from a string id or a QuickBooks entity map.
func idFromFlexible(v any) string {
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	case map[string]any:
		if id, ok := t["Id"].(string); ok {
			return id
		}
		if id, ok := t["id"].(string); ok {
			return id
		}
	}
	return strings.TrimSpace(fmt.Sprint(v))
}

// DeleteCustomerOrDeactivate tries QBO delete; on failure marks Active=false (matches TS fallback).
func (c *Client) DeleteCustomerOrDeactivate(ctx context.Context, idOrEntity any) (map[string]any, error) {
	raw, delErr := c.Delete(ctx, "customer", idOrEntity)
	if delErr == nil {
		if m, e := UnwrapEntity(raw, "customer"); e == nil {
			return m, nil
		}
		return map[string]any{"status": "deleted", "detail": string(raw)}, nil
	}
	id := idFromFlexible(idOrEntity)
	cust, readErr := c.ReadEntity(ctx, "customer", id)
	if readErr != nil {
		return nil, fmt.Errorf("delete failed: %v; could not load customer %q: %w", delErr, id, readErr)
	}
	cust["Active"] = false
	cust["sparse"] = true
	return c.UpdateEntity(ctx, "customer", cust)
}
