package services

import (
	"myfin-api/internal/dtos"
)

type CashHandlingService interface {
	CreateCashHandlingEntry(entry dtos.CreateCashHandlingEntryDTO) (dtos.CreateCashHandlingEntryDTO, error)
}

type cashHandlingService struct{}

func NewCashHandlingService() CashHandlingService {
	return &cashHandlingService{}
}

func (s *cashHandlingService) CreateCashHandlingEntry(entry dtos.CreateCashHandlingEntryDTO) (dtos.CreateCashHandlingEntryDTO, error) {
	var response = entry

	return response, nil
}
