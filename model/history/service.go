package history

import (
	"ajebackend/model/cafassignment"
	"ajebackend/model/cafassignmentenduser"
	"ajebackend/model/coareport"
	"ajebackend/model/coareportln"
	"ajebackend/model/dmo"
	"ajebackend/model/electricassignment"
	"ajebackend/model/electricassignmentenduser"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/insw"
	"ajebackend/model/jettybalance"
	"ajebackend/model/minerba"
	"ajebackend/model/minerbaln"
	"ajebackend/model/pitloss"
	"ajebackend/model/production"
	"ajebackend/model/reportdmo"
	"ajebackend/model/rkab"
	"ajebackend/model/royaltyrecon"
	"ajebackend/model/royaltyreport"
	"ajebackend/model/transaction"
)

type Service interface {
	CreateTransactionDN(inputTransactionDN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error)
	DeleteTransaction(id int, userId uint, transactionType string, iupopkId int) (bool, error)
	UpdateTransactionDN(idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error)
	UploadDocumentTransaction(idTransaction uint, urlS3 string, userId uint, documentType string, transactionType string, iupopkId int) (transaction.Transaction, error)
	DeleteDocumentTransaction(idTransaction uint, userId uint, documentType string, transactionType string, iupopkId int) (transaction.Transaction, error)
	CreateTransactionLN(inputTransactionLN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error)
	UpdateTransactionLN(id int, inputTransactionLN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error)
	CreateMinerba(period string, updateTransaction []int, userId uint, iupopkId int) (minerba.Minerba, error)
	UpdateMinerba(id int, updateTransaction []int, userId uint, iupopkId int) (minerba.Minerba, error)
	DeleteMinerba(idMinerba int, userId uint, iupopkId int) (bool, error)
	UpdateDocumentMinerba(id int, documentLink minerba.InputUpdateDocumentMinerba, userId uint, iupopkId int) (minerba.Minerba, error)
	CreateDmo(dmoInput dmo.CreateDmoInput, userId uint, iupopkId int) (dmo.Dmo, error)
	DeleteDmo(idDmo int, userId uint, iupopkId int) (bool, error)
	UpdateDocumentDmo(id int, documentLink dmo.InputUpdateDocumentDmo, userId uint, iupopkId int) (dmo.Dmo, error)
	UpdateIsDownloadedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, iupopkId int) (dmo.Dmo, error)
	UpdateTrueIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, location string, iupopkId int) (dmo.Dmo, error)
	UpdateFalseIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, iupopkId int) (dmo.Dmo, error)
	UpdateDmo(dmoUpdateInput dmo.UpdateDmoInput, id int, userId uint, iupopkId int) (dmo.Dmo, error)
	CreateProduction(input production.InputCreateProduction, userId uint, iupopkId int) (production.Production, error)
	UpdateProduction(input production.InputCreateProduction, productionId int, userId uint, iupopkId int) (production.Production, error)
	DeleteProduction(productionId int, userId uint, iupopkId int) (bool, error)
	CreateGroupingVesselDN(inputGrouping groupingvesseldn.InputGroupingVesselDn, userId uint, iupopkId int) (groupingvesseldn.GroupingVesselDn, error)
	EditGroupingVesselDn(id int, editGrouping groupingvesseldn.InputEditGroupingVesselDn, userId uint, iupopkId int) (groupingvesseldn.GroupingVesselDn, error)
	DeleteGroupingVesselDn(id int, userId uint, iupopkId int) (bool, error)
	UploadDocumentGroupingVesselDn(id uint, urlS3 string, userId uint, documentType string, iupopkId int) (groupingvesseldn.GroupingVesselDn, error)
	CreateGroupingVesselLN(inputGrouping groupingvesselln.InputGroupingVesselLn, userId uint, iupopkId int) (groupingvesselln.GroupingVesselLn, error)
	EditGroupingVesselLn(id int, editGrouping groupingvesselln.InputEditGroupingVesselLn, userId uint, iupopkId int) (groupingvesselln.GroupingVesselLn, error)
	UploadDocumentGroupingVesselLn(id uint, urlS3 string, userId uint, documentType string, iupopkId int) (groupingvesselln.GroupingVesselLn, error)
	DeleteGroupingVesselLn(id int, userId uint, iupopkId int) (bool, error)
	CreateMinerbaLn(period string, listTransactions []int, userId uint, iupopkId int) (minerbaln.MinerbaLn, error)
	UpdateMinerbaLn(id int, listTransactions []int, userId uint, iupopkId int) (minerbaln.MinerbaLn, error)
	DeleteMinerbaLn(idMinerbaLn int, userId uint, iupopkId int) (bool, error)
	UpdateDocumentMinerbaLn(id int, documentLink minerbaln.InputUpdateDocumentMinerbaLn, userId uint, iupopkId int) (minerbaln.MinerbaLn, error)
	CreateInsw(month string, year int, userId uint, iupopkId int) (insw.Insw, error)
	DeleteInsw(idInsw int, userId uint, iupopkId int) (bool, error)
	UpdateDocumentInsw(id int, documentLink insw.InputUpdateDocumentInsw, userId uint, iupopkId int) (insw.Insw, error)
	CreateReportDmo(input reportdmo.InputCreateReportDmo, userId uint, iupopkId int) (reportdmo.ReportDmo, error)
	UpdateDocumentReportDmo(id int, documentLink reportdmo.InputUpdateDocumentReportDmo, userId uint, iupopkId int) (reportdmo.ReportDmo, error)
	UpdateTransactionReportDmo(id int, inputUpdate reportdmo.InputUpdateReportDmo, userId uint, iupopkId int) (reportdmo.ReportDmo, error)
	DeleteReportDmo(idReportDmo int, userId uint, iupopkId int) (bool, error)
	CreateCoaReport(dateFrom string, dateTo string, iupopkId int, userId uint) (coareport.CoaReport, error)
	DeleteCoaReport(id int, iupopkId int, userId uint) (bool, error)
	UpdateDocumentCoaReport(id int, documentLink coareport.InputUpdateDocumentCoaReport, userId uint, iupopkId int) (coareport.CoaReport, error)
	CreateCoaReportLn(dateFrom string, dateTo string, iupopkId int, userId uint) (coareportln.CoaReportLn, error)
	DeleteCoaReportLn(id int, iupopkId int, userId uint) (bool, error)
	UpdateDocumentCoaReportLn(id int, documentLink coareportln.InputUpdateDocumentCoaReportLn, userId uint, iupopkId int) (coareportln.CoaReportLn, error)
	CreateRkab(input rkab.RkabInput, iupopkId int, userId uint) (rkab.Rkab, error)
	DeleteRkab(id int, iupopkId int, userId uint) (bool, error)
	UploadDocumentRkab(id uint, urlS3 string, userId uint, iupopkId int) (rkab.Rkab, error)
	CreateElectricAssignment(input electricassignmentenduser.CreateElectricAssignmentInput, userId uint, iupopkId int) (electricassignment.ElectricAssignment, error)
	UploadCreateDocumentElectricAssignment(id uint, urlS3 string, userId uint, iupopkId int) (electricassignment.ElectricAssignment, error)
	UploadUpdateDocumentElectricAssignment(id uint, urlS3 string, userId uint, iupopkId int, typeDocument string) (electricassignment.ElectricAssignment, error)
	UpdateElectricAssignment(id int, input electricassignmentenduser.UpdateElectricAssignmentInput, userId uint, iupopkId int) (electricassignment.ElectricAssignment, error)
	DeleteElectricAssignment(id int, iupopkId int, userId uint) (bool, error)
	CreateCafAssignment(input cafassignmentenduser.CreateCafAssignmentInput, userId uint, iupopkId int) (cafassignment.CafAssignment, error)
	UploadCreateDocumentCafAssignment(id uint, urlS3 string, userId uint, iupopkId int) (cafassignment.CafAssignment, error)
	UploadUpdateDocumentCafAssignment(id uint, urlS3 string, userId uint, iupopkId int, typeDocument string) (cafassignment.CafAssignment, error)
	DeleteCafAssignment(id int, iupopkId int, userId uint) (bool, error)
	UpdateCafAssignment(id int, input cafassignmentenduser.UpdateCafAssignmentInput, userId uint, iupopkId int) (cafassignment.CafAssignment, error)
	CreateRoyaltyRecon(dateFrom string, dateTo string, iupopkId int, userId uint) (royaltyrecon.RoyaltyRecon, error)
	DeleteRoyaltyRecon(id int, iupopkId int, userId uint) (bool, error)
	UpdateDocumentRoyaltyRecon(id int, documentLink royaltyrecon.InputUpdateDocumentRoyaltyRecon, userId uint, iupopkId int) (royaltyrecon.RoyaltyRecon, error)
	CreateRoyaltyReport(dateFrom string, dateTo string, iupopkId int, userId uint) (royaltyreport.RoyaltyReport, error)
	DeleteRoyaltyReport(id int, iupopkId int, userId uint) (bool, error)
	UpdateDocumentRoyaltyReport(id int, documentLink royaltyreport.InputUpdateDocumentRoyaltyReport, userId uint, iupopkId int) (royaltyreport.RoyaltyReport, error)
	CreateJettyBalance(input pitloss.InputJettyPitLoss, iupopkId int, userId uint) (jettybalance.JettyBalance, error)
	UpdateJettyBalance(id int, input pitloss.InputUpdateJettyPitLoss, iupopkId int, userId uint) (jettybalance.JettyBalance, error)
	DeleteJettyBalance(id int, userId uint, iupopkId int) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateTransactionDN(inputTransactionDN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error) {
	transaction, transactionErr := s.repository.CreateTransactionDN(inputTransactionDN, userId, iupopkId)

	return transaction, transactionErr
}

func (s *service) DeleteTransaction(id int, userId uint, transactionType string, iupopkId int) (bool, error) {
	isDeletedTransaction, isDeletedTransactionErr := s.repository.DeleteTransaction(id, userId, transactionType, iupopkId)

	return isDeletedTransaction, isDeletedTransactionErr
}

func (s *service) UpdateTransactionDN(idTransaction int, inputEditTransactionDN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error) {
	updateTransaction, updateTransactionErr := s.repository.UpdateTransactionDN(idTransaction, inputEditTransactionDN, userId, iupopkId)

	return updateTransaction, updateTransactionErr
}

func (s *service) UploadDocumentTransaction(idTransaction uint, urlS3 string, userId uint, documentType string, transactionType string, iupopkId int) (transaction.Transaction, error) {
	uploadedDocument, uploadedDocumentErr := s.repository.UploadDocumentTransaction(idTransaction, urlS3, userId, documentType, transactionType, iupopkId)

	return uploadedDocument, uploadedDocumentErr
}

func (s *service) DeleteDocumentTransaction(idTransaction uint, userId uint, documentType string, transactionType string, iupopkId int) (transaction.Transaction, error) {
	transaction, transactionErr := s.repository.DeleteDocumentTransaction(idTransaction, userId, documentType, transactionType, iupopkId)

	return transaction, transactionErr
}

func (s *service) CreateTransactionLN(inputTransactionLN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error) {
	transactionLn, transactionLnErr := s.repository.CreateTransactionLN(inputTransactionLN, userId, iupopkId)

	return transactionLn, transactionLnErr
}

func (s *service) UpdateTransactionLN(id int, inputTransactionLN transaction.DataTransactionInput, userId uint, iupopkId int) (transaction.Transaction, error) {
	transactionLn, transactionLnErr := s.repository.UpdateTransactionLN(id, inputTransactionLN, userId, iupopkId)

	return transactionLn, transactionLnErr
}

func (s *service) CreateMinerba(period string, updateTransaction []int, userId uint, iupopkId int) (minerba.Minerba, error) {
	createdMinerba, createdMinerbaErr := s.repository.CreateMinerba(period, updateTransaction, userId, iupopkId)

	return createdMinerba, createdMinerbaErr
}

func (s *service) UpdateMinerba(id int, updateTransaction []int, userId uint, iupopkId int) (minerba.Minerba, error) {
	updatedMinerba, updatedMinerbaErr := s.repository.UpdateMinerba(id, updateTransaction, userId, iupopkId)

	return updatedMinerba, updatedMinerbaErr
}

func (s *service) DeleteMinerba(idMinerba int, userId uint, iupopkId int) (bool, error) {
	isDeletedMinerba, isDeletedMinerbaErr := s.repository.DeleteMinerba(idMinerba, userId, iupopkId)

	return isDeletedMinerba, isDeletedMinerbaErr
}

func (s *service) UpdateDocumentMinerba(id int, documentLink minerba.InputUpdateDocumentMinerba, userId uint, iupopkId int) (minerba.Minerba, error) {
	uploadMinerba, uploadMinerbaErr := s.repository.UpdateDocumentMinerba(id, documentLink, userId, iupopkId)

	return uploadMinerba, uploadMinerbaErr
}

func (s *service) CreateDmo(dmoInput dmo.CreateDmoInput, userId uint, iupopkId int) (dmo.Dmo, error) {
	createDmo, createDmoErr := s.repository.CreateDmo(dmoInput, userId, iupopkId)

	return createDmo, createDmoErr
}

func (s *service) DeleteDmo(idDmo int, userId uint, iupopkId int) (bool, error) {
	isDeletedDmo, isDeletedDmoErr := s.repository.DeleteDmo(idDmo, userId, iupopkId)

	return isDeletedDmo, isDeletedDmoErr
}

func (s *service) UpdateDocumentDmo(id int, documentLink dmo.InputUpdateDocumentDmo, userId uint, iupopkId int) (dmo.Dmo, error) {
	updateDocumentDmo, updateDocumentDmoErr := s.repository.UpdateDocumentDmo(id, documentLink, userId, iupopkId)

	return updateDocumentDmo, updateDocumentDmoErr
}

func (s *service) UpdateIsDownloadedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, iupopkId int) (dmo.Dmo, error) {
	updateIsDownloadedDmoDocument, updateIsDownloadedDmoDocumentErr := s.repository.UpdateIsDownloadedDmoDocument(isBast, isStatementLetter, isReconciliationLetter, isReconciliationLetterEndUser, id, userId, iupopkId)

	return updateIsDownloadedDmoDocument, updateIsDownloadedDmoDocumentErr
}

func (s *service) UpdateTrueIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, location string, iupopkId int) (dmo.Dmo, error) {
	updateIsSignedDmoDocument, updateIsSignedDmoDocumentErr := s.repository.UpdateTrueIsSignedDmoDocument(isBast, isStatementLetter, isReconciliationLetter, isReconciliationLetterEndUser, id, userId, location, iupopkId)

	return updateIsSignedDmoDocument, updateIsSignedDmoDocumentErr
}

func (s *service) UpdateFalseIsSignedDmoDocument(isBast bool, isStatementLetter bool, isReconciliationLetter bool, isReconciliationLetterEndUser bool, id int, userId uint, iupopkId int) (dmo.Dmo, error) {
	updateIsSignedDmoDocument, updateIsSignedDmoDocumentErr := s.repository.UpdateFalseIsSignedDmoDocument(isBast, isStatementLetter, isReconciliationLetter, isReconciliationLetterEndUser, id, userId, iupopkId)

	return updateIsSignedDmoDocument, updateIsSignedDmoDocumentErr
}

func (s *service) UpdateDmo(dmoUpdateInput dmo.UpdateDmoInput, id int, userId uint, iupopkId int) (dmo.Dmo, error) {
	updateDmo, updateDmoErr := s.repository.UpdateDmo(dmoUpdateInput, id, userId, iupopkId)

	return updateDmo, updateDmoErr
}

func (s *service) CreateProduction(input production.InputCreateProduction, userId uint, iupopkId int) (production.Production, error) {
	createProduction, createProductionErr := s.repository.CreateProduction(input, userId, iupopkId)

	return createProduction, createProductionErr
}

func (s *service) UpdateProduction(input production.InputCreateProduction, productionId int, userId uint, iupopkId int) (production.Production, error) {
	updateProduction, updateProductionErr := s.repository.UpdateProduction(input, productionId, userId, iupopkId)

	return updateProduction, updateProductionErr
}

func (s *service) DeleteProduction(productionId int, userId uint, iupopkId int) (bool, error) {
	deleteProduction, deleteProductionErr := s.repository.DeleteProduction(productionId, userId, iupopkId)

	return deleteProduction, deleteProductionErr
}

func (s *service) CreateGroupingVesselDN(inputGrouping groupingvesseldn.InputGroupingVesselDn, userId uint, iupopkId int) (groupingvesseldn.GroupingVesselDn, error) {
	createGroupingVesselDn, errCreateGroupingVesselDn := s.repository.CreateGroupingVesselDN(inputGrouping, userId, iupopkId)

	return createGroupingVesselDn, errCreateGroupingVesselDn
}

func (s *service) EditGroupingVesselDn(id int, editGrouping groupingvesseldn.InputEditGroupingVesselDn, userId uint, iupopkId int) (groupingvesseldn.GroupingVesselDn, error) {
	editGroupingVesselDn, errEditGroupingVesselDn := s.repository.EditGroupingVesselDn(id, editGrouping, userId, iupopkId)

	return editGroupingVesselDn, errEditGroupingVesselDn
}

func (s *service) DeleteGroupingVesselDn(id int, userId uint, iupopkId int) (bool, error) {
	deleteGroupingVesselDn, errDeleteGroupingVesselDn := s.repository.DeleteGroupingVesselDn(id, userId, iupopkId)

	return deleteGroupingVesselDn, errDeleteGroupingVesselDn
}

func (s *service) UploadDocumentGroupingVesselDn(id uint, urlS3 string, userId uint, documentType string, iupopkId int) (groupingvesseldn.GroupingVesselDn, error) {
	uploadDocumentGroupingVesselDn, errUploadDocumentGroupingVesselDn := s.repository.UploadDocumentGroupingVesselDn(id, urlS3, userId, documentType, iupopkId)

	return uploadDocumentGroupingVesselDn, errUploadDocumentGroupingVesselDn
}

func (s *service) CreateGroupingVesselLN(inputGrouping groupingvesselln.InputGroupingVesselLn, userId uint, iupopkId int) (groupingvesselln.GroupingVesselLn, error) {
	createGroupingVesselLn, errCreateGroupingVesselLn := s.repository.CreateGroupingVesselLN(inputGrouping, userId, iupopkId)

	return createGroupingVesselLn, errCreateGroupingVesselLn
}

func (s *service) EditGroupingVesselLn(id int, editGrouping groupingvesselln.InputEditGroupingVesselLn, userId uint, iupopkId int) (groupingvesselln.GroupingVesselLn, error) {
	editGroupingVesselLn, errEditGroupingVesselLn := s.repository.EditGroupingVesselLn(id, editGrouping, userId, iupopkId)

	return editGroupingVesselLn, errEditGroupingVesselLn
}

func (s *service) UploadDocumentGroupingVesselLn(id uint, urlS3 string, userId uint, documentType string, iupopkId int) (groupingvesselln.GroupingVesselLn, error) {
	uploadDocumentGroupingVesselLn, uploadDocumentGroupingVesselLnErr := s.repository.UploadDocumentGroupingVesselLn(id, urlS3, userId, documentType, iupopkId)

	return uploadDocumentGroupingVesselLn, uploadDocumentGroupingVesselLnErr
}

func (s *service) DeleteGroupingVesselLn(id int, userId uint, iupopkId int) (bool, error) {
	deleteGroupingVesselLn, deleteGroupingVesselLnErr := s.repository.DeleteGroupingVesselLn(id, userId, iupopkId)

	return deleteGroupingVesselLn, deleteGroupingVesselLnErr
}

func (s *service) CreateMinerbaLn(period string, listTransactions []int, userId uint, iupopkId int) (minerbaln.MinerbaLn, error) {
	createMinerbaLn, createMinerbaLnErr := s.repository.CreateMinerbaLn(period, listTransactions, userId, iupopkId)

	return createMinerbaLn, createMinerbaLnErr
}

func (s *service) UpdateMinerbaLn(id int, listTransactions []int, userId uint, iupopkId int) (minerbaln.MinerbaLn, error) {
	updateMinerbaLn, updateMinerbaLnErr := s.repository.UpdateMinerbaLn(id, listTransactions, userId, iupopkId)

	return updateMinerbaLn, updateMinerbaLnErr
}

func (s *service) DeleteMinerbaLn(idMinerbaLn int, userId uint, iupopkId int) (bool, error) {
	isDeletedMinerbaLn, isDeletedMinerbaLnErr := s.repository.DeleteMinerbaLn(idMinerbaLn, userId, iupopkId)

	return isDeletedMinerbaLn, isDeletedMinerbaLnErr
}

func (s *service) UpdateDocumentMinerbaLn(id int, documentLink minerbaln.InputUpdateDocumentMinerbaLn, userId uint, iupopkId int) (minerbaln.MinerbaLn, error) {
	uploadMinerbaLn, uploadMinerbaLnErr := s.repository.UpdateDocumentMinerbaLn(id, documentLink, userId, iupopkId)

	return uploadMinerbaLn, uploadMinerbaLnErr
}

func (s *service) CreateInsw(month string, year int, userId uint, iupopkId int) (insw.Insw, error) {
	createInsw, createInswErr := s.repository.CreateInsw(month, year, userId, iupopkId)

	return createInsw, createInswErr
}

func (s *service) DeleteInsw(idInsw int, userId uint, iupopkId int) (bool, error) {
	deleteInsw, deleteInswErr := s.repository.DeleteInsw(idInsw, userId, iupopkId)

	return deleteInsw, deleteInswErr
}

func (s *service) UpdateDocumentInsw(id int, documentLink insw.InputUpdateDocumentInsw, userId uint, iupopkId int) (insw.Insw, error) {
	updateDocumentInsw, updateDocumentInswErr := s.repository.UpdateDocumentInsw(id, documentLink, userId, iupopkId)

	return updateDocumentInsw, updateDocumentInswErr
}

func (s *service) CreateReportDmo(input reportdmo.InputCreateReportDmo, userId uint, iupopkId int) (reportdmo.ReportDmo, error) {
	createReportDmo, createReportDmoErr := s.repository.CreateReportDmo(input, userId, iupopkId)

	return createReportDmo, createReportDmoErr
}

func (s *service) UpdateDocumentReportDmo(id int, documentLink reportdmo.InputUpdateDocumentReportDmo, userId uint, iupopkId int) (reportdmo.ReportDmo, error) {
	updateDocumentReportDmo, updateDocumentReportDmoErr := s.repository.UpdateDocumentReportDmo(id, documentLink, userId, iupopkId)

	return updateDocumentReportDmo, updateDocumentReportDmoErr
}

func (s *service) UpdateTransactionReportDmo(id int, inputUpdate reportdmo.InputUpdateReportDmo, userId uint, iupopkId int) (reportdmo.ReportDmo, error) {
	updateTransactionReportDmo, updateTransactionReportDmoErr := s.repository.UpdateTransactionReportDmo(id, inputUpdate, userId, iupopkId)

	return updateTransactionReportDmo, updateTransactionReportDmoErr
}

func (s *service) DeleteReportDmo(idReportDmo int, userId uint, iupopkId int) (bool, error) {
	deleteReportDmo, deleteReportDmoErr := s.repository.DeleteReportDmo(idReportDmo, userId, iupopkId)

	return deleteReportDmo, deleteReportDmoErr
}

func (s *service) CreateCoaReport(dateFrom string, dateTo string, iupopkId int, userId uint) (coareport.CoaReport, error) {
	coaReport, coaReportErr := s.repository.CreateCoaReport(dateFrom, dateTo, iupopkId, userId)

	return coaReport, coaReportErr
}

func (s *service) DeleteCoaReport(id int, iupopkId int, userId uint) (bool, error) {
	isDeletedCoaReport, isDeletedCoaReportErr := s.repository.DeleteCoaReport(id, iupopkId, userId)

	return isDeletedCoaReport, isDeletedCoaReportErr
}

func (s *service) UpdateDocumentCoaReport(id int, documentLink coareport.InputUpdateDocumentCoaReport, userId uint, iupopkId int) (coareport.CoaReport, error) {
	updDocumentCoaReport, updDocumentCoaReportErr := s.repository.UpdateDocumentCoaReport(id, documentLink, userId, iupopkId)

	return updDocumentCoaReport, updDocumentCoaReportErr
}

func (s *service) CreateCoaReportLn(dateFrom string, dateTo string, iupopkId int, userId uint) (coareportln.CoaReportLn, error) {
	coaReportLn, coaReportLnErr := s.repository.CreateCoaReportLn(dateFrom, dateTo, iupopkId, userId)

	return coaReportLn, coaReportLnErr
}
func (s *service) DeleteCoaReportLn(id int, iupopkId int, userId uint) (bool, error) {
	isDeletedCoaReportLn, isDeletedCoaReportLnErr := s.repository.DeleteCoaReportLn(id, iupopkId, userId)

	return isDeletedCoaReportLn, isDeletedCoaReportLnErr
}
func (s *service) UpdateDocumentCoaReportLn(id int, documentLink coareportln.InputUpdateDocumentCoaReportLn, userId uint, iupopkId int) (coareportln.CoaReportLn, error) {
	updDocumentCoaReportLn, updDocumentCoaReportLnErr := s.repository.UpdateDocumentCoaReportLn(id, documentLink, userId, iupopkId)

	return updDocumentCoaReportLn, updDocumentCoaReportLnErr
}

func (s *service) CreateRkab(input rkab.RkabInput, iupopkId int, userId uint) (rkab.Rkab, error) {
	createdRkab, createdRkabErr := s.repository.CreateRkab(input, iupopkId, userId)

	return createdRkab, createdRkabErr
}

func (s *service) DeleteRkab(id int, iupopkId int, userId uint) (bool, error) {
	isDeletedRkab, isDeletedRkabErr := s.repository.DeleteRkab(id, iupopkId, userId)

	return isDeletedRkab, isDeletedRkabErr
}

func (s *service) UploadDocumentRkab(id uint, urlS3 string, userId uint, iupopkId int) (rkab.Rkab, error) {
	updateRkab, updateRkabErr := s.repository.UploadDocumentRkab(id, urlS3, userId, iupopkId)

	return updateRkab, updateRkabErr
}

func (s *service) CreateElectricAssignment(input electricassignmentenduser.CreateElectricAssignmentInput, userId uint, iupopkId int) (electricassignment.ElectricAssignment, error) {
	createElectricAssignment, createElectricAssignmentErr := s.repository.CreateElectricAssignment(input, userId, iupopkId)

	return createElectricAssignment, createElectricAssignmentErr
}

func (s *service) UploadCreateDocumentElectricAssignment(id uint, urlS3 string, userId uint, iupopkId int) (electricassignment.ElectricAssignment, error) {
	uploadElectricAssignment, uploadElectricAssignmentErr := s.repository.UploadCreateDocumentElectricAssignment(id, urlS3, userId, iupopkId)

	return uploadElectricAssignment, uploadElectricAssignmentErr
}

func (s *service) UploadUpdateDocumentElectricAssignment(id uint, urlS3 string, userId uint, iupopkId int, typeDocument string) (electricassignment.ElectricAssignment, error) {
	uploadElectricAssignment, uploadElectricAssignmentErr := s.repository.UploadUpdateDocumentElectricAssignment(id, urlS3, userId, iupopkId, typeDocument)

	return uploadElectricAssignment, uploadElectricAssignmentErr
}

func (s *service) UpdateElectricAssignment(id int, input electricassignmentenduser.UpdateElectricAssignmentInput, userId uint, iupopkId int) (electricassignment.ElectricAssignment, error) {
	updateElectricAssigment, updateElectricAssigmentErr := s.repository.UpdateElectricAssignment(id, input, userId, iupopkId)

	return updateElectricAssigment, updateElectricAssigmentErr
}

func (s *service) DeleteElectricAssignment(id int, iupopkId int, userId uint) (bool, error) {
	isDeleted, isDeletedErr := s.repository.DeleteElectricAssignment(id, iupopkId, userId)

	return isDeleted, isDeletedErr
}

func (s *service) CreateCafAssignment(input cafassignmentenduser.CreateCafAssignmentInput, userId uint, iupopkId int) (cafassignment.CafAssignment, error) {
	cafAssignment, cafAssignmentErr := s.repository.CreateCafAssignment(input, userId, iupopkId)

	return cafAssignment, cafAssignmentErr
}

func (s *service) UploadCreateDocumentCafAssignment(id uint, urlS3 string, userId uint, iupopkId int) (cafassignment.CafAssignment, error) {
	cafAssignment, cafAssignmentErr := s.repository.UploadCreateDocumentCafAssignment(id, urlS3, userId, iupopkId)

	return cafAssignment, cafAssignmentErr
}

func (s *service) UploadUpdateDocumentCafAssignment(id uint, urlS3 string, userId uint, iupopkId int, typeDocument string) (cafassignment.CafAssignment, error) {
	cafAssignment, cafAssignmentErr := s.repository.UploadUpdateDocumentCafAssignment(id, urlS3, userId, iupopkId, typeDocument)

	return cafAssignment, cafAssignmentErr
}

func (s *service) DeleteCafAssignment(id int, iupopkId int, userId uint) (bool, error) {
	isDeletedCafAssignment, isDeletedCafAssignmentErr := s.repository.DeleteCafAssignment(id, iupopkId, userId)

	return isDeletedCafAssignment, isDeletedCafAssignmentErr
}

func (s *service) UpdateCafAssignment(id int, input cafassignmentenduser.UpdateCafAssignmentInput, userId uint, iupopkId int) (cafassignment.CafAssignment, error) {
	cafAssignment, cafAssignmentErr := s.repository.UpdateCafAssignment(id, input, userId, iupopkId)

	return cafAssignment, cafAssignmentErr
}

func (s *service) CreateRoyaltyRecon(dateFrom string, dateTo string, iupopkId int, userId uint) (royaltyrecon.RoyaltyRecon, error) {
	royaltyRecon, royaltyReconErr := s.repository.CreateRoyaltyRecon(dateFrom, dateTo, iupopkId, userId)

	return royaltyRecon, royaltyReconErr
}

func (s *service) DeleteRoyaltyRecon(id int, iupopkId int, userId uint) (bool, error) {
	royaltyRecon, royaltyReconErr := s.repository.DeleteRoyaltyRecon(id, iupopkId, userId)

	return royaltyRecon, royaltyReconErr
}

func (s *service) UpdateDocumentRoyaltyRecon(id int, documentLink royaltyrecon.InputUpdateDocumentRoyaltyRecon, userId uint, iupopkId int) (royaltyrecon.RoyaltyRecon, error) {
	royaltyRecon, royaltyReconErr := s.repository.UpdateDocumentRoyaltyRecon(id, documentLink, userId, iupopkId)

	return royaltyRecon, royaltyReconErr
}

func (s *service) CreateRoyaltyReport(dateFrom string, dateTo string, iupopkId int, userId uint) (royaltyreport.RoyaltyReport, error) {
	royaltyReport, royaltyReportErr := s.repository.CreateRoyaltyReport(dateFrom, dateTo, iupopkId, userId)

	return royaltyReport, royaltyReportErr
}

func (s *service) DeleteRoyaltyReport(id int, iupopkId int, userId uint) (bool, error) {
	royaltyReport, royaltyReportErr := s.repository.DeleteRoyaltyReport(id, iupopkId, userId)

	return royaltyReport, royaltyReportErr
}

func (s *service) UpdateDocumentRoyaltyReport(id int, documentLink royaltyreport.InputUpdateDocumentRoyaltyReport, userId uint, iupopkId int) (royaltyreport.RoyaltyReport, error) {
	royaltyReport, royaltyReportErr := s.repository.UpdateDocumentRoyaltyReport(id, documentLink, userId, iupopkId)

	return royaltyReport, royaltyReportErr
}

func (s *service) CreateJettyBalance(input pitloss.InputJettyPitLoss, iupopkId int, userId uint) (jettybalance.JettyBalance, error) {
	data, dataErr := s.repository.CreateJettyBalance(input, iupopkId, userId)

	return data, dataErr
}

func (s *service) UpdateJettyBalance(id int, input pitloss.InputUpdateJettyPitLoss, iupopkId int, userId uint) (jettybalance.JettyBalance, error) {
	data, dataErr := s.repository.UpdateJettyBalance(id, input, iupopkId, userId)

	return data, dataErr
}

func (s *service) DeleteJettyBalance(id int, userId uint, iupopkId int) (bool, error) {
	data, dataErr := s.repository.DeleteJettyBalance(id, userId, iupopkId)

	return data, dataErr
}
