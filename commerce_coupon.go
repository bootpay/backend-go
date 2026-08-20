package bootpay

import (
	"net/url"
	"strconv"
)

// CouponModule handles coupon-related operations
type CouponModule struct {
	api *CommerceApi
}

// List retrieves user-owned coupon list
func (m *CouponModule) List(params *CouponListParams) (map[string]interface{}, error) {
	query := ""
	if params != nil {
		queryParams := url.Values{}
		if params.Status != "" {
			queryParams.Set("status", params.Status)
		}
		if params.Page > 0 {
			queryParams.Set("page", strconv.Itoa(params.Page))
		}
		if params.Limit > 0 {
			queryParams.Set("limit", strconv.Itoa(params.Limit))
		}
		if len(queryParams) > 0 {
			query = "?" + queryParams.Encode()
		}
	}
	return m.api.Get("coupon" + query)
}

// Available retrieves coupons available for download
func (m *CouponModule) Available() (map[string]interface{}, error) {
	return m.api.Get("coupon/available")
}

// Download downloads a coupon from a template
func (m *CouponModule) Download(params CouponDownloadParams) (map[string]interface{}, error) {
	return m.api.Post("coupon/download", params)
}
