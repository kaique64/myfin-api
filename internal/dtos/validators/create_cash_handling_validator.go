package validators

import (
	"myfin-api/internal/dtos"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateCreateCashHandlingEntry(ctx *gin.Context) (*dtos.CreateCashHandlingEntryDTO, bool) {
	var entry dtos.CreateCashHandlingEntryDTO

	if err := ctx.ShouldBindJSON(&entry); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, fieldError := range validationErrors {
				errors[fieldError.Field()] = getCreateCashHandlingValidationMessage(fieldError)
			}
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": errors,
			})
			return nil, false
		}

		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid JSON format",
			"details": err.Error(),
		})
		return nil, false
	}

	return &entry, true
}

// getCreateCashHandlingValidationMessage returns appropriate error messages for CreateCashHandlingEntryDTO validation
func getCreateCashHandlingValidationMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "gt":
		return "Value must be greater than 0"
	case "len":
		return "Must be exactly 3 characters"
	case "min":
		return "Value is too short"
	case "oneof":
		return "Must be either 'income' or 'expense'"
	case "datetime":
		return "Date must be in DD/MM/YYYY format (e.g., 31/12/2025)"
	default:
		return "Invalid value"
	}
}
