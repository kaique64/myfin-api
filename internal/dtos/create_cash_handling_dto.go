package dtos

type CreateCashHandlingEntryDTO struct {
	Amount        float64 `bson:"amount"`
	Currency      string  `bson:"currency"`
	Type          string  `bson:"type"`
	Category      string  `bson:"category"`
	PaymentMethod string  `bson:"paymentMethod"`
	Description   string  `bson:"description,omitempty"`
	Date          string  `bson:"date"`
}
