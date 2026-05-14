package bootpay

import "testing"

func TestCommerceUserToken(t *testing.T) {
	commerce := CreateCommerceApi()
	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.User.Token("test_user_123")
	if err != nil {
		t.Fatalf("Commerce User.Token failed: %v", err)
	}

	t.Logf("Commerce User.Token response: %+v", result)
}

func TestCommerceUserList(t *testing.T) {
	commerce := CreateCommerceApi()
	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.User.List(&UserListParams{
		ListParams: ListParams{
			Page:  1,
			Limit: 10,
		},
	})
	if err != nil {
		t.Fatalf("Commerce User.List failed: %v", err)
	}

	t.Logf("Commerce User.List response: %+v", result)
}

func TestCommerceUserDetail(t *testing.T) {
	t.Skip("Skipping: requires a valid user_id in the Commerce system")

	commerce := CreateCommerceApi()
	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.User.Detail("test_user_123")
	if err != nil {
		t.Fatalf("Commerce User.Detail failed: %v", err)
	}

	t.Logf("Commerce User.Detail response: %+v", result)
}

func TestCommerceUserCheckExist(t *testing.T) {
	commerce := CreateCommerceApi()
	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	result, err := commerce.User.CheckExist("login_id", "test@example.com")
	if err != nil {
		t.Fatalf("Commerce User.CheckExist failed: %v", err)
	}

	t.Logf("Commerce User.CheckExist response: %+v", result)
}

func TestCommerceUserJoinAndDelete(t *testing.T) {
	t.Skip("Skipping: creates and deletes a real user in the Commerce system")

	commerce := CreateCommerceApi()
	_, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("Commerce GetAccessToken failed: %v", err)
	}

	newUser := CommerceUser{
		LoginId: "integration_test_user@example.com",
		LoginPw: "test_password_123",
		Name:    "Integration Test User",
		Email:   "integration_test_user@example.com",
		Phone:   "010-0000-0000",
	}

	result, err := commerce.User.Join(newUser)
	if err != nil {
		t.Fatalf("Commerce User.Join failed: %v", err)
	}

	t.Logf("Commerce User.Join response: %+v", result)

	// Cleanup: delete the created user if we have a user_id
	if userId, ok := result["user_id"].(string); ok && userId != "" {
		deleteResult, err := commerce.User.Delete(userId)
		if err != nil {
			t.Logf("Commerce User.Delete cleanup failed: %v", err)
		} else {
			t.Logf("Commerce User.Delete cleanup response: %+v", deleteResult)
		}
	}
}
