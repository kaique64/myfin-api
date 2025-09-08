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

// MockCashHandlingService é um mock do serviço para testes
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

// setupRouter configura o router do Gin para testes
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	return router
}

// TestSaveHandler testa o método Save do handler
func TestSaveHandler(t *testing.T) {
	// Testes para o método Save
	t.Run("successful_creation", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		// Configurar a rota para o teste
		router.POST("/cash-handling", func(c *gin.Context) {
			handler.Save(c)
		})

		// Criar um payload válido
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

		// Configurar o mock para retornar uma resposta de sucesso
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

		// Converter o payload para JSON
		jsonPayload, _ := json.Marshal(validEntry)

		// Criar uma requisição HTTP de teste
		req, _ := http.NewRequest("POST", "/cash-handling", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Executar a requisição
		router.ServeHTTP(w, req)

		// Verificar o resultado
		assert.Equal(t, http.StatusCreated, w.Code)

		// Verificar o corpo da resposta
		var response dtos.CashHandlingEntryResponseDTO
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Amount, response.Amount)
		assert.Equal(t, expectedResponse.Title, response.Title)

		// Verificar que o mock foi chamado conforme esperado
		mockService.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		// Configurar a rota para o teste
		router.POST("/cash-handling", func(c *gin.Context) {
			handler.Save(c)
		})

		// Criar um payload válido
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

		// Configurar o mock para retornar um erro
		expectedError := errors.New("database error")
		mockService.On("CreateCashHandlingEntry", mock.AnythingOfType("dtos.CreateCashHandlingEntryDTO")).Return(dtos.CashHandlingEntryResponseDTO{}, expectedError)

		// Converter o payload para JSON
		jsonPayload, _ := json.Marshal(validEntry)

		// Criar uma requisição HTTP de teste
		req, _ := http.NewRequest("POST", "/cash-handling", bytes.NewBuffer(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Executar a requisição
		router.ServeHTTP(w, req)

		// Verificar o resultado
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Verificar o corpo da resposta
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedError.Error(), response["error"])

		// Verificar que o mock foi chamado conforme esperado
		mockService.AssertExpectations(t)
	})

	t.Run("invalid_request_body", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		// Configurar a rota para o teste
		router.POST("/cash-handling", func(c *gin.Context) {
			handler.Save(c)
		})

		// Criar um payload inválido (JSON malformado)
		invalidJSON := []byte(`{"amount": "invalid", "title": 123}`)

		// Criar uma requisição HTTP de teste
		req, _ := http.NewRequest("POST", "/cash-handling", bytes.NewBuffer(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		// Executar a requisição
		router.ServeHTTP(w, req)

		// Verificar o resultado - deve ser um erro de validação
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// O serviço não deve ser chamado com dados inválidos
		mockService.AssertNotCalled(t, "CreateCashHandlingEntry")
	})
}

// TestGetAllHandler testa o método GetAll do handler
func TestGetAllHandler(t *testing.T) {
	// Testes para o método GetAll
	t.Run("successful_retrieval", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		// Configurar a rota para o teste
		router.GET("/cash-handling", func(c *gin.Context) {
			handler.GetAll(c)
		})

		// Configurar o mock para retornar uma resposta de sucesso
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

		// Criar uma requisição HTTP de teste
		req, _ := http.NewRequest("GET", "/cash-handling?limit=10&skip=0", nil)
		w := httptest.NewRecorder()

		// Executar a requisição
		router.ServeHTTP(w, req)

		// Verificar o resultado
		assert.Equal(t, http.StatusOK, w.Code)

		// Verificar o corpo da resposta
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verificar a estrutura da resposta
		data, ok := response["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 2)

		pagination, ok := response["pagination"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(10), pagination["limit"])
		assert.Equal(t, float64(0), pagination["skip"])
		assert.Equal(t, float64(2), pagination["count"])

		// Verificar que o mock foi chamado conforme esperado
		mockService.AssertExpectations(t)
	})

	t.Run("with_filters", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		// Configurar a rota para o teste
		router.GET("/cash-handling", func(c *gin.Context) {
			handler.GetAll(c)
		})

		// Configurar o mock para retornar uma resposta filtrada
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

		// Criar uma requisição HTTP de teste com filtros
		req, _ := http.NewRequest("GET", "/cash-handling?limit=10&skip=0&title=lunch&category=food", nil)
		w := httptest.NewRecorder()

		// Executar a requisição
		router.ServeHTTP(w, req)

		// Verificar o resultado
		assert.Equal(t, http.StatusOK, w.Code)

		// Verificar o corpo da resposta
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Verificar a estrutura da resposta
		data, ok := response["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 1)

		filters, ok := response["filters"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, "lunch", filters["title"])
		assert.Equal(t, "food", filters["category"])

		// Verificar que o mock foi chamado conforme esperado
		mockService.AssertExpectations(t)
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		// Configurar a rota para o teste
		router.GET("/cash-handling", func(c *gin.Context) {
			handler.GetAll(c)
		})

		// Configurar o mock para retornar um erro
		expectedError := errors.New("database error")
		mockService.On("GetAllCashHandlingEntries", 10, 0, "", "").Return([]dtos.CashHandlingEntryResponseDTO{}, expectedError)

		// Criar uma requisição HTTP de teste
		req, _ := http.NewRequest("GET", "/cash-handling?limit=10&skip=0", nil)
		w := httptest.NewRecorder()

		// Executar a requisição
		router.ServeHTTP(w, req)

		// Verificar o resultado
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Verificar o corpo da resposta
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to retrieve entries", response["error"])
		assert.Equal(t, expectedError.Error(), response["details"])

		// Verificar que o mock foi chamado conforme esperado
		mockService.AssertExpectations(t)
	})

	t.Run("invalid_pagination_params", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		// Configurar a rota para o teste
		router.GET("/cash-handling", func(c *gin.Context) {
			handler.GetAll(c)
		})

		// Criar uma requisição HTTP de teste com parâmetros inválidos
		req, _ := http.NewRequest("GET", "/cash-handling?limit=invalid&skip=abc", nil)
		w := httptest.NewRecorder()

		// Executar a requisição
		router.ServeHTTP(w, req)

		// Verificar o resultado - deve ser um erro de validação
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// O serviço não deve ser chamado com parâmetros inválidos
		mockService.AssertNotCalled(t, "GetAllCashHandlingEntries")
	})
}

// TestDeleteHandler testa o método Delete do handler
func TestDeleteHandler(t *testing.T) {
	// Testes para o método Delete
	t.Run("successful_deletion", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		// Configurar a rota para o teste
		router.DELETE("/cash-handling/:id", func(c *gin.Context) {
			handler.Delete(c)
		})

		// ID válido para teste
		validID := "123456789012345678901234"

		// Configurar o mock para retornar sucesso
		mockService.On("DeleteCashHandlingEntry", validID).Return(nil)

		// Criar uma requisição HTTP de teste
		req, _ := http.NewRequest("DELETE", "/cash-handling/"+validID, nil)
		w := httptest.NewRecorder()

		// Executar a requisição
		router.ServeHTTP(w, req)

		// Verificar o resultado
		assert.Equal(t, http.StatusOK, w.Code)

		// Verificar o corpo da resposta
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Entry deleted successfully", response["message"])
		assert.Equal(t, validID, response["id"])

		// Verificar que o mock foi chamado conforme esperado
		mockService.AssertExpectations(t)
	})

	t.Run("missing_id", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		// Configurar a rota para o teste
		router.DELETE("/cash-handling/:id", func(c *gin.Context) {
			handler.Delete(c)
		})

		// Criar uma requisição HTTP de teste sem ID
		req, _ := http.NewRequest("DELETE", "/cash-handling/", nil)
		w := httptest.NewRecorder()

		// Executar a requisição
		router.ServeHTTP(w, req)

		// Verificar o resultado - deve ser um erro 404 Not Found (rota não encontrada)
		assert.Equal(t, http.StatusNotFound, w.Code)

		// O serviço não deve ser chamado sem ID
		mockService.AssertNotCalled(t, "DeleteCashHandlingEntry")
	})

	t.Run("empty_id", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		// Configurar uma rota especial para testar ID vazio
		router.DELETE("/cash-handling/", func(c *gin.Context) {
			c.Params = append(c.Params, gin.Param{Key: "id", Value: ""})
			handler.Delete(c)
		})

		// Criar uma requisição HTTP de teste com ID vazio
		req, _ := http.NewRequest("DELETE", "/cash-handling/", nil)
		w := httptest.NewRecorder()

		// Executar a requisição
		router.ServeHTTP(w, req)

		// Verificar o resultado - deve ser um erro de validação
		assert.Equal(t, http.StatusBadRequest, w.Code)

		// Verificar o corpo da resposta
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "ID is required", response["error"])

		// O serviço não deve ser chamado com ID vazio
		mockService.AssertNotCalled(t, "DeleteCashHandlingEntry")
	})

	t.Run("service_error", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		// Configurar a rota para o teste
		router.DELETE("/cash-handling/:id", func(c *gin.Context) {
			handler.Delete(c)
		})

		// ID válido para teste
		validID := "123456789012345678901234"

		// Configurar o mock para retornar um erro
		expectedError := errors.New("database error")
		mockService.On("DeleteCashHandlingEntry", validID).Return(expectedError)

		// Criar uma requisição HTTP de teste
		req, _ := http.NewRequest("DELETE", "/cash-handling/"+validID, nil)
		w := httptest.NewRecorder()

		// Executar a requisição
		router.ServeHTTP(w, req)

		// Verificar o resultado
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Verificar o corpo da resposta
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to delete entry", response["error"])
		assert.Equal(t, expectedError.Error(), response["details"])

		// Verificar que o mock foi chamado conforme esperado
		mockService.AssertExpectations(t)
	})

	t.Run("invalid_id_format", func(t *testing.T) {
		mockService := new(MockCashHandlingService)
		handler := NewCashHandlingHandler(mockService)
		router := setupRouter()

		// Configurar a rota para o teste
		router.DELETE("/cash-handling/:id", func(c *gin.Context) {
			handler.Delete(c)
		})

		// ID inválido para teste
		invalidID := "invalid-id-format"

		// Configurar o mock para retornar um erro específico de formato inválido
		expectedError := errors.New("the provided hex string is not a valid ObjectID")
		mockService.On("DeleteCashHandlingEntry", invalidID).Return(expectedError)

		// Criar uma requisição HTTP de teste
		req, _ := http.NewRequest("DELETE", "/cash-handling/"+invalidID, nil)
		w := httptest.NewRecorder()

		// Executar a requisição
		router.ServeHTTP(w, req)

		// Verificar o resultado
		assert.Equal(t, http.StatusInternalServerError, w.Code)

		// Verificar o corpo da resposta
		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to delete entry", response["error"])
		assert.Equal(t, expectedError.Error(), response["details"])

		// Verificar que o mock foi chamado conforme esperado
		mockService.AssertExpectations(t)
	})
}

