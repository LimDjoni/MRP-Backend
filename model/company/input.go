package company

type InputCreateUpdateCompany struct {
	CompanyName  string `json:"company_name" validate:"required"`
	IndustryType string `json:"industry_type"`
	Address      string `json:"address"`
	Province     string `json:"province" validate:"required"`
	PhoneNumber  string `json:"phone_number"`
	FaxNumber    string `json:"fax_number"`
}
