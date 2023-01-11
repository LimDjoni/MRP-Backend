package transaction

import (
	"ajebackend/model/dmo"
	"ajebackend/model/minerba"
	"ajebackend/model/production"
	"errors"
	"fmt"
	"strconv"
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
	ListDataDN(page int, sortFilter SortAndFilter) (Pagination, error)
	DetailTransactionDN(id int) (Transaction, error)
	ListDataDNWithoutMinerba() ([]Transaction, error)
	CheckDataDnAndMinerba(listData []int) (bool, error)
	CheckDataDnAndMinerbaUpdate(listData []int, idMinerba int) ([]Transaction, error)
	GetDetailMinerba(id int) (DetailMinerba, error)
	ListDataDNWithoutDmo() (ChooseTransactionDmo, error)
	CheckDataDnAndDmo(listData []int) ([]Transaction, error)
	GetDetailDmo(id int) (DetailDmo, error)
	CheckDataUnique(inputTrans DataTransactionInput) (bool, bool, bool, bool)
	GetReport(year int) (ReportRecapOutput, ReportDetailOutput, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// Transaction

func (r *repository) ListDataDN(page int, sortFilter SortAndFilter) (Pagination, error) {
	var transactions []Transaction
	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	defaultSort := "id desc"
	sortString := fmt.Sprintf("%s %s", sortFilter.Field, sortFilter.Sort)
	if sortFilter.Field == "" || sortFilter.Sort == "" {
		sortString = defaultSort
	}

	queryFilter := fmt.Sprintf("transaction_type = '%s' ", "DN")

	if sortFilter.TugboatName != "" {
		queryFilter = queryFilter + " AND tugboat_name = '" + sortFilter.TugboatName + "'"
	}

	if sortFilter.BargeName != "" {
		queryFilter = queryFilter + " AND barge_name = '" + sortFilter.BargeName + "'"
	}

	if sortFilter.VesselName != "" {
		queryFilter = queryFilter + " AND vessel_name = '" + sortFilter.VesselName + "'"
	}

	if sortFilter.ShippingFrom != "" {
		queryFilter = queryFilter + " AND shipping_date >= '" + sortFilter.ShippingFrom + "'"
	}

	if sortFilter.ShippingTo != "" {
		queryFilter = queryFilter + " AND shipping_date <= '" + sortFilter.ShippingTo + "T23:59:59'"
	}

	if sortFilter.Quantity != 0 {
		quantity := fmt.Sprintf("%v", sortFilter.Quantity)
		queryFilter = queryFilter + " AND cast(quantity AS TEXT) LIKE '%" + quantity + "%'"
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

	errFind := r.db.Preload(clause.Associations).Where(queryFilter).Order(sortString).Scopes(paginateDataDN(transactions, &pagination, r.db, queryFilter)).Find(&transactions).Error

	if errFind != nil {
		errWithoutOrder := r.db.Preload(clause.Associations).Where(queryFilter).Order(defaultSort).Scopes(paginateDataDN(transactions, &pagination, r.db, queryFilter)).Find(&transactions).Error

		if errWithoutOrder != nil {
			pagination.Data = transactions
			return pagination, errWithoutOrder
		}
	}

	pagination.Data = transactions

	return pagination, nil
}

func (r *repository) DetailTransactionDN(id int) (Transaction, error) {
	var transaction Transaction

	errFind := r.db.Where("id = ?", id).First(&transaction).Error

	return transaction, errFind
}

func (r *repository) ListDataDNWithoutMinerba() ([]Transaction, error) {
	var listDataDnWithoutMinerba []Transaction

	errFind := r.db.Order("id desc").Where("minerba_id is NULL AND transaction_type = ? AND is_not_claim = ? AND is_migration = ? AND is_finance_check = ?", "DN", false, false, true).Find(&listDataDnWithoutMinerba).Error

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

func (r *repository) CheckDataDnAndMinerba(listData []int) (bool, error) {
	var listDnValid []Transaction

	errFindValid := r.db.Where("id IN ?", listData).Find(&listDnValid).Error

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

func (r *repository) CheckDataDnAndMinerbaUpdate(listData []int, idMinerba int) ([]Transaction, error) {
	var listDnValid []Transaction

	errFindValid := r.db.Where("id IN ?", listData).Find(&listDnValid).Error

	if errFindValid != nil {
		return listDnValid, errFindValid
	}

	if len(listData) != len(listDnValid) {
		return listDnValid, errors.New("please check there is transaction not found")
	}

	var listDn []Transaction

	errFind := r.db.Where("id IN ?", listData).Find(&listDn).Error

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

func (r *repository) GetDetailMinerba(id int) (DetailMinerba, error) {

	var detailMinerba DetailMinerba

	var minerba minerba.Minerba
	var transactions []Transaction

	minerbaFindErr := r.db.Where("id = ?", id).First(&minerba).Error

	if minerbaFindErr != nil {
		return detailMinerba, minerbaFindErr
	}

	detailMinerba.Detail = minerba

	transactionFindErr := r.db.Order("id desc").Where("minerba_id = ?", id).Find(&transactions).Error

	if transactionFindErr != nil {
		return detailMinerba, transactionFindErr
	}

	detailMinerba.List = transactions
	return detailMinerba, nil
}

// DMO

func (r *repository) ListDataDNWithoutDmo() (ChooseTransactionDmo, error) {
	var listDataDnBargeDmo []Transaction
	var listDataDnVesselDmo []Transaction
	var listDataDnForDmo ChooseTransactionDmo

	errFindBarge := r.db.Order("id desc").Where("dmo_id is NULL AND transaction_type = ? AND is_not_claim = ? AND is_migration = ? AND vessel_name = ? AND is_finance_check = ?", "DN", false, false, "", true).Find(&listDataDnBargeDmo).Error

	if errFindBarge != nil {
		return listDataDnForDmo, errFindBarge
	}

	errFindVessel := r.db.Order("id desc").Where("dmo_id is NULL AND transaction_type = ? AND is_not_claim = ? AND is_migration = ? AND vessel_name != ? AND is_finance_check = ?", "DN", false, false, "", true).Find(&listDataDnVesselDmo).Error

	if errFindVessel != nil {
		return listDataDnForDmo, errFindVessel
	}

	listDataDnForDmo.BargeTransaction = listDataDnBargeDmo
	listDataDnForDmo.VesselTransaction = listDataDnVesselDmo
	return listDataDnForDmo, nil
}

func (r *repository) CheckDataDnAndDmo(listData []int) ([]Transaction, error) {
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

func (r *repository) GetDetailDmo(id int) (DetailDmo, error) {

	var detailDmo DetailDmo

	var dmoData dmo.Dmo
	var transactions []Transaction

	dmoFindErr := r.db.Where("id = ?", id).First(&dmoData).Error

	if dmoFindErr != nil {
		return detailDmo, dmoFindErr
	}

	detailDmo.Detail = dmoData

	transactionFindErr := r.db.Order("id desc").Where("dmo_id = ?", id).Find(&transactions).Error

	if transactionFindErr != nil {
		return detailDmo, transactionFindErr
	}

	detailDmo.List = transactions
	return detailDmo, nil
}

// Report

func (r *repository) GetReport(year int) (ReportRecapOutput, ReportDetailOutput, error) {
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

	queryFilter := "minerba_id IS NOT NULL AND shipping_date >= '" + startFilter + "' AND shipping_date <= '" + endFilter + "'"
	queryFilterProduction := "production_date >= '" + startFilter + "' AND production_date <= '" + endFilter + "'"
	errFind := r.db.Where(queryFilter).Order("id ASC").Find(&listTransactions).Error
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

		if v.DmoCategory == "ELECTRICITY" {
			reportRecap.ElectricityTotal += v.Quantity
			reportRecap.Total += v.Quantity
		} else if v.DmoCategory == "NON ELECTRICITY" {
			reportRecap.NonElectricityTotal += v.Quantity
			reportRecap.Total += v.Quantity
		}

		if v.IsNotClaim == false {
			switch int(month) {
			case 1:
				if v.DmoCategory == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.January[v.DmoBuyerName]; ok {
						reportDetail.Electricity.January[v.DmoBuyerName] += v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					} else {
						if !helperString(companyElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyerName)
						}
						reportDetail.Electricity.January[v.DmoBuyerName] = v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					}
				} else if v.DmoCategory == "NON ELECTRICITY" {
					if _, ok := reportDetail.Electricity.January[v.DmoBuyerName]; ok {
						reportDetail.NonElectricity.January[v.DmoBuyerName] += v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyerName)
						}
						reportDetail.NonElectricity.January[v.DmoBuyerName] = v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					}
				}
			case 2:
				if v.DmoCategory == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.February[v.DmoBuyerName]; ok {
						reportDetail.Electricity.February[v.DmoBuyerName] += v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					} else {
						if !helperString(companyElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyerName)
						}
						reportDetail.Electricity.February[v.DmoBuyerName] = v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					}
				} else if v.DmoCategory == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.February[v.DmoBuyerName]; ok {
						reportDetail.NonElectricity.February[v.DmoBuyerName] += v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyerName)
						}
						reportDetail.NonElectricity.February[v.DmoBuyerName] = v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					}
				}
			case 3:
				if v.DmoCategory == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.March[v.DmoBuyerName]; ok {
						reportDetail.Electricity.March[v.DmoBuyerName] += v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					} else {
						if !helperString(companyElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyerName)
						}
						reportDetail.Electricity.March[v.DmoBuyerName] = v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					}
				} else if v.DmoCategory == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.March[v.DmoBuyerName]; ok {
						reportDetail.NonElectricity.March[v.DmoBuyerName] += v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyerName)
						}
						reportDetail.NonElectricity.March[v.DmoBuyerName] = v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					}
				}
			case 4:
				if v.DmoCategory == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.April[v.DmoBuyerName]; ok {
						reportDetail.Electricity.April[v.DmoBuyerName] += v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					} else {
						if !helperString(companyElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyerName)
						}
						reportDetail.Electricity.April[v.DmoBuyerName] = v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					}
				} else if v.DmoCategory == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.April[v.DmoBuyerName]; ok {
						reportDetail.NonElectricity.April[v.DmoBuyerName] += v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyerName)
						}
						reportDetail.NonElectricity.April[v.DmoBuyerName] = v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					}
				}
			case 5:
				if v.DmoCategory == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.May[v.DmoBuyerName]; ok {
						reportDetail.Electricity.May[v.DmoBuyerName] += v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					} else {
						if !helperString(companyElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyerName)
						}
						reportDetail.Electricity.May[v.DmoBuyerName] = v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					}
				} else if v.DmoCategory == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.May[v.DmoBuyerName]; ok {
						reportDetail.NonElectricity.May[v.DmoBuyerName] += v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyerName)
						}
						reportDetail.NonElectricity.May[v.DmoBuyerName] = v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					}
				}
			case 6:
				if v.DmoCategory == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.June[v.DmoBuyerName]; ok {
						reportDetail.Electricity.June[v.DmoBuyerName] += v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					} else {
						if !helperString(companyElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyerName)
						}
						reportDetail.Electricity.June[v.DmoBuyerName] = v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					}
				} else if v.DmoCategory == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.June[v.DmoBuyerName]; ok {
						reportDetail.NonElectricity.June[v.DmoBuyerName] += v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyerName)
						}
						reportDetail.NonElectricity.June[v.DmoBuyerName] = v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					}
				}
			case 7:
				if v.DmoCategory == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.July[v.DmoBuyerName]; ok {
						reportDetail.Electricity.July[v.DmoBuyerName] += v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					} else {
						if !helperString(companyElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyerName)
						}
						reportDetail.Electricity.July[v.DmoBuyerName] = v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					}
				} else if v.DmoCategory == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.July[v.DmoBuyerName]; ok {
						reportDetail.NonElectricity.July[v.DmoBuyerName] += v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyerName)
						}
						reportDetail.NonElectricity.July[v.DmoBuyerName] = v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					}
				}
			case 8:
				if v.DmoCategory == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.August[v.DmoBuyerName]; ok {
						reportDetail.Electricity.August[v.DmoBuyerName] += v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					} else {
						if !helperString(companyElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyerName)
						}
						reportDetail.Electricity.August[v.DmoBuyerName] = v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					}
				} else if v.DmoCategory == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.August[v.DmoBuyerName]; ok {
						reportDetail.NonElectricity.August[v.DmoBuyerName] += v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyerName)
						}
						reportDetail.NonElectricity.August[v.DmoBuyerName] = v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					}
				}
			case 9:
				if v.DmoCategory == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.September[v.DmoBuyerName]; ok {
						reportDetail.Electricity.September[v.DmoBuyerName] += v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					} else {
						if !helperString(companyElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyerName)
						}
						reportDetail.Electricity.September[v.DmoBuyerName] = v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					}
				} else if v.DmoCategory == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.September[v.DmoBuyerName]; ok {
						reportDetail.NonElectricity.September[v.DmoBuyerName] += v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyerName)
						}
						reportDetail.NonElectricity.September[v.DmoBuyerName] = v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					}
				}
			case 10:
				if v.DmoCategory == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.October[v.DmoBuyerName]; ok {
						reportDetail.Electricity.October[v.DmoBuyerName] += v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					} else {
						if !helperString(companyElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyerName)
						}
						reportDetail.Electricity.October[v.DmoBuyerName] = v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					}
				} else if v.DmoCategory == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.October[v.DmoBuyerName]; ok {
						reportDetail.NonElectricity.October[v.DmoBuyerName] += v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyerName)
						}
						reportDetail.NonElectricity.October[v.DmoBuyerName] = v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					}
				}
			case 11:
				if v.DmoCategory == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.November[v.DmoBuyerName]; ok {
						reportDetail.Electricity.November[v.DmoBuyerName] += v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					} else {
						if !helperString(companyElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyerName)
						}
						reportDetail.Electricity.November[v.DmoBuyerName] = v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					}
				} else if v.DmoCategory == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.November[v.DmoBuyerName]; ok {
						reportDetail.NonElectricity.November[v.DmoBuyerName] += v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					} else {
						reportDetail.NonElectricity.November[v.DmoBuyerName] = v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
						if !helperString(companyNonElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyerName)
						}
					}
				}
			case 12:
				if v.DmoCategory == "ELECTRICITY" {
					if _, ok := reportDetail.Electricity.December[v.DmoBuyerName]; ok {
						reportDetail.Electricity.December[v.DmoBuyerName] += v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					} else {
						if !helperString(companyElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyElectricity = append(companyElectricity, v.DmoBuyerName)
						}
						reportDetail.Electricity.December[v.DmoBuyerName] = v.Quantity
						reportDetail.Electricity.Total += v.Quantity
					}
				} else if v.DmoCategory == "NON ELECTRICITY" {
					if _, ok := reportDetail.NonElectricity.December[v.DmoBuyerName]; ok {
						reportDetail.NonElectricity.December[v.DmoBuyerName] += v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					} else {
						if !helperString(companyNonElectricity, v.DmoBuyerName) && v.DmoBuyerName != "" {
							companyNonElectricity = append(companyNonElectricity, v.DmoBuyerName)
						}
						reportDetail.NonElectricity.December[v.DmoBuyerName] = v.Quantity
						reportDetail.NonElectricity.Total += v.Quantity
					}
				}
			}
		} else {
			switch int(month) {
			case 1:
				reportDetail.NotClaimable.January += v.Quantity
				reportDetail.NotClaimable.Total += v.Quantity
			case 2:
				reportDetail.NotClaimable.February += v.Quantity
				reportDetail.NotClaimable.Total += v.Quantity
			case 3:
				reportDetail.NotClaimable.March += v.Quantity
				reportDetail.NotClaimable.Total += v.Quantity
			case 4:
				reportDetail.NotClaimable.April += v.Quantity
				reportDetail.NotClaimable.Total += v.Quantity
			case 5:
				reportDetail.NotClaimable.May += v.Quantity
				reportDetail.NotClaimable.Total += v.Quantity
			case 6:
				reportDetail.NotClaimable.June += v.Quantity
				reportDetail.NotClaimable.Total += v.Quantity
			case 7:
				reportDetail.NotClaimable.July += v.Quantity
				reportDetail.NotClaimable.Total += v.Quantity
			case 8:
				reportDetail.NotClaimable.August += v.Quantity
				reportDetail.NotClaimable.Total += v.Quantity
			case 9:
				reportDetail.NotClaimable.September += v.Quantity
				reportDetail.NotClaimable.Total += v.Quantity
			case 10:
				reportDetail.NotClaimable.October += v.Quantity
				reportDetail.NotClaimable.Total += v.Quantity
			case 11:
				reportDetail.NotClaimable.November += v.Quantity
				reportDetail.NotClaimable.Total += v.Quantity
			case 12:
				reportDetail.NotClaimable.December += v.Quantity
				reportDetail.NotClaimable.Total += v.Quantity
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
