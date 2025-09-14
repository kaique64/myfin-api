package validators

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ValidateGetAllPaginationParams(ctx *gin.Context) (int, int, bool) {
	limitStr := ctx.DefaultQuery("limit", "10")
	skipStr := ctx.DefaultQuery("skip", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid limit parameter",
			"details": "Limit must be a valid integer",
		})

		return 0, 0, false
	}

	skip, err := strconv.Atoi(skipStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid skip parameter",
			"details": "Skip must be a valid integer",
		})

		return 0, 0, false
	}

	return limit, skip, true
}
