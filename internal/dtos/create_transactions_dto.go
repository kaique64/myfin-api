package dtos

type CreateTransactionsEntryDTO struct {
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Title         string  `json:"title" binding:"required"`
	Currency      string  `json:"currency" binding:"required,len=3"`
	Type          string  `json:"type" binding:"required,oneof=income expense"`
	Category      string  `json:"category" binding:"required,min=1"`
	PaymentMethod string  `json:"paymentMethod" binding:"required,min=1"`
	Description   string  `json:"description" binding:"omitempty"`
	Date          string  `json:"date" binding:"required,datetime=02/01/2006"`
}
