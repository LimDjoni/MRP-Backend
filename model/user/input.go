package user

type RegisterUserInput struct {
	EmployeeId uint   `json:"employee_id" validate:"required"`
	Username   string `json:"username" validate:"required"`
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8,max=12"`
	CodeEmp    uint   `json:"code_emp" validate:"required"`
}

type LoginUserInput struct {
	Data     string `json:"data"`
	Password string `json:"password"`
}

type ChangePasswordInput struct {
	OldPassword     string `json:"old_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8,max=12"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=8,max=12"`
}

type ResetPasswordInput struct {
	Email string `json:"email"`
}
