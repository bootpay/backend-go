package bootpay

// EasyUserTokenPayload is the legacy request type for GetUserToken.
// Deprecated: Use UserTokenRequest with RequestUserToken instead.
type EasyUserTokenPayload struct {
	RestConfig
	UserId string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"username"`
	Gender int    `json:"gender"`
	Birth  string `json:"birth"`
	Phone  string `json:"phone"`
}

// GetUserToken requests a user token.
// Deprecated: Use RequestUserToken instead. This method unnecessarily sends
// application_id and private_key in the POST body via RestConfig embedding.
// GetUserToken is kept for backward compatibility and delegates to RequestUserToken.
func (api *Api) GetUserToken(userToken EasyUserTokenPayload) (APIResponse, error) {
	return api.RequestUserToken(UserTokenRequest{
		UserId:   userToken.UserId,
		Email:    userToken.Email,
		Username: userToken.Name,
		Gender:   userToken.Gender,
		Birth:    userToken.Birth,
		Phone:    userToken.Phone,
	})
}
