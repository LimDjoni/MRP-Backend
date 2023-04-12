package user

type RegisterUserInput struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=12"`
}

type LoginUserInput struct {
	Data     string `json:"data"`
	Password string `json:"password"`
}

type ChangePasswordInput struct {
	OldPassword     string `json:"old_password" validate:"required,min=6,max=12"`
	NewPassword     string `json:"new_password" validate:"required,min=6,max=12"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=12"`
}

type ResetPasswordInput struct {
	Email string `json:"email"`
}
