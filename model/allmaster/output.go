package allmaster

import (
	"mrpbackend/model/master/brand"
	"mrpbackend/model/master/department"
	"mrpbackend/model/master/departmentform"
	"mrpbackend/model/master/form"
	"mrpbackend/model/master/heavyequipment"
	"mrpbackend/model/master/role"
	"mrpbackend/model/master/roleform"
	"mrpbackend/model/master/series"
	"mrpbackend/model/master/userrole"
)

type MasterData struct {
	Brand          []brand.Brand                   `json:"brand"`
	HeavyEquipment []heavyequipment.HeavyEquipment `json:"heavy_equipment"`
	Series         []series.Series                 `json:"series"`
	Department     []department.Department         `json:"department"`
	DepartmentForm []departmentform.DepartmentForm `json:"department_form"`
	Form           []form.Form                     `json:"form"`
	Role           []role.Role                     `json:"role"`
	RoleForm       []roleform.RoleForm             `json:"role_form"`
	UserRole       []userrole.UserRole             `json:"user_role"`
}

