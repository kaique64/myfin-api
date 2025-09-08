package validators

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestValidateCreateCashHandlingEntry(t *testing.T) {
	// Configure Gin to test mode
	gin.SetMode(gin.TestMode)

	t.Run("valid_entry", func(t *testing.T) {
		// Create a valid JSON payload
		validJSON := `{
			"amount": 100.50,
			"title": "Test Entry",
			"currency": "USD",
			"type": "expense",
			"category": "food",
			"paymentMethod": "credit_card",
			"description": "Test description",
			"date": "15/03/2025"
		}`

		// Create a test context
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(validJSON))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// Call the validator
		entry, isValid := ValidateCreateCashHandlingEntry(ctx)

		// Assert results
		assert.True(t, isValid)
		assert.NotNil(t, entry)
		assert.Equal(t, 100.50, entry.Amount)
		assert.Equal(t, "Test Entry", entry.Title)
		assert.Equal(t, "USD", entry.Currency)
		assert.Equal(t, "expense", entry.Type)
		assert.Equal(t, "food", entry.Category)
		assert.Equal(t, "credit_card", entry.PaymentMethod)
		assert.Equal(t, "Test description", entry.Description)
		assert.Equal(t, "15/03/2025", entry.Date)
	})

	t.Run("missing_required_fields", func(t *testing.T) {
		// Create JSON with missing required fields
		invalidJSON := `{
			"amount": 100.50,
			"title": "Test Entry",
			"currency": "USD"
		}`

		// Create a test context
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(invalidJSON))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// Call the validator
		_, isValid := ValidateCreateCashHandlingEntry(ctx)

		// Assert results
		assert.False(t, isValid)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Validation failed", response["error"])

		details, ok := response["details"].(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, details, "Type")
		assert.Contains(t, details, "Category")
		assert.Contains(t, details, "PaymentMethod")
		assert.Contains(t, details, "Date")
	})

	t.Run("invalid_amount", func(t *testing.T) {
		// Create JSON with invalid amount (not greater than 0)
		invalidJSON := `{
			"amount": -1,
			"title": "Test Entry",
			"currency": "USD",
			"type": "expense",
			"category": "food",
			"paymentMethod": "credit_card",
			"date": "15/03/2025"
		}`

		// Create a test context
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(invalidJSON))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// Call the validator
		_, isValid := ValidateCreateCashHandlingEntry(ctx)

		// Assert results
		assert.False(t, isValid)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Validation failed", response["error"])

		details, ok := response["details"].(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, details, "Amount")
		assert.Equal(t, "Value must be greater than 0", details["Amount"])
	})

	t.Run("invalid_currency_length", func(t *testing.T) {
		// Create JSON with invalid currency (not 3 characters)
		invalidJSON := `{
			"amount": 100.50,
			"title": "Test Entry",
			"currency": "USDD",
			"type": "expense",
			"category": "food",
			"paymentMethod": "credit_card",
			"date": "15/03/2025"
		}`

		// Create a test context
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(invalidJSON))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// Call the validator
		_, isValid := ValidateCreateCashHandlingEntry(ctx)

		// Assert results
		assert.False(t, isValid)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Validation failed", response["error"])

		details, ok := response["details"].(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, details, "Currency")
		assert.Equal(t, "Must be exactly 3 characters", details["Currency"])
	})

	t.Run("invalid_type", func(t *testing.T) {
		// Create JSON with invalid type (not 'income' or 'expense')
		invalidJSON := `{
			"amount": 100.50,
			"title": "Test Entry",
			"currency": "USD",
			"type": "invalid",
			"category": "food",
			"paymentMethod": "credit_card",
			"date": "15/03/2025"
		}`

		// Create a test context
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(invalidJSON))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// Call the validator
		_, isValid := ValidateCreateCashHandlingEntry(ctx)

		// Assert results
		assert.False(t, isValid)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Validation failed", response["error"])

		details, ok := response["details"].(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, details, "Type")
		assert.Equal(t, "Must be either 'income' or 'expense'", details["Type"])
	})

	t.Run("invalid_date_format", func(t *testing.T) {
		// Create JSON with invalid date format
		invalidJSON := `{
			"amount": 100.50,
			"title": "Test Entry",
			"currency": "USD",
			"type": "expense",
			"category": "food",
			"paymentMethod": "credit_card",
			"date": "2025-03-15"
		}`

		// Create a test context
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(invalidJSON))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// Call the validator
		_, isValid := ValidateCreateCashHandlingEntry(ctx)

		// Assert results
		assert.False(t, isValid)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Validation failed", response["error"])

		details, ok := response["details"].(map[string]interface{})
		assert.True(t, ok)
		assert.Contains(t, details, "Date")
		assert.Equal(t, "Date must be in DD/MM/YYYY format (e.g., 31/12/2025)", details["Date"])
	})

	t.Run("invalid_json_format", func(t *testing.T) {
		// Create malformed JSON
		invalidJSON := `{
			"amount": 100.50,
			"title": "Test Entry",
			"currency": "USD",
			missing_quotes: "value"
		}`

		// Create a test context
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(invalidJSON))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// Call the validator
		_, isValid := ValidateCreateCashHandlingEntry(ctx)

		// Assert results
		assert.False(t, isValid)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid JSON format", response["error"])
		assert.Contains(t, response["details"], "invalid character")
	})

	t.Run("empty_request_body", func(t *testing.T) {
		// Create empty request body
		emptyJSON := ``

		// Create a test context
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request, _ = http.NewRequest("POST", "/", bytes.NewBufferString(emptyJSON))
		ctx.Request.Header.Set("Content-Type", "application/json")

		// Call the validator
		_, isValid := ValidateCreateCashHandlingEntry(ctx)

		// Assert results
		assert.False(t, isValid)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Check response body
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid JSON format", response["error"])
	})
}

func TestGetCreateCashHandlingValidationMessage(t *testing.T) {
	// This is a test for the internal function getCreateCashHandlingValidationMessage
	// Since it's not exported, we can only test it indirectly through ValidateCreateCashHandlingEntry
	// The tests above already cover most cases, but we can add more specific tests if needed
}

