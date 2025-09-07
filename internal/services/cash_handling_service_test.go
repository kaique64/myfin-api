package services

import (
	"errors"
	"testing"
	"time"

	"myfin-api/internal/dtos"
	"myfin-api/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Mock repository
type MockCashHandlingRepository struct {
	mock.Mock
}

func (m *MockCashHandlingRepository) Create(entry *model.CashHandlingEntryModel) (*model.CashHandlingEntryModel, error) {
	args := m.Called(entry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CashHandlingEntryModel), args.Error(1)
}

func (m *MockCashHandlingRepository) GetAll(limit, skip int) ([]*model.CashHandlingEntryModel, error) {
	args := m.Called(limit, skip)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.CashHandlingEntryModel), args.Error(1)
}

// Tests for CreateCashHandlingEntry
func TestCashHandlingService_CreateCashHandlingEntry_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	inputDTO := dtos.CreateCashHandlingEntryDTO{
		Amount:        150.75,
		Currency:      "BRL",
		Type:          "expense",
		Category:      "food",
		PaymentMethod: "credit_card",
		Description:   "Lunch at restaurant",
		Date:          "06/09/2025",
	}

	expectedDate, _ := time.Parse("02/01/2006", "06/09/2025")
	objectID := primitive.NewObjectID()
	createdTime := time.Now().UTC()

	expectedModel := &model.CashHandlingEntryModel{
		ID:            objectID,
		Amount:        150.75,
		Currency:      "BRL",
		Type:          "expense",
		Category:      "food",
		PaymentMethod: "credit_card",
		Description:   "Lunch at restaurant",
		Date:          expectedDate,
		Timestamp:     createdTime.Unix(),
		CreatedAt:     createdTime,
		UpdatedAt:     createdTime,
	}

	mockRepo.On("Create", mock.AnythingOfType("*model.CashHandlingEntryModel")).Return(expectedModel, nil)

	// Act
	result, err := service.CreateCashHandlingEntry(inputDTO)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, objectID.Hex(), result.ID)
	assert.Equal(t, float64(150.75), result.Amount)
	assert.Equal(t, "BRL", result.Currency)
	assert.Equal(t, "expense", result.Type)
	assert.Equal(t, "food", result.Category)
	assert.Equal(t, "credit_card", result.PaymentMethod)
	assert.Equal(t, "Lunch at restaurant", result.Description)
	assert.Equal(t, "06/09/2025", result.Date)
	assert.Equal(t, createdTime.Unix(), result.Timestamp)
	assert.Equal(t, createdTime.UTC().Format(time.RFC3339), result.CreatedAt)
	assert.Equal(t, createdTime.UTC().Format(time.RFC3339), result.UpdatedAt)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingService_CreateCashHandlingEntry_InvalidDate(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	inputDTO := dtos.CreateCashHandlingEntryDTO{
		Amount:        150.75,
		Currency:      "BRL",
		Type:          "expense",
		Category:      "food",
		PaymentMethod: "credit_card",
		Description:   "Lunch at restaurant",
		Date:          "invalid-date",
	}

	// Act
	result, err := service.CreateCashHandlingEntry(inputDTO)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, dtos.CashHandlingEntryResponseDTO{}, result)

	// Repository should not be called
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCashHandlingService_CreateCashHandlingEntry_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	inputDTO := dtos.CreateCashHandlingEntryDTO{
		Amount:        150.75,
		Currency:      "BRL",
		Type:          "expense",
		Category:      "food",
		PaymentMethod: "credit_card",
		Description:   "Lunch at restaurant",
		Date:          "06/09/2025",
	}

	expectedError := errors.New("database connection failed")
	mockRepo.On("Create", mock.AnythingOfType("*model.CashHandlingEntryModel")).Return(nil, expectedError)

	// Act
	result, err := service.CreateCashHandlingEntry(inputDTO)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dtos.CashHandlingEntryResponseDTO{}, result)

	mockRepo.AssertExpectations(t)
}

// Tests for GetAllCashHandlingEntries
func TestCashHandlingService_GetAllCashHandlingEntries_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID1 := primitive.NewObjectID()
	objectID2 := primitive.NewObjectID()
	createdTime1 := time.Date(2025, 9, 6, 14, 30, 0, 0, time.UTC)
	createdTime2 := time.Date(2025, 8, 30, 9, 0, 0, 0, time.UTC)

	mockEntries := []*model.CashHandlingEntryModel{
		{
			ID:            objectID1,
			Amount:        150.75,
			Currency:      "BRL",
			Type:          "expense",
			Category:      "food",
			PaymentMethod: "credit_card",
			Description:   "Lunch at restaurant",
			Date:          createdTime1,
			Timestamp:     createdTime1.Unix(),
			CreatedAt:     createdTime1,
			UpdatedAt:     createdTime1,
		},
		{
			ID:            objectID2,
			Amount:        2500.0,
			Currency:      "BRL",
			Type:          "income",
			Category:      "salary",
			PaymentMethod: "bank_transfer",
			Description:   "Monthly salary",
			Date:          createdTime2,
			Timestamp:     createdTime2.Unix(),
			CreatedAt:     createdTime2,
			UpdatedAt:     createdTime2,
		},
	}

	mockRepo.On("GetAll", 10, 0).Return(mockEntries, nil)

	// Act
	result, err := service.GetAllCashHandlingEntries(10, 0)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify first entry
	assert.Equal(t, objectID1.Hex(), result[0].ID)
	assert.Equal(t, float64(150.75), result[0].Amount)
	assert.Equal(t, "BRL", result[0].Currency)
	assert.Equal(t, "expense", result[0].Type)
	assert.Equal(t, "food", result[0].Category)
	assert.Equal(t, "credit_card", result[0].PaymentMethod)
	assert.Equal(t, "Lunch at restaurant", result[0].Description)
	assert.Equal(t, createdTime1.Format("02/01/2006"), result[0].Date)
	assert.Equal(t, createdTime1.Unix(), result[0].Timestamp)
	assert.Equal(t, createdTime1.UTC().Format(time.RFC3339), result[0].CreatedAt)
	assert.Equal(t, createdTime1.UTC().Format(time.RFC3339), result[0].UpdatedAt)

	// Verify second entry
	assert.Equal(t, objectID2.Hex(), result[1].ID)
	assert.Equal(t, float64(2500.0), result[1].Amount)
	assert.Equal(t, "income", result[1].Type)
	assert.Equal(t, "salary", result[1].Category)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingService_GetAllCashHandlingEntries_EmptyResult(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	mockRepo.On("GetAll", 10, 0).Return([]*model.CashHandlingEntryModel{}, nil)

	// Act
	result, err := service.GetAllCashHandlingEntries(10, 0)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingService_GetAllCashHandlingEntries_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	expectedError := errors.New("database connection failed")
	mockRepo.On("GetAll", 5, 10).Return(nil, expectedError)

	// Act
	_, err := service.GetAllCashHandlingEntries(5, 10)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingService_GetAllCashHandlingEntries_WithPagination(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()
	createdTime := time.Date(2025, 8, 28, 20, 15, 0, 0, time.UTC)

	mockEntries := []*model.CashHandlingEntryModel{
		{
			ID:            objectID,
			Amount:        89.99,
			Currency:      "BRL",
			Type:          "expense",
			Category:      "entertainment",
			PaymentMethod: "credit_card",
			Description:   "Netflix subscription",
			Date:          createdTime,
			Timestamp:     createdTime.Unix(),
			CreatedAt:     createdTime,
			UpdatedAt:     createdTime,
		},
	}

	mockRepo.On("GetAll", 1, 5).Return(mockEntries, nil)

	// Act
	result, err := service.GetAllCashHandlingEntries(1, 5)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, objectID.Hex(), result[0].ID)
	assert.Equal(t, float64(89.99), result[0].Amount)
	assert.Equal(t, "entertainment", result[0].Category)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingService_GetAllCashHandlingEntries_NoPagination(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()
	createdTime := time.Date(2025, 8, 27, 14, 22, 0, 0, time.UTC)

	mockEntries := []*model.CashHandlingEntryModel{
		{
			ID:            objectID,
			Amount:        320.50,
			Currency:      "BRL",
			Type:          "expense",
			Category:      "utilities",
			PaymentMethod: "boleto",
			Description:   "Electric bill",
			Date:          createdTime,
			Timestamp:     createdTime.Unix(),
			CreatedAt:     createdTime,
			UpdatedAt:     createdTime,
		},
	}

	// Test with 0 limit and 0 skip (no pagination)
	mockRepo.On("GetAll", 0, 0).Return(mockEntries, nil)

	// Act
	result, err := service.GetAllCashHandlingEntries(0, 0)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, float64(320.50), result[0].Amount)
	assert.Equal(t, "utilities", result[0].Category)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingService_GetAllCashHandlingEntries_DateFormatting(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()
	// Specific date to test formatting: March 15, 2025
	testDate := time.Date(2025, 3, 15, 10, 30, 0, 0, time.UTC)

	mockEntries := []*model.CashHandlingEntryModel{
		{
			ID:            objectID,
			Amount:        100.0,
			Currency:      "USD",
			Type:          "income",
			Category:      "freelance",
			PaymentMethod: "pix",
			Description:   "Project payment",
			Date:          testDate,
			Timestamp:     testDate.Unix(),
			CreatedAt:     testDate,
			UpdatedAt:     testDate,
		},
	}

	mockRepo.On("GetAll", 10, 0).Return(mockEntries, nil)

	// Act
	result, err := service.GetAllCashHandlingEntries(10, 0)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	// Test date formatting - should be in DD/MM/YYYY format
	assert.Equal(t, "15/03/2025", result[0].Date)

	// Test timestamp formatting - should be in RFC3339 format
	assert.Equal(t, testDate.UTC().Format(time.RFC3339), result[0].CreatedAt)
	assert.Equal(t, testDate.UTC().Format(time.RFC3339), result[0].UpdatedAt)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingService_GetAllCashHandlingEntries_NilEntries(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	// Mock repository returns nil (unusual case)
	mockRepo.On("GetAll", 10, 0).Return(nil, nil)

	// Act
	result, err := service.GetAllCashHandlingEntries(10, 0)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}
