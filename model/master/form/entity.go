package form

import (
	"gorm.io/gorm"
)

type Form struct {
	gorm.Model
	ID         uint    `json:"id"`
	FormName   string  `json:"form_name"`
	Path       *string `json:"path,omitempty"`      // Nullable, because some items are just parents
	ParentID   *uint   `json:"parent_id,omitempty"` // Nullable, for root-level items
	Sequence   uint    `json:"sequence"`
	Children   []Form  `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	CreateFlag bool    `json:"create_flag"`
	UpdateFlag bool    `json:"update_flag"`
	ReadFlag   bool    `json:"read_flag"`
	DeleteFlag bool    `json:"delete_flag"`
}
