package dtos

type CreateCashHandlingEntryDTO struct {
	Amount        float64 `json:"amount" validate:"required,gt=0" binding:"required"`
	Title         string  `json:"title" binding:"required"`
	Currency      string  `json:"currency" validate:"required,len=3" binding:"required"`
	Type          string  `json:"type" validate:"required,oneof=income expense" binding:"required"`
	Category      string  `json:"category" validate:"required,min=1" binding:"required"`
	PaymentMethod string  `json:"paymentMethod" validate:"required,min=1" binding:"required"`
	Description   string  `json:"description,omitempty" validate:"omitempty,min=1"`
	Date          string  `json:"date" validate:"required,datetime=02/01/2006" binding:"required"`
}
