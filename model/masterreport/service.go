package masterreport

import (
	"fmt"
	"reflect"

	"ajebackend/model/master/iupopk"

	"github.com/nleeper/goment"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Service interface {
	RecapDmo(year string, iupopkId int) (ReportDmoOutput, error)
	RealizationReport(year string, iupopkId int) (RealizationOutput, error)
	SaleDetailReport(year string, iupopkId int) (SaleDetail, error)
	CreateReportRecapDmo(year string, reportRecapDmo ReportDmoOutput, iupopk iupopk.Iupopk, file *excelize.File, sheetName string) (*excelize.File, error)
	CreateReportRealization(year string, reportRealization RealizationOutput, iupopk iupopk.Iupopk, file *excelize.File) (*excelize.File, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func createTemplateRealization(file *excelize.File, sheetName []string, year string) (*excelize.File, error) {
	boldStyleCenter, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
			WrapText:   true,
		},
	})

	border := []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}

	boldStyleBorder, _ := file.NewStyle(&excelize.Style{
		Border: border,
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	for _, v := range sheetName {
		mergeErr1 := file.MergeCell(v, "B2", "H2")
		if mergeErr1 != nil {
			return file, mergeErr1
		}

		mergeErr2 := file.MergeCell(v, "B3", "H3")
		if mergeErr2 != nil {
			return file, mergeErr2
		}

		mergeErr3 := file.MergeCell(v, "D5", "F5")
		if mergeErr3 != nil {
			return file, mergeErr3
		}

		mergeErr4 := file.MergeCell(v, "B5", "B6")
		if mergeErr4 != nil {
			return file, mergeErr4
		}

		mergeErr5 := file.MergeCell(v, "C5", "C6")
		if mergeErr5 != nil {
			return file, mergeErr5
		}

		mergeErr6 := file.MergeCell(v, "G5", "G6")
		if mergeErr6 != nil {
			return file, mergeErr6
		}

		mergeErr7 := file.MergeCell(v, "H5", "H6")
		if mergeErr7 != nil {
			return file, mergeErr7
		}

		mergeErr8 := file.MergeCell(v, "J5", "J6")
		if mergeErr8 != nil {
			return file, mergeErr8
		}

		var bulan string
		switch v {
		case "JAN":
			bulan = "Januari"
		case "FEB":
			bulan = "Februari"
		case "MAR":
			bulan = "Maret"
		case "APR":
			bulan = "April"
		case "MEI":
			bulan = "Mei"
		case "JUN":
			bulan = "Juni"
		case "JUL":
			bulan = "Juli"
		case "AGU":
			bulan = "Agustus"
		case "SEP":
			bulan = "September"
		case "OKT":
			bulan = "Oktober"
		case "NOV":
			bulan = "November"
		case "DES":
			bulan = "Desember"
		}

		file.SetColWidth(v, "A", "B", float64(3))
		file.SetColWidth(v, "I", "I", float64(3))

		file.SetColWidth(v, "C", "C", float64(15))
		file.SetColWidth(v, "F", "F", float64(20))

		file.SetColWidth(v, "D", "E", float64(40))
		file.SetColWidth(v, "G", "H", float64(15))

		file.SetColWidth(v, "J", "J", float64(10))

		file.SetCellValue(v, "B2", "REALISASI PEMENUHAN KEBUTUHAN BATUBARA DALAM NEGERI PT ANGSANA JAYA ENERGI")
		file.SetCellValue(v, "B3", fmt.Sprintf("Bulan %s Tahun %s", bulan, year))
		file.SetCellStyle(v, "B2", "H3", boldStyleCenter)
		file.SetCellStyle(v, "J5", "J6", boldStyleCenter)
		file.SetCellStyle(v, "B5", "H6", boldStyleBorder)

		file.SetCellValue(v, "B5", "No")
		file.SetCellValue(v, "C5", "Tanggal")
		file.SetCellValue(v, "D5", "Pembeli")
		file.SetCellValue(v, "D6", "Trader")
		file.SetCellValue(v, "E6", "End User")
		file.SetCellValue(v, "F6", "Bidang Usaha")
		file.SetCellValue(v, "G5", "Kalori CV (Gar)")
		file.SetCellValue(v, "H5", "Jumlah (Ton)")
		file.SetCellValue(v, "J5", "Berita Acara")
	}

	return file, nil
}

func (s *service) RecapDmo(year string, iupopkId int) (ReportDmoOutput, error) {
	recapDmo, recapDmoErr := s.repository.RecapDmo(year, iupopkId)

	return recapDmo, recapDmoErr
}

func (s *service) RealizationReport(year string, iupopkId int) (RealizationOutput, error) {
	realizationReport, realizationReportErr := s.repository.RealizationReport(year, iupopkId)

	return realizationReport, realizationReportErr
}

func (s *service) SaleDetailReport(year string, iupopkId int) (SaleDetail, error) {
	saleDetailReport, saleDetailReportErr := s.repository.SaleDetailReport(year, iupopkId)

	return saleDetailReport, saleDetailReportErr
}

func (s *service) CreateReportRecapDmo(year string, reportRecapDmo ReportDmoOutput, iupopk iupopk.Iupopk, file *excelize.File, sheetName string) (*excelize.File, error) {

	file.SetCellValue(sheetName, "A1", iupopk.Name)
	file.SetCellValue(sheetName, "A2", "Rekap DMO")

	period := fmt.Sprintf("Januari - Desember %s", year)

	file.SetCellValue(sheetName, "A3", period)

	mergeErr1 := file.MergeCell(sheetName, "A1", "C1")
	if mergeErr1 != nil {
		return file, mergeErr1
	}
	mergeErr2 := file.MergeCell(sheetName, "A2", "C2")
	if mergeErr2 != nil {
		return file, mergeErr2
	}
	mergeErr3 := file.MergeCell(sheetName, "A3", "C3")
	if mergeErr3 != nil {
		return file, mergeErr3
	}

	categories := map[string]string{
		"A5":  "Bulan",
		"A6":  "Januari",
		"A7":  "Februari",
		"A8":  "Maret",
		"A9":  "April",
		"A10": "Mei",
		"A11": "Juni",
		"A12": "Juli",
		"A13": "Agustus",
		"A14": "September",
		"A15": "Oktober",
		"A16": "November",
		"A17": "Desember",
		"A18": "TOTAL",
		"B5":  "Kelistrikan",
		"C5":  "Non Kelistrikan",
		"D5":  "Jumlah",
		"F5":  "Produksi",
		"G5":  "Tidak bisa Claim DMO",
	}

	for k, v := range categories {
		file.SetCellValue(sheetName, k, v)
	}

	recapElectricity := reflect.ValueOf(reportRecapDmo.RecapElectricity)

	values := make(map[string]interface{})

	for i := 0; i < recapElectricity.NumField(); i++ {
		switch recapElectricity.Type().Field(i).Name {
		case "January":
			if recapElectricity.Field(i).Interface().(float64) > 0 {
				values["B6"] = recapElectricity.Field(i).Interface().(float64)
				values["D6"] = recapElectricity.Field(i).Interface().(float64)
			} else {
				values["B6"] = "-"
			}
		case "February":
			if recapElectricity.Field(i).Interface().(float64) > 0 {
				values["B7"] = recapElectricity.Field(i).Interface().(float64)
				values["D7"] = recapElectricity.Field(i).Interface().(float64)
			} else {
				values["B7"] = "-"
			}
		case "March":
			if recapElectricity.Field(i).Interface().(float64) > 0 {
				values["B8"] = recapElectricity.Field(i).Interface().(float64)
				values["D8"] = recapElectricity.Field(i).Interface().(float64)
			} else {
				values["B8"] = "-"
			}
		case "April":
			if recapElectricity.Field(i).Interface().(float64) > 0 {
				values["B9"] = recapElectricity.Field(i).Interface().(float64)
				values["D9"] = recapElectricity.Field(i).Interface().(float64)
			} else {
				values["B9"] = "-"
			}
		case "May":
			if recapElectricity.Field(i).Interface().(float64) > 0 {
				values["B10"] = recapElectricity.Field(i).Interface().(float64)
				values["D10"] = recapElectricity.Field(i).Interface().(float64)
			} else {
				values["B10"] = "-"
			}
		case "June":
			if recapElectricity.Field(i).Interface().(float64) > 0 {
				values["B11"] = recapElectricity.Field(i).Interface().(float64)
				values["D11"] = recapElectricity.Field(i).Interface().(float64)
			} else {
				values["B11"] = "-"
			}
		case "July":
			if recapElectricity.Field(i).Interface().(float64) > 0 {
				values["B12"] = recapElectricity.Field(i).Interface().(float64)
				values["D12"] = recapElectricity.Field(i).Interface().(float64)
			} else {
				values["B12"] = "-"
			}
		case "August":
			if recapElectricity.Field(i).Interface().(float64) > 0 {
				values["B13"] = recapElectricity.Field(i).Interface().(float64)
				values["D13"] = recapElectricity.Field(i).Interface().(float64)
			} else {
				values["B13"] = "-"
			}
		case "September":
			if recapElectricity.Field(i).Interface().(float64) > 0 {
				values["B14"] = recapElectricity.Field(i).Interface().(float64)
				values["D14"] = recapElectricity.Field(i).Interface().(float64)
			} else {
				values["B14"] = "-"
			}
		case "October":
			if recapElectricity.Field(i).Interface().(float64) > 0 {
				values["B15"] = recapElectricity.Field(i).Interface().(float64)
				values["D15"] = recapElectricity.Field(i).Interface().(float64)
			} else {
				values["B15"] = "-"
			}
		case "November":
			if recapElectricity.Field(i).Interface().(float64) > 0 {
				values["B16"] = recapElectricity.Field(i).Interface().(float64)
				values["D16"] = recapElectricity.Field(i).Interface().(float64)
			} else {
				values["B16"] = "-"
			}
		case "December":
			if recapElectricity.Field(i).Interface().(float64) > 0 {
				values["B17"] = recapElectricity.Field(i).Interface().(float64)
				values["D17"] = recapElectricity.Field(i).Interface().(float64)
			} else {
				values["B17"] = "-"
			}
		case "Total":
			if recapElectricity.Field(i).Interface().(float64) > 0 {
				values["B18"] = recapElectricity.Field(i).Interface().(float64)
				values["D18"] = recapElectricity.Field(i).Interface().(float64)
			} else {
				values["B18"] = "-"
			}
		}
	}

	recapNonElectricity := reflect.ValueOf(reportRecapDmo.RecapNonElectricity)

	for i := 0; i < recapNonElectricity.NumField(); i++ {
		switch recapNonElectricity.Type().Field(i).Name {
		case "January":
			if recapNonElectricity.Field(i).Interface().(float64) > 0 {
				values["C6"] = recapNonElectricity.Field(i).Interface().(float64)
				values["D6"] = values["D6"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
			} else {
				if values["D6"] == nil {
					values["D6"] = "-"
				}
				values["C6"] = "-"
			}
		case "February":
			if recapNonElectricity.Field(i).Interface().(float64) > 0 {
				values["C7"] = recapNonElectricity.Field(i).Interface().(float64)
				values["D7"] = values["D7"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
			} else {
				if values["D7"] == nil {
					values["D7"] = "-"
				}
				values["C7"] = "-"
			}
		case "March":
			if recapNonElectricity.Field(i).Interface().(float64) > 0 {
				values["C8"] = recapNonElectricity.Field(i).Interface().(float64)
				values["D8"] = values["D8"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
			} else {
				if values["D8"] == nil {
					values["D8"] = "-"
				}
				values["C8"] = "-"
			}
		case "April":
			if recapNonElectricity.Field(i).Interface().(float64) > 0 {
				values["C9"] = recapNonElectricity.Field(i).Interface().(float64)
				values["D9"] = values["D9"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
			} else {
				if values["D9"] == nil {
					values["D9"] = "-"
				}
				values["C9"] = "-"
			}
		case "May":
			if recapNonElectricity.Field(i).Interface().(float64) > 0 {
				values["C10"] = recapNonElectricity.Field(i).Interface().(float64)
				values["D10"] = values["D10"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
			} else {
				if values["D10"] == nil {
					values["D10"] = "-"
				}
				values["C10"] = "-"
			}
		case "June":
			if recapNonElectricity.Field(i).Interface().(float64) > 0 {
				values["C11"] = recapNonElectricity.Field(i).Interface().(float64)
				values["D11"] = values["D11"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
			} else {
				if values["D11"] == nil {
					values["D11"] = "-"
				}
				values["C11"] = "-"
			}
		case "July":
			if recapNonElectricity.Field(i).Interface().(float64) > 0 {
				values["C12"] = recapNonElectricity.Field(i).Interface().(float64)
				values["D12"] = values["D12"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
			} else {
				if values["D12"] == nil {
					values["D12"] = "-"
				}
				values["C12"] = "-"
			}
		case "August":
			if recapNonElectricity.Field(i).Interface().(float64) > 0 {
				values["C13"] = recapNonElectricity.Field(i).Interface().(float64)
				values["D13"] = values["D13"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
			} else {
				if values["D13"] == nil {
					values["D13"] = "-"
				}
				values["C13"] = "-"
			}
		case "September":
			if recapNonElectricity.Field(i).Interface().(float64) > 0 {
				values["C14"] = recapNonElectricity.Field(i).Interface().(float64)
				values["D14"] = values["D14"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
			} else {
				if values["D14"] == nil {
					values["D14"] = "-"
				}
				values["C14"] = "-"
			}
		case "October":
			if recapNonElectricity.Field(i).Interface().(float64) > 0 {
				values["C15"] = recapNonElectricity.Field(i).Interface().(float64)
				values["D15"] = values["D15"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
			} else {
				if values["D15"] == nil {
					values["D15"] = "-"
				}
				values["C15"] = "-"
			}
		case "November":
			if recapNonElectricity.Field(i).Interface().(float64) > 0 {
				values["C16"] = recapNonElectricity.Field(i).Interface().(float64)
				values["D16"] = values["D16"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
			} else {
				if values["D16"] == nil {
					values["D16"] = "-"
				}
				values["C16"] = "-"
			}
		case "December":
			if recapNonElectricity.Field(i).Interface().(float64) > 0 {
				values["C17"] = recapNonElectricity.Field(i).Interface().(float64)
				values["D17"] = values["D17"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
			} else {
				if values["D17"] == nil {
					values["D17"] = "-"
				}
				values["C17"] = "-"
			}
		case "Total":
			if recapNonElectricity.Field(i).Interface().(float64) > 0 {
				values["C18"] = recapNonElectricity.Field(i).Interface().(float64)
				values["D18"] = values["D18"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
			} else {
				if values["D18"] == nil {
					values["D18"] = "-"
				}
				values["C18"] = "-"
			}
		}
	}

	recapProduction := reflect.ValueOf(reportRecapDmo.Production)

	for i := 0; i < recapProduction.NumField(); i++ {
		switch recapProduction.Type().Field(i).Name {
		case "January":
			if recapProduction.Field(i).Interface().(float64) > 0 {
				values["F6"] = recapProduction.Field(i).Interface().(float64)
			} else {
				values["F6"] = "-"
			}
		case "February":
			if recapProduction.Field(i).Interface().(float64) > 0 {
				values["F7"] = recapProduction.Field(i).Interface().(float64)
			} else {
				values["F7"] = "-"
			}
		case "March":
			if recapProduction.Field(i).Interface().(float64) > 0 {
				values["F8"] = recapProduction.Field(i).Interface().(float64)
			} else {
				values["F8"] = "-"
			}
		case "April":
			if recapProduction.Field(i).Interface().(float64) > 0 {
				values["F9"] = recapProduction.Field(i).Interface().(float64)
			} else {
				values["F9"] = "-"
			}
		case "May":
			if recapProduction.Field(i).Interface().(float64) > 0 {
				values["F10"] = recapProduction.Field(i).Interface().(float64)
			} else {
				values["F10"] = "-"
			}
		case "June":
			if recapProduction.Field(i).Interface().(float64) > 0 {
				values["F11"] = recapProduction.Field(i).Interface().(float64)
			} else {
				values["F11"] = "-"
			}
		case "July":
			if recapProduction.Field(i).Interface().(float64) > 0 {
				values["F12"] = recapProduction.Field(i).Interface().(float64)
			} else {
				values["F12"] = "-"
			}
		case "August":
			if recapProduction.Field(i).Interface().(float64) > 0 {
				values["F13"] = recapProduction.Field(i).Interface().(float64)
			} else {
				values["F13"] = "-"
			}
		case "September":
			if recapProduction.Field(i).Interface().(float64) > 0 {
				values["F14"] = recapProduction.Field(i).Interface().(float64)
			} else {
				values["F14"] = "-"
			}
		case "October":
			if recapProduction.Field(i).Interface().(float64) > 0 {
				values["F15"] = recapProduction.Field(i).Interface().(float64)
			} else {
				values["F15"] = "-"
			}
		case "November":
			if recapProduction.Field(i).Interface().(float64) > 0 {
				values["F16"] = recapProduction.Field(i).Interface().(float64)
			} else {
				values["F16"] = "-"
			}
		case "December":
			if recapProduction.Field(i).Interface().(float64) > 0 {
				values["F17"] = recapProduction.Field(i).Interface().(float64)
			} else {
				values["F17"] = "-"
			}
		case "Total":
			if recapProduction.Field(i).Interface().(float64) > 0 {
				values["F18"] = recapProduction.Field(i).Interface().(float64)
			} else {
				values["F18"] = "-"
			}
		}
	}

	recapNotClaim := reflect.ValueOf(reportRecapDmo.NotClaimable)

	for i := 0; i < recapNotClaim.NumField(); i++ {
		switch recapNotClaim.Type().Field(i).Name {
		case "January":
			if recapNotClaim.Field(i).Interface().(float64) > 0 {
				values["G6"] = recapNotClaim.Field(i).Interface().(float64)
			} else {
				values["G6"] = "-"
			}
		case "February":
			if recapNotClaim.Field(i).Interface().(float64) > 0 {
				values["G7"] = recapNotClaim.Field(i).Interface().(float64)
			} else {
				values["G7"] = "-"
			}
		case "March":
			if recapNotClaim.Field(i).Interface().(float64) > 0 {
				values["G8"] = recapNotClaim.Field(i).Interface().(float64)
			} else {
				values["G8"] = "-"
			}
		case "April":
			if recapNotClaim.Field(i).Interface().(float64) > 0 {
				values["G9"] = recapNotClaim.Field(i).Interface().(float64)
			} else {
				values["G9"] = "-"
			}
		case "May":
			if recapNotClaim.Field(i).Interface().(float64) > 0 {
				values["G10"] = recapNotClaim.Field(i).Interface().(float64)
			} else {
				values["G10"] = "-"
			}
		case "June":
			if recapNotClaim.Field(i).Interface().(float64) > 0 {
				values["G11"] = recapNotClaim.Field(i).Interface().(float64)
			} else {
				values["G11"] = "-"
			}
		case "July":
			if recapNotClaim.Field(i).Interface().(float64) > 0 {
				values["G12"] = recapNotClaim.Field(i).Interface().(float64)
			} else {
				values["G12"] = "-"
			}
		case "August":
			if recapNotClaim.Field(i).Interface().(float64) > 0 {
				values["G13"] = recapNotClaim.Field(i).Interface().(float64)
			} else {
				values["G13"] = "-"
			}
		case "September":
			if recapNotClaim.Field(i).Interface().(float64) > 0 {
				values["G14"] = recapNotClaim.Field(i).Interface().(float64)
			} else {
				values["G14"] = "-"
			}
		case "October":
			if recapNotClaim.Field(i).Interface().(float64) > 0 {
				values["G15"] = recapNotClaim.Field(i).Interface().(float64)
			} else {
				values["G15"] = "-"
			}
		case "November":
			if recapNotClaim.Field(i).Interface().(float64) > 0 {
				values["G16"] = recapNotClaim.Field(i).Interface().(float64)
			} else {
				values["G16"] = "-"
			}
		case "December":
			if recapNotClaim.Field(i).Interface().(float64) > 0 {
				values["G17"] = recapNotClaim.Field(i).Interface().(float64)
			} else {
				values["G17"] = "-"
			}
		case "Total":
			if recapNotClaim.Field(i).Interface().(float64) > 0 {
				values["G18"] = recapNotClaim.Field(i).Interface().(float64)
			} else {
				values["G18"] = "-"
			}
		}
	}

	for k, v := range values {
		file.SetCellValue(sheetName, k, v)
	}
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
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})

	fontTitleBoldStyle, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
		},
	})

	borderStyle, _ := file.NewStyle(&excelize.Style{
		Border:       border,
		CustomNumFmt: &custFmt,
	})

	boldNumberStyle, _ := file.NewStyle(&excelize.Style{
		Border: border,
		Font: &excelize.Font{
			Bold: true,
		},
		CustomNumFmt: &custFmt,
	})

	boldNumberOnlyStyle, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		CustomNumFmt: &custFmt,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})

	customNumberRightStyle, _ := file.NewStyle(&excelize.Style{
		Border:       border,
		CustomNumFmt: &custFmt,
		Alignment: &excelize.Alignment{
			Horizontal: "right",
		},
	})

	errStyleFontTitle := file.SetCellStyle(sheetName, "A1", "A3", fontTitleBoldStyle)

	if errStyleFontTitle != nil {
		return file, errStyleFontTitle
	}

	errStyleBorderBold1 := file.SetCellStyle(sheetName, "A5", "D5", boldStyle)

	if errStyleBorderBold1 != nil {
		return file, errStyleBorderBold1
	}

	errStyleBorderBold2 := file.SetCellStyle(sheetName, "F5", "G5", boldStyle)

	if errStyleBorderBold2 != nil {
		return file, errStyleBorderBold2
	}

	errFmtNmbr1 := file.SetCellStyle(sheetName, "A6", "D17", borderStyle)

	if errFmtNmbr1 != nil {
		return file, errFmtNmbr1
	}

	errFmtNmbr2 := file.SetCellStyle(sheetName, "F6", "G17", borderStyle)

	if errFmtNmbr2 != nil {
		return file, errFmtNmbr2
	}

	errBoldNmbr1 := file.SetCellStyle(sheetName, "A18", "D18", boldNumberStyle)

	if errBoldNmbr1 != nil {
		return file, errBoldNmbr1
	}

	errBoldNmbr2 := file.SetCellStyle(sheetName, "F18", "G18", boldNumberStyle)

	if errBoldNmbr2 != nil {
		return file, errBoldNmbr2
	}

	errCustomNmbrRight1 := file.SetCellStyle(sheetName, "B6", "D17", customNumberRightStyle)

	if errCustomNmbrRight1 != nil {
		return file, errCustomNmbrRight1
	}

	errCustomNmbrRight2 := file.SetCellStyle(sheetName, "F6", "G17", customNumberRightStyle)

	if errCustomNmbrRight2 != nil {
		return file, errCustomNmbrRight2
	}

	file.SetColWidth(sheetName, "A", "A", float64(15))
	file.SetColWidth(sheetName, "B", "D", float64(25))
	file.SetColWidth(sheetName, "E", "E", float64(5))
	file.SetColWidth(sheetName, "F", "F", float64(25))
	file.SetColWidth(sheetName, "G", "G", float64(30))

	mergeErr4 := file.MergeCell(sheetName, "A20", "C20")
	if mergeErr4 != nil {
		return file, mergeErr4
	}

	file.SetCellValue(sheetName, "A20", "% Pemenuhan DMO terhadap REALISASI PRODUKSI")

	file.SetCellFormula(sheetName, "D20", "D18/F18")

	percentageBoldStyle, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		NumFmt: 10,
	})

	errPercentageBoldStyle := file.SetCellStyle(sheetName, "A20", "G60", percentageBoldStyle)
	errBoldNmbrOnly1 := file.SetCellStyle(sheetName, "F20", "F60", boldNumberOnlyStyle)

	if errPercentageBoldStyle != nil {
		return file, errPercentageBoldStyle
	}

	if errBoldNmbrOnly1 != nil {
		return file, errBoldNmbrOnly1
	}

	var color []string

	color = append(color, "#FFFF01")

	normalStyle, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: false,
		},
		Fill: excelize.Fill{Color: color},
	})

	goment.SetLocale("id")

	for idx, rkab := range reportRecapDmo.Rkabs {

		dateFormat, errDate := goment.New(rkab.DateOfIssue)
		if errDate != nil {
			return file, errDate
		}

		file.SetCellValue(sheetName, fmt.Sprintf("A%v", 23+(8*idx)), fmt.Sprintf("Rencana Produksi	Disetujui Tanggal %s", dateFormat.Format("DD MMMM YYYY", "id")))

		numberCommas := message.NewPrinter(language.Indonesian)

		withCommaThousandSep := numberCommas.Sprintf("RKAB: %.0f", rkab.ProductionQuota)
		file.SetCellValue(sheetName, fmt.Sprintf("D%v", 23+(8*idx)), withCommaThousandSep)
		file.SetCellStyle(sheetName, fmt.Sprintf("A%v", 23+(8*idx)), fmt.Sprintf("G%v", 23+(8*idx)), normalStyle)

		file.SetCellValue(sheetName, fmt.Sprintf("A%v", 24+(8*idx)), "% Pemenuhan DMO terhadap RENCANA PRODUKSI")
		file.SetCellFormula(sheetName, fmt.Sprintf("D%v", 24+(8*idx)), fmt.Sprintf("D18/F%v", 24+(8*idx)))
		file.SetCellValue(sheetName, fmt.Sprintf("F%v", 24+(8*idx)), rkab.ProductionQuota)
		file.SetCellValue(sheetName, fmt.Sprintf("G%v", 24+(8*idx)), "Quota RKAB")

		file.SetCellValue(sheetName, fmt.Sprintf("A%v", 25+(8*idx)), fmt.Sprintf("disetujui tgl %s", dateFormat.Format("DD MMMM YYYY", "id")))

		file.SetCellValue(sheetName, fmt.Sprintf("A%v", 27+(8*idx)), fmt.Sprintf("%% Pemenuhan DMO terhadap kewajiban pemenuhan DMO %.0f%%", rkab.DmoObligation))
		file.SetCellFormula(sheetName, fmt.Sprintf("D%v", 27+(8*idx)), fmt.Sprintf("D%v/%.2f%%", 24+(8*idx), rkab.DmoObligation))

		file.SetCellValue(sheetName, fmt.Sprintf("A%v", 29+(8*idx)), "% Pemenuhan DMO terhadap Rencana Produksi (Prorata 12 bulan)")
		file.SetCellFormula(sheetName, fmt.Sprintf("D%v", 29+(8*idx)), fmt.Sprintf("D18/F%v", 29+(8*idx)))
		file.SetCellValue(sheetName, fmt.Sprintf("F%v", 29+(8*idx)), rkab.ProductionQuota)
		file.SetCellValue(sheetName, fmt.Sprintf("G%v", 29+(8*idx)), "prorata 12 bulan")

	}

	return file, nil
}

func (s *service) CreateReportRealization(year string, reportRealization RealizationOutput, iupopk iupopk.Iupopk, file *excelize.File) (*excelize.File, error) {
	var sheetName []string

	sheetName = append(sheetName,
		"JAN",
		"FEB",
		"MAR",
		"APR",
		"MEI",
		"JUN",
		"JUL",
		"AGU",
		"SEP",
		"OKT",
		"NOV",
		"DES",
	)

	newFile, err := createTemplateRealization(file, sheetName, year)

	if err != nil {
		return file, err
	}
	// for _, v := range sheetName {

	// 	switch v {
	// 	case "JAN":

	// 	case "FEB":

	// 	case "MAR":

	// 	case "APR":

	// 	case "MEI":

	// 	case "JUN":

	// 	case "JUL":

	// 	case "AGU":

	// 	case "SEP":

	// 	case "OKT":

	// 	case "NOV":

	// 	case "DES":

	// 	}

	// }

	return newFile, nil
}
