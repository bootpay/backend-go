package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)
type Shipping struct {
	//RestConfig
	ReceiptId      string `json:"receipt_id"`
	DeliveryCorp   string `json:"delivery_corp"`
	TrackingNumber string `json:"tracking_number"`
}

type ShippingUser struct {
	Username          string 					`json:"username"`
	Phone          	  string 					`json:"phone"`
	Zipcode           string 					`json:"zipcode"`
	Address           string 					`json:"address"`
}


type ShippingCompany struct {
	Name          	  string 					`json:"name"`
	Phone          	  string 					`json:"phone"`
	Zipcode           string 					`json:"zipcode"`
	Addr1             string 					`json:"addr1"`
	Addr2             string 					`json:"addr2"`
}


func (api *Api) ShippingStart(payload Shipping) (APIResponse, error) {
	if payload.ApplicationId == "" {
		payload.ApplicationId = api.applicationId
	}
	if payload.PrivateKey == "" {
		payload.PrivateKey = api.privateKey
	}
	postBody, _ := json.Marshal(payload)
	body := bytes.NewBuffer(postBody)
	req, err := api.NewRequest(http.MethodPut, "/subscribe/payment", body)
	if err != nil {
		errors.New("bootpay: ShippingStart error: " + err.Error())
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)
	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}