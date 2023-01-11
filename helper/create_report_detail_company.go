package helper

import (
	"fmt"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"
)

func CreateReportDetailCompany(companies []string, sheetName string, file *excelize.File, data map[string]map[string]float64, title string) (*excelize.File, error) {

	categories := map[string]string{
		"B1": "Januari",
		"C1": "Februari",
		"D1": "Maret",
		"E1": "April",
		"F1": "Mei",
		"G1": "Juni",
		"H1": "Juli",
		"I1": "Agustus",
		"J1": "September",
		"K1": "Oktober",
		"L1": "November",
		"M1": "Desember",
	}

	mappingCompanyExcel := make(map[string]int)
	values := make(map[string]float64)
	for idx, v := range companies {
		formatIndex := fmt.Sprintf("A%v", idx+2)
		categories[formatIndex] = v
		mappingCompanyExcel[v] = idx + 2
	}
	for k, v := range categories {
		file.SetCellValue(sheetName, k, v)
	}

	for key, _ := range data {
		switch key {
		case "january":
			for _, companyName := range companies {
				formatIndexMonth := fmt.Sprintf("B%v", mappingCompanyExcel[companyName])
				if data["january"][companyName] > 0 {
					values[formatIndexMonth] = data["january"][companyName]
				} else {
					values[formatIndexMonth] = 0
				}
			}
		case "february":
			for _, companyName := range companies {
				formatIndexMonth := fmt.Sprintf("C%v", mappingCompanyExcel[companyName])
				if data["february"][companyName] > 0 {
					values[formatIndexMonth] = data["february"][companyName]
				} else {
					values[formatIndexMonth] = 0
				}
			}
		case "march":
			for _, companyName := range companies {
				formatIndexMonth := fmt.Sprintf("D%v", mappingCompanyExcel[companyName])
				if data["march"][companyName] > 0 {
					values[formatIndexMonth] = data["march"][companyName]
				} else {
					values[formatIndexMonth] = 0
				}
			}
		case "april":
			for _, companyName := range companies {
				formatIndexMonth := fmt.Sprintf("E%v", mappingCompanyExcel[companyName])
				if data["april"][companyName] > 0 {
					values[formatIndexMonth] = data["april"][companyName]
				} else {
					values[formatIndexMonth] = 0
				}
			}
		case "may":
			for _, companyName := range companies {
				formatIndexMonth := fmt.Sprintf("F%v", mappingCompanyExcel[companyName])
				if data["may"][companyName] > 0 {
					values[formatIndexMonth] = data["may"][companyName]
				} else {
					values[formatIndexMonth] = 0
				}
			}
		case "june":
			for _, companyName := range companies {
				formatIndexMonth := fmt.Sprintf("G%v", mappingCompanyExcel[companyName])
				if data["june"][companyName] > 0 {
					values[formatIndexMonth] = data["june"][companyName]
				} else {
					values[formatIndexMonth] = 0
				}
			}
		case "july":
			for _, companyName := range companies {
				formatIndexMonth := fmt.Sprintf("H%v", mappingCompanyExcel[companyName])
				if data["july"][companyName] > 0 {
					values[formatIndexMonth] = data["july"][companyName]
				} else {
					values[formatIndexMonth] = 0
				}
			}
		case "august":
			for _, companyName := range companies {
				formatIndexMonth := fmt.Sprintf("I%v", mappingCompanyExcel[companyName])
				if data["august"][companyName] > 0 {
					values[formatIndexMonth] = data["august"][companyName]
				} else {
					values[formatIndexMonth] = 0
				}
			}
		case "september":
			for _, companyName := range companies {
				formatIndexMonth := fmt.Sprintf("J%v", mappingCompanyExcel[companyName])
				if data["september"][companyName] > 0 {
					values[formatIndexMonth] = data["september"][companyName]
				} else {
					values[formatIndexMonth] = 0
				}
			}
		case "october":
			for _, companyName := range companies {
				formatIndexMonth := fmt.Sprintf("K%v", mappingCompanyExcel[companyName])
				if data["october"][companyName] > 0 {
					values[formatIndexMonth] = data["october"][companyName]
				} else {
					values[formatIndexMonth] = 0
				}
			}
		case "november":
			for _, companyName := range companies {
				formatIndexMonth := fmt.Sprintf("L%v", mappingCompanyExcel[companyName])
				if data["november"][companyName] > 0 {
					values[formatIndexMonth] = data["november"][companyName]
				} else {
					values[formatIndexMonth] = 0
				}
			}
		case "december":
			for _, companyName := range companies {
				formatIndexMonth := fmt.Sprintf("M%v", mappingCompanyExcel[companyName])
				if data["december"][companyName] > 0 {
					values[formatIndexMonth] = data["december"][companyName]
				} else {
					values[formatIndexMonth] = 0
				}
			}
		}
	}

	for k, v := range values {
		file.SetCellValue(sheetName, k, v)
	}

	var series string
	for idx, _ := range companies {
		if idx < len(companies)-1 {
			series += fmt.Sprintf(`{
		            "name": "%v!$A$%v",
		            "categories": "%v!$B$1:$M$1",
		            "values": "%v!$B$%v:$M$%v",
		            "line": {
		                "smooth": true,
		                "width": 1
		            },
								"marker": {
									"symbol": "square"
								}
		        },`, sheetName, idx+2, sheetName, sheetName, idx+2, idx+2)
		} else {
			series += fmt.Sprintf(`{
								"name": "%v!$A$%v",
								"categories": "%v!$B$1:$M$1",
								"values": "%v!$B$%v:$M$%v",
								"line": {
										"smooth": true,
										"width": 1
								},
								"marker": {
									"symbol": "square"
								} 
						}`, sheetName, idx+2, sheetName, sheetName, idx+2, idx+2)
		}
	}

	valueSheet := fmt.Sprintf(`{
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
	        "name": "%v"
	    },
	    "show_blanks_as": "zero"
	}`, series, title)

	custFmt := "#,##0.000"

	border := []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}

	boldStyle, _ := file.NewStyle(&excelize.Style{
		Border: border,
		Font: &excelize.Font{
			Bold: true,
		},
	})

	borderStyle, _ := file.NewStyle(&excelize.Style{
		Border:       border,
		CustomNumFmt: &custFmt,
	})

	errStyleBold := file.SetCellStyle(sheetName, "A1", "M1", boldStyle)

	if errStyleBold != nil {
		return file, errStyleBold
	}

	lastRow := fmt.Sprintf("A%v", len(companies)+1)

	lastColumnRow := fmt.Sprintf("M%v", len(companies)+1)

	errStyleBoldColumnA := file.SetCellStyle(sheetName, "A2", lastRow, boldStyle)

	if errStyleBoldColumnA != nil {
		return file, errStyleBoldColumnA
	}

	errStyleBorder := file.SetCellStyle(sheetName, "B2", lastColumnRow, borderStyle)

	if errStyleBorder != nil {
		return file, errStyleBorder
	}

	largestWidth := 0
	for _, value := range companies {
		cellWidth := utf8.RuneCountInString(value) + 2 // + 1 for margin
		if cellWidth > largestWidth {
			largestWidth = cellWidth
		}
	}

	file.SetColWidth(sheetName, "A", "A", float64(largestWidth))

	file.SetColWidth(sheetName, "B", "M", float64(16))

	if err := file.AddChart(sheetName, "N1", valueSheet); err != nil {
		return file, err
	}

	return file, nil
}
