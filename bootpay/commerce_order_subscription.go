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
func (m *OrderSubscriptionRequestIngModule) Pause(params OrderSubscriptionPauseParams) (map[string]interface{}, error) {
	return m.api.Post("order_subscriptions/requests/ing/pause", params)
}

// Resume resumes a paused subscription
func (m *OrderSubscriptionRequestIngModule) Resume(params OrderSubscriptionResumeParams) (map[string]interface{}, error) {
	return m.api.Put("order_subscriptions/requests/ing/resume", params)
}

// CalculateTerminationFee calculates termination fee
func (m *OrderSubscriptionRequestIngModule) CalculateTerminationFee(orderSubscriptionId string, orderNumber string) (map[string]interface{}, error) {
	if orderSubscriptionId == "" && orderNumber == "" {
		return nil, fmt.Errorf("orderSubscriptionId or orderNumber is required")
	}

	queryParams := url.Values{}
	if orderSubscriptionId != "" {
		queryParams.Set("order_subscription_id", orderSubscriptionId)
	} else if orderNumber != "" {
		queryParams.Set("order_number", orderNumber)
	}

	return m.api.Get("order_subscriptions/requests/ing/calculate_termination_fee?" + queryParams.Encode())
}

// CalculateTerminationFeeByOrderNumber calculates termination fee by order number
func (m *OrderSubscriptionRequestIngModule) CalculateTerminationFeeByOrderNumber(orderNumber string) (map[string]interface{}, error) {
	return m.CalculateTerminationFee("", orderNumber)
}

// Termination terminates a subscription
func (m *OrderSubscriptionRequestIngModule) Termination(params OrderSubscriptionTerminationParams) (map[string]interface{}, error) {
	return m.api.Post("order_subscriptions/requests/ing/termination", params)
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
		if params.UserId != "" {
			queryParams.Set("user_id", params.UserId)
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

// Update updates order subscription
func (m *OrderSubscriptionModule) Update(params OrderSubscriptionUpdateParams) (map[string]interface{}, error) {
	if params.OrderSubscriptionId == "" {
		return nil, fmt.Errorf("order_subscription_id is required")
	}
	return m.api.Put(fmt.Sprintf("order_subscriptions/%s", params.OrderSubscriptionId), params)
}
