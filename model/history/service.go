package history

import (
	"ajebackend/model/dmo"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/insw"
	"ajebackend/model/minerba"
	"ajebackend/model/minerbaln"
	"ajebackend/model/production"
	"ajebackend/model/reportdmo"
	"ajebackend/model/transaction"
)

type Service interface {
	CreateTransactionDN(inputTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	DeleteTransaction(id int, userId uint, transactionType string) (bool, error)
	UpdateTransactionDN(idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	UploadDocumentTransaction(idTransaction uint, urlS3 string, userId uint, documentType string, transactionType string) (transaction.Transaction, error)
	CreateTransactionLN(inputTransactionLN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
	UpdateTransactionLN(id int, inputTransactionLN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error)
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
	UpdateDmo(dmoUpdateInput dmo.UpdateDmoInput, id int, userId uint) (dmo.Dmo, error)
	CreateProduction(input production.InputCreateProduction, userId uint) (production.Production, error)
	UpdateProduction(input production.InputCreateProduction, productionId int, userId uint) (production.Production, error)
	DeleteProduction(productionId int, userId uint) (bool, error)
	CreateGroupingVesselDN(inputGrouping groupingvesseldn.InputGroupingVesselDn, userId uint) (groupingvesseldn.GroupingVesselDn, error)
	EditGroupingVesselDn(id int, editGrouping groupingvesseldn.InputEditGroupingVesselDn, userId uint) (groupingvesseldn.GroupingVesselDn, error)
	DeleteGroupingVesselDn(id int, userId uint) (bool, error)
	UploadDocumentGroupingVesselDn(id uint, urlS3 string, userId uint, documentType string) (groupingvesseldn.GroupingVesselDn, error)
	CreateGroupingVesselLN(inputGrouping groupingvesselln.InputGroupingVesselLn, userId uint) (groupingvesselln.GroupingVesselLn, error)
	EditGroupingVesselLn(id int, editGrouping groupingvesselln.InputEditGroupingVesselLn, userId uint) (groupingvesselln.GroupingVesselLn, error)
	UploadDocumentGroupingVesselLn(id uint, urlS3 string, userId uint, documentType string) (groupingvesselln.GroupingVesselLn, error)
	DeleteGroupingVesselLn(id int, userId uint) (bool, error)
	CreateMinerbaLn(period string, baseIdNumber string, listTransactions []int, userId uint) (minerbaln.MinerbaLn, error)
	UpdateMinerbaLn(id int, listTransactions []int, userId uint) (minerbaln.MinerbaLn, error)
	DeleteMinerbaLn(idMinerbaLn int, userId uint) (bool, error)
	UpdateDocumentMinerbaLn(id int, documentLink minerbaln.InputUpdateDocumentMinerbaLn, userId uint) (minerbaln.MinerbaLn, error)
	CreateInsw(month string, year int, baseIdNumber string, userId uint) (insw.Insw, error)
	DeleteInsw(idInsw int, userId uint) (bool, error)
	UpdateDocumentInsw(id int, documentLink insw.InputUpdateDocumentInsw, userId uint) (insw.Insw, error)
	CreateReportDmo(input reportdmo.InputCreateReportDmo, baseIdNumber string, userId uint) (reportdmo.ReportDmo, error)
	UpdateDocumentReportDmo(id int, documentLink reportdmo.InputUpdateDocumentReportDmo, userId uint) (reportdmo.ReportDmo, error)
	UpdateTransactionReportDmo(id int, inputUpdate reportdmo.InputUpdateReportDmo, userId uint) (reportdmo.ReportDmo, error)
	DeleteReportDmo(idReportDmo int, userId uint) (bool, error)
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

func (s *service) DeleteTransaction(id int, userId uint, transactionType string) (bool, error) {
	isDeletedTransaction, isDeletedTransactionErr := s.repository.DeleteTransaction(id, userId, transactionType)

	return isDeletedTransaction, isDeletedTransactionErr
}

func (s *service) UpdateTransactionDN(idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	updateTransaction, updateTransactionErr := s.repository.UpdateTransactionDN(idTransaction, inputEditTransactionDN, userId)

	return updateTransaction, updateTransactionErr
}

func (s *service) UploadDocumentTransaction(idTransaction uint, urlS3 string, userId uint, documentType string, transactionType string) (transaction.Transaction, error) {
	uploadedDocument, uploadedDocumentErr := s.repository.UploadDocumentTransaction(idTransaction, urlS3, userId, documentType, transactionType)

	return uploadedDocument, uploadedDocumentErr
}

func (s *service) CreateTransactionLN(inputTransactionLN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	transactionLn, transactionLnErr := s.repository.CreateTransactionLN(inputTransactionLN, userId)

	return transactionLn, transactionLnErr
}

func (s *service) UpdateTransactionLN(id int, inputTransactionLN transaction.DataTransactionInput, userId uint) (transaction.Transaction, error) {
	transactionLn, transactionLnErr := s.repository.UpdateTransactionLN(id, inputTransactionLN, userId)

	return transactionLn, transactionLnErr
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

func (s *service) UpdateDmo(dmoUpdateInput dmo.UpdateDmoInput, id int, userId uint) (dmo.Dmo, error) {
	updateDmo, updateDmoErr := s.repository.UpdateDmo(dmoUpdateInput, id, userId)

	return updateDmo, updateDmoErr
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

func (s *service) CreateGroupingVesselDN(inputGrouping groupingvesseldn.InputGroupingVesselDn, userId uint) (groupingvesseldn.GroupingVesselDn, error) {
	createGroupingVesselDn, errCreateGroupingVesselDn := s.repository.CreateGroupingVesselDN(inputGrouping, userId)

	return createGroupingVesselDn, errCreateGroupingVesselDn
}

func (s *service) EditGroupingVesselDn(id int, editGrouping groupingvesseldn.InputEditGroupingVesselDn, userId uint) (groupingvesseldn.GroupingVesselDn, error) {
	editGroupingVesselDn, errEditGroupingVesselDn := s.repository.EditGroupingVesselDn(id, editGrouping, userId)

	return editGroupingVesselDn, errEditGroupingVesselDn
}

func (s *service) DeleteGroupingVesselDn(id int, userId uint) (bool, error) {
	deleteGroupingVesselDn, errDeleteGroupingVesselDn := s.repository.DeleteGroupingVesselDn(id, userId)

	return deleteGroupingVesselDn, errDeleteGroupingVesselDn
}

func (s *service) UploadDocumentGroupingVesselDn(id uint, urlS3 string, userId uint, documentType string) (groupingvesseldn.GroupingVesselDn, error) {
	uploadDocumentGroupingVesselDn, errUploadDocumentGroupingVesselDn := s.repository.UploadDocumentGroupingVesselDn(id, urlS3, userId, documentType)

	return uploadDocumentGroupingVesselDn, errUploadDocumentGroupingVesselDn
}

func (s *service) CreateGroupingVesselLN(inputGrouping groupingvesselln.InputGroupingVesselLn, userId uint) (groupingvesselln.GroupingVesselLn, error) {
	createGroupingVesselLn, errCreateGroupingVesselLn := s.repository.CreateGroupingVesselLN(inputGrouping, userId)

	return createGroupingVesselLn, errCreateGroupingVesselLn
}

func (s *service) EditGroupingVesselLn(id int, editGrouping groupingvesselln.InputEditGroupingVesselLn, userId uint) (groupingvesselln.GroupingVesselLn, error) {
	editGroupingVesselLn, errEditGroupingVesselLn := s.repository.EditGroupingVesselLn(id, editGrouping, userId)

	return editGroupingVesselLn, errEditGroupingVesselLn
}

func (s *service) UploadDocumentGroupingVesselLn(id uint, urlS3 string, userId uint, documentType string) (groupingvesselln.GroupingVesselLn, error) {
	uploadDocumentGroupingVesselLn, uploadDocumentGroupingVesselLnErr := s.repository.UploadDocumentGroupingVesselLn(id, urlS3, userId, documentType)

	return uploadDocumentGroupingVesselLn, uploadDocumentGroupingVesselLnErr
}

func (s *service) DeleteGroupingVesselLn(id int, userId uint) (bool, error) {
	deleteGroupingVesselLn, deleteGroupingVesselLnErr := s.repository.DeleteGroupingVesselLn(id, userId)

	return deleteGroupingVesselLn, deleteGroupingVesselLnErr
}

func (s *service) CreateMinerbaLn(period string, baseIdNumber string, listTransactions []int, userId uint) (minerbaln.MinerbaLn, error) {
	createMinerbaLn, createMinerbaLnErr := s.repository.CreateMinerbaLn(period, baseIdNumber, listTransactions, userId)

	return createMinerbaLn, createMinerbaLnErr
}

func (s *service) UpdateMinerbaLn(id int, listTransactions []int, userId uint) (minerbaln.MinerbaLn, error) {
	updateMinerbaLn, updateMinerbaLnErr := s.repository.UpdateMinerbaLn(id, listTransactions, userId)

	return updateMinerbaLn, updateMinerbaLnErr
}

func (s *service) DeleteMinerbaLn(idMinerbaLn int, userId uint) (bool, error) {
	isDeletedMinerbaLn, isDeletedMinerbaLnErr := s.repository.DeleteMinerbaLn(idMinerbaLn, userId)

	return isDeletedMinerbaLn, isDeletedMinerbaLnErr
}

func (s *service) UpdateDocumentMinerbaLn(id int, documentLink minerbaln.InputUpdateDocumentMinerbaLn, userId uint) (minerbaln.MinerbaLn, error) {
	uploadMinerbaLn, uploadMinerbaLnErr := s.repository.UpdateDocumentMinerbaLn(id, documentLink, userId)

	return uploadMinerbaLn, uploadMinerbaLnErr
}

func (s *service) CreateInsw(month string, year int, baseIdNumber string, userId uint) (insw.Insw, error) {
	createInsw, createInswErr := s.repository.CreateInsw(month, year, baseIdNumber, userId)

	return createInsw, createInswErr
}

func (s *service) DeleteInsw(idInsw int, userId uint) (bool, error) {
	deleteInsw, deleteInswErr := s.repository.DeleteInsw(idInsw, userId)

	return deleteInsw, deleteInswErr
}

func (s *service) UpdateDocumentInsw(id int, documentLink insw.InputUpdateDocumentInsw, userId uint) (insw.Insw, error) {
	updateDocumentInsw, updateDocumentInswErr := s.repository.UpdateDocumentInsw(id, documentLink, userId)

	return updateDocumentInsw, updateDocumentInswErr
}

func (s *service) CreateReportDmo(input reportdmo.InputCreateReportDmo, baseIdNumber string, userId uint) (reportdmo.ReportDmo, error) {
	createReportDmo, createReportDmoErr := s.repository.CreateReportDmo(input, baseIdNumber, userId)

	return createReportDmo, createReportDmoErr
}

func (s *service) UpdateDocumentReportDmo(id int, documentLink reportdmo.InputUpdateDocumentReportDmo, userId uint) (reportdmo.ReportDmo, error) {
	updateDocumentReportDmo, updateDocumentReportDmoErr := s.repository.UpdateDocumentReportDmo(id, documentLink, userId)

	return updateDocumentReportDmo, updateDocumentReportDmoErr
}

func (s *service) UpdateTransactionReportDmo(id int, inputUpdate reportdmo.InputUpdateReportDmo, userId uint) (reportdmo.ReportDmo, error) {
	updateTransactionReportDmo, updateTransactionReportDmoErr := s.repository.UpdateTransactionReportDmo(id, inputUpdate, userId)

	return updateTransactionReportDmo, updateTransactionReportDmoErr
}

func (s *service) DeleteReportDmo(idReportDmo int, userId uint) (bool, error) {
	deleteReportDmo, deleteReportDmoErr := s.repository.DeleteReportDmo(idReportDmo, userId)

	return deleteReportDmo, deleteReportDmoErr
}
