package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"myfin-api/internal/dtos"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionsService struct {
	mock.Mock
}

func (m *MockTransactionsService) CreateTransactionsEntry(entry dtos.CreateTransactionsEntryDTO) (dtos.TransactionsEntryResponseDTO, error) {
	args := m.Called(entry)
	return args.Get(0).(dtos.TransactionsEntryResponseDTO), args.Error(1)
}

func (m *MockTransactionsService) GetAllTransactionsEntries(limit, skip int, titleFilter, categoryFilter string) ([]dtos.TransactionsEntryResponseDTO, error) {
	args := m.Called(limit, skip, titleFilter, categoryFilter)
	return args.Get(0).([]dtos.TransactionsEntryResponseDTO), args.Error(1)
}

func (m *MockTransactionsService) DeleteTransactionsEntry(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTransactionsService) UpdateTransactionsEntry(id string, entry dtos.UpdateTransactionsEntryDTO) (dtos.TransactionsEntryResponseDTO, error) {
	args := m.Called(id, entry)
	return args.Get(0).(dtos.TransactionsEntryResponseDTO), args.Error(1)
}

func (m *MockTransactionsService) GetTransactionsEntryByID(id string) (dtos.TransactionsEntryResponseDTO, error) {
	args := m.Called(id)
	return args.Get(0).(dtos.TransactionsEntryResponseDTO), args.Error(1)
}

func (m *MockTransactionsService) GetTransactionDashboardData() (dtos.TransactionDashboardResponseDTO, error) {
	args := m.Called()
	return args.Get(0).(dtos.TransactionDashboardResponseDTO), args.Error(1)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestSaveHandler(t *testing.T) {
	t.Run("successful_creation", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.POST("/transactions", func(c *gin.Context) {
			handler.Save(c)
		})

		validEntry := dtos.CreateTransactionsEntryDTO{
			Amount:        100.0,
			Title:         "Test Entry",
			Currency:      "USD",
			Type:          "expense",
			Category:      "food",
			PaymentMethod: "credit_card",
			Description:   "Test description",
			Date:          "15/03/2025",
		}

		expectedResponse := dtos.TransactionsEntryResponseDTO{
			ID:            "123456789012345678901234",
			Amount:        100.0,
			Title:         "Test Entry",
			Currency:      "USD",
			Type:          "expense",
			Category:      "food",
			PaymentMethod: "credit_card",
			Description:   "Test description",
			Date:          "15/03/2025",
			Timestamp:     1234567890,
			CreatedAt:     "2025-03-15T10:30:00Z",
			UpdatedAt:     "2025-03-15T10:30:00Z",
		}

		mockService.On("CreateTransactionsEntry", mock.AnythingOfType("dtos.CreateTransactionsEntryDTO")).Return(expectedResponse, nil)

		jsonPayload, _ := json.Marshal(validEntry)

		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response dtos.TransactionsEntryResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Amount, response.Amount)
		assert.Equal(t, expectedResponse.Title, response.Title)

		mockService.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.POST("/transactions", func(c *gin.Context) {
			handler.Save(c)
		})

		validEntry := dtos.CreateTransactionsEntryDTO{
			Amount:        100.0,
			Title:         "Test Entry",
			Currency:      "USD",
			Type:          "expense",
			Category:      "food",
			PaymentMethod: "credit_card",
			Description:   "Test description",
			Date:          "15/03/2025",
		}

		expectedError := errors.New("database error")
		mockService.On("CreateTransactionsEntry", mock.AnythingOfType("dtos.CreateTransactionsEntryDTO")).Return(dtos.TransactionsEntryResponseDTO{}, expectedError)

		jsonPayload, _ := json.Marshal(validEntry)

		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedError.Error(), response["error"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.POST("/transactions", func(c *gin.Context) {
			handler.Save(c)
		})

		invalidJSON := []byte(`{"amount": "invalid", "title": 123}`)

		req, _ := http.NewRequest("POST", "/transactions", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockService.AssertNotCalled(t, "CreateTransactionsEntry")
	})
}

func TestGetAllHandler(t *testing.T) {
	t.Run("successful_retrieval", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.GET("/transactions", func(c *gin.Context) {
			handler.GetAll(c)
		})

		expectedEntries := []dtos.TransactionsEntryResponseDTO{
			{
				ID:            "123456789012345678901234",
				Amount:        100.0,
				Title:         "Entry 1",
				Currency:      "USD",
				Type:          "expense",
				Category:      "food",
				PaymentMethod: "credit_card",
				Description:   "Description 1",
				Date:          "15/03/2025",
				Timestamp:     1234567890,
				CreatedAt:     "2025-03-15T10:30:00Z",
				UpdatedAt:     "2025-03-15T10:30:00Z",
			},
			{
				ID:            "234567890123456789012345",
				Amount:        200.0,
				Title:         "Entry 2",
				Currency:      "EUR",
				Type:          "income",
				Category:      "salary",
				PaymentMethod: "bank_transfer",
				Description:   "Description 2",
				Date:          "16/03/2025",
				Timestamp:     1234567891,
				CreatedAt:     "2025-03-16T10:30:00Z",
				UpdatedAt:     "2025-03-16T10:30:00Z",
			},
		}

		mockService.On("GetAllTransactionsEntries", 10, 0, "", "").Return(expectedEntries, nil)

		req, _ := http.NewRequest("GET", "/transactions?limit=10&skip=0", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data, ok := response["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 2)

		pagination, ok := response["pagination"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(10), pagination["limit"])
		assert.Equal(t, float64(0), pagination["skip"])
		assert.Equal(t, float64(2), pagination["count"])

		mockService.AssertExpectations(t)
	})

	t.Run("with_filters", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.GET("/transactions", func(c *gin.Context) {
			handler.GetAll(c)
		})

		filteredEntries := []dtos.TransactionsEntryResponseDTO{
			{
				ID:            "123456789012345678901234",
				Amount:        100.0,
				Title:         "Lunch",
				Currency:      "USD",
				Type:          "expense",
				Category:      "food",
				PaymentMethod: "credit_card",
				Description:   "Lunch expense",
				Date:          "15/03/2025",
				Timestamp:     1234567890,
				CreatedAt:     "2025-03-15T10:30:00Z",
				UpdatedAt:     "2025-03-15T10:30:00Z",
			},
		}

		mockService.On("GetAllTransactionsEntries", 10, 0, "lunch", "food").Return(filteredEntries, nil)

		req, _ := http.NewRequest("GET", "/transactions?limit=10&skip=0&title=lunch&category=food", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		data, ok := response["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 1)

		filters, ok := response["filters"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "lunch", filters["title"])
		assert.Equal(t, "food", filters["category"])

		mockService.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.GET("/transactions", func(c *gin.Context) {
			handler.GetAll(c)
		})

		expectedError := errors.New("database error")
		mockService.On("GetAllTransactionsEntries", 10, 0, "", "").Return([]dtos.TransactionsEntryResponseDTO{}, expectedError)

		req, _ := http.NewRequest("GET", "/transactions?limit=10&skip=0", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve entries", response["error"])
		assert.Equal(t, expectedError.Error(), response["details"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid_pagination_params", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.GET("/transactions", func(c *gin.Context) {
			handler.GetAll(c)
		})

		req, _ := http.NewRequest("GET", "/transactions?limit=invalid&skip=abc", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockService.AssertNotCalled(t, "GetAllTransactionsEntries")
	})
}

func TestDeleteHandler(t *testing.T) {
	t.Run("successful_deletion", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.DELETE("/transactions/:id", func(c *gin.Context) {
			handler.Delete(c)
		})

		validID := "123456789012345678901234"

		mockService.On("DeleteTransactionsEntry", validID).Return(nil)

		req, _ := http.NewRequest("DELETE", "/transactions/"+validID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Entry deleted successfully", response["message"])
		assert.Equal(t, validID, response["id"])

		mockService.AssertExpectations(t)
	})

	t.Run("missing_id", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.DELETE("/transactions/:id", func(c *gin.Context) {
			handler.Delete(c)
		})

		req, _ := http.NewRequest("DELETE", "/transactions/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertNotCalled(t, "DeleteTransactionsEntry")
	})

	t.Run("empty_id", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.DELETE("/transactions/", func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "id", Value: ""})
			handler.Delete(c)
		})

		req, _ := http.NewRequest("DELETE", "/transactions/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID is required", response["error"])

		mockService.AssertNotCalled(t, "DeleteTransactionsEntry")
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.DELETE("/transactions/:id", func(c *gin.Context) {
			handler.Delete(c)
		})

		validID := "123456789012345678901234"

		expectedError := errors.New("database error")
		mockService.On("DeleteTransactionsEntry", validID).Return(expectedError)

		req, _ := http.NewRequest("DELETE", "/transactions/"+validID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to delete entry", response["error"])
		assert.Equal(t, expectedError.Error(), response["details"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid_id_format", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.DELETE("/transactions/:id", func(c *gin.Context) {
			handler.Delete(c)
		})

		invalidID := "invalid-id-format"

		expectedError := errors.New("the provided hex string is not a valid ObjectID")
		mockService.On("DeleteTransactionsEntry", invalidID).Return(expectedError)

		req, _ := http.NewRequest("DELETE", "/transactions/"+invalidID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to delete entry", response["error"])
		assert.Equal(t, expectedError.Error(), response["details"])

		mockService.AssertExpectations(t)
	})
}

func TestUpdateHandler(t *testing.T) {
	t.Run("successful_update", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.PUT("/transactions/:id", func(c *gin.Context) {
			handler.Update(c)
		})

		validID := "123456789012345678901234"

		validEntry := dtos.UpdateTransactionsEntryDTO{
			Amount:        200.50,
			Title:         "Updated Entry",
			Currency:      "USD",
			Type:          "expense",
			Category:      "entertainment",
			PaymentMethod: "debit_card",
			Description:   "Updated description",
			Date:          "15/10/2025",
		}

		expectedResponse := dtos.TransactionsEntryResponseDTO{
			ID:            validID,
			Amount:        200.50,
			Title:         "Updated Entry",
			Currency:      "USD",
			Type:          "expense",
			Category:      "entertainment",
			PaymentMethod: "debit_card",
			Description:   "Updated description",
			Date:          "15/10/2025",
			Timestamp:     1234567890,
			CreatedAt:     "2025-03-15T10:30:00Z",
			UpdatedAt:     "2025-10-15T10:30:00Z",
		}

		mockService.On("UpdateTransactionsEntry", validID, mock.AnythingOfType("dtos.UpdateTransactionsEntryDTO")).Return(expectedResponse, nil)

		jsonPayload, _ := json.Marshal(validEntry)

		req, _ := http.NewRequest("PUT", "/transactions/"+validID, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Entry updated successfully", response["message"])

		data, ok := response["data"].(map[string]interface{})
		assert.True(t, ok)

		assert.Equal(t, validID, data["id"])
		assert.Equal(t, 200.50, data["amount"])
		assert.Equal(t, "Updated Entry", data["title"])
		assert.Equal(t, "USD", data["currency"])
		assert.Equal(t, "expense", data["type"])
		assert.Equal(t, "entertainment", data["category"])
		assert.Equal(t, "debit_card", data["paymentMethod"])
		assert.Equal(t, "Updated description", data["description"])
		assert.Equal(t, "15/10/2025", data["date"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.PUT("/transactions/:id", func(c *gin.Context) {
			handler.Update(c)
		})

		validID := "123456789012345678901234"

		invalidJSON := []byte(`{"amount": "invalid", "title": 123}`)

		req, _ := http.NewRequest("PUT", "/transactions/"+validID, bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockService.AssertNotCalled(t, "UpdateTransactionsEntry")
	})

	t.Run("missing_id", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.PUT("/transactions/:id", func(c *gin.Context) {
			handler.Update(c)
		})

		validEntry := dtos.UpdateTransactionsEntryDTO{
			Amount:        200.50,
			Title:         "Updated Entry",
			Currency:      "USD",
			Type:          "expense",
			Category:      "entertainment",
			PaymentMethod: "debit_card",
			Description:   "Updated description",
			Date:          "15/10/2025",
		}

		jsonPayload, _ := json.Marshal(validEntry)

		req, _ := http.NewRequest("PUT", "/transactions/", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertNotCalled(t, "UpdateTransactionsEntry")
	})

	t.Run("empty_id", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.PUT("/transactions/", func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "id", Value: ""})
			handler.Update(c)
		})

		validEntry := dtos.UpdateTransactionsEntryDTO{
			Amount:        200.50,
			Title:         "Updated Entry",
			Currency:      "USD",
			Type:          "expense",
			Category:      "entertainment",
			PaymentMethod: "debit_card",
			Description:   "Updated description",
			Date:          "15/10/2025",
		}

		jsonPayload, _ := json.Marshal(validEntry)

		req, _ := http.NewRequest("PUT", "/transactions/", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockService.AssertNotCalled(t, "UpdateTransactionsEntry")
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.PUT("/transactions/:id", func(c *gin.Context) {
			handler.Update(c)
		})

		validID := "123456789012345678901234"

		validEntry := dtos.UpdateTransactionsEntryDTO{
			Amount:        200.50,
			Title:         "Updated Entry",
			Currency:      "USD",
			Type:          "expense",
			Category:      "entertainment",
			PaymentMethod: "debit_card",
			Description:   "Updated description",
			Date:          "15/10/2025",
		}

		expectedError := errors.New("database error")
		mockService.On("UpdateTransactionsEntry", validID, mock.AnythingOfType("dtos.UpdateTransactionsEntryDTO")).Return(dtos.TransactionsEntryResponseDTO{}, expectedError)

		jsonPayload, _ := json.Marshal(validEntry)

		req, _ := http.NewRequest("PUT", "/transactions/"+validID, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to update entry", response["error"])
		assert.Equal(t, expectedError.Error(), response["details"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid_date_format", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.PUT("/transactions/:id", func(c *gin.Context) {
			handler.Update(c)
		})

		validID := "123456789012345678901234"

		invalidEntry := map[string]interface{}{
			"amount":        200.50,
			"title":         "Updated Entry",
			"currency":      "USD",
			"type":          "expense",
			"category":      "entertainment",
			"paymentMethod": "debit_card",
			"description":   "Updated description",
			"date":          "2025-10-15", // Wrong format, should be DD/MM/YYYY
		}

		jsonPayload, _ := json.Marshal(invalidEntry)

		req, _ := http.NewRequest("PUT", "/transactions/"+validID, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Date must be in format")

		mockService.AssertNotCalled(t, "UpdateTransactionsEntry")
	})
}

func TestGetByIDHandler(t *testing.T) {
	t.Run("successful_retrieval", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.GET("/transactions/:id", func(c *gin.Context) {
			handler.GetByID(c)
		})

		validID := "123456789012345678901234"

		expectedResponse := dtos.TransactionsEntryResponseDTO{
			ID:            validID,
			Amount:        150.75,
			Title:         "Lunch at restaurant",
			Currency:      "BRL",
			Type:          "expense",
			Category:      "food",
			PaymentMethod: "credit_card",
			Description:   "Lunch at restaurant",
			Date:          "06/09/2025",
			Timestamp:     1234567890,
			CreatedAt:     "2025-09-06T14:30:00Z",
			UpdatedAt:     "2025-09-06T14:30:00Z",
		}

		mockService.On("GetTransactionsEntryByID", validID).Return(expectedResponse, nil)

		req, _ := http.NewRequest("GET", "/transactions/"+validID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dtos.TransactionsEntryResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Amount, response.Amount)
		assert.Equal(t, expectedResponse.Title, response.Title)
		assert.Equal(t, expectedResponse.Currency, response.Currency)
		assert.Equal(t, expectedResponse.Type, response.Type)
		assert.Equal(t, expectedResponse.Category, response.Category)
		assert.Equal(t, expectedResponse.PaymentMethod, response.PaymentMethod)
		assert.Equal(t, expectedResponse.Description, response.Description)
		assert.Equal(t, expectedResponse.Date, response.Date)
		assert.Equal(t, expectedResponse.Timestamp, response.Timestamp)
		assert.Equal(t, expectedResponse.CreatedAt, response.CreatedAt)
		assert.Equal(t, expectedResponse.UpdatedAt, response.UpdatedAt)

		mockService.AssertExpectations(t)
	})

	t.Run("missing_id", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.GET("/transactions/:id", func(c *gin.Context) {
			handler.GetByID(c)
		})

		req, _ := http.NewRequest("GET", "/transactions/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertNotCalled(t, "GetTransactionsEntryByID")
	})

	t.Run("empty_id", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.GET("/transactions/", func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "id", Value: ""})
			handler.GetByID(c)
		})

		req, _ := http.NewRequest("GET", "/transactions/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID is required", response["error"])

		mockService.AssertNotCalled(t, "GetTransactionsEntryByID")
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.GET("/transactions/:id", func(c *gin.Context) {
			handler.GetByID(c)
		})

		validID := "123456789012345678901234"

		expectedError := errors.New("database error")
		mockService.On("GetTransactionsEntryByID", validID).Return(dtos.TransactionsEntryResponseDTO{}, expectedError)

		req, _ := http.NewRequest("GET", "/transactions/"+validID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve entry", response["error"])
		assert.Equal(t, expectedError.Error(), response["details"])

		mockService.AssertExpectations(t)
	})

	t.Run("not_found_error", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.GET("/transactions/:id", func(c *gin.Context) {
			handler.GetByID(c)
		})

		validID := "123456789012345678901234"

		expectedError := errors.New("entry not found")
		mockService.On("GetTransactionsEntryByID", validID).Return(dtos.TransactionsEntryResponseDTO{}, expectedError)

		req, _ := http.NewRequest("GET", "/transactions/"+validID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve entry", response["error"])
		assert.Equal(t, expectedError.Error(), response["details"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid_id_format", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.GET("/transactions/:id", func(c *gin.Context) {
			handler.GetByID(c)
		})

		invalidID := "invalid-id-format"

		expectedError := errors.New("the provided hex string is not a valid ObjectID")
		mockService.On("GetTransactionsEntryByID", invalidID).Return(dtos.TransactionsEntryResponseDTO{}, expectedError)

		req, _ := http.NewRequest("GET", "/transactions/"+invalidID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve entry", response["error"])
		assert.Equal(t, expectedError.Error(), response["details"])

		mockService.AssertExpectations(t)
	})
}

func TestTransactionsHandlerGetTransactionDashboardData(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.GET("/transactions/dashboard", func(c *gin.Context) {
			handler.GetTransactionDashboardData(c)
		})

		expectedResponse := dtos.TransactionDashboardResponseDTO{
			IncomeAmount:  1000.50,
			ExpenseAmount: 250.75,
			TotalAmount:   749.75,
		}

		mockService.On("GetTransactionDashboardData").Return(expectedResponse, nil)

		req, _ := http.NewRequest("GET", "/transactions/dashboard", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dtos.TransactionDashboardResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.IncomeAmount, response.IncomeAmount)
		assert.Equal(t, expectedResponse.ExpenseAmount, response.ExpenseAmount)
		assert.Equal(t, expectedResponse.TotalAmount, response.TotalAmount)

		mockService.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.GET("/transactions/dashboard", func(c *gin.Context) {
			handler.GetTransactionDashboardData(c)
		})

		expectedError := errors.New("database connection failed")
		mockService.On("GetTransactionDashboardData").Return(dtos.TransactionDashboardResponseDTO{}, expectedError)

		req, _ := http.NewRequest("GET", "/transactions/dashboard", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve dashboard data", response["error"])
		assert.Equal(t, expectedError.Error(), response["details"])

		mockService.AssertExpectations(t)
	})

	t.Run("empty_dashboard_data", func(t *testing.T) {
		mockService := new(MockTransactionsService)
		handler := NewTransactionsHandler(mockService)
		router := setupRouter()

		router.GET("/transactions/dashboard", func(c *gin.Context) {
			handler.GetTransactionDashboardData(c)
		})

		expectedResponse := dtos.TransactionDashboardResponseDTO{
			IncomeAmount:  0.0,
			ExpenseAmount: 0.0,
			TotalAmount:   0.0,
		}

		mockService.On("GetTransactionDashboardData").Return(expectedResponse, nil)

		req, _ := http.NewRequest("GET", "/transactions/dashboard", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dtos.TransactionDashboardResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 0.0, response.IncomeAmount)
		assert.Equal(t, 0.0, response.ExpenseAmount)
		assert.Equal(t, 0.0, response.TotalAmount)

		mockService.AssertExpectations(t)
	})
}
