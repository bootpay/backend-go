package bootpay

// WebhookModule handles webhook-related operations
type WebhookModule struct {
	api *CommerceApi
}

// SendTest sends a test webhook to the registered webhook URL
// POST /v1/webhook/test
// params may be nil (header_content_type falls back to the server default).
func (m *WebhookModule) SendTest(params *SendTestWebhookParams) (map[string]interface{}, error) {
	body := &SendTestWebhookParams{}
	idempotencyKey := ""
	if params != nil {
		body = params
		idempotencyKey = params.IdempotencyKey
	}
	return m.api.postWithHeaders("webhook/test", body, commerceIdempotencyHeaders(idempotencyKey))
}
