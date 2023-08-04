package allmaster

type InputBarge struct {
	Name            string  `json:"name"`
	Height          float64 `json:"height"`
	Deadweight      float64 `json:"deadweight"`
	MinimumQuantity float64 `json:"minimum_quantity"`
	MaximumQuantity float64 `json:"maximum_quantity"`
}

type InputTugboat struct {
	Name            string  `json:"name"`
	Height          float64 `json:"height"`
	Deadweight      float64 `json:"deadweight"`
	MinimumQuantity float64 `json:"minimum_quantity"`
	MaximumQuantity float64 `json:"maximum_quantity"`
}

type InputVessel struct {
	Name            string  `json:"name"`
	Deadweight      float64 `json:"deadweight"`
	MinimumQuantity float64 `json:"minimum_quantity"`
	MaximumQuantity float64 `json:"maximum_quantity"`
}

type InputPortLocation struct {
	Name string `json:"name"`
}

type InputPort struct {
	Name                 string `json:"name"`
	PortLocationId       uint   `json:"port_location_id"`
	IsLoadingPort        bool   `json:"is_loading_port"`
	IsUnloadingPort      bool   `json:"is_unloading_port"`
	IsDmoDestinationPort bool   `json:"is_dmo_destination_port"`
}

type InputCompany struct {
	CompanyName    string `json:"company_name" validate:"required"`
	IndustryTypeId *uint  `json:"industry_type_id"`
	Address        string `json:"address"`
	Province       string `json:"province" validate:"required"`
	PhoneNumber    string `json:"phone_number"`
	FaxNumber      string `json:"fax_number"`
	IsTrader       bool   `json:"is_trader"`
	IsEndUser      bool   `json:"is_end_user"`
}

type InputTrader struct {
	TraderName string  `json:"trader_name" validate:"required"`
	Position   string  `json:"position" validate:"required"`
	Email      *string `json:"email" validate:"omitempty,email"`
	CompanyId  int     `json:"company_id" validate:"required"`
}

type InputIndustryType struct {
	Name                   string `json:"name" validate:"required"`
	CategoryIndustryTypeId uint   `json:"category_industry_type_id" validate:"required"`
	SystemCategory         string `json:"system_category" validate:"required"`
}
