package ports

import (
	"ajebackend/model/master/portlocation"

	"gorm.io/gorm"
)

type Port struct {
	gorm.Model
	Name                 string                    `json:"name" gorm:"UNIQUE"`
	PortLocationId       uint                      `json:"port_location_id"`
	PortLocation         portlocation.PortLocation `json:"port_location"`
	IsLoadingPort        bool                      `json:"is_loading_port"`
	IsUnloadingPort      bool                      `json:"is_unloading_port"`
	IsDmoDestinationPort bool                      `json:"is_dmo_destination_port"`
}
