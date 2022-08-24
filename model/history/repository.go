package history

import (
	"ajebackend/helper"
	"ajebackend/model/transaction"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Repository interface {
	CreateTransactionDN (inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	DeleteTransactionDN(id int, userId uint) (bool, error)
	UpdateTransactionDN (idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	UploadDocument (idTransaction uint, urlS3 string, userId uint, documentType string) (transaction.Transaction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateTransactionDN (inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	var createdTransaction transaction.Transaction
	var totalCount int64
	year, month, _ := time.Now().Date()
	startDate := fmt.Sprintf("%v-%v-01  00:00:00", year, int(month))
	endDate := fmt.Sprintf("%v-%v-31  00:00:00", year, int(month))

	tx := r.db.Begin()

	findErr := tx.Unscoped().Where("created_at BETWEEN ? AND ?", startDate, endDate).Model(transaction.Transaction{}).Count(&totalCount).Error

	if findErr != nil {
		return createdTransaction, findErr
	}


	if inputTransactionDN.DpRoyaltyCurrency == "" {
		inputTransactionDN.DpRoyaltyCurrency = "IDR"
	}

	if inputTransactionDN.PaymentDpRoyaltyCurrency == "" {
		inputTransactionDN.PaymentDpRoyaltyCurrency = "IDR"
	}

	inputTransactionJson, errMarshal := json.Marshal(inputTransactionDN)

	if errMarshal != nil {
		return  createdTransaction, errMarshal
	}
	inputMap := make(map[string]interface{})
	errUnmarshal := json.Unmarshal([]byte(inputTransactionJson), &inputMap)

	if errUnmarshal != nil {
		return  createdTransaction, errUnmarshal
	}

	inputMap["dmo_id"] = nil
	inputMap["id_number"] = fmt.Sprintf("DN-%v-%v-%v", year, int(month), helper.CreateIdNumber(int(totalCount + 1)))
	inputMap["transaction_type"] = "DN"

	createTransactionErr := tx.Model(&createdTransaction).Create(inputMap).Error

	if createTransactionErr != nil {
		tx.Rollback()
		return createdTransaction, createTransactionErr
	}

	var history History

	history.TransactionId = &createdTransaction.ID
	history.Status = "Created"
	history.UserId = userId

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return createdTransaction, createHistoryErr
	}

	tx.Commit()
	return createdTransaction, createHistoryErr
}

func (r *repository) DeleteTransactionDN(id int, userId uint) (bool, error) {
	var transaction transaction.Transaction

	tx := r.db.Begin()

	errFind := r.db.Where("id = ?", id).First(&transaction).Error

	if errFind != nil {
		return false, errFind
	}

	errDelete := r.db.Where("id = ?", id).Delete(&transaction).Error

	if errDelete != nil {
		tx.Rollback()
		return false, errDelete
	}

	var history History

	uId := uint(id)
	history.TransactionId = &uId
	history.UserId = userId
	history.Status = "Deleted"

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return false, createHistoryErr
	}

	tx.Commit()
	return true, createHistoryErr
}

func (r *repository) UpdateTransactionDN (idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	var transaction transaction.Transaction

	tx := r.db.Begin()

	errFind := r.db.Where("id = ?", idTransaction).First(&transaction).Error

	if errFind != nil {
		tx.Rollback()
		return transaction, errFind
	}

	beforeData , errorBeforeDataJsonMarshal := json.Marshal(transaction)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return transaction, errorBeforeDataJsonMarshal
	}

	dataInput, errorMarshal := json.Marshal(inputEditTransactionDN)

	if errorMarshal != nil {
		tx.Rollback()
		return  transaction, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		tx.Rollback()
		return  transaction, errorUnmarshal
	}

	updateErr := tx.Model(&transaction).Updates(dataInputMapString).Error

	if updateErr != nil {
		tx.Rollback()
		return  transaction, updateErr
	}

	afterData , errorAfterDataJsonMarshal := json.Marshal(transaction)

	if errorBeforeDataJsonMarshal != nil {
		tx.Rollback()
		return transaction, errorAfterDataJsonMarshal
	}

	var history History

	history.TransactionId = &transaction.ID
	history.UserId = userId
	history.Status = "Updated"
	history.BeforeData = beforeData
	history.AfterData = afterData

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return transaction, createHistoryErr
	}

	tx.Commit()
	return transaction, nil
}

func (r *repository) UploadDocument (idTransaction uint, urlS3 string, userId uint, documentType string) (transaction.Transaction, error) {
	var uploadedTransaction transaction.Transaction

	tx := r.db.Begin()

	errFind := tx.Where("id = ?", idTransaction).First(&uploadedTransaction).Error

	if errFind != nil {
		return uploadedTransaction, errFind
	}
	var isReupload = false
	editData := make(map[string]interface{})

	switch documentType {
		case "skb":
			if uploadedTransaction.SkbDocumentLink != "" {
				isReupload = true
			}
			editData["skb_document_link"] = urlS3
		case "skab":
			if uploadedTransaction.SkabDocumentLink != "" {
				isReupload = true
			}
			editData["skab_document_link"] = urlS3
		case "bl":
			if uploadedTransaction.BLDocumentLink != "" {
				isReupload = true
			}
			editData["bl_document_link"] = urlS3
		case "royalti_provision":
			if uploadedTransaction.RoyaltiProvisionDocumentLink != "" {
				isReupload = true
			}
			editData["royalti_provision_document_link"] = urlS3
		case "royalti_final":
			if uploadedTransaction.RoyaltiFinalDocumentLink != "" {
				isReupload = true
			}
			editData["royalti_final_document_link"] = urlS3
		case "cow":
			if uploadedTransaction.COWDocumentLink != "" {
				isReupload = true
			}
			editData["cow_document_link"] = urlS3
		case "coa":
			if uploadedTransaction.COADocumentLink != "" {
				isReupload = true
			}
			editData["coa_document_link"] = urlS3
		case "invoice":
			if uploadedTransaction.InvoiceAndContractDocumentLink != "" {
				isReupload = true
			}
			editData["invoice_and_contract_document_link"] = urlS3
		case "lhv":
			if uploadedTransaction.LHVDocumentLink != "" {
				isReupload = true
			}
			editData["lhv_document_link"] = urlS3
	}

	errEdit := tx.Model(&uploadedTransaction).Updates(editData).Error

	if errEdit != nil {
		tx.Rollback()
		return uploadedTransaction, errEdit
	}
	var history History

	history.TransactionId = &uploadedTransaction.ID
	history.UserId = userId
	if isReupload == false {
		history.Status = fmt.Sprintf("Uploaded %s document", documentType)
	}

	if isReupload == true {
		history.Status = fmt.Sprintf("Reupload %s document", documentType)
	}

	createHistoryErr := tx.Create(&history).Error

	if createHistoryErr != nil {
		tx.Rollback()
		return uploadedTransaction, createHistoryErr
	}

	tx.Commit()
	return uploadedTransaction, nil
}
