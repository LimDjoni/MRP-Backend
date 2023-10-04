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

	shippingDateFrom := fmt.Sprintf("%s-01-01", cafAssignment.Year)
	shippingDateTo := fmt.Sprintf("%s-12-31", cafAssignment.Year)

	var listRealization []ListRealization

	if cafAssignment.LetterNumber != "" {
		var listAssignment []CafAssignmentEndUser
		var listRealizationTemp ListRealization

		listRealizationTemp.Order = 1
		listRealizationTemp.LetterNumber = cafAssignment.LetterNumber

		errFindList := r.db.Preload(clause.Associations).Where("caf_assignment_id = ? AND letter_number = ?", id, cafAssignment.LetterNumber).Find(&listAssignment).Error

		if errFindList != nil {
			return detailCafAssignment, errFindList
		}
		for _, value := range listAssignment {
			var realization RealizationEndUser
			realization.AverageCalories = value.AverageCalories
			realization.Quantity = value.Quantity
			realization.EndUserId = value.EndUserId
			realization.EndUser = value.EndUser
			realization.EndUserString = value.EndUserString
			realization.ID = value.ID
			realization.LetterNumber = value.LetterNumber
			var transactionRealization Realization
			var tempAssignment []CafAssignmentEndUser

			query := fmt.Sprintf("transactions.transaction_type = '%s' AND transactions.seller_id = %v AND transactions.is_not_claim = false AND transactions.shipping_date >= '%s' AND transactions.shipping_date <= '%s' AND company.company_name = '%s' AND transactions.dmo_id IS NOT NULL AND transactions.report_dmo_id IS NOT NULL AND transactions.grouping_vessel_dn_id IS NULL", "DN", iupopkId, shippingDateFrom, shippingDateTo, value.EndUserString)

			errTrRealization := r.db.Table("transactions").Select("SUM(transactions.quantity) as realization_quantity, AVG(transactions.quality_calories_ar) as realization_average_calories ").Joins("left join companies company on company.id = transactions.dmo_buyer_id").Where(query).Scan(&transactionRealization).Error

			if errTrRealization != nil {
				return detailCafAssignment, errTrRealization
			}

			var groupingRealization Realization

			var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns gvd
				LEFT JOIN companies c on c.id = gvd.buyer_id
where gvd.id in (select grouping_vessel_dn_id from transactions where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and transaction_type = 'DN' and is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND gvd.bl_date >= '%s' AND gvd.bl_date <= '%s' AND gvd.iupopk_id = %v and c.company_name = '%s' AND gvd.report_dmo_id IS NOT NULL
				`, shippingDateFrom, shippingDateTo, iupopkId, value.EndUserString)

			errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

			if errGroupingRealization != nil {
				return detailCafAssignment, errGroupingRealization
			}

			transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
			transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

			errFindTemp := r.db.Where("end_user_string = ? AND caf_assignment_id = ?", value.EndUserString, cafAssignment.ID).Find(&tempAssignment).Error

			if errFindTemp != nil {
				return detailCafAssignment, errFindTemp
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

	if cafAssignment.LetterNumber2 != "" {
		var listAssignment []CafAssignmentEndUser
		var listRealizationTemp ListRealization

		listRealizationTemp.Order = 2
		listRealizationTemp.LetterNumber = cafAssignment.LetterNumber2

		errFindList := r.db.Preload(clause.Associations).Where("caf_assignment_id = ? AND letter_number = ?", id, cafAssignment.LetterNumber2).Find(&listAssignment).Error

		if errFindList != nil {
			return detailCafAssignment, errFindList
		}
		for _, value := range listAssignment {
			var realization RealizationEndUser
			realization.AverageCalories = value.AverageCalories
			realization.Quantity = value.Quantity
			realization.EndUserId = value.EndUserId
			realization.EndUser = value.EndUser
			realization.EndUserString = value.EndUserString
			realization.ID = value.ID
			realization.LetterNumber = value.LetterNumber
			var transactionRealization Realization
			var tempAssignment []CafAssignmentEndUser

			query := fmt.Sprintf("transactions.transaction_type = '%s' AND transactions.seller_id = %v AND transactions.is_not_claim = false AND transactions.shipping_date >= '%s' AND transactions.shipping_date <= '%s' AND company.company_name = '%s' AND transactions.dmo_id IS NOT NULL AND transactions.report_dmo_id IS NOT NULL", "DN", iupopkId, shippingDateFrom, shippingDateTo, value.EndUserString)

			errTrRealization := r.db.Table("transactions").Select("SUM(transactions.quantity) as realization_quantity, AVG(transactions.quality_calories_ar) as realization_average_calories ").Joins("left join companies company on company.id = transactions.dmo_buyer_id").Where(query).Scan(&transactionRealization).Error

			if errTrRealization != nil {
				return detailCafAssignment, errTrRealization
			}

			var groupingRealization Realization

			var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns gvd
				LEFT JOIN companies c on c.id = gvd.buyer_id
where gvd.id in (select grouping_vessel_dn_id from transactions where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and transaction_type = 'DN' and is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND gvd.bl_date >= '%s' AND gvd.bl_date <= '%s' AND gvd.iupopk_id = %v and c.company_name = '%s' AND gvd.report_dmo_id IS NOT NULL
				`, shippingDateFrom, shippingDateTo, iupopkId, value.EndUserString)

			errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

			if errGroupingRealization != nil {
				return detailCafAssignment, errGroupingRealization
			}

			transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
			transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

			errFindTemp := r.db.Where("end_user_string = ? AND caf_assignment_id = ?", value.EndUserString, cafAssignment.ID).Find(&tempAssignment).Error

			if errFindTemp != nil {
				return detailCafAssignment, errFindTemp
			}

			var quantity = transactionRealization.RealizationQuantity
			for _, val := range listCafAssignment {
				if val.EndUserString == value.EndUserString && val.LetterNumber == cafAssignment.LetterNumber {
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
					if temp.LetterNumber == cafAssignment.LetterNumber3 || temp.LetterNumber == cafAssignment.LetterNumber4 {
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

	if cafAssignment.LetterNumber3 != "" {
		var listAssignment []CafAssignmentEndUser
		var listRealizationTemp ListRealization

		listRealizationTemp.Order = 3
		listRealizationTemp.LetterNumber = cafAssignment.LetterNumber3

		errFindList := r.db.Preload(clause.Associations).Where("caf_assignment_id = ? AND letter_number = ?", id, cafAssignment.LetterNumber3).Find(&listAssignment).Error

		if errFindList != nil {
			return detailCafAssignment, errFindList
		}
		for _, value := range listAssignment {
			var realization RealizationEndUser
			realization.AverageCalories = value.AverageCalories
			realization.Quantity = value.Quantity
			realization.EndUserId = value.EndUserId
			realization.EndUser = value.EndUser
			realization.EndUserString = value.EndUserString
			realization.ID = value.ID
			realization.LetterNumber = value.LetterNumber
			var transactionRealization Realization
			var tempAssignment []CafAssignmentEndUser

			query := fmt.Sprintf("transactions.transaction_type = '%s' AND transactions.seller_id = %v AND transactions.is_not_claim = false AND transactions.shipping_date >= '%s' AND transactions.shipping_date <= '%s' AND company.company_name = '%s' AND transactions.dmo_id IS NOT NULL AND transactions.report_dmo_id IS NOT NULL", "DN", iupopkId, shippingDateFrom, shippingDateTo, value.EndUserString)

			errTrRealization := r.db.Table("transactions").Select("SUM(transactions.quantity) as realization_quantity, AVG(transactions.quality_calories_ar) as realization_average_calories ").Joins("left join companies company on company.id = transactions.dmo_buyer_id").Where(query).Scan(&transactionRealization).Error

			if errTrRealization != nil {
				return detailCafAssignment, errTrRealization
			}

			var groupingRealization Realization

			var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns gvd
				LEFT JOIN companies c on c.id = gvd.buyer_id
where gvd.id in (select grouping_vessel_dn_id from transactions where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and transaction_type = 'DN' and is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND gvd.bl_date >= '%s' AND gvd.bl_date <= '%s' AND gvd.iupopk_id = %v and c.company_name = '%s' AND gvd.report_dmo_id IS NOT NULL
				`, shippingDateFrom, shippingDateTo, iupopkId, value.EndUserString)

			errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

			if errGroupingRealization != nil {
				return detailCafAssignment, errGroupingRealization
			}

			transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
			transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

			errFindTemp := r.db.Where("end_user_string = ? AND caf_assignment_id = ?", value.EndUserString, cafAssignment.ID).Find(&tempAssignment).Error

			if errFindTemp != nil {
				return detailCafAssignment, errFindTemp
			}

			var quantity = transactionRealization.RealizationQuantity
			for _, val := range listCafAssignment {
				if val.EndUserString == value.EndUserString && (val.LetterNumber == cafAssignment.LetterNumber || val.LetterNumber == cafAssignment.LetterNumber2) {
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
					if temp.LetterNumber == cafAssignment.LetterNumber4 {
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

	if cafAssignment.LetterNumber4 != "" {
		var listAssignment []CafAssignmentEndUser
		var listRealizationTemp ListRealization

		listRealizationTemp.Order = 4
		listRealizationTemp.LetterNumber = cafAssignment.LetterNumber4

		errFindList := r.db.Preload(clause.Associations).Where("caf_assignment_id = ? AND letter_number = ?", id, cafAssignment.LetterNumber4).Find(&listAssignment).Error

		if errFindList != nil {
			return detailCafAssignment, errFindList
		}
		for _, value := range listAssignment {
			var realization RealizationEndUser
			realization.AverageCalories = value.AverageCalories
			realization.Quantity = value.Quantity
			realization.EndUserId = value.EndUserId
			realization.EndUser = value.EndUser
			realization.EndUserString = value.EndUserString
			realization.ID = value.ID
			realization.LetterNumber = value.LetterNumber
			var transactionRealization Realization

			query := fmt.Sprintf("transactions.transaction_type = '%s' AND transactions.seller_id = %v AND transactions.is_not_claim = false AND transactions.shipping_date >= '%s' AND transactions.shipping_date <= '%s' AND company.company_name = '%s' AND transactions.dmo_id IS NOT NULL AND transactions.report_dmo_id IS NOT NULL", "DN", iupopkId, shippingDateFrom, shippingDateTo, value.EndUserString)

			errTrRealization := r.db.Table("transactions").Select("SUM(transactions.quantity) as realization_quantity, AVG(transactions.quality_calories_ar) as realization_average_calories ").Joins("left join companies company on company.id = transactions.dmo_buyer_id").Where(query).Scan(&transactionRealization).Error

			if errTrRealization != nil {
				return detailCafAssignment, errTrRealization
			}

			var groupingRealization Realization

			var rawQuery = fmt.Sprintf(`select SUM(grand_total_quantity) as realization_quantity, AVG(quality_calories_ar) as realization_average_calories from grouping_vessel_dns gvd
				LEFT JOIN companies c on c.id = gvd.buyer_id
where gvd.id in (select grouping_vessel_dn_id from transactions where dmo_id IS NOT NULL and grouping_vessel_dn_id IS NOT NULL and transaction_type = 'DN' and is_not_claim = false
GROUP BY grouping_vessel_dn_id) AND gvd.bl_date >= '%s' AND gvd.bl_date <= '%s' AND gvd.iupopk_id = %v and c.company_name = '%s' AND gvd.report_dmo_id IS NOT NULL
				`, shippingDateFrom, shippingDateTo, iupopkId, value.EndUserString)

			errGroupingRealization := r.db.Raw(rawQuery).Scan(&groupingRealization).Error

			if errGroupingRealization != nil {
				return detailCafAssignment, errGroupingRealization
			}

			transactionRealization.RealizationQuantity += groupingRealization.RealizationQuantity
			transactionRealization.RealizationAverageCalories = (transactionRealization.RealizationAverageCalories + groupingRealization.RealizationAverageCalories) / 2

			var quantity = transactionRealization.RealizationQuantity
			for _, val := range listCafAssignment {
				if val.EndUserString == value.EndUserString && (val.LetterNumber == cafAssignment.LetterNumber || val.LetterNumber == cafAssignment.LetterNumber2 || val.LetterNumber == cafAssignment.LetterNumber3) {
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

	detailCafAssignment.ListRealization = listRealization
	return detailCafAssignment, nil
}
