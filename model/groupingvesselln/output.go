package groupingvesselln

import "ajebackend/model/insw"

type DetailInsw struct {
	Detail               insw.Insw          `json:"detail"`
	ListGroupingVesselLn []GroupingVesselLn `json:"list_grouping_vessel_ln"`
}
