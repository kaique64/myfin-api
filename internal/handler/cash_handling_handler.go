package handlers

import (
	"myfin-api/internal/dtos"
	"myfin-api/internal/services"

	"github.com/gin-gonic/gin"
)

type CashHandlingHandler interface {
	Save(ctx *gin.Context) dtos.CreateCashHandlingEntryDTO
}

type cashHandlingHandler struct {
	cashHandlingService services.CashHandlingService
}

func NewCashHandlingHandler(cashHandlingService services.CashHandlingService) CashHandlingHandler {
	return &cashHandlingHandler{
		cashHandlingService: cashHandlingService,
	}
}

func (h *cashHandlingHandler) Save(ctx *gin.Context) dtos.CreateCashHandlingEntryDTO {
	var entry dtos.CreateCashHandlingEntryDTO

	ctx.BindJSON(&entry)

	h.cashHandlingService.CreateCashHandlingEntry(entry)

	return entry
}
