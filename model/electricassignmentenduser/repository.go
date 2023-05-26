package electricassignmentenduser

import (
	"ajebackend/model/electricassignment"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	DetailElectricAssignment(id int, iupopkId int) (DetailElectricAssignment, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) DetailElectricAssignment(id int, iupopkId int) (DetailElectricAssignment, error) {
	var detailElectricAssignment DetailElectricAssignment

	var electricAssignment electricassignment.ElectricAssignment

	var listElectricAssignment []ElectricAssignmentEndUser
	// RealizationListEndUser
	errFind := r.db.Preload(clause.Associations).Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&electricAssignment).Error

	if errFind != nil {
		return detailElectricAssignment, errFind
	}

	detailElectricAssignment.Detail = electricAssignment

	errFindList := r.db.Preload(clause.Associations).Preload("Port.PortLocation").Where("electric_assignment_id = ?", id).Find(&listElectricAssignment).Error

	if errFindList != nil {
		return detailElectricAssignment, errFindList
	}

	detailElectricAssignment.ListEndUser = listElectricAssignment

	for _, value := range listElectricAssignment {
		var realization RealizationEndUser
		realization.PortId = value.PortId
		realization.Port = value.Port
		realization.Supplier = value.Supplier
		realization.AverageCalories = value.AverageCalories
		realization.Quantity = value.Quantity
		realization.EndUser = value.EndUser
		realization.ID = value.ID
		realization.LetterNumber = value.LetterNumber
		var transactionRealization Realization

		shippingDateFrom := fmt.Sprintf("%s-01-01", electricAssignment.Year)
		shippingDateTo := fmt.Sprintf("%s-12-31", electricAssignment.Year)

		errTrRealization := r.db.Table("transactions").Select("SUM(quantity_unloading) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories ").Where("transaction_type = ? AND seller_id = ? AND is_not_claim = ? AND dmo_destination_port_id = ? AND shipping_date >= ? AND shipping_date <= ? AND dmo_id IS NOT NULL", "DN", iupopkId, false, value.PortId, shippingDateFrom, shippingDateTo).Scan(&transactionRealization).Error

		if errTrRealization != nil {
			return detailElectricAssignment, errTrRealization
		}

		realization.RealizationQuantity = transactionRealization.RealizationQuantity
		realization.RealizationAverageCalories = transactionRealization.RealizationAverageCalories

		detailElectricAssignment.ListRealizationEndUser = append(detailElectricAssignment.ListRealizationEndUser, realization)
	}

	return detailElectricAssignment, nil
}
