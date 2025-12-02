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
func (m *OrderSubscriptionBillModule) List(params *OrderSubscriptionBillListParams) (map[string]interface{}, error) {
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
		if params.OrderSubscriptionId != "" {
			queryParams.Set("order_subscription_id", params.OrderSubscriptionId)
		}
		if len(params.Status) > 0 {
			statusStrs := make([]string, len(params.Status))
			for i, s := range params.Status {
				statusStrs[i] = strconv.Itoa(s)
			}
			queryParams.Set("status", strings.Join(statusStrs, ","))
		}
		if len(queryParams) > 0 {
			query = "?" + queryParams.Encode()
		}
	}
	return m.api.Get("order_subscription_bills" + query)
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
