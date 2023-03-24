package transaction

import (
	"ajebackend/model/dmo"
	"ajebackend/model/dmovessel"
	"ajebackend/model/groupingvesseldn"
	"ajebackend/model/groupingvesselln"
	"ajebackend/model/master/salessystem"
	"ajebackend/model/master/trader"
	"ajebackend/model/minerba"
	"ajebackend/model/minerbaln"
	"ajebackend/model/production"
	"ajebackend/model/reportdmo"
	"ajebackend/model/traderdmo"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func helperString(listString []string, dataString string) bool {
	for _, v := range listString {
		if v == dataString {
			return true
		}
	}
	return false
}

type Repository interface {
	ListData(page int, sortFilter SortAndFilter, transactionType string, iupopkId int) (Pagination, error)
	DetailTransaction(id int, transactionType string, iupopkId int) (Transaction, error)
	ListDataDNWithoutMinerba(iupopkId int) ([]Transaction, error)
	CheckDataDnAndMinerba(listData []int, iupopkId int) (bool, error)
	CheckDataDnAndMinerbaUpdate(listData []int, idMinerba int, iupopkId int) ([]Transaction, error)
	GetDetailMinerba(id int, iupopkId int) (DetailMinerba, error)
	ListDataDNBargeWithoutVessel(iupopkId int) ([]Transaction, error)
	ListDataDNBargeWithVessel(iupopkId int) ([]Transaction, error)
	ListDataDNVessel(iupopkId int) ([]Transaction, error)
	CheckDataDnAndDmo(listData []int, iupopkId int) ([]Transaction, error)
	CheckGroupingVesselAndDmo(listData []int, iupopkId int) ([]dmovessel.DmoVessel, error)
	GetDetailDmo(id int, iupopkId int) (DetailDmo, error)
	GetDataDmo(id uint, iupopkId int) (ListTransactionDmoBackgroundJob, error)
	GetDetailReportDmo(id int, iupopkId int) (DetailReportDmo, error)
	CheckDataUnique(inputTrans DataTransactionInput) (bool, bool, bool, bool)
	GetReport(year int, iupopkId int) (ReportRecapOutput, ReportDetailOutput, error)
	GetListForReport(iupopkId int) (ListForCreatingReportDmoOutput, error)
	GetDetailGroupingVesselDn(id int, iupopkId int) (DetailGroupingVesselDn, error)
	ListDataDnWithoutGroup(iupopkId int) (ListTransactionNotHaveGroupingVessel, error)
	GetDetailGroupingVesselLn(id int, iupopkId int) (DetailGroupingVesselLn, error)
	ListDataLnWithoutGroup(iupopkId int) ([]Transaction, error)
	ListDataLNWithoutMinerba(iupopkId int) ([]Transaction, error)
	GetDetailMinerbaLn(id int, iupopkId int) (DetailMinerbaLn, error)
	CheckDataLnAndMinerbaLnUpdate(listData []int, idMinerba int, iupopkId int) ([]Transaction, error)
	CheckDataLnAndMinerbaLn(listData []int, iupopkId int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// Transaction

func (r *repository) ListData(page int, sortFilter SortAndFilter, transactionType string, iupopkId int) (Pagination, error) {
	var transactions []Transaction
	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	defaultSort := "id desc"
	sortString := fmt.Sprintf("%s %s", sortFilter.Field, sortFilter.Sort)
	if sortFilter.Field == "" || sortFilter.Sort == "" {
		sortString = defaultSort
	}

	if sortFilter.Field != "" || sortFilter.Sort != "" {
		if sortFilter.Field == "tugboat_id" {
			sortString = "tugboats.name " + sortFilter.Sort
		}

		if sortFilter.Field == "barge_id" {
			sortString = "barges.name " + sortFilter.Sort
		}

		if sortFilter.Field == "vessel_id" {
			sortString = "vessels.name " + sortFilter.Sort
		}

		if sortFilter.Field == "customer_id" {
			sortString = "companies.company_name " + sortFilter.Sort
		}

		if sortFilter.Field == "shipping_date" {
			sortString = "shipping_date " + sortFilter.Sort
		}

		if sortFilter.Field == "quantity" {
			sortString = "quantity_unloading " + sortFilter.Sort
		}
	}

	queryFilter := fmt.Sprintf("transaction_type = '%s' AND seller_id = %v", strings.ToUpper(transactionType), iupopkId)

	if sortFilter.TugboatId != "" {
		queryFilter = queryFilter + " AND tugboat_id = " + sortFilter.TugboatId
	}

	if sortFilter.BargeId != "" {
		queryFilter = queryFilter + " AND barge_id = " + sortFilter.BargeId
	}

	if sortFilter.VesselId != "" {
		queryFilter = queryFilter + " AND vessel_id = " + sortFilter.VesselId
	}

	if sortFilter.ShippingStart != "" {
		queryFilter = queryFilter + " AND shipping_date >= '" + sortFilter.ShippingStart + "'"
	}

	if sortFilter.ShippingEnd != "" {
		queryFilter = queryFilter + " AND shipping_date <= '" + sortFilter.ShippingEnd + "T23:59:59'"
	}

	if sortFilter.Quantity != "" {
		queryFilter = queryFilter + " AND cast(quantity_unloading AS TEXT) LIKE '%" + sortFilter.Quantity + "%'"
	}

	if sortFilter.VerificationFilter == "Belum diverifikasi" {
		queryFilter = queryFilter + " AND is_finance_check = FALSE"
	}

	if sortFilter.VerificationFilter == "Sudah diverifikasi" {
		queryFilter = queryFilter + " AND is_finance_check = TRUE AND is_coa_finish IS NOT TRUE AND is_royalty_final_finish IS NOT TRUE"
	}

	if sortFilter.VerificationFilter == "Data belum lengkap" {
		queryFilter = queryFilter + " AND is_finance_check = TRUE AND ((is_coa_finish = TRUE AND is_royalty_final_finish = FALSE) OR (is_coa_finish IS NOT TRUE AND is_royalty_final_finish = TRUE))"
	}

	if sortFilter.VerificationFilter == "Data lengkap" {
		queryFilter = queryFilter + " AND is_finance_check = TRUE AND is_coa_finish = TRUE AND is_royalty_final_finish = TRUE"
	}

	errFind := r.db.Table("transactions").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Select("transactions.*").Joins("LEFT JOIN tugboats tugboats on transactions.tugboat_id = tugboats.id").Joins("LEFT JOIN barges barges on transactions.barge_id = barges.id").Joins("LEFT JOIN vessels vessels on transactions.vessel_id = vessels.id").Joins("LEFT JOIN companies companies on transactions.customer_id = companies.id").Order(sortString).Where(queryFilter).Scopes(paginateData(transactions, &pagination, r.db, queryFilter)).Find(&transactions).Error

	if errFind != nil {

		errWithoutOrder := r.db.Table("transactions").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Order(defaultSort).Where(queryFilter).Scopes(paginateData(transactions, &pagination, r.db, queryFilter)).Find(&transactions).Error

		if errWithoutOrder != nil {
			pagination.Data = transactions
			return pagination, errWithoutOrder
		}
	}

	pagination.Data = transactions

	return pagination, nil
}

func (r *repository) DetailTransaction(id int, transactionType string, iupopkId int) (Transaction, error) {
	var transaction Transaction

	errFind := r.db.Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("id = ? AND transaction_type = ? AND seller_id = ?", id, transactionType, iupopkId).First(&transaction).Error

	return transaction, errFind
}

func (r *repository) ListDataDNWithoutMinerba(iupopkId int) ([]Transaction, error) {
	var listDataDnWithoutMinerba []Transaction

	errFind := r.db.Order("id desc").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("minerba_id is NULL AND transaction_type = ? AND is_not_claim = ? AND is_migration = ? AND is_finance_check = ? AND seller_id = ?", "DN", false, false, true, iupopkId).Find(&listDataDnWithoutMinerba).Error

	return listDataDnWithoutMinerba, errFind
}

func (r *repository) CheckDataUnique(inputTrans DataTransactionInput) (bool, bool, bool, bool) {
	isDpRoyaltyNtpnUnique := false
	isDpRoyaltyBillingCodeUnique := false
	isPaymentDpRoyaltyNtpnUnique := false
	isPaymentDpRoyaltyBillingCodeUnique := false

	if inputTrans.DpRoyaltyNtpn != nil {
		var checkDpRoyaltyNtpn Transaction

		errCheckDpRoyaltyNtpn := r.db.Where("dp_royalty_ntpn = ?", inputTrans.DpRoyaltyNtpn).First(&checkDpRoyaltyNtpn).Error

		if errCheckDpRoyaltyNtpn == nil {
			isDpRoyaltyNtpnUnique = true
		}
	}

	if inputTrans.DpRoyaltyBillingCode != nil {
		var checkDpRoyaltyBillingCode Transaction

		errCheckDpRoyaltyBillingCode := r.db.Where("dp_royalty_billing_code", inputTrans.DpRoyaltyBillingCode).First(&checkDpRoyaltyBillingCode).Error

		if errCheckDpRoyaltyBillingCode == nil {
			isDpRoyaltyBillingCodeUnique = true
		}
	}

	if inputTrans.PaymentDpRoyaltyNtpn != nil {
		var checkPaymentDpRoyaltyNtpn Transaction

		errCheckPaymentDpRoyaltyNtpn := r.db.Where("payment_dp_royalty_ntpn", inputTrans.PaymentDpRoyaltyNtpn).First(&checkPaymentDpRoyaltyNtpn).Error

		if errCheckPaymentDpRoyaltyNtpn == nil {
			isPaymentDpRoyaltyNtpnUnique = true
		}
	}

	if inputTrans.PaymentDpRoyaltyBillingCode != nil {
		var checkPaymentDpRoyaltyBillingCode Transaction

		errCheckPaymentDpRoyaltyBillingCode := r.db.Where("payment_dp_royalty_billing_code", inputTrans.PaymentDpRoyaltyBillingCode).First(&checkPaymentDpRoyaltyBillingCode).Error

		if errCheckPaymentDpRoyaltyBillingCode == nil {
			isPaymentDpRoyaltyBillingCodeUnique = true
		}
	}

	return isDpRoyaltyNtpnUnique, isDpRoyaltyBillingCodeUnique, isPaymentDpRoyaltyNtpnUnique, isPaymentDpRoyaltyBillingCodeUnique
}

// Minerba

func (r *repository) CheckDataDnAndMinerba(listData []int, iupopkId int) (bool, error) {
	var listDnValid []Transaction

	errFindValid := r.db.Where("id IN ? AND seller_id = ?", listData, iupopkId).Find(&listDnValid).Error

	if errFindValid != nil {
		return false, errFindValid
	}

	if len(listData) != len(listDnValid) {
		return false, errors.New("please check there is transaction not found")
	}

	var listDn []Transaction

	errFind := r.db.Where("minerba_id is NULL AND id IN ?", listData).Find(&listDn).Error

	if errFind != nil {
		return false, errFind
	}

	if len(listDn) != len(listData) {
		return false, errors.New("please check there is transaction already in report")
	}

	return true, nil
}

func (r *repository) CheckDataDnAndMinerbaUpdate(listData []int, idMinerba int, iupopkId int) ([]Transaction, error) {
	var listDnValid []Transaction

	errFindValid := r.db.Where("id IN ? AND seller_id = ?", listData, iupopkId).Find(&listDnValid).Error

	if errFindValid != nil {
		return listDnValid, errFindValid
	}

	if len(listData) != len(listDnValid) {
		return listDnValid, errors.New("please check there is transaction not found")
	}

	var listDn []Transaction

	errFind := r.db.Where("id IN ? AND seller_id = ?", listData, iupopkId).Find(&listDn).Error

	if errFind != nil {
		return listDn, errFind
	}

	uintIdMinerba := uint(idMinerba)

	for _, v := range listDn {
		if v.MinerbaId != nil && *v.MinerbaId != uintIdMinerba {
			return listDn, errors.New("please check there is transaction already in report")
		}
	}

	return listDn, nil
}

func (r *repository) GetDetailMinerba(id int, iupopkId int) (DetailMinerba, error) {

	var detailMinerba DetailMinerba

	var minerba minerba.Minerba
	var transactions []Transaction

	minerbaFindErr := r.db.Preload(clause.Associations).Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&minerba).Error

	if minerbaFindErr != nil {
		return detailMinerba, minerbaFindErr
	}

	detailMinerba.Detail = minerba

	transactionFindErr := r.db.Order("shipping_date asc").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("minerba_id = ? AND seller_id = ?", id, iupopkId).Find(&transactions).Error

	if transactionFindErr != nil {
		return detailMinerba, transactionFindErr
	}

	detailMinerba.List = transactions
	return detailMinerba, nil
}

// DMO

func (r *repository) ListDataDNBargeWithoutVessel(iupopkId int) ([]Transaction, error) {
	var listDataDnBargeDmo []Transaction

	var salesSystem []salessystem.SalesSystem
	var salesSystemId []uint

	errFindSalesSystem := r.db.Where("name ILIKE '%Barge'").Find(&salesSystem).Error

	if errFindSalesSystem != nil {
		return listDataDnBargeDmo, errFindSalesSystem
	}

	for _, v := range salesSystem {
		salesSystemId = append(salesSystemId, v.ID)
	}

	errFindBarge := r.db.Order("id desc").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("dmo_id is NULL AND transaction_type = ? AND is_not_claim = ? AND is_migration = ? AND vessel_id IS NULL AND sales_system_id IN ? AND grouping_vessel_dn_id is NULL AND seller_id = ?", "DN", false, false, salesSystemId, iupopkId).Find(&listDataDnBargeDmo).Error

	if errFindBarge != nil {
		return listDataDnBargeDmo, errFindBarge
	}

	return listDataDnBargeDmo, nil
}

func (r *repository) ListDataDNBargeWithVessel(iupopkId int) ([]Transaction, error) {
	var listDataDnBargeDmo []Transaction

	var salesSystem []salessystem.SalesSystem
	var salesSystemId []uint

	errFindSalesSystem := r.db.Where("name ILIKE '%Barge'").Find(&salesSystem).Error

	if errFindSalesSystem != nil {
		return listDataDnBargeDmo, errFindSalesSystem
	}

	for _, v := range salesSystem {
		salesSystemId = append(salesSystemId, v.ID)
	}

	errFindBarge := r.db.Order("id desc").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("dmo_id is NULL AND transaction_type = ? AND is_not_claim = ? AND is_migration = ? AND vessel_id IS NOT NULL AND sales_system_id IN ? AND grouping_vessel_dn_id is NULL AND seller_id = ?", "DN", false, false, salesSystemId, iupopkId).Find(&listDataDnBargeDmo).Error

	if errFindBarge != nil {
		return listDataDnBargeDmo, errFindBarge
	}

	return listDataDnBargeDmo, nil
}

func (r *repository) ListDataDNVessel(iupopkId int) ([]Transaction, error) {
	var listDataDnVessel []Transaction

	var salesSystem []salessystem.SalesSystem
	var salesSystemId []uint

	errFindSalesSystem := r.db.Where("name ILIKE '%Vessel'").Find(&salesSystem).Error

	if errFindSalesSystem != nil {
		return listDataDnVessel, errFindSalesSystem
	}

	for _, v := range salesSystem {
		salesSystemId = append(salesSystemId, v.ID)
	}

	errFindBarge := r.db.Order("id desc").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("dmo_id is NULL AND transaction_type = ? AND is_not_claim = ? AND is_migration = ? AND vessel_id IS NOT NULL AND is_finance_check = ? AND sales_system IN ? AND grouping_vessel_dn_id is NULL AND seller_id = ?", "DN", false, false, true, salesSystemId, iupopkId).Find(&listDataDnVessel).Error

	if errFindBarge != nil {
		return listDataDnVessel, errFindBarge
	}

	return listDataDnVessel, nil
}

func (r *repository) CheckDataDnAndDmo(listData []int, iupopkId int) ([]Transaction, error) {
	var listDnValid []Transaction

	errFindValid := r.db.Where("id IN ?", listData).Find(&listDnValid).Error

	if errFindValid != nil {
		return listDnValid, errFindValid
	}

	if len(listData) != len(listDnValid) {
		return listDnValid, errors.New("please check there is transaction not found")
	}

	var listDn []Transaction

	errFind := r.db.Where("dmo_id is NULL AND id IN ?", listData).Find(&listDn).Error

	if errFind != nil {
		return listDn, errFind
	}

	if len(listDn) != len(listData) {
		return listDn, errors.New("please check there is transaction already in report")
	}

	return listDn, nil
}

func (r *repository) CheckGroupingVesselAndDmo(listData []int, iupopkId int) ([]dmovessel.DmoVessel, error) {
	var listGroupingVessel []dmovessel.DmoVessel

	errFind := r.db.Where("grouping_vessel_dn_id in ?", listData).Find(&listGroupingVessel).Error

	if len(listGroupingVessel) > 0 {
		return listGroupingVessel, errors.New("please check grouping vessel already in report")
	}

	return listGroupingVessel, errFind
}

func (r *repository) GetDetailDmo(id int, iupopkId int) (DetailDmo, error) {

	var detailDmo DetailDmo

	var dmoData dmo.Dmo
	var transactions []Transaction
	var groupingVessel []groupingvesseldn.GroupingVesselDn

	dmoFindErr := r.db.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&dmoData).Error

	if dmoFindErr != nil {
		return detailDmo, dmoFindErr
	}

	detailDmo.Detail = dmoData

	transactionFindErr := r.db.Order("id desc").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("dmo_id = ? and grouping_vessel_dn_id IS NULL AND seller_id = ?", id, iupopkId).Find(&transactions).Error

	if transactionFindErr != nil {
		return detailDmo, transactionFindErr
	}

	detailDmo.Transactions = transactions

	var dmoVessel []dmovessel.DmoVessel

	dmoVesselFindErr := r.db.Preload(clause.Associations).Preload("GroupingVesselDn.Buyer").Preload("GroupingVesselDn.Destination").Preload("GroupingVesselDn.Vessel").Preload("GroupingVesselDn.DmoDestinationPort.PortLocation").Where("dmo_id = ?", id).Find(&dmoVessel).Error

	if dmoVesselFindErr != nil {
		return detailDmo, dmoVesselFindErr
	}

	for _, v := range dmoVessel {
		groupingVessel = append(groupingVessel, v.GroupingVesselDn)
	}

	detailDmo.GroupingVessels = groupingVessel

	var traderData []trader.Trader
	var endUser trader.Trader
	var traderDmo []traderdmo.TraderDmo
	traderDmoFindErr := r.db.Order(`"order" asc`).Preload(clause.Associations).Preload("Trader.Company.IndustryType").Where("dmo_id = ?", id).Find(&traderDmo).Error

	if traderDmoFindErr != nil {
		return detailDmo, traderDmoFindErr
	}

	for _, v := range traderDmo {
		if v.IsEndUser {
			endUser = v.Trader
		} else {
			traderData = append(traderData, v.Trader)
		}
	}

	detailDmo.Trader = traderData

	detailDmo.EndUser = endUser

	return detailDmo, nil
}

func (r *repository) GetDataDmo(id uint, iupopkId int) (ListTransactionDmoBackgroundJob, error) {
	var listTransactionDmoBackgroundJob ListTransactionDmoBackgroundJob

	var transactionBarge []Transaction
	var transactionGroupingVessel []Transaction
	var groupingVessel []groupingvesseldn.GroupingVesselDn

	errFindTransactionBarge := r.db.Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("dmo_id = ? AND grouping_vessel_dn_id is NULL", id).Order("shipping_date asc").Find(&transactionBarge).Error

	if errFindTransactionBarge != nil {
		return listTransactionDmoBackgroundJob, errFindTransactionBarge
	}

	listTransactionDmoBackgroundJob.ListTransactionBarge = transactionBarge

	errFindTransactionGroupingVessel := r.db.Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("dmo_id = ? AND grouping_vessel_dn_id IS NOT NULL", id).Order("shipping_date asc").Find(&transactionGroupingVessel).Error

	if errFindTransactionGroupingVessel != nil {
		return listTransactionDmoBackgroundJob, errFindTransactionGroupingVessel
	}

	listTransactionDmoBackgroundJob.ListTransactionGroupingVessel = transactionGroupingVessel

	var dmoVessel []dmovessel.DmoVessel

	errFindDmoVessel := r.db.Preload(clause.Associations).Where("dmo_id = ?", id).Find(&dmoVessel).Error

	if errFindDmoVessel != nil {
		return listTransactionDmoBackgroundJob, errFindDmoVessel
	}

	var groupingVesselId []uint

	for _, v := range dmoVessel {
		groupingVesselId = append(groupingVesselId, v.GroupingVesselDnId)
	}

	errFindGroupingVessel := r.db.Preload(clause.Associations).Where("id in ?", groupingVesselId).Order("bl_date asc").Find(&groupingVessel).Error

	if errFindGroupingVessel != nil {
		return listTransactionDmoBackgroundJob, errFindGroupingVessel
	}

	listTransactionDmoBackgroundJob.ListGroupingVessel = groupingVessel

	return listTransactionDmoBackgroundJob, nil
}

func (r *repository) GetDetailReportDmo(id int, iupopkId int) (DetailReportDmo, error) {
	var detailReportDmo DetailReportDmo

	var reportDmoData reportdmo.ReportDmo
	var transactions []Transaction
	var groupingVesselDn []groupingvesseldn.GroupingVesselDn

	var salesSystem []salessystem.SalesSystem
	var salesSystemId []uint

	errFindSalesSystem := r.db.Where("name ILIKE '%Barge'").Find(&salesSystem).Error

	if errFindSalesSystem != nil {
		return detailReportDmo, errFindSalesSystem
	}

	for _, v := range salesSystem {
		salesSystemId = append(salesSystemId, v.ID)
	}

	errFindTransactions := r.db.Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Order("shipping_date asc").Where("report_dmo_id = ? AND sales_system_id IN ? AND seller_id = ?", id, salesSystemId, iupopkId).Find(&transactions).Error

	if errFindTransactions != nil {
		return detailReportDmo, errFindTransactions
	}

	errFindGroupingVessel := r.db.Preload(clause.Associations).Order("bl_date asc").Where("report_dmo_id = ? AND iupopk_id = ?", id, iupopkId).Find(&groupingVesselDn).Error

	if errFindGroupingVessel != nil {
		return detailReportDmo, errFindGroupingVessel
	}

	errFindReportDmo := r.db.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&reportDmoData).Error

	if errFindReportDmo != nil {
		return detailReportDmo, errFindReportDmo
	}

	detailReportDmo.Detail = reportDmoData
	detailReportDmo.Transactions = transactions
	detailReportDmo.GroupingVessels = groupingVesselDn

	return detailReportDmo, nil
}

// Report

func (r *repository) GetReport(year int, iupopkId int) (ReportRecapOutput, ReportDetailOutput, error) {
	var reportRecap ReportRecapOutput
	var reportDetail ReportDetailOutput

	var caloriesMinimum float64
	var caloriesMaximum float64
	var listTransactions []Transaction
	var listProduction []production.Production
	var companyElectricity []string
	var companyNonElectricity []string

	startFilter := fmt.Sprintf("%v-01-01", year)
	endFilter := fmt.Sprintf("%v-12-31", year)

	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)
	queryFilter += " AND transaction_type = DN" + "minerba_id IS NOT NULL AND shipping_date >= '" + startFilter + "' AND shipping_date <= '" + endFilter + "'"
	queryFilterProduction := "production_date >= '" + startFilter + "' AND production_date <= '" + endFilter + "'"
	errFind := r.db.Preload("Company.IndustryType").Where(queryFilter).Order("id ASC").Find(&listTransactions).Error
	errFindProduction := r.db.Where(queryFilterProduction).Order("id ASC").Find(&listProduction).Error

	if errFind != nil {
		return reportRecap, reportDetail, errFind
	}

	if errFindProduction != nil {
		return reportRecap, reportDetail, errFindProduction
	}

	reportDetail.Electricity.January = make(map[string]float64)
	reportDetail.NonElectricity.January = make(map[string]float64)
	reportDetail.Electricity.February = make(map[string]float64)
	reportDetail.NonElectricity.February = make(map[string]float64)
	reportDetail.Electricity.March = make(map[string]float64)
	reportDetail.NonElectricity.March = make(map[string]float64)
	reportDetail.Electricity.April = make(map[string]float64)
	reportDetail.NonElectricity.April = make(map[string]float64)
	reportDetail.Electricity.May = make(map[string]float64)
	reportDetail.NonElectricity.May = make(map[string]float64)
	reportDetail.Electricity.June = make(map[string]float64)
	reportDetail.NonElectricity.June = make(map[string]float64)
	reportDetail.Electricity.July = make(map[string]float64)
	reportDetail.NonElectricity.July = make(map[string]float64)
	reportDetail.Electricity.August = make(map[string]float64)
	reportDetail.NonElectricity.August = make(map[string]float64)
	reportDetail.Electricity.September = make(map[string]float64)
	reportDetail.NonElectricity.September = make(map[string]float64)
	reportDetail.Electricity.October = make(map[string]float64)
	reportDetail.NonElectricity.October = make(map[string]float64)
	reportDetail.Electricity.November = make(map[string]float64)
	reportDetail.NonElectricity.November = make(map[string]float64)
	reportDetail.Electricity.December = make(map[string]float64)
	reportDetail.NonElectricity.December = make(map[string]float64)

	for i, v := range listTransactions {
		date, _ := time.Parse("2006-01-02T00:00:00Z", *v.ShippingDate)
		_, month, _ := date.Date()

		if i == 0 {
			caloriesMinimum = v.QualityCaloriesAr
			caloriesMaximum = v.QualityCaloriesAr
		} else {
			if v.QualityCaloriesAr < caloriesMinimum {
				caloriesMinimum = v.QualityCaloriesAr
			}

			if v.QualityCaloriesAr > caloriesMaximum {
				caloriesMaximum = v.QualityCaloriesAr
			}
		}

		if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
			reportRecap.ElectricityTotal += v.QuantityUnloading
			reportRecap.Total += v.QuantityUnloading
		} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
			reportRecap.NonElectricityTotal += v.QuantityUnloading
			reportRecap.Total += v.QuantityUnloading
		}

		if v.IsNotClaim == false {
			switch int(month) {
			case 1:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.January[v.DmoBuyer.CompanyName]; ok {
						reportDetail.Electricity.January[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.Electricity.January[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := reportDetail.Electricity.January[v.DmoBuyer.CompanyName]; ok {
						reportDetail.NonElectricity.January[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.NonElectricity.January[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					}
				}
			case 2:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.February[v.DmoBuyer.CompanyName]; ok {
						reportDetail.Electricity.February[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.Electricity.February[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.February[v.DmoBuyer.CompanyName]; ok {
						reportDetail.NonElectricity.February[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.NonElectricity.February[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					}
				}
			case 3:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.March[v.DmoBuyer.CompanyName]; ok {
						reportDetail.Electricity.March[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.Electricity.March[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.March[v.DmoBuyer.CompanyName]; ok {
						reportDetail.NonElectricity.March[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.NonElectricity.March[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					}
				}
			case 4:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.April[v.DmoBuyer.CompanyName]; ok {
						reportDetail.Electricity.April[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.Electricity.April[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.April[v.DmoBuyer.CompanyName]; ok {
						reportDetail.NonElectricity.April[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.NonElectricity.April[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					}
				}
			case 5:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.May[v.DmoBuyer.CompanyName]; ok {
						reportDetail.Electricity.May[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.Electricity.May[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.May[v.DmoBuyer.CompanyName]; ok {
						reportDetail.NonElectricity.May[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.NonElectricity.May[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					}
				}
			case 6:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.June[v.DmoBuyer.CompanyName]; ok {
						reportDetail.Electricity.June[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.Electricity.June[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.June[v.DmoBuyer.CompanyName]; ok {
						reportDetail.NonElectricity.June[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.NonElectricity.June[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					}
				}
			case 7:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.July[v.DmoBuyer.CompanyName]; ok {
						reportDetail.Electricity.July[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.Electricity.July[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.July[v.DmoBuyer.CompanyName]; ok {
						reportDetail.NonElectricity.July[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.NonElectricity.July[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					}
				}
			case 8:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.August[v.DmoBuyer.CompanyName]; ok {
						reportDetail.Electricity.August[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.Electricity.August[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.August[v.DmoBuyer.CompanyName]; ok {
						reportDetail.NonElectricity.August[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.NonElectricity.August[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					}
				}
			case 9:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.September[v.DmoBuyer.CompanyName]; ok {
						reportDetail.Electricity.September[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.Electricity.September[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.September[v.DmoBuyer.CompanyName]; ok {
						reportDetail.NonElectricity.September[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.NonElectricity.September[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					}
				}
			case 10:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.October[v.DmoBuyer.CompanyName]; ok {
						reportDetail.Electricity.October[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.Electricity.October[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.October[v.DmoBuyer.CompanyName]; ok {
						reportDetail.NonElectricity.October[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.NonElectricity.October[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					}
				}
			case 11:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.November[v.DmoBuyer.CompanyName]; ok {
						reportDetail.Electricity.November[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.Electricity.November[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.November[v.DmoBuyer.CompanyName]; ok {
						reportDetail.NonElectricity.November[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						reportDetail.NonElectricity.November[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
					}
				}
			case 12:
				if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.December[v.DmoBuyer.CompanyName]; ok {
						reportDetail.Electricity.December[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.Electricity.December[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.Electricity.Total += v.QuantityUnloading
					}
				} else if v.DmoBuyer != nil && v.DmoBuyer.IndustryType != nil && v.DmoBuyer.IndustryType.Category == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.December[v.DmoBuyer.CompanyName]; ok {
						reportDetail.NonElectricity.December[v.DmoBuyer.CompanyName] += v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyer.CompanyName) && v.DmoBuyer.CompanyName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyer.CompanyName)
						}
						reportDetail.NonElectricity.December[v.DmoBuyer.CompanyName] = v.QuantityUnloading
						reportDetail.NonElectricity.Total += v.QuantityUnloading
					}
				}
			}
		} else {
			switch int(month) {
			case 1:
				reportDetail.NotClaimable.January += v.QuantityUnloading
				reportDetail.NotClaimable.Total += v.QuantityUnloading
			case 2:
				reportDetail.NotClaimable.February += v.QuantityUnloading
				reportDetail.NotClaimable.Total += v.QuantityUnloading
			case 3:
				reportDetail.NotClaimable.March += v.QuantityUnloading
				reportDetail.NotClaimable.Total += v.QuantityUnloading
			case 4:
				reportDetail.NotClaimable.April += v.QuantityUnloading
				reportDetail.NotClaimable.Total += v.QuantityUnloading
			case 5:
				reportDetail.NotClaimable.May += v.QuantityUnloading
				reportDetail.NotClaimable.Total += v.QuantityUnloading
			case 6:
				reportDetail.NotClaimable.June += v.QuantityUnloading
				reportDetail.NotClaimable.Total += v.QuantityUnloading
			case 7:
				reportDetail.NotClaimable.July += v.QuantityUnloading
				reportDetail.NotClaimable.Total += v.QuantityUnloading
			case 8:
				reportDetail.NotClaimable.August += v.QuantityUnloading
				reportDetail.NotClaimable.Total += v.QuantityUnloading
			case 9:
				reportDetail.NotClaimable.September += v.QuantityUnloading
				reportDetail.NotClaimable.Total += v.QuantityUnloading
			case 10:
				reportDetail.NotClaimable.October += v.QuantityUnloading
				reportDetail.NotClaimable.Total += v.QuantityUnloading
			case 11:
				reportDetail.NotClaimable.November += v.QuantityUnloading
				reportDetail.NotClaimable.Total += v.QuantityUnloading
			case 12:
				reportDetail.NotClaimable.December += v.QuantityUnloading
				reportDetail.NotClaimable.Total += v.QuantityUnloading
			}
		}
	}

	for _, v := range listProduction {
		date, _ := time.Parse("2006-01-02T00:00:00Z", v.ProductionDate)
		_, month, _ := date.Date()
		switch int(month) {
		case 1:
			reportDetail.Production.January += v.Quantity
			reportDetail.Production.Total += v.Quantity
		case 2:
			reportDetail.Production.February += v.Quantity
			reportDetail.Production.Total += v.Quantity
		case 3:
			reportDetail.Production.March += v.Quantity
			reportDetail.Production.Total += v.Quantity
		case 4:
			reportDetail.Production.April += v.Quantity
			reportDetail.Production.Total += v.Quantity
		case 5:
			reportDetail.Production.May += v.Quantity
			reportDetail.Production.Total += v.Quantity
		case 6:
			reportDetail.Production.June += v.Quantity
			reportDetail.Production.Total += v.Quantity
		case 7:
			reportDetail.Production.July += v.Quantity
			reportDetail.Production.Total += v.Quantity
		case 8:
			reportDetail.Production.August += v.Quantity
			reportDetail.Production.Total += v.Quantity
		case 9:
			reportDetail.Production.September += v.Quantity
			reportDetail.Production.Total += v.Quantity
		case 10:
			reportDetail.Production.October += v.Quantity
			reportDetail.Production.Total += v.Quantity
		case 11:
			reportDetail.Production.November += v.Quantity
			reportDetail.Production.Total += v.Quantity
		case 12:
			reportDetail.Production.December += v.Quantity
			reportDetail.Production.Total += v.Quantity
		}
	}

	stringTempElectricityTotal := fmt.Sprintf("%.3f", reportDetail.Electricity.Total)
	parseTempElectricityTotal, _ := strconv.ParseFloat(stringTempElectricityTotal, 64)

	stringTempNonElectricityTotal := fmt.Sprintf("%.3f", reportDetail.NonElectricity.Total)
	parseTempNonElectricityTotal, _ := strconv.ParseFloat(stringTempNonElectricityTotal, 64)

	stringTempNotClaimableTotal := fmt.Sprintf("%.3f", reportDetail.NotClaimable.Total)
	parseTempNotClaimableTotal, _ := strconv.ParseFloat(stringTempNotClaimableTotal, 64)

	stringTempProductionTotal := fmt.Sprintf("%.3f", reportDetail.Production.Total)
	parseTempProductionTotal, _ := strconv.ParseFloat(stringTempProductionTotal, 64)

	reportDetail.Electricity.Total = parseTempElectricityTotal
	reportDetail.NonElectricity.Total = parseTempNonElectricityTotal
	reportDetail.NotClaimable.Total = parseTempNotClaimableTotal
	reportDetail.Production.Total = parseTempProductionTotal
	reportDetail.ElectricityCompany = companyElectricity
	reportDetail.NonElectricityCompany = companyNonElectricity

	stringTempRecapElectricityTotal := fmt.Sprintf("%.3f", reportRecap.ElectricityTotal)
	parseTempRecapElectricityTotal, _ := strconv.ParseFloat(stringTempRecapElectricityTotal, 64)

	stringTempRecapNonElectricityTotal := fmt.Sprintf("%.3f", reportRecap.NonElectricityTotal)
	parseTempRecapNonElectricityTotal, _ := strconv.ParseFloat(stringTempRecapNonElectricityTotal, 64)

	stringTempRecapTotal := fmt.Sprintf("%.3f", reportRecap.Total)
	parseTempRecapTotal, _ := strconv.ParseFloat(stringTempRecapTotal, 64)

	reportRecap.ElectricityTotal = parseTempRecapElectricityTotal
	reportRecap.NonElectricityTotal = parseTempRecapNonElectricityTotal
	reportRecap.Total = parseTempRecapTotal

	var productionReality float64

	r.db.Model(production.Production{}).Where(queryFilterProduction).Select("sum(quantity)").Row().Scan(&productionReality)

	reportRecap.TotalProduction = productionReality
	reportRecap.FulfillmentOfProductionRealization = fmt.Sprintf("%.2f%%", reportRecap.Total/productionReality*100)
	reportRecap.RateCalories = fmt.Sprintf("%v - %v GAR", caloriesMinimum, caloriesMaximum)
	return reportRecap, reportDetail, nil
}

func (r *repository) GetListForReport(iupopkId int) (ListForCreatingReportDmoOutput, error) {
	var list ListForCreatingReportDmoOutput

	var transactions []Transaction

	var groupingVessel []groupingvesseldn.GroupingVesselDn

	var salesSystemBarge []salessystem.SalesSystem
	var salesSystemBargeId []uint

	errFindSalesSystemBarge := r.db.Where("name ILIKE '%Barge'").Find(&salesSystemBarge).Error

	if errFindSalesSystemBarge != nil {
		return list, errFindSalesSystemBarge
	}

	for _, v := range salesSystemBarge {
		salesSystemBargeId = append(salesSystemBargeId, v.ID)
	}

	errFindTransaction := r.db.Preload(clause.Associations).Where("report_dmo_id IS NULL AND sales_system_id IN ? AND is_finance_check = ? AND transaction_type = ? AND seller_id = ?", salesSystemBargeId, true, "DN", iupopkId).Find(&transactions).Error

	if errFindTransaction != nil {
		return list, errFindTransaction
	}

	errFindGroupingVessel := r.db.Preload(clause.Associations).Where("report_dmo_id IS NULL AND sales_system = ? AND iupopk_id = ?", "Vessel", iupopkId).Find(&groupingVessel).Error

	if errFindGroupingVessel != nil {
		return list, errFindGroupingVessel
	}

	list.Transactions = transactions
	list.GroupingVessels = groupingVessel

	return list, nil
}

// Grouping Vessel Dn
func (r *repository) GetDetailGroupingVesselDn(id int, iupopkId int) (DetailGroupingVesselDn, error) {

	var detailGroupingVesselDn DetailGroupingVesselDn

	var groupingVesselDn groupingvesseldn.GroupingVesselDn
	var transactions []Transaction

	findGroupingVesselDnErr := r.db.Preload(clause.Associations).Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&groupingVesselDn).Error

	if findGroupingVesselDnErr != nil {
		return detailGroupingVesselDn, findGroupingVesselDnErr
	}

	detailGroupingVesselDn.Detail = groupingVesselDn

	transactionFindErr := r.db.Order("id desc").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("grouping_vessel_dn_id = ? AND seller_id = ?", id, iupopkId).Find(&transactions).Error

	if transactionFindErr != nil {
		return detailGroupingVesselDn, transactionFindErr
	}

	detailGroupingVesselDn.ListTransactions = transactions

	return detailGroupingVesselDn, nil
}

func (r *repository) ListDataDnWithoutGroup(iupopkId int) (ListTransactionNotHaveGroupingVessel, error) {
	var listGroup ListTransactionNotHaveGroupingVessel
	var transactionBarge []Transaction
	var transactionVessel []Transaction

	var salesSystemBarge []salessystem.SalesSystem
	var salesSystemBargeId []uint

	errFindSalesSystemBarge := r.db.Where("name ILIKE '%Barge'").Find(&salesSystemBarge).Error

	if errFindSalesSystemBarge != nil {
		return listGroup, errFindSalesSystemBarge
	}

	for _, v := range salesSystemBarge {
		salesSystemBargeId = append(salesSystemBargeId, v.ID)
	}

	var salesSystemVessel []salessystem.SalesSystem
	var salesSystemVesselId []uint

	errFindSalesSystemVessel := r.db.Where("name ILIKE '%Vessel'").Find(&salesSystemVessel).Error

	if errFindSalesSystemVessel != nil {
		return listGroup, errFindSalesSystemVessel
	}

	for _, v := range salesSystemVessel {
		salesSystemVesselId = append(salesSystemVesselId, v.ID)
	}

	findTransactionBargeErr := r.db.Order("id desc").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("transaction_type = ? AND is_not_claim = ? AND is_migration = ? AND grouping_vessel_dn_id is NULL AND sales_system_id IN ? AND vessel_id IS NOT NULL AND seller_id = ? AND dmo_id IS NULL", "DN", false, false, salesSystemBargeId, iupopkId).Find(&transactionBarge).Error

	if findTransactionBargeErr != nil {
		return listGroup, findTransactionBargeErr
	}

	findTransactionVesselErr := r.db.Order("id desc").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("transaction_type = ? AND is_not_claim = ? AND is_migration = ? AND grouping_vessel_dn_id is NULL AND sales_system_id IN ? AND vessel_id IS NOT NULL AND seller_id = ? AND dmo_id IS NULL", "DN", false, false, salesSystemVesselId, iupopkId).Find(&transactionVessel).Error

	if findTransactionVesselErr != nil {
		return listGroup, findTransactionVesselErr
	}
	listGroup.TransactionBarge = transactionBarge
	listGroup.TransactionVessel = transactionVessel

	return listGroup, nil
}

// Grouping Vessel Ln
func (r *repository) GetDetailGroupingVesselLn(id int, iupopkId int) (DetailGroupingVesselLn, error) {

	var detailGroupingVesselLn DetailGroupingVesselLn

	var groupingVesselLn groupingvesselln.GroupingVesselLn
	var transactions []Transaction

	findGroupingVesselLnErr := r.db.Preload(clause.Associations).Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&groupingVesselLn).Error

	if findGroupingVesselLnErr != nil {
		return detailGroupingVesselLn, findGroupingVesselLnErr
	}

	detailGroupingVesselLn.Detail = groupingVesselLn

	transactionFindErr := r.db.Order("id desc").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("grouping_vessel_ln_id = ? AND seller_id = ?", id, iupopkId).Find(&transactions).Error

	if transactionFindErr != nil {
		return detailGroupingVesselLn, transactionFindErr
	}

	detailGroupingVesselLn.ListTransactions = transactions

	return detailGroupingVesselLn, nil
}

func (r *repository) ListDataLnWithoutGroup(iupopkId int) ([]Transaction, error) {
	var listDataLnWithoutGrouping []Transaction

	errListDataLnWithoutGrouping := r.db.Order("id desc").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("transaction_type = ? AND is_not_claim = ? AND is_migration = ? AND grouping_vessel_ln_id is NULL AND seller_id = ? AND dmo_id IS NULL", "LN", false, false, iupopkId).Find(&listDataLnWithoutGrouping).Error

	if errListDataLnWithoutGrouping != nil {
		return listDataLnWithoutGrouping, errListDataLnWithoutGrouping
	}

	return listDataLnWithoutGrouping, nil
}

// Minerba LN

func (r *repository) ListDataLNWithoutMinerba(iupopkId int) ([]Transaction, error) {
	var listDataLnWithoutMinerba []Transaction

	errFind := r.db.Order("id desc").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("minerba_ln_id is NULL AND transaction_type = ? AND is_not_claim = ? AND is_migration = ? AND is_finance_check = ? AND seller_id = ?", "LN", false, false, true, iupopkId).Find(&listDataLnWithoutMinerba).Error

	return listDataLnWithoutMinerba, errFind
}

func (r *repository) GetDetailMinerbaLn(id int, iupopkId int) (DetailMinerbaLn, error) {

	var detailMinerbaLn DetailMinerbaLn

	var minerbaLn minerbaln.MinerbaLn
	var transactions []Transaction

	minerbaLnFindErr := r.db.Preload(clause.Associations).Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&minerbaLn).Error

	if minerbaLnFindErr != nil {
		return detailMinerbaLn, minerbaLnFindErr
	}

	detailMinerbaLn.Detail = minerbaLn

	transactionFindErr := r.db.Order("shipping_date asc").Preload(clause.Associations).Preload("LoadingPort.PortLocation").Preload("UnloadingPort.PortLocation").Preload("DmoBuyer.IndustryType").Where("minerba_ln_id = ? AND seller_id = ?", id, iupopkId).Find(&transactions).Error

	if transactionFindErr != nil {
		return detailMinerbaLn, transactionFindErr
	}

	detailMinerbaLn.List = transactions
	return detailMinerbaLn, nil
}

func (r *repository) CheckDataLnAndMinerbaLnUpdate(listData []int, idMinerba int, iupopkId int) ([]Transaction, error) {
	var listLnValid []Transaction

	errFindValid := r.db.Where("id IN ? AND seller_id = ?", listData, iupopkId).Find(&listLnValid).Error

	if errFindValid != nil {
		return listLnValid, errFindValid
	}

	if len(listData) != len(listLnValid) {
		return listLnValid, errors.New("please check there is transaction not found")
	}

	var listLn []Transaction

	errFind := r.db.Where("id IN ? AND seller_id = ?", listData, iupopkId).Find(&listLn).Error

	if errFind != nil {
		return listLn, errFind
	}

	uintIdMinerba := uint(idMinerba)

	for _, v := range listLn {
		if v.MinerbaId != nil && *v.MinerbaId != uintIdMinerba {
			return listLn, errors.New("please check there is transaction already in report")
		}
	}

	return listLn, nil
}

func (r *repository) CheckDataLnAndMinerbaLn(listData []int, iupopkId int) (bool, error) {
	var listLnValid []Transaction

	errFindValid := r.db.Where("id IN ? AND seller_id = ?", listData, iupopkId).Find(&listLnValid).Error

	if errFindValid != nil {
		return false, errFindValid
	}

	if len(listData) != len(listLnValid) {
		return false, errors.New("please check there is transaction not found")
	}

	var listLn []Transaction

	errFind := r.db.Where("minerba_ln_id is NULL AND id IN ? AND seller_id = ?", listData, iupopkId).Find(&listLn).Error

	if errFind != nil {
		return false, errFind
	}

	if len(listLn) != len(listData) {
		return false, errors.New("please check there is transaction already in report")
	}

	return true, nil
}
