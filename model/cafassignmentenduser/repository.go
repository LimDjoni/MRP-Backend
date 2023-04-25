package cafassignmentenduser

import (
	"ajebackend/model/cafassignment"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	DetailCafAssignment(id int, iupopkId int) (DetailCafAssignment, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) DetailCafAssignment(id int, iupopkId int) (DetailCafAssignment, error) {
	var detailCafAssignment DetailCafAssignment

	var cafAssignment cafassignment.CafAssignment

	var listCafAssignment []CafAssignmentEndUser
	// RealizationListEndUser
	errFind := r.db.Preload(clause.Associations).Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&cafAssignment).Error

	if errFind != nil {
		return detailCafAssignment, errFind
	}

	detailCafAssignment.Detail = cafAssignment

	errFindList := r.db.Preload(clause.Associations).Where("caf_assignment_id = ?", id).Find(&listCafAssignment).Error

	if errFindList != nil {
		return detailCafAssignment, errFindList
	}

	detailCafAssignment.ListEndUser = listCafAssignment

	for _, value := range listCafAssignment {
		var realization RealizationEndUser
		realization.AverageCalories = value.AverageCalories
		realization.Quantity = value.Quantity
		realization.EndUserId = value.EndUserId
		realization.EndUser = value.EndUser
		realization.EndUserString = value.EndUserString
		realization.ID = value.ID
		var transactionRealization Realization

		shippingDateFrom := fmt.Sprintf("%s-01-01", cafAssignment.Year)
		shippingDateTo := fmt.Sprintf("%s-12-31", cafAssignment.Year)

		query := fmt.Sprintf("transactions.transaction_type = '%s' AND transactions.seller_id = %v AND transactions.is_not_claim = false AND transactions.shipping_date >= '%s' AND transactions.shipping_date <= '%s' AND company.company_name = '%s'", "DN", iupopkId, shippingDateFrom, shippingDateTo, value.EndUserString)

		errTrRealization := r.db.Table("transactions").Select("SUM(transactions.quantity_unloading) as realization_quantity, AVG(transactions.quality_calories_ar) as realization_average_calories ").Joins("left join companies company on company.id = transactions.dmo_buyer_id").Where(query).Scan(&transactionRealization).Error

		if errTrRealization != nil {
			return detailCafAssignment, errTrRealization
		}

		realization.RealizationQuantity = transactionRealization.RealizationQuantity
		realization.RealizationAverageCalories = transactionRealization.RealizationAverageCalories

		detailCafAssignment.ListRealizationEndUser = append(detailCafAssignment.ListRealizationEndUser, realization)
	}

	return detailCafAssignment, nil
}
