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

	API_VERSION string = "5.1.0"
	SDK_VERSION string = "2.7.0"
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
	ApplicationId string `json:"application_id,omitempty"`
	PrivateKey    string `json:"private_key,omitempty"`
	ClientKey     string `json:"client_key,omitempty"`
	SecretKey     string `json:"secret_key,omitempty"`
}

type Api struct {
	token         string
	applicationId string
	privateKey    string
	clientKey     string
	secretKey     string
	baseUrl       string
	client        *http.Client
}

func (api Api) validateCredentials() error {
	if api.clientKey != "" || api.secretKey != "" {
		if api.clientKey == "" || api.secretKey == "" {
			return errors.New("bootpay: client_key and secret_key must be provided together")
		}
		return nil
	}
	if api.applicationId == "" || api.privateKey == "" {
		return errors.New("bootpay: application_id and private_key must be provided together")
	}
	return nil
}

func (api Api) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	if err := api.validateCredentials(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, api.baseUrl+url, body)
	if err != nil {
		return nil, err
	}
	if api.clientKey != "" && api.secretKey != "" {
		credentials := base64.StdEncoding.EncodeToString([]byte(api.clientKey + ":" + api.secretKey))
		req.Header.Set("Authorization", "Basic "+credentials)
	} else if api.token != "" {
		req.Header.Set("Authorization", "Bearer "+api.token)
	}
	req.Header.Set("BOOTPAY-API-VERSION", API_VERSION)
	req.Header.Set("BOOTPAY-SDK-VERSION", SDK_VERSION)
	req.Header.Set("BOOTPAY-SDK-TYPE", "305")

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Charset", "utf-8")
	return req, nil
}

func newAPIBase(client *http.Client, mode string) (string, *http.Client) {
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
	return baseUrl, client
}

// NewAPI creates a new PG API instance using legacy application_id/private_key.
// Deprecated credentials remain supported for backward compatibility.
func NewAPI(applicationId string, privateKey string, client *http.Client, mode string) *Api {
	baseUrl, httpClient := newAPIBase(client, mode)
	return &Api{
		applicationId: applicationId,
		privateKey:    privateKey,
		baseUrl:       baseUrl,
		client:        httpClient,
	}
}

// NewAPIWithClientKey creates a PG API instance using client_key/secret_key.
// When client_key is present this Basic Auth flow takes precedence over legacy token auth.
func NewAPIWithClientKey(clientKey string, secretKey string, client *http.Client, mode string) *Api {
	baseUrl, httpClient := newAPIBase(client, mode)
	return &Api{
		clientKey: clientKey,
		secretKey: secretKey,
		baseUrl:   baseUrl,
		client:    httpClient,
	}
}

// New creates a new PG API instance (deprecated: use NewAPI instead)
func (api Api) New(applicationId string, privateKey string, client *http.Client, mode string) *Api {
	return NewAPI(applicationId, privateKey, client, mode)
}
