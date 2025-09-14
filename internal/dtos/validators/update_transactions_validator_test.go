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

func TestValidateUpdateTransactionsEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		id             string
		requestBody    map[string]interface{}
		expectedResult bool
		expectedStatus int
	}{
		{
			name: "Valid request",
			id:   "123",
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
			name:           "Missing ID",
			id:             "",
			requestBody:    map[string]interface{}{},
			expectedResult: false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid amount",
			id:   "123",
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
			id:   "123",
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
			name: "Invalid type",
			id:   "123",
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
			name: "Invalid JSON",
			id:   "123",
			requestBody: map[string]interface{}{
				"amount": "not-a-number",
			},
			expectedResult: false,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPut, "/transactions/"+tt.id, bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			ctx.Params = []gin.Param{{Key: "id", Value: tt.id}}

			entry, id, result := ValidateUpdateTransactionsEntry(ctx)

			assert.Equal(t, tt.expectedResult, result)

			if tt.expectedResult {
				assert.NotNil(t, entry)
				assert.Equal(t, tt.id, id)
			} else {
				assert.Equal(t, tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestGetUpdateTransactionsValidationMessage(t *testing.T) {
	validate := validator.New()

	type TestStruct struct {
		Amount        float64 `validate:"required,gt=0"`
		Title         string  `validate:"required"`
		Currency      string  `validate:"required,len=3"`
		Type          string  `validate:"required,oneof=income expense"`
		Category      string  `validate:"required"`
		PaymentMethod string  `validate:"required"`
		Description   string  `validate:"min=1"`
		Date          string  `validate:"required"`
		Test          string  `validate:"required"`
	}

	tests := []struct {
		name           string
		testStruct     TestStruct
		expectedField  string
		expectedErrMsg string
	}{
		{
			name: "Invalid amount",
			testStruct: TestStruct{
				Amount:        0,
				Title:         "Test",
				Currency:      "USD",
				Type:          "income",
				Category:      "Salary",
				PaymentMethod: "Credit Card",
				Date:          "01/01/2023",
			},
			expectedField:  "Amount",
			expectedErrMsg: "Amount must be greater than 0",
		},
		{
			name: "Missing title",
			testStruct: TestStruct{
				Amount:        100,
				Title:         "",
				Currency:      "USD",
				Type:          "income",
				Category:      "Salary",
				PaymentMethod: "Credit Card",
				Date:          "01/01/2023",
			},
			expectedField:  "Title",
			expectedErrMsg: "Title is required",
		},
		{
			name: "Invalid currency",
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
			expectedErrMsg: "Currency must be a 3-letter code",
		},
		{
			name: "Invalid type",
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
			expectedErrMsg: "Type must be either 'income' or 'expense'",
		},
		{
			name: "Missing category",
			testStruct: TestStruct{
				Amount:        100,
				Title:         "Test",
				Currency:      "USD",
				Type:          "income",
				Category:      "",
				PaymentMethod: "Credit Card",
				Date:          "01/01/2023",
			},
			expectedField:  "Category",
			expectedErrMsg: "Category is required",
		},
		{
			name: "Missing payment method",
			testStruct: TestStruct{
				Amount:        100,
				Title:         "Test",
				Currency:      "USD",
				Type:          "income",
				Category:      "Salary",
				PaymentMethod: "",
				Date:          "01/01/2023",
			},
			expectedField:  "PaymentMethod",
			expectedErrMsg: "Payment method is required",
		},
		{
			name: "Empty description",
			testStruct: TestStruct{
				Amount:        100,
				Title:         "Test",
				Currency:      "USD",
				Type:          "income",
				Category:      "Salary",
				PaymentMethod: "Credit Card",
				Description:   "",
				Date:          "01/01/2023",
			},
			expectedField:  "Description",
			expectedErrMsg: "Description must not be empty",
		},
		{
			name: "Missing date",
			testStruct: TestStruct{
				Amount:        100,
				Title:         "Test",
				Currency:      "USD",
				Type:          "income",
				Category:      "Salary",
				PaymentMethod: "Credit Card",
				Date:          "",
			},
			expectedField:  "Date",
			expectedErrMsg: "Date must be in format DD/MM/YYYY",
		},
		{
			name: "Invalid value",
			testStruct: TestStruct{
				Amount:        100,
				Title:         "Test",
				Currency:      "USD",
				Type:          "income",
				Category:      "Salary",
				PaymentMethod: "Credit Card",
				Description:   "fsdadfsfsfsd",
				Date:          "01/01/2023",
			},
			expectedField:  "Test",
			expectedErrMsg: "Key: 'TestStruct.Test' Error:Field validation for 'Test' failed on the 'required' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.testStruct)
			if err != nil {
				validationErrors := err.(validator.ValidationErrors)
				for _, fieldError := range validationErrors {
					if fieldError.Field() == tt.expectedField {
						message := getUpdateTransactionsValidationMessage(fieldError)
						assert.Equal(t, tt.expectedErrMsg, message)
						return
					}
				}
				t.Errorf("Expected validation error for field %s not found", tt.expectedField)
			} else {
				t.Error("Expected validation error but got none")
			}
		})
	}
}
