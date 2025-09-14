package repository

import (
	"context"
	"time"

	"myfin-api/internal/model"
	"myfin-api/internal/repository/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionsEntryRepository interface {
	Create(entry *model.TransactionsEntryModel) (*model.TransactionsEntryModel, error)
	GetAll(limit, skip int) ([]*model.TransactionsEntryModel, error)
	GetAllWithFilter(limit, skip int, filter types.FilterOptions) ([]*model.TransactionsEntryModel, error)
	Delete(id string) error
	Update(id string, entry *model.TransactionsEntryModel) (*model.TransactionsEntryModel, error)
	GetByID(id string) (*model.TransactionsEntryModel, error)
}

type transactionsEntryRepository struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewTransactionsEntryRepository(database *mongo.Database) TransactionsEntryRepository {
	collection := database.Collection("transactions_entries")
	return &transactionsEntryRepository{
		database:   database,
		collection: collection,
	}
}

func (r *transactionsEntryRepository) Create(entry *model.TransactionsEntryModel) (*model.TransactionsEntryModel, error) {
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

func (r *transactionsEntryRepository) GetAll(limit, skip int) ([]*model.TransactionsEntryModel, error) {
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

	entries := make([]*model.TransactionsEntryModel, 0)

	for cursor.Next(ctx) {
		var entry model.TransactionsEntryModel
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}

	return entries, nil
}

func (r *transactionsEntryRepository) GetAllWithFilter(limit, skip int, filter types.FilterOptions) ([]*model.TransactionsEntryModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := bson.M{}

	if filter.Title != "" {
		query["title"] = bson.M{"$regex": filter.Title, "$options": "i"}
	}

	if filter.Category != "" {
		query["category"] = bson.M{"$regex": "^" + filter.Category + "$", "$options": "i"}
	}

	options := options.Find()

	if limit > 0 {
		options.SetLimit(int64(limit))
	}

	if skip > 0 {
		options.SetSkip(int64(skip))
	}

	options.SetSort(bson.D{{Key: "date", Value: -1}})

	cursor, err := r.collection.Find(ctx, query, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	entries := make([]*model.TransactionsEntryModel, 0)

	for cursor.Next(ctx) {
		var entry model.TransactionsEntryModel
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}

	return entries, cursor.Err()
}

func (r *transactionsEntryRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}

func (r *transactionsEntryRepository) Update(id string, entry *model.TransactionsEntryModel) (*model.TransactionsEntryModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	entry.ID = objectID
	entry.UpdatedAt = time.Now().UTC().Local()

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"amount":         entry.Amount,
			"title":          entry.Title,
			"currency":       entry.Currency,
			"type":           entry.Type,
			"category":       entry.Category,
			"payment_method": entry.PaymentMethod,
			"description":    entry.Description,
			"date":           entry.Date,
			"updated_at":     entry.UpdatedAt,
		},
	}

	_, err = r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *transactionsEntryRepository) GetByID(id string) (*model.TransactionsEntryModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var entry model.TransactionsEntryModel
	err = r.collection.FindOne(ctx, filter).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}
