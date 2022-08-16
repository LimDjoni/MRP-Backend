package user

type TokenUser struct {
	Username string `json:"username"`
	Email string `json:"email"`
	Token string  `json:"token"`
}
