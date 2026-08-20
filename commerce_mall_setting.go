package bootpay

// MallSettingModule handles mall setting operations (supervisor scope only)
type MallSettingModule struct {
	api *CommerceApi
}

// GetMallSetting retrieves the mall setting
// GET /v1/mall-setting (supervisor scope only)
// idempotencyKey is auto-generated when empty.
func (m *MallSettingModule) GetMallSetting(idempotencyKey string) (map[string]interface{}, error) {
	return m.api.getWithHeaders("mall-setting", commerceRoleHeaders("supervisor", idempotencyKey))
}

// Detail is an alias of GetMallSetting
func (m *MallSettingModule) Detail(idempotencyKey string) (map[string]interface{}, error) {
	return m.GetMallSetting(idempotencyKey)
}

// UpdateMallSetting updates the mall setting
// PUT /v1/mall-setting (supervisor scope only)
// The body is flat and only provided (non-nil / non-zero) values are sent.
// idempotencyKey is auto-generated when empty.
func (m *MallSettingModule) UpdateMallSetting(params MallSettingUpdateParams, idempotencyKey string) (map[string]interface{}, error) {
	return m.api.putWithHeaders("mall-setting", params, commerceRoleHeaders("supervisor", idempotencyKey))
}

// Update is an alias of UpdateMallSetting
func (m *MallSettingModule) Update(params MallSettingUpdateParams, idempotencyKey string) (map[string]interface{}, error) {
	return m.UpdateMallSetting(params, idempotencyKey)
}
