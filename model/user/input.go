package user

type RegisterUserInput struct {
	Username string `json:"username" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=12"`
}

type LoginUserInput struct {
	Data string `json:"data"`
	Password string `json:"password"`
}
