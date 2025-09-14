package services

import (
	"time"

	"myfin-api/internal/dtos"
	"myfin-api/internal/model"
	"myfin-api/internal/repository"
	"myfin-api/internal/repository/types"
)

const DateFormat = "02/01/2006"

type TransactionsService interface {
	CreateTransactionsEntry(entry dtos.CreateTransactionsEntryDTO) (dtos.TransactionsEntryResponseDTO, error)
	GetAllTransactionsEntries(limit, skip int, titleFilter, categoryFilter string) ([]dtos.TransactionsEntryResponseDTO, error)
	DeleteTransactionsEntry(id string) error
	UpdateTransactionsEntry(id string, entry dtos.UpdateTransactionsEntryDTO) (dtos.TransactionsEntryResponseDTO, error)
	GetTransactionsEntryByID(id string) (dtos.TransactionsEntryResponseDTO, error)
}

type transactionsService struct {
	transactionsRepo repository.TransactionsEntryRepository
}

func NewTransactionsService(transactionsRepo repository.TransactionsEntryRepository) TransactionsService {
	return &transactionsService{
		transactionsRepo: transactionsRepo,
	}
}

func (s *transactionsService) CreateTransactionsEntry(entry dtos.CreateTransactionsEntryDTO) (dtos.TransactionsEntryResponseDTO, error) {
	parsedDate, err := time.Parse(DateFormat, entry.Date)
	if err != nil {
		return dtos.TransactionsEntryResponseDTO{}, err
	}

	transactionsEntry := &model.TransactionsEntryModel{
		Amount:        entry.Amount,
		Title:         entry.Title,
		Currency:      entry.Currency,
		Type:          entry.Type,
		Category:      entry.Category,
		PaymentMethod: entry.PaymentMethod,
		Description:   entry.Description,
		Date:          parsedDate,
	}

	createdEntry, err := s.transactionsRepo.Create(transactionsEntry)
	if err != nil {
		return dtos.TransactionsEntryResponseDTO{}, err
	}

	response := dtos.TransactionsEntryResponseDTO{
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

func (s *transactionsService) GetAllTransactionsEntries(limit, skip int, titleFilter, categoryFilter string) ([]dtos.TransactionsEntryResponseDTO, error) {
	if limit < 0 {
		limit = 10
	}

	if limit > 100 {
		limit = 100
	}

	if skip < 0 {
		skip = 0
	}

	var entries []*model.TransactionsEntryModel
	var err error

	if titleFilter != "" || categoryFilter != "" {
		filter := types.FilterOptions{
			Title:    titleFilter,
			Category: categoryFilter,
		}
		entries, err = s.transactionsRepo.GetAllWithFilter(limit, skip, filter)
	} else {
		entries, err = s.transactionsRepo.GetAll(limit, skip)
	}

	if err != nil {
		return nil, err
	}

	response := make([]dtos.TransactionsEntryResponseDTO, 0, len(entries))
	for _, entry := range entries {
		response = append(response, dtos.TransactionsEntryResponseDTO{
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

func (s *transactionsService) DeleteTransactionsEntry(id string) error {
	return s.transactionsRepo.Delete(id)
}

func (s *transactionsService) UpdateTransactionsEntry(id string, entry dtos.UpdateTransactionsEntryDTO) (dtos.TransactionsEntryResponseDTO, error) {
	parsedDate, err := time.Parse(DateFormat, entry.Date)
	if err != nil {
		return dtos.TransactionsEntryResponseDTO{}, err
	}

	existingEntry, err := s.transactionsRepo.GetByID(id)
	if err != nil {
		return dtos.TransactionsEntryResponseDTO{}, err
	}

	transactionsEntry := &model.TransactionsEntryModel{
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

	updatedEntry, err := s.transactionsRepo.Update(id, transactionsEntry)
	if err != nil {
		return dtos.TransactionsEntryResponseDTO{}, err
	}

	response := dtos.TransactionsEntryResponseDTO{
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

func (s *transactionsService) GetTransactionsEntryByID(id string) (dtos.TransactionsEntryResponseDTO, error) {
	entry, err := s.transactionsRepo.GetByID(id)
	if err != nil {
		return dtos.TransactionsEntryResponseDTO{}, err
	}

	response := dtos.TransactionsEntryResponseDTO{
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
