package pitloss

import (
	"ajebackend/model/jettybalance"
	"ajebackend/model/master/pit"

	"gorm.io/gorm"
)

type PitLoss struct {
	gorm.Model
	PitId                 uint                      `json:"pit_id"`
	Pit                   pit.Pit                   `json:"pit"`
	JanuaryLossQuantity   float64                   `json:"january_loss_quantity"`
	FebruaryLossQuantity  float64                   `json:"february_loss_quantity"`
	MarchLossQuantity     float64                   `json:"march_loss_quantity"`
	AprilLossQuantity     float64                   `json:"april_loss_quantity"`
	MayLossQuantity       float64                   `json:"may_loss_quantity"`
	JuneLossQuantity      float64                   `json:"june_loss_quantity"`
	JulyLossQuantity      float64                   `json:"july_loss_quantity"`
	AugustLossQuantity    float64                   `json:"august_loss_quantity"`
	SeptemberLossQuantity float64                   `json:"september_loss_quantity"`
	OctoberLossQuantity   float64                   `json:"october_loss_quantity"`
	NovemberLossQuantity  float64                   `json:"november_loss_quantity"`
	DecemberLossQuantity  float64                   `json:"december_loss_quantity"`
	JettyBalanceId        uint                      `json:"jetty_balance_id"`
	JettyBalance          jettybalance.JettyBalance `json:"jetty_balance" gorm:"constraint:OnDelete:CASCADE;"`
}
