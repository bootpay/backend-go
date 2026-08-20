package bootpay

import (
	"fmt"
	"net/url"
	"strconv"
)

// InvoiceModule handles invoice-related operations
type InvoiceModule struct {
	api *CommerceApi
}

// List retrieves invoice list
// GET /v1/invoices — kept for backward compatibility; use ListWithParams for the full filter set.
// ⚠️ The response data is { list: [...], count: N } — not { items, total }.
func (m *InvoiceModule) List(params *ListParams) (map[string]interface{}, error) {
	converted := &InvoiceListParams{}
	if params != nil {
		converted.ListParams = *params
	}
	return m.ListWithParams(converted)
}

// ListWithParams retrieves invoice list with the full filter set
// GET /v1/invoices
// page/limit default to 1/24 when unset (24 matches the server default).
// ⚠️ The response data is { list: [...], count: N } — not { items, total }.
func (m *InvoiceModule) ListWithParams(params *InvoiceListParams) (map[string]interface{}, error) {
	if params == nil {
		params = &InvoiceListParams{}
	}
	queryParams := url.Values{}
	page := params.Page
	if page <= 0 {
		page = 1
	}
	limit := params.Limit
	if limit <= 0 {
		limit = 24
	}
	queryParams.Set("page", strconv.Itoa(page))
	queryParams.Set("limit", strconv.Itoa(limit))
	if params.Keyword != "" {
		queryParams.Set("keyword", params.Keyword)
	}
	if params.CsType != "" {
		queryParams.Set("cs_type", params.CsType)
	}
	if params.UserId != "" {
		queryParams.Set("user_id", params.UserId)
	}
	if params.ProductType != nil {
		queryParams.Set("product_type", strconv.Itoa(*params.ProductType))
	}
	if params.CssAt != "" {
		queryParams.Set("css_at", params.CssAt)
	}
	if params.CseAt != "" {
		queryParams.Set("cse_at", params.CseAt)
	}
	return m.api.getWithHeaders("invoices?"+queryParams.Encode(), commerceRoleHeaders("user", params.IdempotencyKey))
}

// Create creates a new invoice
func (m *InvoiceModule) Create(invoice CommerceInvoice) (map[string]interface{}, error) {
	return m.api.Post("invoices", invoice)
}

// Notify sends invoice notification
// POST /v1/invoices/{invoice_id}/notify
// sendTypes may be nil — the key is then omitted and the server treats it as an empty array.
// ⚠️ Real customer notifications are sent — be careful with test calls.
func (m *InvoiceModule) Notify(invoiceId string, sendTypes []int) (map[string]interface{}, error) {
	data := map[string]interface{}{}
	if sendTypes != nil {
		data["send_types"] = sendTypes
	}
	return m.api.postWithHeaders(fmt.Sprintf("invoices/%s/notify", invoiceId), data, commerceRoleHeaders("user", ""))
}

// Detail retrieves invoice details
// GET /v1/invoices/{invoice_id}
func (m *InvoiceModule) Detail(invoiceId string) (map[string]interface{}, error) {
	return m.api.getWithHeaders(fmt.Sprintf("invoices/%s", invoiceId), commerceRoleHeaders("user", ""))
}
