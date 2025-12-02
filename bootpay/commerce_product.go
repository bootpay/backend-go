package bootpay

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

// ProductModule handles product-related operations
type ProductModule struct {
	api *CommerceApi
}

// List retrieves product list
func (m *ProductModule) List(params *ProductListParams) (map[string]interface{}, error) {
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
		if params.Type > 0 {
			queryParams.Set("type", strconv.Itoa(params.Type))
		}
		if params.PeriodType != "" {
			queryParams.Set("period_type", params.PeriodType)
		}
		if params.SAt != "" {
			queryParams.Set("s_at", params.SAt)
		}
		if params.EAt != "" {
			queryParams.Set("e_at", params.EAt)
		}
		if params.CategoryCode != "" {
			queryParams.Set("category_code", params.CategoryCode)
		}
		if len(queryParams) > 0 {
			query = "?" + queryParams.Encode()
		}
	}
	return m.api.Get("products" + query)
}

// Create creates a new product with optional images
func (m *ProductModule) Create(product CommerceProduct, imagePaths []string) (map[string]interface{}, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Write product data as individual fields
	if product.Name != "" {
		writer.WriteField("name", product.Name)
	}
	if product.Description != "" {
		writer.WriteField("description", product.Description)
	}
	if product.Type > 0 {
		writer.WriteField("type", strconv.Itoa(product.Type))
	}
	if product.DisplayPrice > 0 {
		writer.WriteField("display_price", strconv.Itoa(product.DisplayPrice))
	}
	if product.TaxFreePrice > 0 {
		writer.WriteField("tax_free_price", strconv.Itoa(product.TaxFreePrice))
	}
	// Add more fields as needed...

	// For complex objects, marshal as JSON
	if product.Options != nil {
		optionsJson, _ := json.Marshal(product.Options)
		writer.WriteField("options", string(optionsJson))
	}
	if product.SubscriptionSetting != nil {
		settingJson, _ := json.Marshal(product.SubscriptionSetting)
		writer.WriteField("subscription_setting", string(settingJson))
	}

	// Add image files
	if imagePaths != nil {
		for _, imagePath := range imagePaths {
			file, err := os.Open(imagePath)
			if err != nil {
				return nil, err
			}
			defer file.Close()

			part, err := writer.CreateFormFile("images", filepath.Base(imagePath))
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(part, file)
			if err != nil {
				return nil, err
			}
		}
	}

	writer.Close()

	// Make the request
	req, err := m.api.newRequest("POST", "products", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := m.api.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	result := make(map[string]interface{})
	json.NewDecoder(res.Body).Decode(&result)

	// If no images, use simple JSON post
	if imagePaths == nil || len(imagePaths) == 0 {
		return m.api.Post("products", product)
	}

	return result, nil
}

// CreateSimple creates a new product without images (JSON only)
func (m *ProductModule) CreateSimple(product CommerceProduct) (map[string]interface{}, error) {
	return m.api.Post("products", product)
}

// Detail retrieves product details
func (m *ProductModule) Detail(productId string) (map[string]interface{}, error) {
	return m.api.Get(fmt.Sprintf("products/%s", productId))
}

// Update updates product information
func (m *ProductModule) Update(product CommerceProduct) (map[string]interface{}, error) {
	if product.ProductId == "" {
		return nil, fmt.Errorf("product_id is required")
	}
	return m.api.Put(fmt.Sprintf("products/%s", product.ProductId), product)
}

// Status changes product status
func (m *ProductModule) Status(params ProductStatusParams) (map[string]interface{}, error) {
	if params.ProductId == "" {
		return nil, fmt.Errorf("product_id is required")
	}
	return m.api.Put(fmt.Sprintf("products/%s/status", params.ProductId), params)
}

// Delete deletes a product
func (m *ProductModule) Delete(productId string) (map[string]interface{}, error) {
	return m.api.Delete(fmt.Sprintf("products/%s", productId))
}

// Ensure unused import doesn't cause issues
var _ = string([]byte{})
