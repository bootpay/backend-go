package bootpay

// 현재 환경: "production" 또는 "development"
const CurrentEnv = "production"

// PG API 키
const (
	// Production 환경
	ProductionApplicationId = "5b8f6a4d396fa665fdc2b5ea"
	ProductionPrivateKey    = "rm6EYECr6aroQVG2ntW0A6LpWnkTgP4uQ3H18sDDUYw="

	// Development 환경
	DevApplicationId = "59bfc738e13f337dbd6ca48a"
	DevPrivateKey    = "pDc0NwlkEX3aSaHTp/PPL/i8vn5E/CqRChgyEp/gHD0="
)

// Commerce API 키
const (
	// Production 환경
	ProductionClientKey = "sEN72kYZBiyMNytA8nUGxQ"
	ProductionSecretKey = "rnZLJamENRgfwTccwmI_Uu9cxsPpAV9X2W-Htg73yfU="

	// Development 환경
	DevClientKey = "hxS-Up--5RvT6oU6QJE0JA"
	DevSecretKey = "r5zxvDcQJiAP2PBQ0aJjSHQtblNmYFt6uFoEMhti_mg="
)

// 테스트 데이터
const (
	TestReceiptId          = "628b2206d01c7e00209b6087"
	TestReceiptIdConfirm   = "62876963d01c7e00209b6028"
	TestReceiptIdCash      = "62e0f11f1fc192036b1b3c92"
	TestReceiptIdEscrow    = "628ae7ffd01c7e001e9b6066"
	TestReceiptIdBilling   = "62c7ccebcf9f6d001b3adcd4"
	TestReceiptIdTransfer  = "66541bc4ca4517e69343e24c"
	TestBillingKey         = "628b2644d01c7e00209b6092"
	TestBillingKey2        = "66542dfb4d18d5fc7b43e1b6"
	TestReserveId          = "6490149ca575b40024f0b70d"
	TestReserveId2         = "628b316cd01c7e00219b6081"
	TestUserId             = "1234"
	TestCertificateReceiptId = "61b009aaec81b4057e7f6ecd"
)

// PG API 키 가져오기
func GetPgKeys() (string, string) {
	if CurrentEnv == "production" {
		return ProductionApplicationId, ProductionPrivateKey
	}
	return DevApplicationId, DevPrivateKey
}

// Commerce API 키 가져오기
func GetCommerceKeys() (string, string) {
	if CurrentEnv == "production" {
		return ProductionClientKey, ProductionSecretKey
	}
	return DevClientKey, DevSecretKey
}

// PG API 인스턴스 생성
func CreatePgApi() *Api {
	appId, privateKey := GetPgKeys()
	mode := ""
	if CurrentEnv == "development" {
		mode = "development"
	}
	return Api{}.New(appId, privateKey, nil, mode)
}

// Commerce API 인스턴스 생성
func CreateCommerceApi() *CommerceApi {
	clientKey, secretKey := GetCommerceKeys()
	mode := ""
	if CurrentEnv == "development" {
		mode = "development"
	}
	return CommerceApi{}.New(clientKey, secretKey, nil, mode)
}
