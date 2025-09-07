package repository_test

import (
	"myfin-api/internal/model"
	"myfin-api/internal/repository"
	"myfin-api/internal/repository/types"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCashHandlingEntryRepositoryCreate(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("successful_creation_without_timestamp", func(mt *mtest.T) {
		objectID := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "ok", Value: 1},
			bson.E{Key: "n", Value: 1},
			bson.E{Key: "id", Value: objectID},
		))

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		entry := &model.CashHandlingEntryModel{
			Amount:      100.0,
			Description: "Test entry",
		}

		beforeCreate := time.Now().Unix()
		result, err := repo.Create(entry)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.GreaterOrEqual(t, result.Timestamp, beforeCreate)
		assert.NotZero(t, result.CreatedAt)
		assert.NotZero(t, result.UpdatedAt)
	})

	mt.Run("successful_creation_with_timestamp", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "ok", Value: 1},
			bson.E{Key: "n", Value: 1},
		))

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		customTimestamp := int64(1234567890)
		entry := &model.CashHandlingEntryModel{
			Amount:      200.0,
			Description: "Test entry with timestamp",
			Timestamp:   customTimestamp,
		}

		result, err := repo.Create(entry)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, customTimestamp, result.Timestamp)
		assert.NotZero(t, result.CreatedAt)
		assert.NotZero(t, result.UpdatedAt)
	})

	mt.Run("database_error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "duplicate key error",
		}))

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		entry := &model.CashHandlingEntryModel{
			Amount:      300.0,
			Description: "Test entry that will fail",
		}

		result, err := repo.Create(entry)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestCashHandlingEntryRepositoryGetAll(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	mt.Run("successful_retrieval_with_data", func(mt *mtest.T) {
		objectID1 := primitive.NewObjectID()
		objectID2 := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{"_id", objectID1},
			{"amount", 150.75},
			{"title", "Lunch at restaurant"},
			{"currency", "BRL"},
			{"type", "expense"},
			{"category", "food"},
			{"payment_method", "credit_card"},
			{"description", "Lunch at restaurant"},
			{"date", time.Date(2025, 9, 6, 14, 30, 0, 0, time.UTC)},
			{"timestamp", int64(1725635400)},
			{"created_at", time.Date(2025, 9, 6, 14, 30, 0, 0, time.UTC)},
			{"updated_at", time.Date(2025, 9, 6, 14, 30, 0, 0, time.UTC)},
		})

		second := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.NextBatch, bson.D{
			{"_id", objectID2},
			{"amount", 2500.0},
			{"title", "Monthly salary"},
			{"currency", "BRL"},
			{"type", "income"},
			{"category", "salary"},
			{"payment_method", "bank_transfer"},
			{"description", "Monthly salary"},
			{"date", time.Date(2025, 8, 30, 9, 0, 0, 0, time.UTC)},
			{"timestamp", int64(1725008400)},
			{"created_at", time.Date(2025, 8, 30, 9, 0, 0, 0, time.UTC)},
			{"updated_at", time.Date(2025, 8, 30, 9, 0, 0, 0, time.UTC)},
		})

		killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		result, err := repo.GetAll(0, 0)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)

		assert.Equal(t, objectID1, result[0].ID)
		assert.Equal(t, 150.75, result[0].Amount)
		assert.Equal(t, "BRL", result[0].Currency)
		assert.Equal(t, "expense", result[0].Type)
		assert.Equal(t, "food", result[0].Category)
		assert.Equal(t, "credit_card", result[0].PaymentMethod)
		assert.Equal(t, "Lunch at restaurant", result[0].Description)

		assert.Equal(t, objectID2, result[1].ID)
		assert.Equal(t, 2500.0, result[1].Amount)
		assert.Equal(t, "income", result[1].Type)
		assert.Equal(t, "salary", result[1].Category)
	})

	mt.Run("with_limit_and_skip", func(mt *mtest.T) {
		objectID := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{"_id", objectID},
			{"amount", 89.99},
			{"title", "Lunch at restaurant"},
			{"currency", "BRL"},
			{"type", "expense"},
			{"category", "entertainment"},
			{"payment_method", "credit_card"},
			{"description", "Netflix subscription"},
			{"date", time.Date(2025, 8, 28, 20, 15, 0, 0, time.UTC)},
			{"timestamp", int64(1724873700)},
			{"created_at", time.Date(2025, 8, 28, 20, 15, 0, 0, time.UTC)},
			{"updated_at", time.Date(2025, 8, 28, 20, 15, 0, 0, time.UTC)},
		})

		killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		result, err := repo.GetAll(1, 5)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, objectID, result[0].ID)
		assert.Equal(t, 89.99, result[0].Amount)
		assert.Equal(t, "entertainment", result[0].Category)
	})

	mt.Run("database_connection_error", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    2,
			Message: "BadValue",
		}))

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		result, err := repo.GetAll(10, 0)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "BadValue")
	})
	mt.Run("cursor_decode_error", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{"_id", "invalid_object_id"},
			{"amount", "invalid_amount"},
		})

		killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		result, err := repo.GetAll(10, 0)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
	mt.Run("only_limit_parameter", func(mt *mtest.T) {
		objectID := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{"_id", objectID},
			{"amount", 45.90},
			{"title", "Lunch at restaurant"},
			{"currency", "BRL"},
			{"type", "expense"},
			{"category", "transport"},
			{"payment_method", "pix"},
			{"description", "Gas station"},
			{"date", time.Date(2025, 8, 29, 18, 45, 0, 0, time.UTC)},
			{"timestamp", int64(1724954700)},
		})

		killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		result, err := repo.GetAll(5, 0)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, 45.90, result[0].Amount)
		assert.Equal(t, "transport", result[0].Category)
	})
	mt.Run("only_skip_parameter", func(mt *mtest.T) {
		objectID := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{"_id", objectID},
			{"amount", 320.50},
			{"title", "Lunch at restaurant"},
			{"currency", "BRL"},
			{"type", "expense"},
			{"category", "utilities"},
			{"payment_method", "boleto"},
			{"description", "Electric bill"},
			{"date", time.Date(2025, 8, 27, 14, 22, 0, 0, time.UTC)},
			{"timestamp", int64(1724765320)},
		})

		killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		result, err := repo.GetAll(0, 3)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, 320.50, result[0].Amount)
		assert.Equal(t, "utilities", result[0].Category)
	})
	mt.Run("negative_parameters", func(mt *mtest.T) {
		objectID := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{"_id", objectID},
			{"title", "Lunch at restaurant"},
			{"amount", 100.0},
			{"currency", "USD"},
			{"type", "income"},
		})

		killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		result, err := repo.GetAll(-5, -10)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
	})
}

func TestCashHandlingRepositoryGetAllWithFilter(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()
	mt.Run("filter_by_title_only", func(mt *mtest.T) {
		// Create mock data
		objectID := primitive.NewObjectID()

		// Set up the expected MongoDB response
		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: objectID},
			{Key: "amount", Value: 50.0},
			{Key: "title", Value: "Coffee Shop Visit"},
			{Key: "currency", Value: "BRL"},
			{Key: "type", Value: "expense"},
			{Key: "category", Value: "food"},
			{Key: "payment_method", Value: "credit_card"},
			{Key: "description", Value: "Morning coffee"},
			{Key: "date", Value: time.Date(2025, 9, 7, 10, 0, 0, 0, time.UTC)},
			{Key: "timestamp", Value: int64(1725700800)},
			{Key: "created_at", Value: time.Date(2025, 9, 7, 10, 0, 0, 0, time.UTC)},
			{Key: "updated_at", Value: time.Date(2025, 9, 7, 10, 0, 0, 0, time.UTC)},
		})

		killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		// Create the repository
		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		// Create filter with title only
		filter := types.FilterOptions{
			Title:    "coffee",
			Category: "", // Empty category
		}

		// Test GetAllWithFilter
		result, err := repo.GetAllWithFilter(10, 0, filter)

		// Verify the result
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 1)
		assert.Equal(t, objectID, result[0].ID)
		assert.Equal(t, "Coffee Shop Visit", result[0].Title)
		assert.Equal(t, "food", result[0].Category)
	})

	// Test case: Filter by category only
	mt.Run("filter_by_category_only", func(mt *mtest.T) {
		objectID := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: objectID},
			{Key: "amount", Value: 25.0},
			{Key: "title", Value: "Bus Ticket"},
			{Key: "currency", Value: "BRL"},
			{Key: "type", Value: "expense"},
			{Key: "category", Value: "transport"},
			{Key: "payment_method", Value: "cash"},
			{Key: "description", Value: "Public transport"},
			{Key: "date", Value: time.Date(2025, 9, 7, 14, 30, 0, 0, time.UTC)},
			{Key: "timestamp", Value: int64(1725717000)},
		})

		killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		filter := types.FilterOptions{
			Title:    "", // Empty title
			Category: "transport",
		}

		result, err := repo.GetAllWithFilter(5, 0, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "transport", result[0].Category)
		assert.Equal(t, "Bus Ticket", result[0].Title)
	})

	// Test case: Filter by both title and category
	mt.Run("filter_by_both_title_and_category", func(mt *mtest.T) {
		objectID1 := primitive.NewObjectID()
		objectID2 := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{"_id", objectID1},
			{"amount", 120.0},
			{"title", "Lunch at Restaurant"},
			{"currency", "BRL"},
			{"type", "expense"},
			{"category", "food"},
			{"payment_method", "credit_card"},
			{"description", "Business lunch"},
			{"date", time.Date(2025, 9, 7, 12, 0, 0, 0, time.UTC)},
		})

		second := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.NextBatch, bson.D{
			{"_id", objectID2},
			{"amount", 35.0},
			{"title", "Quick lunch"},
			{"currency", "BRL"},
			{"type", "expense"},
			{"category", "food"},
			{"payment_method", "cash"},
			{"description", "Fast food"},
			{"date", time.Date(2025, 9, 6, 13, 0, 0, 0, time.UTC)},
		})

		killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, second, killCursors)

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		filter := types.FilterOptions{
			Title:    "lunch",
			Category: "food",
		}

		result, err := repo.GetAllWithFilter(10, 0, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		// Verify both entries match the filter criteria
		for _, entry := range result {
			assert.Equal(t, "food", entry.Category)
			// Check case-insensitive contains for title
			assert.Contains(t, strings.ToLower(entry.Title), "lunch") // Should contain "lunch" (case insensitive)
		}
	})

	// Test case: No filters (empty strings)
	mt.Run("no_filters_empty_strings", func(mt *mtest.T) {
		objectID := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{"_id", objectID},
			{"amount", 200.0},
			{"title", "Any Title"},
			{"category", "any_category"},
		})

		killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		filter := types.FilterOptions{
			Title:    "", // Empty filters should return all entries
			Category: "",
		}

		result, err := repo.GetAllWithFilter(10, 0, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Any Title", result[0].Title)
	})

	// Test case: Zero and negative parameters with filter
	mt.Run("zero_and_negative_params_with_filter", func(mt *mtest.T) {
		objectID := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{"_id", objectID},
			{"amount", 150.0},
			{"title", "Test Entry"},
			{"category", "test"},
		})

		killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		filter := types.FilterOptions{
			Title:    "test",
			Category: "",
		}

		// Test with limit=0 and skip=-1 to ensure these branches are covered
		result, err := repo.GetAllWithFilter(0, -1, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Test Entry", result[0].Title)
	})

	// Test case: Database error with filters
	mt.Run("database_error_with_filters", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    2,
			Message: "BadValue in filter query",
		}))

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		filter := types.FilterOptions{
			Title:    "test",
			Category: "error",
		}

		result, err := repo.GetAllWithFilter(10, 0, filter)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "BadValue")
	})

	// Test case: Cursor iteration error after successful query
	mt.Run("cursor_iteration_error_with_filters", func(mt *mtest.T) {
		// Create a response that will cause cursor iteration issues
		// This simulates a scenario where Find() succeeds but cursor.Next() or cursor.Err() fails
		objectID := primitive.NewObjectID()

		// First response with valid data
		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: objectID},
			{Key: "amount", Value: 100.0},
			{Key: "title", Value: "Test Entry"},
			{Key: "currency", Value: "BRL"},
			{Key: "type", Value: "expense"},
			{Key: "category", Value: "test"},
			{Key: "payment_method", Value: "cash"},
			{Key: "date", Value: time.Date(2025, 9, 7, 10, 0, 0, 0, time.UTC)},
		})

		// Second response that triggers a cursor error
		cursorError := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, cursorError)

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		filter := types.FilterOptions{
			Title:    "test",
			Category: "",
		}

		result, err := repo.GetAllWithFilter(10, 0, filter)

		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, result)
		} else {
			assert.NoError(t, err)
			assert.NotNil(t, result)
		}
	})

	mt.Run("comprehensive_coverage_test", func(mt *mtest.T) {
		objectID := primitive.NewObjectID()

		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{"_id", objectID},
			{"amount", 75.0},
			{"title", "Complete Test"},
			{"category", "coverage"},
			{"currency", "BRL"},
			{"type", "expense"},
			{"payment_method", "cash"},
			{"date", time.Date(2025, 9, 7, 15, 0, 0, 0, time.UTC)},
		})

		killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		filter := types.FilterOptions{
			Title:    "Complete",
			Category: "coverage",
		}

		result, err := repo.GetAllWithFilter(5, 2, filter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "Complete Test", result[0].Title)
		assert.Equal(t, "coverage", result[0].Category)
	})

	mt.Run("cursor_decode_error_with_filter", func(mt *mtest.T) {
		first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
			{"_id", "invalid-object-id"},
			{"amount", "invalid-amount"},
			{"title", 123},              
		})

		killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

		mt.AddMockResponses(first, killCursors)

		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		filter := types.FilterOptions{
			Title:    "test",
			Category: "",
		}

		result, err := repo.GetAllWithFilter(10, 0, filter)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
