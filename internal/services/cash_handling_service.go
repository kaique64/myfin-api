package services

import (
	"myfin-api/internal/dtos"
	"myfin-api/internal/model"
	"myfin-api/internal/repository"
	"time"
)

type CashHandlingService interface {
	CreateCashHandlingEntry(entry dtos.CreateCashHandlingEntryDTO) (dtos.CashHandlingEntryResponseDTO, error)
}

type cashHandlingService struct {
	cashHandlingRepo repository.CashHandlingEntryRepository
}

func NewCashHandlingService(cashHandlingRepo repository.CashHandlingEntryRepository) CashHandlingService {
	return &cashHandlingService{
		cashHandlingRepo: cashHandlingRepo,
	}
}

func (s *cashHandlingService) CreateCashHandlingEntry(entry dtos.CreateCashHandlingEntryDTO) (dtos.CashHandlingEntryResponseDTO, error) {
	parsedDate, err := time.Parse("02/01/2006", entry.Date)
	if err != nil {
		return dtos.CashHandlingEntryResponseDTO{}, err
	}

	cashHandlingEntry := &model.CashHandlingEntryModel{
		Amount:        entry.Amount,
		Currency:      entry.Currency,
		Type:          entry.Type,
		Category:      entry.Category,
		PaymentMethod: entry.PaymentMethod,
		Description:   entry.Description,
		Date:          parsedDate,
	}

	createdEntry, err := s.cashHandlingRepo.Create(cashHandlingEntry)
	if err != nil {
		return dtos.CashHandlingEntryResponseDTO{}, err
	}

	dateformat := "02/01/2006"

	response := dtos.CashHandlingEntryResponseDTO{
		ID:            createdEntry.ID.Hex(),
		Amount:        createdEntry.Amount,
		Currency:      createdEntry.Currency,
		Type:          createdEntry.Type,
		Category:      createdEntry.Category,
		PaymentMethod: createdEntry.PaymentMethod,
		Description:   createdEntry.Description,
		Date:          createdEntry.Date.Format(dateformat),
		Timestamp:     createdEntry.Timestamp,
		CreatedAt:     createdEntry.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     createdEntry.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return response, nil
}
