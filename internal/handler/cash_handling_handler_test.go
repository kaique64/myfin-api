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

type MockCashHandlingService struct {
	mock.Mock
}

func (m *MockCashHandlingService) CreateCashHandlingEntry(entry dtos.CreateCashHandlingEntryDTO) (dtos.CashHandlingEntryResponseDTO, error) {
	args := m.Called(entry)
	return args.Get(0).(dtos.CashHandlingEntryResponseDTO), args.Error(1)
}

func (m *MockCashHandlingService) GetAllCashHandlingEntries(limit, skip int, titleFilter, categoryFilter string) ([]dtos.CashHandlingEntryResponseDTO, error) {
	args := m.Called(limit, skip, titleFilter, categoryFilter)
	return args.Get(0).([]dtos.CashHandlingEntryResponseDTO), args.Error(1)
}

func (m *MockCashHandlingService) DeleteCashHandlingEntry(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCashHandlingService) UpdateCashHandlingEntry(id string, entry dtos.UpdateCashHandlingEntryDTO) (dtos.CashHandlingEntryResponseDTO, error) {
	args := m.Called(id, entry)
	return args.Get(0).(dtos.CashHandlingEntryResponseDTO), args.Error(1)
}

func (m *MockCashHandlingService) GetCashHandlingEntryByID(id string) (dtos.CashHandlingEntryResponseDTO, error) {
	args := m.Called(id)
	return args.Get(0).(dtos.CashHandlingEntryResponseDTO), args.Error(1)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

func TestSaveHandler(t *testing.T) {
	t.Run("successful_creation", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.POST("/cash-handling", func(c *gin.Context) {
			handler.Save(c)
		})

		validEntry := dtos.CreateCashHandlingEntryDTO{
			Amount:        100.0,
			Title:         "Test Entry",
			Currency:      "USD",
			Type:          "expense",
			Category:      "food",
			PaymentMethod: "credit_card",
			Description:   "Test description",
			Date:          "15/03/2025",
		}

		expectedResponse := dtos.CashHandlingEntryResponseDTO{
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

		mockService.On("CreateCashHandlingEntry", mock.AnythingOfType("dtos.CreateCashHandlingEntryDTO")).Return(expectedResponse, nil)

		jsonPayload, _ := json.Marshal(validEntry)

		req, _ := http.NewRequest("POST", "/cash-handling", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response dtos.CashHandlingEntryResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Amount, response.Amount)
		assert.Equal(t, expectedResponse.Title, response.Title)

		mockService.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.POST("/cash-handling", func(c *gin.Context) {
			handler.Save(c)
		})

		validEntry := dtos.CreateCashHandlingEntryDTO{
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
		mockService.On("CreateCashHandlingEntry", mock.AnythingOfType("dtos.CreateCashHandlingEntryDTO")).Return(dtos.CashHandlingEntryResponseDTO{}, expectedError)

		jsonPayload, _ := json.Marshal(validEntry)

		req, _ := http.NewRequest("POST", "/cash-handling", bytes.NewBuffer(jsonPayload))
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
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.POST("/cash-handling", func(c *gin.Context) {
			handler.Save(c)
		})

		invalidJSON := []byte(`{"amount": "invalid", "title": 123}`)

		req, _ := http.NewRequest("POST", "/cash-handling", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockService.AssertNotCalled(t, "CreateCashHandlingEntry")
	})
}

func TestGetAllHandler(t *testing.T) {
	t.Run("successful_retrieval", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.GET("/cash-handling", func(c *gin.Context) {
			handler.GetAll(c)
		})

		expectedEntries := []dtos.CashHandlingEntryResponseDTO{
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

		mockService.On("GetAllCashHandlingEntries", 10, 0, "", "").Return(expectedEntries, nil)

		req, _ := http.NewRequest("GET", "/cash-handling?limit=10&skip=0", nil)
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
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.GET("/cash-handling", func(c *gin.Context) {
			handler.GetAll(c)
		})

		filteredEntries := []dtos.CashHandlingEntryResponseDTO{
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

		mockService.On("GetAllCashHandlingEntries", 10, 0, "lunch", "food").Return(filteredEntries, nil)

		req, _ := http.NewRequest("GET", "/cash-handling?limit=10&skip=0&title=lunch&category=food", nil)
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
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.GET("/cash-handling", func(c *gin.Context) {
			handler.GetAll(c)
		})

		expectedError := errors.New("database error")
		mockService.On("GetAllCashHandlingEntries", 10, 0, "", "").Return([]dtos.CashHandlingEntryResponseDTO{}, expectedError)

		req, _ := http.NewRequest("GET", "/cash-handling?limit=10&skip=0", nil)
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
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.GET("/cash-handling", func(c *gin.Context) {
			handler.GetAll(c)
		})

		req, _ := http.NewRequest("GET", "/cash-handling?limit=invalid&skip=abc", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockService.AssertNotCalled(t, "GetAllCashHandlingEntries")
	})
}

func TestDeleteHandler(t *testing.T) {
	t.Run("successful_deletion", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.DELETE("/cash-handling/:id", func(c *gin.Context) {
			handler.Delete(c)
		})

		validID := "123456789012345678901234"

		mockService.On("DeleteCashHandlingEntry", validID).Return(nil)

		req, _ := http.NewRequest("DELETE", "/cash-handling/"+validID, nil)
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
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.DELETE("/cash-handling/:id", func(c *gin.Context) {
			handler.Delete(c)
		})

		req, _ := http.NewRequest("DELETE", "/cash-handling/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertNotCalled(t, "DeleteCashHandlingEntry")
	})

	t.Run("empty_id", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.DELETE("/cash-handling/", func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "id", Value: ""})
			handler.Delete(c)
		})

		req, _ := http.NewRequest("DELETE", "/cash-handling/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID is required", response["error"])

		mockService.AssertNotCalled(t, "DeleteCashHandlingEntry")
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.DELETE("/cash-handling/:id", func(c *gin.Context) {
			handler.Delete(c)
		})

		validID := "123456789012345678901234"

		expectedError := errors.New("database error")
		mockService.On("DeleteCashHandlingEntry", validID).Return(expectedError)

		req, _ := http.NewRequest("DELETE", "/cash-handling/"+validID, nil)
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
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.DELETE("/cash-handling/:id", func(c *gin.Context) {
			handler.Delete(c)
		})

		invalidID := "invalid-id-format"

		expectedError := errors.New("the provided hex string is not a valid ObjectID")
		mockService.On("DeleteCashHandlingEntry", invalidID).Return(expectedError)

		req, _ := http.NewRequest("DELETE", "/cash-handling/"+invalidID, nil)
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
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.PUT("/cash-handling/:id", func(c *gin.Context) {
			handler.Update(c)
		})

		validID := "123456789012345678901234"

		validEntry := dtos.UpdateCashHandlingEntryDTO{
			Amount:        200.50,
			Title:         "Updated Entry",
			Currency:      "USD",
			Type:          "expense",
			Category:      "entertainment",
			PaymentMethod: "debit_card",
			Description:   "Updated description",
			Date:          "15/10/2025",
		}

		expectedResponse := dtos.CashHandlingEntryResponseDTO{
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

		mockService.On("UpdateCashHandlingEntry", validID, mock.AnythingOfType("dtos.UpdateCashHandlingEntryDTO")).Return(expectedResponse, nil)

		jsonPayload, _ := json.Marshal(validEntry)

		req, _ := http.NewRequest("PUT", "/cash-handling/"+validID, bytes.NewBuffer(jsonPayload))
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

		assert.Equal(t, validID, data["ID"])
		assert.Equal(t, 200.50, data["Amount"])
		assert.Equal(t, "Updated Entry", data["Title"])
		assert.Equal(t, "USD", data["Currency"])
		assert.Equal(t, "expense", data["Type"])
		assert.Equal(t, "entertainment", data["Category"])
		assert.Equal(t, "debit_card", data["PaymentMethod"])
		assert.Equal(t, "Updated description", data["Description"])
		assert.Equal(t, "15/10/2025", data["Date"])

		mockService.AssertExpectations(t)
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.PUT("/cash-handling/:id", func(c *gin.Context) {
			handler.Update(c)
		})

		validID := "123456789012345678901234"

		invalidJSON := []byte(`{"amount": "invalid", "title": 123}`)

		req, _ := http.NewRequest("PUT", "/cash-handling/"+validID, bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockService.AssertNotCalled(t, "UpdateCashHandlingEntry")
	})

	t.Run("missing_id", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.PUT("/cash-handling/:id", func(c *gin.Context) {
			handler.Update(c)
		})

		validEntry := dtos.UpdateCashHandlingEntryDTO{
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

		req, _ := http.NewRequest("PUT", "/cash-handling/", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertNotCalled(t, "UpdateCashHandlingEntry")
	})

	t.Run("empty_id", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.PUT("/cash-handling/", func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "id", Value: ""})
			handler.Update(c)
		})

		validEntry := dtos.UpdateCashHandlingEntryDTO{
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

		req, _ := http.NewRequest("PUT", "/cash-handling/", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		mockService.AssertNotCalled(t, "UpdateCashHandlingEntry")
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.PUT("/cash-handling/:id", func(c *gin.Context) {
			handler.Update(c)
		})

		validID := "123456789012345678901234"

		validEntry := dtos.UpdateCashHandlingEntryDTO{
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
		mockService.On("UpdateCashHandlingEntry", validID, mock.AnythingOfType("dtos.UpdateCashHandlingEntryDTO")).Return(dtos.CashHandlingEntryResponseDTO{}, expectedError)

		jsonPayload, _ := json.Marshal(validEntry)

		req, _ := http.NewRequest("PUT", "/cash-handling/"+validID, bytes.NewBuffer(jsonPayload))
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
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.PUT("/cash-handling/:id", func(c *gin.Context) {
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

		req, _ := http.NewRequest("PUT", "/cash-handling/"+validID, bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response["error"], "Date must be in format")

		mockService.AssertNotCalled(t, "UpdateCashHandlingEntry")
	})
}

func TestGetByIDHandler(t *testing.T) {
	t.Run("successful_retrieval", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.GET("/cash-handling/:id", func(c *gin.Context) {
			handler.GetByID(c)
		})

		validID := "123456789012345678901234"

		expectedResponse := dtos.CashHandlingEntryResponseDTO{
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

		mockService.On("GetCashHandlingEntryByID", validID).Return(expectedResponse, nil)

		req, _ := http.NewRequest("GET", "/cash-handling/"+validID, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response dtos.CashHandlingEntryResponseDTO
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
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.GET("/cash-handling/:id", func(c *gin.Context) {
			handler.GetByID(c)
		})

		req, _ := http.NewRequest("GET", "/cash-handling/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertNotCalled(t, "GetCashHandlingEntryByID")
	})

	t.Run("empty_id", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.GET("/cash-handling/", func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "id", Value: ""})
			handler.GetByID(c)
		})

		req, _ := http.NewRequest("GET", "/cash-handling/", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID is required", response["error"])

		mockService.AssertNotCalled(t, "GetCashHandlingEntryByID")
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.GET("/cash-handling/:id", func(c *gin.Context) {
			handler.GetByID(c)
		})

		validID := "123456789012345678901234"

		expectedError := errors.New("database error")
		mockService.On("GetCashHandlingEntryByID", validID).Return(dtos.CashHandlingEntryResponseDTO{}, expectedError)

		req, _ := http.NewRequest("GET", "/cash-handling/"+validID, nil)
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
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.GET("/cash-handling/:id", func(c *gin.Context) {
			handler.GetByID(c)
		})

		validID := "123456789012345678901234"

		expectedError := errors.New("entry not found")
		mockService.On("GetCashHandlingEntryByID", validID).Return(dtos.CashHandlingEntryResponseDTO{}, expectedError)

		req, _ := http.NewRequest("GET", "/cash-handling/"+validID, nil)
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
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		router.GET("/cash-handling/:id", func(c *gin.Context) {
			handler.GetByID(c)
		})

		invalidID := "invalid-id-format"

		expectedError := errors.New("the provided hex string is not a valid ObjectID")
		mockService.On("GetCashHandlingEntryByID", invalidID).Return(dtos.CashHandlingEntryResponseDTO{}, expectedError)

		req, _ := http.NewRequest("GET", "/cash-handling/"+invalidID, nil)
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
