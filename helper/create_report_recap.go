package helper

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func CreateReportRecap(sheetName string, file *excelize.File, data map[string]map[string]map[string]float64, dataProduction map[string]map[string]float64, recap map[string]interface{}) (*excelize.File, error) {

	categories := map[string]string{
		"A1":  "Bulan",
		"A2":  "Januari",
		"A3":  "Februari",
		"A4":  "Maret",
		"A5":  "April",
		"A6":  "Mei",
		"A7":  "Juni",
		"A8":  "Juli",
		"A9":  "Agustus",
		"A10": "September",
		"A11": "Oktober",
		"A12": "November",
		"A13": "Desember",
		"A14": "TOTAL",
		"B1":  "Kelistrikan",
		"C1":  "Non-Kelistrikan",
		"D1":  "Jumlah",
		"F1":  "Produksi",
		"G1":  "Tidak Bisa Klaim DMO",
	}

	values := make(map[string]float64)

	for key, value := range data {
		switch key {
		case "electricity":
			for keyValue, v := range value {
				switch keyValue {
				case "january":
					if len(v) == 0 {
						values["B14"] += 0
						values["B2"] += 0
						values["D2"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							fmt.Println("here")
							values["B14"] += value
							values["B2"] += value
							values["D2"] += value
							values["D14"] += value
						}
					}

				case "february":
					if len(v) == 0 {
						values["B14"] += 0
						values["B3"] += 0
						values["D3"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["B14"] += value
							values["B3"] += value
							values["D3"] += value
							values["D14"] += value
						}
					}
				case "march":
					if len(v) == 0 {
						values["B14"] += 0
						values["B4"] += 0
						values["D4"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["B14"] += value
							values["B4"] += value
							values["D4"] += value
							values["D14"] += value
						}
					}
				case "april":
					if len(v) == 0 {
						values["B14"] += 0
						values["B5"] += 0
						values["D5"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["B14"] += value
							values["B5"] += value
							values["D5"] += value
							values["D14"] += value
						}
					}
				case "may":
					if len(v) == 0 {
						values["B14"] += 0
						values["B6"] += 0
						values["D6"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["B14"] += value
							values["B6"] += value
							values["D6"] += value
							values["D14"] += value

						}
					}
				case "june":
					if len(v) == 0 {
						values["B14"] += 0
						values["B7"] += 0
						values["D7"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["B14"] += value
							values["B7"] += value
							values["D7"] += value
							values["D14"] += value
						}
					}
				case "july":
					if len(v) == 0 {
						values["B14"] += 0
						values["B8"] += 0
						values["D8"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["B14"] += value
							values["B8"] += value
							values["D8"] += value
							values["D14"] += value
						}
					}
				case "august":
					if len(v) == 0 {
						values["B14"] += 0
						values["B9"] += 0
						values["D9"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["B14"] += value
							values["B9"] += value
							values["D9"] += value
							values["D14"] += value
						}
					}
				case "september":
					if len(v) == 0 {
						values["B10"] += 0
						values["D10"] += 0
						values["D14"] += 0
						values["B14"] += 0
					} else {
						for _, value := range v {
							values["B10"] += value
							values["D10"] += value
							values["D14"] += value
							values["B14"] += value
						}
					}
				case "october":
					if len(v) == 0 {
						values["B14"] += 0
						values["B11"] += 0
						values["D11"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["B14"] += value
							values["B11"] += value
							values["D11"] += value
							values["D14"] += value
						}
					}
				case "november":
					if len(v) == 0 {
						values["B14"] += 0
						values["B12"] += 0
						values["D12"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["B14"] += value
							values["B12"] += value
							values["D12"] += value
							values["D14"] += value
						}
					}
				case "december":
					if len(v) == 0 {
						values["B14"] += 0
						values["B13"] += 0
						values["D13"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["B14"] += value
							values["B13"] += value
							values["D13"] += value
							values["D14"] += value
						}
					}
				}
			}

		case "non_electricity":
			for keyValue, v := range value {
				switch keyValue {
				case "january":
					if len(v) == 0 {
						values["C14"] += 0
						values["C2"] += 0
						values["D2"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["C14"] += value
							values["C2"] += value
							values["D2"] += value
							values["D14"] += value
						}
					}
				case "february":
					if len(v) == 0 {
						values["C14"] += 0
						values["C3"] += 0
						values["D3"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["C14"] += value
							values["C3"] += value
							values["D3"] += value
							values["D14"] += value
						}
					}
				case "march":
					if len(v) == 0 {
						values["C14"] += 0
						values["C4"] += 0
						values["D4"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["C14"] += value
							values["C4"] += value
							values["D4"] += value
							values["D14"] += value
						}
					}
				case "april":
					if len(v) == 0 {
						values["C14"] += 0
						values["C5"] += 0
						values["D5"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["C14"] += value
							values["C5"] += value
							values["D5"] += value
							values["D14"] += value
						}
					}
				case "may":
					if len(v) == 0 {
						values["C14"] += 0
						values["C6"] += 0
						values["D6"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["C14"] += value
							values["C6"] += value
							values["D6"] += value
							values["D14"] += value
						}
					}
				case "june":
					if len(v) == 0 {
						values["C14"] += 0
						values["C7"] += 0
						values["D7"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["C14"] += value
							values["C7"] += value
							values["D7"] += value
							values["D14"] += value
						}
					}
				case "july":
					if len(v) == 0 {
						values["C14"] += 0
						values["C8"] += 0
						values["D8"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["C14"] += value
							values["C8"] += value
							values["D8"] += value
							values["D14"] += value
						}
					}
				case "august":
					if len(v) == 0 {
						values["C14"] += 0
						values["C9"] += 0
						values["D9"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["C14"] += value
							values["C9"] += value
							values["D9"] += value
							values["D14"] += value
						}
					}
				case "september":
					if len(v) == 0 {
						values["C14"] += 0
						values["C10"] += 0
						values["D10"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["C14"] += value
							values["C10"] += value
							values["D10"] += value
							values["D14"] += value
						}
					}
				case "october":
					if len(v) == 0 {
						values["C14"] += 0
						values["C11"] += 0
						values["D11"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["C14"] += value
							values["C11"] += value
							values["D11"] += value
							values["D14"] += value
						}
					}
				case "november":
					if len(v) == 0 {
						values["C14"] += 0
						values["C12"] += 0
						values["D12"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["C14"] += value
							values["C12"] += value
							values["D12"] += value
							values["D14"] += value
						}
					}
				case "december":
					if len(v) == 0 {
						values["C14"] += 0
						values["C13"] += 0
						values["D13"] += 0
						values["D14"] += 0
					} else {
						for _, value := range v {
							values["C14"] += value
							values["C13"] += value
							values["D13"] += value
							values["D14"] += value
						}
					}
				}
			}
		}
	}

	for key, value := range dataProduction {
		switch key {
		case "production":
			for keyValue, v := range value {
				switch keyValue {
				case "january":
					values["F2"] += v
					values["F14"] += v
				case "february":
					values["F3"] += v
					values["F14"] += v
				case "march":
					values["F4"] += v
					values["F14"] += v
				case "april":
					values["F5"] += v
					values["F14"] += v
				case "may":
					values["F6"] += v
					values["F14"] += v
				case "june":
					values["F7"] += v
					values["F14"] += v
				case "july":
					values["F8"] += v
					values["F14"] += v
				case "august":
					values["F9"] += v
					values["F14"] += v
				case "september":
					values["F10"] += v
					values["F14"] += v
				case "october":
					values["F11"] += v
					values["F14"] += v
				case "november":
					values["F12"] += v
					values["F14"] += v
				case "december":
					values["F13"] += v
					values["F14"] += v
				}
			}
		case "not_claimable":
			for keyValue, v := range value {
				switch keyValue {
				case "january":
					values["G2"] += v
					values["G14"] += v
				case "february":
					values["G3"] += v
					values["G14"] += v
				case "march":
					values["G4"] += v
					values["G14"] += v
				case "april":
					values["G5"] += v
					values["G14"] += v
				case "may":
					values["G6"] += v
					values["G14"] += v
				case "june":
					values["G7"] += v
					values["G14"] += v
				case "july":
					values["G8"] += v
					values["G14"] += v
				case "august":
					values["G9"] += v
					values["G14"] += v
				case "september":
					values["G10"] += v
					values["G14"] += v
				case "october":
					values["G11"] += v
					values["G14"] += v
				case "november":
					values["G12"] += v
					values["G14"] += v
				case "december":
					values["G13"] += v
					values["G14"] += v
				}
			}
		}
	}

	for k, v := range categories {
		file.SetCellValue(sheetName, k, v)
	}

	for k, v := range values {
		file.SetCellValue(sheetName, k, v)
	}

	file.SetColWidth(sheetName, "A", "A", float64(16))
	file.SetColWidth(sheetName, "B", "D", float64(26))
	file.SetColWidth(sheetName, "F", "G", float64(26))

	custFmt := "#,##0.000"
	style, _ := file.NewStyle(&excelize.Style{CustomNumFmt: &custFmt})

	errStyle := file.SetColStyle(sheetName, "B:G", style)

	if errStyle != nil {
		return file, errStyle
	}

	border := []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}

	bold, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		CustomNumFmt: &custFmt,
	})

	boldStyle, _ := file.NewStyle(&excelize.Style{
		Border: border,
		Font: &excelize.Font{
			Bold: true,
		},
		CustomNumFmt: &custFmt,
	})

	borderStyle, _ := file.NewStyle(&excelize.Style{
		Border:       border,
		CustomNumFmt: &custFmt,
	})

	percentFmt := "0.00##"
	percentFmt += `\%`
	percentFmt += ";[Red](0.00##"
	percentFmt += `\%)`
	percentStyle, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		CustomNumFmt: &percentFmt,
	})

	errStyleBold := file.SetCellStyle(sheetName, "A1", "D1", boldStyle)

	if errStyleBold != nil {
		return file, errStyleBold
	}

	errStyleBoldProduction := file.SetCellStyle(sheetName, "F1", "G1", boldStyle)

	if errStyleBoldProduction != nil {
		return file, errStyleBoldProduction
	}

	errStyleBold14 := file.SetCellStyle(sheetName, "A14", "D14", boldStyle)

	if errStyleBold14 != nil {
		return file, errStyleBold14
	}

	errStyleBold14Production := file.SetCellStyle(sheetName, "F14", "G14", boldStyle)

	if errStyleBold14Production != nil {
		return file, errStyleBold14Production
	}

	errStyleBorder := file.SetCellStyle(sheetName, "A2", "D13", borderStyle)

	if errStyleBorder != nil {
		return file, errStyleBorder
	}

	errStyleBorderProduction := file.SetCellStyle(sheetName, "F2", "G13", borderStyle)

	if errStyleBorderProduction != nil {
		return file, errStyleBorderProduction
	}

	errBold := file.SetCellStyle(sheetName, "F16", "G22", bold)

	if errBold != nil {
		return file, errBold
	}

	errPercent := file.SetCellStyle(sheetName, "A16", "D22", percentStyle)

	if errPercent != nil {
		return file, errPercent
	}

	var seriesElectricAndNonElectric string

	seriesElectricAndNonElectric += fmt.Sprintf(`{
						"name": "%v!$B$1",
						"categories": "%v!$A$2:$A$13",
						"values": "%v!$B$2:$B$13",
						"marker": {
									"symbol": "square"
								} 
				},`, sheetName, sheetName, sheetName)

	seriesElectricAndNonElectric += fmt.Sprintf(`{
						"name": "%v!$C$1",
						"categories": "%v!$A$2:$A$13",
						"values": "%v!$C$2:$C$13",
						"marker": {
									"symbol": "square"
								} 
				}`, sheetName, sheetName, sheetName)

	valueSheetElectricAndNonElectric := fmt.Sprintf(`{
	    "type": "col",
	    "series": [%v],
			"dimension": {	
				"width": 720
			},
	    "format":
	    {
	        "x_scale": 1.0,
	        "y_scale": 1.0,
	        "x_offset": 15,
	        "y_offset": 10,
	        "print_obj": true,
	        "lock_aspect_ratio": false,
	        "locked": false
	    },
	    "legend":
	    {
	        "position": "top",
	        "show_legend_key": false
	    },
	    "title":
	    {
	        "name": "Jumlah Listrik & Non-Kelistrikan"
	    },
	    "show_blanks_as": "zero"
	}`, seriesElectricAndNonElectric)

	if err := file.AddChart(sheetName, "I1", valueSheetElectricAndNonElectric); err != nil {
		return file, err
	}

	var seriesTotalElectricAndNonElectric string

	seriesTotalElectricAndNonElectric += fmt.Sprintf(`{
					"name": "%v!$D$1",
					"categories": "%v!$A$2:$A$13",
					"values": "%v!$D$2:$D$13",
					"line": {
							"smooth": true,
							"width": 1
					},
					"marker": {
						"symbol": "square"
					} 
			}`, sheetName, sheetName, sheetName)

	valueSheetTotalElectricAndNonElectric := fmt.Sprintf(`{
	    "type": "line",
	    "series": [%v],
			"dimension": {	
				"width": 720
			},
	    "format":
	    {
	        "x_scale": 1.0,
	        "y_scale": 1.0,
	        "x_offset": 15,
	        "y_offset": 10,
	        "print_obj": true,
	        "lock_aspect_ratio": false,
	        "locked": false
	    },
	    "legend":
	    {
	        "position": "top",
	        "show_legend_key": false
	    },
	    "title":
	    {
	        "name": "Jumlah Kelistrikan & Non-Kelistrikan"
	    },
	    "show_blanks_as": "zero"
	}`, seriesTotalElectricAndNonElectric)

	if err := file.AddChart(sheetName, "I18", valueSheetTotalElectricAndNonElectric); err != nil {
		return file, err
	}

	var seriesProductionAndNotClaim string

	seriesProductionAndNotClaim += fmt.Sprintf(`{
						"name": "%v!$F$1",
						"categories": "%v!$A$2:$A$13",
						"values": "%v!$F$2:$F$13",
						"line": {
								"smooth": true,
								"width": 1
						},
						"marker": {
							"symbol": "square"
						} 
				},`, sheetName, sheetName, sheetName)

	seriesProductionAndNotClaim += fmt.Sprintf(`{
						"name": "%v!$G$1",
						"categories": "%v!$A$2:$A$13",
						"values": "%v!$G$2:$G$13",
						"line": {
								"smooth": true,
								"width": 1
						},
						"marker": {
							"symbol": "square"
						} 
				}`, sheetName, sheetName, sheetName)

	valueSheetProductionAndNotClaim := fmt.Sprintf(`{
	    "type": "line",
	    "series": [%v],
			"dimension": {	
				"width": 720
			},
	    "format":
	    {
	        "x_scale": 1.0,
	        "y_scale": 1.0,
	        "x_offset": 15,
	        "y_offset": 10,
	        "print_obj": true,
	        "lock_aspect_ratio": false,
	        "locked": false
	    },
	    "legend":
	    {
	        "position": "top",
	        "show_legend_key": false
	    },
	    "title":
	    {
	        "name": "Produksi & Tidak Bisa Klaim DMO"
	    },
	    "show_blanks_as": "zero"
	}`, seriesProductionAndNotClaim)

	if err := file.AddChart(sheetName, "I35", valueSheetProductionAndNotClaim); err != nil {
		return file, err
	}

	file.MergeCell(sheetName, "A16", "C16")
	file.MergeCell(sheetName, "A18", "C18")
	file.MergeCell(sheetName, "A20", "C20")
	file.MergeCell(sheetName, "A22", "C22")

	s := "% Pemenuhan DMO terhadap kewajiban pemenuhan DMO "
	s += fmt.Sprintf("%v", recap["percentage_production_obligation"])
	s += "%"
	file.SetCellValue(sheetName, "A16", "% Pemenuhan DMO terhadap realisasi produksi")
	file.SetCellValue(sheetName, "A18", "% Pemenuhan DMO terhadap rencana produksi")
	file.SetCellValue(sheetName, "A20", "% Pemenuhan DMO terhadap rencana produksi (prorata 12 bulan)")
	file.SetCellValue(sheetName, "A22", s)

	fullfilmentOfProductionPlan := strings.Replace(recap["fulfillment_of_production_plan"].(string), "%", "", -1)
	fullfilmentOfProductionPlanRealization := strings.Replace(recap["fulfillment_of_production_realization"].(string), "%", "", -1)
	fullfilmentPercentageProductionObligation := strings.Replace(recap["fulfillment_percentage_production_obligation"].(string), "%", "", -1)
	prorateProductionPlan := strings.Replace(recap["prorate_production_plan"].(string), "%", "", -1)

	fullfilmentOfProductionPlanFloat, _ := strconv.ParseFloat(fullfilmentOfProductionPlan, 64)
	fullfilmentOfProductionPlanRealizationFloat, _ := strconv.ParseFloat(fullfilmentOfProductionPlanRealization, 64)
	fullfilmentPercentageProductionObligationFloat, _ := strconv.ParseFloat(fullfilmentPercentageProductionObligation, 64)
	prorateProductionPlanFloat, _ := strconv.ParseFloat(prorateProductionPlan, 64)

	file.SetCellValue(sheetName, "D16", fullfilmentOfProductionPlanRealizationFloat)
	file.SetCellValue(sheetName, "D18", fullfilmentOfProductionPlanFloat)
	file.SetCellValue(sheetName, "D20", prorateProductionPlanFloat)
	file.SetCellValue(sheetName, "D22", fullfilmentPercentageProductionObligationFloat)

	file.SetCellValue(sheetName, "F18", recap["production_plan"])
	file.SetCellValue(sheetName, "F22", recap["production_plan"])

	file.SetCellValue(sheetName, "G18", "Quota RKAB")
	file.SetCellValue(sheetName, "G22", "prorata 12 bulan")

	return file, nil
}
