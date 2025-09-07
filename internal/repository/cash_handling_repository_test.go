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