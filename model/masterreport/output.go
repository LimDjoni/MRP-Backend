package masterreport

import (
	"ajebackend/model/cafassignment"
	"ajebackend/model/electricassignment"
	"ajebackend/model/master/company"
	"ajebackend/model/rkab"
)

// Report Recap DMO Output

type ReportDmoOutput struct {
	RecapElectricity              RecapElectricity                      `json:"recap_electricity"`
	RecapNonElectricity           RecapNonElectricity                   `json:"recap_non_electricity"`
	NotClaimable                  NotClaimable                          `json:"not_claimable"`
	Production                    QuantityProduction                    `json:"production"`
	Rkabs                         []rkab.Rkab                           `json:"rkabs"`
	CafAssignment                 cafassignment.CafAssignment           `json:"caf_assignment"`
	RealizationCafAssignment      float64                               `json:"realization_caf_assignment"`
	ElectricAssignment            electricassignment.ElectricAssignment `json:"electric_assignment"`
	RealizationElectricAssignment float64                               `json:"realization_electric_assignment"`
}

type RecapElectricity struct {
	January   float64 `json:"january"`
	February  float64 `json:"february"`
	March     float64 `json:"march"`
	April     float64 `json:"april"`
	May       float64 `json:"may"`
	June      float64 `json:"june"`
	July      float64 `json:"july"`
	August    float64 `json:"august"`
	September float64 `json:"september"`
	October   float64 `json:"october"`
	November  float64 `json:"november"`
	December  float64 `json:"december"`
	Total     float64 `json:"total"`
}

type RecapNonElectricity struct {
	January   float64 `json:"january"`
	February  float64 `json:"february"`
	March     float64 `json:"march"`
	April     float64 `json:"april"`
	May       float64 `json:"may"`
	June      float64 `json:"june"`
	July      float64 `json:"july"`
	August    float64 `json:"august"`
	September float64 `json:"september"`
	October   float64 `json:"october"`
	November  float64 `json:"november"`
	December  float64 `json:"december"`
	Total     float64 `json:"total"`
}

type NotClaimable struct {
	January   float64 `json:"january"`
	February  float64 `json:"february"`
	March     float64 `json:"march"`
	April     float64 `json:"april"`
	May       float64 `json:"may"`
	June      float64 `json:"june"`
	July      float64 `json:"july"`
	August    float64 `json:"august"`
	September float64 `json:"september"`
	October   float64 `json:"october"`
	November  float64 `json:"november"`
	December  float64 `json:"december"`
	Total     float64 `json:"total"`
}

type QuantityProduction struct {
	January   float64 `json:"january"`
	February  float64 `json:"february"`
	March     float64 `json:"march"`
	April     float64 `json:"april"`
	May       float64 `json:"may"`
	June      float64 `json:"june"`
	July      float64 `json:"july"`
	August    float64 `json:"august"`
	September float64 `json:"september"`
	October   float64 `json:"october"`
	November  float64 `json:"november"`
	December  float64 `json:"december"`
	Total     float64 `json:"total"`
}

// Report Realization Output

type RealizationOutput struct {
	Electric RealizationElectric `json:"electric"`
	Caf      RealizationCaf      `json:"caf"`
}

type RealizationElectric struct {
	January   []RealizationTransaction `json:"january"`
	February  []RealizationTransaction `json:"february"`
	March     []RealizationTransaction `json:"march"`
	April     []RealizationTransaction `json:"april"`
	May       []RealizationTransaction `json:"may"`
	June      []RealizationTransaction `json:"june"`
	July      []RealizationTransaction `json:"july"`
	August    []RealizationTransaction `json:"august"`
	September []RealizationTransaction `json:"september"`
	October   []RealizationTransaction `json:"october"`
	November  []RealizationTransaction `json:"november"`
	December  []RealizationTransaction `json:"december"`
}

type RealizationCaf struct {
	January   []RealizationTransaction `json:"january"`
	February  []RealizationTransaction `json:"february"`
	March     []RealizationTransaction `json:"march"`
	April     []RealizationTransaction `json:"april"`
	May       []RealizationTransaction `json:"may"`
	June      []RealizationTransaction `json:"june"`
	July      []RealizationTransaction `json:"july"`
	August    []RealizationTransaction `json:"august"`
	September []RealizationTransaction `json:"september"`
	October   []RealizationTransaction `json:"october"`
	November  []RealizationTransaction `json:"november"`
	December  []RealizationTransaction `json:"december"`
}

type RealizationTransaction struct {
	ShippingDate      string          `json:"shipping_date"`
	Trader            company.Company `json:"trader"`
	EndUser           company.Company `json:"end_user"`
	QualityCaloriesAr float64         `json:"quality_calories_ar"`
	Quantity          float64         `json:"quantity"`
	IsBastOk          bool            `json:"is_bast_ok"`
}

// Report Detail Output

type SaleDetail struct {
	Rkabs                               []rkab.Rkab                           `json:"rkabs"`
	Electricity                         Electricity                           `json:"electricity"`
	NonElectricity                      NonElectricity                        `json:"non_electricity"`
	RecapElectricity                    RecapElectricity                      `json:"recap_electricity"`
	RecapNonElectricity                 RecapNonElectricity                   `json:"recap_non_electricity"`
	NotClaimable                        NotClaimable                          `json:"not_claimable"`
	Production                          QuantityProduction                    `json:"production"`
	Domestic                            Domestic                              `json:"domestic"`
	Export                              Export                                `json:"export"`
	CafAssignment                       cafassignment.CafAssignment           `json:"caf_assignment"`
	RealizationCafAssignment            float64                               `json:"realization_caf_assignment"`
	ElectricAssignment                  electricassignment.ElectricAssignment `json:"electric_assignment"`
	RealizationElectricAssignment       float64                               `json:"realization_electric_assignment"`
	CompanyElectricity                  []string                              `json:"company_electricity"`
	CompanyNonElectricity               []string                              `json:"company_non_electricity"`
	DetailRealizationCafAssignment      DetailRealizationCafAssignment        `json:"detail_realization_caf_assignment"`
	DetailRealizationElectricAssignment DetailRealizationElectricAssignment   `json:"detail_realization_electric_assignment"`
	RealizationCompanyElectricity       []string                              `json:"realization_company_electricity"`
	RealizationCompanyCaf               []string                              `json:"realization_company_caf"`
}

type DetailRealizationCafAssignment struct {
	January   map[string]float64 `json:"january"`
	February  map[string]float64 `json:"february"`
	March     map[string]float64 `json:"march"`
	April     map[string]float64 `json:"april"`
	May       map[string]float64 `json:"may"`
	June      map[string]float64 `json:"june"`
	July      map[string]float64 `json:"july"`
	August    map[string]float64 `json:"august"`
	September map[string]float64 `json:"september"`
	October   map[string]float64 `json:"october"`
	November  map[string]float64 `json:"november"`
	December  map[string]float64 `json:"december"`
	Total     float64            `json:"total"`
}

type DetailRealizationElectricAssignment struct {
	January   map[string]float64 `json:"january"`
	February  map[string]float64 `json:"february"`
	March     map[string]float64 `json:"march"`
	April     map[string]float64 `json:"april"`
	May       map[string]float64 `json:"may"`
	June      map[string]float64 `json:"june"`
	July      map[string]float64 `json:"july"`
	August    map[string]float64 `json:"august"`
	September map[string]float64 `json:"september"`
	October   map[string]float64 `json:"october"`
	November  map[string]float64 `json:"november"`
	December  map[string]float64 `json:"december"`
	Total     float64            `json:"total"`
}

type Electricity struct {
	January   map[string]float64 `json:"january"`
	February  map[string]float64 `json:"february"`
	March     map[string]float64 `json:"march"`
	April     map[string]float64 `json:"april"`
	May       map[string]float64 `json:"may"`
	June      map[string]float64 `json:"june"`
	July      map[string]float64 `json:"july"`
	August    map[string]float64 `json:"august"`
	September map[string]float64 `json:"september"`
	October   map[string]float64 `json:"october"`
	November  map[string]float64 `json:"november"`
	December  map[string]float64 `json:"december"`
	Total     float64            `json:"total"`
}

type NonElectricity struct {
	January   map[string]float64 `json:"january"`
	February  map[string]float64 `json:"february"`
	March     map[string]float64 `json:"march"`
	April     map[string]float64 `json:"april"`
	May       map[string]float64 `json:"may"`
	June      map[string]float64 `json:"june"`
	July      map[string]float64 `json:"july"`
	August    map[string]float64 `json:"august"`
	September map[string]float64 `json:"september"`
	October   map[string]float64 `json:"october"`
	November  map[string]float64 `json:"november"`
	December  map[string]float64 `json:"december"`
	Total     float64            `json:"total"`
}

type Domestic struct {
	January   float64 `json:"january"`
	February  float64 `json:"february"`
	March     float64 `json:"march"`
	April     float64 `json:"april"`
	May       float64 `json:"may"`
	June      float64 `json:"june"`
	July      float64 `json:"july"`
	August    float64 `json:"august"`
	September float64 `json:"september"`
	October   float64 `json:"october"`
	November  float64 `json:"november"`
	December  float64 `json:"december"`
	Total     float64 `json:"total"`
}

type Export struct {
	January   float64 `json:"january"`
	February  float64 `json:"february"`
	March     float64 `json:"march"`
	April     float64 `json:"april"`
	May       float64 `json:"may"`
	June      float64 `json:"june"`
	July      float64 `json:"july"`
	August    float64 `json:"august"`
	September float64 `json:"september"`
	October   float64 `json:"october"`
	November  float64 `json:"november"`
	December  float64 `json:"december"`
	Total     float64 `json:"total"`
}
