package qbo

import (
	"context"
	"encoding/json"
	"fmt"
)

// UnwrapEntity extracts the entity object from a typical QBO JSON envelope.
func UnwrapEntity(raw json.RawMessage, entity string) (map[string]any, error) {
	key := CapitalizeEntity(entity)
	var w map[string]json.RawMessage
	if err := json.Unmarshal(raw, &w); err != nil {
		return nil, err
	}
	body, ok := w[key]
	if !ok || string(body) == "null" {
		return nil, fmt.Errorf("missing %s in response: %s", key, string(raw))
	}
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		return nil, err
	}
	return m, nil
}

// ReadEntity loads an entity by Id and returns the inner object.
func (c *Client) ReadEntity(ctx context.Context, entity, id string) (map[string]any, error) {
	raw, err := c.Read(ctx, entity, id)
	if err != nil {
		return nil, err
	}
	return UnwrapEntity(raw, entity)
}

// CreateEntity POST-creates an entity and returns the persisted object.
func (c *Client) CreateEntity(ctx context.Context, entity string, body any) (map[string]any, error) {
	raw, err := c.Create(ctx, entity, body)
	if err != nil {
		return nil, err
	}
	return UnwrapEntity(raw, entity)
}

// UpdateEntity POST-updates (sparse) and returns the updated object.
func (c *Client) UpdateEntity(ctx context.Context, entity string, body any) (map[string]any, error) {
	raw, err := c.Update(ctx, entity, body)
	if err != nil {
		return nil, err
	}
	return UnwrapEntity(raw, entity)
}
