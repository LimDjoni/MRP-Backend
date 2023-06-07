package masterreport

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

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
	CreateReportSalesDetail(year string, reportSaleDetail SaleDetail, iupopk iupopk.Iupopk, file *excelize.File, sheetName string, chartSheetName string) (*excelize.File, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func createTemplateRealization(file *excelize.File, sheetName []string, year string, iupopkName string) (*excelize.File, error) {
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

		file.SetCellValue(v, "B2", fmt.Sprintf("REALISASI PEMENUHAN KEBUTUHAN BATUBARA DALAM NEGERI %s", strings.ToUpper(iupopkName)))
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

func insertTransactionRealization(file *excelize.File, sheetName string, transactionElectric []RealizationTransaction, transactionNonElectric []RealizationTransaction) (*excelize.File, error) {
	custFmt := "_(* #,##0_);_(* \\(#,##0\\);_(* \"-\"_);_(@_)"

	custDateFmt := "d-mmm-yyyy"

	border := []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}

	borderStyle, _ := file.NewStyle(&excelize.Style{
		Border:       border,
		CustomNumFmt: &custFmt,
	})

	borderCenterStyle, _ := file.NewStyle(&excelize.Style{
		Border: border,
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})

	borderCenterDateStyle, _ := file.NewStyle(&excelize.Style{
		Border: border,
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
		CustomNumFmt: &custDateFmt,
	})

	centerStyle, _ := file.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "center",
		},
	})

	file.SetCellStyle(sheetName, "B7", fmt.Sprintf("F%v", 7+len(transactionElectric)+len(transactionNonElectric)), borderCenterStyle)

	file.SetCellStyle(sheetName, "C7", fmt.Sprintf("C%v", 7+len(transactionElectric)+len(transactionNonElectric)), borderCenterDateStyle)

	file.SetCellStyle(sheetName, "G7", fmt.Sprintf("H%v", 7+len(transactionElectric)+len(transactionNonElectric)), borderStyle)

	file.SetCellStyle(sheetName, "J7", fmt.Sprintf("J%v", 7+len(transactionElectric)+len(transactionNonElectric)), centerStyle)

	for idx, nonElectric := range transactionNonElectric {
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", idx+7), idx+1)
		date, err := time.Parse("2006-01-02", nonElectric.ShippingDate)

		if err != nil {
			return file, err
		}

		file.SetCellValue(sheetName, fmt.Sprintf("C%v", idx+7), date)
		if nonElectric.Trader != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", idx+7), strings.ToUpper(nonElectric.Trader.CompanyName))
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", idx+7), "-")
		}

		if nonElectric.EndUser != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("E%v", idx+7), strings.ToUpper(nonElectric.EndUser.CompanyName))
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("E%v", idx+7), "-")
		}

		if nonElectric.EndUser != nil {
			if nonElectric.EndUser.IndustryType != nil {
				if nonElectric.EndUser.IndustryType.Category == "NON ELECTRICITY" {
					file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7), "NON KELISTRIKAN")
				}
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7), "-")
			}
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7), "-")
		}

		file.SetCellValue(sheetName, fmt.Sprintf("G%v", idx+7), nonElectric.QualityCaloriesAr)
		file.SetCellValue(sheetName, fmt.Sprintf("H%v", idx+7), nonElectric.Quantity)

		if nonElectric.IsBastOk {
			file.SetCellValue(sheetName, fmt.Sprintf("J%v", idx+7), "OK")
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("J%v", idx+7), "-")
		}
	}

	for idx, electric := range transactionElectric {
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", idx+7+len(transactionNonElectric)), idx+1+len(transactionNonElectric))
		date, err := time.Parse("2006-01-02", electric.ShippingDate)

		if err != nil {
			return file, err
		}
		file.SetCellValue(sheetName, fmt.Sprintf("C%v", idx+7+len(transactionNonElectric)), date)
		if electric.Trader != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", idx+7+len(transactionNonElectric)), strings.ToUpper(electric.Trader.CompanyName))
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", idx+7+len(transactionNonElectric)), "-")
		}

		if electric.EndUser != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("E%v", idx+7+len(transactionNonElectric)), strings.ToUpper(electric.EndUser.CompanyName))
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("E%v", idx+7+len(transactionNonElectric)), "-")
		}

		if electric.EndUser != nil {
			if electric.EndUser.IndustryType != nil {
				if electric.EndUser.IndustryType.Category == "ELECTRICITY" {
					file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)), "KELISTRIKAN")
				} else if electric.EndUser.IndustryType.Category == "NON ELECTRICITY" {
					file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)), "NON KELISTRIKAN")
				}
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)), "-")
			}
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)), "-")
		}

		file.SetCellValue(sheetName, fmt.Sprintf("G%v", idx+7+len(transactionNonElectric)), electric.QualityCaloriesAr)
		file.SetCellValue(sheetName, fmt.Sprintf("H%v", idx+7+len(transactionNonElectric)), electric.Quantity)

		if electric.IsBastOk {
			file.SetCellValue(sheetName, fmt.Sprintf("J%v", idx+7+len(transactionNonElectric)), "OK")
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("J%v", idx+7+len(transactionNonElectric)), "-")
		}
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
			values["B6"] = recapElectricity.Field(i).Interface().(float64)
			values["D6"] = recapElectricity.Field(i).Interface().(float64)
		case "February":
			values["B7"] = recapElectricity.Field(i).Interface().(float64)
			values["D7"] = recapElectricity.Field(i).Interface().(float64)

		case "March":
			values["B8"] = recapElectricity.Field(i).Interface().(float64)
			values["D8"] = recapElectricity.Field(i).Interface().(float64)

		case "April":
			values["B9"] = recapElectricity.Field(i).Interface().(float64)
			values["D9"] = recapElectricity.Field(i).Interface().(float64)

		case "May":
			values["B10"] = recapElectricity.Field(i).Interface().(float64)
			values["D10"] = recapElectricity.Field(i).Interface().(float64)

		case "June":
			values["B11"] = recapElectricity.Field(i).Interface().(float64)
			values["D11"] = recapElectricity.Field(i).Interface().(float64)

		case "July":
			values["B12"] = recapElectricity.Field(i).Interface().(float64)
			values["D12"] = recapElectricity.Field(i).Interface().(float64)

		case "August":
			values["B13"] = recapElectricity.Field(i).Interface().(float64)
			values["D13"] = recapElectricity.Field(i).Interface().(float64)

		case "September":
			values["B14"] = recapElectricity.Field(i).Interface().(float64)
			values["D14"] = recapElectricity.Field(i).Interface().(float64)

		case "October":
			values["B15"] = recapElectricity.Field(i).Interface().(float64)
			values["D15"] = recapElectricity.Field(i).Interface().(float64)

		case "November":
			values["B16"] = recapElectricity.Field(i).Interface().(float64)
			values["D16"] = recapElectricity.Field(i).Interface().(float64)

		case "December":
			values["B17"] = recapElectricity.Field(i).Interface().(float64)
			values["D17"] = recapElectricity.Field(i).Interface().(float64)

		case "Total":
			values["B18"] = recapElectricity.Field(i).Interface().(float64)
			values["D18"] = recapElectricity.Field(i).Interface().(float64)

		}
	}

	recapNonElectricity := reflect.ValueOf(reportRecapDmo.RecapNonElectricity)

	for i := 0; i < recapNonElectricity.NumField(); i++ {
		switch recapNonElectricity.Type().Field(i).Name {
		case "January":
			values["C6"] = recapNonElectricity.Field(i).Interface().(float64)
			values["D6"] = values["D6"].(float64) + recapNonElectricity.Field(i).Interface().(float64)

		case "February":
			values["C7"] = recapNonElectricity.Field(i).Interface().(float64)
			values["D7"] = values["D7"].(float64) + recapNonElectricity.Field(i).Interface().(float64)

		case "March":
			values["C8"] = recapNonElectricity.Field(i).Interface().(float64)
			values["D8"] = values["D8"].(float64) + recapNonElectricity.Field(i).Interface().(float64)

		case "April":
			values["C9"] = recapNonElectricity.Field(i).Interface().(float64)
			values["D9"] = values["D9"].(float64) + recapNonElectricity.Field(i).Interface().(float64)

		case "May":
			values["C10"] = recapNonElectricity.Field(i).Interface().(float64)
			values["D10"] = values["D10"].(float64) + recapNonElectricity.Field(i).Interface().(float64)

		case "June":
			values["C11"] = recapNonElectricity.Field(i).Interface().(float64)
			values["D11"] = values["D11"].(float64) + recapNonElectricity.Field(i).Interface().(float64)

		case "July":
			values["C12"] = recapNonElectricity.Field(i).Interface().(float64)
			values["D12"] = values["D12"].(float64) + recapNonElectricity.Field(i).Interface().(float64)

		case "August":
			values["C13"] = recapNonElectricity.Field(i).Interface().(float64)
			values["D13"] = values["D13"].(float64) + recapNonElectricity.Field(i).Interface().(float64)

		case "September":
			values["C14"] = recapNonElectricity.Field(i).Interface().(float64)
			values["D14"] = values["D14"].(float64) + recapNonElectricity.Field(i).Interface().(float64)

		case "October":
			values["C15"] = recapNonElectricity.Field(i).Interface().(float64)
			values["D15"] = values["D15"].(float64) + recapNonElectricity.Field(i).Interface().(float64)

		case "November":
			values["C16"] = recapNonElectricity.Field(i).Interface().(float64)
			values["D16"] = values["D16"].(float64) + recapNonElectricity.Field(i).Interface().(float64)

		case "December":
			values["C17"] = recapNonElectricity.Field(i).Interface().(float64)
			values["D17"] = values["D17"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
		case "Total":
			values["C18"] = recapNonElectricity.Field(i).Interface().(float64)
			values["D18"] = values["D18"].(float64) + recapNonElectricity.Field(i).Interface().(float64)
		}
	}

	recapProduction := reflect.ValueOf(reportRecapDmo.Production)

	for i := 0; i < recapProduction.NumField(); i++ {
		switch recapProduction.Type().Field(i).Name {
		case "January":
			values["F6"] = recapProduction.Field(i).Interface().(float64)
		case "February":
			values["F7"] = recapProduction.Field(i).Interface().(float64)
		case "March":
			values["F8"] = recapProduction.Field(i).Interface().(float64)
		case "April":
			values["F9"] = recapProduction.Field(i).Interface().(float64)
		case "May":
			values["F10"] = recapProduction.Field(i).Interface().(float64)
		case "June":
			values["F11"] = recapProduction.Field(i).Interface().(float64)
		case "July":
			values["F12"] = recapProduction.Field(i).Interface().(float64)
		case "August":
			values["F13"] = recapProduction.Field(i).Interface().(float64)
		case "September":
			values["F14"] = recapProduction.Field(i).Interface().(float64)
		case "October":
			values["F15"] = recapProduction.Field(i).Interface().(float64)
		case "November":
			values["F16"] = recapProduction.Field(i).Interface().(float64)
		case "December":
			values["F17"] = recapProduction.Field(i).Interface().(float64)
		case "Total":
			values["F18"] = recapProduction.Field(i).Interface().(float64)
		}
	}

	recapNotClaim := reflect.ValueOf(reportRecapDmo.NotClaimable)

	for i := 0; i < recapNotClaim.NumField(); i++ {
		switch recapNotClaim.Type().Field(i).Name {
		case "January":
			values["G6"] = recapNotClaim.Field(i).Interface().(float64)
		case "February":
			values["G7"] = recapNotClaim.Field(i).Interface().(float64)
		case "March":
			values["G8"] = recapNotClaim.Field(i).Interface().(float64)
		case "April":
			values["G9"] = recapNotClaim.Field(i).Interface().(float64)
		case "May":
			values["G10"] = recapNotClaim.Field(i).Interface().(float64)
		case "June":
			values["G11"] = recapNotClaim.Field(i).Interface().(float64)
		case "July":
			values["G12"] = recapNotClaim.Field(i).Interface().(float64)
		case "August":
			values["G13"] = recapNotClaim.Field(i).Interface().(float64)
		case "September":
			values["G14"] = recapNotClaim.Field(i).Interface().(float64)
		case "October":
			values["G15"] = recapNotClaim.Field(i).Interface().(float64)
		case "November":
			values["G16"] = recapNotClaim.Field(i).Interface().(float64)
		case "December":
			values["G17"] = recapNotClaim.Field(i).Interface().(float64)
		case "Total":
			values["G18"] = recapNotClaim.Field(i).Interface().(float64)
		}
	}

	for k, v := range values {
		file.SetCellValue(sheetName, k, v)
	}
	custFmt := "_(* #,##0_);_(* \\(#,##0\\);_(* \"-\"_);_(@_)"

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

	t := time.Now()

	month := t.Month()
	yearNow := t.Year()
	yearString := strconv.Itoa(yearNow)

	var maxDmoObligationQuota float64

	for _, rkab := range reportRecapDmo.Rkabs {
		if maxDmoObligationQuota < rkab.ProductionQuota*rkab.DmoObligation/100 {
			maxDmoObligationQuota = rkab.ProductionQuota * rkab.DmoObligation / 100
		}
	}

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
		file.SetCellFormula(sheetName, fmt.Sprintf("D%v", 27+(8*idx)), fmt.Sprintf("D18/F%v", 27+(8*idx)))
		file.SetCellValue(sheetName, fmt.Sprintf("F%v", 27+(8*idx)), maxDmoObligationQuota)

		if yearString == year {
			file.SetCellValue(sheetName, fmt.Sprintf("A%v", 29+(8*idx)), fmt.Sprintf("%% Pemenuhan DMO terhadap Rencana Produksi (Prorata %v bulan)", int(month)))
			file.SetCellFormula(sheetName, fmt.Sprintf("D%v", 29+(8*idx)), fmt.Sprintf("D18/F%v", 29+(8*idx)))
			file.SetCellValue(sheetName, fmt.Sprintf("F%v", 29+(8*idx)), rkab.ProductionQuota*float64(int(month))/float64(12))
			file.SetCellValue(sheetName, fmt.Sprintf("G%v", 29+(8*idx)), fmt.Sprintf("prorata %v bulan", int(month)))
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("A%v", 29+(8*idx)), "% Pemenuhan DMO terhadap Rencana Produksi (Prorata 12 bulan)")
			file.SetCellFormula(sheetName, fmt.Sprintf("D%v", 29+(8*idx)), fmt.Sprintf("D18/F%v", 29+(8*idx)))
			file.SetCellValue(sheetName, fmt.Sprintf("F%v", 29+(8*idx)), rkab.ProductionQuota)
			file.SetCellValue(sheetName, fmt.Sprintf("G%v", 29+(8*idx)), "prorata 12 bulan")
		}

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

	newFile, err := createTemplateRealization(file, sheetName, year, iupopk.Name)

	if err != nil {
		return file, err
	}

	var insertedFile *excelize.File
	var errFile error

	insertedFile = newFile
	for _, v := range sheetName {

		switch v {
		case "JAN":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.January, reportRealization.NonElectric.January)

			if errFile != nil {
				return insertedFile, errFile
			}
		case "FEB":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.February, reportRealization.NonElectric.February)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "MAR":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.March, reportRealization.NonElectric.March)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "APR":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.April, reportRealization.NonElectric.April)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "MEI":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.May, reportRealization.NonElectric.May)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "JUN":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.June, reportRealization.NonElectric.June)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "JUL":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.July, reportRealization.NonElectric.July)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "AGU":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.August, reportRealization.NonElectric.August)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "SEP":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.September, reportRealization.NonElectric.September)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "OKT":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.October, reportRealization.NonElectric.October)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "NOV":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.November, reportRealization.NonElectric.November)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "DES":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.December, reportRealization.NonElectric.December)
			if errFile != nil {
				return insertedFile, errFile
			}
		}

	}

	return insertedFile, nil
}

func (s *service) CreateReportSalesDetail(year string, reportSaleDetail SaleDetail, iupopk iupopk.Iupopk, file *excelize.File, sheetName string, chartSheetName string) (*excelize.File, error) {

	custFmt := "_(* #,##0_);_(* \\(#,##0\\);_(* \"-\"_);_(@_)"
	custMtFmt := "#,##0\\ \"mt\""

	custDateFmt := "d-mmm-yyyy"
	border := []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}

	mtFmtStyle, _ := file.NewStyle(&excelize.Style{
		CustomNumFmt: &custMtFmt,
		Alignment: &excelize.Alignment{
			Horizontal: "left",
		},
	})

	boldTitleStyle, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
		},
	})

	dateStyle, _ := file.NewStyle(&excelize.Style{
		CustomNumFmt: &custDateFmt,
	})

	boldOnlyStyle, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
	})

	percentTableStyle, _ := file.NewStyle(&excelize.Style{
		Border: border,
		NumFmt: 10,
	})

	boldPercentStyle, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		NumFmt: 10,
	})

	borderStyle, _ := file.NewStyle(&excelize.Style{
		Border: border,
	})

	centerStyle, _ := file.NewStyle(&excelize.Style{
		Border: border,
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	boldTitleTableStyle, _ := file.NewStyle(&excelize.Style{
		Border: border,
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	boldNumberStyle, _ := file.NewStyle(&excelize.Style{
		Border: border,
		Font: &excelize.Font{
			Bold: true,
		},
		CustomNumFmt: &custFmt,
		Alignment: &excelize.Alignment{
			Vertical: "center",
		},
	})

	formatNumberStyle, _ := file.NewStyle(&excelize.Style{
		Border:       border,
		CustomNumFmt: &custFmt,
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
	})

	var monthString []string

	monthString = append(monthString,
		"January",
		"February",
		"March",
		"April",
		"May",
		"June",
		"July",
		"August",
		"September",
		"October",
		"November",
		"December",
	)

	var maxProductionQuota float64

	file.SetColWidth(sheetName, "A", "A", float64(5))
	file.SetColWidth(sheetName, "B", "B", float64(25))
	file.SetColWidth(sheetName, "C", "C", float64(2))
	file.SetColWidth(sheetName, "D", "D", float64(25))
	file.SetColWidth(sheetName, "E", "E", float64(20))
	file.SetColWidth(sheetName, "F", "Q", float64(16))
	file.SetColWidth(sheetName, "H", "H", float64(20))

	file.SetCellValue(sheetName, "A1", iupopk.Name)

	file.SetColWidth(chartSheetName, "A", "A", float64(5))
	file.SetColWidth(chartSheetName, "B", "B", float64(25))
	file.SetColWidth(chartSheetName, "C", "E", float64(20))

	file.SetCellValue(chartSheetName, "A1", iupopk.Name)
	errTitleIupopk := file.SetCellStyle(sheetName, "A1", "A1", boldTitleStyle)
	if errTitleIupopk != nil {
		return file, errTitleIupopk
	}

	errChartTitleIupopk := file.SetCellStyle(chartSheetName, "A1", "A1", boldTitleStyle)
	if errChartTitleIupopk != nil {
		return file, errChartTitleIupopk
	}

	errYearStyle := file.SetCellStyle(sheetName, "A2", "A2", boldOnlyStyle)
	if errYearStyle != nil {
		return file, errYearStyle
	}

	errChartYearStyle := file.SetCellStyle(chartSheetName, "A2", "A2", boldOnlyStyle)
	if errChartYearStyle != nil {
		return file, errChartYearStyle
	}

	file.SetCellValue(chartSheetName, "A2", fmt.Sprintf("Tahun %s", year))

	file.SetCellValue(sheetName, "A2", fmt.Sprintf("Tahun %s", year))

	file.SetCellValue(sheetName, "A5", "RKAB")
	errTitleRkab := file.SetCellStyle(sheetName, "A5", "A5", boldTitleStyle)

	if errTitleRkab != nil {
		return file, errTitleRkab
	}
	goment.SetLocale("id")

	for idx, v := range reportSaleDetail.Rkabs {

		if maxProductionQuota < v.ProductionQuota {
			maxProductionQuota = v.ProductionQuota
		}

		file.SetCellValue(sheetName, fmt.Sprintf("B%v", 6+(idx*4)), "No. Surat")
		file.SetCellValue(sheetName, fmt.Sprintf("C%v", 6+(idx*4)), ":")
		file.SetCellValue(sheetName, fmt.Sprintf("D%v", 6+(idx*4)), v.LetterNumber)

		dateFormat, errDate := goment.New(v.DateOfIssue)
		if errDate != nil {
			return file, errDate
		}

		file.SetCellValue(sheetName, fmt.Sprintf("B%v", 7+(idx*4)), "Tanggal")
		file.SetCellValue(sheetName, fmt.Sprintf("C%v", 7+(idx*4)), ":")
		file.SetCellValue(sheetName, fmt.Sprintf("D%v", 7+(idx*4)), dateFormat.Format("DD MMMM YYYY", "id"))

		errDateRkab := file.SetCellStyle(sheetName, fmt.Sprintf("D%v", 7+(idx*4)), fmt.Sprintf("D%v", 7+(idx*4)), dateStyle)

		if errDateRkab != nil {
			return file, errDateRkab
		}

		file.SetCellValue(sheetName, fmt.Sprintf("B%v", 8+(idx*4)), "Quota Produksi")
		file.SetCellValue(sheetName, fmt.Sprintf("C%v", 8+(idx*4)), ":")

		errProductionQuota := file.SetCellStyle(sheetName, fmt.Sprintf("D%v", 8+(idx*4)), fmt.Sprintf("D%v", 8+(idx*4)), mtFmtStyle)

		if errProductionQuota != nil {
			return file, errProductionQuota
		}
		file.SetCellValue(sheetName, fmt.Sprintf("D%v", 8+(idx*4)), v.ProductionQuota)
	}

	var startRkab int

	startRkab = 5 + 2 + (len(reportSaleDetail.Rkabs) * 4)

	file.SetCellValue(sheetName, fmt.Sprintf("A%v", startRkab), "DMO")

	errTitleDmo := file.SetCellStyle(sheetName, fmt.Sprintf("A%v", startRkab), fmt.Sprintf("A%v", startRkab), boldTitleStyle)

	if errTitleDmo != nil {
		return file, errTitleDmo
	}

	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startRkab+1), "Kewajiban DMO")
	file.SetCellValue(sheetName, fmt.Sprintf("C%v", startRkab+1), ":")
	file.SetCellValue(sheetName, fmt.Sprintf("D%v", startRkab+1), maxProductionQuota*reportSaleDetail.Rkabs[0].DmoObligation/100)

	errDmoObligation := file.SetCellStyle(sheetName, fmt.Sprintf("D%v", startRkab+1), fmt.Sprintf("D%v", startRkab+1), mtFmtStyle)

	if errDmoObligation != nil {
		return file, errDmoObligation
	}

	var startProduction int

	startProduction = startRkab + 4

	file.SetCellValue(sheetName, fmt.Sprintf("A%v", startProduction), "PRODUKSI")
	errTitleProduction := file.SetCellStyle(sheetName, fmt.Sprintf("A%v", startProduction), fmt.Sprintf("A%v", startProduction), boldTitleStyle)

	if errTitleProduction != nil {
		return file, errTitleProduction
	}

	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+1), fmt.Sprintf("C%v", startProduction+1))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+2), fmt.Sprintf("C%v", startProduction+2))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+3), fmt.Sprintf("C%v", startProduction+3))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+4), fmt.Sprintf("C%v", startProduction+4))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+5), fmt.Sprintf("C%v", startProduction+5))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+6), fmt.Sprintf("C%v", startProduction+6))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+7), fmt.Sprintf("C%v", startProduction+7))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+8), fmt.Sprintf("C%v", startProduction+8))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+9), fmt.Sprintf("C%v", startProduction+9))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+10), fmt.Sprintf("C%v", startProduction+10))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+11), fmt.Sprintf("C%v", startProduction+11))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+12), fmt.Sprintf("C%v", startProduction+12))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+13), fmt.Sprintf("C%v", startProduction+13))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startProduction+14), fmt.Sprintf("C%v", startProduction+14))
	file.MergeCell(sheetName, fmt.Sprintf("E%v", startProduction+2), fmt.Sprintf("E%v", startProduction+14))

	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startProduction+1), "Bulan")
	file.SetCellValue(sheetName, fmt.Sprintf("D%v", startProduction+1), "Produksi")
	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startProduction+1), "RKAB")

	errTitleTable1 := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startProduction+1), fmt.Sprintf("E%v", startProduction+1), boldTitleTableStyle)

	if errTitleTable1 != nil {
		return file, errTitleTable1
	}

	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startProduction+2), maxProductionQuota)

	for idx, v := range monthString {
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startProduction+2+idx), v)
		switch v {
		case "January":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startProduction+2+idx), reportSaleDetail.Production.January)
		case "February":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startProduction+2+idx), reportSaleDetail.Production.February)
		case "March":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startProduction+2+idx), reportSaleDetail.Production.March)
		case "April":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startProduction+2+idx), reportSaleDetail.Production.April)
		case "May":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startProduction+2+idx), reportSaleDetail.Production.May)
		case "June":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startProduction+2+idx), reportSaleDetail.Production.June)
		case "July":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startProduction+2+idx), reportSaleDetail.Production.July)
		case "August":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startProduction+2+idx), reportSaleDetail.Production.August)
		case "September":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startProduction+2+idx), reportSaleDetail.Production.September)
		case "October":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startProduction+2+idx), reportSaleDetail.Production.October)
		case "November":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startProduction+2+idx), reportSaleDetail.Production.November)
		case "December":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startProduction+2+idx), reportSaleDetail.Production.December)

		}
	}
	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startProduction+14), "TOTAL")
	file.SetCellFormula(sheetName, fmt.Sprintf("D%v", startProduction+14), fmt.Sprintf("SUM(D%v:D%v)", startProduction+2, startProduction+13))

	file.SetCellFormula(sheetName, fmt.Sprintf("D%v", startProduction+15), fmt.Sprintf("D%v", startProduction+14))
	file.SetCellFormula(sheetName, fmt.Sprintf("E%v", startProduction+15), fmt.Sprintf("E%v", startProduction+2))

	file.SetCellValue(chartSheetName, "A4", "PRODUKSI")

	errChartTitle := file.SetColStyle(chartSheetName, "A", boldTitleStyle)
	if errChartTitle != nil {
		return file, errChartTitle
	}

	file.SetCellValue(chartSheetName, "B5", "RKAB")
	file.SetCellValue(chartSheetName, "B6", "PRODUKSI")

	file.SetCellFormula(chartSheetName, "C5", fmt.Sprintf("%v!E%v", sheetName, startProduction+2))
	file.SetCellFormula(chartSheetName, "C6", fmt.Sprintf("%v!D%v", sheetName, startProduction+14))

	errChartProductionRkabStyleTable := file.SetCellStyle(chartSheetName, "B5", "C6", formatNumberStyle)

	if errChartProductionRkabStyleTable != nil {
		return file, errChartProductionRkabStyleTable
	}

	var seriesChartProduction string

	seriesChartProduction += fmt.Sprintf(`{
						"name": "RKAB and Production",
						"categories": "%v!$B$5:$B$6",
						"values": "%v!$C$5:$C$6",
						"marker": {
									"symbol": "square"
								} 
				}`, chartSheetName, chartSheetName)

	valueChartProduction := fmt.Sprintf(`{
	    "type": "col",
	    "series": [%v],
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
	        "show_legend_key": true
	    },
	    "title":
	    {
	        "name": "RENCANA PRODUKSI DAN REALISASI PRODUKSI"
	    },
	    "show_blanks_as": "zero"
	}`, seriesChartProduction)

	if err := file.AddChart(chartSheetName, "G4", valueChartProduction); err != nil {
		return file, err
	}

	errProductionMonthStyleTable := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startProduction+2), fmt.Sprintf("B%v", startProduction+13), borderStyle)

	if errProductionMonthStyleTable != nil {
		return file, errProductionMonthStyleTable
	}

	errProductionStyleTable := file.SetCellStyle(sheetName, fmt.Sprintf("D%v", startProduction+2), fmt.Sprintf("D%v", startProduction+13), formatNumberStyle)

	if errProductionStyleTable != nil {
		return file, errProductionStyleTable
	}

	errProductionTotalStyle := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startProduction+14), fmt.Sprintf("B%v", startProduction+14), boldNumberStyle)

	if errProductionTotalStyle != nil {
		return file, errProductionTotalStyle
	}

	errProductionTotalNumberStyle := file.SetCellStyle(sheetName, fmt.Sprintf("D%v", startProduction+14), fmt.Sprintf("D%v", startProduction+14), boldNumberStyle)

	if errProductionTotalNumberStyle != nil {
		return file, errProductionTotalNumberStyle
	}

	errProductionRkabStyle := file.SetCellStyle(sheetName, fmt.Sprintf("E%v", startProduction+2), fmt.Sprintf("E%v", startProduction+14), boldNumberStyle)

	if errProductionRkabStyle != nil {
		return file, errProductionRkabStyle
	}

	var startSale int

	startSale = startProduction + 18

	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+1), fmt.Sprintf("C%v", startSale+1))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+2), fmt.Sprintf("C%v", startSale+2))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+3), fmt.Sprintf("C%v", startSale+3))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+4), fmt.Sprintf("C%v", startSale+4))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+5), fmt.Sprintf("C%v", startSale+5))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+6), fmt.Sprintf("C%v", startSale+6))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+7), fmt.Sprintf("C%v", startSale+7))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+8), fmt.Sprintf("C%v", startSale+8))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+9), fmt.Sprintf("C%v", startSale+9))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+10), fmt.Sprintf("C%v", startSale+10))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+11), fmt.Sprintf("C%v", startSale+11))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+12), fmt.Sprintf("C%v", startSale+12))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+13), fmt.Sprintf("C%v", startSale+13))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSale+14), fmt.Sprintf("C%v", startSale+14))

	file.SetCellValue(sheetName, fmt.Sprintf("A%v", startSale), "PENJUALAN")

	file.SetCellValue(chartSheetName, "A40", "REKAP DMO PER BULAN BERDASARKAN JENIS INDUSTRI")

	errTitlePenjualan := file.SetCellStyle(sheetName, fmt.Sprintf("A%v", startSale), fmt.Sprintf("A%v", startSale), boldTitleStyle)

	if errTitlePenjualan != nil {
		return file, errTitlePenjualan
	}

	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+1), "Bulan")
	file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+1), "Kelistrikan")
	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+1), "Non Kelistrikan")
	file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+1), "Jumlah")
	file.SetCellValue(sheetName, fmt.Sprintf("H%v", startSale+1), "Tidak Bisa Claim DMO")

	file.SetCellValue(chartSheetName, "B41", "Bulan")
	file.SetCellValue(chartSheetName, "C41", "Kelistrikan")
	file.SetCellValue(chartSheetName, "D41", "Non Kelistrikan")
	file.SetCellValue(chartSheetName, "E41", "Jumlah")

	errChartTitleRecap := file.SetCellStyle(chartSheetName, "B41", "E41", boldTitleTableStyle)

	if errChartTitleRecap != nil {
		return file, errChartTitleRecap
	}

	errTitleTablePenjualan := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startSale+1), fmt.Sprintf("F%v", startSale+1), boldTitleTableStyle)

	if errTitleTablePenjualan != nil {
		return file, errTitleTablePenjualan
	}

	errTitleTablePenjualan2 := file.SetCellStyle(sheetName, fmt.Sprintf("H%v", startSale+1), fmt.Sprintf("H%v", startSale+1), boldTitleTableStyle)

	if errTitleTablePenjualan2 != nil {
		return file, errTitleTablePenjualan2
	}

	for idx, v := range monthString {
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+2+idx), v)
		file.SetCellValue(chartSheetName, fmt.Sprintf("B%v", 42+idx), v)
		switch v {
		case "January":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.January)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.January)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.January+reportSaleDetail.RecapElectricity.January)

			file.SetCellValue(sheetName, fmt.Sprintf("H%v", startSale+2+idx), reportSaleDetail.NotClaimable.January)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.January)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapNonElectricity.January)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("E%v", 42+idx), fmt.Sprintf("SUM(C%v:D%v)", 42+idx, 42+idx))

		case "February":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.February)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.February)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.February+reportSaleDetail.RecapElectricity.February)

			file.SetCellValue(sheetName, fmt.Sprintf("H%v", startSale+2+idx), reportSaleDetail.NotClaimable.February)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.February)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapNonElectricity.February)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("E%v", 42+idx), fmt.Sprintf("SUM(C%v:D%v)", 42+idx, 42+idx))

		case "March":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.March)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.March)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.March+reportSaleDetail.RecapElectricity.March)

			file.SetCellValue(sheetName, fmt.Sprintf("H%v", startSale+2+idx), reportSaleDetail.NotClaimable.March)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.March)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapNonElectricity.March)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("E%v", 42+idx), fmt.Sprintf("SUM(C%v:D%v)", 42+idx, 42+idx))

		case "April":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.April)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.April)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.April+reportSaleDetail.RecapElectricity.April)

			file.SetCellValue(sheetName, fmt.Sprintf("H%v", startSale+2+idx), reportSaleDetail.NotClaimable.April)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.April)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapNonElectricity.April)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("E%v", 42+idx), fmt.Sprintf("SUM(C%v:D%v)", 42+idx, 42+idx))

		case "May":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.May)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.May)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.May+reportSaleDetail.RecapElectricity.May)

			file.SetCellValue(sheetName, fmt.Sprintf("H%v", startSale+2+idx), reportSaleDetail.NotClaimable.May)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.May)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapNonElectricity.May)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("E%v", 42+idx), fmt.Sprintf("SUM(C%v:D%v)", 42+idx, 42+idx))

		case "June":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.June)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.June)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.June+reportSaleDetail.RecapElectricity.June)

			file.SetCellValue(sheetName, fmt.Sprintf("H%v", startSale+2+idx), reportSaleDetail.NotClaimable.June)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.June)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapNonElectricity.June)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("E%v", 42+idx), fmt.Sprintf("SUM(C%v:D%v)", 42+idx, 42+idx))

		case "July":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.July)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.July)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.July+reportSaleDetail.RecapElectricity.July)

			file.SetCellValue(sheetName, fmt.Sprintf("H%v", startSale+2+idx), reportSaleDetail.NotClaimable.July)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.July)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapNonElectricity.July)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("E%v", 42+idx), fmt.Sprintf("SUM(C%v:D%v)", 42+idx, 42+idx))

		case "August":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.August)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.August)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.August+reportSaleDetail.RecapElectricity.August)

			file.SetCellValue(sheetName, fmt.Sprintf("H%v", startSale+2+idx), reportSaleDetail.NotClaimable.August)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.August)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapNonElectricity.August)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("E%v", 42+idx), fmt.Sprintf("SUM(C%v:D%v)", 42+idx, 42+idx))

		case "September":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.September)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.September)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.September+reportSaleDetail.RecapElectricity.September)

			file.SetCellValue(sheetName, fmt.Sprintf("H%v", startSale+2+idx), reportSaleDetail.NotClaimable.September)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.September)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapNonElectricity.September)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("E%v", 42+idx), fmt.Sprintf("SUM(C%v:D%v)", 42+idx, 42+idx))

		case "October":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.October)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.October)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.October+reportSaleDetail.RecapElectricity.October)

			file.SetCellValue(sheetName, fmt.Sprintf("H%v", startSale+2+idx), reportSaleDetail.NotClaimable.October)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.October)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapNonElectricity.October)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("E%v", 42+idx), fmt.Sprintf("SUM(C%v:D%v)", 42+idx, 42+idx))

		case "November":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.November)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.November)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.November+reportSaleDetail.RecapElectricity.November)

			file.SetCellValue(sheetName, fmt.Sprintf("H%v", startSale+2+idx), reportSaleDetail.NotClaimable.November)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.November)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapNonElectricity.November)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("E%v", 42+idx), fmt.Sprintf("SUM(C%v:D%v)", 42+idx, 42+idx))

		case "December":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.December)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.December)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.December+reportSaleDetail.RecapElectricity.December)

			file.SetCellValue(sheetName, fmt.Sprintf("H%v", startSale+2+idx), reportSaleDetail.NotClaimable.December)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.December)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapNonElectricity.December)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("E%v", 42+idx), fmt.Sprintf("SUM(C%v:D%v)", 42+idx, 42+idx))
		}
	}

	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+14), "TOTAL")
	file.SetCellFormula(sheetName, fmt.Sprintf("D%v", startSale+14), fmt.Sprintf("SUM(D%v:D%v)", startSale+2, startSale+13))
	file.SetCellFormula(sheetName, fmt.Sprintf("E%v", startSale+14), fmt.Sprintf("SUM(E%v:E%v)", startSale+2, startSale+13))
	file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startSale+14), fmt.Sprintf("SUM(F%v:F%v)", startSale+2, startSale+13))
	file.SetCellFormula(sheetName, fmt.Sprintf("H%v", startSale+14), fmt.Sprintf("SUM(H%v:H%v)", startSale+2, startSale+13))

	file.SetCellValue(chartSheetName, "B54", "TOTAL")
	file.SetCellFormula(chartSheetName, "C54", "SUM(C42:C53)")
	file.SetCellFormula(chartSheetName, "D54", "SUM(D42:D53)")
	file.SetCellFormula(chartSheetName, "E54", "SUM(E42:E53)")

	errChartRecapStyleTable := file.SetCellStyle(chartSheetName, "B42", "E53", formatNumberStyle)

	if errChartRecapStyleTable != nil {
		return file, errChartRecapStyleTable
	}

	errChartTotalRecapStyleTable := file.SetCellStyle(chartSheetName, "B54", "E54", boldNumberStyle)

	if errChartTotalRecapStyleTable != nil {
		return file, errChartTotalRecapStyleTable
	}

	var seriesChartRecapSale string

	seriesChartRecapSale += fmt.Sprintf(`{
						"name": "%v!$C$%v",
						"categories": "%v!$B$42:$B$53",
						"values": "%v!$C$42:$C$53",
						"marker": {
									"symbol": "square"
								} 
				},`, chartSheetName, 41, chartSheetName, chartSheetName)

	seriesChartRecapSale += fmt.Sprintf(`{
						"name": "%v!$D$%v",
						"categories": "%v!$B$42:$B$53",
						"values": "%v!$D$42:$D$53",
						"marker": {
									"symbol": "square"
								} 
				}`, chartSheetName, 41, chartSheetName, chartSheetName)

	valueChartRecapSale := fmt.Sprintf(`{
	    "type": "col",
	    "series": [%v],
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
	        "show_legend_key": true
	    },
	    "title":
	    {
	        "name": "REKAP DMO PER BULAN BERDASARKAN JENIS INDUSTRI"
	    },
	    "show_blanks_as": "zero"
	}`, seriesChartRecapSale)

	if err := file.AddChart(chartSheetName, "G40", valueChartRecapSale); err != nil {
		return file, err
	}

	errPenjualanTotalStyle := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startSale+14), fmt.Sprintf("B%v", startSale+14), boldNumberStyle)

	if errPenjualanTotalStyle != nil {
		return file, errPenjualanTotalStyle
	}

	errPenjualanTotalNumberStyle := file.SetCellStyle(sheetName, fmt.Sprintf("D%v", startSale+14), fmt.Sprintf("F%v", startSale+14), boldNumberStyle)

	if errPenjualanTotalNumberStyle != nil {
		return file, errPenjualanTotalNumberStyle
	}

	errPenjualanTotalNumberStyle2 := file.SetCellStyle(sheetName, fmt.Sprintf("H%v", startSale+14), fmt.Sprintf("H%v", startSale+14), boldNumberStyle)

	if errPenjualanTotalNumberStyle2 != nil {
		return file, errPenjualanTotalNumberStyle2
	}

	erPenjualanMonthStyleTable := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startSale+2), fmt.Sprintf("B%v", startSale+13), borderStyle)

	if erPenjualanMonthStyleTable != nil {
		return file, erPenjualanMonthStyleTable
	}

	errPenjualanNumberStyle := file.SetCellStyle(sheetName, fmt.Sprintf("D%v", startSale+2), fmt.Sprintf("F%v", startSale+13), formatNumberStyle)

	if errPenjualanNumberStyle != nil {
		return file, errPenjualanNumberStyle
	}

	errPenjualanNumberStyle2 := file.SetCellStyle(sheetName, fmt.Sprintf("H%v", startSale+2), fmt.Sprintf("H%v", startSale+13), formatNumberStyle)

	if errPenjualanNumberStyle2 != nil {
		return file, errPenjualanNumberStyle2
	}

	t := time.Now()

	month := t.Month()
	yearNow := t.Year()
	yearString := strconv.Itoa(yearNow)

	for idx, rkab := range reportSaleDetail.Rkabs {
		dateFormat, errDate := goment.New(rkab.DateOfIssue)
		if errDate != nil {
			return file, errDate
		}
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+16+(idx*10)), "% Pemenuhan DMO terhadap REALISASI PRODUKSI")
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+18+(idx*10)), "% Pemenuhan DMO terhadap RENCANA PRODUKSI")
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+19+(idx*10)), fmt.Sprintf("disetujui tgl %s", dateFormat.Format("DD MMMM YYYY", "id")))
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+21+(idx*10)), fmt.Sprintf("%% Pemenuhan DMO terhadap kewajiban pemenuhan DMO %.0f%%", rkab.DmoObligation))
		if yearString == year {
			file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+23+(idx*10)), fmt.Sprintf("%% Pemenuhan DMO terhadap Rencana Produksi (Prorata %v bulan)", int(month)))
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+23+(idx*10)), "% Pemenuhan DMO terhadap Rencana Produksi (Prorata 12 bulan)")
		}

		errBoldOnlyStyle := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startSale+16+(idx*10)), fmt.Sprintf("B%v", startSale+23+(idx*10)), boldOnlyStyle)

		if errBoldOnlyStyle != nil {
			return file, errBoldOnlyStyle
		}

		file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startSale+16+(idx*10)), fmt.Sprintf("F%v/D%v", startSale+14, startProduction+14))
		file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startSale+18+(idx*10)), fmt.Sprintf("F%v/D%v", startSale+14, 8+idx*4))
		file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startSale+21+(idx*10)), fmt.Sprintf("F%v/D%v", startSale+14, startRkab+1))
		if yearString == year {
			file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startSale+23+(idx*10)), fmt.Sprintf("F%v/(D%v*%v/12)", startSale+14, 8+idx*4, int(month)))
		} else {
			file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startSale+23+(idx*10)), fmt.Sprintf("F%v/(D%v*12/12)", startSale+14, 8+idx*4))
		}

		errBoldPercentStyle := file.SetCellStyle(sheetName, fmt.Sprintf("F%v", startSale+16+(idx*10)), fmt.Sprintf("F%v", startSale+23+(idx*10)), boldPercentStyle)

		if errBoldPercentStyle != nil {
			return file, errBoldPercentStyle
		}

	}

	file.SetCellValue(chartSheetName, "A61", "KEWAJIBAN DMO")
	file.SetCellValue(chartSheetName, "B62", fmt.Sprintf("Kewajiban DMO(%v%%)", reportSaleDetail.Rkabs[0].DmoObligation))
	file.SetCellValue(chartSheetName, "B63", "Pemenuhan DMO")

	file.SetCellValue(chartSheetName, "C62", 1)
	file.SetCellFormula(chartSheetName, "C63", "C79/C78")

	var seriesChartRecap string

	seriesChartRecap += fmt.Sprintf(`{
						"name": "Kewajiban DMO",
						"categories": "%v!$B$62:$B$63",
						"values": "%v!$C$62:$C$63",
						"marker": {
									"symbol": "square"
								} 
				}`, chartSheetName, chartSheetName)

	valueChartRecap := fmt.Sprintf(`{
	    "type": "col",
	    "series": [%v],
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
	        "show_legend_key": true
	    },
	    "title":
	    {
	        "name": "Pemenuhan DMO Tahun %s"
	    },
	    "show_blanks_as": "zero"
	}`, seriesChartRecap, year)

	if err := file.AddChart(chartSheetName, "G61", valueChartRecap); err != nil {
		return file, err
	}

	errPercentStyle := file.SetCellStyle(chartSheetName, "B62", "C63", percentTableStyle)

	if errPercentStyle != nil {
		return file, errPercentStyle
	}

	file.SetCellValue(chartSheetName, "B77", fmt.Sprintf("Pemenuhan DMO terhadap kewajiban pemenuhan DMO %v%%", reportSaleDetail.Rkabs[0].DmoObligation))

	errTitleDmo1 := file.SetCellStyle(chartSheetName, "B77", "B77", boldOnlyStyle)

	if errTitleDmo1 != nil {
		return file, errTitleDmo1
	}

	file.SetCellValue(chartSheetName, "B78", "Kewajiban DMO")
	file.SetCellValue(chartSheetName, "B79", "Realisasi DMO")
	file.SetCellValue(chartSheetName, "B80", "% Pemenuhan DMO")

	file.SetCellFormula(chartSheetName, "C78", fmt.Sprintf("Detail!D16"))
	file.SetCellFormula(chartSheetName, "C79", fmt.Sprintf("E54"))
	file.SetCellFormula(chartSheetName, "C80", fmt.Sprintf("C79/C78"))

	errChartPemenuhanDmo := file.SetCellStyle(chartSheetName, "B78", "C80", formatNumberStyle)

	if errChartPemenuhanDmo != nil {
		return file, errChartPemenuhanDmo
	}

	errChartPemenuhanPercentDmo := file.SetCellStyle(chartSheetName, "C80", "C80", percentTableStyle)

	if errChartPemenuhanPercentDmo != nil {
		return file, errChartPemenuhanPercentDmo
	}

	var seriesChartPemenuhanDmoTerhadapKewajibanDmo string

	seriesChartPemenuhanDmoTerhadapKewajibanDmo += fmt.Sprintf(`{
						"name": "Pemenuhan DMO Terhadap Kewajiban DMO",
						"categories": "%v!$B$78:$B$79",
						"values": "%v!$C$78:$C$79",
						"marker": {
									"symbol": "square"
								} 
				}`, chartSheetName, chartSheetName)

	valueChartPemenuhanDmoTerhadapKewajibanDmo := fmt.Sprintf(`{
	    "type": "col",
	    "series": [%v],
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
	        "show_legend_key": true
	    },
	    "title":
	    {
	        "name": "Pemenuhan DMO terhadap Kewajiban Pemenuhan DMO"
	    },
	    "show_blanks_as": "zero"
	}`, seriesChartPemenuhanDmoTerhadapKewajibanDmo)

	if err := file.AddChart(chartSheetName, "G77", valueChartPemenuhanDmoTerhadapKewajibanDmo); err != nil {
		return file, err
	}

	file.SetCellValue(chartSheetName, "B95", "Pemenuhan DMO terhadap Realisasi Produksi")

	errTitleDmo2 := file.SetCellStyle(chartSheetName, "B95", "B95", boldOnlyStyle)

	if errTitleDmo2 != nil {
		return file, errTitleDmo2
	}

	file.SetCellValue(chartSheetName, "B96", "Realisasi Produksi")
	file.SetCellValue(chartSheetName, "B97", "Realisasi DMO")
	file.SetCellValue(chartSheetName, "B98", "% Pemenuhan DMO terhadap Realisasi Produksi")

	file.SetCellFormula(chartSheetName, "C96", "C6")
	file.SetCellFormula(chartSheetName, "C97", "C79")
	file.SetCellFormula(chartSheetName, "C98", fmt.Sprintf("C97/C96"))

	errChartPemenuhanDmoTerhadapRealisasiProduksi := file.SetCellStyle(chartSheetName, "B96", "C98", formatNumberStyle)

	if errChartPemenuhanDmoTerhadapRealisasiProduksi != nil {
		return file, errChartPemenuhanDmoTerhadapRealisasiProduksi
	}

	errChartPercentPemenuhanDmoTerhadapRealisasiProduksi := file.SetCellStyle(chartSheetName, "C98", "C98", percentTableStyle)

	if errChartPercentPemenuhanDmoTerhadapRealisasiProduksi != nil {
		return file, errChartPercentPemenuhanDmoTerhadapRealisasiProduksi
	}

	var seriesChartPemenuhanDmoTerhadapRealisasiProduksi string

	seriesChartPemenuhanDmoTerhadapRealisasiProduksi += fmt.Sprintf(`{
						"name": "Pemenuhan DMO Terhadap Realisasi Produksi",
						"categories": "%v!$B$96:$B$97",
						"values": "%v!$C$96:$C$97",
						"marker": {
									"symbol": "square"
								} 
				}`, chartSheetName, chartSheetName)

	valueChartPemenuhanDmoTerhadapRealisasiProduksi := fmt.Sprintf(`{
	    "type": "col",
	    "series": [%v],
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
	        "show_legend_key": true
	    },
	    "title":
	    {
	        "name": "Pemenuhan DMO terhadap Realisasi Produksi"
	    },
	    "show_blanks_as": "zero"
	}`, seriesChartPemenuhanDmoTerhadapRealisasiProduksi)

	if err := file.AddChart(chartSheetName, "G95", valueChartPemenuhanDmoTerhadapRealisasiProduksi); err != nil {
		return file, err
	}

	file.SetCellValue(chartSheetName, "B111", "Pemenuhan DMO terhadap Rencana Produksi")

	errTitleDmo3 := file.SetCellStyle(chartSheetName, "B111", "B111", boldOnlyStyle)

	if errTitleDmo3 != nil {
		return file, errTitleDmo3
	}

	file.SetCellValue(chartSheetName, "B112", "Rencana Produksi")
	file.SetCellValue(chartSheetName, "B113", "Realisasi DMO")
	file.SetCellValue(chartSheetName, "B114", "% Pemenuhan DMO terhadap Rencana Produksi")

	file.SetCellFormula(chartSheetName, "C112", fmt.Sprintf("%s!E21", sheetName))
	file.SetCellFormula(chartSheetName, "C113", "C79")
	file.SetCellFormula(chartSheetName, "C114", fmt.Sprintf("C113/C112"))

	errChartPemenuhanDmoTerhadapRencanaProduksi := file.SetCellStyle(chartSheetName, "B112", "C114", formatNumberStyle)

	if errChartPemenuhanDmoTerhadapRencanaProduksi != nil {
		return file, errChartPemenuhanDmoTerhadapRencanaProduksi
	}

	errChartPercentPemenuhanDmoTerhadapRencanaProduksi := file.SetCellStyle(chartSheetName, "C114", "C114", percentTableStyle)

	if errChartPercentPemenuhanDmoTerhadapRencanaProduksi != nil {
		return file, errChartPercentPemenuhanDmoTerhadapRencanaProduksi
	}

	var seriesChartPemenuhanDmoTerhadapRencanaProduksi string

	seriesChartPemenuhanDmoTerhadapRencanaProduksi += fmt.Sprintf(`{
						"name": "Pemenuhan DMO Terhadap Rencana Produksi",
						"categories": "%v!$B$112:$B$113",
						"values": "%v!$C$112:$C$113",
						"marker": {
									"symbol": "square"
								} 
				}`, chartSheetName, chartSheetName)

	valueChartPemenuhanDmoTerhadapRencanaProduksi := fmt.Sprintf(`{
	    "type": "col",
	    "series": [%v],
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
	        "show_legend_key": true
	    },
	    "title":
	    {
	        "name": "Pemenuhan DMO terhadap Rencana Produksi"
	    },
	    "show_blanks_as": "zero"
	}`, seriesChartPemenuhanDmoTerhadapRencanaProduksi)

	if err := file.AddChart(chartSheetName, "G111", valueChartPemenuhanDmoTerhadapRencanaProduksi); err != nil {
		return file, err
	}

	file.SetCellValue(chartSheetName, "B126", fmt.Sprintf("Pemenuhan DMO terhadap Rencana Produksi (Prorata %v bulan)", int(month)))

	errTitleDmo4 := file.SetCellStyle(chartSheetName, "B126", "B126", boldOnlyStyle)

	if errTitleDmo4 != nil {
		return file, errTitleDmo4
	}

	file.SetCellValue(chartSheetName, "B127", "Rencana Produksi")

	file.SetCellValue(chartSheetName, "B129", "Prorata Produksi")
	file.SetCellValue(chartSheetName, "B130", "Realisasi DMO")
	file.SetCellValue(chartSheetName, "B131", "% Pemenuhan DMO terhadap Rencana Produksi (Prorata)")

	file.SetCellFormula(chartSheetName, "C127", fmt.Sprintf("%s!E21", sheetName))

	if year == yearString {
		file.SetCellFormula(chartSheetName, "C129", fmt.Sprintf("%s!E21*%v/12", sheetName, int(month)))
	} else {
		file.SetCellFormula(chartSheetName, "C129", fmt.Sprintf("%s!E21*12/12", sheetName))
	}
	file.SetCellFormula(chartSheetName, "C130", "C79")
	file.SetCellFormula(chartSheetName, "C131", fmt.Sprintf("C130/C129"))

	errChartPemenuhanDmoTerhadapRencanaProduksiProrata := file.SetCellStyle(chartSheetName, "B129", "C131", formatNumberStyle)

	if errChartPemenuhanDmoTerhadapRencanaProduksiProrata != nil {
		return file, errChartPemenuhanDmoTerhadapRencanaProduksiProrata
	}

	errChartPemenuhanDmoTerhadapRencanaProduksiProrata2 := file.SetCellStyle(chartSheetName, "B127", "C127", formatNumberStyle)

	if errChartPemenuhanDmoTerhadapRencanaProduksiProrata2 != nil {
		return file, errChartPemenuhanDmoTerhadapRencanaProduksiProrata2
	}

	errChartPercentPemenuhanDmoTerhadapRencanaProduksiProrata := file.SetCellStyle(chartSheetName, "C131", "C131", percentTableStyle)

	if errChartPercentPemenuhanDmoTerhadapRencanaProduksiProrata != nil {
		return file, errChartPercentPemenuhanDmoTerhadapRencanaProduksiProrata
	}

	var seriesChartPemenuhanDmoTerhadapRencanaProduksiProrata string

	seriesChartPemenuhanDmoTerhadapRencanaProduksiProrata += fmt.Sprintf(`{
						"name": "Pemenuhan DMO Terhadap Rencana Produksi Prorata (12 Bulan)",
						"categories": "%v!$B$129:$B$130",
						"values": "%v!$C$129:$C$130",
						"marker": {
									"symbol": "square"
								} 
				}`, chartSheetName, chartSheetName)

	valueChartPemenuhanDmoTerhadapRencanaProduksiProrata := fmt.Sprintf(`{
	    "type": "col",
	    "series": [%v],
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
	        "show_legend_key": true
	    },
	    "title":
	    {
	        "name": "Pemenuhan DMO terhadap Rencana Produksi (Prorata 12 bulan)"
	    },
	    "show_blanks_as": "zero"
	}`, seriesChartPemenuhanDmoTerhadapRencanaProduksiProrata)

	if err := file.AddChart(chartSheetName, "G126", valueChartPemenuhanDmoTerhadapRencanaProduksiProrata); err != nil {
		return file, err
	}

	file.SetCellValue(chartSheetName, "B143", "REKAP TOTAL DMO BERDASARKAN JENIS INDUSTRI")

	errTitleDmo5 := file.SetCellStyle(chartSheetName, "B143", "B143", boldOnlyStyle)

	if errTitleDmo5 != nil {
		return file, errTitleDmo5
	}

	file.SetCellValue(chartSheetName, "B144", "Kelistrikan")
	file.SetCellValue(chartSheetName, "B145", "Non Kelistrikan")

	file.SetCellFormula(chartSheetName, "C144", "C54")
	file.SetCellFormula(chartSheetName, "C145", "D54")

	errChartRecapJenisIndustri := file.SetCellStyle(chartSheetName, "B144", "C145", formatNumberStyle)

	if errChartRecapJenisIndustri != nil {
		return file, errChartRecapJenisIndustri
	}

	var seriesChartRecapJenisIndustri string

	seriesChartRecapJenisIndustri += fmt.Sprintf(`{
						"name": "REKAP TOTAL DMO BERDASARKAN JENIS INDUSTRI",
						"categories": "%v!$B$144:$B$145",
						"values": "%v!$C$144:$C$145",
						"marker": {
									"symbol": "square"
								} 
				}`, chartSheetName, chartSheetName)

	valueChartRecapJenisIndustri := fmt.Sprintf(`{
	    "type": "pie",
	    "series": [%v],
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
	        "show_legend_key": true
	    },
			"plotarea": {
				"show_percent": true
			},
	    "title":
	    {
	        "name": "Penjualan Kelistrikan & Non Kelistrikan"
	    },
	    "show_blanks_as": "zero"
	}`, seriesChartRecapJenisIndustri)

	if err := file.AddChart(chartSheetName, "G143", valueChartRecapJenisIndustri); err != nil {
		return file, err
	}

	var startDetail int

	startDetail = startSale + 17 + (len(reportSaleDetail.Rkabs) * 10)

	file.SetCellValue(sheetName, fmt.Sprintf("A%v", startDetail), "DETAIL PENJUALAN")

	errTitleDetailPenjualan := file.SetCellStyle(sheetName, fmt.Sprintf("A%v", startDetail), fmt.Sprintf("A%v", startDetail), boldTitleStyle)

	if errTitleDetailPenjualan != nil {
		return file, errTitleDetailPenjualan
	}

	file.SetCellValue(sheetName, fmt.Sprintf("A%v", startDetail+1), "Data Realisasi Penjualan Batubara Untuk Memenuhi Kebutuhan Batubara Bagi Penyediaan Tenaga Listrik Untuk Kepentingan Umum / Kelistrikan")

	file.MergeCell(sheetName, fmt.Sprintf("A%v", startDetail+2), fmt.Sprintf("A%v", startDetail+3))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startDetail+2), fmt.Sprintf("D%v", startDetail+3))

	file.MergeCell(sheetName, fmt.Sprintf("E%v", startDetail+2), fmt.Sprintf("Q%v", startDetail+2))

	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startDetail+2), "Realisasi Penjualan Batubara (Pengguna Akhir Batubara) (Ton)")

	file.SetCellValue(sheetName, fmt.Sprintf("A%v", startDetail+2), "No.")
	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startDetail+2), "END USER")
	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startDetail+3), "January")
	file.SetCellValue(sheetName, fmt.Sprintf("F%v", startDetail+3), "February")
	file.SetCellValue(sheetName, fmt.Sprintf("G%v", startDetail+3), "March")
	file.SetCellValue(sheetName, fmt.Sprintf("H%v", startDetail+3), "April")
	file.SetCellValue(sheetName, fmt.Sprintf("I%v", startDetail+3), "May")
	file.SetCellValue(sheetName, fmt.Sprintf("J%v", startDetail+3), "June")
	file.SetCellValue(sheetName, fmt.Sprintf("K%v", startDetail+3), "July")
	file.SetCellValue(sheetName, fmt.Sprintf("L%v", startDetail+3), "August")
	file.SetCellValue(sheetName, fmt.Sprintf("M%v", startDetail+3), "September")
	file.SetCellValue(sheetName, fmt.Sprintf("N%v", startDetail+3), "October")
	file.SetCellValue(sheetName, fmt.Sprintf("O%v", startDetail+3), "November")
	file.SetCellValue(sheetName, fmt.Sprintf("P%v", startDetail+3), "December")
	file.SetCellValue(sheetName, fmt.Sprintf("Q%v", startDetail+3), "TOTAL")

	file.SetCellValue(chartSheetName, "B160", "DETAIL REKAP TOTAL DMO BERDASARKAN JENIS INDUSTRI - KELISTRIKAN")

	errTitleDmo6 := file.SetCellStyle(chartSheetName, "B160", "B160", boldOnlyStyle)

	if errTitleDmo6 != nil {
		return file, errTitleDmo6
	}

	var countPltu int
	var numberElectric int = 1
	for k, value := range reportSaleDetail.CompanyElectricity {

		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startDetail+(4+numberElectric-1)), k)

		for _, v := range value {

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startDetail+(4+numberElectric-1)), v)
			file.SetCellValue(sheetName, fmt.Sprintf("A%v", startDetail+(4+numberElectric-1)), numberElectric)

			if _, ok := reportSaleDetail.Electricity.January[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("E%v", startDetail+(4+numberElectric-1)), reportSaleDetail.Electricity.January[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("E%v", startDetail+(4+numberElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.Electricity.February[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("F%v", startDetail+(4+numberElectric-1)), reportSaleDetail.Electricity.February[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("F%v", startDetail+(4+numberElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.Electricity.March[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("G%v", startDetail+(4+numberElectric-1)), reportSaleDetail.Electricity.March[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("G%v", startDetail+(4+numberElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.Electricity.April[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("H%v", startDetail+(4+numberElectric-1)), reportSaleDetail.Electricity.April[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("H%v", startDetail+(4+numberElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.Electricity.May[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("I%v", startDetail+(4+numberElectric-1)), reportSaleDetail.Electricity.May[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("I%v", startDetail+(4+numberElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.Electricity.June[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("J%v", startDetail+(4+numberElectric-1)), reportSaleDetail.Electricity.June[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("J%v", startDetail+(4+numberElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.Electricity.July[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("K%v", startDetail+(4+numberElectric-1)), reportSaleDetail.Electricity.July[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("K%v", startDetail+(4+numberElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.Electricity.August[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("L%v", startDetail+(4+numberElectric-1)), reportSaleDetail.Electricity.August[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("L%v", startDetail+(4+numberElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.Electricity.September[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("M%v", startDetail+(4+numberElectric-1)), reportSaleDetail.Electricity.September[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("M%v", startDetail+(4+numberElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.Electricity.October[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("N%v", startDetail+(4+numberElectric-1)), reportSaleDetail.Electricity.October[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("N%v", startDetail+(4+numberElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.Electricity.November[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("O%v", startDetail+(4+numberElectric-1)), reportSaleDetail.Electricity.November[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("O%v", startDetail+(4+numberElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.Electricity.December[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("P%v", startDetail+(4+numberElectric-1)), reportSaleDetail.Electricity.December[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("P%v", startDetail+(4+numberElectric-1)), 0)
			}

			file.SetCellFormula(sheetName, fmt.Sprintf("Q%v", startDetail+(4+numberElectric-1)), fmt.Sprintf("SUM(E%v:P%v)", startDetail+(4+numberElectric-1), startDetail+(4+numberElectric-1)))
			numberElectric += 1

			file.SetCellFormula(chartSheetName, fmt.Sprintf("B%v", 161+countPltu), fmt.Sprintf("Detail!D%v", startDetail+(4+numberElectric-2)))
			file.SetCellFormula(chartSheetName, fmt.Sprintf("C%v", 161+countPltu), fmt.Sprintf("Detail!Q%v", startDetail+(4+numberElectric-2)))
			countPltu += 1
		}

		file.MergeCell(sheetName, fmt.Sprintf("B%v", startDetail+numberElectric+1), fmt.Sprintf("C%v", startDetail+numberElectric+len(value)))

	}

	errChartSaleElectric := file.SetCellStyle(chartSheetName, "B161", fmt.Sprintf("C%v", 161+countPltu-1), formatNumberStyle)

	if errChartSaleElectric != nil {
		return file, errChartSaleElectric
	}

	var seriesChartCompanyElectric string

	seriesChartCompanyElectric += fmt.Sprintf(`{
						"name": "Sales Kelistrikan",
						"categories": "%v!$B$161:$B$%v",
						"values": "%v!$C$161:$C$%v",
						"marker": {
									"symbol": "square"
								} 
				}`, chartSheetName, 161+countPltu-1, chartSheetName, 161+countPltu-1)

	valueChartElectric := fmt.Sprintf(`{
	    "type": "pie",
	    "series": [%v],
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
	        "show_legend_key": true
	    },
			"plotarea": 
			{
					"show_percent": true
			},
	    "title":
	    {
	        "name": "Penjualan berdasarkan JENIS INDUSTRI (KELISTRIKAN)"
	    },
	    "show_blanks_as": "zero"
	}`, seriesChartCompanyElectric)

	if err := file.AddChart(chartSheetName, "G160", valueChartElectric); err != nil {
		return file, err
	}

	errNoStyle := file.SetCellStyle(sheetName, fmt.Sprintf("A%v", startDetail+4), fmt.Sprintf("A%v", startDetail+3+numberElectric), centerStyle)

	if errNoStyle != nil {
		return file, errNoStyle
	}

	errTitleTableElectricStyle := file.SetCellStyle(sheetName, fmt.Sprintf("A%v", startDetail+2), fmt.Sprintf("Q%v", startDetail+3), boldTitleTableStyle)

	if errTitleTableElectricStyle != nil {
		return file, errTitleTableElectricStyle
	}

	errNumberElectricStyle := file.SetCellStyle(sheetName, fmt.Sprintf("E%v", startDetail+4), fmt.Sprintf("P%v", startDetail+3+numberElectric-1), formatNumberStyle)

	if errNumberElectricStyle != nil {
		return file, errNumberElectricStyle
	}

	errNumberTotalElectricStyle := file.SetCellStyle(sheetName, fmt.Sprintf("Q%v", startDetail+4), fmt.Sprintf("Q%v", startDetail+3+numberElectric-1), boldNumberStyle)

	if errNumberTotalElectricStyle != nil {
		return file, errNumberTotalElectricStyle
	}

	errCompanyElectricStyle := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startDetail+3), fmt.Sprintf("D%v", startDetail+3+numberElectric-1), boldTitleTableStyle)

	if errCompanyElectricStyle != nil {
		return file, errCompanyElectricStyle
	}

	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startDetail+3+numberElectric), "TOTAL")
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startDetail+3+numberElectric), fmt.Sprintf("D%v", startDetail+3+numberElectric))

	errSaleDetailElectricTotalStyle := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startDetail+3+numberElectric), fmt.Sprintf("B%v", startDetail+3+numberElectric), boldNumberStyle)

	if errSaleDetailElectricTotalStyle != nil {
		return file, errSaleDetailElectricTotalStyle
	}

	errSaleDetailElectricTotalNumberStyle := file.SetCellStyle(sheetName, fmt.Sprintf("E%v", startDetail+3+numberElectric), fmt.Sprintf("Q%v", startDetail+3+numberElectric), boldNumberStyle)

	if errSaleDetailElectricTotalNumberStyle != nil {
		return file, errSaleDetailElectricTotalNumberStyle
	}

	file.SetCellFormula(sheetName, fmt.Sprintf("E%v", startDetail+3+numberElectric), fmt.Sprintf("SUM(E%v:E%v)", startDetail+3, startDetail+2+numberElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startDetail+3+numberElectric), fmt.Sprintf("SUM(F%v:F%v)", startDetail+3, startDetail+2+numberElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("G%v", startDetail+3+numberElectric), fmt.Sprintf("SUM(G%v:G%v)", startDetail+3, startDetail+2+numberElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("H%v", startDetail+3+numberElectric), fmt.Sprintf("SUM(H%v:H%v)", startDetail+3, startDetail+2+numberElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("I%v", startDetail+3+numberElectric), fmt.Sprintf("SUM(I%v:I%v)", startDetail+3, startDetail+2+numberElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("J%v", startDetail+3+numberElectric), fmt.Sprintf("SUM(J%v:J%v)", startDetail+3, startDetail+2+numberElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("K%v", startDetail+3+numberElectric), fmt.Sprintf("SUM(K%v:K%v)", startDetail+3, startDetail+2+numberElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("L%v", startDetail+3+numberElectric), fmt.Sprintf("SUM(L%v:L%v)", startDetail+3, startDetail+2+numberElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("M%v", startDetail+3+numberElectric), fmt.Sprintf("SUM(M%v:M%v)", startDetail+3, startDetail+2+numberElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("N%v", startDetail+3+numberElectric), fmt.Sprintf("SUM(N%v:N%v)", startDetail+3, startDetail+2+numberElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("O%v", startDetail+3+numberElectric), fmt.Sprintf("SUM(O%v:O%v)", startDetail+3, startDetail+2+numberElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("P%v", startDetail+3+numberElectric), fmt.Sprintf("SUM(P%v:P%v)", startDetail+3, startDetail+2+numberElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("Q%v", startDetail+3+numberElectric), fmt.Sprintf("SUM(Q%v:Q%v)", startDetail+3, startDetail+2+numberElectric))

	var startDetailNonElectric int
	startDetailNonElectric = startDetail + 3 + numberElectric + 3

	file.SetCellValue(sheetName, fmt.Sprintf("A%v", startDetailNonElectric), "Data Realisasi Penjualan Batubara Untuk Memenuhi Kebutuhan Batubara Untuk Industri / Non Kelistrikan Umum")

	file.MergeCell(sheetName, fmt.Sprintf("A%v", startDetailNonElectric+1), fmt.Sprintf("A%v", startDetailNonElectric+2))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startDetailNonElectric+1), fmt.Sprintf("D%v", startDetailNonElectric+2))

	file.MergeCell(sheetName, fmt.Sprintf("E%v", startDetailNonElectric+1), fmt.Sprintf("Q%v", startDetailNonElectric+1))

	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startDetailNonElectric+1), "Realisasi Penjualan Batubara (Pengguna Akhir Batubara) (Ton)")

	file.SetCellValue(sheetName, fmt.Sprintf("A%v", startDetailNonElectric+2), "No.")
	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startDetailNonElectric+2), "END USER")
	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startDetailNonElectric+2), "January")
	file.SetCellValue(sheetName, fmt.Sprintf("F%v", startDetailNonElectric+2), "February")
	file.SetCellValue(sheetName, fmt.Sprintf("G%v", startDetailNonElectric+2), "March")
	file.SetCellValue(sheetName, fmt.Sprintf("H%v", startDetailNonElectric+2), "April")
	file.SetCellValue(sheetName, fmt.Sprintf("I%v", startDetailNonElectric+2), "May")
	file.SetCellValue(sheetName, fmt.Sprintf("J%v", startDetailNonElectric+2), "June")
	file.SetCellValue(sheetName, fmt.Sprintf("K%v", startDetailNonElectric+2), "July")
	file.SetCellValue(sheetName, fmt.Sprintf("L%v", startDetailNonElectric+2), "August")
	file.SetCellValue(sheetName, fmt.Sprintf("M%v", startDetailNonElectric+2), "September")
	file.SetCellValue(sheetName, fmt.Sprintf("N%v", startDetailNonElectric+2), "October")
	file.SetCellValue(sheetName, fmt.Sprintf("O%v", startDetailNonElectric+2), "November")
	file.SetCellValue(sheetName, fmt.Sprintf("P%v", startDetailNonElectric+2), "December")
	file.SetCellValue(sheetName, fmt.Sprintf("Q%v", startDetailNonElectric+2), "TOTAL")

	var countIndustry int
	if countPltu > 16 {
		countIndustry = countPltu - 16
	}

	var base int

	base = 178 + countIndustry
	file.SetCellValue(chartSheetName, fmt.Sprintf("B%v", 178+countIndustry), "DETAIL REKAP TOTAL DMO BERDASARKAN JENIS INDUSTRI - NON KELISTRIKAN")

	errTitleDmo7 := file.SetCellStyle(chartSheetName, fmt.Sprintf("B%v", 178+countIndustry), fmt.Sprintf("B%v", 178+countIndustry), boldOnlyStyle)

	if errTitleDmo7 != nil {
		return file, errTitleDmo7
	}

	var numberNonElectric int = 1

	for k, value := range reportSaleDetail.CompanyNonElectricity {

		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startDetailNonElectric+(3+numberNonElectric-1)), k)

		for _, v := range value {
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startDetailNonElectric+(3+numberNonElectric-1)), v)
			file.SetCellValue(sheetName, fmt.Sprintf("A%v", startDetailNonElectric+(3+numberNonElectric-1)), numberNonElectric)

			if _, ok := reportSaleDetail.NonElectricity.January[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("E%v", startDetailNonElectric+(3+numberNonElectric-1)), reportSaleDetail.NonElectricity.January[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("E%v", startDetailNonElectric+(3+numberNonElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.NonElectricity.February[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("F%v", startDetailNonElectric+(3+numberNonElectric-1)), reportSaleDetail.NonElectricity.February[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("F%v", startDetailNonElectric+(3+numberNonElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.NonElectricity.March[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("G%v", startDetailNonElectric+(3+numberNonElectric-1)), reportSaleDetail.NonElectricity.March[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("G%v", startDetailNonElectric+(3+numberNonElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.NonElectricity.April[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("H%v", startDetailNonElectric+(3+numberNonElectric-1)), reportSaleDetail.NonElectricity.April[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("H%v", startDetailNonElectric+(3+numberNonElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.NonElectricity.May[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("I%v", startDetailNonElectric+(3+numberNonElectric-1)), reportSaleDetail.NonElectricity.May[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("I%v", startDetailNonElectric+(3+numberNonElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.NonElectricity.June[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("J%v", startDetailNonElectric+(3+numberNonElectric-1)), reportSaleDetail.NonElectricity.June[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("J%v", startDetailNonElectric+(3+numberNonElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.NonElectricity.July[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("K%v", startDetailNonElectric+(3+numberNonElectric-1)), reportSaleDetail.NonElectricity.July[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("K%v", startDetailNonElectric+(3+numberNonElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.NonElectricity.August[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("L%v", startDetailNonElectric+(3+numberNonElectric-1)), reportSaleDetail.NonElectricity.August[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("L%v", startDetailNonElectric+(3+numberNonElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.NonElectricity.September[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("M%v", startDetailNonElectric+(3+numberNonElectric-1)), reportSaleDetail.NonElectricity.September[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("M%v", startDetailNonElectric+(3+numberNonElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.NonElectricity.October[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("N%v", startDetailNonElectric+(3+numberNonElectric-1)), reportSaleDetail.NonElectricity.October[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("N%v", startDetailNonElectric+(3+numberNonElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.NonElectricity.November[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("O%v", startDetailNonElectric+(3+numberNonElectric-1)), reportSaleDetail.NonElectricity.November[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("O%v", startDetailNonElectric+(3+numberNonElectric-1)), 0)
			}

			if _, ok := reportSaleDetail.NonElectricity.December[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("P%v", startDetailNonElectric+(3+numberNonElectric-1)), reportSaleDetail.NonElectricity.December[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("P%v", startDetailNonElectric+(3+numberNonElectric-1)), 0)
			}

			file.SetCellFormula(sheetName, fmt.Sprintf("Q%v", startDetailNonElectric+(3+numberNonElectric-1)), fmt.Sprintf("SUM(E%v:P%v)", startDetailNonElectric+(3+numberNonElectric-1), startDetailNonElectric+(3+numberNonElectric-1)))

			file.SetCellFormula(chartSheetName, fmt.Sprintf("B%v", base+1+countIndustry), fmt.Sprintf("Detail!D%v", startDetailNonElectric+(3+numberNonElectric-1)))
			file.SetCellFormula(chartSheetName, fmt.Sprintf("C%v", base+1+countIndustry), fmt.Sprintf("Detail!Q%v", startDetailNonElectric+(3+numberNonElectric-1)))
			countIndustry += 1
			numberNonElectric += 1
		}

		file.MergeCell(sheetName, fmt.Sprintf("B%v", startDetailNonElectric+numberNonElectric), fmt.Sprintf("C%v", startDetailNonElectric+numberNonElectric+len(value)-1))
	}

	errChartSaleNonElectric := file.SetCellStyle(chartSheetName, fmt.Sprintf("B%v", base+1), fmt.Sprintf("C%v", base+countIndustry), formatNumberStyle)

	if errChartSaleNonElectric != nil {
		return file, errChartSaleNonElectric
	}

	var seriesChartCompanyNonElectric string

	seriesChartCompanyNonElectric += fmt.Sprintf(`{
						"name": "Sales Non Kelistrikan",
						"categories": "%v!$B$%v:$B$%v",
						"values": "%v!$C$%v:$C$%v",
						"marker": {
									"symbol": "square"
								} 
				}`, chartSheetName, base+1, base+1+countIndustry, chartSheetName, base+1, base+1+countIndustry)

	valueChartNonElectric := fmt.Sprintf(`{
	    "type": "pie",
	    "series": [%v],
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
	        "show_legend_key": true
	    },
			"plotarea": 
			{
					"show_percent": true
			},
	    "title":
	    {
	        "name": "Penjualan berdasarkan JENIS INDUSTRI (NON KELISTRIKAN)"
	    },
	    "show_blanks_as": "zero"
	}`, seriesChartCompanyNonElectric)

	if err := file.AddChart(chartSheetName, fmt.Sprintf("G%v", base), valueChartNonElectric); err != nil {
		return file, err
	}

	errNoNonElectricStyle := file.SetCellStyle(sheetName, fmt.Sprintf("A%v", startDetailNonElectric+3), fmt.Sprintf("A%v", startDetailNonElectric+3+numberNonElectric-1), centerStyle)

	if errNoNonElectricStyle != nil {
		return file, errNoNonElectricStyle
	}

	errTitleTableNonElectricStyle := file.SetCellStyle(sheetName, fmt.Sprintf("A%v", startDetailNonElectric+1), fmt.Sprintf("Q%v", startDetailNonElectric+2), boldTitleTableStyle)

	if errTitleTableNonElectricStyle != nil {
		return file, errTitleTableNonElectricStyle
	}

	errNumberNonElectricStyle := file.SetCellStyle(sheetName, fmt.Sprintf("E%v", startDetailNonElectric+3), fmt.Sprintf("P%v", startDetailNonElectric+3+numberNonElectric-1), formatNumberStyle)

	if errNumberNonElectricStyle != nil {
		return file, errNumberNonElectricStyle
	}

	errNumberTotalNonElectricStyle := file.SetCellStyle(sheetName, fmt.Sprintf("Q%v", startDetailNonElectric+3), fmt.Sprintf("Q%v", startDetailNonElectric+3+numberNonElectric-1), boldNumberStyle)

	if errNumberTotalNonElectricStyle != nil {
		return file, errNumberTotalNonElectricStyle
	}

	errCompanyNonElectricStyle := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startDetailNonElectric+3), fmt.Sprintf("D%v", startDetailNonElectric+3+numberNonElectric-1), boldTitleTableStyle)

	if errCompanyNonElectricStyle != nil {
		return file, errCompanyNonElectricStyle
	}

	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startDetail+3+numberNonElectric), "TOTAL")
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startDetail+3+numberNonElectric), fmt.Sprintf("D%v", startDetail+3+numberNonElectric))

	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startDetailNonElectric+2+numberNonElectric), "TOTAL")
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("D%v", startDetailNonElectric+2+numberNonElectric))

	errSaleDetailNonElectricTotalStyle := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("B%v", startDetailNonElectric+2+numberNonElectric), boldNumberStyle)

	if errSaleDetailNonElectricTotalStyle != nil {
		return file, errSaleDetailNonElectricTotalStyle
	}

	errSaleDetailNonElectricTotalNumberStyle := file.SetCellStyle(sheetName, fmt.Sprintf("E%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("Q%v", startDetailNonElectric+2+numberNonElectric), boldNumberStyle)

	if errSaleDetailNonElectricTotalNumberStyle != nil {
		return file, errSaleDetailNonElectricTotalNumberStyle
	}

	file.SetCellFormula(sheetName, fmt.Sprintf("E%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(E%v:E%v)", startDetailNonElectric+3, startDetailNonElectric+numberNonElectric+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(F%v:F%v)", startDetailNonElectric+3, startDetailNonElectric+numberNonElectric+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("G%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(G%v:G%v)", startDetailNonElectric+3, startDetailNonElectric+numberNonElectric+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("H%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(H%v:H%v)", startDetailNonElectric+3, startDetailNonElectric+numberNonElectric+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("I%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(I%v:I%v)", startDetailNonElectric+3, startDetailNonElectric+numberNonElectric+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("J%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(J%v:J%v)", startDetailNonElectric+3, startDetailNonElectric+numberNonElectric+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("K%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(K%v:K%v)", startDetailNonElectric+3, startDetailNonElectric+numberNonElectric+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("L%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(L%v:L%v)", startDetailNonElectric+3, startDetailNonElectric+numberNonElectric+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("M%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(M%v:M%v)", startDetailNonElectric+3, startDetailNonElectric+numberNonElectric+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("N%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(N%v:N%v)", startDetailNonElectric+3, startDetailNonElectric+numberNonElectric+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("O%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(O%v:O%v)", startDetailNonElectric+3, startDetailNonElectric+numberNonElectric+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("P%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(P%v:P%v)", startDetailNonElectric+3, startDetailNonElectric+numberNonElectric+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("Q%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(Q%v:Q%v)", startDetailNonElectric+3, startDetailNonElectric+numberNonElectric+1))

	var startSaleExportImport int
	startSaleExportImport = startDetailNonElectric + 6 + numberNonElectric

	file.SetCellValue(sheetName, fmt.Sprintf("A%v", startSaleExportImport), "PENJUALAN BATUBARA")
	errTitleDetailPenjualanBatubara := file.SetCellStyle(sheetName, fmt.Sprintf("A%v", startSaleExportImport), fmt.Sprintf("A%v", startSaleExportImport), boldTitleStyle)

	if errTitleDetailPenjualanBatubara != nil {
		return file, errTitleDetailPenjualanBatubara
	}

	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+1), fmt.Sprintf("C%v", startSaleExportImport+1))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+2), fmt.Sprintf("C%v", startSaleExportImport+2))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+3), fmt.Sprintf("C%v", startSaleExportImport+3))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+4), fmt.Sprintf("C%v", startSaleExportImport+4))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+5), fmt.Sprintf("C%v", startSaleExportImport+5))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+6), fmt.Sprintf("C%v", startSaleExportImport+6))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+7), fmt.Sprintf("C%v", startSaleExportImport+7))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+8), fmt.Sprintf("C%v", startSaleExportImport+8))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+9), fmt.Sprintf("C%v", startSaleExportImport+9))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+10), fmt.Sprintf("C%v", startSaleExportImport+10))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+11), fmt.Sprintf("C%v", startSaleExportImport+11))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+12), fmt.Sprintf("C%v", startSaleExportImport+12))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+13), fmt.Sprintf("C%v", startSaleExportImport+13))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startSaleExportImport+14), fmt.Sprintf("C%v", startSaleExportImport+14))

	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSaleExportImport+1), "Periode")
	file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSaleExportImport+1), "Domestik (MT)")
	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSaleExportImport+1), "Ekspor (MT)")

	errTitleTable2 := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startSaleExportImport+1), fmt.Sprintf("E%v", startSaleExportImport+1), boldTitleTableStyle)

	if errTitleTable2 != nil {
		return file, errTitleTable2
	}

	for idx, v := range monthString {
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSaleExportImport+2+idx), v)
		switch v {
		case "January":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSaleExportImport+2+idx), reportSaleDetail.Domestic.January)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSaleExportImport+2+idx), reportSaleDetail.Export.January)
		case "February":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSaleExportImport+2+idx), reportSaleDetail.Domestic.February)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSaleExportImport+2+idx), reportSaleDetail.Export.February)

		case "March":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSaleExportImport+2+idx), reportSaleDetail.Domestic.March)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSaleExportImport+2+idx), reportSaleDetail.Export.March)
		case "April":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSaleExportImport+2+idx), reportSaleDetail.Domestic.April)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSaleExportImport+2+idx), reportSaleDetail.Export.April)
		case "May":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSaleExportImport+2+idx), reportSaleDetail.Domestic.May)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSaleExportImport+2+idx), reportSaleDetail.Export.May)

		case "June":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSaleExportImport+2+idx), reportSaleDetail.Domestic.June)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSaleExportImport+2+idx), reportSaleDetail.Export.June)

		case "July":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSaleExportImport+2+idx), reportSaleDetail.Domestic.July)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSaleExportImport+2+idx), reportSaleDetail.Export.July)
		case "August":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSaleExportImport+2+idx), reportSaleDetail.Domestic.August)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSaleExportImport+2+idx), reportSaleDetail.Export.August)

		case "September":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSaleExportImport+2+idx), reportSaleDetail.Domestic.September)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSaleExportImport+2+idx), reportSaleDetail.Export.September)
		case "October":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSaleExportImport+2+idx), reportSaleDetail.Domestic.October)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSaleExportImport+2+idx), reportSaleDetail.Export.October)

		case "November":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSaleExportImport+2+idx), reportSaleDetail.Domestic.November)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSaleExportImport+2+idx), reportSaleDetail.Export.November)
		case "December":
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSaleExportImport+2+idx), reportSaleDetail.Domestic.December)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSaleExportImport+2+idx), reportSaleDetail.Export.December)
		}
	}

	erBorderDomesticExportStyle := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startSaleExportImport+2), fmt.Sprintf("B%v", startSaleExportImport+13), borderStyle)

	if erBorderDomesticExportStyle != nil {
		return file, erBorderDomesticExportStyle
	}

	errDomesticExportNumberStyle := file.SetCellStyle(sheetName, fmt.Sprintf("D%v", startSaleExportImport+2), fmt.Sprintf("E%v", startSaleExportImport+13), formatNumberStyle)

	if errDomesticExportNumberStyle != nil {
		return file, errDomesticExportNumberStyle
	}

	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSaleExportImport+14), "TOTAL")
	file.SetCellFormula(sheetName, fmt.Sprintf("D%v", startSaleExportImport+14), fmt.Sprintf("SUM(D%v:D%v)", startSaleExportImport+1, startSaleExportImport+13))
	file.SetCellFormula(sheetName, fmt.Sprintf("E%v", startSaleExportImport+14), fmt.Sprintf("SUM(E%v:E%v)", startSaleExportImport+1, startSaleExportImport+13))

	errDomesticExportTotalStyle := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startSaleExportImport+14), fmt.Sprintf("E%v", startSaleExportImport+14), boldNumberStyle)

	if errDomesticExportTotalStyle != nil {
		return file, errDomesticExportTotalStyle
	}

	file.SetCellValue(chartSheetName, "A22", "PENJUALAN")

	file.SetCellValue(chartSheetName, "B23", "DALAM NEGERI")
	file.SetCellValue(chartSheetName, "B24", "LUAR NEGERI")
	file.SetCellFormula(chartSheetName, "C23", fmt.Sprintf("%v!D%v", sheetName, startSaleExportImport+14))
	file.SetCellFormula(chartSheetName, "C24", fmt.Sprintf("%v!E%v", sheetName, startSaleExportImport+14))

	errChartDomesticExportStyleTable := file.SetCellStyle(chartSheetName, "B23", "C24", formatNumberStyle)

	if errChartDomesticExportStyleTable != nil {
		return file, errChartDomesticExportStyleTable
	}

	var seriesChartDomesticExport string

	seriesChartDomesticExport += fmt.Sprintf(`{
						"name": "Amount",
						"categories": "%v!$B$23:$B$24",
						"values": "%v!$C$23:$C$24",
						"marker": {
									"symbol": "square"
								} 
				}`, chartSheetName, chartSheetName)

	valueChartDomesticExport := fmt.Sprintf(`{
	    "type": "pie",
	    "series": [%v],
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
	        "show_legend_key": true
	    },
			"plotarea": 
			{
					"show_percent": true
			},
	    "title":
	    {
	        "name": "Penjualan Dalam & Luar Negeri"
	    },
	    "show_blanks_as": "zero"
	}`, seriesChartDomesticExport)

	if err := file.AddChart(chartSheetName, "G22", valueChartDomesticExport); err != nil {
		return file, err
	}

	var startAssignment int
	startAssignment = startSaleExportImport + 14 + 3

	file.SetCellValue(sheetName, fmt.Sprintf("A%v", startAssignment), "Surat Penugasan")
	errTitleAssignment := file.SetCellStyle(sheetName, fmt.Sprintf("A%v", startAssignment), fmt.Sprintf("A%v", startAssignment), boldTitleStyle)

	if errTitleAssignment != nil {
		return file, errTitleAssignment
	}

	file.MergeCell(sheetName, fmt.Sprintf("B%v", startAssignment+1), fmt.Sprintf("C%v", startAssignment+1))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startAssignment+2), fmt.Sprintf("C%v", startAssignment+2))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startAssignment+3), fmt.Sprintf("C%v", startAssignment+3))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startAssignment+4), fmt.Sprintf("C%v", startAssignment+4))

	file.SetCellValue(sheetName, fmt.Sprintf("D%v", startAssignment+1), "Kebutuhan")
	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startAssignment+1), "Realisasi")
	file.SetCellValue(sheetName, fmt.Sprintf("F%v", startAssignment+1), "Sisa")

	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startAssignment+2), "Kelistrikan")
	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startAssignment+3), "Non Kelistrikan")

	file.SetCellValue(sheetName, fmt.Sprintf("D%v", startAssignment+2), reportSaleDetail.ElectricAssignment.Quantity)
	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startAssignment+2), reportSaleDetail.ElectricAssignment.RealizationQuantity)
	file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startAssignment+2), fmt.Sprintf("D%v-E%v", startAssignment+2, startAssignment+2))

	file.SetCellValue(sheetName, fmt.Sprintf("D%v", startAssignment+3), reportSaleDetail.CafAssignment.Quantity)
	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startAssignment+3), reportSaleDetail.CafAssignment.RealizationQuantity)
	file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startAssignment+3), fmt.Sprintf("D%v-E%v", startAssignment+3, startAssignment+3))

	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startAssignment+4), "TOTAL")
	file.SetCellFormula(sheetName, fmt.Sprintf("D%v", startAssignment+4), fmt.Sprintf("SUM(D%v:D%v)", startAssignment+2, startAssignment+3))
	file.SetCellFormula(sheetName, fmt.Sprintf("E%v", startAssignment+4), fmt.Sprintf("SUM(E%v:E%v)", startAssignment+2, startAssignment+3))
	file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startAssignment+4), fmt.Sprintf("SUM(F%v:F%v)", startAssignment+2, startAssignment+3))

	errTitleAssignment1 := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startAssignment+1), fmt.Sprintf("F%v", startAssignment+1), boldTitleTableStyle)

	if errTitleAssignment1 != nil {
		return file, errTitleAssignment1
	}

	errNumberAssignmentStyle := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startAssignment+2), fmt.Sprintf("F%v", startAssignment+3), formatNumberStyle)

	if errNumberAssignmentStyle != nil {
		return file, errNumberAssignmentStyle
	}

	errNumberTotalAssignment := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startAssignment+4), fmt.Sprintf("F%v", startAssignment+4), boldNumberStyle)

	if errNumberTotalAssignment != nil {
		return file, errNumberTotalAssignment
	}

	return file, nil
}
