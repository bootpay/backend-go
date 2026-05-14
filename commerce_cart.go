package bootpay

// CartModule handles cart-related operations
type CartModule struct {
	api *CommerceApi
}

// OrderPreview previews an order (authoritative shipping/discount calculation)
// POST /v1/cart/order-preview
//
// member_mode='guest' (default): cart_items required
// member_mode='member': uses server cart (user token required)
func (m *CartModule) OrderPreview(params *OrderPreviewParams) (map[string]interface{}, error) {
	body := params
	if body == nil {
		body = &OrderPreviewParams{}
	}
	return m.api.Post("cart/order-preview", body)
}
