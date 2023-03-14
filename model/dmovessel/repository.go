package dmovessel

import (
	"ajebackend/model/groupingvesseldn"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	GetDataDmoVessel(id uint, iupopkId int) ([]DmoVessel, error)
	ListGroupingVesselWithoutDmo(iupopkId int) ([]groupingvesseldn.GroupingVesselDn, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetDataDmoVessel(id uint, iupopkId int) ([]DmoVessel, error) {
	var listDmoVessel []DmoVessel

	errFind := r.db.Preload(clause.Associations).Select("dmo_vessels.*").Joins("LEFT JOIN dmos dmos on dmo_vessels.dmo_id = dmos.id").Where("dmo_vessels.dmo_id = ? AND dmos.iupopk_id = ?", id, iupopkId).Find(&listDmoVessel).Error

	return listDmoVessel, errFind
}

func (r *repository) ListGroupingVesselWithoutDmo(iupopkId int) ([]groupingvesseldn.GroupingVesselDn, error) {

	var listGroupingVesselWithoutDmo []groupingvesseldn.GroupingVesselDn

	var dmoVessel []DmoVessel

	findDmoVesselErr := r.db.Find(&dmoVessel).Error

	if findDmoVesselErr != nil {
		return listGroupingVesselWithoutDmo, findDmoVesselErr
	}

	var listIdGroupingVessel []uint

	for _, v := range dmoVessel {
		listIdGroupingVessel = append(listIdGroupingVessel, v.GroupingVesselDnId)
	}

	if len(listIdGroupingVessel) == 0 {
		listIdGroupingVessel = append(listIdGroupingVessel, 0)
	}

	findListGroupingWithoutDmoErr := r.db.Preload(clause.Associations).Where("id NOT IN ? AND iupopk_id = ?", listIdGroupingVessel, iupopkId).Find(&listGroupingVesselWithoutDmo).Error

	if findListGroupingWithoutDmoErr != nil {
		return listGroupingVesselWithoutDmo, findListGroupingWithoutDmoErr
	}

	return listGroupingVesselWithoutDmo, nil
}
