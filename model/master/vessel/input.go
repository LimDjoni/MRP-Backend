package vessel

type InputVessel struct {
	Name            string   `json:"vessel_name"`
	Deadweight      *float64 `json:"deadweight"`
	MinimumQuantity *float64 `json:"minimum_quantity"`
	MaximumQuantity *float64 `json:"maximum_quantity"`
}
