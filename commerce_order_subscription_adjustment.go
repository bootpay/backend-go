package bootpay

import (
	"fmt"
)

// OrderSubscriptionAdjustmentModule handles order subscription adjustment-related operations
type OrderSubscriptionAdjustmentModule struct {
	api *CommerceApi
}

// Create creates a new order subscription adjustment
func (m *OrderSubscriptionAdjustmentModule) Create(orderSubscriptionId string, adjustment CommerceOrderSubscriptionAdjustment) (map[string]interface{}, error) {
	return m.api.Post(fmt.Sprintf("order_subscriptions/%s/adjustments", orderSubscriptionId), adjustment)
}

// Update updates an order subscription adjustment
func (m *OrderSubscriptionAdjustmentModule) Update(params OrderSubscriptionAdjustmentUpdateParams) (map[string]interface{}, error) {
	if params.OrderSubscriptionId == "" {
		return nil, fmt.Errorf("order_subscription_id is required")
	}
	return m.api.Put(fmt.Sprintf("order_subscriptions/%s/adjustments", params.OrderSubscriptionId), params)
}

// Delete deletes an order subscription adjustment
func (m *OrderSubscriptionAdjustmentModule) Delete(orderSubscriptionId string, orderSubscriptionAdjustmentId string) (map[string]interface{}, error) {
	return m.api.Delete(fmt.Sprintf("order_subscriptions/%s/adjustments?order_subscription_adjustment_id=%s", orderSubscriptionId, orderSubscriptionAdjustmentId))
}
