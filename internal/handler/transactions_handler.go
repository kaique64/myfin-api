package handlers

import (
	"net/http"

	"myfin-api/internal/dtos/validators"
	"myfin-api/internal/services"

	"github.com/gin-gonic/gin"
)

type TransactionsHandler interface {
	Save(ctx *gin.Context)
	GetAll(ctx *gin.Context)
	Delete(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetByID(ctx *gin.Context)
}

type transactionsHandler struct {
	transactionsService services.TransactionsService
}

func NewTransactionsHandler(transactionsService services.TransactionsService) TransactionsHandler {
	return &transactionsHandler{
		transactionsService: transactionsService,
	}
}

func (h *transactionsHandler) Save(ctx *gin.Context) {
	entry, isValid := validators.ValidateCreateTransactionsEntry(ctx)
	if !isValid {
		return
	}

	response, err := h.transactionsService.CreateTransactionsEntry(*entry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (h *transactionsHandler) GetAll(ctx *gin.Context) {
	limit, skip, isValid := validators.ValidateGetAllPaginationParams(ctx)
	if !isValid {
		return
	}

	titleFilter := ctx.Query("title")
	categoryFilter := ctx.Query("category")

	entries, err := h.transactionsService.GetAllTransactionsEntries(limit, skip, titleFilter, categoryFilter)
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

func (h *transactionsHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID is required",
		})
		return
	}

	err := h.transactionsService.DeleteTransactionsEntry(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete entry",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Entry deleted successfully",
		"id":      id,
	})
}

func (h *transactionsHandler) Update(ctx *gin.Context) {
	entry, id, isValid := validators.ValidateUpdateTransactionsEntry(ctx)
	if !isValid {
		return
	}

	response, err := h.transactionsService.UpdateTransactionsEntry(id, *entry)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update entry",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Entry updated successfully",
		"data":    response,
	})
}

func (h *transactionsHandler) GetByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "ID is required",
		})
		return
	}

	entry, err := h.transactionsService.GetTransactionsEntryByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve entry",
			"details": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, entry)
}

