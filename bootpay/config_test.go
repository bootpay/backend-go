package bootpay

import "os"

// getEnv returns the current environment from BOOTPAY_ENV env var.
// Defaults to "development" if not set.
func getEnv() string {
	env := os.Getenv("BOOTPAY_ENV")
	if env == "" {
		return "development"
	}
	return env
}

// PG API keys
const (
	// Production
	ProductionApplicationId = "5b8f6a4d396fa665fdc2b5ea"
	ProductionPrivateKey    = "rm6EYECr6aroQVG2ntW0A6LpWnkTgP4uQ3H18sDDUYw="

	// Development
	DevApplicationId = "59bfc738e13f337dbd6ca48a"
	DevPrivateKey    = "pDc0NwlkEX3aSaHTp/PPL/i8vn5E/CqRChgyEp/gHD0="
)

// Commerce API keys
const (
	// Production
	ProductionClientKey = "sEN72kYZBiyMNytA8nUGxQ"
	ProductionSecretKey = "rnZLJamENRgfwTccwmI_Uu9cxsPpAV9X2W-Htg73yfU="

	// Development
	DevClientKey = "hxS-Up--5RvT6oU6QJE0JA"
	DevSecretKey = "r5zxvDcQJiAP2PBQ0aJjSHQtblNmYFt6uFoEMhti_mg="
)

// Test data
const (
	TestReceiptId            = "628b2206d01c7e00209b6087"
	TestReceiptIdConfirm     = "62876963d01c7e00209b6028"
	TestReceiptIdCash        = "62e0f11f1fc192036b1b3c92"
	TestReceiptIdEscrow      = "628ae7ffd01c7e001e9b6066"
	TestReceiptIdBilling     = "62c7ccebcf9f6d001b3adcd4"
	TestReceiptIdTransfer    = "66541bc4ca4517e69343e24c"
	TestBillingKey           = "628b2644d01c7e00209b6092"
	TestBillingKey2          = "66542dfb4d18d5fc7b43e1b6"
	TestReserveId            = "6490149ca575b40024f0b70d"
	TestReserveId2           = "628b316cd01c7e00219b6081"
	TestUserId               = "1234"
	TestCertificateReceiptId = "61b009aaec81b4057e7f6ecd"
)

// GetPgKeys returns PG API keys based on environment
func GetPgKeys() (string, string) {
	if getEnv() == "production" {
		return ProductionApplicationId, ProductionPrivateKey
	}
	return DevApplicationId, DevPrivateKey
}

// GetCommerceKeys returns Commerce API keys based on environment
func GetCommerceKeys() (string, string) {
	if getEnv() == "production" {
		return ProductionClientKey, ProductionSecretKey
	}
	return DevClientKey, DevSecretKey
}

// CreatePgApi creates a PG API instance based on environment
func CreatePgApi() *Api {
	appId, privateKey := GetPgKeys()
	mode := ""
	if getEnv() == "development" {
		mode = "development"
	}
	return NewAPI(appId, privateKey, nil, mode)
}

// CreateCommerceApi creates a Commerce API instance based on environment
func CreateCommerceApi() *CommerceApi {
	clientKey, secretKey := GetCommerceKeys()
	mode := ""
	if getEnv() == "development" {
		mode = "development"
	}
	return NewCommerceAPI(clientKey, secretKey, nil, mode)
}
