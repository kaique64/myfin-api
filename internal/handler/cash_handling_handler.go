package handlers

import (
	"myfin-api/internal/dtos/validators"
	"myfin-api/internal/services"
	"net/http"

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
	entry, isValid := validators.ValidateCreateCashHandlingEntry(ctx)
	if !isValid {
		return
	}

	response, err := h.cashHandlingService.CreateCashHandlingEntry(*entry)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *cashHandlingHandler) GetAll(ctx *gin.Context) {
	limit, skip, isValid := validators.ValidateGetAllPaginationParams(ctx)
	if !isValid {
		return
	}

	titleFilter := ctx.Query("title")
	categoryFilter := ctx.Query("category")

	entries, err := h.cashHandlingService.GetAllCashHandlingEntries(limit, skip, titleFilter, categoryFilter)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve entries",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": entries,
		"pagination": gin.H{
			"limit": limit,
			"skip":  skip,
			"count": len(entries),
		},
		"filters": gin.H{
			"title":    titleFilter,
			"category": categoryFilter,
		},
	})
}
