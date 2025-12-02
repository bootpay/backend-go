package bootpay

import (
	"fmt"
	"net/url"
)

// OrderCancelModule handles order cancel-related operations
type OrderCancelModule struct {
	api *CommerceApi
}

// List retrieves order cancel request list
func (m *OrderCancelModule) List(params *OrderCancelListParams) (map[string]interface{}, error) {
	query := ""
	if params != nil {
		queryParams := url.Values{}
		if params.OrderId != "" {
			queryParams.Set("order_id", params.OrderId)
		}
		if params.OrderNumber != "" {
			queryParams.Set("order_number", params.OrderNumber)
		}
		if len(queryParams) > 0 {
			query = "?" + queryParams.Encode()
		}
	}
	return m.api.Get("order/cancel" + query)
}

// Request creates a cancel request
func (m *OrderCancelModule) Request(params OrderCancelParams) (map[string]interface{}, error) {
	return m.api.Post("order/cancel", params)
}

// Withdraw withdraws a cancel request
func (m *OrderCancelModule) Withdraw(orderCancelRequestHistoryId string) (map[string]interface{}, error) {
	return m.api.Put(fmt.Sprintf("order/cancel/%s/withdraw", orderCancelRequestHistoryId), map[string]interface{}{})
}

// Approve approves a cancel request
func (m *OrderCancelModule) Approve(params OrderCancelActionParams) (map[string]interface{}, error) {
	if params.OrderCancelRequestHistoryId == "" {
		return nil, fmt.Errorf("order_cancel_request_history_id is required")
	}
	return m.api.Put(fmt.Sprintf("order/cancel/%s/approve", params.OrderCancelRequestHistoryId), params)
}

// Reject rejects a cancel request
func (m *OrderCancelModule) Reject(params OrderCancelActionParams) (map[string]interface{}, error) {
	if params.OrderCancelRequestHistoryId == "" {
		return nil, fmt.Errorf("order_cancel_request_history_id is required")
	}
	return m.api.Put(fmt.Sprintf("order/cancel/%s/reject", params.OrderCancelRequestHistoryId), params)
}
