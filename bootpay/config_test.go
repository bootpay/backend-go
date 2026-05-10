package bootpay

import (
	"fmt"
	"os"
	"strings"
)

func init() {
	loadDotEnv(".env")
	loadDotEnv("../.env")
}

func loadDotEnv(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	for _, raw := range strings.Split(string(data), "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}

func envOrDefault(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getEnv returns the current environment from BOOTPAY_ENV env var.
// Defaults to "production" if not set.
func getEnv() string {
	return envOrDefault("BOOTPAY_ENV", "production")
}

// getAuthMode returns the active PG auth mode: "new" (ck/sk) 또는 "legacy" (application_id/private_key).
// 매 실행 시 BOOTPAY_AUTH_MODE 환경변수로 토글한다.
func getAuthMode() string {
	mode := strings.ToLower(envOrDefault("BOOTPAY_AUTH_MODE", "new"))
	if mode == "" {
		return "new"
	}
	return mode
}

// PG legacy application_id/private_key (kept for legacy auth tests; ck/sk는 .env 로 주입)
const (
	// Production
	ProductionApplicationId = "5b8f6a4d396fa665fdc2b5ea"
	ProductionPrivateKey    = "rm6EYECr6aroQVG2ntW0A6LpWnkTgP4uQ3H18sDDUYw="

	// Development
	DevApplicationId = "59bfc738e13f337dbd6ca48a"
	DevPrivateKey    = "pDc0NwlkEX3aSaHTp/PPL/i8vn5E/CqRChgyEp/gHD0="
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
	TestCertificateReceiptId = "69fd7187564d1f550535538c"
)

// GetPgKeys returns PG API keys based on environment
func GetPgKeys() (string, string) {
	if getEnv() == "production" {
		return ProductionApplicationId, ProductionPrivateKey
	}
	return DevApplicationId, DevPrivateKey
}

// GetPgClientKeys returns recommended PG client_key/secret_key credentials based on environment.
// Keys must be provided via .env / environment variables (see .env.example).
func GetPgClientKeys() (string, string) {
	if getEnv() == "production" {
		return envOrDefault("BOOTPAY_PG_CLIENT_KEY_PROD", ""),
			envOrDefault("BOOTPAY_PG_SECRET_KEY_PROD", "")
	}
	return envOrDefault("BOOTPAY_PG_CLIENT_KEY_DEV", ""),
		envOrDefault("BOOTPAY_PG_SECRET_KEY_DEV", "")
}

// GetCommerceKeys returns Commerce API keys based on environment
// Keys must be provided via .env / environment variables (see .env.example).
func GetCommerceKeys() (string, string) {
	if getEnv() == "production" {
		return envOrDefault("BOOTPAY_COMMERCE_CLIENT_KEY_PROD", ""),
			envOrDefault("BOOTPAY_COMMERCE_SECRET_KEY_PROD", "")
	}
	return envOrDefault("BOOTPAY_COMMERCE_CLIENT_KEY_DEV", ""),
		envOrDefault("BOOTPAY_COMMERCE_SECRET_KEY_DEV", "")
}

// CreatePgApi creates a PG API instance based on environment.
// Defaults to client_key/secret_key. Toggle to legacy via BOOTPAY_AUTH_MODE=legacy.
func CreatePgApi() *Api {
	mode := ""
	if getEnv() == "development" {
		mode = "development"
	}
	envLabel := getEnv()
	if getAuthMode() == "legacy" {
		fmt.Printf("[BOOTPAY_AUTH_MODE=legacy] PG: application_id/private_key (Bearer) | env=%s\n", envLabel)
		appId, privateKey := GetPgKeys()
		return NewAPI(appId, privateKey, nil, mode)
	}
	fmt.Printf("[BOOTPAY_AUTH_MODE=new] PG: client_key/secret_key (Basic Auth) | env=%s\n", envLabel)
	clientKey, secretKey := GetPgClientKeys()
	return NewAPIWithClientKey(clientKey, secretKey, nil, mode)
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
