package handlers

import (
	"myfin-api/internal/dtos"
	"myfin-api/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CashHandlingHandler interface {
	Save(ctx *gin.Context)
	GetAll(ctx *gin.Context)
}

type cashHandlingHandler struct {
	cashHandlingService services.CashHandlingService
}

func NewCashHandlingHandler(cashHandlingService services.CashHandlingService) CashHandlingHandler {
	return &cashHandlingHandler{
		cashHandlingService: cashHandlingService,
	}
}

func (h *cashHandlingHandler) Save(ctx *gin.Context) {
	var entry dtos.CreateCashHandlingEntryDTO

	ctx.BindJSON(&entry)

	response, err := h.cashHandlingService.CreateCashHandlingEntry(entry)

	if err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *cashHandlingHandler) GetAll(ctx *gin.Context) {
	// Get pagination parameters from query string with defaults
	limitStr := ctx.DefaultQuery("limit", "10")
	skipStr := ctx.DefaultQuery("skip", "0")

	// Parse limit parameter
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid limit parameter",
			"details": "Limit must be a valid integer",
		})
		return
	}

	// Parse skip parameter
	skip, err := strconv.Atoi(skipStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid skip parameter",
			"details": "Skip must be a valid integer",
		})
		return
	}

	// Validate parameters (optional - set reasonable limits)
	if limit < 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Max 100 entries per request
	}
	if skip < 0 {
		skip = 0
	}

	// Call service
	entries, err := h.cashHandlingService.GetAllCashHandlingEntries(limit, skip)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve entries",
			"details": err.Error(),
		})
		return
	}

	// Return response with pagination info
	ctx.JSON(http.StatusOK, gin.H{
		"data": entries,
		"pagination": gin.H{
			"limit": limit,
			"skip":  skip,
			"count": len(entries),
		},
	})
}
