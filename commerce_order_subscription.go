package bootpay

import (
	"fmt"
	"net/url"
	"strconv"
)

// OrderSubscriptionRequestIngModule handles subscription request ing operations
type OrderSubscriptionRequestIngModule struct {
	api *CommerceApi
}

// Pause pauses a subscription
// POST /v1/order_subscriptions/requests/ing/pause (user scope — buyer-side request)
func (m *OrderSubscriptionRequestIngModule) Pause(params OrderSubscriptionPauseParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("order_subscriptions/requests/ing/pause", params, commerceRoleHeaders("user", params.IdempotencyKey))
}

// Resume resumes a paused subscription
// PUT /v1/order_subscriptions/requests/ing/resume
// ⚠️ The only PUT in the requests/ing family — do not "fix" it to POST.
func (m *OrderSubscriptionRequestIngModule) Resume(params OrderSubscriptionResumeParams) (map[string]interface{}, error) {
	return m.api.putWithHeaders("order_subscriptions/requests/ing/resume", params, commerceRoleHeaders("user", params.IdempotencyKey))
}

// Purchase requests a mid-contract purchase
// POST /v1/order_subscriptions/requests/ing/purchase
func (m *OrderSubscriptionRequestIngModule) Purchase(params OrderSubscriptionPurchaseParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("order_subscriptions/requests/ing/purchase", params, commerceRoleHeaders("user", params.IdempotencyKey))
}

// Transfer requests a subscription transfer/succession
// POST /v1/order_subscriptions/requests/ing/transfer
func (m *OrderSubscriptionRequestIngModule) Transfer(params OrderSubscriptionTransferParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("order_subscriptions/requests/ing/transfer", params, commerceRoleHeaders("user", params.IdempotencyKey))
}

// CalculateTerminationFee calculates termination fee
func (m *OrderSubscriptionRequestIngModule) CalculateTerminationFee(orderSubscriptionId string, orderNumber string) (map[string]interface{}, error) {
	if orderSubscriptionId == "" && orderNumber == "" {
		return nil, fmt.Errorf("orderSubscriptionId or orderNumber is required")
	}

	queryParams := url.Values{}
	if orderSubscriptionId != "" {
		queryParams.Set("order_subscription_id", orderSubscriptionId)
	}
	if orderNumber != "" {
		queryParams.Set("order_number", orderNumber)
	}

	return m.api.getWithHeaders("order_subscriptions/requests/ing/calculate_termination_fee?"+queryParams.Encode(), commerceRoleHeaders("user", ""))
}

// CalculateTerminationFeeByOrderNumber calculates termination fee by order number
func (m *OrderSubscriptionRequestIngModule) CalculateTerminationFeeByOrderNumber(orderNumber string) (map[string]interface{}, error) {
	return m.CalculateTerminationFee("", orderNumber)
}

// Termination terminates a subscription
// POST /v1/order_subscriptions/requests/ing/termination
func (m *OrderSubscriptionRequestIngModule) Termination(params OrderSubscriptionTerminationParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("order_subscriptions/requests/ing/termination", params, commerceRoleHeaders("user", params.IdempotencyKey))
}

// OrderSubscriptionModule handles order subscription-related operations
type OrderSubscriptionModule struct {
	api        *CommerceApi
	RequestIng *OrderSubscriptionRequestIngModule
}

// List retrieves order subscription list
func (m *OrderSubscriptionModule) List(params *OrderSubscriptionListParams) (map[string]interface{}, error) {
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
		if params.SearchDateFrom != "" {
			queryParams.Set("search_date_from", params.SearchDateFrom)
		}
		if params.SearchDateTo != "" {
			queryParams.Set("search_date_to", params.SearchDateTo)
		}
		if params.SAt != "" {
			queryParams.Set("s_at", params.SAt)
		}
		if params.EAt != "" {
			queryParams.Set("e_at", params.EAt)
		}
		if params.RequestType != "" {
			queryParams.Set("request_type", params.RequestType)
		}
		if params.UserGroupId != "" {
			queryParams.Set("user_group_id", params.UserGroupId)
		}
		if params.Status != nil {
			queryParams.Set("status", strconv.Itoa(*params.Status))
		}
		if params.UserId != "" {
			queryParams.Set("user_id", params.UserId)
		}
		if params.OrderNumber != "" {
			queryParams.Set("order_number", params.OrderNumber)
		}
		if len(queryParams) > 0 {
			query = "?" + queryParams.Encode()
		}
	}
	return m.api.Get("order_subscriptions" + query)
}

// Detail retrieves order subscription details
func (m *OrderSubscriptionModule) Detail(orderSubscriptionId string) (map[string]interface{}, error) {
	return m.api.Get(fmt.Sprintf("order_subscriptions/%s", orderSubscriptionId))
}

// Update updates the subscription contract (supervisor scope)
// PUT /v1/order_subscriptions/{order_subscription_id}
// Only changed values need to be sent (the server keeps the rest as-is).
//
// Price is the base charge amount per billing cycle. Changing it immediately recalculates
// the amount of the READY (scheduled) cycle, and every cycle created afterwards uses it too.
// Already-paid cycles are untouched. Values of 0 or less are rejected by the server.
// To add/subtract on specific cycles only, use OrderSubscriptionAdjustment.Create instead.
// (Same implementation as the amount change in the admin console.)
//
// Memo is the reason recorded on the change history (SUBSCRIPTION_ACTION_UPDATE).
func (m *OrderSubscriptionModule) Update(params OrderSubscriptionUpdateParams) (map[string]interface{}, error) {
	if params.OrderSubscriptionId == "" {
		return nil, fmt.Errorf("order_subscription_id is required")
	}
	return m.api.putWithHeaders(fmt.Sprintf("order_subscriptions/%s", params.OrderSubscriptionId), params, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// SupervisorApprove approves a subscription request (supervisor role)
// ⚠️ The server requires supervisor scope (scope_invalid!) — BOOTPAY-ROLE must be sent.
func (m *OrderSubscriptionModule) SupervisorApprove(orderSubscriptionId string, params *SupervisorOrderSubscriptionApproveParams) (map[string]interface{}, error) {
	if params == nil {
		params = &SupervisorOrderSubscriptionApproveParams{}
	}
	return m.api.putWithHeaders(fmt.Sprintf("order_subscriptions/%s/approve", orderSubscriptionId), params,
		commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// SupervisorReject rejects a subscription request (supervisor role)
// ⚠️ The server requires supervisor scope (scope_invalid!) — BOOTPAY-ROLE must be sent.
func (m *OrderSubscriptionModule) SupervisorReject(orderSubscriptionId string, params *SupervisorOrderSubscriptionRejectParams) (map[string]interface{}, error) {
	if params == nil {
		params = &SupervisorOrderSubscriptionRejectParams{}
	}
	return m.api.putWithHeaders(fmt.Sprintf("order_subscriptions/%s/reject", orderSubscriptionId), params,
		commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// SupervisorTerminate terminates a subscription (supervisor role)
// ⚠️ The server requires supervisor scope (scope_invalid!) — BOOTPAY-ROLE must be sent.
func (m *OrderSubscriptionModule) SupervisorTerminate(orderSubscriptionId string, params *SupervisorOrderSubscriptionTerminateParams) (map[string]interface{}, error) {
	if params == nil {
		params = &SupervisorOrderSubscriptionTerminateParams{}
	}
	return m.api.putWithHeaders(fmt.Sprintf("order_subscriptions/%s/terminate", orderSubscriptionId), params,
		commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// SupervisorPause pauses a subscription (supervisor role)
// ⚠️ The server requires supervisor scope (scope_invalid!) — BOOTPAY-ROLE must be sent.
func (m *OrderSubscriptionModule) SupervisorPause(orderSubscriptionId string, params SupervisorOrderSubscriptionPauseParams) (map[string]interface{}, error) {
	return m.api.putWithHeaders(fmt.Sprintf("order_subscriptions/%s/pause", orderSubscriptionId), params,
		commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// SupervisorResume resumes a subscription (supervisor role)
// ⚠️ The server requires supervisor scope (scope_invalid!) — BOOTPAY-ROLE must be sent.
func (m *OrderSubscriptionModule) SupervisorResume(orderSubscriptionId string, params *SupervisorOrderSubscriptionResumeParams) (map[string]interface{}, error) {
	if params == nil {
		params = &SupervisorOrderSubscriptionResumeParams{}
	}
	return m.api.putWithHeaders(fmt.Sprintf("order_subscriptions/%s/resume", orderSubscriptionId), params,
		commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// SupervisorCharge performs an on-demand charge_key payment (supervisor scope)
// POST /v1/order_subscriptions/charge
// ⚠️ charge_key is sent only in the body — never in the URL/query (access log exposure).
func (m *OrderSubscriptionModule) SupervisorCharge(params SupervisorOrderSubscriptionChargeParams) (map[string]interface{}, error) {
	if params.ChargeKey == "" {
		return nil, fmt.Errorf("charge_key is required")
	}
	return m.api.postWithHeaders("order_subscriptions/charge", params, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// SupervisorChargeRevoke revokes an on-demand charge_key (supervisor scope)
// DELETE /v1/order_subscriptions/charge — charge_key is sent in the body
// After revoking, the key can never be charged again.
func (m *OrderSubscriptionModule) SupervisorChargeRevoke(params SupervisorOrderSubscriptionChargeRevokeParams) (map[string]interface{}, error) {
	if params.ChargeKey == "" {
		return nil, fmt.Errorf("charge_key is required")
	}
	return m.api.deleteWithHeaders("order_subscriptions/charge", params, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}
