package services

import (
	"errors"
	"testing"
	"time"

	"myfin-api/internal/dtos"
	"myfin-api/internal/model"
	"myfin-api/internal/repository/types"

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

func (m *MockCashHandlingRepository) GetAllWithFilter(limit, skip int, filter types.FilterOptions) ([]*model.CashHandlingEntryModel, error) {
	args := m.Called(limit, skip, filter)
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
		Title:         "Lunch at restaurant",
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
		Title:         "Lunch at restaurant",
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
	assert.Equal(t, "Lunch at restaurant", result.Title)
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
		Title:         "Lunch at restaurant",
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
		Title:         "Lunch at restaurant",
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
			Title:         "Lunch at restaurant",
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
			Title:         "Lunch at restaurant",
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
	result, err := service.GetAllCashHandlingEntries(10, 0, "", "")

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
	assert.Equal(t, "Lunch at restaurant", result[0].Title)
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
	result, err := service.GetAllCashHandlingEntries(10, 0, "", "")

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
	_, err := service.GetAllCashHandlingEntries(5, 10, "", "")

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
	result, err := service.GetAllCashHandlingEntries(1, 5, "", "")

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
			Title:         "Electric bill",
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
	result, err := service.GetAllCashHandlingEntries(0, 0, "", "")

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
			Title:         "Electric bill",
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
	result, err := service.GetAllCashHandlingEntries(10, 0, "", "")

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
	result, err := service.GetAllCashHandlingEntries(10, 0, "", "")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllParameterValidation(t *testing.T) {
	// Test cases for parameter validation logic
	testCases := []struct {
		name          string
		inputLimit    int
		inputSkip     int
		expectedLimit int
		expectedSkip  int
		description   string
	}{
		{
			name:          "negative_limit_defaults_to_10",
			inputLimit:    -5,
			inputSkip:     0,
			expectedLimit: 10,
			expectedSkip:  0,
			description:   "Negative limit should default to 10",
		},
		{
			name:          "limit_over_100_capped_to_100",
			inputLimit:    150,
			inputSkip:     0,
			expectedLimit: 100,
			expectedSkip:  0,
			description:   "Limit over 100 should be capped to 100",
		},
		{
			name:          "negative_skip_defaults_to_0",
			inputLimit:    10,
			inputSkip:     -3,
			expectedLimit: 10,
			expectedSkip:  0,
			description:   "Negative skip should default to 0",
		},
		{
			name:          "valid_parameters_unchanged",
			inputLimit:    20,
			inputSkip:     5,
			expectedLimit: 20,
			expectedSkip:  5,
			description:   "Valid parameters should remain unchanged",
		},
		{
			name:          "zero_limit_unchanged",
			inputLimit:    0,
			inputSkip:     10,
			expectedLimit: 0,
			expectedSkip:  10,
			description:   "Zero limit should remain unchanged (no limit)",
		},
		{
			name:          "edge_case_limit_exactly_100",
			inputLimit:    100,
			inputSkip:     0,
			expectedLimit: 100,
			expectedSkip:  0,
			description:   "Limit exactly 100 should remain unchanged",
		},
		{
			name:          "multiple_validations_applied",
			inputLimit:    -10,
			inputSkip:     -5,
			expectedLimit: 10,
			expectedSkip:  0,
			description:   "Both negative limit and skip should be corrected",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := new(MockCashHandlingRepository)
			service := NewCashHandlingService(mockRepo)

			objectID := primitive.NewObjectID()
			createdTime := time.Date(2025, 9, 7, 10, 0, 0, 0, time.UTC)

			mockEntries := []*model.CashHandlingEntryModel{
				{
					ID:            objectID,
					Amount:        100.0,
					Title:         "Test Entry",
					Currency:      "BRL",
					Type:          "expense",
					Category:      "test",
					PaymentMethod: "cash",
					Description:   "Test description",
					Date:          createdTime,
					Timestamp:     createdTime.Unix(),
					CreatedAt:     createdTime,
					UpdatedAt:     createdTime,
				},
			}

			// Mock should be called with the expected (validated) parameters
			mockRepo.On("GetAll", tc.expectedLimit, tc.expectedSkip).Return(mockEntries, nil)

			// Act
			result, err := service.GetAllCashHandlingEntries(tc.inputLimit, tc.inputSkip, "", "")

			// Assert
			assert.NoError(t, err, tc.description)
			assert.Len(t, result, 1, tc.description)
			assert.Equal(t, objectID.Hex(), result[0].ID, tc.description)

			// Verify that the mock was called with the expected parameters
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCashHandlingServiceGetAllBoundaryValues(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()
	createdTime := time.Now().UTC()

	mockEntries := []*model.CashHandlingEntryModel{
		{
			ID:            objectID,
			Amount:        50.0,
			Title:         "Boundary Test",
			Currency:      "USD",
			Type:          "income",
			Category:      "test",
			PaymentMethod: "bank_transfer",
			Description:   "Boundary value test",
			Date:          createdTime,
			Timestamp:     createdTime.Unix(),
			CreatedAt:     createdTime,
			UpdatedAt:     createdTime,
		},
	}

	// Test boundary values
	testCases := []struct {
		name          string
		limit         int
		skip          int
		expectedLimit int
		expectedSkip  int
	}{
		{"limit_101_capped_to_100", 101, 0, 100, 0},
		{"limit_99_unchanged", 99, 0, 99, 0},
		{"limit_1_unchanged", 1, 0, 1, 0},
		{"skip_large_positive_unchanged", 50, 1000, 50, 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.On("GetAll", tc.expectedLimit, tc.expectedSkip).Return(mockEntries, nil).Once()

			// Act
			result, err := service.GetAllCashHandlingEntries(tc.limit, tc.skip, "", "")

			// Assert
			assert.NoError(t, err)
			assert.Len(t, result, 1)
		})
	}

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceParameterValidationWithError(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	expectedError := errors.New("repository error after parameter validation")

	// Even with invalid input parameters, they should be validated before repository call
	// Input: limit=-10, skip=-5
	// Expected after validation: limit=10, skip=0
	mockRepo.On("GetAll", 10, 0).Return(nil, expectedError)

	// Act
	result, err := service.GetAllCashHandlingEntries(-10, -5, "", "")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	// Verify that repository was called with validated parameters, not original ones
	mockRepo.AssertExpectations(t)
}

// Tests for GetAllCashHandlingEntries with filters
func TestCashHandlingServiceGetAllWithFilter(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID1 := primitive.NewObjectID()
	objectID2 := primitive.NewObjectID()
	createdTime := time.Date(2025, 9, 7, 10, 0, 0, 0, time.UTC)

	// Mock entries that match filter criteria
	mockEntries := []*model.CashHandlingEntryModel{
		{
			ID:            objectID1,
			Amount:        150.75,
			Title:         "Lunch Restaurant",
			Currency:      "BRL",
			Type:          "expense",
			Category:      "food",
			PaymentMethod: "credit_card",
			Description:   "Restaurant meal",
			Date:          createdTime,
			Timestamp:     createdTime.Unix(),
			CreatedAt:     createdTime,
			UpdatedAt:     createdTime,
		},
		{
			ID:            objectID2,
			Amount:        25.50,
			Title:         "Fast food lunch",
			Currency:      "BRL",
			Type:          "expense",
			Category:      "food",
			PaymentMethod: "cash",
			Description:   "Quick lunch",
			Date:          createdTime,
			Timestamp:     createdTime.Unix(),
			CreatedAt:     createdTime,
			UpdatedAt:     createdTime,
		},
	}

	expectedFilter := types.FilterOptions{
		Title:    "lunch",
		Category: "food",
	}

	mockRepo.On("GetAllWithFilter", 10, 0, expectedFilter).Return(mockEntries, nil)

	// Act
	result, err := service.GetAllCashHandlingEntries(10, 0, "lunch", "food")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	// Verify first entry
	assert.Equal(t, objectID1.Hex(), result[0].ID)
	assert.Equal(t, "Lunch Restaurant", result[0].Title)
	assert.Equal(t, "food", result[0].Category)

	// Verify second entry
	assert.Equal(t, objectID2.Hex(), result[1].ID)
	assert.Equal(t, "Fast food lunch", result[1].Title)
	assert.Equal(t, "food", result[1].Category)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllWithFilterTitleOnly(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()
	createdTime := time.Date(2025, 9, 7, 10, 0, 0, 0, time.UTC)

	mockEntries := []*model.CashHandlingEntryModel{
		{
			ID:            objectID,
			Amount:        100.0,
			Title:         "Coffee Shop",
			Currency:      "BRL",
			Type:          "expense",
			Category:      "food",
			PaymentMethod: "credit_card",
			Description:   "Morning coffee",
			Date:          createdTime,
			Timestamp:     createdTime.Unix(),
			CreatedAt:     createdTime,
			UpdatedAt:     createdTime,
		},
	}

	expectedFilter := types.FilterOptions{
		Title:    "coffee",
		Category: "", // Empty category filter
	}

	mockRepo.On("GetAllWithFilter", 5, 2, expectedFilter).Return(mockEntries, nil)

	// Act
	result, err := service.GetAllCashHandlingEntries(5, 2, "coffee", "")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Coffee Shop", result[0].Title)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllWithFilterCategoryOnly(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()
	createdTime := time.Date(2025, 9, 7, 10, 0, 0, 0, time.UTC)

	mockEntries := []*model.CashHandlingEntryModel{
		{
			ID:            objectID,
			Amount:        50.0,
			Title:         "Bus Ticket",
			Currency:      "BRL",
			Type:          "expense",
			Category:      "transport",
			PaymentMethod: "cash",
			Description:   "Public transport",
			Date:          createdTime,
			Timestamp:     createdTime.Unix(),
			CreatedAt:     createdTime,
			UpdatedAt:     createdTime,
		},
	}

	expectedFilter := types.FilterOptions{
		Title:    "", // Empty title filter
		Category: "transport",
	}

	mockRepo.On("GetAllWithFilter", 10, 0, expectedFilter).Return(mockEntries, nil)

	// Act
	result, err := service.GetAllCashHandlingEntries(10, 0, "", "transport")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "transport", result[0].Category)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllWithFilterParameterValidation(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()
	createdTime := time.Date(2025, 9, 7, 10, 0, 0, 0, time.UTC)

	mockEntries := []*model.CashHandlingEntryModel{
		{
			ID:            objectID,
			Amount:        200.0,
			Title:         "Test Entry",
			Currency:      "USD",
			Type:          "income",
			Category:      "salary",
			PaymentMethod: "bank_transfer",
			Description:   "Test income",
			Date:          createdTime,
			Timestamp:     createdTime.Unix(),
			CreatedAt:     createdTime,
			UpdatedAt:     createdTime,
		},
	}

	expectedFilter := types.FilterOptions{
		Title:    "test",
		Category: "salary",
	}

	// Parameters should be validated: limit=-5 -> 10, skip=-3 -> 0
	mockRepo.On("GetAllWithFilter", 10, 0, expectedFilter).Return(mockEntries, nil)

	// Act
	result, err := service.GetAllCashHandlingEntries(-5, -3, "test", "salary")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, result, 1)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllWithFilterRepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	expectedError := errors.New("database filter query failed")
	expectedFilter := types.FilterOptions{
		Title:    "error",
		Category: "test",
	}

	mockRepo.On("GetAllWithFilter", 10, 0, expectedFilter).Return(nil, expectedError)

	// Act
	result, err := service.GetAllCashHandlingEntries(10, 0, "error", "test")

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllWithFilterEmptyResult(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	expectedFilter := types.FilterOptions{
		Title:    "nonexistent",
		Category: "unknown",
	}

	// Repository returns empty slice for filters that don't match anything
	mockRepo.On("GetAllWithFilter", 10, 0, expectedFilter).Return([]*model.CashHandlingEntryModel{}, nil)

	// Act
	result, err := service.GetAllCashHandlingEntries(10, 0, "nonexistent", "unknown")

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.NotNil(t, result) // Should return empty slice, not nil

	mockRepo.AssertExpectations(t)
}
