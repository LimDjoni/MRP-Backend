package royaltyrecon

import (
	"ajebackend/model/master/barge"
	"ajebackend/model/master/company"
	"ajebackend/model/master/tugboat"
)

type RoyaltyReconData struct {
	ShippingDate                string           `json:"shipping_date"`
	Quantity                    string           `json:"quantity"`
	TugboatId                   *uint            `json:"tugboat_id"`
	Tugboat                     *tugboat.Tugboat `json:"tugboat"`
	BargeId                     *uint            `json:"barge_id"`
	Barge                       *barge.Barge     `json:"barge"`
	CustomerId                  *uint            `json:"customer_id"`
	Customer                    *company.Company `json:"customer"`
	DmoBuyerId                  *uint            `json:"dmo_buyer_id"`
	DmoBuyer                    *company.Company `json:"dmo_buyer"`
	RoyaltyRate                 float64          `json:"royalty_rate"`
	DpRoyaltyDate               *string          `json:"dp_royalty_date"`
	DpRoyaltyNtpn               *string          `json:"dp_royalty_ntpn"`
	DpRoyaltyBillingCode        *string          `json:"dp_royalty_billing_code"`
	DpRoyaltyTotal              float64          `json:"dp_royalty_total"`
	PaymentDpRoyaltyDate        *string          `json:"payment_dp_royalty_date"`
	PaymentDpRoyaltyNtpn        *string          `json:"payment_dp_royalty_ntpn"`
	PaymentDpRoyaltyBillingCode *string          `json:"payment_dp_royalty_billing_code"`
	PaymentDpRoyaltyTotal       float64          `json:"payment_dp_royalty_total"`
}

type RoyaltyReconDetail struct {
	Detail          RoyaltyRecon       `json:"detail"`
	ListTransaction []RoyaltyReconData `json:"list_transaction"`
}
