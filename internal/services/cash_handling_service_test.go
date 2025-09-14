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

func (m *MockCashHandlingRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCashHandlingRepository) GetByID(id string) (*model.CashHandlingEntryModel, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CashHandlingEntryModel), args.Error(1)
}

func (m *MockCashHandlingRepository) Update(id string, entry *model.CashHandlingEntryModel) (*model.CashHandlingEntryModel, error) {
	args := m.Called(id, entry)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.CashHandlingEntryModel), args.Error(1)
}

func TestCashHandlingServiceDeleteCashHandlingEntrySuccess(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()

	mockRepo.On("Delete", objectID.Hex()).Return(nil)

	err := service.DeleteCashHandlingEntry(objectID.Hex())

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceDeleteCashHandlingEntryRepositoryError(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()

	expectedError := errors.New("database error")
	mockRepo.On("Delete", objectID.Hex()).Return(expectedError)

	err := service.DeleteCashHandlingEntry(objectID.Hex())

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceCreateCashHandlingEntrySuccess(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	inputDTO := dtos.CreateCashHandlingEntryDTO{
		Amount:        150.75,
		Title:         "Lunch at restaurant",
		Currency:      "BRL",
		Type:          "expense",
		Category:      "food",
		PaymentMethod: "creditcard",
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
		PaymentMethod: "creditcard",
		Description:   "Lunch at restaurant",
		Date:          expectedDate,
		Timestamp:     createdTime.Unix(),
		CreatedAt:     createdTime,
		UpdatedAt:     createdTime,
	}

	mockRepo.On("Create", mock.AnythingOfType("*model.CashHandlingEntryModel")).Return(expectedModel, nil)

	result, err := service.CreateCashHandlingEntry(inputDTO)

	assert.NoError(t, err)
	assert.Equal(t, objectID.Hex(), result.ID)
	assert.Equal(t, float64(150.75), result.Amount)
	assert.Equal(t, "BRL", result.Currency)
	assert.Equal(t, "expense", result.Type)
	assert.Equal(t, "food", result.Category)
	assert.Equal(t, "creditcard", result.PaymentMethod)
	assert.Equal(t, "Lunch at restaurant", result.Title)
	assert.Equal(t, "Lunch at restaurant", result.Description)
	assert.Equal(t, "06/09/2025", result.Date)
	assert.Equal(t, createdTime.Unix(), result.Timestamp)
	assert.Equal(t, createdTime.UTC().Format(time.RFC3339), result.CreatedAt)
	assert.Equal(t, createdTime.UTC().Format(time.RFC3339), result.UpdatedAt)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceCreateCashHandlingEntryInvalidDate(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	inputDTO := dtos.CreateCashHandlingEntryDTO{
		Amount:        150.75,
		Title:         "Lunch at restaurant",
		Currency:      "BRL",
		Type:          "expense",
		Category:      "food",
		PaymentMethod: "creditcard",
		Description:   "Lunch at restaurant",
		Date:          "invalid-date",
	}

	result, err := service.CreateCashHandlingEntry(inputDTO)

	assert.Error(t, err)
	assert.Equal(t, dtos.CashHandlingEntryResponseDTO{}, result)

	mockRepo.AssertNotCalled(t, "Create")
}

func TestCashHandlingServiceCreateCashHandlingEntryRepositoryError(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	inputDTO := dtos.CreateCashHandlingEntryDTO{
		Amount:        150.75,
		Title:         "Lunch at restaurant",
		Currency:      "BRL",
		Type:          "expense",
		Category:      "food",
		PaymentMethod: "creditcard",
		Description:   "Lunch at restaurant",
		Date:          "06/09/2025",
	}

	expectedError := errors.New("database connection failed")
	mockRepo.On("Create", mock.AnythingOfType("*model.CashHandlingEntryModel")).Return(nil, expectedError)

	result, err := service.CreateCashHandlingEntry(inputDTO)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dtos.CashHandlingEntryResponseDTO{}, result)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllCashHandlingEntriesSuccess(t *testing.T) {
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
			PaymentMethod: "creditcard",
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
			PaymentMethod: "banktransfer",
			Description:   "Monthly salary",
			Date:          createdTime2,
			Timestamp:     createdTime2.Unix(),
			CreatedAt:     createdTime2,
			UpdatedAt:     createdTime2,
		},
	}

	mockRepo.On("GetAll", 10, 0).Return(mockEntries, nil)

	result, err := service.GetAllCashHandlingEntries(10, 0, "", "")

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	assert.Equal(t, objectID1.Hex(), result[0].ID)
	assert.Equal(t, float64(150.75), result[0].Amount)
	assert.Equal(t, "BRL", result[0].Currency)
	assert.Equal(t, "expense", result[0].Type)
	assert.Equal(t, "food", result[0].Category)
	assert.Equal(t, "creditcard", result[0].PaymentMethod)
	assert.Equal(t, "Lunch at restaurant", result[0].Title)
	assert.Equal(t, "Lunch at restaurant", result[0].Description)
	assert.Equal(t, createdTime1.Format("02/01/2006"), result[0].Date)
	assert.Equal(t, createdTime1.Unix(), result[0].Timestamp)
	assert.Equal(t, createdTime1.UTC().Format(time.RFC3339), result[0].CreatedAt)
	assert.Equal(t, createdTime1.UTC().Format(time.RFC3339), result[0].UpdatedAt)

	assert.Equal(t, objectID2.Hex(), result[1].ID)
	assert.Equal(t, float64(2500.0), result[1].Amount)
	assert.Equal(t, "income", result[1].Type)
	assert.Equal(t, "salary", result[1].Category)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllCashHandlingEntriesEmptyResult(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	mockRepo.On("GetAll", 10, 0).Return([]*model.CashHandlingEntryModel{}, nil)

	result, err := service.GetAllCashHandlingEntries(10, 0, "", "")

	assert.NoError(t, err)
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllCashHandlingEntriesRepositoryError(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	expectedError := errors.New("database connection failed")
	mockRepo.On("GetAll", 5, 10).Return(nil, expectedError)

	_, err := service.GetAllCashHandlingEntries(5, 10, "", "")

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllCashHandlingEntriesWithPagination(t *testing.T) {
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
			PaymentMethod: "creditcard",
			Description:   "Netflix subscription",
			Date:          createdTime,
			Timestamp:     createdTime.Unix(),
			CreatedAt:     createdTime,
			UpdatedAt:     createdTime,
		},
	}

	mockRepo.On("GetAll", 1, 5).Return(mockEntries, nil)

	result, err := service.GetAllCashHandlingEntries(1, 5, "", "")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, objectID.Hex(), result[0].ID)
	assert.Equal(t, float64(89.99), result[0].Amount)
	assert.Equal(t, "entertainment", result[0].Category)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllCashHandlingEntriesNoPagination(t *testing.T) {
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

	mockRepo.On("GetAll", 0, 0).Return(mockEntries, nil)

	result, err := service.GetAllCashHandlingEntries(0, 0, "", "")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, float64(320.50), result[0].Amount)
	assert.Equal(t, "utilities", result[0].Category)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllCashHandlingEntriesDateFormatting(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()

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

	result, err := service.GetAllCashHandlingEntries(10, 0, "", "")

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	assert.Equal(t, "15/03/2025", result[0].Date)

	assert.Equal(t, testDate.UTC().Format(time.RFC3339), result[0].CreatedAt)
	assert.Equal(t, testDate.UTC().Format(time.RFC3339), result[0].UpdatedAt)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllCashHandlingEntriesNilEntries(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	mockRepo.On("GetAll", 10, 0).Return(nil, nil)

	result, err := service.GetAllCashHandlingEntries(10, 0, "", "")

	assert.NoError(t, err)
	assert.Empty(t, result)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllParameterValidation(t *testing.T) {
	testCases := []struct {
		name          string
		inputLimit    int
		inputSkip     int
		expectedLimit int
		expectedSkip  int
		description   string
	}{
		{
			name:          "negativelimitdefaultsto10",
			inputLimit:    -5,
			inputSkip:     0,
			expectedLimit: 10,
			expectedSkip:  0,
			description:   "Negative limit should default to 10",
		},
		{
			name:          "limitover100cappedto100",
			inputLimit:    150,
			inputSkip:     0,
			expectedLimit: 100,
			expectedSkip:  0,
			description:   "Limit over 100 should be capped to 100",
		},
		{
			name:          "negativeskipdefaultsto0",
			inputLimit:    10,
			inputSkip:     -3,
			expectedLimit: 10,
			expectedSkip:  0,
			description:   "Negative skip should default to 0",
		},
		{
			name:          "validparametersunchanged",
			inputLimit:    20,
			inputSkip:     5,
			expectedLimit: 20,
			expectedSkip:  5,
			description:   "Valid parameters should remain unchanged",
		},
		{
			name:          "zerolimitunchanged",
			inputLimit:    0,
			inputSkip:     10,
			expectedLimit: 0,
			expectedSkip:  10,
			description:   "Zero limit should remain unchanged (no limit)",
		},
		{
			name:          "edgecaselimitexactly100",
			inputLimit:    100,
			inputSkip:     0,
			expectedLimit: 100,
			expectedSkip:  0,
			description:   "Limit exactly 100 should remain unchanged",
		},
		{
			name:          "multiplevalidationsapplied",
			inputLimit:    -10,
			inputSkip:     -5,
			expectedLimit: 10,
			expectedSkip:  0,
			description:   "Both negative limit and skip should be corrected",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

			mockRepo.On("GetAll", tc.expectedLimit, tc.expectedSkip).Return(mockEntries, nil)

			result, err := service.GetAllCashHandlingEntries(tc.inputLimit, tc.inputSkip, "", "")

			assert.NoError(t, err, tc.description)
			assert.Len(t, result, 1, tc.description)
			assert.Equal(t, objectID.Hex(), result[0].ID, tc.description)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCashHandlingServiceGetAllBoundaryValues(t *testing.T) {
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
			PaymentMethod: "banktransfer",
			Description:   "Boundary value test",
			Date:          createdTime,
			Timestamp:     createdTime.Unix(),
			CreatedAt:     createdTime,
			UpdatedAt:     createdTime,
		},
	}

	testCases := []struct {
		name          string
		limit         int
		skip          int
		expectedLimit int
		expectedSkip  int
	}{
		{"limit101cappedto100", 101, 0, 100, 0},
		{"limit99unchanged", 99, 0, 99, 0},
		{"limit1unchanged", 1, 0, 1, 0},
		{"skiplargepositiveunchanged", 50, 1000, 50, 1000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo.On("GetAll", tc.expectedLimit, tc.expectedSkip).Return(mockEntries, nil).Once()

			result, err := service.GetAllCashHandlingEntries(tc.limit, tc.skip, "", "")

			assert.NoError(t, err)
			assert.Len(t, result, 1)
		})
	}

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceParameterValidationWithError(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	expectedError := errors.New("repository error after parameter validation")

	mockRepo.On("GetAll", 10, 0).Return(nil, expectedError)

	result, err := service.GetAllCashHandlingEntries(-10, -5, "", "")

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllWithFilter(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID1 := primitive.NewObjectID()
	objectID2 := primitive.NewObjectID()
	createdTime := time.Date(2025, 9, 7, 10, 0, 0, 0, time.UTC)

	mockEntries := []*model.CashHandlingEntryModel{
		{
			ID:            objectID1,
			Amount:        150.75,
			Title:         "Lunch Restaurant",
			Currency:      "BRL",
			Type:          "expense",
			Category:      "food",
			PaymentMethod: "creditcard",
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

	result, err := service.GetAllCashHandlingEntries(10, 0, "lunch", "food")

	assert.NoError(t, err)
	assert.Len(t, result, 2)

	assert.Equal(t, objectID1.Hex(), result[0].ID)
	assert.Equal(t, "Lunch Restaurant", result[0].Title)
	assert.Equal(t, "food", result[0].Category)

	assert.Equal(t, objectID2.Hex(), result[1].ID)
	assert.Equal(t, "Fast food lunch", result[1].Title)
	assert.Equal(t, "food", result[1].Category)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllWithFilterTitleOnly(t *testing.T) {
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
			PaymentMethod: "creditcard",
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

	result, err := service.GetAllCashHandlingEntries(5, 2, "coffee", "")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Coffee Shop", result[0].Title)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllWithFilterCategoryOnly(t *testing.T) {
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

	result, err := service.GetAllCashHandlingEntries(10, 0, "", "transport")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "transport", result[0].Category)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllWithFilterParameterValidation(t *testing.T) {
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
			PaymentMethod: "banktransfer",
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

	mockRepo.On("GetAllWithFilter", 10, 0, expectedFilter).Return(mockEntries, nil)

	result, err := service.GetAllCashHandlingEntries(-5, -3, "test", "salary")

	assert.NoError(t, err)
	assert.Len(t, result, 1)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllWithFilterRepositoryError(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	expectedError := errors.New("database filter query failed")
	expectedFilter := types.FilterOptions{
		Title:    "error",
		Category: "test",
	}

	mockRepo.On("GetAllWithFilter", 10, 0, expectedFilter).Return(nil, expectedError)

	result, err := service.GetAllCashHandlingEntries(10, 0, "error", "test")

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetAllWithFilterEmptyResult(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	expectedFilter := types.FilterOptions{
		Title:    "nonexistent",
		Category: "unknown",
	}

	mockRepo.On("GetAllWithFilter", 10, 0, expectedFilter).Return([]*model.CashHandlingEntryModel{}, nil)

	result, err := service.GetAllCashHandlingEntries(10, 0, "nonexistent", "unknown")

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.NotNil(t, result) // Should return empty slice, not nil

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceUpdateCashHandlingEntrySuccess(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()
	createdTime := time.Date(2025, 9, 6, 14, 30, 0, 0, time.UTC)
	updatedTime := time.Now().UTC()

	existingEntry := &model.CashHandlingEntryModel{
		ID:            objectID,
		Amount:        150.75,
		Title:         "Old Title",
		Currency:      "BRL",
		Type:          "expense",
		Category:      "food",
		PaymentMethod: "creditcard",
		Description:   "Old Description",
		Date:          createdTime,
		Timestamp:     createdTime.Unix(),
		CreatedAt:     createdTime,
		UpdatedAt:     createdTime,
	}

	updatedEntry := &model.CashHandlingEntryModel{
		ID:            objectID,
		Amount:        200.50,
		Title:         "Updated Title",
		Currency:      "USD",
		Type:          "expense",
		Category:      "entertainment",
		PaymentMethod: "debitcard",
		Description:   "Updated Description",
		Date:          time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC),
		Timestamp:     createdTime.Unix(), // Timestamp should remain the same
		CreatedAt:     createdTime,        // CreatedAt should remain the same
		UpdatedAt:     updatedTime,        // UpdatedAt should be updated
	}

	updateDTO := dtos.UpdateCashHandlingEntryDTO{
		Amount:        200.50,
		Title:         "Updated Title",
		Currency:      "USD",
		Type:          "expense",
		Category:      "entertainment",
		PaymentMethod: "debitcard",
		Description:   "Updated Description",
		Date:          "15/10/2025",
	}

	mockRepo.On("GetByID", objectID.Hex()).Return(existingEntry, nil)
	mockRepo.On("Update", objectID.Hex(), mock.AnythingOfType("*model.CashHandlingEntryModel")).Return(updatedEntry, nil)

	result, err := service.UpdateCashHandlingEntry(objectID.Hex(), updateDTO)

	assert.NoError(t, err)
	assert.Equal(t, objectID.Hex(), result.ID)
	assert.Equal(t, float64(200.50), result.Amount)
	assert.Equal(t, "Updated Title", result.Title)
	assert.Equal(t, "USD", result.Currency)
	assert.Equal(t, "expense", result.Type)
	assert.Equal(t, "entertainment", result.Category)
	assert.Equal(t, "debitcard", result.PaymentMethod)
	assert.Equal(t, "Updated Description", result.Description)
	assert.Equal(t, "15/10/2025", result.Date)
	assert.Equal(t, createdTime.Unix(), result.Timestamp) // Timestamp should remain unchanged
	assert.Equal(t, createdTime.UTC().Format(time.RFC3339), result.CreatedAt)
	assert.Equal(t, updatedTime.UTC().Format(time.RFC3339), result.UpdatedAt)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceUpdateCashHandlingEntryInvalidDate(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()

	updateDTO := dtos.UpdateCashHandlingEntryDTO{
		Amount:        200.50,
		Title:         "Updated Title",
		Currency:      "USD",
		Type:          "expense",
		Category:      "entertainment",
		PaymentMethod: "debitcard",
		Description:   "Updated Description",
		Date:          "invalid-date", // Invalid date format
	}

	result, err := service.UpdateCashHandlingEntry(objectID.Hex(), updateDTO)

	assert.Error(t, err)
	assert.Equal(t, dtos.CashHandlingEntryResponseDTO{}, result)

	mockRepo.AssertNotCalled(t, "GetByID")
	mockRepo.AssertNotCalled(t, "Update")
}

func TestCashHandlingServiceUpdateCashHandlingEntryGetByIDError(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()
	expectedError := errors.New("entry not found")

	updateDTO := dtos.UpdateCashHandlingEntryDTO{
		Amount:        200.50,
		Title:         "Updated Title",
		Currency:      "USD",
		Type:          "expense",
		Category:      "entertainment",
		PaymentMethod: "debitcard",
		Description:   "Updated Description",
		Date:          "15/10/2025",
	}

	mockRepo.On("GetByID", objectID.Hex()).Return(nil, expectedError)

	result, err := service.UpdateCashHandlingEntry(objectID.Hex(), updateDTO)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dtos.CashHandlingEntryResponseDTO{}, result)

	mockRepo.AssertNotCalled(t, "Update")
}

func TestCashHandlingServiceUpdateCashHandlingEntryUpdateError(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()
	createdTime := time.Date(2025, 9, 6, 14, 30, 0, 0, time.UTC)
	expectedError := errors.New("database update error")

	existingEntry := &model.CashHandlingEntryModel{
		ID:            objectID,
		Amount:        150.75,
		Title:         "Old Title",
		Currency:      "BRL",
		Type:          "expense",
		Category:      "food",
		PaymentMethod: "creditcard",
		Description:   "Old Description",
		Date:          createdTime,
		Timestamp:     createdTime.Unix(),
		CreatedAt:     createdTime,
		UpdatedAt:     createdTime,
	}

	updateDTO := dtos.UpdateCashHandlingEntryDTO{
		Amount:        200.50,
		Title:         "Updated Title",
		Currency:      "USD",
		Type:          "expense",
		Category:      "entertainment",
		PaymentMethod: "debitcard",
		Description:   "Updated Description",
		Date:          "15/10/2025",
	}

	mockRepo.On("GetByID", objectID.Hex()).Return(existingEntry, nil)
	mockRepo.On("Update", objectID.Hex(), mock.AnythingOfType("*model.CashHandlingEntryModel")).Return(nil, expectedError)

	result, err := service.UpdateCashHandlingEntry(objectID.Hex(), updateDTO)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dtos.CashHandlingEntryResponseDTO{}, result)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetCashHandlingEntryByIDSuccess(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()
	createdTime := time.Date(2025, 9, 6, 14, 30, 0, 0, time.UTC)

	mockEntry := &model.CashHandlingEntryModel{
		ID:            objectID,
		Amount:        150.75,
		Title:         "Lunch at restaurant",
		Currency:      "BRL",
		Type:          "expense",
		Category:      "food",
		PaymentMethod: "creditcard",
		Description:   "Lunch at restaurant",
		Date:          createdTime,
		Timestamp:     createdTime.Unix(),
		CreatedAt:     createdTime,
		UpdatedAt:     createdTime,
	}

	mockRepo.On("GetByID", objectID.Hex()).Return(mockEntry, nil)

	result, err := service.GetCashHandlingEntryByID(objectID.Hex())

	assert.NoError(t, err)
	assert.Equal(t, objectID.Hex(), result.ID)
	assert.Equal(t, float64(150.75), result.Amount)
	assert.Equal(t, "Lunch at restaurant", result.Title)
	assert.Equal(t, "BRL", result.Currency)
	assert.Equal(t, "expense", result.Type)
	assert.Equal(t, "food", result.Category)
	assert.Equal(t, "creditcard", result.PaymentMethod)
	assert.Equal(t, "Lunch at restaurant", result.Description)
	assert.Equal(t, "06/09/2025", result.Date)
	assert.Equal(t, createdTime.Unix(), result.Timestamp)
	assert.Equal(t, createdTime.UTC().Format(time.RFC3339), result.CreatedAt)
	assert.Equal(t, createdTime.UTC().Format(time.RFC3339), result.UpdatedAt)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetCashHandlingEntryByIDNotFound(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	objectID := primitive.NewObjectID()
	expectedError := errors.New("entry not found")

	mockRepo.On("GetByID", objectID.Hex()).Return(nil, expectedError)

	result, err := service.GetCashHandlingEntryByID(objectID.Hex())

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dtos.CashHandlingEntryResponseDTO{}, result)

	mockRepo.AssertExpectations(t)
}

func TestCashHandlingServiceGetCashHandlingEntryByIDInvalidID(t *testing.T) {
	mockRepo := new(MockCashHandlingRepository)
	service := NewCashHandlingService(mockRepo)

	invalidID := "invalid-id"
	expectedError := errors.New("invalid ID format")

	mockRepo.On("GetByID", invalidID).Return(nil, expectedError)

	result, err := service.GetCashHandlingEntryByID(invalidID)

	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, dtos.CashHandlingEntryResponseDTO{}, result)

	mockRepo.AssertExpectations(t)
}
