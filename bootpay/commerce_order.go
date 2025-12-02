package bootpay

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// OrderModule handles order-related operations
type OrderModule struct {
	api *CommerceApi
}

// List retrieves order list
func (m *OrderModule) List(params *OrderListParams) (map[string]interface{}, error) {
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
		if params.UserId != "" {
			queryParams.Set("user_id", params.UserId)
		}
		if params.UserGroupId != "" {
			queryParams.Set("user_group_id", params.UserGroupId)
		}
		if params.CsType != "" {
			queryParams.Set("cs_type", params.CsType)
		}
		if params.CssAt != "" {
			queryParams.Set("css_at", params.CssAt)
		}
		if params.CseAt != "" {
			queryParams.Set("cse_at", params.CseAt)
		}
		if params.SubscriptionBillingType > 0 {
			queryParams.Set("subscription_billing_type", strconv.Itoa(params.SubscriptionBillingType))
		}
		if len(params.Status) > 0 {
			statusStrs := make([]string, len(params.Status))
			for i, s := range params.Status {
				statusStrs[i] = strconv.Itoa(s)
			}
			queryParams.Set("status", strings.Join(statusStrs, ","))
		}
		if len(params.PaymentStatus) > 0 {
			paymentStatusStrs := make([]string, len(params.PaymentStatus))
			for i, s := range params.PaymentStatus {
				paymentStatusStrs[i] = strconv.Itoa(s)
			}
			queryParams.Set("payment_status", strings.Join(paymentStatusStrs, ","))
		}
		if len(params.OrderSubscriptionIds) > 0 {
			queryParams.Set("order_subscription_ids", strings.Join(params.OrderSubscriptionIds, ","))
		}
		if len(queryParams) > 0 {
			query = "?" + queryParams.Encode()
		}
	}
	return m.api.Get("orders" + query)
}

// Detail retrieves order details
func (m *OrderModule) Detail(orderId string) (map[string]interface{}, error) {
	return m.api.Get(fmt.Sprintf("orders/%s", orderId))
}

// Month retrieves monthly orders
func (m *OrderModule) Month(userGroupId string, searchDate string) (map[string]interface{}, error) {
	queryParams := url.Values{}
	queryParams.Set("user_group_id", userGroupId)
	queryParams.Set("search_date", searchDate)
	return m.api.Get("orders/month?" + queryParams.Encode())
}
