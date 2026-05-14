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
func (m *CategoryModule) Create(params CategoryCreateParams) (map[string]interface{}, error) {
	return m.api.Post("categories", params)
}

// Update updates a category
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
	return m.api.Put(fmt.Sprintf("categories/%s", categoryId), body)
}

// Destroy deletes a category
func (m *CategoryModule) Destroy(categoryId string) (map[string]interface{}, error) {
	return m.api.Delete(fmt.Sprintf("categories/%s", categoryId))
}
