package repository

import (
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/rudyraharjo/emurojaah/admin/quote"
	"github.com/rudyraharjo/emurojaah/models"
)

// postgreQuoteRepository struct
type postgreQuoteRepository struct {
	DbConn *gorm.DB
}

// NewQuoteRepository DB
func NewQuoteRepository(DbConn *gorm.DB) quote.Repository {
	return &postgreQuoteRepository{DbConn}
}

func (repo *postgreQuoteRepository) GetQuotesListAll() ([]models.Quote, error) {
	var quote []models.Quote

	query := fmt.Sprintf("select * from quotes where deleted_at is null order by created_at desc")

	err := repo.DbConn.Raw(query).Scan(&quote).Error

	if err != nil {
		return nil, err
	}
	return quote, nil
}

func (repo *postgreQuoteRepository) AddQuote(quote models.Quote) ([]models.Quote, error) {
	var quoteSc []models.Quote

	err := repo.DbConn.Table("quotes").Create(&quote).Scan(&quoteSc).Error
	if err != nil {
		return quoteSc, err
	}

	return quoteSc, nil
}

func (repo *postgreQuoteRepository) UpdateQuote(quote models.Quote) ([]models.Quote, error) {

	var quoteSc []models.Quote

	err := repo.DbConn.Table("quotes").Where("id = ? ", quote.Id).Update(&quote).Scan(&quoteSc).Error
	if err != nil {
		fmt.Print(err)
		return quoteSc, err
	}
	return quoteSc, nil

}

func (repo *postgreQuoteRepository) GetQuoteByID(paramID int) ([]models.Quote, error) {

	var Quote []models.Quote

	err := repo.DbConn.Table("quotes").Where("id = ?", paramID).Scan(&Quote).Error

	if err != nil {
		return nil, err
	}
	return Quote, nil

}

func (repo *postgreQuoteRepository) UpdateStatus(paramID int) error {
	var quoteSc []models.Quote

	err := repo.DbConn.Raw(`update quotes set is_active = case when (is_active = 1) then 0 else 1 end where id = ?`, paramID).Scan(&quoteSc).Error

	if err != nil {
		return err
	}
	return nil
}

func (repo *postgreQuoteRepository) DeletedQuote(paramID int) error {

	var quoteSc []models.Quote
	now := time.Now()

	err := repo.DbConn.Raw(`update quotes set deleted_at = ? where id = ?`, now, paramID).Scan(&quoteSc).Error

	if err != nil {
		return err
	}
	return err
}
