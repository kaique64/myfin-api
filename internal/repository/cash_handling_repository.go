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

type CashHandlingEntryRepository interface {
	Create(entry *model.CashHandlingEntryModel) (*model.CashHandlingEntryModel, error)
	GetAll(limit, skip int) ([]*model.CashHandlingEntryModel, error)
	GetAllWithFilter(limit, skip int, filter types.FilterOptions) ([]*model.CashHandlingEntryModel, error)
	Delete(id string) error
	Update(id string, entry *model.CashHandlingEntryModel) (*model.CashHandlingEntryModel, error)
	GetByID(id string) (*model.CashHandlingEntryModel, error)
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

	entries := make([]*model.CashHandlingEntryModel, 0)

	for cursor.Next(ctx) {
		var entry model.CashHandlingEntryModel
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}

	return entries, nil
}

func (r *cashHandlingEntryRepository) GetAllWithFilter(limit, skip int, filter types.FilterOptions) ([]*model.CashHandlingEntryModel, error) {
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

	entries := make([]*model.CashHandlingEntryModel, 0)

	for cursor.Next(ctx) {
		var entry model.CashHandlingEntryModel
		if err := cursor.Decode(&entry); err != nil {
			return nil, err
		}
		entries = append(entries, &entry)
	}

	return entries, cursor.Err()
}

func (r *cashHandlingEntryRepository) Delete(id string) error {
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

func (r *cashHandlingEntryRepository) Update(id string, entry *model.CashHandlingEntryModel) (*model.CashHandlingEntryModel, error) {
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

func (r *cashHandlingEntryRepository) GetByID(id string) (*model.CashHandlingEntryModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}
	var entry model.CashHandlingEntryModel
	err = r.collection.FindOne(ctx, filter).Decode(&entry)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}
