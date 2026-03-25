package qbo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

// BuildSearchCriteria mirrors src/helpers/build-quickbooks-search-criteria.ts
func BuildSearchCriteria(input any) any {
	if input == nil {
		return map[string]any{}
	}
	if arr, ok := input.([]any); ok {
		return arr
	}
	m, ok := input.(map[string]any)
	if !ok {
		return input
	}
	advancedKeys := map[string]struct{}{
		"filters": {}, "asc": {}, "desc": {}, "limit": {}, "offset": {},
		"count": {}, "fetchAll": {},
	}
	isAdvanced := false
	for k := range m {
		if _, ok := advancedKeys[k]; ok {
			isAdvanced = true
			break
		}
	}
	if !isAdvanced {
		return m
	}
	var criteria []map[string]any
	if filters, ok := m["filters"].([]any); ok {
		for _, f := range filters {
			fm, ok := f.(map[string]any)
			if !ok {
				continue
			}
			entry := map[string]any{"field": fm["field"], "value": fm["value"]}
			if op, ok := fm["operator"]; ok {
				entry["operator"] = op
			}
			criteria = append(criteria, entry)
		}
	}
	if asc, ok := m["asc"]; ok && asc != nil {
		criteria = append(criteria, map[string]any{"field": "asc", "value": asc})
	}
	if desc, ok := m["desc"]; ok && desc != nil {
		criteria = append(criteria, map[string]any{"field": "desc", "value": desc})
	}
	if lim, ok := m["limit"]; ok {
		criteria = append(criteria, map[string]any{"field": "limit", "value": lim})
	}
	if off, ok := m["offset"]; ok {
		criteria = append(criteria, map[string]any{"field": "offset", "value": off})
	}
	if c, ok := m["count"]; ok && truthy(c) {
		criteria = append(criteria, map[string]any{"field": "count", "value": true})
	}
	if f, ok := m["fetchAll"]; ok && truthy(f) {
		criteria = append(criteria, map[string]any{"field": "fetchAll", "value": true})
	}
	if len(criteria) == 0 {
		return map[string]any{}
	}
	return criteria
}

func truthy(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return strings.EqualFold(t, "true")
	default:
		return false
	}
}

func deepCopyCriteria(criteria any) any {
	b, err := json.Marshal(criteria)
	if err != nil {
		return criteria
	}
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return criteria
	}
	return v
}

func extractFetchAll(criteria any) bool {
	switch c := criteria.(type) {
	case map[string]any:
		return truthy(c["fetchAll"])
	case []any:
		for _, e := range c {
			em, ok := e.(map[string]any)
			if !ok {
				continue
			}
			if strings.EqualFold(fmt.Sprint(em["field"]), "fetchall") && truthy(em["value"]) {
				return true
			}
		}
	}
	return false
}

// --- Query SQL (node-quickbooks module.query / criteriaToString) ---

func buildQuerySQL(entity string, criteria any) (string, error) {
	criteria = deepCopyCriteria(criteria)
	if criteria == nil {
		criteria = map[string]any{}
	}

	base := "select * from " + entity
	if hasCountFlag(criteria) {
		base = strings.Replace(base, "select * from", "select count(*) from", 1)
		criteria = stripCountFlag(criteria)
	}

	limit, offset, fetchAll := 1000, 1, extractFetchAll(criteria)
	_ = fetchAll // pagination handled in Query()

	body, lim, off, err := normalizeCriteriaLimits(criteria, limit, offset)
	if err != nil {
		return "", err
	}
	if lim != nil {
		limit = *lim
	}
	if off != nil {
		offset = *off
	}

	where, order, err := criteriaToWhereOrder(body)
	if err != nil {
		return "", err
	}

	sql := base
	if where != "" {
		sql += " where " + where
	}
	sql += order
	sql += fmt.Sprintf(" startposition %d", offset)
	sql += fmt.Sprintf(" maxresults %d", limit)
	return sql, nil
}

func hasCountFlag(criteria any) bool {
	switch c := criteria.(type) {
	case map[string]any:
		return truthy(c["count"])
	case []any:
		for _, e := range c {
			em, ok := e.(map[string]any)
			if !ok {
				continue
			}
			if strings.EqualFold(fmt.Sprint(em["field"]), "count") && truthy(em["value"]) {
				return true
			}
		}
	}
	return false
}

func stripCountFlag(criteria any) any {
	switch c := criteria.(type) {
	case map[string]any:
		out := map[string]any{}
		for k, v := range c {
			if strings.EqualFold(k, "count") {
				continue
			}
			out[k] = v
		}
		return out
	case []any:
		var out []any
		for _, e := range c {
			em, ok := e.(map[string]any)
			if ok && strings.EqualFold(fmt.Sprint(em["field"]), "count") {
				continue
			}
			out = append(out, e)
		}
		return out
	default:
		return criteria
	}
}

func normalizeCriteriaLimits(criteria any, defaultLimit, defaultOffset int) (body any, lim, off *int, err error) {
	limit := defaultLimit
	offset := defaultOffset

	switch c := criteria.(type) {
	case map[string]any:
		m := deepCopyCriteria(c).(map[string]any)
		if v, ok := m["limit"]; ok {
			limit = intFromAny(v)
			delete(m, "limit")
		}
		if v, ok := m["offset"]; ok {
			offset = intFromAny(v)
			delete(m, "offset")
		}
		delete(m, "fetchAll")
		delete(m, "count")
		// asc/desc stay in map for where builder — actually they should NOT be in WHERE.
		// Remove and handle via array path only; for map, strip asc/desc to order.
		return m, &limit, &offset, nil

	case []any:
		arr := deepCopyCriteria(c).([]any)
		var filtered []any
		for _, e := range arr {
			em, ok := e.(map[string]any)
			if !ok {
				continue
			}
			f := fmt.Sprint(em["field"])
			switch {
			case strings.EqualFold(f, "limit"):
				limit = intFromAny(em["value"])
			case strings.EqualFold(f, "offset"):
				offset = intFromAny(em["value"])
			default:
				filtered = append(filtered, e)
			}
		}
		// node adds offset 1 if missing
		hasOff := false
		for _, e := range arr {
			em, ok := e.(map[string]any)
			if ok && strings.EqualFold(fmt.Sprint(em["field"]), "offset") {
				hasOff = true
			}
		}
		if !hasOff {
			filtered = append(filtered, map[string]any{"field": "offset", "value": 1})
			offset = 1
		}
		return filtered, &limit, &offset, nil

	default:
		return criteria, &limit, &offset, nil
	}
}

func intFromAny(v any) int {
	switch t := v.(type) {
	case float64:
		return int(t)
	case int:
		return t
	case json.Number:
		i, _ := t.Int64()
		return int(i)
	default:
		return 0
	}
}

func criteriaToWhereOrder(criteria any) (where string, order string, err error) {
	var ascVal, descVal *string

	switch c := criteria.(type) {
	case map[string]any:
		// Pull sort keys
		if v, ok := c["asc"]; ok {
			s := fmt.Sprint(v)
			ascVal = &s
			delete(c, "asc")
		}
		if v, ok := c["desc"]; ok {
			s := fmt.Sprint(v)
			descVal = &s
			delete(c, "desc")
		}
		flat := toCriterionMap(c)
		where, err = buildWhereClause(flat)
		if err != nil {
			return "", "", err
		}
	case []any:
		var flat []map[string]any
		for _, e := range c {
			em, ok := e.(map[string]any)
			if !ok {
				continue
			}
			f := fmt.Sprint(em["field"])
			if strings.EqualFold(f, "asc") {
				s := fmt.Sprint(em["value"])
				ascVal = &s
				continue
			}
			if strings.EqualFold(f, "desc") {
				s := fmt.Sprint(em["value"])
				descVal = &s
				continue
			}
			if strings.EqualFold(f, "fetchall") {
				continue
			}
			flat = append(flat, em)
		}
		flat2 := flattenCriterionEntries(flat)
		where, err = buildWhereClause(flat2)
		if err != nil {
			return "", "", err
		}
	default:
		return "", "", nil
	}

	if ascVal != nil {
		order += " orderby " + *ascVal + " asc"
	}
	if descVal != nil {
		order += " orderby " + *descVal + " desc"
	}
	return where, order, nil
}

func toCriterionMap(c map[string]any) []map[string]any {
	var out []map[string]any
	for k, v := range c {
		op := "="
		if v != nil && reflect.TypeOf(v).Kind() == reflect.Slice {
			op = "IN"
		}
		out = append(out, map[string]any{"field": k, "value": v, "operator": op})
	}
	return out
}

func flattenCriterionEntries(entries []map[string]any) []map[string]any {
	var out []map[string]any
	for _, em := range entries {
		if _, ok := em["field"]; ok {
			if _, ok2 := em["value"]; ok2 {
				op := "="
				if o, ok := em["operator"]; ok && o != nil {
					op = fmt.Sprint(o)
				}
				out = append(out, map[string]any{
					"field": em["field"], "value": em["value"], "operator": op,
				})
				continue
			}
		}
		// nested object form
		for k, v := range em {
			op := "="
			if v != nil && reflect.TypeOf(v).Kind() == reflect.Slice {
				op = "IN"
			}
			out = append(out, map[string]any{"field": k, "value": v, "operator": op})
		}
	}
	return out
}

func quoteValue(x any) string {
	switch v := x.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "\\'") + "'"
	default:
		return fmt.Sprint(v)
	}
}

func buildWhereClause(flat []map[string]any) (string, error) {
	var parts []string
	for _, criterion := range flat {
		f := fmt.Sprint(criterion["field"])
		op := fmt.Sprint(criterion["operator"])
		if op == "" {
			op = "="
		}
		val := criterion["value"]
		var seg string
		if arr, ok := val.([]any); ok && op == "IN" {
			var qs []string
			for _, x := range arr {
				qs = append(qs, quoteValue(x))
			}
			seg = f + " IN (" + strings.Join(qs, ",") + ")"
		} else {
			seg = f + " " + op + " " + quoteValue(val)
		}
		parts = append(parts, seg)
	}
	return strings.Join(parts, " and "), nil
}

// Query runs a single-page query. If criteria requests fetchAll, use queryFetchAll instead.
func (c *Client) Query(ctx context.Context, entity string, criteria any) (json.RawMessage, error) {
	if extractFetchAll(criteria) {
		return c.queryFetchAll(ctx, entity, criteria)
	}
	sql, err := buildQuerySQL(entity, criteria)
	if err != nil {
		return nil, err
	}
	q := urlValuesQuery(sql)
	return c.do(ctx, http.MethodGet, "/query", q, nil)
}

func urlValuesQuery(sql string) url.Values {
	v := url.Values{}
	v.Set("query", sql)
	return v
}

func (c *Client) queryFetchAll(ctx context.Context, entity string, criteria any) (json.RawMessage, error) {
	criteria = deepCopyCriteria(criteria)
	limit := 1000
	if arr, ok := criteria.([]any); ok {
		if lmt := findFieldValue(arr, "limit"); lmt != nil {
			limit = intFromAny(lmt)
		}
	}
	if m, ok := criteria.(map[string]any); ok {
		if v, ok := m["limit"]; ok {
			limit = intFromAny(v)
		}
	}

	offset := 1
	var merged []any
	var lastRaw json.RawMessage

	for {
		pageCrit := applyPageCriteria(deepCopyCriteria(criteria), limit, offset)

		sql, err := buildQuerySQL(entity, pageCrit)
		if err != nil {
			return nil, err
		}
		raw, err := c.do(ctx, http.MethodGet, "/query", urlValuesQuery(sql), nil)
		if err != nil {
			return nil, err
		}
		lastRaw = raw

		key := CapitalizeEntity(entity)
		items, maxResults, err := ExtractQueryEntities(raw, key)
		if err != nil {
			return nil, err
		}
		merged = append(merged, items...)
		if maxResults < limit {
			break
		}
		offset += limit
	}

	var wrap map[string]any
	_ = json.Unmarshal(lastRaw, &wrap)
	qr, _ := wrap["QueryResponse"].(map[string]any)
	if qr == nil {
		qr = map[string]any{}
	}
	qr[CapitalizeEntity(entity)] = merged
	qr["maxResults"] = len(merged)
	wrap["QueryResponse"] = qr
	out, _ := json.Marshal(wrap)
	return out, nil
}

func findFieldValue(arr []any, name string) any {
	for _, e := range arr {
		em, ok := e.(map[string]any)
		if !ok {
			continue
		}
		if strings.EqualFold(fmt.Sprint(em["field"]), name) {
			return em["value"]
		}
	}
	return nil
}

// applyPageCriteria returns a copy with limit/offset set and fetchAll removed (for paged fetch).
func applyPageCriteria(criteria any, limit, offset int) any {
	switch c := criteria.(type) {
	case map[string]any:
		m := deepCopyCriteria(c).(map[string]any)
		m["limit"] = limit
		m["offset"] = offset
		delete(m, "fetchAll")
		return m
	case []any:
		arr := deepCopyCriteria(c).([]any)
		var out []any
		for _, e := range arr {
			em, ok := e.(map[string]any)
			if !ok {
				out = append(out, e)
				continue
			}
			f := fmt.Sprint(em["field"])
			if strings.EqualFold(f, "fetchall") {
				continue
			}
			if strings.EqualFold(f, "limit") || strings.EqualFold(f, "offset") {
				continue
			}
			out = append(out, e)
		}
		out = append(out,
			map[string]any{"field": "limit", "value": limit},
			map[string]any{"field": "offset", "value": offset},
		)
		return out
	default:
		return criteria
	}
}

// ExtractQueryEntities reads QueryResponse entity list (or single object as one-element slice).
func ExtractQueryEntities(raw []byte, entityKey string) (items []any, maxResults int, err error) {
	var wrap map[string]json.RawMessage
	if err := json.Unmarshal(raw, &wrap); err != nil {
		return nil, 0, err
	}
	qrRaw, ok := wrap["QueryResponse"]
	if !ok {
		return nil, 0, nil
	}
	var qr map[string]json.RawMessage
	if err := json.Unmarshal(qrRaw, &qr); err != nil {
		return nil, 0, err
	}
	_ = json.Unmarshal(qr["maxResults"], &maxResults)
	entRaw, ok := qr[entityKey]
	if !ok || string(entRaw) == "null" {
		return nil, maxResults, nil
	}
	if len(entRaw) > 0 && entRaw[0] == '[' {
		var arr []any
		_ = json.Unmarshal(entRaw, &arr)
		return arr, maxResults, nil
	}
	var one any
	_ = json.Unmarshal(entRaw, &one)
	return []any{one}, maxResults, nil
}
