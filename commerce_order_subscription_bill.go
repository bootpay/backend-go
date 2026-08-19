package bootpay

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// OrderSubscriptionBillModule handles order subscription bill-related operations
type OrderSubscriptionBillModule struct {
	api *CommerceApi
}

// List retrieves order subscription bill list
// GET /v1/order_subscription_bills
// ⚠️ The path is order_subscription_bills — underscores (not hyphens).
// page/limit default to 1/20 when unset.
func (m *OrderSubscriptionBillModule) List(params *OrderSubscriptionBillListParams) (map[string]interface{}, error) {
	if params == nil {
		params = &OrderSubscriptionBillListParams{}
	}
	queryParams := url.Values{}
	if params.OrderSubscriptionId != "" {
		queryParams.Set("order_subscription_id", params.OrderSubscriptionId)
	}
	page := params.Page
	if page <= 0 {
		page = 1
	}
	limit := params.Limit
	if limit <= 0 {
		limit = 20
	}
	queryParams.Set("page", strconv.Itoa(page))
	queryParams.Set("limit", strconv.Itoa(limit))
	if params.Keyword != "" {
		queryParams.Set("keyword", params.Keyword)
	}
	if len(params.Status) > 0 {
		statusStrs := make([]string, len(params.Status))
		for i, s := range params.Status {
			statusStrs[i] = strconv.Itoa(s)
		}
		queryParams.Set("status", strings.Join(statusStrs, ","))
	}
	return m.api.getWithHeaders("order_subscription_bills?"+queryParams.Encode(), commerceRoleHeaders("user", params.IdempotencyKey))
}

// Detail retrieves order subscription bill details
func (m *OrderSubscriptionBillModule) Detail(orderSubscriptionBillId string) (map[string]interface{}, error) {
	return m.api.Get(fmt.Sprintf("order_subscription_bills/%s", orderSubscriptionBillId))
}

// Update updates order subscription bill
func (m *OrderSubscriptionBillModule) Update(orderSubscriptionBill CommerceOrderSubscriptionBill) (map[string]interface{}, error) {
	if orderSubscriptionBill.OrderSubscriptionBillId == "" {
		return nil, fmt.Errorf("order_subscription_bill_id is required")
	}
	return m.api.Put(fmt.Sprintf("order_subscription_bills/%s", orderSubscriptionBill.OrderSubscriptionBillId), orderSubscriptionBill)
}
