package pitloss

type InputJettyPitLoss struct {
	Year         string         `json:"year" validate:"required"`
	JettyId      uint           `json:"jetty_id" validate:"required"`
	StartBalance float64        `json:"start_balance"`
	TotalLoss    float64        `json:"total_loss"`
	InputPitLoss []InputPitLoss `json:"input_pit_loss" validate:"required,min=1"`
}

type InputPitLoss struct {
	PitId                 uint    `json:"pit_id"`
	JanuaryLossQuantity   float64 `json:"january_loss_quantity"`
	FebruaryLossQuantity  float64 `json:"february_loss_quantity"`
	MarchLossQuantity     float64 `json:"march_loss_quantity"`
	AprilLossQuantity     float64 `json:"april_loss_quantity"`
	MayLossQuantity       float64 `json:"may_loss_quantity"`
	JuneLossQuantity      float64 `json:"june_loss_quantity"`
	JulyLossQuantity      float64 `json:"july_loss_quantity"`
	AugustLossQuantity    float64 `json:"august_loss_quantity"`
	SeptemberLossQuantity float64 `json:"september_loss_quantity"`
	OctoberLossQuantity   float64 `json:"october_loss_quantity"`
	NovemberLossQuantity  float64 `json:"november_loss_quantity"`
	DecemberLossQuantity  float64 `json:"december_loss_quantity"`
}

type InputUpdateJettyPitLoss struct {
	StartBalance float64   `json:"start_balance"`
	TotalLoss    float64   `json:"total_loss"`
	InputPitLoss []PitLoss `json:"input_pit_loss"`
}
