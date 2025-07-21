package roleform

import (
	"gorm.io/gorm"
)

type RoleForm struct {
	gorm.Model
	DepartmentFormId uint `json:"department_form_id"`
	FormId           uint `json:"form_id"`
	CreateFlag       bool `json:"create_flag"`
	UpdateFlag       bool `json:"update_flag"`
	ReadFlag         bool `json:"read_flag"`
	DeleteFlag       bool `json:"delete_flag"`
}
