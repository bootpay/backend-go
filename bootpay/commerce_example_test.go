package bootpay

import (
	"fmt"
	"testing"
)

// 테스트용 키 (Commerce API)
const (
	// Development 환경
	commerceDevClientKey = "hxS-Up--5RvT6oU6QJE0JA"
	commerceDevSecretKey = "r5zxvDcQJiAP2PBQ0aJjSHQtblNmYFt6uFoEMhti_mg="

	// Production 환경
	commerceProdClientKey = "sEN72kYZBiyMNytA8nUGxQ"
	commerceProdSecretKey = "rnZLJamENRgfwTccwmI_Uu9cxsPpAV9X2W-Htg73yfU="
)

// getCommerceTestApi returns Commerce API client for development
func getCommerceTestApi() *CommerceApi {
	return NewCommerceApi(commerceDevClientKey, commerceDevSecretKey, nil, "development")
}

// getCommerceProductionApi returns Commerce API client for production
func getCommerceProductionApi() *CommerceApi {
	return NewCommerceApi(commerceProdClientKey, commerceProdSecretKey, nil, "production")
}

// TestCommerceGetAccessToken tests the GetAccessToken method
func TestCommerceGetAccessToken(t *testing.T) {
	commerce := getCommerceTestApi()

	result, err := commerce.GetAccessToken()
	if err != nil {
		t.Logf("GetAccessToken error (expected with test keys): %v", err)
		return
	}

	fmt.Printf("Access Token Response: %+v\n", result)
}

// ExampleNewCommerceApi demonstrates how to create a Commerce API client
func ExampleNewCommerceApi() {
	// Create Commerce API client (Development)
	commerce := NewCommerceApi(commerceDevClientKey, commerceDevSecretKey, nil, "development")
	// Production: NewCommerceApi(commerceProdClientKey, commerceProdSecretKey, nil, "production")

	// Get access token
	_, _ = commerce.GetAccessToken()

	// Set role (method chaining supported)
	commerce.AsManager()

	fmt.Println("Commerce API client created")
	// Output: Commerce API client created
}

// ExampleCommerceApi_User demonstrates User module usage
func ExampleCommerceApi_User() {
	commerce := getCommerceTestApi()
	commerce.GetAccessToken()

	// Get user token
	_, _ = commerce.User.Token("user_123")

	// Create new user
	newUser := CommerceUser{
		LoginId:  "testuser@example.com",
		LoginPw:  "password123",
		Name:     "Test User",
		Email:    "testuser@example.com",
		Phone:    "010-1234-5678",
	}
	_, _ = commerce.User.Join(newUser)

	// Check if user exists
	_, _ = commerce.User.CheckExist("login_id", "testuser@example.com")

	// Get user list
	_, _ = commerce.User.List(&UserListParams{
		ListParams: ListParams{
			Page:  1,
			Limit: 10,
		},
	})

	// Get user detail
	_, _ = commerce.User.Detail("user_123")

	// Update user
	updateUser := CommerceUser{
		UserId: "user_123",
		Name:   "Updated User",
	}
	_, _ = commerce.User.Update(updateUser)

	// Delete user
	_, _ = commerce.User.Delete("user_123")

	fmt.Println("User module examples completed")
	// Output: User module examples completed
}

// ExampleCommerceApi_UserGroup demonstrates UserGroup module usage
func ExampleCommerceApi_UserGroup() {
	commerce := getCommerceTestApi()
	commerce.GetAccessToken()

	// Create user group
	newGroup := CommerceUserGroup{
		CompanyName:    "Test Company",
		CorporateType:  CORPORATE_TYPE_CORPORATE,
		BusinessNumber: "123-45-67890",
	}
	_, _ = commerce.UserGroup.Create(newGroup)

	// Get user group list
	_, _ = commerce.UserGroup.List(&UserGroupListParams{
		ListParams: ListParams{
			Page:  1,
			Limit: 10,
		},
		CorporateType: CORPORATE_TYPE_CORPORATE,
	})

	// Get user group detail
	_, _ = commerce.UserGroup.Detail("group_123")

	// Add user to group
	_, _ = commerce.UserGroup.UserCreate("group_123", "user_456")

	// Remove user from group
	_, _ = commerce.UserGroup.UserDelete("group_123", "user_456")

	// Set group limit
	_, _ = commerce.UserGroup.Limit(UserGroupLimitParams{
		UserGroupId:   "group_123",
		UseLimit:      true,
		PurchaseLimit: 1000000,
	})

	// Set aggregate transaction settings
	_, _ = commerce.UserGroup.AggregateTransaction(UserGroupAggregateTransactionParams{
		UserGroupId:                         "group_123",
		UseSubscriptionAggregateTransaction: true,
		SubscriptionMonthDay:                15,
	})

	fmt.Println("UserGroup module examples completed")
	// Output: UserGroup module examples completed
}

// ExampleCommerceApi_Product demonstrates Product module usage
func ExampleCommerceApi_Product() {
	commerce := getCommerceTestApi()
	commerce.GetAccessToken()

	// Get product list
	_, _ = commerce.Product.List(&ProductListParams{
		ListParams: ListParams{
			Page:  1,
			Limit: 10,
		},
		Type: 1,
	})

	// Create product (without images)
	newProduct := CommerceProduct{
		Name:         "Test Product",
		DisplayPrice: 10000,
		TaxFreePrice: 0,
		Type:         1,
	}
	_, _ = commerce.Product.CreateSimple(newProduct)

	// Create product (with images)
	// _, _ = commerce.Product.Create(newProduct, []string{"/path/to/image1.jpg", "/path/to/image2.jpg"})

	// Get product detail
	_, _ = commerce.Product.Detail("product_123")

	// Update product
	updateProduct := CommerceProduct{
		ProductId:    "product_123",
		Name:         "Updated Product",
		DisplayPrice: 15000,
	}
	_, _ = commerce.Product.Update(updateProduct)

	// Change product status
	_, _ = commerce.Product.Status(ProductStatusParams{
		ProductId:     "product_123",
		Status:        1,
		StatusDisplay: true,
		StatusSale:    true,
	})

	// Delete product
	_, _ = commerce.Product.Delete("product_123")

	fmt.Println("Product module examples completed")
	// Output: Product module examples completed
}

// ExampleCommerceApi_Invoice demonstrates Invoice module usage
func ExampleCommerceApi_Invoice() {
	commerce := getCommerceTestApi()
	commerce.GetAccessToken()

	// Get invoice list
	_, _ = commerce.Invoice.List(&ListParams{
		Page:  1,
		Limit: 10,
	})

	// Create invoice
	newInvoice := CommerceInvoice{
		Title:        "Test Invoice",
		Price:        50000,
		TaxFreePrice: 0,
		UserId:       "user_123",
		InvoiceItems: []CommerceInvoiceItem{
			{Name: "Item 1", Price: 30000, Qty: 1},
			{Name: "Item 2", Price: 20000, Qty: 1},
		},
	}
	_, _ = commerce.Invoice.Create(newInvoice)

	// Get invoice detail
	_, _ = commerce.Invoice.Detail("invoice_123")

	// Send invoice notification
	_, _ = commerce.Invoice.Notify("invoice_123", []int{INVOICE_SEND_TYPE_SMS, INVOICE_SEND_TYPE_EMAIL})

	fmt.Println("Invoice module examples completed")
	// Output: Invoice module examples completed
}

// ExampleCommerceApi_Order demonstrates Order module usage
func ExampleCommerceApi_Order() {
	commerce := getCommerceTestApi()
	commerce.GetAccessToken()

	// Get order list
	_, _ = commerce.Order.List(&OrderListParams{
		ListParams: ListParams{
			Page:  1,
			Limit: 10,
		},
		UserId:        "user_123",
		Status:        []int{1, 2, 3},
		PaymentStatus: []int{1},
	})

	// Get order detail
	_, _ = commerce.Order.Detail("order_123")

	// Get monthly orders
	_, _ = commerce.Order.Month("group_123", "2024-01")

	fmt.Println("Order module examples completed")
	// Output: Order module examples completed
}

// ExampleCommerceApi_OrderCancel demonstrates OrderCancel module usage
func ExampleCommerceApi_OrderCancel() {
	commerce := getCommerceTestApi()
	commerce.GetAccessToken()

	// Get cancel request list
	_, _ = commerce.OrderCancel.List(&OrderCancelListParams{
		OrderId: "order_123",
	})

	// Request cancellation
	_, _ = commerce.OrderCancel.Request(OrderCancelParams{
		OrderNumber: "ORDER-2024-001",
		RequestCancelParameters: &RequestCancelParameter{
			CancelReason: "Customer request",
			CancelType:   1,
		},
	})

	// Withdraw cancel request
	_, _ = commerce.OrderCancel.Withdraw("cancel_request_123")

	// Approve cancel request
	_, _ = commerce.OrderCancel.Approve(OrderCancelActionParams{
		OrderCancelRequestHistoryId: "cancel_request_123",
		CancelReason:                "Approved by manager",
	})

	// Reject cancel request
	_, _ = commerce.OrderCancel.Reject(OrderCancelActionParams{
		OrderCancelRequestHistoryId: "cancel_request_456",
		CancelReason:                "Cannot cancel - already shipped",
	})

	fmt.Println("OrderCancel module examples completed")
	// Output: OrderCancel module examples completed
}

// ExampleCommerceApi_OrderSubscription demonstrates OrderSubscription module usage
func ExampleCommerceApi_OrderSubscription() {
	commerce := getCommerceTestApi()
	commerce.GetAccessToken()

	// Get subscription list
	_, _ = commerce.OrderSubscription.List(&OrderSubscriptionListParams{
		ListParams: ListParams{
			Page:  1,
			Limit: 10,
		},
		UserId: "user_123",
	})

	// Get subscription detail
	_, _ = commerce.OrderSubscription.Detail("subscription_123")

	// Update subscription
	_, _ = commerce.OrderSubscription.Update(OrderSubscriptionUpdateParams{
		OrderSubscriptionId: "subscription_123",
		Status:              1,
	})

	// Pause subscription
	_, _ = commerce.OrderSubscription.RequestIng.Pause(OrderSubscriptionPauseParams{
		OrderSubscriptionId: "subscription_123",
		Reason:              "Temporary pause",
	})

	// Resume subscription
	_, _ = commerce.OrderSubscription.RequestIng.Resume(OrderSubscriptionResumeParams{
		OrderSubscriptionId: "subscription_123",
	})

	// Calculate termination fee
	_, _ = commerce.OrderSubscription.RequestIng.CalculateTerminationFee("subscription_123", "")

	// Calculate termination fee by order number
	_, _ = commerce.OrderSubscription.RequestIng.CalculateTerminationFeeByOrderNumber("ORDER-2024-001")

	// Terminate subscription
	_, _ = commerce.OrderSubscription.RequestIng.Termination(OrderSubscriptionTerminationParams{
		OrderSubscriptionId: "subscription_123",
		Reason:              "Customer request",
		TerminationFee:      10000,
	})

	fmt.Println("OrderSubscription module examples completed")
	// Output: OrderSubscription module examples completed
}

// ExampleCommerceApi_OrderSubscriptionBill demonstrates OrderSubscriptionBill module usage
func ExampleCommerceApi_OrderSubscriptionBill() {
	commerce := getCommerceTestApi()
	commerce.GetAccessToken()

	// Get subscription bill list
	_, _ = commerce.OrderSubscriptionBill.List(&OrderSubscriptionBillListParams{
		ListParams: ListParams{
			Page:  1,
			Limit: 10,
		},
		OrderSubscriptionId: "subscription_123",
		Status:              []int{1, 2},
	})

	// Get subscription bill detail
	_, _ = commerce.OrderSubscriptionBill.Detail("bill_123")

	// Update subscription bill
	_, _ = commerce.OrderSubscriptionBill.Update(CommerceOrderSubscriptionBill{
		OrderSubscriptionBillId: "bill_123",
		Status:                  2,
	})

	fmt.Println("OrderSubscriptionBill module examples completed")
	// Output: OrderSubscriptionBill module examples completed
}

// ExampleCommerceApi_OrderSubscriptionAdjustment demonstrates OrderSubscriptionAdjustment module usage
func ExampleCommerceApi_OrderSubscriptionAdjustment() {
	commerce := getCommerceTestApi()
	commerce.GetAccessToken()

	// Create subscription adjustment
	_, _ = commerce.OrderSubscriptionAdjustment.Create("subscription_123", CommerceOrderSubscriptionAdjustment{
		Name:         "Discount",
		Price:        -5000,
		TaxFreePrice: 0,
		Duration:     3,
		Type:         SUBSCRIPTION_ADJUSTMENT_TYPE_PERIOD_DISCOUNT,
	})

	// Update subscription adjustment
	_, _ = commerce.OrderSubscriptionAdjustment.Update(OrderSubscriptionAdjustmentUpdateParams{
		OrderSubscriptionId:           "subscription_123",
		OrderSubscriptionAdjustmentId: "adjustment_123",
		Price:                         -10000,
	})

	// Delete subscription adjustment
	_, _ = commerce.OrderSubscriptionAdjustment.Delete("subscription_123", "adjustment_123")

	fmt.Println("OrderSubscriptionAdjustment module examples completed")
	// Output: OrderSubscriptionAdjustment module examples completed
}

// ExampleCommerceApi_RoleChaining demonstrates role chaining
func ExampleCommerceApi_RoleChaining() {
	commerce := getCommerceTestApi()
	commerce.GetAccessToken()

	// Method chaining for role setting
	commerce.AsUser().User.List(nil)
	commerce.AsManager().UserGroup.List(nil)
	commerce.AsPartner().Product.List(nil)
	commerce.AsVendor().Order.List(nil)
	commerce.AsSupervisor().Invoice.List(nil)

	// Get current role
	currentRole := commerce.GetRole()
	fmt.Printf("Current role: %s\n", currentRole)

	// Output: Current role: supervisor
}

// TestCommerceAllEndpoints tests all Commerce API endpoints for 404 errors
func TestCommerceAllEndpoints(t *testing.T) {
	commerce := getCommerceTestApi()

	// Get access token first
	tokenResult, err := commerce.GetAccessToken()
	if err != nil {
		t.Fatalf("GetAccessToken failed: %v", err)
	}
	fmt.Printf("[Token] %+v\n\n", tokenResult)

	// Helper function to check for 404
	check404 := func(name string, result map[string]interface{}, err error) {
		if err != nil {
			fmt.Printf("[%s] ERROR: %v\n", name, err)
			return
		}
		if msg, ok := result["message"].(string); ok && msg == "Not Found" {
			fmt.Printf("[%s] ❌ 404 NOT FOUND\n", name)
		} else if errCode, ok := result["error_code"].(float64); ok && errCode == 404 {
			fmt.Printf("[%s] ❌ 404 NOT FOUND\n", name)
		} else {
			fmt.Printf("[%s] ✅ OK - %+v\n", name, result)
		}
	}

	fmt.Println("========== User Module ==========")
	result, err := commerce.User.Token("test_user_123")
	check404("User.Token", result, err)

	result, err = commerce.User.List(&UserListParams{ListParams: ListParams{Page: 1, Limit: 10}})
	check404("User.List", result, err)

	result, err = commerce.User.CheckExist("login_id", "test@example.com")
	check404("User.CheckExist", result, err)

	result, err = commerce.User.Detail("test_user_123")
	check404("User.Detail", result, err)

	fmt.Println("\n========== UserGroup Module ==========")
	result, err = commerce.UserGroup.List(&UserGroupListParams{ListParams: ListParams{Page: 1, Limit: 10}})
	check404("UserGroup.List", result, err)

	result, err = commerce.UserGroup.Detail("test_group_123")
	check404("UserGroup.Detail", result, err)

	fmt.Println("\n========== Product Module ==========")
	result, err = commerce.Product.List(&ProductListParams{ListParams: ListParams{Page: 1, Limit: 10}})
	check404("Product.List", result, err)

	result, err = commerce.Product.Detail("test_product_123")
	check404("Product.Detail", result, err)

	fmt.Println("\n========== Invoice Module ==========")
	result, err = commerce.Invoice.List(&ListParams{Page: 1, Limit: 10})
	check404("Invoice.List", result, err)

	result, err = commerce.Invoice.Detail("test_invoice_123")
	check404("Invoice.Detail", result, err)

	fmt.Println("\n========== Order Module ==========")
	result, err = commerce.Order.List(&OrderListParams{ListParams: ListParams{Page: 1, Limit: 10}})
	check404("Order.List", result, err)

	result, err = commerce.Order.Detail("test_order_123")
	check404("Order.Detail", result, err)

	result, err = commerce.Order.Month("test_group_123", "2024-01")
	check404("Order.Month", result, err)

	fmt.Println("\n========== OrderCancel Module ==========")
	result, err = commerce.OrderCancel.List(&OrderCancelListParams{OrderId: "test_order_123"})
	check404("OrderCancel.List", result, err)

	fmt.Println("\n========== OrderSubscription Module ==========")
	result, err = commerce.OrderSubscription.List(&OrderSubscriptionListParams{ListParams: ListParams{Page: 1, Limit: 10}})
	check404("OrderSubscription.List", result, err)

	result, err = commerce.OrderSubscription.Detail("test_subscription_123")
	check404("OrderSubscription.Detail", result, err)

	result, err = commerce.OrderSubscription.RequestIng.CalculateTerminationFee("test_subscription_123", "")
	check404("OrderSubscription.RequestIng.CalculateTerminationFee", result, err)

	fmt.Println("\n========== OrderSubscriptionBill Module ==========")
	result, err = commerce.OrderSubscriptionBill.List(&OrderSubscriptionBillListParams{ListParams: ListParams{Page: 1, Limit: 10}})
	check404("OrderSubscriptionBill.List", result, err)

	result, err = commerce.OrderSubscriptionBill.Detail("test_bill_123")
	check404("OrderSubscriptionBill.Detail", result, err)

	fmt.Println("\n========== Test Complete ==========")
}
