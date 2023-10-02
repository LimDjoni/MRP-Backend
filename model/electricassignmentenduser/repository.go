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

	shippingDateFrom := fmt.Sprintf("%s-01-01", electricAssignment.Year)
	shippingDateTo := fmt.Sprintf("%s-12-31", electricAssignment.Year)

	var listRealization []ListRealization
	if electricAssignment.LetterNumber != "" {
		var listAssignment []ElectricAssignmentEndUser
		var listRealizationTemp ListRealization

		listRealizationTemp.Order = 1
		listRealizationTemp.LetterNumber = electricAssignment.LetterNumber

		errFindList := r.db.Preload(clause.Associations).Preload("Port.PortLocation").Where("electric_assignment_id = ? AND letter_number = ?", id, electricAssignment.LetterNumber).Find(&listAssignment).Error

		if errFindList != nil {
			return detailElectricAssignment, errFindList
		}

		for _, value := range listAssignment {
			var realization RealizationEndUser
			realization.PortId = value.PortId
			realization.Port = value.Port
			realization.AverageCalories = value.AverageCalories
			realization.Quantity = value.Quantity
			realization.EndUser = value.EndUser
			realization.ID = value.ID
			realization.LetterNumber = value.LetterNumber
			var tempAssignment []ElectricAssignmentEndUser
			var transactionRealization Realization

			errTrRealization := r.db.Table("transactions").Select("SUM(quantity_unloading) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories").Where("transaction_type = ? AND seller_id = ? AND is_not_claim = ? AND dmo_destination_port_id = ? AND shipping_date >= ? AND shipping_date <= ? AND dmo_id IS NOT NULL AND grouping_vessel_dn_id IS NULL AND report_dmo_id IS NOT NULL", "DN", iupopkId, false, value.PortId, shippingDateFrom, shippingDateTo).Scan(&transactionRealization).Error

			if errTrRealization != nil {
				return detailElectricAssignment, errTrRealization
			}

			var groupingRealization Realization

			var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns
where id in (select grouping_vessel_dn_id from transactions where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and transaction_type = 'DN' and is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND dmo_destination_port_id = %v AND bl_date >= '%s' AND bl_date <= '%s' AND iupopk_id = %v AND report_dmo_id IS NOT NULL
				`, value.PortId, shippingDateFrom, shippingDateTo, iupopkId)

			errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

			if errGroupingRealization != nil {
				return detailElectricAssignment, errGroupingRealization
			}

			transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
			transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

			errFindTemp := r.db.Where("port_id = ? AND electric_assignment_id = ?", value.PortId, electricAssignment.ID).Find(&tempAssignment).Error

			if errFindTemp != nil {
				return detailElectricAssignment, errFindTemp
			}

			if transactionRealization.RealizationQuantity > value.Quantity {

				if len(tempAssignment) > 1 {
					realization.RealizationQuantity = value.Quantity
				} else {
					realization.RealizationQuantity = transactionRealization.RealizationQuantity
				}
			} else {
				realization.RealizationQuantity = transactionRealization.RealizationQuantity
			}

			realization.RealizationAverageCalories = transactionRealization.RealizationAverageCalories

			listRealizationTemp.ListRealizationEndUser = append(listRealizationTemp.ListRealizationEndUser, realization)
		}

		listRealization = append(listRealization, listRealizationTemp)
	}

	if electricAssignment.LetterNumber2 != "" {
		var listAssignment []ElectricAssignmentEndUser
		var listRealizationTemp ListRealization

		listRealizationTemp.Order = 2
		listRealizationTemp.LetterNumber = electricAssignment.LetterNumber2

		errFindList := r.db.Preload(clause.Associations).Preload("Port.PortLocation").Where("electric_assignment_id = ? AND letter_number = ?", id, electricAssignment.LetterNumber2).Find(&listAssignment).Error

		if errFindList != nil {
			return detailElectricAssignment, errFindList
		}

		for _, value := range listAssignment {
			var realization RealizationEndUser
			realization.PortId = value.PortId
			realization.Port = value.Port
			realization.AverageCalories = value.AverageCalories
			realization.Quantity = value.Quantity
			realization.EndUser = value.EndUser
			realization.ID = value.ID
			realization.LetterNumber = value.LetterNumber

			var transactionRealization Realization
			var tempAssignment []ElectricAssignmentEndUser

			errTrRealization := r.db.Table("transactions").Select("SUM(quantity_unloading) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories ").Where("transaction_type = ? AND seller_id = ? AND is_not_claim = ? AND dmo_destination_port_id = ? AND shipping_date >= ? AND shipping_date <= ? AND dmo_id IS NOT NULL AND customer_id IS NULL", "DN", iupopkId, false, value.PortId, shippingDateFrom, shippingDateTo).Scan(&transactionRealization).Error

			if errTrRealization != nil {
				return detailElectricAssignment, errTrRealization
			}

			var groupingRealization Realization

			var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns
where id in (select grouping_vessel_dn_id from transactions where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and transaction_type = 'DN' and is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND dmo_destination_port_id = %v AND bl_date >= '%s' AND bl_date <= '%s' AND iupopk_id = %v
				`, value.PortId, shippingDateFrom, shippingDateTo, iupopkId)

			errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

			if errGroupingRealization != nil {
				return detailElectricAssignment, errGroupingRealization
			}

			transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
			transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

			errFindTemp := r.db.Where("port_id = ? AND electric_assignment_id = ?", value.PortId, electricAssignment.ID).Find(&tempAssignment).Error

			if errFindTemp != nil {
				return detailElectricAssignment, errFindTemp
			}

			var quantity = transactionRealization.RealizationQuantity
			for _, val := range listElectricAssignment {
				if val.PortId == value.PortId && val.LetterNumber == electricAssignment.LetterNumber {
					if quantity-val.Quantity > 0 {
						quantity = quantity - val.Quantity
					} else {
						quantity = 0
					}
				}
			}

			if quantity > value.Quantity {
				var isOverQuantity = false

				for _, temp := range tempAssignment {
					if temp.LetterNumber == electricAssignment.LetterNumber3 || temp.LetterNumber == electricAssignment.LetterNumber4 {
						isOverQuantity = true
					}
				}

				if isOverQuantity {
					realization.RealizationQuantity = quantity
				} else {
					realization.RealizationQuantity = value.Quantity
				}
			} else {
				realization.RealizationQuantity = quantity
			}
			realization.RealizationAverageCalories = transactionRealization.RealizationAverageCalories

			listRealizationTemp.ListRealizationEndUser = append(listRealizationTemp.ListRealizationEndUser, realization)
		}

		listRealization = append(listRealization, listRealizationTemp)
	}

	if electricAssignment.LetterNumber3 != "" {
		var listAssignment []ElectricAssignmentEndUser
		var listRealizationTemp ListRealization

		listRealizationTemp.Order = 3
		listRealizationTemp.LetterNumber = electricAssignment.LetterNumber3

		errFindList := r.db.Preload(clause.Associations).Preload("Port.PortLocation").Where("electric_assignment_id = ? AND letter_number = ?", id, electricAssignment.LetterNumber3).Find(&listAssignment).Error

		if errFindList != nil {
			return detailElectricAssignment, errFindList
		}

		for _, value := range listAssignment {
			var realization RealizationEndUser
			realization.PortId = value.PortId
			realization.Port = value.Port
			realization.AverageCalories = value.AverageCalories
			realization.Quantity = value.Quantity
			realization.EndUser = value.EndUser
			realization.ID = value.ID
			realization.LetterNumber = value.LetterNumber

			var transactionRealization Realization
			var tempAssignment []ElectricAssignmentEndUser

			errTrRealization := r.db.Table("transactions").Select("SUM(quantity_unloading) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories ").Where("transaction_type = ? AND seller_id = ? AND is_not_claim = ? AND dmo_destination_port_id = ? AND shipping_date >= ? AND shipping_date <= ? AND dmo_id IS NOT NULL AND customer_id IS NULL", "DN", iupopkId, false, value.PortId, shippingDateFrom, shippingDateTo).Scan(&transactionRealization).Error

			if errTrRealization != nil {
				return detailElectricAssignment, errTrRealization
			}

			var groupingRealization Realization

			var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns
where id in (select grouping_vessel_dn_id from transactions where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and transaction_type = 'DN' and is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND dmo_destination_port_id = %v AND bl_date >= '%s' AND bl_date <= '%s' AND iupopk_id = %v
				`, value.PortId, shippingDateFrom, shippingDateTo, iupopkId)

			errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

			if errGroupingRealization != nil {
				return detailElectricAssignment, errGroupingRealization
			}

			transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
			transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

			errFindTemp := r.db.Where("port_id = ? AND electric_assignment_id = ?", value.PortId, electricAssignment.ID).Find(&tempAssignment).Error

			if errFindTemp != nil {
				return detailElectricAssignment, errFindTemp
			}

			var quantity = transactionRealization.RealizationQuantity
			for _, val := range listElectricAssignment {
				if val.PortId == value.PortId && (val.LetterNumber == electricAssignment.LetterNumber || val.LetterNumber == electricAssignment.LetterNumber2) {
					if quantity-val.Quantity > 0 {
						quantity = quantity - val.Quantity
					} else {
						quantity = 0
					}
				}
			}

			if quantity > value.Quantity {

				var isOverQuantity = false

				for _, temp := range tempAssignment {
					if temp.LetterNumber == electricAssignment.LetterNumber4 {
						isOverQuantity = true
					}
				}

				if isOverQuantity {
					realization.RealizationQuantity = quantity
				} else {
					realization.RealizationQuantity = value.Quantity
				}
			} else {
				realization.RealizationQuantity = quantity
			}
			realization.RealizationAverageCalories = transactionRealization.RealizationAverageCalories

			listRealizationTemp.ListRealizationEndUser = append(listRealizationTemp.ListRealizationEndUser, realization)
		}

		listRealization = append(listRealization, listRealizationTemp)
	}

	if electricAssignment.LetterNumber4 != "" {
		var listAssignment []ElectricAssignmentEndUser
		var listRealizationTemp ListRealization

		listRealizationTemp.Order = 4
		listRealizationTemp.LetterNumber = electricAssignment.LetterNumber4

		errFindList := r.db.Preload(clause.Associations).Preload("Port.PortLocation").Where("electric_assignment_id = ? AND letter_number = ?", id, electricAssignment.LetterNumber4).Find(&listAssignment).Error

		if errFindList != nil {
			return detailElectricAssignment, errFindList
		}

		for _, value := range listAssignment {
			var realization RealizationEndUser
			realization.PortId = value.PortId
			realization.Port = value.Port
			realization.AverageCalories = value.AverageCalories
			realization.Quantity = value.Quantity
			realization.EndUser = value.EndUser
			realization.ID = value.ID
			realization.LetterNumber = value.LetterNumber

			var transactionRealization Realization
			errTrRealization := r.db.Table("transactions").Select("SUM(quantity_unloading) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories ").Where("transaction_type = ? AND seller_id = ? AND is_not_claim = ? AND dmo_destination_port_id = ? AND shipping_date >= ? AND shipping_date <= ? AND dmo_id IS NOT NULL AND customer_id IS NULL", "DN", iupopkId, false, value.PortId, shippingDateFrom, shippingDateTo).Scan(&transactionRealization).Error

			if errTrRealization != nil {
				return detailElectricAssignment, errTrRealization
			}

			var groupingRealization Realization

			var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns
where id in (select grouping_vessel_dn_id from transactions where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and transaction_type = 'DN' and is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND dmo_destination_port_id = %v AND bl_date >= '%s' AND bl_date <= '%s' AND iupopk_id = %v
				`, value.PortId, shippingDateFrom, shippingDateTo, iupopkId)

			errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

			if errGroupingRealization != nil {
				return detailElectricAssignment, errGroupingRealization
			}

			transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
			transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

			var quantity = transactionRealization.RealizationQuantity
			for _, val := range listElectricAssignment {
				if val.PortId == value.PortId && (val.LetterNumber == electricAssignment.LetterNumber || val.LetterNumber == electricAssignment.LetterNumber2 || val.LetterNumber == electricAssignment.LetterNumber3) {
					if quantity-val.Quantity > 0 {
						quantity = quantity - val.Quantity
					} else {
						quantity = 0
					}
				}
			}

			realization.RealizationQuantity = quantity
			realization.RealizationAverageCalories = transactionRealization.RealizationAverageCalories

			listRealizationTemp.ListRealizationEndUser = append(listRealizationTemp.ListRealizationEndUser, realization)
		}

		listRealization = append(listRealization, listRealizationTemp)
	}

	var realizationSupplier []RealizationSupplier

	for _, v := range listElectricAssignment {
		var realizationSupplierTemp []RealizationSupplier

		var groupingRealizationSupplierTemp []RealizationSupplier

		var rawQuery = fmt.Sprintf(`select SUM(quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories, t.customer_id as supplier_id, t.dmo_destination_port_id as port_id from transactions t
                                LEFT JOIN companies c on c.id = t.customer_id
																where  t.shipping_date >= '%s' AND t.shipping_date <= '%s' AND t.seller_id = %v and t.dmo_destination_port_id = %v and t.dmo_id IS NOT NULL and t.grouping_vessel_dn_id IS NULL and t.report_dmo_id IS NOT NULL
																group by t.customer_id , t.dmo_destination_port_id, c.id, p.id
				`, shippingDateFrom, shippingDateTo, iupopkId, v.PortId, v.PortId)

		errRealizationSupplier := r.db.Preload(clause.Associations).Raw(rawQuery).Find(&realizationSupplierTemp).Error

		if errRealizationSupplier != nil {
			return detailElectricAssignment, errRealizationSupplier
		}

		var groupingRawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories, gvd.buyer_id as supplier_id, gvd.dmo_destination_port_id as port_id from grouping_vessel_dns gvd
                              LEFT JOIN companies c on c.id = gvd.buyer_id
																where  gvd.bl_date >= '%s' AND gvd.bl_date <= '%s' AND gvd.iupopk_id = %v and gvd.dmo_destination_port_id = %v and gvd.report_dmo_id IS NOT NULL
																group by gvd.buyer_id , gvd.dmo_destination_port_id
				`, shippingDateFrom, shippingDateTo, iupopkId, v.PortId, v.PortId)

		errGroupingRealizationSupplier := r.db.Preload(clause.Associations).Raw(groupingRawQuery).Find(&groupingRealizationSupplierTemp).Error

		if errGroupingRealizationSupplier != nil {
			return detailElectricAssignment, errGroupingRealizationSupplier
		}

		for _, v := range realizationSupplierTemp {
			realizationSupplier = append(realizationSupplier, v)
		}

		for _, v := range groupingRealizationSupplierTemp {
			realizationSupplier = append(realizationSupplier, v)
		}
	}

	detailElectricAssignment.ListRealizationSupplier = realizationSupplier
	detailElectricAssignment.ListRealization = listRealization

	return detailElectricAssignment, nil
}
