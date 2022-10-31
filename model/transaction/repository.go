package transaction

import (
	"ajebackend/model/dmo"
	"ajebackend/model/minerba"
	"ajebackend/model/production"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Repository interface {
	ListDataDN(page int, sortFilter SortAndFilter) (Pagination, error)
	DetailTransactionDN(id int) (Transaction, error)
	ListDataDNWithoutMinerba() ([]Transaction, error)
	CheckDataDnAndMinerba(listData []int)(bool, error)
	CheckDataDnAndMinerbaUpdate(listData []int, idMinerba int)([]Transaction, error)
	GetDetailMinerba(id int)(DetailMinerba, error)
	ListDataDNWithoutDmo() ([]Transaction, error)
	CheckDataDnAndDmo(listData []int)([]Transaction, error)
	GetDetailDmo(id int)(DetailDmo, error)
	CheckDataUnique(inputTrans DataTransactionInput) (bool,bool,bool,bool)
	GetReportDetail(year int) (ReportDetailOutput, error)
	GetReportRecap(year int) (ReportRecapOutput, error)
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
	pagination.Limit = 10
	pagination.Page = page
	defaultSort := "id desc"
	sortString := fmt.Sprintf("%s %s", sortFilter.Field, sortFilter.Sort)
	if sortFilter.Field == "" || sortFilter.Sort == "" {
		sortString = defaultSort
	}

	queryFilter := fmt.Sprintf("transaction_type = '%s' ", "DN")

	if sortFilter.TugboatName != "" {
		queryFilter = queryFilter + " AND tugboat_name ILIKE '%" + sortFilter.TugboatName + "%'"
	}

	if sortFilter.BargeName != "" {
		queryFilter = queryFilter + " AND barge_name ILIKE '%" + sortFilter.BargeName + "%'"
	}

	if sortFilter.VesselName != "" {
		queryFilter = queryFilter + " AND vessel_name ILIKE '%" + sortFilter.VesselName + "%'"
	}

	if sortFilter.ShippingFrom != "" {
		queryFilter = queryFilter + " AND shipping_date >= '" + sortFilter.ShippingFrom + "'"
	}

	if sortFilter.ShippingTo != "" {
		queryFilter = queryFilter + " AND shipping_date <= '" + sortFilter.ShippingTo + "T23:59:59'"
	}

	if sortFilter.Quantity != 0 {
		quantity := fmt.Sprintf("%v", sortFilter.Quantity)
		queryFilter = queryFilter + " AND cast(quantity AS TEXT) LIKE '%" +  quantity + "%'"
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

	errFind := r.db.Where("minerba_id is NULL AND transaction_type = ? AND is_not_claim = ?", "DN", false).Find(&listDataDnWithoutMinerba).Error

	return listDataDnWithoutMinerba, errFind
}

func (r *repository) CheckDataUnique(inputTrans DataTransactionInput) (bool,bool,bool,bool) {
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

func (r *repository) CheckDataDnAndMinerba(listData []int)(bool, error) {
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

func (r *repository) CheckDataDnAndMinerbaUpdate(listData []int, idMinerba int)([]Transaction, error) {
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

func(r *repository) GetDetailMinerba(id int)(DetailMinerba, error) {

	var detailMinerba DetailMinerba

	var minerba minerba.Minerba
	var transactions []Transaction

	minerbaFindErr := r.db.Where("id = ?", id).First(&minerba).Error

	if minerbaFindErr != nil {
		return detailMinerba, minerbaFindErr
	}

	detailMinerba.Detail = minerba

	transactionFindErr := r.db.Where("minerba_id = ?", id).Find(&transactions).Error

	if transactionFindErr != nil {
		return detailMinerba, transactionFindErr
	}

	detailMinerba.List = transactions
	return detailMinerba, nil
}

// DMO

func (r *repository) ListDataDNWithoutDmo() ([]Transaction, error) {
	var listDataDnWithoutDmo []Transaction

	errFind := r.db.Where("dmo_id is NULL AND transaction_type = ? AND is_not_claim = ?", "DN", false).Find(&listDataDnWithoutDmo).Error

	return listDataDnWithoutDmo, errFind
}

func (r *repository) CheckDataDnAndDmo(listData []int)([]Transaction, error) {
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

func (r *repository) GetDetailDmo(id int)(DetailDmo, error) {

	var detailDmo DetailDmo

	var dmoData dmo.Dmo
	var transactions []Transaction

	dmoFindErr := r.db.Where("id = ?", id).First(&dmoData).Error

	if dmoFindErr != nil {
		return detailDmo, dmoFindErr
	}

	detailDmo.Detail = dmoData

	transactionFindErr := r.db.Where("dmo_id = ?", id).Find(&transactions).Error

	if transactionFindErr != nil {
		return detailDmo, transactionFindErr
	}

	detailDmo.List = transactions
	return detailDmo, nil
}

// Report

func (r *repository) GetReportDetail(year int) (ReportDetailOutput, error) {
	var report ReportDetailOutput

	var listTransactions []Transaction
	var listProduction []production.Production
	startFilter := fmt.Sprintf("%v-01-01", year)
	endFilter := fmt.Sprintf("%v-12-31", year)

	queryFilter := "minerba_id IS NOT NULL OR is_not_claim = true AND shipping_date >= '" + startFilter + "' AND shipping_date <= '" + endFilter + "'"
	queryFilterProduction := "production_date >= '" + startFilter + "' AND production_date <= '" + endFilter + "'"
	errFind := r.db.Where(queryFilter).Order("id ASC").Find(&listTransactions).Error
	errFindProduction := r.db.Where(queryFilterProduction).Order("id ASC").Find(&listProduction).Error

	if errFind != nil {
		return  report, errFind
	}

	if errFindProduction != nil {
		return  report, errFindProduction
	}

	report.Electricity.January = make(map[string]float64)
	report.NonElectricity.January = make(map[string]float64)
	report.Electricity.February = make(map[string]float64)
	report.NonElectricity.February = make(map[string]float64)
	report.Electricity.March = make(map[string]float64)
	report.NonElectricity.March = make(map[string]float64)
	report.Electricity.April = make(map[string]float64)
	report.NonElectricity.April = make(map[string]float64)
	report.Electricity.May = make(map[string]float64)
	report.NonElectricity.May = make(map[string]float64)
	report.Electricity.June = make(map[string]float64)
	report.NonElectricity.June = make(map[string]float64)
	report.Electricity.July = make(map[string]float64)
	report.NonElectricity.July = make(map[string]float64)
	report.Electricity.August = make(map[string]float64)
	report.NonElectricity.August = make(map[string]float64)
	report.Electricity.September = make(map[string]float64)
	report.NonElectricity.September = make(map[string]float64)
	report.Electricity.October = make(map[string]float64)
	report.NonElectricity.October = make(map[string]float64)
	report.Electricity.November = make(map[string]float64)
	report.NonElectricity.November = make(map[string]float64)
	report.Electricity.December = make(map[string]float64)
	report.NonElectricity.December = make(map[string]float64)

	for _, v := range listTransactions {
		date, _ := time.Parse("2006-01-02T00:00:00Z", *v.ShippingDate)
		_, month, _ := date.Date()

		if v.IsNotClaim == false {
			switch int(month) {
				case 1:
					if v.DmoCategory == "ELECTRICITY" {
						if _, ok := report.Electricity.January[v.DmoBuyerName]; ok {
							report.Electricity.January[v.DmoBuyerName] += v.Quantity
							report.Electricity.Total += v.Quantity
						} else {
							report.Electricity.January[v.DmoBuyerName] = v.Quantity
							report.Electricity.Total += v.Quantity
						}
					} else if v.DmoCategory == "NON ELECTRICITY" {
						if _, ok := report.Electricity.January[v.DmoBuyerName]; ok {
							report.NonElectricity.January[v.DmoBuyerName] += v.Quantity
							report.NonElectricity.Total += v.Quantity
						} else {
							report.NonElectricity.January[v.DmoBuyerName] = v.Quantity
							report.NonElectricity.Total += v.Quantity
						}
					}
				case 2:
					if v.DmoCategory == "ELECTRICITY" {
						if _, ok := report.Electricity.February[v.DmoBuyerName]; ok {
							report.Electricity.February[v.DmoBuyerName] += v.Quantity
							report.Electricity.Total += v.Quantity
						} else {
							report.Electricity.February[v.DmoBuyerName] = v.Quantity
							report.Electricity.Total += v.Quantity
						}
					} else if v.DmoCategory == "NON ELECTRICITY" {
						if _, ok := report.NonElectricity.February[v.DmoBuyerName]; ok {
							report.NonElectricity.February[v.DmoBuyerName] += v.Quantity
							report.NonElectricity.Total += v.Quantity
						} else {
							report.NonElectricity.February[v.DmoBuyerName] = v.Quantity
							report.NonElectricity.Total += v.Quantity
						}
					}
				case 3:
					if v.DmoCategory == "ELECTRICITY" {
						if _, ok := report.Electricity.March[v.DmoBuyerName]; ok {
							report.Electricity.March[v.DmoBuyerName] += v.Quantity
							report.Electricity.Total += v.Quantity
						} else {
							report.Electricity.March[v.DmoBuyerName] = v.Quantity
							report.Electricity.Total += v.Quantity
						}
					} else if v.DmoCategory == "NON ELECTRICITY" {
						if _, ok := report.NonElectricity.March[v.DmoBuyerName]; ok {
							report.NonElectricity.March[v.DmoBuyerName] += v.Quantity
							report.NonElectricity.Total += v.Quantity
						} else {
							report.NonElectricity.March[v.DmoBuyerName] = v.Quantity
							report.NonElectricity.Total += v.Quantity
						}
					}
				case 4:
					if v.DmoCategory == "ELECTRICITY" {
						if _, ok := report.Electricity.April[v.DmoBuyerName]; ok {
							 report.Electricity.April[v.DmoBuyerName] += v.Quantity
							report.Electricity.Total += v.Quantity
						} else {
							report.Electricity.April[v.DmoBuyerName] = v.Quantity
							report.Electricity.Total += v.Quantity
						}
					} else if v.DmoCategory == "NON ELECTRICITY" {
						if _, ok := report.NonElectricity.April[v.DmoBuyerName]; ok {
							report.NonElectricity.April[v.DmoBuyerName] += v.Quantity
							report.NonElectricity.Total += v.Quantity
						} else {
							report.NonElectricity.April[v.DmoBuyerName] = v.Quantity
							report.NonElectricity.Total += v.Quantity
						}
					}
				case 5:
					if v.DmoCategory == "ELECTRICITY" {
						if _, ok := report.Electricity.May[v.DmoBuyerName]; ok {
							report.Electricity.May[v.DmoBuyerName] += v.Quantity
							report.Electricity.Total += v.Quantity
						} else {
							report.Electricity.May[v.DmoBuyerName] = v.Quantity
							report.Electricity.Total += v.Quantity
						}
					} else if v.DmoCategory == "NON ELECTRICITY" {
						if _, ok := report.NonElectricity.May[v.DmoBuyerName]; ok {
							report.NonElectricity.May[v.DmoBuyerName] += v.Quantity
							report.NonElectricity.Total += v.Quantity
						} else {
							report.NonElectricity.May[v.DmoBuyerName] = v.Quantity
							report.NonElectricity.Total += v.Quantity
						}
					}
				case 6:
					if v.DmoCategory == "ELECTRICITY" {
						if _, ok := report.Electricity.June[v.DmoBuyerName]; ok {
							report.Electricity.June[v.DmoBuyerName] += v.Quantity
							report.Electricity.Total += v.Quantity
						} else {
							report.Electricity.June[v.DmoBuyerName] = v.Quantity
							report.Electricity.Total += v.Quantity
						}
					} else if v.DmoCategory == "NON ELECTRICITY" {
						if _, ok := report.NonElectricity.June[v.DmoBuyerName]; ok {
							report.NonElectricity.June[v.DmoBuyerName] += v.Quantity
							report.NonElectricity.Total += v.Quantity
						} else {
							report.NonElectricity.June[v.DmoBuyerName] = v.Quantity
							report.NonElectricity.Total += v.Quantity
						}
					}
				case 7:
					if v.DmoCategory == "ELECTRICITY" {
						if _, ok := report.Electricity.July[v.DmoBuyerName]; ok {
							report.Electricity.July[v.DmoBuyerName] += v.Quantity
							report.Electricity.Total += v.Quantity
						} else {
							report.Electricity.July[v.DmoBuyerName] = v.Quantity
							report.Electricity.Total += v.Quantity
						}
					} else if v.DmoCategory == "NON ELECTRICITY" {
						if _, ok := report.NonElectricity.July[v.DmoBuyerName]; ok {
							report.NonElectricity.July[v.DmoBuyerName] += v.Quantity
							report.NonElectricity.Total += v.Quantity
						} else {
							report.NonElectricity.July[v.DmoBuyerName] = v.Quantity
							report.NonElectricity.Total += v.Quantity
						}
					}
				case 8:
					if v.DmoCategory == "ELECTRICITY" {
						if _, ok := report.Electricity.August[v.DmoBuyerName]; ok {
							report.Electricity.August[v.DmoBuyerName] += v.Quantity
							report.Electricity.Total += v.Quantity
						} else {
							report.Electricity.August[v.DmoBuyerName] = v.Quantity
							report.Electricity.Total += v.Quantity
						}
					} else if v.DmoCategory == "NON ELECTRICITY" {
						if _, ok := report.NonElectricity.August[v.DmoBuyerName]; ok {
							report.NonElectricity.August[v.DmoBuyerName] += v.Quantity
							report.NonElectricity.Total += v.Quantity
						} else {
							report.NonElectricity.August[v.DmoBuyerName] = v.Quantity
							report.NonElectricity.Total += v.Quantity
						}
					}
				case 9:
					if v.DmoCategory == "ELECTRICITY" {
						if _, ok := report.Electricity.September[v.DmoBuyerName]; ok {
							report.Electricity.September[v.DmoBuyerName] += v.Quantity
							report.Electricity.Total += v.Quantity
						} else {
							report.Electricity.September[v.DmoBuyerName] = v.Quantity
							report.Electricity.Total += v.Quantity
						}
					} else if v.DmoCategory == "NON ELECTRICITY" {
						if _, ok := report.NonElectricity.September[v.DmoBuyerName]; ok {
							report.NonElectricity.September[v.DmoBuyerName] += v.Quantity
							report.NonElectricity.Total += v.Quantity
						} else {
							report.NonElectricity.September[v.DmoBuyerName] = v.Quantity
							report.NonElectricity.Total += v.Quantity
						}
					}
				case 10:
					if v.DmoCategory == "ELECTRICITY" {
						if _, ok := report.Electricity.October[v.DmoBuyerName]; ok {
							report.Electricity.October[v.DmoBuyerName] += v.Quantity
							report.Electricity.Total += v.Quantity
						} else {
							report.Electricity.October[v.DmoBuyerName] = v.Quantity
							report.Electricity.Total += v.Quantity
						}
					} else if v.DmoCategory == "NON ELECTRICITY" {
						if _, ok := report.NonElectricity.October[v.DmoBuyerName]; ok {
							report.NonElectricity.October[v.DmoBuyerName] += v.Quantity
							report.NonElectricity.Total += v.Quantity
						} else {
							report.NonElectricity.October[v.DmoBuyerName] = v.Quantity
							report.NonElectricity.Total += v.Quantity
						}
					}
				case 11:
					if v.DmoCategory == "ELECTRICITY" {
						if _, ok := report.Electricity.November[v.DmoBuyerName]; ok {
							report.Electricity.November[v.DmoBuyerName] += v.Quantity
							report.Electricity.Total += v.Quantity
						} else {
							report.Electricity.November[v.DmoBuyerName] = v.Quantity
							report.Electricity.Total += v.Quantity
						}
					} else if v.DmoCategory == "NON ELECTRICITY" {
						if _, ok := report.NonElectricity.November[v.DmoBuyerName]; ok {
							report.NonElectricity.November[v.DmoBuyerName] += v.Quantity
							report.NonElectricity.Total += v.Quantity
						} else {
							report.NonElectricity.November[v.DmoBuyerName] = v.Quantity
							report.NonElectricity.Total += v.Quantity
						}
					}
				case 12:
					if v.DmoCategory == "ELECTRICITY" {
						if _, ok := report.Electricity.December[v.DmoBuyerName]; ok {
							report.Electricity.December[v.DmoBuyerName] += v.Quantity
							report.Electricity.Total += v.Quantity
						} else {
							report.Electricity.December[v.DmoBuyerName] = v.Quantity
							report.Electricity.Total += v.Quantity
						}
					} else if v.DmoCategory == "NON ELECTRICITY" {
						if _, ok := report.NonElectricity.December[v.DmoBuyerName]; ok {
							report.NonElectricity.December[v.DmoBuyerName] += v.Quantity
							report.NonElectricity.Total += v.Quantity
						} else {
							report.NonElectricity.December[v.DmoBuyerName] = v.Quantity
							report.NonElectricity.Total += v.Quantity
						}
					}
			}
		} else {
			switch int(month) {
			case 1:
				report.NotClaimable.January += v.Quantity
				report.NotClaimable.Total += v.Quantity
			case 2:
				report.NotClaimable.February += v.Quantity
				report.NotClaimable.Total += v.Quantity
			case 3:
				report.NotClaimable.March += v.Quantity
				report.NotClaimable.Total += v.Quantity
			case 4:
				report.NotClaimable.April += v.Quantity
				report.NotClaimable.Total += v.Quantity
			case 5:
				report.NotClaimable.May += v.Quantity
				report.NotClaimable.Total += v.Quantity
			case 6:
				report.NotClaimable.June += v.Quantity
				report.NotClaimable.Total += v.Quantity
			case 7:
				report.NotClaimable.July += v.Quantity
				report.NotClaimable.Total += v.Quantity
			case 8:
				report.NotClaimable.August += v.Quantity
				report.NotClaimable.Total += v.Quantity
			case 9:
				report.NotClaimable.September += v.Quantity
				report.NotClaimable.Total += v.Quantity
			case 10:
				report.NotClaimable.October += v.Quantity
				report.NotClaimable.Total += v.Quantity
			case 11:
				report.NotClaimable.November += v.Quantity
				report.NotClaimable.Total += v.Quantity
			case 12:
				report.NotClaimable.December += v.Quantity
				report.NotClaimable.Total += v.Quantity
			}
		}
	}

	for _, v := range listProduction {
		date, _ := time.Parse("2006-01-02T00:00:00Z", v.ProductionDate)
		_, month, _ := date.Date()
		switch int(month) {
		case 1:
			report.Production.January += v.Quantity
			report.Production.Total += v.Quantity
		case 2:
			report.Production.February += v.Quantity
			report.Production.Total += v.Quantity
		case 3:
			report.Production.March += v.Quantity
			report.Production.Total += v.Quantity
		case 4:
			report.Production.April += v.Quantity
			report.Production.Total += v.Quantity
		case 5:
			report.Production.May += v.Quantity
			report.Production.Total += v.Quantity
		case 6:
			report.Production.June += v.Quantity
			report.Production.Total += v.Quantity
		case 7:
			report.Production.July += v.Quantity
			report.Production.Total += v.Quantity
		case 8:
			report.Production.August += v.Quantity
			report.Production.Total += v.Quantity
		case 9:
			report.Production.September += v.Quantity
			report.Production.Total += v.Quantity
		case 10:
			report.Production.October += v.Quantity
			report.Production.Total += v.Quantity
		case 11:
			report.Production.November += v.Quantity
			report.Production.Total += v.Quantity
		case 12:
			report.Production.December += v.Quantity
			report.Production.Total += v.Quantity
		}
	}
	return report, nil
}

func (r *repository) GetReportRecap(year int) (ReportRecapOutput, error) {
	var report ReportRecapOutput

	var caloriesMinimum float64
	var caloriesMaximum float64
	var listTransactions []Transaction

	startFilter := fmt.Sprintf("%v-01-01", year)
	endFilter := fmt.Sprintf("%v-12-31", year)

	queryFilter := "minerba_id IS NOT NULL AND shipping_date >= '" + startFilter + "' AND shipping_date <= '" + endFilter + "'"
	queryFilterProduction := "production_date >= '" + startFilter + "' AND production_date <= '" + endFilter + "'"
	errFind := r.db.Where(queryFilter).Order("id ASC").Find(&listTransactions).Error

	if errFind != nil {
		return  report, errFind
	}

	for i, v := range listTransactions {
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
			report.ElectricityTotal += v.Quantity
			report.Total += v.Quantity
		} else if v.DmoCategory == "NON ELECTRICITY" {
			report.NonElectricityTotal += v.Quantity
			report.Total += v.Quantity
		}

	}

	var productionReality float64

	r.db.Model(production.Production{}).Where(queryFilterProduction).Select("sum(quantity)").Row().Scan(&productionReality)

	report.TotalProduction = productionReality
	report.FulfillmentOfProductionRealization = fmt.Sprintf("%.2f%%", report.Total / productionReality * 100)
	report.RateCalories = fmt.Sprintf("%v - %v GAR",caloriesMinimum, caloriesMaximum )
	return report, nil
}
