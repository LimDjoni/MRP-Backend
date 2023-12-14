package production

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetListProduction(page int, filter FilterListProduction, iupopkId int) (Pagination, error)
	DetailProduction(id int, iupopkId int) (Production, error)
	SummaryProduction(year string, iupopkId int) (OutputSummaryProduction, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func helperString(listString []string, dataString string) bool {
	for _, v := range listString {
		if v == dataString {
			return true
		}
	}
	return false
}

func (r *repository) GetListProduction(page int, filter FilterListProduction, iupopkId int) (Pagination, error) {
	var listProduction []Production

	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)
	querySort := "id desc"

	if filter.Field != "" && filter.Sort != "" {
		querySort = filter.Field + " " + filter.Sort
	}

	if filter.ProductionDateStart != "" {
		queryFilter = queryFilter + " AND production_date >= '" + filter.ProductionDateStart + "'"
	}

	if filter.ProductionDateEnd != "" {
		queryFilter = queryFilter + " AND production_date <= '" + filter.ProductionDateEnd + "T23:59:59'"
	}

	if filter.Quantity != "" {
		queryFilter = queryFilter + " AND cast(quantity AS TEXT) LIKE '%" + filter.Quantity + "%'"
	}

	if filter.PitId != "" {
		queryFilter = queryFilter + " AND pit_id = " + filter.PitId
	}

	if filter.JettyId != "" {
		queryFilter = queryFilter + " AND jetty_id = " + filter.JettyId
	}

	errFind := r.db.Preload(clause.Associations).Where(queryFilter).Order(querySort).Scopes(paginateProduction(listProduction, &pagination, r.db, queryFilter)).Find(&listProduction).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listProduction

	return pagination, nil
}

func (r *repository) DetailProduction(id int, iupopkId int) (Production, error) {
	var detailProduction Production

	errFind := r.db.Preload(clause.Associations).Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&detailProduction).Error

	return detailProduction, errFind
}

func (r *repository) SummaryProduction(year string, iupopkId int) (OutputSummaryProduction, error) {
	var summary OutputSummaryProduction

	summary.January = make(map[string]map[string]map[string]float64)
	summary.February = make(map[string]map[string]map[string]float64)
	summary.March = make(map[string]map[string]map[string]float64)
	summary.April = make(map[string]map[string]map[string]float64)
	summary.May = make(map[string]map[string]map[string]float64)
	summary.June = make(map[string]map[string]map[string]float64)
	summary.July = make(map[string]map[string]map[string]float64)
	summary.August = make(map[string]map[string]map[string]float64)
	summary.September = make(map[string]map[string]map[string]float64)
	summary.October = make(map[string]map[string]map[string]float64)
	summary.November = make(map[string]map[string]map[string]float64)
	summary.December = make(map[string]map[string]map[string]float64)

	var productionJanuary []GroupProduction
	var productionFebruary []GroupProduction
	var productionMarch []GroupProduction
	var productionApril []GroupProduction
	var productionMay []GroupProduction
	var productionJune []GroupProduction
	var productionJuly []GroupProduction
	var productionAugust []GroupProduction
	var productionSeptember []GroupProduction
	var productionOctober []GroupProduction
	var productionNovember []GroupProduction
	var productionDecember []GroupProduction

	summary.ListJettyPit = make(map[string][]string)

	errJanuary := r.db.Table("productions").Select("pit_id, isp_id, jetty_id, SUM(quantity) as quantity, SUM(ritase_quantity) as ritase_quantity").Preload(clause.Associations).Where("iupopk_id = ? AND production_date >= ? AND production_date < ?", iupopkId, fmt.Sprintf("%v-01-01", year), fmt.Sprintf("%v-02-01", year)).Group("pit_id, isp_id, jetty_id").Find(&productionJanuary).Error

	if errJanuary != nil {
		return summary, errJanuary
	}

	for _, v := range productionJanuary {
		if v.Jetty != nil {

			if _, ok := summary.ListJettyPit[v.Jetty.Name]; !ok {
				summary.ListJettyPit[v.Jetty.Name] = []string{}
			}

			if v.Pit != nil {
				if !helperString(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name) {
					summary.ListJettyPit[v.Jetty.Name] = append(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name)
				}

				if _, ok := summary.January[v.Jetty.Name]; ok {
					if _, ok2 := summary.January[v.Jetty.Name][v.Pit.Name]; ok2 {
						summary.January[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.January[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					} else {
						summary.January[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
						summary.January[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.January[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					}
				} else {
					summary.January[v.Jetty.Name] = make(map[string]map[string]float64)
					summary.January[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
					summary.January[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
					summary.January[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
				}
			}
		}
	}

	errFebruary := r.db.Table("productions").Select("pit_id, isp_id, jetty_id, SUM(quantity) as quantity, SUM(ritase_quantity) as ritase_quantity").Preload(clause.Associations).Where("iupopk_id = ? AND production_date >= ? AND production_date < ?", iupopkId, fmt.Sprintf("%v-02-01", year), fmt.Sprintf("%v-03-01", year)).Group("pit_id, isp_id, jetty_id").Find(&productionFebruary).Error

	if errFebruary != nil {
		return summary, errFebruary
	}

	for _, v := range productionFebruary {
		if v.Jetty != nil {
			if _, ok := summary.ListJettyPit[v.Jetty.Name]; !ok {
				summary.ListJettyPit[v.Jetty.Name] = []string{}
			}

			if v.Pit != nil {
				if !helperString(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name) {
					summary.ListJettyPit[v.Jetty.Name] = append(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name)
				}
				if _, ok := summary.February[v.Jetty.Name]; ok {
					if _, ok2 := summary.February[v.Jetty.Name][v.Pit.Name]; ok2 {
						summary.February[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.February[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					} else {
						summary.February[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
						summary.February[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.February[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					}
				} else {
					summary.February[v.Jetty.Name] = make(map[string]map[string]float64)
					summary.February[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
					summary.February[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
					summary.February[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
				}
			}
		}
	}
	errMarch := r.db.Table("productions").Select("pit_id, isp_id, jetty_id, SUM(quantity) as quantity, SUM(ritase_quantity) as ritase_quantity").Preload(clause.Associations).Where("iupopk_id = ? AND production_date >= ? AND production_date < ?", iupopkId, fmt.Sprintf("%v-03-01", year), fmt.Sprintf("%v-04-01", year)).Group("pit_id, isp_id, jetty_id").Find(&productionMarch).Error

	if errMarch != nil {
		return summary, errMarch
	}

	for _, v := range productionMarch {
		if v.Jetty != nil {
			if _, ok := summary.ListJettyPit[v.Jetty.Name]; !ok {
				summary.ListJettyPit[v.Jetty.Name] = []string{}
			}

			if v.Pit != nil {
				if !helperString(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name) {
					summary.ListJettyPit[v.Jetty.Name] = append(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name)
				}
				if _, ok := summary.March[v.Jetty.Name]; ok {
					if _, ok2 := summary.March[v.Jetty.Name][v.Pit.Name]; ok2 {
						summary.March[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.March[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					} else {
						summary.March[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
						summary.March[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.March[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					}
				} else {
					summary.March[v.Jetty.Name] = make(map[string]map[string]float64)
					summary.March[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
					summary.March[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
					summary.March[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
				}
			}
		}
	}

	errApril := r.db.Table("productions").Select("pit_id, isp_id, jetty_id, SUM(quantity) as quantity, SUM(ritase_quantity) as ritase_quantity").Preload(clause.Associations).Where("iupopk_id = ? AND production_date >= ? AND production_date < ?", iupopkId, fmt.Sprintf("%v-04-01", year), fmt.Sprintf("%v-05-01", year)).Group("pit_id, isp_id, jetty_id").Find(&productionApril).Error

	if errApril != nil {
		return summary, errApril
	}

	for _, v := range productionApril {
		if v.Jetty != nil {
			if _, ok := summary.ListJettyPit[v.Jetty.Name]; !ok {
				summary.ListJettyPit[v.Jetty.Name] = []string{}
			}

			if v.Pit != nil {
				if !helperString(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name) {
					summary.ListJettyPit[v.Jetty.Name] = append(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name)
				}

				if _, ok := summary.April[v.Jetty.Name]; ok {
					if _, ok2 := summary.April[v.Jetty.Name][v.Pit.Name]; ok2 {
						summary.April[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.April[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					} else {
						summary.April[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
						summary.April[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.April[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					}
				} else {
					summary.April[v.Jetty.Name] = make(map[string]map[string]float64)
					summary.April[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
					summary.April[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
					summary.April[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
				}
			}
		}
	}

	errMay := r.db.Table("productions").Select("pit_id, isp_id, jetty_id, SUM(quantity) as quantity, SUM(ritase_quantity) as ritase_quantity").Preload(clause.Associations).Where("iupopk_id = ? AND production_date >= ? AND production_date < ?", iupopkId, fmt.Sprintf("%v-05-01", year), fmt.Sprintf("%v-06-01", year)).Group("pit_id, isp_id, jetty_id").Find(&productionMay).Error

	if errMay != nil {
		return summary, errMay
	}

	for _, v := range productionMay {
		if v.Jetty != nil {
			if _, ok := summary.ListJettyPit[v.Jetty.Name]; !ok {
				summary.ListJettyPit[v.Jetty.Name] = []string{}
			}

			if v.Pit != nil {
				if !helperString(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name) {
					summary.ListJettyPit[v.Jetty.Name] = append(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name)
				}
				if _, ok := summary.May[v.Jetty.Name]; ok {
					if _, ok2 := summary.May[v.Jetty.Name][v.Pit.Name]; ok2 {
						summary.May[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.May[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					} else {
						summary.May[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
						summary.May[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.May[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					}
				} else {
					summary.May[v.Jetty.Name] = make(map[string]map[string]float64)
					summary.May[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
					summary.May[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
					summary.May[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
				}
			}
		}
	}

	errJune := r.db.Table("productions").Select("pit_id, isp_id, jetty_id, SUM(quantity) as quantity, SUM(ritase_quantity) as ritase_quantity").Preload(clause.Associations).Where("iupopk_id = ? AND production_date >= ? AND production_date < ?", iupopkId, fmt.Sprintf("%v-06-01", year), fmt.Sprintf("%v-07-01", year)).Group("pit_id, isp_id, jetty_id").Find(&productionJune).Error

	if errJune != nil {
		return summary, errJune
	}

	for _, v := range productionJune {
		if v.Jetty != nil {
			if _, ok := summary.ListJettyPit[v.Jetty.Name]; !ok {
				summary.ListJettyPit[v.Jetty.Name] = []string{}
			}

			if v.Pit != nil {
				if !helperString(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name) {
					summary.ListJettyPit[v.Jetty.Name] = append(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name)
				}
				if _, ok := summary.June[v.Jetty.Name]; ok {
					if _, ok2 := summary.June[v.Jetty.Name][v.Pit.Name]; ok2 {
						summary.June[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.June[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					} else {
						summary.June[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
						summary.June[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.June[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					}
				} else {
					summary.June[v.Jetty.Name] = make(map[string]map[string]float64)
					summary.June[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
					summary.June[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
					summary.June[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
				}
			}
		}
	}

	errJuly := r.db.Table("productions").Select("pit_id, isp_id, jetty_id, SUM(quantity) as quantity, SUM(ritase_quantity) as ritase_quantity").Preload(clause.Associations).Where("iupopk_id = ? AND production_date >= ? AND production_date < ?", iupopkId, fmt.Sprintf("%v-07-01", year), fmt.Sprintf("%v-08-01", year)).Group("pit_id, isp_id, jetty_id").Find(&productionJuly).Error

	if errJuly != nil {
		return summary, errJuly
	}

	for _, v := range productionJuly {
		if v.Jetty != nil {
			if _, ok := summary.ListJettyPit[v.Jetty.Name]; !ok {
				summary.ListJettyPit[v.Jetty.Name] = []string{}
			}

			if v.Pit != nil {
				if !helperString(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name) {
					summary.ListJettyPit[v.Jetty.Name] = append(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name)
				}
				if _, ok := summary.July[v.Jetty.Name]; ok {
					if _, ok2 := summary.July[v.Jetty.Name][v.Pit.Name]; ok2 {
						summary.July[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.July[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					} else {
						summary.July[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
						summary.July[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.July[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					}
				} else {
					summary.July[v.Jetty.Name] = make(map[string]map[string]float64)
					summary.July[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
					summary.July[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
					summary.July[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
				}
			}
		}
	}

	errAugust := r.db.Table("productions").Select("pit_id, isp_id, jetty_id, SUM(quantity) as quantity, SUM(ritase_quantity) as ritase_quantity").Preload(clause.Associations).Where("iupopk_id = ? AND production_date >= ? AND production_date < ?", iupopkId, fmt.Sprintf("%v-08-01", year), fmt.Sprintf("%v-09-01", year)).Group("pit_id, isp_id, jetty_id").Find(&productionAugust).Error

	if errAugust != nil {
		return summary, errAugust
	}

	for _, v := range productionAugust {
		if v.Jetty != nil {
			if _, ok := summary.ListJettyPit[v.Jetty.Name]; !ok {
				summary.ListJettyPit[v.Jetty.Name] = []string{}
			}

			if v.Pit != nil {
				if !helperString(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name) {
					summary.ListJettyPit[v.Jetty.Name] = append(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name)
				}
				if _, ok := summary.August[v.Jetty.Name]; ok {
					if _, ok2 := summary.August[v.Jetty.Name][v.Pit.Name]; ok2 {
						summary.August[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.August[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					} else {
						summary.August[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
						summary.August[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.August[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					}
				} else {
					summary.August[v.Jetty.Name] = make(map[string]map[string]float64)
					summary.August[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
					summary.August[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
					summary.August[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
				}
			}
		}
	}

	errSeptember := r.db.Table("productions").Select("pit_id, isp_id, jetty_id, SUM(quantity) as quantity, SUM(ritase_quantity) as ritase_quantity").Preload(clause.Associations).Where("iupopk_id = ? AND production_date >= ? AND production_date < ?", iupopkId, fmt.Sprintf("%v-09-01", year), fmt.Sprintf("%v-10-01", year)).Group("pit_id, isp_id, jetty_id").Find(&productionSeptember).Error

	if errSeptember != nil {
		return summary, errSeptember
	}

	for _, v := range productionSeptember {
		if v.Jetty != nil {
			if _, ok := summary.ListJettyPit[v.Jetty.Name]; !ok {
				summary.ListJettyPit[v.Jetty.Name] = []string{}
			}

			if v.Pit != nil {
				if !helperString(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name) {
					summary.ListJettyPit[v.Jetty.Name] = append(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name)
				}
				if _, ok := summary.September[v.Jetty.Name]; ok {
					if _, ok2 := summary.September[v.Jetty.Name][v.Pit.Name]; ok2 {
						summary.September[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.September[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					} else {
						summary.September[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
						summary.September[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.September[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					}
				} else {
					summary.September[v.Jetty.Name] = make(map[string]map[string]float64)
					summary.September[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
					summary.September[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
					summary.September[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
				}
			}
		}
	}

	errOctober := r.db.Table("productions").Select("pit_id, isp_id, jetty_id, SUM(quantity) as quantity, SUM(ritase_quantity) as ritase_quantity").Preload(clause.Associations).Where("iupopk_id = ? AND production_date >= ? AND production_date < ?", iupopkId, fmt.Sprintf("%v-10-01", year), fmt.Sprintf("%v-11-01", year)).Group("pit_id, isp_id, jetty_id").Find(&productionOctober).Error

	if errOctober != nil {
		return summary, errOctober
	}

	for _, v := range productionOctober {
		if v.Jetty != nil {
			if _, ok := summary.ListJettyPit[v.Jetty.Name]; !ok {
				summary.ListJettyPit[v.Jetty.Name] = []string{}
			}

			if v.Pit != nil {
				if !helperString(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name) {
					summary.ListJettyPit[v.Jetty.Name] = append(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name)
				}
				if _, ok := summary.October[v.Jetty.Name]; ok {
					if _, ok2 := summary.October[v.Jetty.Name][v.Pit.Name]; ok2 {
						summary.October[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.October[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					} else {
						summary.October[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
						summary.October[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.October[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					}
				} else {
					summary.October[v.Jetty.Name] = make(map[string]map[string]float64)
					summary.October[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
					summary.October[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
					summary.October[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
				}
			}
		}
	}

	errNovember := r.db.Table("productions").Select("pit_id, isp_id, jetty_id, SUM(quantity) as quantity, SUM(ritase_quantity) as ritase_quantity").Preload(clause.Associations).Where("iupopk_id = ? AND production_date >= ? AND production_date < ?", iupopkId, fmt.Sprintf("%v-11-01", year), fmt.Sprintf("%v-12-01", year)).Group("pit_id, isp_id, jetty_id").Find(&productionNovember).Error

	if errNovember != nil {
		return summary, errNovember
	}

	for _, v := range productionNovember {
		if v.Jetty != nil {
			if _, ok := summary.ListJettyPit[v.Jetty.Name]; !ok {
				summary.ListJettyPit[v.Jetty.Name] = []string{}
			}

			if v.Pit != nil {
				if !helperString(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name) {
					summary.ListJettyPit[v.Jetty.Name] = append(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name)
				}
				if _, ok := summary.November[v.Jetty.Name]; ok {
					if _, ok2 := summary.November[v.Jetty.Name][v.Pit.Name]; ok2 {
						summary.November[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.November[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					} else {
						summary.November[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
						summary.November[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.November[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					}
				} else {
					summary.November[v.Jetty.Name] = make(map[string]map[string]float64)
					summary.November[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
					summary.November[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
					summary.November[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
				}
			}
		}
	}

	errDecember := r.db.Table("productions").Select("pit_id, isp_id, jetty_id, SUM(quantity) as quantity, SUM(ritase_quantity) as ritase_quantity").Preload(clause.Associations).Where("iupopk_id = ? AND production_date >= ? AND production_date < ?", iupopkId, fmt.Sprintf("%v-12-01", year), fmt.Sprintf("%v-12-31", year)).Group("pit_id, isp_id, jetty_id").Find(&productionDecember).Error

	if errDecember != nil {
		return summary, errDecember
	}

	for _, v := range productionDecember {
		if v.Jetty != nil {
			if _, ok := summary.ListJettyPit[v.Jetty.Name]; !ok {
				summary.ListJettyPit[v.Jetty.Name] = []string{}
			}

			if v.Pit != nil {
				if !helperString(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name) {
					summary.ListJettyPit[v.Jetty.Name] = append(summary.ListJettyPit[v.Jetty.Name], v.Pit.Name)
				}
				if _, ok := summary.December[v.Jetty.Name]; ok {
					if _, ok2 := summary.December[v.Jetty.Name][v.Pit.Name]; ok2 {
						summary.December[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.December[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					} else {
						summary.December[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
						summary.December[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
						summary.December[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
					}
				} else {
					summary.December[v.Jetty.Name] = make(map[string]map[string]float64)
					summary.December[v.Jetty.Name][v.Pit.Name] = make(map[string]float64)
					summary.December[v.Jetty.Name][v.Pit.Name]["quantity"] = v.Quantity
					summary.December[v.Jetty.Name][v.Pit.Name]["ritase_quantity"] = v.RitaseQuantity
				}
			}
		}
	}

	return summary, nil
}
