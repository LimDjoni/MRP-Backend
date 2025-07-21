package user

import "mrpbackend/model/employee"

type TokenUser struct {
	ID         uint              `json:"id"`
	EmployeeId uint              `json:"employee_id"`
	Username   string            `json:"username"`
	Email      string            `json:"email"`
	Token      string            `json:"token"`
	Role       []string          `json:"role"`
	Employee   employee.Employee `json:"employee"`
	CodeEmp    uint              `json:"code_emp"`
}
