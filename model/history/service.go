package history

import (
	"ajebackend/model/dmo"
	"ajebackend/model/minerba"
	"ajebackend/model/production"
	"ajebackend/model/transaction"
)

type Service interface {
	CreateTransactionDN(inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	DeleteTransactionDN(id int, userId uint) (bool, error)
	UpdateTransactionDN(idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	UploadDocumentTransactionDN(idTransaction uint, urlS3 string, userId uint, documentType string) (transaction.Transaction, error)
	CreateMinerba(period string, baseIdNumber string, updateTransaction []int, userId uint) (minerba.Minerba, error)
	UpdateMinerba(id int, updateTransaction []int, userId uint) (minerba.Minerba, error)
	DeleteMinerba(idMinerba int, userId uint) (bool, error)
	UpdateDocumentMinerba(id int, documentLink minerba.InputUpdateDocumentMinerba, userId uint) (minerba.Minerba, error)
	CreateDmo(dmoInput dmo.CreateDmoInput, baseIdNumber string, userId uint) (dmo.Dmo, error)
	DeleteDmo(idDmo int, userId uint) (bool, error)
	UpdateDocumentDmo(id int, documentLink dmo.InputUpdateDocumentDmo, userId uint) (dmo.Dmo, error)
	UpdateIsDownloadedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint) (dmo.Dmo, error)
	UpdateTrueIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, location string) (dmo.Dmo, error)
	UpdateFalseIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint) (dmo.Dmo, error)
	CreateProduction(input production.InputCreateProduction, userId uint) (production.Production, error)
	UpdateProduction(input production.InputCreateProduction, productionId int, userId uint) (production.Production, error)
	DeleteProduction(productionId int, userId uint) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateTransactionDN(inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	transaction, transactionErr := s.repository.CreateTransactionDN(inputTransactionDN, userId)

	return transaction, transactionErr
}

func (s *service) DeleteTransactionDN(id int, userId uint) (bool, error) {
	isDeletedTransaction, isDeletedTransactionErr := s.repository.DeleteTransactionDN(id, userId)

	return isDeletedTransaction, isDeletedTransactionErr
}

func (s *service) UpdateTransactionDN(idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	updateTransaction, updateTransactionErr := s.repository.UpdateTransactionDN(idTransaction, inputEditTransactionDN, userId)

	return updateTransaction, updateTransactionErr
}

func (s *service) UploadDocumentTransactionDN(idTransaction uint, urlS3 string, userId uint, documentType string) (transaction.Transaction, error) {
	uploadedDocument, uploadedDocumentErr := s.repository.UploadDocumentTransactionDN(idTransaction, urlS3, userId, documentType)

	return uploadedDocument, uploadedDocumentErr
}

func (s *service) CreateMinerba(period string, baseIdNumber string, updateTransaction []int, userId uint) (minerba.Minerba, error) {
	createdMinerba, createdMinerbaErr := s.repository.CreateMinerba(period, baseIdNumber, updateTransaction, userId)

	return createdMinerba, createdMinerbaErr
}

func (s *service) UpdateMinerba(id int, updateTransaction []int, userId uint) (minerba.Minerba, error) {
	updatedMinerba, updatedMinerbaErr := s.repository.UpdateMinerba(id, updateTransaction, userId)

	return updatedMinerba, updatedMinerbaErr
}

func (s *service) DeleteMinerba(idMinerba int, userId uint) (bool, error) {
	isDeletedMinerba, isDeletedMinerbaErr := s.repository.DeleteMinerba(idMinerba, userId)

	return isDeletedMinerba, isDeletedMinerbaErr
}

func (s *service) UpdateDocumentMinerba(id int, documentLink minerba.InputUpdateDocumentMinerba, userId uint) (minerba.Minerba, error) {
	uploadMinerba, uploadMinerbaErr := s.repository.UpdateDocumentMinerba(id, documentLink, userId)

	return uploadMinerba, uploadMinerbaErr
}

func (s *service) CreateDmo(dmoInput dmo.CreateDmoInput, baseIdNumber string, userId uint) (dmo.Dmo, error) {
	createDmo, createDmoErr := s.repository.CreateDmo(dmoInput, baseIdNumber, userId)

	return createDmo, createDmoErr
}

func (s *service) DeleteDmo(idDmo int, userId uint) (bool, error) {
	isDeletedDmo, isDeletedDmoErr := s.repository.DeleteDmo(idDmo, userId)

	return isDeletedDmo, isDeletedDmoErr
}

func (s *service) UpdateDocumentDmo(id int, documentLink dmo.InputUpdateDocumentDmo, userId uint) (dmo.Dmo, error) {
	updateDocumentDmo, updateDocumentDmoErr := s.repository.UpdateDocumentDmo(id, documentLink, userId)

	return updateDocumentDmo, updateDocumentDmoErr
}

func (s *service) UpdateIsDownloadedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint) (dmo.Dmo, error) {
	updateIsDownloadedDmoDocument, updateIsDownloadedDmoDocumentErr := s.repository.UpdateIsDownloadedDmoDocument(isBast, isStatementLetter, isReconciliationLetter, isReconciliationLetterEndUser, id, userId)

	return updateIsDownloadedDmoDocument, updateIsDownloadedDmoDocumentErr
}

func (s *service) UpdateTrueIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, location string) (dmo.Dmo, error) {
	updateIsSignedDmoDocument, updateIsSignedDmoDocumentErr := s.repository.UpdateTrueIsSignedDmoDocument(isBast, isStatementLetter, isReconciliationLetter, isReconciliationLetterEndUser, id, userId, location)

	return updateIsSignedDmoDocument, updateIsSignedDmoDocumentErr
}

func (s *service) UpdateFalseIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint) (dmo.Dmo, error) {
	updateIsSignedDmoDocument, updateIsSignedDmoDocumentErr := s.repository.UpdateFalseIsSignedDmoDocument(isBast, isStatementLetter, isReconciliationLetter, isReconciliationLetterEndUser, id, userId)

	return updateIsSignedDmoDocument, updateIsSignedDmoDocumentErr
}

func (s *service) CreateProduction(input production.InputCreateProduction, userId uint) (production.Production, error) {
	createProduction, createProductionErr := s.repository.CreateProduction(input, userId)

	return createProduction, createProductionErr
}

func (s *service) UpdateProduction(input production.InputCreateProduction, productionId int, userId uint) (production.Production, error) {
	updateProduction, updateProductionErr := s.repository.UpdateProduction(input, productionId, userId)

	return updateProduction, updateProductionErr
}

func (s *service) DeleteProduction(productionId int, userId uint) (bool, error) {
	deleteProduction, deleteProductionErr := s.repository.DeleteProduction(productionId, userId)

	return deleteProduction, deleteProductionErr
}
