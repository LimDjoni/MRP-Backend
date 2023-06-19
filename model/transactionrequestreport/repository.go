package transactionrequestreport

import (
	"ajebackend/model/master/iupopk"
	masterreport "command-line-arguments/Users/toktok2/Documents/toktok/Deli-AJE-Backend/model/masterreport/input.go"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateTransactionReport(input masterreport.TransactionReportInput, iupopkId int) (TransactionRequestReport, error)
	UpdateTransactionReport(DnDocumentLink string, LnDocumentLink string, id int, iupopkId int) (TransactionRequestReport, error)
	DeleteTransactionReport() (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateTransactionReport(input masterreport.TransactionReportInput, iupopkId int) (TransactionRequestReport, error) {
	var createTransactionReport TransactionRequestReport

	createTransactionReport.DateFrom = input.DateFrom
	createTransactionReport.DateTo = input.DateTo
	createTransactionReport.IupopkId = iupopkId

	errCreate := r.db.Create(&createTransactionReport).Error

	if errCreate != nil {
		return createTransactionReport, errCreate
	}

	var iup iupopk.Iupopk

	errFindIupopk := r.db.Where("id = ?", iupopkId).First(&iup).Error

	if errFindIupopk != nil {
		return createTransactionReport, errFindIupopk
	}

	year, month, _ := time.Now().Date()

	monthNumber := strconv.Itoa(int(month))

	if len([]rune(monthNumber)) < 2 {
		monthNumber = "0" + monthNumber
	}

	idNumber := fmt.Sprintf("TRR-%v-%v-%v-%v", iup.Code, monthNumber, year%1e2, createTransactionReport.ID)

	updateTransactionsReqErr := r.db.Model(&createTransactionReport).Where("id = ?", createTransactionReport.ID).Update("id_number", idNumber).Error

	if updateTransactionsReqErr != nil {
		return createTransactionReport, updateTransactionsReqErr
	}

	return createTransactionReport, nil
}

func (r *repository) UpdateTransactionReport(DnDocumentLink string, LnDocumentLink string, id int, iupopkId int) (TransactionRequestReport, error) {
	var updTransactionReport TransactionRequestReport

	errFind := r.db.Where("id = ? and iupopk_id = ?", id, iupopkId).First(&updTransactionReport).Error

	if errFind != nil {
		return updTransactionReport, errFind
	}

	updInput := make(map[string]interface{})

	updInput["document_dn_link"] = DnDocumentLink
	updInput["document_ln_link"] = LnDocumentLink

	updErr := r.db.Model(&r.UpdateTransactionReport).Updates(updInput).Error

	if updErr != nil {
		return updTransactionReport, updErr
	}

	return updTransactionReport, nil
}

func (r *repository) DeleteTransactionReport() (bool, error) {
	var transactionRequest TransactionRequestReport

	errFind := r.db.Find(&transactionRequest).Error

	if errFind != nil {
		return false, errFind
	}

	currentTime := time.Now()

	subtractedCurrentTime := currentTime.Add(-time.Hour * 24 * 7)
	errDelete := r.db.Unscoped().Where("updated_at <= ?", subtractedCurrentTime).Delete(&transactionRequest).Error

	if errDelete != nil {
		return false, errDelete
	}

	return true, nil
}
