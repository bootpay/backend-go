package bootpay

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)


type Shipping struct {
	ReceiptId string `json:"receipt_id"`
	DeliveryCorp string `json:"delivery_corp"`
	TrackingNumber string `json:"tracking_number"`
	//Data VerifyData `json:"data"`
	ShippingPrepayment bool `json:"shipping_prepayment"`
	ShippingDay bool `json:"shipping_day"`
	User ShippingUser `json:"user"`
	Company ShippingCompany `json:"company"`
}


type ShippingUser struct {
	Username string `json:"username"`
	Phone string `json:"phone"`
	Zipcode string `json:"zipcode"`
	Address string `json:"address"`
}


type ShippingCompany struct {
	Name string `json:"name"`
	Phone string `json:"phone"`
	Zipcode string `json:"zipcode"`
	Addr1 string `json:"addr1"`
	Addr2 string `json:"addr2"`
}


func (api *Api) putShippingStart(shipping Shipping) (APIResponse, error) {
	putBody, _ := json.Marshal(shipping)
	body := bytes.NewBuffer(putBody)

	req, err := api.NewRequest(http.MethodPut, "/escrow/shipping/start/" + shipping.ReceiptId, body)
	if err != nil {
		errors.New("bootpay: putShippingStart error: " + err.Error())
		return APIResponse{}, err
	}
	res, err := api.client.Do(req)

	defer res.Body.Close()

	result := APIResponse{}
	json.NewDecoder(res.Body).Decode(&result)
	return result, nil
}
 