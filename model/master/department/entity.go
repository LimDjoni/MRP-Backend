package department

import (
	"gorm.io/gorm"
)

type Department struct {
	gorm.Model
	DepartmentName string `json:"department_name"`
}
