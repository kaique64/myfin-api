package repository

import (
	"context"
	"myfin-api/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CashHandlingEntryRepository interface {
	Create(entry *model.CashHandlingEntryModel) (*model.CashHandlingEntryModel, error)
}

type cashHandlingEntryRepository struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewCashHandlingEntryRepository(database *mongo.Database) CashHandlingEntryRepository {
	collection := database.Collection("cash_handling_entries")
	return &cashHandlingEntryRepository{
		database:   database,
		collection: collection,
	}
}

func (r *cashHandlingEntryRepository) Create(entry *model.CashHandlingEntryModel) (*model.CashHandlingEntryModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if entry.Timestamp == 0 {
		entry.Timestamp = time.Now().Unix()
	}

	entry.CreatedAt = time.Now().UTC().Local()
	entry.UpdatedAt = time.Now().UTC().Local()

	result, err := r.collection.InsertOne(ctx, entry)
	if err != nil {
		return nil, err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		entry.ID = oid
	}

	return entry, nil
}
