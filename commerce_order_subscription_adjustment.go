package bootpay

import (
	"fmt"
)

// OrderSubscriptionAdjustmentModule handles order subscription adjustment-related operations
//
// ⚠️ POST · PUT · DELETE all share the single /adjustments path —
// do not infer the method from the path alone.
// The server requires supervisor scope for all three.
type OrderSubscriptionAdjustmentModule struct {
	api *CommerceApi
}

// Create creates a new order subscription adjustment (supervisor scope)
// POST /v1/order_subscriptions/{order_subscription_id}/adjustments
// price/duration/tax_free_price default to 0 / 1 / 0 when unset.
// When type is omitted the server auto-detects: price > 0 → SETUP_PRICE, otherwise PERIOD_DISCOUNT.
//
// Three ways to target the billing cycles (widest last):
//   - Duration: 5                        → the 5th cycle only
//   - DurationFrom: 3, DurationTo: 7     → cycles 3~7, one record each (5 records)
//   - DurationFrom: 3, IsUnlimited: true → from cycle 3 to the end of the contract
//     (a single record, DurationTo is ignored)
//
// The upper bound is the contract's total duration; contracts with an unlimited total
// stop at cycle 60. Already-paid cycles are rejected, and if any single cycle in the
// range would end up with a negative amount the whole request is rejected (no partial apply).
func (m *OrderSubscriptionAdjustmentModule) Create(orderSubscriptionId string, adjustment CommerceOrderSubscriptionAdjustment) (map[string]interface{}, error) {
	duration := adjustment.Duration
	if duration == 0 {
		duration = 1
	}
	body := map[string]interface{}{
		"price":          adjustment.Price,
		"duration":       duration,
		"tax_free_price": adjustment.TaxFreePrice,
	}
	if adjustment.Name != "" {
		body["name"] = adjustment.Name
	}
	if adjustment.Type != 0 {
		body["type"] = adjustment.Type
	}
	if adjustment.DurationFrom != 0 {
		body["duration_from"] = adjustment.DurationFrom
	}
	if adjustment.DurationTo != 0 {
		body["duration_to"] = adjustment.DurationTo
	}
	if adjustment.IsUnlimited != nil {
		body["is_unlimited"] = *adjustment.IsUnlimited
	}
	return m.api.postWithHeaders(fmt.Sprintf("order_subscriptions/%s/adjustments", orderSubscriptionId), body, commerceRoleHeaders("supervisor", ""))
}

// Update replaces the adjustments of a given duration as a whole (supervisor scope)
// PUT /v1/order_subscriptions/{order_subscription_id}/adjustments
// The server swaps the adjustments array per Duration (defaults to 1 when unset).
func (m *OrderSubscriptionAdjustmentModule) Update(params OrderSubscriptionAdjustmentUpdateParams) (map[string]interface{}, error) {
	if params.OrderSubscriptionId == "" {
		return nil, fmt.Errorf("order_subscription_id is required")
	}
	if params.Duration == 0 {
		params.Duration = 1
	}
	return m.api.putWithHeaders(fmt.Sprintf("order_subscriptions/%s/adjustments", params.OrderSubscriptionId), params, commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// Delete deletes an order subscription adjustment (supervisor scope)
// DELETE /v1/order_subscriptions/{order_subscription_id}/adjustments
// ⚠️ The target id is sent in the body, not the query.
func (m *OrderSubscriptionAdjustmentModule) Delete(orderSubscriptionId string, orderSubscriptionAdjustmentId string) (map[string]interface{}, error) {
	body := map[string]interface{}{
		"order_subscription_adjustment_id": orderSubscriptionAdjustmentId,
	}
	return m.api.deleteWithHeaders(fmt.Sprintf("order_subscriptions/%s/adjustments", orderSubscriptionId), body, commerceRoleHeaders("supervisor", ""))
}
