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

func TestCashHandlingServiceCreateCashHandlingEntrySuccess(t *testing.T) {
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

func TestCashHandlingServiceCreateCashHandlingEntryInvalidDate(t *testing.T) {
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

func TestCashHandlingServiceCreateCashHandlingEntryRepositoryError(t *testing.T) {
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

func TestCashHandlingServiceCreateCashHandlingEntryEmptyDescription(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	inputDTO := dtos.CreateCashHandlingEntryDTO{
		Amount:        100.0,
		Currency:      "USD",
		Type:          "income",
		Category:      "salary",
		PaymentMethod: "bank_transfer",
		Description:   "", // Empty description
		Date:          "01/09/2025",
	}

	expectedDate, _ := time.Parse("02/01/2006", "01/09/2025")
	objectID := primitive.NewObjectID()
	createdTime := time.Now().UTC()

	expectedModel := &model.CashHandlingEntryModel{
		ID:            objectID,
		Amount:        100.0,
		Currency:      "USD",
		Type:          "income",
		Category:      "salary",
		PaymentMethod: "bank_transfer",
		Description:   "",
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
	assert.Equal(t, "", result.Description)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceCreateCashHandlingEntryDateParsing(t *testing.T) {
	// Test different date formats - should only work with DD/MM/YYYY
	testCases := []struct {
		name        string
		inputDate   string
		shouldError bool
	}{
		{"Valid DD/MM/YYYY", "15/03/2025", false},
		{"Invalid format YYYY-MM-DD", "2025-03-15", true},
		{"Invalid format MM/DD/YYYY", "03/15/2025", true},
		{"Empty date", "", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockCashHandlingRepository)
			service := NewCashHandlingService(mockRepo)

			inputDTO := dtos.CreateCashHandlingEntryDTO{
				Amount:        50.0,
				Currency:      "EUR",
				Type:          "expense",
				Category:      "transport",
				PaymentMethod: "cash",
				Description:   "Bus ticket",
				Date:          tc.inputDate,
			}

			if !tc.shouldError {
				objectID := primitive.NewObjectID()
				createdTime := time.Now().UTC()
				expectedModel := &model.CashHandlingEntryModel{
					ID:            objectID,
					Amount:        50.0,
					Currency:      "EUR",
					Type:          "expense",
					Category:      "transport",
					PaymentMethod: "cash",
					Description:   "Bus ticket",
					Date:          time.Time{}, // Will be set by parser
					Timestamp:     createdTime.Unix(),
					CreatedAt:     createdTime,
					UpdatedAt:     createdTime,
				}
				mockRepo.On("Create", mock.AnythingOfType("*model.CashHandlingEntryModel")).Return(expectedModel, nil)
			}

			// Act
			result, err := service.CreateCashHandlingEntry(inputDTO)

			// Assert
			if tc.shouldError {
				assert.Error(t, err)
				assert.Equal(t, dtos.CashHandlingEntryResponseDTO{}, result)
				mockRepo.AssertNotCalled(t, "Create")
			} else {
				assert.NoError(t, err)
				assert.NotEqual(t, dtos.CashHandlingEntryResponseDTO{}, result)
				mockRepo.AssertExpectations(t)
			}
		})
	}
}
