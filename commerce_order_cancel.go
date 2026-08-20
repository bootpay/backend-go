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
// GET /v1/order/cancel — filtered by order_number or order_id (all when both empty).
// The order_cancellation_request_id to pass to Approve / Reject / Withdraw comes from here.
func (m *OrderCancelModule) List(params *OrderCancelListParams) (map[string]interface{}, error) {
	query := ""
	idempotencyKey := ""
	if params != nil {
		idempotencyKey = params.IdempotencyKey
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
	return m.api.getWithHeaders("order/cancel"+query, commerceRoleHeaders("user", idempotencyKey))
}

// Request creates a cancel request
func (m *OrderCancelModule) Request(params OrderCancelParams) (map[string]interface{}, error) {
	return m.api.Post("order/cancel", params)
}

// Withdraw withdraws a cancel request (buyer side)
// PUT /v1/order/cancel/{order_cancellation_request_id}/withdraw
// ⚠️ A different route from DELETE /v1/order/cancel/{id} — the documented one is withdraw.
func (m *OrderCancelModule) Withdraw(orderCancelRequestHistoryId string) (map[string]interface{}, error) {
	return m.api.putWithHeaders(fmt.Sprintf("order/cancel/%s/withdraw", orderCancelRequestHistoryId), map[string]interface{}{}, commerceRoleHeaders("user", ""))
}

// Approve approves a cancel request (supervisor scope)
// PUT /v1/order/cancel/{order_cancellation_request_id}/approve
func (m *OrderCancelModule) Approve(params OrderCancelActionParams) (map[string]interface{}, error) {
	cancellationId := orderCancellationId(params)
	if cancellationId == "" {
		return nil, fmt.Errorf("order_cancellation_request_id is required")
	}
	return m.api.putWithHeaders(fmt.Sprintf("order/cancel/%s/approve", cancellationId), params, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// Reject rejects a cancel request (supervisor scope)
// PUT /v1/order/cancel/{order_cancellation_request_id}/reject
func (m *OrderCancelModule) Reject(params OrderCancelActionParams) (map[string]interface{}, error) {
	cancellationId := orderCancellationId(params)
	if cancellationId == "" {
		return nil, fmt.Errorf("order_cancellation_request_id is required")
	}
	return m.api.putWithHeaders(fmt.Sprintf("order/cancel/%s/reject", cancellationId), params, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// orderCancellationId resolves the cancel request history id.
// The official name is OrderCancellationRequestId; the old OrderCancelRequestHistoryId
// keeps working (when both are set, the official one wins).
func orderCancellationId(params OrderCancelActionParams) string {
	if params.OrderCancellationRequestId != "" {
		return params.OrderCancellationRequestId
	}
	return params.OrderCancelRequestHistoryId
}
