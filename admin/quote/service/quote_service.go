package service

import (
	"fmt"
	"time"

	"github.com/rudyraharjo/emurojaah/admin/quote"
	"github.com/rudyraharjo/emurojaah/models"
)

type quoteService struct {
	quoteRepo quote.Repository
}

// NewQuoteService func
func NewQuoteService(repoQuote quote.Repository) quote.Service {
	return &quoteService{repoQuote}
}

func (s *quoteService) GetQuotesList() []models.Quote {

	res, err := s.quoteRepo.GetQuotesListAll()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return res
}

func (s *quoteService) AddQuote(quote models.Quote) ([]models.Quote, error) {

	now := time.Now()

	addQuote := models.Quote{
		Message:   quote.Message,
		Author:    quote.Author,
		IsActive:  1,
		CreatedAt: now,
	}

	data, err := s.quoteRepo.AddQuote(addQuote)
	if err != nil {
		fmt.Println(err)
		return data, err
	}

	return data, nil
}

func (s *quoteService) UpdateQuote(quote models.Quote) ([]models.Quote, error) {

	now := time.Now()

	updateQuote := models.Quote{
		Id:        quote.Id,
		Message:   quote.Message,
		Author:    quote.Author,
		IsActive:  1,
		UpdatedAt: now,
	}

	data, err := s.quoteRepo.UpdateQuote(updateQuote)

	if err != nil {
		fmt.Println(err)
		return data, err
	}
	return data, nil
}

func (s *quoteService) UpdateStatus(paramID int) error {

	err := s.quoteRepo.UpdateStatus(paramID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (s *quoteService) DeletedQuote(paramID int) error {

	err := s.quoteRepo.DeletedQuote(paramID)
	if err != nil {
		return err
	}

	return err

}

func (s *quoteService) GetQuoteByID(paramID int) ([]models.Quote, error) {
	data, err := s.quoteRepo.GetQuoteByID(paramID)
	if err != nil {
		return data, err
	}

	return data, nil
}
