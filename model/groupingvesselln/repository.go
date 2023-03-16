package groupingvesselln

import (
	"ajebackend/helper"
	"ajebackend/model/insw"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	ListGroupingVesselLn(page int, sortFilter SortFilterGroupingVesselLn) (Pagination, error)
	ListGroupingVesselLnWithPeriod(month string, year int) ([]GroupingVesselLn, error)
	DetailInsw(id int) (DetailInsw, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) ListGroupingVesselLn(page int, sortFilter SortFilterGroupingVesselLn) (Pagination, error) {
	var listGroupingVesselLn []GroupingVesselLn
	var pagination Pagination

	pagination.Limit = 7
	pagination.Page = page
	querySort := "id desc"
	queryFilter := ""

	if sortFilter.Field != "" && sortFilter.Sort != "" {

		if sortFilter.Field == "vessel_id" {
			querySort = "vessels.name " + sortFilter.Sort
		}

		if querySort == "id desc" {
			querySort = fmt.Sprintf("%s %s", sortFilter.Field, sortFilter.Sort)
		}
	}

	if sortFilter.Quantity != "" {
		queryFilter = queryFilter + "cast(grand_total_quantity AS TEXT) LIKE '%" + sortFilter.Quantity + "%'"
	}

	if sortFilter.VesselId != "" {
		if queryFilter != "" {
			queryFilter += "AND vessel_id = " + sortFilter.VesselId
		} else {
			queryFilter = "vessel_id = " + sortFilter.VesselId
		}
	}

	if sortFilter.BlDateStart != "" {
		if queryFilter != "" {
			queryFilter += "AND bl_date >= '" + sortFilter.BlDateStart + "'"
		} else {
			queryFilter = "bl_date >= '" + sortFilter.BlDateStart + "'"
		}
	}

	if sortFilter.BlDateEnd != "" {
		if queryFilter != "" {
			queryFilter += "AND bl_date <= '" + sortFilter.BlDateEnd + "T23:59:59'"
		} else {
			queryFilter = "bl_date <= '" + sortFilter.BlDateEnd + "T23:59:59'"
		}
	}

	errFind := r.db.Preload(clause.Associations).Select("grouping_vessel_lns.*").Joins("LEFT JOIN vessels vessels on grouping_vessel_lns.vessel_id = vessels.id").Where(queryFilter).Order(querySort).Scopes(paginateData(listGroupingVesselLn, &pagination, r.db, queryFilter)).Find(&listGroupingVesselLn).Error

	if errFind != nil {
		errWithoutOrder := r.db.Preload(clause.Associations).Where(queryFilter).Order(querySort).Scopes(paginateData(listGroupingVesselLn, &pagination, r.db, queryFilter)).Find(&listGroupingVesselLn).Error

		if errWithoutOrder != nil {
			pagination.Data = listGroupingVesselLn
			return pagination, errWithoutOrder
		}
	}

	pagination.Data = listGroupingVesselLn

	return pagination, nil
}

func (r *repository) ListGroupingVesselLnWithPeriod(month string, year int) ([]GroupingVesselLn, error) {
	var listGroupingVesselLn []GroupingVesselLn
	var checkInsw insw.Insw

	errFindInsw := r.db.Where("month = ? AND year = ?", month, year).First(&checkInsw).Error

	if errFindInsw == nil {
		return listGroupingVesselLn, errors.New("Laporan INSW sudah pernah dibuat")
	}

	firstOfMonth := time.Date(year, time.Month(helper.MonthLongToNumber(month)), 1, 0, 0, 0, 0, time.Local)

	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	errFind := r.db.Preload(clause.Associations).Order("peb_register_date asc").Where("peb_register_date >= ? AND peb_register_date <= ? AND insw_id is NULL", firstOfMonth, lastOfMonth).Find(&listGroupingVesselLn).Error

	if errFind != nil {
		return listGroupingVesselLn, errFind
	}

	return listGroupingVesselLn, nil
}

func (r *repository) DetailInsw(id int) (DetailInsw, error) {
	var detailInsw DetailInsw

	var inswData insw.Insw
	errFindInsw := r.db.Where("id = ?", id).First(&inswData).Error

	if errFindInsw != nil {
		return detailInsw, errFindInsw
	}

	detailInsw.Detail = inswData

	var listGroupingVessel []GroupingVesselLn

	errFindListGroupingVessel := r.db.Preload(clause.Associations).Order("peb_register_date asc").Where("insw_id = ?", id).Find(&listGroupingVessel).Error

	if errFindListGroupingVessel != nil {
		return detailInsw, errFindListGroupingVessel
	}

	detailInsw.ListGroupingVesselLn = listGroupingVessel

	return detailInsw, nil
}
