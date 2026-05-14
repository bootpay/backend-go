package bootpay

// StoreModule handles store-related operations
type StoreModule struct {
	api *CommerceApi
}

// Info retrieves store basic information
func (m *StoreModule) Info() (map[string]interface{}, error) {
	return m.api.Get("store")
}

// Detail retrieves store detailed information
func (m *StoreModule) Detail() (map[string]interface{}, error) {
	return m.api.Get("store/detail")
}
