package bootpay

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"time"
)

const (
	DEVELOPMENT string = "https://dev-api.bootpay.co.kr/v2"
	TEST        string = "https://test-api.bootpay.co.kr/v2"
	STAGE       string = "https://stage-api.bootpay.co.kr/v2"
	PRODUCTION  string = "https://api.bootpay.co.kr/v2"
)
const defaultHTTPTimeout = 10 * time.Second

type APIResponse = map[string]interface{}

//type APIResponse struct {
//	Data map[string]interface{}
//	//Status  int    `json:"status"`
//	//Code    int    `json:"code"`
//	//Message string `json:"message"`
//	//Data    map[string]interface{} `json:"data"`
//}

type RestConfig struct {
	ApplicationId string `json:"application_id"`
	PrivateKey    string `json:"private_key"`
}

type Api struct {
	token         string
	applicationId string
	privateKey    string
	baseUrl       string
	client        *http.Client
}

func (api Api) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, api.baseUrl+url, body)
	if err != nil {
		errors.New("Cannot create Bootpay request: " + err.Error())
		return nil, err
	}
	if api.token != "" {
		req.Header.Set("Authorization", "Bearer " + api.token)
	} else if api.applicationId != "" && api.privateKey != "" {
		credentials := api.applicationId + ":" + api.privateKey
		encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
		req.Header.Set("Authorization", "Basic "+encoded)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	return req, nil
}

func (api Api) New(applicationId string, privateKey string, client *http.Client, mode string) *Api {
	if client == nil {
		client = &http.Client{
			Timeout: defaultHTTPTimeout,
			Transport: &http.Transport{
				TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
			},
		}
	}
	baseUrl := PRODUCTION
	if mode == "development" {
		baseUrl = DEVELOPMENT
	} else if mode == "test" {
		baseUrl = TEST
	} else if mode == "stage" {
		baseUrl = STAGE
	}
	return &Api{
		applicationId: applicationId,
		privateKey:    privateKey,
		baseUrl:       baseUrl,
		client:    	   client,
	}
}
