package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/qboapi/qbo-mcp-server/internal/config"
	"github.com/qboapi/qbo-mcp-server/internal/qbo"
)

// RegisterAll registers the same QuickBooks tools as the TypeScript server (see src/index.ts).
func RegisterAll(s *mcp.Server, cfg *config.Config, userBearer string, personaID string) {
	schema := json.RawMessage(`{"type":"object","properties":{"params":{"type":"object","description":"Tool parameters (mirrors TS zod schemas in src/tools)"}},"required":["params"]}`)
	ub := userBearer

	add := func(name, desc string, run func(context.Context, map[string]any) (*mcp.CallToolResult, error)) {
		s.AddTool(&mcp.Tool{Name: name, Description: desc, InputSchema: schema}, func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			ctx = qbo.WithUserBearer(ctx, ub)
			if personaID != "" {
				ctx = qbo.WithPersonaID(ctx, personaID)
			}
			var top map[string]any
			if err := json.Unmarshal(req.Params.Arguments, &top); err != nil {
				return errResult(fmt.Sprintf("invalid arguments: %v", err)), nil
			}
			params, _ := top["params"].(map[string]any)
			if params == nil {
				params = map[string]any{}
			}
			return run(ctx, params)
		})
	}

	// --- Customers ---
	add("quickbook_create_customer", "Create a customer in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		cust, ok := toMap(p["customer"])
		if !ok {
			return errResult("params.customer object required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.CreateEntity(ctx, "customer", cust)
		if err != nil {
			return errResult(fmt.Sprintf("Error creating customer: %v", err)), nil
		}
		return okResult("Customer created:", jsonPretty(m)), nil
	})
	add("quickbook_get_customer", "Get a customer by Id from QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		id, _ := p["id"].(string)
		if id == "" {
			return errResult("id required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.ReadEntity(ctx, "customer", id)
		if err != nil {
			return errResult(fmt.Sprintf("Error getting customer: %v", err)), nil
		}
		return okResult("Customer:", jsonPretty(m)), nil
	})
	add("quickbook_update_customer", "Update a customer in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		cust, ok := toMap(p["customer"])
		if !ok {
			return errResult("params.customer object required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.UpdateEntity(ctx, "customer", cust)
		if err != nil {
			return errResult(fmt.Sprintf("Error updating customer: %v", err)), nil
		}
		return okResult("Customer updated:", jsonPretty(m)), nil
	})
	add("quickbook_delete_customer", "Delete (make inactive) a customer in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		qb := qbo.NewClient(cfg)
		m, err := qb.DeleteCustomerOrDeactivate(ctx, p["idOrEntity"])
		if err != nil {
			return errResult(fmt.Sprintf("Error deleting customer: %v", err)), nil
		}
		return okResult("Customer deleted:", jsonPretty(m)), nil
	})
	add("quickbook_search_customers", "Search customers in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		return runSearch(ctx, cfg, "customer", p)
	})

	// --- Estimates ---
	add("quickbook_create_estimate", "Create an estimate in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["estimate"])
		if !ok {
			return errResult("params.estimate object required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.CreateEntity(ctx, "estimate", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error creating estimate: %v", err)), nil
		}
		return okResult("Estimate created:", jsonPretty(m)), nil
	})
	add("quickbook_get_estimate", "Get an estimate by Id.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		id, _ := p["id"].(string)
		if id == "" {
			return errResult("id required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.ReadEntity(ctx, "estimate", id)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_update_estimate", "Update an estimate in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["estimate"])
		if !ok {
			return errResult("params.estimate object required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.UpdateEntity(ctx, "estimate", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult("Estimate updated:", jsonPretty(m)), nil
	})
	add("quickbook_delete_estimate", "Delete an estimate in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		qb := qbo.NewClient(cfg)
		raw, err := qb.Delete(ctx, "estimate", p["idOrEntity"])
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		m, _ := qbo.UnwrapEntity(raw, "estimate")
		return okResult("Deleted:", jsonPretty(m)), nil
	})
	add("quickbook_search_estimates", "Search estimates in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		return runSearch(ctx, cfg, "estimate", p)
	})

	// --- Bills ---
	add("quickbook_create_bill", "Create a bill in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["bill"])
		if !ok {
			return errResult("params.bill object required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.CreateEntity(ctx, "bill", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error creating bill: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_update_bill", "Update a bill in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["bill"])
		if !ok {
			return errResult("params.bill object required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.UpdateEntity(ctx, "bill", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_delete_bill", "Delete a bill in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["bill"])
		if !ok {
			return errResult("params.bill object with Id and SyncToken required"), nil
		}
		qb := qbo.NewClient(cfg)
		raw, err := qb.Delete(ctx, "bill", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		m, _ := qbo.UnwrapEntity(raw, "bill")
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_get_bill", "Get a bill by Id.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		id, _ := p["id"].(string)
		if id == "" {
			return errResult("id required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.ReadEntity(ctx, "bill", id)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_search_bills", "Search bills in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		return runSearch(ctx, cfg, "bill", p)
	})

	// --- Invoices ---
	add("quickbook_read_invoice", "Read a single invoice by ID.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		id, _ := p["invoice_id"].(string)
		if id == "" {
			return errResult("invoice_id required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.ReadEntity(ctx, "invoice", id)
		if err != nil {
			return errResult(fmt.Sprintf("Error reading invoice: %v", err)), nil
		}
		return okResult(fmt.Sprintf("Invoice details for ID %s:", id), jsonPretty(m)), nil
	})
	add("quickbook_search_invoices", "Search invoices in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		return runSearch(ctx, cfg, "invoice", p)
	})
	add("quickbook_create_invoice", "Create an invoice in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		customerRef, _ := p["customer_ref"].(string)
		if customerRef == "" {
			return errResult("customer_ref required"), nil
		}
		lines, _ := p["line_items"].([]any)
		if len(lines) == 0 {
			return errResult("line_items required"), nil
		}
		var lineSlice []any
		for i, li := range lines {
			lm, ok := li.(map[string]any)
			if !ok {
				continue
			}
			itemRef, _ := lm["item_ref"].(string)
			qty := num(lm["qty"])
			price := num(lm["unit_price"])
			desc, _ := lm["description"].(string)
			amt := qty * price
			lineSlice = append(lineSlice, map[string]any{
				"Id":          fmt.Sprintf("%d", i+1),
				"LineNum":     i + 1,
				"Description": desc,
				"Amount":      amt,
				"DetailType":  "SalesItemLineDetail",
				"SalesItemLineDetail": map[string]any{
					"ItemRef":   map[string]any{"value": itemRef},
					"Qty":       qty,
					"UnitPrice": price,
				},
			})
		}
		payload := map[string]any{
			"CustomerRef": map[string]any{"value": customerRef},
			"Line":        lineSlice,
		}
		if d, ok := p["doc_number"].(string); ok && d != "" {
			payload["DocNumber"] = d
		}
		if t, ok := p["txn_date"].(string); ok && t != "" {
			payload["TxnDate"] = t
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.CreateEntity(ctx, "invoice", payload)
		if err != nil {
			return errResult(fmt.Sprintf("Error creating invoice: %v", err)), nil
		}
		return okResult("Invoice created:", jsonPretty(m)), nil
	})
	add("quickbook_update_invoice", "Update an invoice by ID (sparse merge).", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		id, _ := p["invoice_id"].(string)
		if id == "" {
			return errResult("invoice_id required"), nil
		}
		patch, ok := toMap(p["patch"])
		if !ok {
			return errResult("patch object required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := mergeUpdate(ctx, qb, "invoice", id, patch)
		if err != nil {
			return errResult(fmt.Sprintf("Error updating invoice: %v", err)), nil
		}
		return okResult("Invoice updated:", jsonPretty(m)), nil
	})

	// --- Accounts ---
	add("quickbook_create_account", "Create an account in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		name, _ := p["name"].(string)
		typ, _ := p["type"].(string)
		if name == "" || typ == "" {
			return errResult("name and type required"), nil
		}
		pl := map[string]any{"Name": name, "AccountType": typ}
		if s, ok := p["sub_type"].(string); ok && s != "" {
			pl["AccountSubType"] = s
		}
		if s, ok := p["description"].(string); ok && s != "" {
			pl["Description"] = s
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.CreateEntity(ctx, "account", pl)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult("Account created:", jsonPretty(m)), nil
	})
	add("quickbook_update_account", "Update an account by ID (sparse merge).", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		id, _ := p["account_id"].(string)
		if id == "" {
			return errResult("account_id required"), nil
		}
		patch, ok := toMap(p["patch"])
		if !ok {
			return errResult("patch object required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := mergeUpdate(ctx, qb, "account", id, patch)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult("Account updated:", jsonPretty(m)), nil
	})
	add("quickbook_search_accounts", "Search accounts in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		return runSearch(ctx, cfg, "account", p)
	})

	// --- Items ---
	add("quickbook_read_item", "Read an item by ID.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		id, _ := p["item_id"].(string)
		if id == "" {
			return errResult("item_id required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.ReadEntity(ctx, "item", id)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(fmt.Sprintf("Item %s:", id), jsonPretty(m)), nil
	})
	add("quickbook_search_items", "Search items in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		return runSearch(ctx, cfg, "item", p)
	})
	add("quickbook_create_item", "Create an item in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		name, _ := p["name"].(string)
		typ, _ := p["type"].(string)
		inc, _ := p["income_account_ref"].(string)
		if name == "" || typ == "" || inc == "" {
			return errResult("name, type, and income_account_ref required"), nil
		}
		pl := map[string]any{
			"Name":             name,
			"Type":             typ,
			"IncomeAccountRef": map[string]any{"value": inc},
		}
		if e, ok := p["expense_account_ref"].(string); ok && e != "" {
			pl["ExpenseAccountRef"] = map[string]any{"value": e}
		}
		if p["unit_price"] != nil {
			pl["UnitPrice"] = num(p["unit_price"])
		}
		if d, ok := p["description"].(string); ok && d != "" {
			pl["Description"] = d
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.CreateEntity(ctx, "item", pl)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult("Item created:", jsonPretty(m)), nil
	})
	add("quickbook_update_item", "Update an item by ID (sparse merge).", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		id, _ := p["item_id"].(string)
		if id == "" {
			return errResult("item_id required"), nil
		}
		patch, ok := toMap(p["patch"])
		if !ok {
			return errResult("patch required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := mergeUpdate(ctx, qb, "item", id, patch)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult("Item updated:", jsonPretty(m)), nil
	})

	// --- Vendors ---
	add("quickbook_create_vendor", "Create a vendor in QuickBooks Online.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["vendor"])
		if !ok {
			return errResult("params.vendor object required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.CreateEntity(ctx, "vendor", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_update_vendor", "Update a vendor.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["vendor"])
		if !ok {
			return errResult("params.vendor required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.UpdateEntity(ctx, "vendor", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_delete_vendor", "Delete a vendor.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		qb := qbo.NewClient(cfg)
		raw, err := qb.Delete(ctx, "vendor", p["idOrEntity"])
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		m, _ := qbo.UnwrapEntity(raw, "vendor")
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_get_vendor", "Get a vendor by ID.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		id, _ := p["id"].(string)
		if id == "" {
			return errResult("id required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.ReadEntity(ctx, "vendor", id)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_search_vendors", "Search vendors.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		return runSearch(ctx, cfg, "vendor", p)
	})

	// --- Employees ---
	add("quickbook_create_employee", "Create an employee.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["employee"])
		if !ok {
			return errResult("params.employee required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.CreateEntity(ctx, "employee", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult("Employee created:", jsonPretty(m)), nil
	})
	add("quickbook_get_employee", "Get an employee by ID.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		id, _ := p["id"].(string)
		if id == "" {
			return errResult("id required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.ReadEntity(ctx, "employee", id)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_update_employee", "Update an employee.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["employee"])
		if !ok {
			return errResult("params.employee required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.UpdateEntity(ctx, "employee", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_search_employees", "Search employees.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		return runSearch(ctx, cfg, "employee", p)
	})

	// --- Journal entries ---
	add("quickbook_create_journal_entry", "Create a journal entry.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["journalEntry"])
		if !ok {
			obj, ok = toMap(p["journal_entry"])
		}
		if !ok {
			return errResult("params.journalEntry object required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.CreateEntity(ctx, "journalEntry", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_get_journal_entry", "Get a journal entry by ID.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		id, _ := p["id"].(string)
		if id == "" {
			return errResult("id required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.ReadEntity(ctx, "journalEntry", id)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_update_journal_entry", "Update a journal entry.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["journalEntry"])
		if !ok {
			obj, ok = toMap(p["journal_entry"])
		}
		if !ok {
			return errResult("params.journalEntry required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.UpdateEntity(ctx, "journalEntry", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_delete_journal_entry", "Delete a journal entry.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		qb := qbo.NewClient(cfg)
		raw, err := qb.Delete(ctx, "journalEntry", p["idOrEntity"])
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		m, _ := qbo.UnwrapEntity(raw, "journalEntry")
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_search_journal_entries", "Search journal entries.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		return runSearch(ctx, cfg, "journalEntry", p)
	})

	// --- Bill payments ---
	add("quickbook_create_bill_payment", "Create a bill payment.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["billPayment"])
		if !ok {
			obj, ok = toMap(p["bill_payment"])
		}
		if !ok {
			return errResult("params.billPayment required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.CreateEntity(ctx, "billPayment", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_get_bill_payment", "Get a bill payment by ID.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		id, _ := p["id"].(string)
		if id == "" {
			return errResult("id required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.ReadEntity(ctx, "billPayment", id)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_update_bill_payment", "Update a bill payment.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["billPayment"])
		if !ok {
			obj, ok = toMap(p["bill_payment"])
		}
		if !ok {
			return errResult("params.billPayment required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.UpdateEntity(ctx, "billPayment", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_delete_bill_payment", "Delete a bill payment.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		qb := qbo.NewClient(cfg)
		raw, err := qb.Delete(ctx, "billPayment", p["idOrEntity"])
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		m, _ := qbo.UnwrapEntity(raw, "billPayment")
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_search_bill_payments", "Search bill payments.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		return runSearch(ctx, cfg, "billPayment", p)
	})

	// --- Purchases ---
	add("quickbook_create_purchase", "Create a purchase.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["purchase"])
		if !ok {
			return errResult("params.purchase required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.CreateEntity(ctx, "purchase", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_get_purchase", "Get a purchase by ID.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		id, _ := p["id"].(string)
		if id == "" {
			return errResult("id required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.ReadEntity(ctx, "purchase", id)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_update_purchase", "Update a purchase.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		obj, ok := toMap(p["purchase"])
		if !ok {
			return errResult("params.purchase required"), nil
		}
		qb := qbo.NewClient(cfg)
		m, err := qb.UpdateEntity(ctx, "purchase", obj)
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_delete_purchase", "Delete a purchase.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		qb := qbo.NewClient(cfg)
		raw, err := qb.Delete(ctx, "purchase", p["idOrEntity"])
		if err != nil {
			return errResult(fmt.Sprintf("Error: %v", err)), nil
		}
		m, _ := qbo.UnwrapEntity(raw, "purchase")
		return okResult(jsonPretty(m)), nil
	})
	add("quickbook_search_purchases", "Search purchases.", func(ctx context.Context, p map[string]any) (*mcp.CallToolResult, error) {
		return runSearch(ctx, cfg, "purchase", p)
	})
}

func runSearch(ctx context.Context, cfg *config.Config, qEntity string, p map[string]any) (*mcp.CallToolResult, error) {
	qb := qbo.NewClient(cfg)
	crit := qbo.BuildSearchCriteria(p["criteria"])
	raw, err := qb.Query(ctx, qEntity, crit)
	if err != nil {
		return errResult(fmt.Sprintf("Error searching: %v", err)), nil
	}
	items, _, err := qbo.ExtractQueryEntities(raw, qbo.CapitalizeEntity(qEntity))
	if err != nil {
		return errResult(err.Error()), nil
	}
	return okResult(fmt.Sprintf("Found %d results", len(items)), jsonPretty(items)), nil
}

func mergeUpdate(ctx context.Context, qb *qbo.Client, entity, id string, patch map[string]any) (map[string]any, error) {
	ex, err := qb.ReadEntity(ctx, entity, id)
	if err != nil {
		return nil, err
	}
	for k, v := range patch {
		ex[k] = v
	}
	ex["Id"] = id
	ex["sparse"] = true
	return qb.UpdateEntity(ctx, entity, ex)
}
