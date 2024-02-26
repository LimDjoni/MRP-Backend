package royaltyreport

import (
	"ajebackend/model/contract"
	"ajebackend/model/ici"

	"gorm.io/gorm"
)

type RoyaltyReport struct {
	gorm.Model
	PicName               string             `json:"pic_name"`
	ContractNumber        string             `json:"contract_number" gorm:"UNIQUE"`
	EarlyLycanDate        string             `json:"early_lycan_date" gorm:"type:DATE"`
	LatestLycanDate       string             `json:"latest_lycan_date" gorm:"type:DATE"`
	PaymentType           string             `json:"payment_type"`
	Validity              string             `json:"validity" gorm:"type:DATE"`
	Pricelist             string             `json:"pricelist"`
	PaymentTerms          string             `json:"payment_term"`
	IciId                 *uint              `json:"ici_id"`
	Ici                   *ici.Ici           `json:"ici" gorm:"constraint:OnDelete:SET NULL;"`
	EmailNominationVessel string             `json:"email_nomination_vessel"`
	Product               string             `json:"product"`
	Description           string             `json:"description"`
	Quantity              float64            `json:"quantity"`
	Route                 string             `json:"route"`
	Measure               string             `json:"measure"`
	UnitPerPrice          float64            `json:"unit_per_price"`
	Taxes                 float64            `json:"taxes"`
	Discount              float64            `json:"discount"`
	SubTotal              float64            `json:"sub_total"`
	ContractId            *uint              `json:"contract_id"`
	Contract              *contract.Contract `json:"contract" gorm:"constraint:OnDelete:SET NULL;"`
}

type SortFilterRoyaltyReport struct {
	Field     string
	Sort      string
	DateStart string
	DateEnd   string
}
