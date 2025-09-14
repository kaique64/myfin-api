package services

import (
	"time"

	"myfin-api/internal/dtos"
	"myfin-api/internal/model"
	"myfin-api/internal/repository"
	"myfin-api/internal/repository/types"
)

const DateFormat = "02/01/2006"

type CashHandlingService interface {
	CreateCashHandlingEntry(entry dtos.CreateCashHandlingEntryDTO) (dtos.CashHandlingEntryResponseDTO, error)
	GetAllCashHandlingEntries(limit, skip int, titleFilter, categoryFilter string) ([]dtos.CashHandlingEntryResponseDTO, error)
	DeleteCashHandlingEntry(id string) error
	UpdateCashHandlingEntry(id string, entry dtos.UpdateCashHandlingEntryDTO) (dtos.CashHandlingEntryResponseDTO, error)
	GetCashHandlingEntryByID(id string) (dtos.CashHandlingEntryResponseDTO, error)
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
	parsedDate, err := time.Parse(DateFormat, entry.Date)
	if err != nil {
		return dtos.CashHandlingEntryResponseDTO{}, err
	}

	cashHandlingEntry := &model.CashHandlingEntryModel{
		Amount:        entry.Amount,
		Title:         entry.Title,
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

	response := dtos.CashHandlingEntryResponseDTO{
		ID:            createdEntry.ID.Hex(),
		Amount:        createdEntry.Amount,
		Title:         createdEntry.Title,
		Currency:      createdEntry.Currency,
		Type:          createdEntry.Type,
		Category:      createdEntry.Category,
		PaymentMethod: createdEntry.PaymentMethod,
		Description:   createdEntry.Description,
		Date:          createdEntry.Date.Format(DateFormat),
		Timestamp:     createdEntry.Timestamp,
		CreatedAt:     createdEntry.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     createdEntry.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return response, nil
}

func (s *cashHandlingService) GetAllCashHandlingEntries(limit, skip int, titleFilter, categoryFilter string) ([]dtos.CashHandlingEntryResponseDTO, error) {
	if limit < 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	if skip < 0 {
		skip = 0
	}

	var entries []*model.CashHandlingEntryModel
	var err error

	if titleFilter != "" || categoryFilter != "" {
		filter := types.FilterOptions{
			Title:    titleFilter,
			Category: categoryFilter,
		}
		entries, err = s.cashHandlingRepo.GetAllWithFilter(limit, skip, filter)
	} else {
		entries, err = s.cashHandlingRepo.GetAll(limit, skip)
	}

	if err != nil {
		return nil, err
	}

	response := make([]dtos.CashHandlingEntryResponseDTO, 0, len(entries))
	for _, entry := range entries {
		response = append(response, dtos.CashHandlingEntryResponseDTO{
			ID:            entry.ID.Hex(),
			Amount:        entry.Amount,
			Title:         entry.Title,
			Currency:      entry.Currency,
			Type:          entry.Type,
			Category:      entry.Category,
			PaymentMethod: entry.PaymentMethod,
			Description:   entry.Description,
			Date:          entry.Date.Format(DateFormat),
			Timestamp:     entry.Timestamp,
			CreatedAt:     entry.CreatedAt.UTC().Format(time.RFC3339),
			UpdatedAt:     entry.UpdatedAt.UTC().Format(time.RFC3339),
		})
	}

	return response, nil
}

func (s *cashHandlingService) DeleteCashHandlingEntry(id string) error {
	return s.cashHandlingRepo.Delete(id)
}

func (s *cashHandlingService) UpdateCashHandlingEntry(id string, entry dtos.UpdateCashHandlingEntryDTO) (dtos.CashHandlingEntryResponseDTO, error) {
	parsedDate, err := time.Parse(DateFormat, entry.Date)
	if err != nil {
		return dtos.CashHandlingEntryResponseDTO{}, err
	}

	existingEntry, err := s.cashHandlingRepo.GetByID(id)
	if err != nil {
		return dtos.CashHandlingEntryResponseDTO{}, err
	}

	cashHandlingEntry := &model.CashHandlingEntryModel{
		Amount:        entry.Amount,
		Title:         entry.Title,
		Currency:      entry.Currency,
		Type:          entry.Type,
		Category:      entry.Category,
		PaymentMethod: entry.PaymentMethod,
		Description:   entry.Description,
		Date:          parsedDate,
		Timestamp:     existingEntry.Timestamp,
		CreatedAt:     existingEntry.CreatedAt,
	}

	updatedEntry, err := s.cashHandlingRepo.Update(id, cashHandlingEntry)
	if err != nil {
		return dtos.CashHandlingEntryResponseDTO{}, err
	}

	response := dtos.CashHandlingEntryResponseDTO{
		ID:            updatedEntry.ID.Hex(),
		Amount:        updatedEntry.Amount,
		Title:         updatedEntry.Title,
		Currency:      updatedEntry.Currency,
		Type:          updatedEntry.Type,
		Category:      updatedEntry.Category,
		PaymentMethod: updatedEntry.PaymentMethod,
		Description:   updatedEntry.Description,
		Date:          updatedEntry.Date.Format(DateFormat),
		Timestamp:     updatedEntry.Timestamp,
		CreatedAt:     updatedEntry.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     updatedEntry.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return response, nil
}

func (s *cashHandlingService) GetCashHandlingEntryByID(id string) (dtos.CashHandlingEntryResponseDTO, error) {
	entry, err := s.cashHandlingRepo.GetByID(id)
	if err != nil {
		return dtos.CashHandlingEntryResponseDTO{}, err
	}

	response := dtos.CashHandlingEntryResponseDTO{
		ID:            entry.ID.Hex(),
		Amount:        entry.Amount,
		Title:         entry.Title,
		Currency:      entry.Currency,
		Type:          entry.Type,
		Category:      entry.Category,
		PaymentMethod: entry.PaymentMethod,
		Description:   entry.Description,
		Date:          entry.Date.Format(DateFormat),
		Timestamp:     entry.Timestamp,
		CreatedAt:     entry.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:     entry.UpdatedAt.UTC().Format(time.RFC3339),
	}

	return response, nil
}
