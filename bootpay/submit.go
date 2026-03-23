package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// Submit represents the server submit request payload.
// RestConfig embedding was removed because application_id and private_key
// should not be sent in the POST body; authentication is handled via the
// Authorization header set by NewRequest.
type Submit struct {
	ReceiptId string `json:"receipt_id"`
}

func (api *Api) ServerSubmit(receiptId string) (APIResponse, error) {
	sub := Submit{
		ReceiptId: receiptId,
	}
	postBody, _ := json.Marshal(sub)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPost, "/submit", body)
	if err != nil {
		errors.New("bootpay: Submit error: " + err.Error())
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)

	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	if result == nil { result =  map[string]interface{}{} }
	result["http_status"] = res.StatusCode
	return result, nil
}
