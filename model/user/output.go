package user

import "ajebackend/model/master/iupopk"

type TokenUser struct {
	Username string          `json:"username"`
	Email    string          `json:"email"`
	Token    string          `json:"token"`
	Role     string          `json:"role"`
	Iupopk   []iupopk.Iupopk `json:"iupopk"`
}
