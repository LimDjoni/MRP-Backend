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

			realization.SupplierId = value.SupplierId
			realization.Supplier = value.Supplier
			realization.AverageCalories = value.AverageCalories
			realization.Quantity = value.Quantity
			realization.EndUser = value.EndUser
			realization.ID = value.ID
			realization.LetterNumber = value.LetterNumber
			var tempAssignment []ElectricAssignmentEndUser
			var transactionRealization Realization

			if value.SupplierId != nil {
				errTrRealization := r.db.Table("transactions").Select("SUM(transactions.quantity_unloading) as realization_quantity, AVG(transactions.quality_calories_ar) as realization_average_calories ").Joins("LEFT JOIN companies companies on companies.id = transactions.customer_id").Where("transactions.transaction_type = ? AND transactions.seller_id = ? AND transactions.is_not_claim = ? AND transactions.dmo_destination_port_id = ? AND transactions.shipping_date >= ? AND transactions.shipping_date <= ? AND transactions.dmo_id IS NOT NULL AND companies.company_name = ? AND grouping_vessel_dn_id IS NULL", "DN", iupopkId, false, value.PortId, shippingDateFrom, shippingDateTo, value.Supplier.CompanyName).Scan(&transactionRealization).Error

				if errTrRealization != nil {
					return detailElectricAssignment, errTrRealization
				}

				var groupingRealization Realization

				var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns
where id in (select grouping_vessel_dn_id from transactions t LEFT JOIN companies c on c.id = t.customer_id where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and c.company_name = '%s' and t.transaction_type = 'DN' and t.is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND dmo_destination_port_id = %v AND bl_date >= '%s' AND bl_date <= '%s' AND iupopk_id = %v
				`, value.Supplier.CompanyName, value.PortId, shippingDateFrom, shippingDateTo, iupopkId)

				errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

				if errGroupingRealization != nil {
					return detailElectricAssignment, errGroupingRealization
				}

				transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
				transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

				errFindTemp := r.db.Where("port_id = ? AND supplier_id = ? AND electric_assignment_id = ?", value.PortId, value.SupplierId, electricAssignment.ID).Find(&tempAssignment).Error

				if errFindTemp != nil {
					return detailElectricAssignment, errFindTemp
				}
			} else {
				errTrRealization := r.db.Table("transactions").Select("SUM(quantity_unloading) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories").Where("transaction_type = ? AND seller_id = ? AND is_not_claim = ? AND dmo_destination_port_id = ? AND shipping_date >= ? AND shipping_date <= ? AND dmo_id IS NOT NULL AND customer_id IS NULL AND grouping_vessel_dn_id IS NULL", "DN", iupopkId, false, value.PortId, shippingDateFrom, shippingDateTo).Scan(&transactionRealization).Error

				if errTrRealization != nil {
					return detailElectricAssignment, errTrRealization
				}

				var groupingRealization Realization

				var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns
where id in (select grouping_vessel_dn_id from transactions t LEFT JOIN companies c on c.id = t.customer_id where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and t.customer_id IS NULL and t.transaction_type = 'DN' and t.is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND dmo_destination_port_id = %v AND bl_date >= '%s' AND bl_date <= '%s' AND iupopk_id = %v
				`, value.PortId, shippingDateFrom, shippingDateTo, iupopkId)

				errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

				if errGroupingRealization != nil {
					return detailElectricAssignment, errGroupingRealization
				}

				transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
				transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

				errFindTemp := r.db.Where("port_id = ? AND supplier_id IS NULL AND electric_assignment_id = ?", value.PortId, electricAssignment.ID).Find(&tempAssignment).Error

				if errFindTemp != nil {
					return detailElectricAssignment, errFindTemp
				}
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
			realization.SupplierId = value.SupplierId
			realization.Supplier = value.Supplier
			realization.AverageCalories = value.AverageCalories
			realization.Quantity = value.Quantity
			realization.EndUser = value.EndUser
			realization.ID = value.ID
			realization.LetterNumber = value.LetterNumber

			var transactionRealization Realization
			var tempAssignment []ElectricAssignmentEndUser
			if value.SupplierId != nil {
				errTrRealization := r.db.Table("transactions").Select("SUM(transactions.quantity_unloading) as realization_quantity, AVG(transactions.quality_calories_ar) as realization_average_calories ").Joins("LEFT JOIN companies companies on companies.id = transactions.customer_id").Where("transactions.transaction_type = ? AND transactions.seller_id = ? AND transactions.is_not_claim = ? AND transactions.dmo_destination_port_id = ? AND transactions.shipping_date >= ? AND transactions.shipping_date <= ? AND transactions.dmo_id IS NOT NULL AND companies.company_name = ?", "DN", iupopkId, false, value.PortId, shippingDateFrom, shippingDateTo, value.Supplier.CompanyName).Scan(&transactionRealization).Error

				if errTrRealization != nil {
					return detailElectricAssignment, errTrRealization
				}

				var groupingRealization Realization

				var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns
where id in (select grouping_vessel_dn_id from transactions t LEFT JOIN companies c on c.id = t.customer_id where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and c.company_name = '%s' and t.transaction_type = 'DN' and t.is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND dmo_destination_port_id = %v AND bl_date >= '%s' AND bl_date <= '%s' AND iupopk_id = %v
				`, value.Supplier.CompanyName, value.PortId, shippingDateFrom, shippingDateTo, iupopkId)

				errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

				if errGroupingRealization != nil {
					return detailElectricAssignment, errGroupingRealization
				}

				transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
				transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

				errFindTemp := r.db.Where("port_id = ? AND supplier_id = ? AND electric_assignment_id = ?", value.PortId, value.SupplierId, electricAssignment.ID).Find(&tempAssignment).Error

				if errFindTemp != nil {
					return detailElectricAssignment, errFindTemp
				}
			} else {
				errTrRealization := r.db.Table("transactions").Select("SUM(quantity_unloading) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories ").Where("transaction_type = ? AND seller_id = ? AND is_not_claim = ? AND dmo_destination_port_id = ? AND shipping_date >= ? AND shipping_date <= ? AND dmo_id IS NOT NULL AND customer_id IS NULL", "DN", iupopkId, false, value.PortId, shippingDateFrom, shippingDateTo).Scan(&transactionRealization).Error

				if errTrRealization != nil {
					return detailElectricAssignment, errTrRealization
				}

				var groupingRealization Realization

				var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns
where id in (select grouping_vessel_dn_id from transactions t LEFT JOIN companies c on c.id = t.customer_id where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and t.customer_id IS NULL and t.transaction_type = 'DN' and t.is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND dmo_destination_port_id = %v AND bl_date >= '%s' AND bl_date <= '%s' AND iupopk_id = %v
				`, value.PortId, shippingDateFrom, shippingDateTo, iupopkId)

				errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

				if errGroupingRealization != nil {
					return detailElectricAssignment, errGroupingRealization
				}

				transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
				transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

				errFindTemp := r.db.Where("port_id = ? AND supplier_id IS NULL AND electric_assignment_id = ?", value.PortId, electricAssignment.ID).Find(&tempAssignment).Error

				if errFindTemp != nil {
					return detailElectricAssignment, errFindTemp
				}
			}

			var quantity = transactionRealization.RealizationQuantity
			for _, val := range listElectricAssignment {
				if value.SupplierId != nil {
					if val.SupplierId != nil {
						if val.PortId == value.PortId && val.Supplier.CompanyName == value.Supplier.CompanyName && val.LetterNumber == electricAssignment.LetterNumber {
							if quantity-val.Quantity > 0 {
								quantity = quantity - val.Quantity
							} else {
								quantity = 0
							}
						}
					}
				} else {
					if val.PortId == value.PortId && val.SupplierId == nil && value.SupplierId == nil && val.LetterNumber == electricAssignment.LetterNumber {
						if quantity-val.Quantity > 0 {
							quantity = quantity - val.Quantity
						} else {
							quantity = 0
						}
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
			realization.SupplierId = value.SupplierId
			realization.Supplier = value.Supplier
			realization.AverageCalories = value.AverageCalories
			realization.Quantity = value.Quantity
			realization.EndUser = value.EndUser
			realization.ID = value.ID
			realization.LetterNumber = value.LetterNumber

			var transactionRealization Realization
			var tempAssignment []ElectricAssignmentEndUser

			if value.SupplierId != nil {
				errTrRealization := r.db.Table("transactions").Select("SUM(transactions.quantity_unloading) as realization_quantity, AVG(transactions.quality_calories_ar) as realization_average_calories ").Joins("LEFT JOIN companies companies on companies.id = transactions.customer_id").Where("transactions.transaction_type = ? AND transactions.seller_id = ? AND transactions.is_not_claim = ? AND transactions.dmo_destination_port_id = ? AND transactions.shipping_date >= ? AND transactions.shipping_date <= ? AND transactions.dmo_id IS NOT NULL AND companies.company_name = ?", "DN", iupopkId, false, value.PortId, shippingDateFrom, shippingDateTo, value.Supplier.CompanyName).Scan(&transactionRealization).Error

				if errTrRealization != nil {
					return detailElectricAssignment, errTrRealization
				}

				var groupingRealization Realization

				var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns
where id in (select grouping_vessel_dn_id from transactions t LEFT JOIN companies c on c.id = t.customer_id where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and c.company_name = '%s' and t.transaction_type = 'DN' and t.is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND dmo_destination_port_id = %v AND bl_date >= '%s' AND bl_date <= '%s' AND iupopk_id = %v
				`, value.Supplier.CompanyName, value.PortId, shippingDateFrom, shippingDateTo, iupopkId)

				errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

				if errGroupingRealization != nil {
					return detailElectricAssignment, errGroupingRealization
				}

				transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
				transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

				errFindTemp := r.db.Where("port_id = ? AND supplier_id = ? AND electric_assignment_id = ?", value.PortId, value.SupplierId, electricAssignment.ID).Find(&tempAssignment).Error

				if errFindTemp != nil {
					return detailElectricAssignment, errFindTemp
				}

			} else {
				errTrRealization := r.db.Table("transactions").Select("SUM(quantity_unloading) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories ").Where("transaction_type = ? AND seller_id = ? AND is_not_claim = ? AND dmo_destination_port_id = ? AND shipping_date >= ? AND shipping_date <= ? AND dmo_id IS NOT NULL AND customer_id IS NULL", "DN", iupopkId, false, value.PortId, shippingDateFrom, shippingDateTo).Scan(&transactionRealization).Error

				if errTrRealization != nil {
					return detailElectricAssignment, errTrRealization
				}

				var groupingRealization Realization

				var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns
where id in (select grouping_vessel_dn_id from transactions t LEFT JOIN companies c on c.id = t.customer_id where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and t.customer_id IS NULL and t.transaction_type = 'DN' and t.is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND dmo_destination_port_id = %v AND bl_date >= '%s' AND bl_date <= '%s' AND iupopk_id = %v
				`, value.PortId, shippingDateFrom, shippingDateTo, iupopkId)

				errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

				if errGroupingRealization != nil {
					return detailElectricAssignment, errGroupingRealization
				}

				transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
				transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

				errFindTemp := r.db.Where("port_id = ? AND supplier_id IS NULL AND electric_assignment_id = ?", value.PortId, electricAssignment.ID).Find(&tempAssignment).Error

				if errFindTemp != nil {
					return detailElectricAssignment, errFindTemp
				}
			}

			var quantity = transactionRealization.RealizationQuantity
			for _, val := range listElectricAssignment {
				if value.SupplierId != nil {
					if val.SupplierId != nil {
						if val.PortId == value.PortId && val.Supplier.CompanyName == value.Supplier.CompanyName && (val.LetterNumber == electricAssignment.LetterNumber || val.LetterNumber == electricAssignment.LetterNumber2) {
							if quantity-val.Quantity > 0 {
								quantity = quantity - val.Quantity
							} else {
								quantity = 0
							}
						}
					}
				} else {
					if val.PortId == value.PortId && val.SupplierId == nil && value.SupplierId == nil && (val.LetterNumber == electricAssignment.LetterNumber || val.LetterNumber == electricAssignment.LetterNumber2) {
						if quantity-val.Quantity > 0 {
							quantity = quantity - val.Quantity
						} else {
							quantity = 0
						}
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
			realization.SupplierId = value.SupplierId
			realization.Supplier = value.Supplier
			realization.AverageCalories = value.AverageCalories
			realization.Quantity = value.Quantity
			realization.EndUser = value.EndUser
			realization.ID = value.ID
			realization.LetterNumber = value.LetterNumber

			var transactionRealization Realization
			if value.SupplierId != nil {
				errTrRealization := r.db.Table("transactions").Select("SUM(transactions.quantity_unloading) as realization_quantity, AVG(transactions.quality_calories_ar) as realization_average_calories ").Joins("LEFT JOIN companies companies on companies.id = transactions.customer_id").Where("transactions.transaction_type = ? AND transactions.seller_id = ? AND transactions.is_not_claim = ? AND transactions.dmo_destination_port_id = ? AND transactions.shipping_date >= ? AND transactions.shipping_date <= ? AND transactions.dmo_id IS NOT NULL AND companies.company_name = ?", "DN", iupopkId, false, value.PortId, shippingDateFrom, shippingDateTo, value.Supplier.CompanyName).Scan(&transactionRealization).Error

				if errTrRealization != nil {
					return detailElectricAssignment, errTrRealization
				}

				var groupingRealization Realization

				var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns
where id in (select grouping_vessel_dn_id from transactions t LEFT JOIN companies c on c.id = t.customer_id where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and c.company_name = '%s' and t.transaction_type = 'DN' and t.is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND dmo_destination_port_id = %v AND bl_date >= '%s' AND bl_date <= '%s' AND iupopk_id = %v
				`, value.Supplier.CompanyName, value.PortId, shippingDateFrom, shippingDateTo, iupopkId)

				errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

				if errGroupingRealization != nil {
					return detailElectricAssignment, errGroupingRealization
				}

				transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
				transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

			} else {
				errTrRealization := r.db.Table("transactions").Select("SUM(quantity_unloading) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories ").Where("transaction_type = ? AND seller_id = ? AND is_not_claim = ? AND dmo_destination_port_id = ? AND shipping_date >= ? AND shipping_date <= ? AND dmo_id IS NOT NULL AND customer_id IS NULL", "DN", iupopkId, false, value.PortId, shippingDateFrom, shippingDateTo).Scan(&transactionRealization).Error

				if errTrRealization != nil {
					return detailElectricAssignment, errTrRealization
				}

				var groupingRealization Realization

				var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns
where id in (select grouping_vessel_dn_id from transactions t LEFT JOIN companies c on c.id = t.customer_id where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and t.customer_id IS NULL and t.transaction_type = 'DN' and t.is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND dmo_destination_port_id = %v AND bl_date >= '%s' AND bl_date <= '%s' AND iupopk_id = %v
				`, value.PortId, shippingDateFrom, shippingDateTo, iupopkId)

				errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

				if errGroupingRealization != nil {
					return detailElectricAssignment, errGroupingRealization
				}

				transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
				transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2
			}

			var quantity = transactionRealization.RealizationQuantity
			for _, val := range listElectricAssignment {
				if value.SupplierId != nil {
					if val.SupplierId != nil {
						if val.PortId == value.PortId && val.Supplier.CompanyName == value.Supplier.CompanyName && (val.LetterNumber == electricAssignment.LetterNumber || val.LetterNumber == electricAssignment.LetterNumber2 || val.LetterNumber == electricAssignment.LetterNumber3) {
							if quantity-val.Quantity > 0 {
								quantity = quantity - val.Quantity
							} else {
								quantity = 0
							}
						}
					}
				} else {
					if val.PortId == value.PortId && val.SupplierId == nil && value.SupplierId == nil && (val.LetterNumber == electricAssignment.LetterNumber || val.LetterNumber == electricAssignment.LetterNumber2 || val.LetterNumber == electricAssignment.LetterNumber3) {
						if quantity-val.Quantity > 0 {
							quantity = quantity - val.Quantity
						} else {
							quantity = 0
						}
					}
				}
			}

			realization.RealizationQuantity = quantity
			realization.RealizationAverageCalories = transactionRealization.RealizationAverageCalories

			listRealizationTemp.ListRealizationEndUser = append(listRealizationTemp.ListRealizationEndUser, realization)
		}

		listRealization = append(listRealization, listRealizationTemp)
	}

	detailElectricAssignment.ListRealization = listRealization

	return detailElectricAssignment, nil
}
