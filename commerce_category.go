package bootpay

import "fmt"

// CategoryModule handles category-related operations
type CategoryModule struct {
	api *CommerceApi
}

// List retrieves category tree
func (m *CategoryModule) List() (map[string]interface{}, error) {
	return m.api.Get("categories")
}

// Detail retrieves a single category
func (m *CategoryModule) Detail(categoryId string) (map[string]interface{}, error) {
	return m.api.Get(fmt.Sprintf("categories/%s", categoryId))
}

// Create creates a new category
// ⚠️ The server requires supervisor scope (scope_invalid!) — BOOTPAY-ROLE must be sent.
func (m *CategoryModule) Create(params CategoryCreateParams) (map[string]interface{}, error) {
	return m.api.postWithHeaders("categories", params,
		commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// Update updates a category
// ⚠️ The server requires supervisor scope (scope_invalid!) — BOOTPAY-ROLE must be sent.
func (m *CategoryModule) Update(params CategoryUpdateParams) (map[string]interface{}, error) {
	if params.CategoryId == "" {
		return nil, fmt.Errorf("category_id is required")
	}
	categoryId := params.CategoryId
	// Send body without category_id (kept in URL)
	body := CategoryUpdateBody{
		Name:             params.Name,
		ParentCategoryId: params.ParentCategoryId,
		StatusDisplay:    params.StatusDisplay,
		StatusBest:       params.StatusBest,
		FilterColor:      params.FilterColor,
		FilterSize:       params.FilterSize,
	}
	return m.api.putWithHeaders(fmt.Sprintf("categories/%s", categoryId), body,
		commerceRoleHeaders("supervisor", params.IdempotencyKey))
}

// Destroy deletes a category
// ⚠️ The server requires supervisor scope (scope_invalid!) — BOOTPAY-ROLE must be sent.
// idempotencyKey is optional; when omitted a key is generated per call.
func (m *CategoryModule) Destroy(categoryId string, idempotencyKey ...string) (map[string]interface{}, error) {
	return m.api.deleteWithHeaders(fmt.Sprintf("categories/%s", categoryId), nil,
		commerceRoleHeaders("supervisor", firstOrEmpty(idempotencyKey)))
}
