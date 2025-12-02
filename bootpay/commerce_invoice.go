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
func (m *InvoiceModule) List(params *ListParams) (map[string]interface{}, error) {
	query := ""
	if params != nil {
		queryParams := url.Values{}
		if params.Page > 0 {
			queryParams.Set("page", strconv.Itoa(params.Page))
		}
		if params.Limit > 0 {
			queryParams.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.Keyword != "" {
			queryParams.Set("keyword", params.Keyword)
		}
		if len(queryParams) > 0 {
			query = "?" + queryParams.Encode()
		}
	}
	return m.api.Get("invoices" + query)
}

// Create creates a new invoice
func (m *InvoiceModule) Create(invoice CommerceInvoice) (map[string]interface{}, error) {
	return m.api.Post("invoices", invoice)
}

// Notify sends invoice notification
func (m *InvoiceModule) Notify(invoiceId string, sendTypes []int) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"send_types": sendTypes,
	}
	return m.api.Post(fmt.Sprintf("invoices/%s/notify", invoiceId), data)
}

// Detail retrieves invoice details
func (m *InvoiceModule) Detail(invoiceId string) (map[string]interface{}, error) {
	return m.api.Get(fmt.Sprintf("invoices/%s", invoiceId))
}
