package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionsEntryModel struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Amount        float64            `bson:"amount" json:"amount"`
	Title         string             `bson:"title" json:"title"`
	Currency      string             `bson:"currency" json:"currency"`
	Type          string             `bson:"type" json:"type"`
	Category      string             `bson:"category" json:"category"`
	PaymentMethod string             `bson:"payment_method" json:"payment_method"`
	Description   string             `bson:"description,omitempty" json:"description,omitempty"`
	Date          time.Time          `bson:"date" json:"date"`
	Timestamp     int64              `bson:"timestamp" json:"timestamp"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updated_at"`
}
