package validators

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreateTransactionsEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedResult bool
		expectedStatus int
	}{
		{
			name: "Valid request",
			requestBody: map[string]interface{}{
				"amount":        100.50,
				"title":         "Test Entry",
				"currency":      "USD",
				"type":          "income",
				"category":      "Salary",
				"paymentMethod": "Credit Card",
				"description":   "Monthly salary",
				"date":          "01/01/2023",
			},
			expectedResult: true,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid amount",
			requestBody: map[string]interface{}{
				"amount":        0,
				"title":         "Test Entry",
				"currency":      "USD",
				"type":          "income",
				"category":      "Salary",
				"paymentMethod": "Credit Card",
				"date":          "01/01/2023",
			},
			expectedResult: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing title",
			requestBody: map[string]interface{}{
				"amount":        100.50,
				"currency":      "USD",
				"type":          "income",
				"category":      "Salary",
				"paymentMethod": "Credit Card",
				"date":          "01/01/2023",
			},
			expectedResult: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid currency length",
			requestBody: map[string]interface{}{
				"amount":        100.50,
				"title":         "Test Entry",
				"currency":      "USDD",
				"type":          "income",
				"category":      "Salary",
				"paymentMethod": "Credit Card",
				"date":          "01/01/2023",
			},
			expectedResult: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid type",
			requestBody: map[string]interface{}{
				"amount":        100.50,
				"title":         "Test Entry",
				"currency":      "USD",
				"type":          "invalid",
				"category":      "Salary",
				"paymentMethod": "Credit Card",
				"date":          "01/01/2023",
			},
			expectedResult: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing category",
			requestBody: map[string]interface{}{
				"amount":        100.50,
				"title":         "Test Entry",
				"currency":      "USD",
				"type":          "income",
				"paymentMethod": "Credit Card",
				"date":          "01/01/2023",
			},
			expectedResult: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing payment method",
			requestBody: map[string]interface{}{
				"amount":   100.50,
				"title":    "Test Entry",
				"currency": "USD",
				"type":     "income",
				"category": "Salary",
				"date":     "01/01/2023",
			},
			expectedResult: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid date format",
			requestBody: map[string]interface{}{
				"amount":        100.50,
				"title":         "Test Entry",
				"currency":      "USD",
				"type":          "income",
				"category":      "Salary",
				"paymentMethod": "Credit Card",
				"date":          "2023-01-01",
			},
			expectedResult: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid JSON",
			requestBody: map[string]interface{}{
				"amount": "not-a-number",
			},
			expectedResult: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty request body",
			requestBody:    map[string]interface{}{},
			expectedResult: false,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/transactions", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			entry, result := ValidateCreateTransactionsEntry(ctx)

			assert.Equal(t, tt.expectedResult, result)

			if tt.expectedResult {
				assert.NotNil(t, entry)

				if entry != nil {
					if val, ok := tt.requestBody["amount"].(float64); ok {
						assert.Equal(t, val, entry.Amount)
					}
					if val, ok := tt.requestBody["title"].(string); ok {
						assert.Equal(t, val, entry.Title)
					}
					if val, ok := tt.requestBody["currency"].(string); ok {
						assert.Equal(t, val, entry.Currency)
					}
					if val, ok := tt.requestBody["type"].(string); ok {
						assert.Equal(t, val, entry.Type)
					}
				}
			} else {
				assert.Equal(t, tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetCreateTransactionsValidationMessage(t *testing.T) {
	validate := validator.New()

	type TestStruct struct {
		Amount        float64 `validate:"required,gt=0"`
		Title         string  `validate:"required"`
		Currency      string  `validate:"required,len=3"`
		Type          string  `validate:"required,oneof=income expense"`
		Category      string  `validate:"min=1"`
		PaymentMethod string  `validate:"required,min=1"`
		Description   string  `validate:"min=1,lt=4"`
		Date          string  `validate:"required,datetime=02/01/2006"`
	}

	tests := []struct {
		name           string
		testStruct     TestStruct
		expectedField  string
		expectedTag    string
		expectedErrMsg string
	}{
		{
			name: "Required field missing",
			testStruct: TestStruct{
				Amount:        100,
				Currency:      "USD",
				Type:          "income",
				Category:      "Salary",
				PaymentMethod: "Credit Card",
				Date:          "01/01/2023",
			},
			expectedField:  "Title",
			expectedTag:    "required",
			expectedErrMsg: "This field is required",
		},
		{
			name: "Amount not greater than zero",
			testStruct: TestStruct{
				Amount:        -1,
				Title:         "Test",
				Currency:      "USD",
				Type:          "income",
				Category:      "Salary",
				PaymentMethod: "Credit Card",
				Date:          "01/01/2023",
			},
			expectedField:  "Amount",
			expectedTag:    "gt",
			expectedErrMsg: "Value must be greater than 0",
		},
		{
			name: "Currency not 3 characters",
			testStruct: TestStruct{
				Amount:        100,
				Title:         "Test",
				Currency:      "US",
				Type:          "income",
				Category:      "Salary",
				PaymentMethod: "Credit Card",
				Date:          "01/01/2023",
			},
			expectedField:  "Currency",
			expectedTag:    "len",
			expectedErrMsg: "Must be exactly 3 characters",
		},
		{
			name: "Invalid type value",
			testStruct: TestStruct{
				Amount:        100,
				Title:         "Test",
				Currency:      "USD",
				Type:          "invalid",
				Category:      "Salary",
				PaymentMethod: "Credit Card",
				Date:          "01/01/2023",
			},
			expectedField:  "Type",
			expectedTag:    "oneof",
			expectedErrMsg: "Must be either 'income' or 'expense'",
		},
		{
			name: "Category too short",
			testStruct: TestStruct{
				Amount:        100,
				Title:         "Test",
				Currency:      "USD",
				Type:          "income",
				Category:      "",
				PaymentMethod: "Credit Card",
				Description:   "description",
				Date:          "01/01/2023",
			},
			expectedField:  "Category",
			expectedTag:    "min",
			expectedErrMsg: "Value is too short",
		},
		{
			name: "Invalid date format",
			testStruct: TestStruct{
				Amount:        100,
				Title:         "Test",
				Currency:      "USD",
				Type:          "income",
				Category:      "Salary",
				PaymentMethod: "Credit Card",
				Date:          "2023-01-01",
			},
			expectedField:  "Date",
			expectedTag:    "datetime",
			expectedErrMsg: "Date must be in DD/MM/YYYY format (e.g., 31/12/2025)",
		},
		{
			name: "Invalid field",
			testStruct: TestStruct{
				Amount:        100,
				Title:         "Test",
				Currency:      "USD",
				Type:          "income",
				Category:      "Salary",
				PaymentMethod: "Credit Card",
				Description:   "description",
				Date:          "31/12/2025",
			},
			expectedField:  "Description",
			expectedTag:    "lt",
			expectedErrMsg: "Invalid value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.testStruct)
			if err != nil {
				validationErrors := err.(validator.ValidationErrors)
				for _, fieldError := range validationErrors {
					if fieldError.Field() == tt.expectedField && fieldError.Tag() == tt.expectedTag {
						message := getCreateTransactionsValidationMessage(fieldError)
						assert.Equal(t, tt.expectedErrMsg, message)
						return
					}
				}
				t.Errorf("Expected validation error for field %s with tag %s not found", tt.expectedField, tt.expectedTag)
			} else {
				t.Error("Expected validation error but got none")
			}
		})
	}
}
