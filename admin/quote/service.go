package quote

import "github.com/rudyraharjo/emurojaah/models"

// Service interface
type Service interface {
	GetQuoteByID(paramID int) ([]models.Quote, error)
	GetQuotesList() []models.Quote
	AddQuote(param models.Quote) ([]models.Quote, error)
	DeletedQuote(paramID int) error
	UpdateQuote(param models.Quote) ([]models.Quote, error)
	UpdateStatus(paramID int) error
}
