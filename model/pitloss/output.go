package pitloss

import "ajebackend/model/jettybalance"

type OutputJettyBalancePitLossDetail struct {
	JettyBalance jettybalance.JettyBalance `json:"jetty_balance"`
	PitLoss      []PitLoss                 `json:"pit_loss"`
}
