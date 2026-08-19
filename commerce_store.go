package bootpay

// StoreModule handles store-related operations
type StoreModule struct {
	api *CommerceApi
}

// Info retrieves store basic information (/v1/store)
// An Idempotency-Key header is auto-generated per call; use InfoWithIdempotencyKey to specify one.
func (m *StoreModule) Info() (map[string]interface{}, error) {
	return m.InfoWithIdempotencyKey("")
}

// InfoWithIdempotencyKey retrieves store basic information with an explicit Idempotency-Key
// (auto-generated when empty)
func (m *StoreModule) InfoWithIdempotencyKey(idempotencyKey string) (map[string]interface{}, error) {
	return m.api.getWithHeaders("store", commerceIdempotencyHeaders(idempotencyKey))
}

// Detail retrieves store detailed information (/v1/store/detail)
// An Idempotency-Key header is auto-generated per call; use DetailWithIdempotencyKey to specify one.
func (m *StoreModule) Detail() (map[string]interface{}, error) {
	return m.DetailWithIdempotencyKey("")
}

// DetailWithIdempotencyKey retrieves store detailed information with an explicit Idempotency-Key
// (auto-generated when empty)
func (m *StoreModule) DetailWithIdempotencyKey(idempotencyKey string) (map[string]interface{}, error) {
	return m.api.getWithHeaders("store/detail", commerceIdempotencyHeaders(idempotencyKey))
}
