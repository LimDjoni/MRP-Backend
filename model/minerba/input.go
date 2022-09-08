package minerba

type InputCreateMinerba struct {
	Period string `json:"period" validate:"PeriodValidation,required"`
	ListDataDn []int `json:"list_data_dn" validate:"required,min=1"`
}
