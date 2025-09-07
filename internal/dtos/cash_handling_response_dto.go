package dtos

type CashHandlingEntryResponseDTO struct {
	ID            string  `bson:"_id"`
	Amount        float64 `bson:"amount"`
	Title         string  `bson:"title"`
	Currency      string  `bson:"currency"`
	Type          string  `bson:"type"`
	Category      string  `bson:"category"`
	PaymentMethod string  `bson:"paymentMethod"`
	Description   string  `bson:"description,omitempty"`
	Date          string  `bson:"date"`
	Timestamp     int64   `bson:"timestamp"`
	CreatedAt     string  `bson:"createdAt"`
	UpdatedAt     string  `bson:"updatedAt"`
}
