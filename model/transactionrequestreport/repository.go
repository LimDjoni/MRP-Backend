package transactionrequestreport

import (
	"ajebackend/model/master/iupopk"
	"ajebackend/model/masterreport"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	PreviewTransactionReport(input masterreport.TransactionReportInput, iupopkId int) (TransactionRequestReportPreview, error)
	DetailTransactionReport(id int, iupopkId int) (TransactionRequestReportDetail, error)
	CreateTransactionReport(input masterreport.TransactionReportInput, iupopkId int) (TransactionRequestReport, error)
	UpdateTransactionReport(DnDocumentLink string, LnDocumentLink string, id int, iupopkId int) (TransactionRequestReport, error)
	UpdateTransactionReportError(id int, iupopkId int) (TransactionRequestReport, error)
	DeleteTransactionReport() (bool, error)
	DeleteTransactionReportById(id int, iupopkId int) (bool, error)
	ListTransactionReport(page int, iupopkId int) (Pagination, error)
	ListDeletedTransactionReport() ([]TransactionRequestReport, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) PreviewTransactionReport(input masterreport.TransactionReportInput, iupopkId int) (TransactionRequestReportPreview, error) {
	var previewTransactions TransactionRequestReportPreview

	var dnTransactions []masterreport.TransactionReport
	var lnTransactions []masterreport.TransactionReport
	errFindDn := r.db.Where("shipping_date >= ? and shipping_date <= ? and iupopk_id = ? and transaction_type = ?", input.DateFrom, input.DateTo, iupopkId, "DN").Find(&dnTransactions).Error

	if errFindDn != nil {
		return previewTransactions, errFindDn
	}

	errFindLn := r.db.Where("shipping_date >= ? and shipping_date <= ? and iupopk_id = ? and transaction_type = ?", input.DateFrom, input.DateTo, iupopkId, "LN").Find(&lnTransactions).Error

	if errFindLn != nil {
		return previewTransactions, errFindLn
	}

	previewTransactions.ListDnTransactions = dnTransactions
	previewTransactions.ListLnTransactions = lnTransactions

	return previewTransactions, nil
}

func (r *repository) DetailTransactionReport(id int, iupopkId int) (TransactionRequestReportDetail, error) {

	var detailTransaction TransactionRequestReportDetail
	var transactionReqDetail TransactionRequestReport

	errFindDetail := r.db.Where("id = ? and iupopk_id = ?", id, iupopkId).First(&transactionReqDetail).Error

	if errFindDetail != nil {
		return detailTransaction, errFindDetail
	}

	detailTransaction.Detail = transactionReqDetail

	var transactionsReportDn []masterreport.TransactionReport

	errFindDn := r.db.Table("transactions").Where("shipping_date >= ? and shipping_date <= ? and seller_id = ? and transaction_type = ?", transactionReqDetail.DateFrom, transactionReqDetail.DateTo, iupopkId, "DN").Find(&transactionsReportDn).Error

	if errFindDn != nil {
		return detailTransaction, errFindDn
	}

	detailTransaction.ListDnTransactions = transactionsReportDn

	var transactionsReportLn []masterreport.TransactionReport

	errFindLn := r.db.Table("transactions").Where("shipping_date >= ? and shipping_date <= ? and seller_id = ? and transaction_type = ?", transactionReqDetail.DateFrom, transactionReqDetail.DateTo, iupopkId, "LN").Find(&transactionsReportLn).Error

	if errFindLn != nil {
		return detailTransaction, errFindLn
	}

	detailTransaction.ListLnTransactions = transactionsReportLn

	return detailTransaction, nil
}

func (r *repository) CreateTransactionReport(input masterreport.TransactionReportInput, iupopkId int) (TransactionRequestReport, error) {

	var checkTransactionReport TransactionRequestReport

	errFindCheck := r.db.Where("date_from = ? and date_to = ? and iupopk_id = ?", input.DateFrom, input.DateTo, iupopkId).First(&checkTransactionReport).Error

	if errFindCheck == nil {
		return checkTransactionReport, errors.New("duplicate value")
	}

	var createTransactionReport TransactionRequestReport

	createTransactionReport.DateFrom = input.DateFrom
	createTransactionReport.DateTo = input.DateTo
	createTransactionReport.IupopkId = uint(iupopkId)
	createTransactionReport.Status = "In Progress"
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
	updInput["status"] = "Finished"

	updErr := r.db.Model(&updTransactionReport).Updates(updInput).Error

	if updErr != nil {
		return updTransactionReport, updErr
	}

	return updTransactionReport, nil
}

func (r *repository) UpdateTransactionReportError(id int, iupopkId int) (TransactionRequestReport, error) {
	var updTransactionReport TransactionRequestReport

	errFind := r.db.Where("id = ? and iupopk_id = ?", id, iupopkId).First(&updTransactionReport).Error

	if errFind != nil {
		return updTransactionReport, errFind
	}

	updInput := make(map[string]interface{})

	updInput["status"] = "Error"

	updErr := r.db.Model(&updTransactionReport).Updates(updInput).Error

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

func (r *repository) DeleteTransactionReportById(id int, iupopkId int) (bool, error) {
	var transactionRequest TransactionRequestReport

	errFind := r.db.Where("id = ? and iupopk_id = ?", id, iupopkId).First(&transactionRequest).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Unscoped().Where("id = ? and iupopk_id = ?", id, iupopkId).Delete(&transactionRequest).Error

	if errDelete != nil {
		return false, errDelete
	}

	return true, nil
}

func (r *repository) ListTransactionReport(page int, iupopkId int) (Pagination, error) {
	var listTransactionReport []TransactionRequestReport

	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	defaultSort := "created_at desc"

	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)

	errFind := r.db.Preload(clause.Associations).Order(defaultSort).Where(queryFilter).Scopes(paginateData(listTransactionReport, &pagination, r.db, queryFilter)).Find(&listTransactionReport).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listTransactionReport

	return pagination, nil
}

func (r *repository) ListDeletedTransactionReport() ([]TransactionRequestReport, error) {
	var listTransactionReport []TransactionRequestReport

	currentTime := time.Now()

	subtractedCurrentTime := currentTime.Add(-time.Hour * 24 * 7)

	errFind := r.db.Where("updated <= ?", subtractedCurrentTime).Find(&listTransactionReport).Error

	if errFind != nil {
		return listTransactionReport, errFind
	}

	return listTransactionReport, nil
}
