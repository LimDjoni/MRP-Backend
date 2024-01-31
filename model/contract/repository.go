package contract

import (
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	GetListReportContractAll(page int, filterContract FilterAndSortContract, iupopkId int) (Pagination, error)
	GetDataContract(id int, iupopkId int) (Contract, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetListReportContractAll(page int, filterContract FilterAndSortContract, iupopkId int) (Pagination, error) {
	var listReportContract []Contract

	var pagination Pagination
	pagination.Limit = 7
	pagination.Page = page
	queryFilter := fmt.Sprintf("iupopk_id = %v", iupopkId)
	sortFilter := "id desc"

	if filterContract.Field != "" && filterContract.Sort != "" {
		sortFilter = filterContract.Field + " " + filterContract.Sort
	}

	if filterContract.ValidityStart != "" {
		queryFilter = queryFilter + "AND validity >= '" + filterContract.ValidityStart + "'"
	}

	if filterContract.ValidityEnd != "" {
		queryFilter = queryFilter + " AND validity <= '" + filterContract.ValidityEnd + "T23:59:59'"
	}

	if filterContract.ContractDateStart != "" {
		queryFilter = queryFilter + "AND contract_date >= '" + filterContract.ContractDateStart + "'"
	}

	if filterContract.ContractDateEnd != "" {
		queryFilter = queryFilter + " AND contract_date <= '" + filterContract.ContractDateEnd + "T23:59:59'"
	}

	if filterContract.Quantity != "" {
		quantity := fmt.Sprintf("%v", filterContract.Quantity)
		queryFilter = queryFilter + " AND cast(quantity AS TEXT) LIKE '%" + quantity + "%'"
	}

	if filterContract.ContractNumber != "" {
		queryFilter = queryFilter + " AND contract_number ILIKE '%" + filterContract.ContractNumber + "%'"
	}

	if filterContract.CustomerId != "" {
		queryFilter = queryFilter + " AND customer_id = " + filterContract.CustomerId
	}

	errFind := r.db.Where(queryFilter).Order(sortFilter).Scopes(paginateContract(listReportContract, &pagination, r.db, queryFilter)).Find(&listReportContract).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listReportContract

	return pagination, nil
}

func (r *repository) GetDataContract(id int, iupopkId int) (Contract, error) {
	var Contract Contract

	errFind := r.db.Where("id = ? AND iupopk_id = ?", id, iupopkId).First(&Contract).Error

	return Contract, errFind
}
