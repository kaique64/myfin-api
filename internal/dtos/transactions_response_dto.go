package dtos

type TransactionsEntryResponseDTO struct {
	ID            string  `bson:"_id" json:"id"`
	Amount        float64 `bson:"amount" json:"amount"`
	Title         string  `bson:"title" json:"title"`
	Currency      string  `bson:"currency" json:"currency"`
	Type          string  `bson:"type" json:"type"`
	Category      string  `bson:"category" json:"category"`
	PaymentMethod string  `bson:"paymentMethod" json:"paymentMethod"`
	Description   string  `bson:"description,omitempty" json:"description,omitempty"`
	Date          string  `bson:"date" json:"date"`
	Timestamp     int64   `bson:"timestamp" json:"timestamp"`
	CreatedAt     string  `bson:"createdAt" json:"createdAt"`
	UpdatedAt     string  `bson:"updatedAt" json:"updatedAt"`
}
