package production

import (
	"ajebackend/model/master/isp"
	"ajebackend/model/master/jetty"
	"ajebackend/model/master/pit"
)

type OutputSummaryProduction struct {
	January      map[string]map[string]map[string]float64 `json:"january"`
	February     map[string]map[string]map[string]float64 `json:"february"`
	March        map[string]map[string]map[string]float64 `json:"march"`
	April        map[string]map[string]map[string]float64 `json:"april"`
	May          map[string]map[string]map[string]float64 `json:"may"`
	June         map[string]map[string]map[string]float64 `json:"june"`
	July         map[string]map[string]map[string]float64 `json:"july"`
	August       map[string]map[string]map[string]float64 `json:"august"`
	September    map[string]map[string]map[string]float64 `json:"september"`
	October      map[string]map[string]map[string]float64 `json:"october"`
	November     map[string]map[string]map[string]float64 `json:"november"`
	December     map[string]map[string]map[string]float64 `json:"december"`
	ListJettyPit map[string][]string                      `json:"list_jetty"`
}

type GroupProduction struct {
	RitaseQuantity float64      `json:"ritase_quantity"`
	Quantity       float64      `json:"quantity"`
	PitId          *uint        `json:"pit_id"`
	Pit            *pit.Pit     `json:"pit"`
	IspId          *uint        `json:"isp_id"`
	Isp            *isp.Isp     `json:"isp"`
	JettyId        *uint        `json:"jetty_id"`
	Jetty          *jetty.Jetty `json:"jetty"`
}
