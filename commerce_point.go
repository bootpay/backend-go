package bootpay

import (
	"net/url"
	"strconv"
)

// PointModule handles point-related operations
type PointModule struct {
	api *CommerceApi
}

// Balance retrieves point balance
func (m *PointModule) Balance() (map[string]interface{}, error) {
	return m.api.Get("point/balance")
}

// Transactions retrieves point transaction history
func (m *PointModule) Transactions(params *PointTransactionsParams) (map[string]interface{}, error) {
	query := ""
	if params != nil {
		queryParams := url.Values{}
		if params.Page > 0 {
			queryParams.Set("page", strconv.Itoa(params.Page))
		}
		if params.Limit > 0 {
			queryParams.Set("limit", strconv.Itoa(params.Limit))
		}
		if params.TransactionType != 0 {
			queryParams.Set("transaction_type", strconv.Itoa(params.TransactionType))
		}
		if len(queryParams) > 0 {
			query = "?" + queryParams.Encode()
		}
	}
	return m.api.Get("point/transactions" + query)
}
