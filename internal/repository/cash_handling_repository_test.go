package repository_test

import (
	"myfin-api/internal/model"
	"myfin-api/internal/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestCashHandlingEntryRepositoryCreate(t *testing.T) {
	// Create MongoDB test harness
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Close()

	// Test case: Successful creation without timestamp
	mt.Run("successful_creation_without_timestamp", func(mt *mtest.T) {
		// Set up the expected MongoDB response
		objectID := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "ok", Value: 1},
			bson.E{Key: "n", Value: 1},
			bson.E{Key: "id", Value: objectID},
		))

		// Create the repository with the mock MongoDB client
		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		// Create a test entry
		entry := &model.CashHandlingEntryModel{
			Amount:      100.0,
			Description: "Test entry",
			// Timestamp intentionally omitted
		}

		// Test the Create method
		beforeCreate := time.Now().Unix()
		result, err := repo.Create(entry)

		// Verify the result
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.GreaterOrEqual(t, result.Timestamp, beforeCreate)
		assert.NotZero(t, result.CreatedAt)
		assert.NotZero(t, result.UpdatedAt)
	})

	// Test case: Successful creation with timestamp
	mt.Run("successful_creation_with_timestamp", func(mt *mtest.T) {
		// Set up the expected MongoDB response
		mt.AddMockResponses(mtest.CreateSuccessResponse(
			bson.E{Key: "ok", Value: 1},
			bson.E{Key: "n", Value: 1},
		))

		// Create the repository with the mock MongoDB client
		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		// Create a test entry with timestamp
		customTimestamp := int64(1234567890)
		entry := &model.CashHandlingEntryModel{
			Amount:      200.0,
			Description: "Test entry with timestamp",
			Timestamp:   customTimestamp,
		}

		// Test the Create method
		result, err := repo.Create(entry)

		// Verify the result
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, customTimestamp, result.Timestamp)
		assert.NotZero(t, result.CreatedAt)
		assert.NotZero(t, result.UpdatedAt)
	})

	// Test case: Database error
	mt.Run("database_error", func(mt *mtest.T) {
		// Set up the expected MongoDB error response
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "duplicate key error",
		}))

		// Create the repository with the mock MongoDB client
		repo := repository.NewCashHandlingEntryRepository(mt.DB)

		// Create a test entry
		entry := &model.CashHandlingEntryModel{
			Amount:      300.0,
			Description: "Test entry that will fail",
		}

		// Test the Create method
		result, err := repo.Create(entry)

		// Verify the result
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestCashHandlingEntryRepositoryGetAll(t *testing.T) {
    // Create MongoDB test harness
    mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
    defer mt.Close()

    // Test case: Successful retrieval with data
    mt.Run("successful_retrieval_with_data", func(mt *mtest.T) {
        // Create mock data
        objectID1 := primitive.NewObjectID()
        objectID2 := primitive.NewObjectID()
        
        // Set up the expected MongoDB response
        first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
            {"_id", objectID1},
            {"amount", 150.75},
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

        // Create the repository
        repo := repository.NewCashHandlingEntryRepository(mt.DB)

        // Test GetAll with no limits
        result, err := repo.GetAll(0, 0)

        // Verify the result
        assert.NoError(t, err)
        assert.NotNil(t, result)
        assert.Len(t, result, 2)

        // Verify first entry
        assert.Equal(t, objectID1, result[0].ID)
        assert.Equal(t, 150.75, result[0].Amount)
        assert.Equal(t, "BRL", result[0].Currency)
        assert.Equal(t, "expense", result[0].Type)
        assert.Equal(t, "food", result[0].Category)
        assert.Equal(t, "credit_card", result[0].PaymentMethod)
        assert.Equal(t, "Lunch at restaurant", result[0].Description)

        // Verify second entry
        assert.Equal(t, objectID2, result[1].ID)
        assert.Equal(t, 2500.0, result[1].Amount)
        assert.Equal(t, "income", result[1].Type)
        assert.Equal(t, "salary", result[1].Category)
    })

    // Test case: With limit and skip parameters
    mt.Run("with_limit_and_skip", func(mt *mtest.T) {
        // Create mock data
        objectID := primitive.NewObjectID()
        
        // Set up the expected MongoDB response with one item
        first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
            {"_id", objectID},
            {"amount", 89.99},
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

        // Create the repository
        repo := repository.NewCashHandlingEntryRepository(mt.DB)

        // Test GetAll with limit and skip
        result, err := repo.GetAll(1, 5)

        // Verify the result
        assert.NoError(t, err)
        assert.NotNil(t, result)
        assert.Len(t, result, 1)
        assert.Equal(t, objectID, result[0].ID)
        assert.Equal(t, 89.99, result[0].Amount)
        assert.Equal(t, "entertainment", result[0].Category)
    })

    // Test case: Database connection error
    mt.Run("database_connection_error", func(mt *mtest.T) {
        // Set up command error response
        mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
            Code:    2,
            Message: "BadValue",
        }))

        // Create the repository
        repo := repository.NewCashHandlingEntryRepository(mt.DB)

        // Test GetAll
        result, err := repo.GetAll(10, 0)

        // Verify the result
        assert.Error(t, err)
        assert.Nil(t, result)
        assert.Contains(t, err.Error(), "BadValue")
    })

    // Test case: Cursor decode error
    mt.Run("cursor_decode_error", func(mt *mtest.T) {
        // Set up response with invalid data that will cause decode error
        first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
            {"_id", "invalid_object_id"}, // This should cause a decode error
            {"amount", "invalid_amount"},  // This should also cause issues
        })
        
        killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

        mt.AddMockResponses(first, killCursors)

        // Create the repository
        repo := repository.NewCashHandlingEntryRepository(mt.DB)

        // Test GetAll
        result, err := repo.GetAll(10, 0)

        // Verify the result - this should fail during decode
        assert.Error(t, err)
        assert.Nil(t, result)
    })

    // Test case: Only limit parameter (no skip)
    mt.Run("only_limit_parameter", func(mt *mtest.T) {
        // Create mock data
        objectID := primitive.NewObjectID()
        
        first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
            {"_id", objectID},
            {"amount", 45.90},
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

        // Create the repository
        repo := repository.NewCashHandlingEntryRepository(mt.DB)

        // Test GetAll with only limit (skip = 0)
        result, err := repo.GetAll(5, 0)

        // Verify the result
        assert.NoError(t, err)
        assert.NotNil(t, result)
        assert.Len(t, result, 1)
        assert.Equal(t, 45.90, result[0].Amount)
        assert.Equal(t, "transport", result[0].Category)
    })

    // Test case: Only skip parameter (no limit)
    mt.Run("only_skip_parameter", func(mt *mtest.T) {
        // Create mock data
        objectID := primitive.NewObjectID()
        
        first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
            {"_id", objectID},
            {"amount", 320.50},
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

        // Create the repository
        repo := repository.NewCashHandlingEntryRepository(mt.DB)

        // Test GetAll with only skip (limit = 0, meaning no limit)
        result, err := repo.GetAll(0, 3)

        // Verify the result
        assert.NoError(t, err)
        assert.NotNil(t, result)
        assert.Len(t, result, 1)
        assert.Equal(t, 320.50, result[0].Amount)
        assert.Equal(t, "utilities", result[0].Category)
    })

    // Test case: Negative parameters (should be ignored)
    mt.Run("negative_parameters", func(mt *mtest.T) {
        // Create mock data
        objectID := primitive.NewObjectID()
        
        first := mtest.CreateCursorResponse(1, "cash_handling_entries.entries", mtest.FirstBatch, bson.D{
            {"_id", objectID},
            {"amount", 100.0},
            {"currency", "USD"},
            {"type", "income"},
        })
        
        killCursors := mtest.CreateCursorResponse(0, "cash_handling_entries.entries", mtest.NextBatch)

        mt.AddMockResponses(first, killCursors)

        // Create the repository
        repo := repository.NewCashHandlingEntryRepository(mt.DB)

        // Test GetAll with negative parameters (should be treated as 0)
        result, err := repo.GetAll(-5, -10)

        // Verify the result - should work since negative values are ignored
        assert.NoError(t, err)
        assert.NotNil(t, result)
        assert.Len(t, result, 1)
    })
}