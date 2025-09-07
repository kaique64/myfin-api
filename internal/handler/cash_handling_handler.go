package handlers

import (
	"myfin-api/internal/dtos"
	"myfin-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CashHandlingHandler interface {
	Save(ctx *gin.Context)
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
