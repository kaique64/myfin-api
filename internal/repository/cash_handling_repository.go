package repository

import (
	"context"
	"myfin-api/internal/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CashHandlingEntryRepository interface {
	Create(entry *model.CashHandlingEntryModel) (*model.CashHandlingEntryModel, error)
	GetAll(limit, skip int) ([]*model.CashHandlingEntryModel, error)
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

func (r *cashHandlingEntryRepository) GetAll(limit, skip int) ([]*model.CashHandlingEntryModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	options := options.Find()

	if limit > 0 {
		options.SetLimit(int64(limit))
	}

	if skip > 0 {
		options.SetSkip(int64(skip))
	}

	options.SetSort(bson.D{{Key: "date", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, options)

	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var entries []*model.CashHandlingEntryModel

	for cursor.Next(ctx) {
		var entry model.CashHandlingEntryModel
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}

	return entries, nil
}
