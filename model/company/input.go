package company

type InputCreateUpdateCompany struct {
	CompanyName  string `json:"company_name" validate:"required"`
	Address     string	`json:"address"`
	Email       string	`json:"email" validate:"required,email"`
	Province	string	`json:"province"`
	PhoneNumber	string	`json:"phone_number"`
	FaxNumber	string	`json:"fax_number"`
}
