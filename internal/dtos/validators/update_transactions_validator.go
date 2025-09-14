package validators

import (
	"net/http"

	"myfin-api/internal/dtos"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateUpdateTransactionsEntry(ctx *gin.Context) (*dtos.UpdateTransactionsEntryDTO, string, bool) {
	var entry dtos.UpdateTransactionsEntryDTO

	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID is required",
		})
		return nil, "", false
	}

	if err := ctx.ShouldBindJSON(&entry); err != nil {
		var errorMessage string
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				errorMessage = getUpdateTransactionsValidationMessage(fieldError)
				break
			}
		} else {
			errorMessage = err.Error()
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errorMessage,
		})
		return nil, "", false
	}

	return &entry, id, true
}

func getUpdateTransactionsValidationMessage(fieldError validator.FieldError) string {
	switch fieldError.Field() {
	case "Amount":
		return "Amount must be greater than 0"
	case "Title":
		return "Title is required"
	case "Currency":
		return "Currency must be a 3-letter code"
	case "Type":
		return "Type must be either 'income' or 'expense'"
	case "Category":
		return "Category is required"
	case "PaymentMethod":
		return "Payment method is required"
	case "Description":
		if fieldError.Tag() == "min" {
			return "Description must not be empty"
		}
	case "Date":
		return "Date must be in format DD/MM/YYYY"
	}
	return fieldError.Error()
}
