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
	GetTransactionReport(iupopkId int, input TransactionReportInput, typeTransaction string) ([]TransactionReport, error)
	CreateTransactionReport(file *excelize.File, sheetName string, iupopk iupopk.Iupopk, transactionData []TransactionReport) (*excelize.File, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

var month = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

func getMonthProrate(electric RecapElectricity, nonElectric RecapNonElectricity, cement RecapCement, year string) int {

	t := time.Now()

	monthNow := t.Month()
	yearNow := t.Year()
	yearString := strconv.Itoa(yearNow)

	if yearString != year {
		return 12
	}

	if monthNow.String() == "January" {
		return 1
	}

	electricReflect := reflect.ValueOf(electric)
	nonElectricReflect := reflect.ValueOf(nonElectric)
	cementReflect := reflect.ValueOf(cement)

	valElectric := reflect.Indirect(electricReflect).FieldByName(monthNow.String()).Interface().(float64)
	valNonElectric := reflect.Indirect(nonElectricReflect).FieldByName(monthNow.String()).Interface().(float64)
	valCement := reflect.Indirect(cementReflect).FieldByName(monthNow.String()).Interface().(float64)

	if valElectric != 0 || valNonElectric != 0 || valCement != 0 {
		return int(monthNow)
	}

	return int(monthNow) - 1
}

func getProductionProrate(prorate int, prod QuantityProduction) float64 {
	var prodQuantity float64
	productionReflect := reflect.ValueOf(prod)
	for i := 0; i < prorate; i++ {

		valProd := reflect.Indirect(productionReflect).FieldByName(month[i]).Interface().(float64)

		prodQuantity += valProd
	}

	return prodQuantity
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

func insertTransactionRealization(file *excelize.File, sheetName string, transactionElectric []RealizationTransaction, transactionNonElectric []RealizationTransaction, transactionCement []RealizationTransaction) (*excelize.File, error) {
	custFmt := "#,##0.000"

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

	file.SetCellStyle(sheetName, "B7", fmt.Sprintf("F%v", 7+len(transactionElectric)+len(transactionNonElectric)+len(transactionCement)), borderCenterStyle)

	file.SetCellStyle(sheetName, "C7", fmt.Sprintf("C%v", 7+len(transactionElectric)+len(transactionNonElectric)+len(transactionCement)), borderCenterDateStyle)

	file.SetCellStyle(sheetName, "G7", fmt.Sprintf("H%v", 7+len(transactionElectric)+len(transactionNonElectric)+len(transactionCement)), borderStyle)

	file.SetCellStyle(sheetName, "J7", fmt.Sprintf("J%v", 7+len(transactionElectric)+len(transactionNonElectric)+len(transactionCement)), centerStyle)

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
				if nonElectric.EndUser.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7), "KELISTRIKAN")
				} else if nonElectric.EndUser.IndustryType.CategoryIndustryType.Name == "Semen" {
					file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7), "SEMEN")
				} else if nonElectric.EndUser.IndustryType.CategoryIndustryType.Name == "Smelter" {
					file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7), "SMELTER")
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

	for idx, cement := range transactionCement {
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", idx+7+len(transactionNonElectric)), idx+1+len(transactionNonElectric))
		date, err := time.Parse("2006-01-02", cement.ShippingDate)

		if err != nil {
			return file, err
		}
		file.SetCellValue(sheetName, fmt.Sprintf("C%v", idx+7+len(transactionNonElectric)), date)
		if cement.Trader != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", idx+7+len(transactionNonElectric)), strings.ToUpper(cement.Trader.CompanyName))
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", idx+7+len(transactionNonElectric)), "-")
		}

		if cement.EndUser != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("E%v", idx+7+len(transactionNonElectric)), strings.ToUpper(cement.EndUser.CompanyName))
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("E%v", idx+7+len(transactionNonElectric)), "-")
		}

		if cement.EndUser != nil {
			if cement.EndUser.IndustryType != nil {
				if cement.EndUser.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)), "KELISTRIKAN")
				} else if cement.EndUser.IndustryType.CategoryIndustryType.Name == "Semen" {
					file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)), "SEMEN")
				} else if cement.EndUser.IndustryType.CategoryIndustryType.Name == "Smelter" {
					file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)), "SMELTER")
				}
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)), "-")
			}
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)), "-")
		}

		file.SetCellValue(sheetName, fmt.Sprintf("G%v", idx+7+len(transactionNonElectric)), cement.QualityCaloriesAr)
		file.SetCellValue(sheetName, fmt.Sprintf("H%v", idx+7+len(transactionNonElectric)), cement.Quantity)

		if cement.IsBastOk {
			file.SetCellValue(sheetName, fmt.Sprintf("J%v", idx+7+len(transactionNonElectric)), "OK")
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("J%v", idx+7+len(transactionNonElectric)), "-")
		}
	}

	for idx, electric := range transactionElectric {
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", idx+7+len(transactionNonElectric)+len(transactionCement)), idx+1+len(transactionNonElectric)+len(transactionCement))
		date, err := time.Parse("2006-01-02", electric.ShippingDate)

		if err != nil {
			return file, err
		}
		file.SetCellValue(sheetName, fmt.Sprintf("C%v", idx+7+len(transactionNonElectric)+len(transactionCement)), date)
		if electric.Trader != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", idx+7+len(transactionNonElectric)+len(transactionCement)), strings.ToUpper(electric.Trader.CompanyName))
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", idx+7+len(transactionNonElectric)+len(transactionCement)), "-")
		}

		if electric.EndUser != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("E%v", idx+7+len(transactionNonElectric)+len(transactionCement)), strings.ToUpper(electric.EndUser.CompanyName))
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("E%v", idx+7+len(transactionNonElectric)+len(transactionCement)), "-")
		}

		if electric.EndUser != nil {
			if electric.EndUser.IndustryType != nil {
				if electric.EndUser.IndustryType.CategoryIndustryType.Name == "Kelistrikan" {
					file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)+len(transactionCement)), "KELISTRIKAN")
				} else if electric.EndUser.IndustryType.CategoryIndustryType.Name == "Semen" {
					file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)+len(transactionCement)), "SEMEN")
				} else if electric.EndUser.IndustryType.CategoryIndustryType.Name == "Smelter" {
					file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)+len(transactionCement)), "SMELTER")
				}
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)+len(transactionCement)), "-")
			}
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("F%v", idx+7+len(transactionNonElectric)+len(transactionCement)), "-")
		}

		file.SetCellValue(sheetName, fmt.Sprintf("G%v", idx+7+len(transactionNonElectric)+len(transactionCement)), electric.QualityCaloriesAr)
		file.SetCellValue(sheetName, fmt.Sprintf("H%v", idx+7+len(transactionNonElectric)+len(transactionCement)), electric.Quantity)

		if electric.IsBastOk {
			file.SetCellValue(sheetName, fmt.Sprintf("J%v", idx+7+len(transactionNonElectric)+len(transactionCement)), "OK")
		} else {
			file.SetCellValue(sheetName, fmt.Sprintf("J%v", idx+7+len(transactionNonElectric)+len(transactionCement)), "-")
		}
	}

	file.SetCellValue(sheetName, fmt.Sprintf("G%v", 7+len(transactionElectric)+len(transactionNonElectric)+len(transactionCement)), "Jumlah")

	file.SetCellFormula(sheetName, fmt.Sprintf("H%v", 7+len(transactionElectric)+len(transactionNonElectric)+len(transactionCement)), fmt.Sprintf("SUM(H7:H%v)", 6+len(transactionElectric)+len(transactionNonElectric)+len(transactionCement)))

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
		"C5":  "Smelter",
		"D5":  "Semen",
		"E5":  "Jumlah",
		"G5":  "Produksi",
		"H5":  "Tidak bisa Claim DMO",
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
			values["E6"] = recapElectricity.Field(i).Interface().(float64)
		case "February":
			values["B7"] = recapElectricity.Field(i).Interface().(float64)
			values["E7"] = recapElectricity.Field(i).Interface().(float64)

		case "March":
			values["B8"] = recapElectricity.Field(i).Interface().(float64)
			values["E8"] = recapElectricity.Field(i).Interface().(float64)

		case "April":
			values["B9"] = recapElectricity.Field(i).Interface().(float64)
			values["E9"] = recapElectricity.Field(i).Interface().(float64)

		case "May":
			values["B10"] = recapElectricity.Field(i).Interface().(float64)
			values["E10"] = recapElectricity.Field(i).Interface().(float64)

		case "June":
			values["B11"] = recapElectricity.Field(i).Interface().(float64)
			values["E11"] = recapElectricity.Field(i).Interface().(float64)

		case "July":
			values["B12"] = recapElectricity.Field(i).Interface().(float64)
			values["E12"] = recapElectricity.Field(i).Interface().(float64)

		case "August":
			values["B13"] = recapElectricity.Field(i).Interface().(float64)
			values["E13"] = recapElectricity.Field(i).Interface().(float64)

		case "September":
			values["B14"] = recapElectricity.Field(i).Interface().(float64)
			values["E14"] = recapElectricity.Field(i).Interface().(float64)

		case "October":
			values["B15"] = recapElectricity.Field(i).Interface().(float64)
			values["E15"] = recapElectricity.Field(i).Interface().(float64)

		case "November":
			values["B16"] = recapElectricity.Field(i).Interface().(float64)
			values["E16"] = recapElectricity.Field(i).Interface().(float64)

		case "December":
			values["B17"] = recapElectricity.Field(i).Interface().(float64)
			values["E17"] = recapElectricity.Field(i).Interface().(float64)

		case "Total":
			values["B18"] = recapElectricity.Field(i).Interface().(float64)
			values["E18"] = recapElectricity.Field(i).Interface().(float64)

		}
	}

	recapNonElectricity := reflect.ValueOf(reportRecapDmo.RecapNonElectricity)

	for i := 0; i < recapNonElectricity.NumField(); i++ {
		switch recapNonElectricity.Type().Field(i).Name {
		case "January":
			values["C6"] = recapNonElectricity.Field(i).Interface().(float64)

		case "February":
			values["C7"] = recapNonElectricity.Field(i).Interface().(float64)

		case "March":
			values["C8"] = recapNonElectricity.Field(i).Interface().(float64)

		case "April":
			values["C9"] = recapNonElectricity.Field(i).Interface().(float64)

		case "May":
			values["C10"] = recapNonElectricity.Field(i).Interface().(float64)

		case "June":
			values["C11"] = recapNonElectricity.Field(i).Interface().(float64)

		case "July":
			values["C12"] = recapNonElectricity.Field(i).Interface().(float64)

		case "August":
			values["C13"] = recapNonElectricity.Field(i).Interface().(float64)

		case "September":
			values["C14"] = recapNonElectricity.Field(i).Interface().(float64)

		case "October":
			values["C15"] = recapNonElectricity.Field(i).Interface().(float64)

		case "November":
			values["C16"] = recapNonElectricity.Field(i).Interface().(float64)

		case "December":
			values["C17"] = recapNonElectricity.Field(i).Interface().(float64)
		}
	}

	recapCement := reflect.ValueOf(reportRecapDmo.RecapCement)

	for i := 0; i < recapCement.NumField(); i++ {
		switch recapCement.Type().Field(i).Name {
		case "January":
			values["D6"] = recapCement.Field(i).Interface().(float64)

		case "February":
			values["D7"] = recapCement.Field(i).Interface().(float64)

		case "March":
			values["D8"] = recapCement.Field(i).Interface().(float64)

		case "April":
			values["D9"] = recapCement.Field(i).Interface().(float64)

		case "May":
			values["D10"] = recapCement.Field(i).Interface().(float64)

		case "June":
			values["D11"] = recapCement.Field(i).Interface().(float64)

		case "July":
			values["D12"] = recapCement.Field(i).Interface().(float64)

		case "August":
			values["D13"] = recapCement.Field(i).Interface().(float64)

		case "September":
			values["D14"] = recapCement.Field(i).Interface().(float64)

		case "October":
			values["D15"] = recapCement.Field(i).Interface().(float64)

		case "November":
			values["D16"] = recapCement.Field(i).Interface().(float64)

		case "December":
			values["D17"] = recapCement.Field(i).Interface().(float64)
		}
	}

	recapProduction := reflect.ValueOf(reportRecapDmo.Production)

	for i := 0; i < recapProduction.NumField(); i++ {
		switch recapProduction.Type().Field(i).Name {
		case "January":
			values["G6"] = recapProduction.Field(i).Interface().(float64)
		case "February":
			values["G7"] = recapProduction.Field(i).Interface().(float64)
		case "March":
			values["G8"] = recapProduction.Field(i).Interface().(float64)
		case "April":
			values["G9"] = recapProduction.Field(i).Interface().(float64)
		case "May":
			values["G10"] = recapProduction.Field(i).Interface().(float64)
		case "June":
			values["G11"] = recapProduction.Field(i).Interface().(float64)
		case "July":
			values["G12"] = recapProduction.Field(i).Interface().(float64)
		case "August":
			values["G13"] = recapProduction.Field(i).Interface().(float64)
		case "September":
			values["G14"] = recapProduction.Field(i).Interface().(float64)
		case "October":
			values["G15"] = recapProduction.Field(i).Interface().(float64)
		case "November":
			values["G16"] = recapProduction.Field(i).Interface().(float64)
		case "December":
			values["G17"] = recapProduction.Field(i).Interface().(float64)
		}
	}

	recapNotClaim := reflect.ValueOf(reportRecapDmo.NotClaimable)

	for i := 0; i < recapNotClaim.NumField(); i++ {
		switch recapNotClaim.Type().Field(i).Name {
		case "January":
			values["H6"] = recapNotClaim.Field(i).Interface().(float64)
		case "February":
			values["H7"] = recapNotClaim.Field(i).Interface().(float64)
		case "March":
			values["H8"] = recapNotClaim.Field(i).Interface().(float64)
		case "April":
			values["H9"] = recapNotClaim.Field(i).Interface().(float64)
		case "May":
			values["H10"] = recapNotClaim.Field(i).Interface().(float64)
		case "June":
			values["H11"] = recapNotClaim.Field(i).Interface().(float64)
		case "July":
			values["H12"] = recapNotClaim.Field(i).Interface().(float64)
		case "August":
			values["H13"] = recapNotClaim.Field(i).Interface().(float64)
		case "September":
			values["H14"] = recapNotClaim.Field(i).Interface().(float64)
		case "October":
			values["H15"] = recapNotClaim.Field(i).Interface().(float64)
		case "November":
			values["H16"] = recapNotClaim.Field(i).Interface().(float64)
		case "December":
			values["H17"] = recapNotClaim.Field(i).Interface().(float64)
		}
	}

	for k, v := range values {
		file.SetCellValue(sheetName, k, v)
	}

	file.SetCellFormula(sheetName, "E6", "SUM(B6:D6)")
	file.SetCellFormula(sheetName, "E7", "SUM(B7:D7)")
	file.SetCellFormula(sheetName, "E8", "SUM(B8:D8)")
	file.SetCellFormula(sheetName, "E9", "SUM(B9:D9)")
	file.SetCellFormula(sheetName, "E10", "SUM(B10:D10)")
	file.SetCellFormula(sheetName, "E11", "SUM(B11:D11)")
	file.SetCellFormula(sheetName, "E12", "SUM(B12:D12)")
	file.SetCellFormula(sheetName, "E13", "SUM(B13:D13)")
	file.SetCellFormula(sheetName, "E14", "SUM(B14:D14)")
	file.SetCellFormula(sheetName, "E15", "SUM(B15:D15)")
	file.SetCellFormula(sheetName, "E16", "SUM(B16:D16)")
	file.SetCellFormula(sheetName, "E17", "SUM(B17:D17)")

	file.SetCellFormula(sheetName, "B18", "SUM(B6:B17)")
	file.SetCellFormula(sheetName, "C18", "SUM(C6:C17)")
	file.SetCellFormula(sheetName, "D18", "SUM(D6:D17)")
	file.SetCellFormula(sheetName, "E18", "SUM(E6:E17)")
	file.SetCellFormula(sheetName, "G18", "SUM(G6:G17)")
	file.SetCellFormula(sheetName, "H18", "SUM(H6:H17)")

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

	errStyleBorderBold1 := file.SetCellStyle(sheetName, "A5", "F5", boldStyle)

	if errStyleBorderBold1 != nil {
		return file, errStyleBorderBold1
	}

	errStyleBorderBold2 := file.SetCellStyle(sheetName, "G5", "H5", boldStyle)

	if errStyleBorderBold2 != nil {
		return file, errStyleBorderBold2
	}

	errFmtNmbr1 := file.SetCellStyle(sheetName, "A6", "E17", borderStyle)

	if errFmtNmbr1 != nil {
		return file, errFmtNmbr1
	}

	errFmtNmbr2 := file.SetCellStyle(sheetName, "G6", "H17", borderStyle)

	if errFmtNmbr2 != nil {
		return file, errFmtNmbr2
	}

	errBoldNmbr1 := file.SetCellStyle(sheetName, "A18", "E18", boldNumberStyle)

	if errBoldNmbr1 != nil {
		return file, errBoldNmbr1
	}

	errBoldNmbr2 := file.SetCellStyle(sheetName, "G18", "H18", boldNumberStyle)

	if errBoldNmbr2 != nil {
		return file, errBoldNmbr2
	}

	errCustomNmbrRight1 := file.SetCellStyle(sheetName, "B6", "E17", customNumberRightStyle)

	if errCustomNmbrRight1 != nil {
		return file, errCustomNmbrRight1
	}

	errCustomNmbrRight2 := file.SetCellStyle(sheetName, "G6", "H17", customNumberRightStyle)

	if errCustomNmbrRight2 != nil {
		return file, errCustomNmbrRight2
	}

	file.SetColWidth(sheetName, "A", "A", float64(15))
	file.SetColWidth(sheetName, "B", "E", float64(25))
	file.SetColWidth(sheetName, "F", "F", float64(5))
	file.SetColWidth(sheetName, "G", "G", float64(30))
	file.SetColWidth(sheetName, "H", "H", float64(25))

	mergeErr4 := file.MergeCell(sheetName, "A20", "C20")
	if mergeErr4 != nil {
		return file, mergeErr4
	}

	file.SetCellValue(sheetName, "A20", "% Pemenuhan DMO terhadap REALISASI PRODUKSI")

	file.SetCellFormula(sheetName, "E20", fmt.Sprintf("D18/SUM(G6:G%v)", 6+getMonthProrate(reportRecapDmo.RecapElectricity, reportRecapDmo.RecapNonElectricity, reportRecapDmo.RecapCement, year)-1))

	percentageBoldStyle, _ := file.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
		},
		NumFmt: 10,
	})

	errPercentageBoldStyle := file.SetCellStyle(sheetName, "A20", "H60", percentageBoldStyle)
	errBoldNmbrOnly1 := file.SetCellStyle(sheetName, "G20", "G60", boldNumberOnlyStyle)

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
		file.SetCellValue(sheetName, fmt.Sprintf("E%v", 23+(8*idx)), withCommaThousandSep)
		file.SetCellStyle(sheetName, fmt.Sprintf("A%v", 23+(8*idx)), fmt.Sprintf("G%v", 23+(8*idx)), normalStyle)

		file.SetCellValue(sheetName, fmt.Sprintf("A%v", 24+(8*idx)), "% Pemenuhan DMO terhadap RENCANA PRODUKSI")
		file.SetCellFormula(sheetName, fmt.Sprintf("E%v", 24+(8*idx)), fmt.Sprintf("E18/G%v", 24+(8*idx)))
		file.SetCellValue(sheetName, fmt.Sprintf("G%v", 24+(8*idx)), rkab.ProductionQuota)
		file.SetCellValue(sheetName, fmt.Sprintf("H%v", 24+(8*idx)), "Quota RKAB")

		file.SetCellValue(sheetName, fmt.Sprintf("A%v", 25+(8*idx)), fmt.Sprintf("disetujui tgl %s", dateFormat.Format("DD MMMM YYYY", "id")))

		file.SetCellValue(sheetName, fmt.Sprintf("A%v", 27+(8*idx)), fmt.Sprintf("%% Pemenuhan DMO terhadap kewajiban pemenuhan DMO %.0f%%", rkab.DmoObligation))
		file.SetCellFormula(sheetName, fmt.Sprintf("E%v", 27+(8*idx)), fmt.Sprintf("D18/G%v", 27+(8*idx)))
		file.SetCellValue(sheetName, fmt.Sprintf("G%v", 27+(8*idx)), maxDmoObligationQuota)

		file.SetCellValue(sheetName, fmt.Sprintf("A%v", 29+(8*idx)), fmt.Sprintf("%% Pemenuhan DMO terhadap Rencana Produksi (Prorata %v bulan)", getMonthProrate(reportRecapDmo.RecapElectricity, reportRecapDmo.RecapNonElectricity, reportRecapDmo.RecapCement, year)))
		file.SetCellFormula(sheetName, fmt.Sprintf("E%v", 29+(8*idx)), fmt.Sprintf("D18/G%v", 29+(8*idx)))
		file.SetCellValue(sheetName, fmt.Sprintf("G%v", 29+(8*idx)), rkab.ProductionQuota*float64(getMonthProrate(reportRecapDmo.RecapElectricity, reportRecapDmo.RecapNonElectricity, reportRecapDmo.RecapCement, year))/float64(12))
		file.SetCellValue(sheetName, fmt.Sprintf("H%v", 29+(8*idx)), fmt.Sprintf("prorata %v bulan", getMonthProrate(reportRecapDmo.RecapElectricity, reportRecapDmo.RecapNonElectricity, reportRecapDmo.RecapCement, year)))

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
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.January, reportRealization.NonElectric.January, reportRealization.Cement.January)

			if errFile != nil {
				return insertedFile, errFile
			}
		case "FEB":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.February, reportRealization.NonElectric.February, reportRealization.Cement.February)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "MAR":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.March, reportRealization.NonElectric.March, reportRealization.Cement.March)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "APR":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.April, reportRealization.NonElectric.April, reportRealization.Cement.April)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "MEI":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.May, reportRealization.NonElectric.May, reportRealization.Cement.May)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "JUN":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.June, reportRealization.NonElectric.June, reportRealization.Cement.June)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "JUL":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.July, reportRealization.NonElectric.July, reportRealization.Cement.July)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "AGU":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.August, reportRealization.NonElectric.August, reportRealization.Cement.August)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "SEP":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.September, reportRealization.NonElectric.September, reportRealization.Cement.September)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "OKT":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.October, reportRealization.NonElectric.October, reportRealization.Cement.October)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "NOV":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.November, reportRealization.NonElectric.November, reportRealization.Cement.November)
			if errFile != nil {
				return insertedFile, errFile
			}
		case "DES":
			insertedFile, errFile = insertTransactionRealization(insertedFile, v, reportRealization.Electric.December, reportRealization.NonElectric.December, reportRealization.Cement.December)
			if errFile != nil {
				return insertedFile, errFile
			}
		}
	}

	return insertedFile, nil
}

func (s *service) CreateReportSalesDetail(year string, reportSaleDetail SaleDetail, iupopk iupopk.Iupopk, file *excelize.File, sheetName string, chartSheetName string) (*excelize.File, error) {

	custFmt := "#,##0.000"
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
	file.SetColWidth(sheetName, "I", "I", float64(20))

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
	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+1), "Semen")
	file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+1), "Smelter")
	file.SetCellValue(sheetName, fmt.Sprintf("G%v", startSale+1), "Jumlah")
	file.SetCellValue(sheetName, fmt.Sprintf("I%v", startSale+1), "Tidak Bisa Claim DMO")

	file.SetCellValue(chartSheetName, "B41", "Bulan")
	file.SetCellValue(chartSheetName, "C41", "Kelistrikan")
	file.SetCellValue(chartSheetName, "D41", "Semen")
	file.SetCellValue(chartSheetName, "E41", "Smelter")
	file.SetCellValue(chartSheetName, "F41", "Jumlah")

	errChartTitleRecap := file.SetCellStyle(chartSheetName, "B41", "F41", boldTitleTableStyle)

	if errChartTitleRecap != nil {
		return file, errChartTitleRecap
	}

	errTitleTablePenjualan := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startSale+1), fmt.Sprintf("G%v", startSale+1), boldTitleTableStyle)

	if errTitleTablePenjualan != nil {
		return file, errTitleTablePenjualan
	}

	errTitleTablePenjualan2 := file.SetCellStyle(sheetName, fmt.Sprintf("I%v", startSale+1), fmt.Sprintf("I%v", startSale+1), boldTitleTableStyle)

	if errTitleTablePenjualan2 != nil {
		return file, errTitleTablePenjualan2
	}

	for idx, v := range monthString {
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+2+idx), v)
		file.SetCellValue(chartSheetName, fmt.Sprintf("B%v", 42+idx), v)
		switch v {
		case "January":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.January)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapCement.January)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.January+reportSaleDetail.RecapElectricity.January)

			file.SetCellValue(sheetName, fmt.Sprintf("G%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.January+reportSaleDetail.RecapElectricity.January+reportSaleDetail.RecapCement.January)

			file.SetCellValue(sheetName, fmt.Sprintf("I%v", startSale+2+idx), reportSaleDetail.NotClaimable.January)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.January)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapCement.January)

			file.SetCellValue(chartSheetName, fmt.Sprintf("E%v", 42+idx), reportSaleDetail.RecapNonElectricity.January)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("F%v", 42+idx), fmt.Sprintf("SUM(C%v:E%v)", 42+idx, 42+idx))

		case "February":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.February)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapCement.February)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.February+reportSaleDetail.RecapElectricity.February)

			file.SetCellValue(sheetName, fmt.Sprintf("G%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.February+reportSaleDetail.RecapElectricity.February+reportSaleDetail.RecapCement.February)

			file.SetCellValue(sheetName, fmt.Sprintf("I%v", startSale+2+idx), reportSaleDetail.NotClaimable.February)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.February)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapCement.February)

			file.SetCellValue(chartSheetName, fmt.Sprintf("E%v", 42+idx), reportSaleDetail.RecapNonElectricity.February)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("F%v", 42+idx), fmt.Sprintf("SUM(C%v:E%v)", 42+idx, 42+idx))

		case "March":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.March)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapCement.March)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.March+reportSaleDetail.RecapElectricity.March)

			file.SetCellValue(sheetName, fmt.Sprintf("G%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.March+reportSaleDetail.RecapElectricity.March+reportSaleDetail.RecapCement.March)

			file.SetCellValue(sheetName, fmt.Sprintf("I%v", startSale+2+idx), reportSaleDetail.NotClaimable.March)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.March)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapCement.March)

			file.SetCellValue(chartSheetName, fmt.Sprintf("E%v", 42+idx), reportSaleDetail.RecapNonElectricity.March)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("F%v", 42+idx), fmt.Sprintf("SUM(C%v:E%v)", 42+idx, 42+idx))

		case "April":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.April)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapCement.April)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.April+reportSaleDetail.RecapElectricity.April)

			file.SetCellValue(sheetName, fmt.Sprintf("G%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.April+reportSaleDetail.RecapElectricity.April+reportSaleDetail.RecapCement.April)

			file.SetCellValue(sheetName, fmt.Sprintf("I%v", startSale+2+idx), reportSaleDetail.NotClaimable.April)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.April)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapCement.April)

			file.SetCellValue(chartSheetName, fmt.Sprintf("E%v", 42+idx), reportSaleDetail.RecapNonElectricity.April)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("F%v", 42+idx), fmt.Sprintf("SUM(C%v:E%v)", 42+idx, 42+idx))

		case "May":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.May)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapCement.May)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.May+reportSaleDetail.RecapElectricity.May)

			file.SetCellValue(sheetName, fmt.Sprintf("G%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.May+reportSaleDetail.RecapElectricity.May+reportSaleDetail.RecapCement.May)

			file.SetCellValue(sheetName, fmt.Sprintf("I%v", startSale+2+idx), reportSaleDetail.NotClaimable.May)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.May)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapCement.May)

			file.SetCellValue(chartSheetName, fmt.Sprintf("E%v", 42+idx), reportSaleDetail.RecapNonElectricity.May)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("F%v", 42+idx), fmt.Sprintf("SUM(C%v:E%v)", 42+idx, 42+idx))

		case "June":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.June)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapCement.June)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.June+reportSaleDetail.RecapElectricity.June)

			file.SetCellValue(sheetName, fmt.Sprintf("G%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.June+reportSaleDetail.RecapElectricity.June+reportSaleDetail.RecapCement.June)

			file.SetCellValue(sheetName, fmt.Sprintf("I%v", startSale+2+idx), reportSaleDetail.NotClaimable.June)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.June)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapCement.June)

			file.SetCellValue(chartSheetName, fmt.Sprintf("E%v", 42+idx), reportSaleDetail.RecapNonElectricity.June)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("F%v", 42+idx), fmt.Sprintf("SUM(C%v:E%v)", 42+idx, 42+idx))

		case "July":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.July)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapCement.July)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.July+reportSaleDetail.RecapElectricity.July)

			file.SetCellValue(sheetName, fmt.Sprintf("G%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.July+reportSaleDetail.RecapElectricity.July+reportSaleDetail.RecapCement.July)

			file.SetCellValue(sheetName, fmt.Sprintf("I%v", startSale+2+idx), reportSaleDetail.NotClaimable.July)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.July)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapCement.July)

			file.SetCellValue(chartSheetName, fmt.Sprintf("E%v", 42+idx), reportSaleDetail.RecapNonElectricity.July)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("F%v", 42+idx), fmt.Sprintf("SUM(C%v:E%v)", 42+idx, 42+idx))

		case "August":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.August)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapCement.August)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.August+reportSaleDetail.RecapElectricity.August)

			file.SetCellValue(sheetName, fmt.Sprintf("G%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.August+reportSaleDetail.RecapElectricity.August+reportSaleDetail.RecapCement.August)

			file.SetCellValue(sheetName, fmt.Sprintf("I%v", startSale+2+idx), reportSaleDetail.NotClaimable.August)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.August)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapCement.August)

			file.SetCellValue(chartSheetName, fmt.Sprintf("E%v", 42+idx), reportSaleDetail.RecapNonElectricity.August)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("F%v", 42+idx), fmt.Sprintf("SUM(C%v:E%v)", 42+idx, 42+idx))

		case "September":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.September)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapCement.September)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.September+reportSaleDetail.RecapElectricity.September)

			file.SetCellValue(sheetName, fmt.Sprintf("G%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.September+reportSaleDetail.RecapElectricity.September+reportSaleDetail.RecapCement.September)

			file.SetCellValue(sheetName, fmt.Sprintf("I%v", startSale+2+idx), reportSaleDetail.NotClaimable.September)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.September)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapCement.September)

			file.SetCellValue(chartSheetName, fmt.Sprintf("E%v", 42+idx), reportSaleDetail.RecapNonElectricity.September)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("F%v", 42+idx), fmt.Sprintf("SUM(C%v:E%v)", 42+idx, 42+idx))

		case "October":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.October)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapCement.October)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.October+reportSaleDetail.RecapElectricity.October)

			file.SetCellValue(sheetName, fmt.Sprintf("G%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.October+reportSaleDetail.RecapElectricity.October+reportSaleDetail.RecapCement.October)

			file.SetCellValue(sheetName, fmt.Sprintf("I%v", startSale+2+idx), reportSaleDetail.NotClaimable.October)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.October)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapCement.October)

			file.SetCellValue(chartSheetName, fmt.Sprintf("E%v", 42+idx), reportSaleDetail.RecapNonElectricity.October)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("F%v", 42+idx), fmt.Sprintf("SUM(C%v:E%v)", 42+idx, 42+idx))
		case "November":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.November)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapCement.November)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.November+reportSaleDetail.RecapElectricity.November)

			file.SetCellValue(sheetName, fmt.Sprintf("G%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.November+reportSaleDetail.RecapElectricity.November+reportSaleDetail.RecapCement.November)

			file.SetCellValue(sheetName, fmt.Sprintf("I%v", startSale+2+idx), reportSaleDetail.NotClaimable.November)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.November)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapCement.November)

			file.SetCellValue(chartSheetName, fmt.Sprintf("E%v", 42+idx), reportSaleDetail.RecapNonElectricity.November)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("F%v", 42+idx), fmt.Sprintf("SUM(C%v:E%v)", 42+idx, 42+idx))

		case "December":

			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startSale+2+idx), reportSaleDetail.RecapElectricity.December)

			file.SetCellValue(sheetName, fmt.Sprintf("E%v", startSale+2+idx), reportSaleDetail.RecapCement.December)

			file.SetCellValue(sheetName, fmt.Sprintf("F%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.December+reportSaleDetail.RecapElectricity.December)

			file.SetCellValue(sheetName, fmt.Sprintf("G%v", startSale+2+idx), reportSaleDetail.RecapNonElectricity.December+reportSaleDetail.RecapElectricity.December+reportSaleDetail.RecapCement.December)

			file.SetCellValue(sheetName, fmt.Sprintf("I%v", startSale+2+idx), reportSaleDetail.NotClaimable.December)

			file.SetCellValue(chartSheetName, fmt.Sprintf("C%v", 42+idx), reportSaleDetail.RecapElectricity.December)

			file.SetCellValue(chartSheetName, fmt.Sprintf("D%v", 42+idx), reportSaleDetail.RecapCement.December)

			file.SetCellValue(chartSheetName, fmt.Sprintf("E%v", 42+idx), reportSaleDetail.RecapNonElectricity.December)

			file.SetCellFormula(chartSheetName, fmt.Sprintf("F%v", 42+idx), fmt.Sprintf("SUM(C%v:E%v)", 42+idx, 42+idx))
		}
	}

	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+14), "TOTAL")
	file.SetCellFormula(sheetName, fmt.Sprintf("D%v", startSale+14), fmt.Sprintf("SUM(D%v:D%v)", startSale+2, startSale+13))
	file.SetCellFormula(sheetName, fmt.Sprintf("E%v", startSale+14), fmt.Sprintf("SUM(E%v:E%v)", startSale+2, startSale+13))
	file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startSale+14), fmt.Sprintf("SUM(F%v:F%v)", startSale+2, startSale+13))
	file.SetCellFormula(sheetName, fmt.Sprintf("G%v", startSale+14), fmt.Sprintf("SUM(G%v:G%v)", startSale+2, startSale+13))
	file.SetCellFormula(sheetName, fmt.Sprintf("I%v", startSale+14), fmt.Sprintf("SUM(I%v:I%v)", startSale+2, startSale+13))

	file.SetCellValue(chartSheetName, "B54", "TOTAL")
	file.SetCellFormula(chartSheetName, "C54", "SUM(C42:C53)")
	file.SetCellFormula(chartSheetName, "D54", "SUM(D42:D53)")
	file.SetCellFormula(chartSheetName, "E54", "SUM(E42:E53)")
	file.SetCellFormula(chartSheetName, "F54", "SUM(F42:F53)")

	errChartRecapStyleTable := file.SetCellStyle(chartSheetName, "B42", "F53", formatNumberStyle)

	if errChartRecapStyleTable != nil {
		return file, errChartRecapStyleTable
	}

	errChartTotalRecapStyleTable := file.SetCellStyle(chartSheetName, "B54", "F54", boldNumberStyle)

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
				},`, chartSheetName, 41, chartSheetName, chartSheetName)

	seriesChartRecapSale += fmt.Sprintf(`{
					"name": "%v!$E$%v",
					"categories": "%v!$B$42:$B$53",
					"values": "%v!$E$42:$E$53",
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

	errPenjualanTotalNumberStyle := file.SetCellStyle(sheetName, fmt.Sprintf("D%v", startSale+14), fmt.Sprintf("G%v", startSale+14), boldNumberStyle)

	if errPenjualanTotalNumberStyle != nil {
		return file, errPenjualanTotalNumberStyle
	}

	errPenjualanTotalNumberStyle2 := file.SetCellStyle(sheetName, fmt.Sprintf("I%v", startSale+14), fmt.Sprintf("I%v", startSale+14), boldNumberStyle)

	if errPenjualanTotalNumberStyle2 != nil {
		return file, errPenjualanTotalNumberStyle2
	}

	erPenjualanMonthStyleTable := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startSale+2), fmt.Sprintf("B%v", startSale+13), borderStyle)

	if erPenjualanMonthStyleTable != nil {
		return file, erPenjualanMonthStyleTable
	}

	errPenjualanNumberStyle := file.SetCellStyle(sheetName, fmt.Sprintf("D%v", startSale+2), fmt.Sprintf("G%v", startSale+13), formatNumberStyle)

	if errPenjualanNumberStyle != nil {
		return file, errPenjualanNumberStyle
	}

	errPenjualanNumberStyle2 := file.SetCellStyle(sheetName, fmt.Sprintf("I%v", startSale+2), fmt.Sprintf("I%v", startSale+13), formatNumberStyle)

	if errPenjualanNumberStyle2 != nil {
		return file, errPenjualanNumberStyle2
	}

	for idx, rkab := range reportSaleDetail.Rkabs {
		dateFormat, errDate := goment.New(rkab.DateOfIssue)
		if errDate != nil {
			return file, errDate
		}
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+16+(idx*10)), "% Pemenuhan DMO terhadap REALISASI PRODUKSI")
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+18+(idx*10)), "% Pemenuhan DMO terhadap RENCANA PRODUKSI")
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+19+(idx*10)), fmt.Sprintf("disetujui tgl %s", dateFormat.Format("DD MMMM YYYY", "id")))
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+21+(idx*10)), fmt.Sprintf("%% Pemenuhan DMO terhadap kewajiban pemenuhan DMO %.0f%%", rkab.DmoObligation))
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startSale+23+(idx*10)), fmt.Sprintf("%% Pemenuhan DMO terhadap Rencana Produksi (Prorata %v bulan)", getMonthProrate(reportSaleDetail.RecapElectricity, reportSaleDetail.RecapNonElectricity, reportSaleDetail.RecapCement, year)))

		errBoldOnlyStyle := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startSale+16+(idx*10)), fmt.Sprintf("B%v", startSale+23+(idx*10)), boldOnlyStyle)

		if errBoldOnlyStyle != nil {
			return file, errBoldOnlyStyle
		}

		file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startSale+16+(idx*10)), fmt.Sprintf("G%v/SUM(D%v:D%v)", startSale+14, startProduction+2, startProduction+2+getMonthProrate(reportSaleDetail.RecapElectricity, reportSaleDetail.RecapNonElectricity, reportSaleDetail.RecapCement, year)-1))
		file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startSale+18+(idx*10)), fmt.Sprintf("G%v/D%v", startSale+14, 8+idx*4))
		file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startSale+21+(idx*10)), fmt.Sprintf("G%v/D%v", startSale+14, startRkab+1))
		file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startSale+23+(idx*10)), fmt.Sprintf("G%v/(D%v*%v/12)", startSale+14, 8+idx*4, getMonthProrate(reportSaleDetail.RecapElectricity, reportSaleDetail.RecapNonElectricity, reportSaleDetail.RecapCement, year)))

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
	file.SetCellFormula(chartSheetName, "C79", fmt.Sprintf("F54"))
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

	file.SetCellValue(chartSheetName, "B126", fmt.Sprintf("Pemenuhan DMO terhadap Rencana Produksi (Prorata %v bulan)", getMonthProrate(reportSaleDetail.RecapElectricity, reportSaleDetail.RecapNonElectricity, reportSaleDetail.RecapCement, year)))

	errTitleDmo4 := file.SetCellStyle(chartSheetName, "B126", "B126", boldOnlyStyle)

	if errTitleDmo4 != nil {
		return file, errTitleDmo4
	}

	file.SetCellValue(chartSheetName, "B127", "Rencana Produksi")

	file.SetCellValue(chartSheetName, "B129", "Prorata Produksi")
	file.SetCellValue(chartSheetName, "B130", "Realisasi DMO")
	file.SetCellValue(chartSheetName, "B131", "% Pemenuhan DMO terhadap Rencana Produksi (Prorata)")

	file.SetCellFormula(chartSheetName, "C127", fmt.Sprintf("%s!E21", sheetName))

	file.SetCellFormula(chartSheetName, "C129", fmt.Sprintf("%s!E21*%v/12", sheetName, getMonthProrate(reportSaleDetail.RecapElectricity, reportSaleDetail.RecapNonElectricity, reportSaleDetail.RecapCement, year)))

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
						"name": "Pemenuhan DMO Terhadap Rencana Produksi Prorata (%v Bulan)",
						"categories": "%v!$B$129:$B$130",
						"values": "%v!$C$129:$C$130",
						"marker": {
									"symbol": "square"
								} 
				}`, getMonthProrate(reportSaleDetail.RecapElectricity, reportSaleDetail.RecapNonElectricity, reportSaleDetail.RecapCement, year), chartSheetName, chartSheetName)

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
	        "name": "Pemenuhan DMO terhadap Rencana Produksi (Prorata %v bulan)"
	    },
	    "show_blanks_as": "zero"
	}`, seriesChartPemenuhanDmoTerhadapRencanaProduksiProrata, getMonthProrate(reportSaleDetail.RecapElectricity, reportSaleDetail.RecapNonElectricity, reportSaleDetail.RecapCement, year))

	if err := file.AddChart(chartSheetName, "G126", valueChartPemenuhanDmoTerhadapRencanaProduksiProrata); err != nil {
		return file, err
	}

	file.SetCellValue(chartSheetName, "B143", "REKAP TOTAL DMO BERDASARKAN JENIS INDUSTRI")

	errTitleDmo5 := file.SetCellStyle(chartSheetName, "B143", "B143", boldOnlyStyle)

	if errTitleDmo5 != nil {
		return file, errTitleDmo5
	}

	file.SetCellValue(chartSheetName, "B144", "Kelistrikan")
	file.SetCellValue(chartSheetName, "B145", "Semen")
	file.SetCellValue(chartSheetName, "B146", "Smelter")

	file.SetCellFormula(chartSheetName, "C144", "C54")
	file.SetCellFormula(chartSheetName, "C145", "D54")
	file.SetCellFormula(chartSheetName, "C146", "E54")

	errChartRecapJenisIndustri := file.SetCellStyle(chartSheetName, "B144", "C146", formatNumberStyle)

	if errChartRecapJenisIndustri != nil {
		return file, errChartRecapJenisIndustri
	}

	var seriesChartRecapJenisIndustri string

	seriesChartRecapJenisIndustri += fmt.Sprintf(`{
						"name": "REKAP TOTAL DMO BERDASARKAN JENIS INDUSTRI",
						"categories": "%v!$B$144:$B$146",
						"values": "%v!$C$144:$C$146",
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
	        "name": "Penjualan Berdasarkan Jenis Industri"
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

	errNumberElectricStyle := file.SetCellStyle(sheetName, fmt.Sprintf("E%v", startDetail+4), fmt.Sprintf("P%v", startDetail+3+numberElectric), formatNumberStyle)

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

	file.SetCellValue(sheetName, fmt.Sprintf("A%v", startDetailNonElectric), "Data Realisasi Penjualan Batubara Untuk Memenuhi Kebutuhan Batubara Untuk Industri Smelter")

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
	file.SetCellValue(chartSheetName, fmt.Sprintf("B%v", 178+countIndustry), "DETAIL REKAP TOTAL DMO BERDASARKAN JENIS INDUSTRI - Smelter")

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
						"name": "Sales Smelter",
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
	        "name": "Penjualan berdasarkan JENIS INDUSTRI (Smelter)"
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

	file.SetCellFormula(sheetName, fmt.Sprintf("E%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(E%v:E%v)", startDetailNonElectric+2, startDetailNonElectric+1+numberNonElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(F%v:F%v)", startDetailNonElectric+2, startDetailNonElectric+1+numberNonElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("G%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(G%v:G%v)", startDetailNonElectric+2, startDetailNonElectric+1+numberNonElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("H%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(H%v:H%v)", startDetailNonElectric+2, startDetailNonElectric+1+numberNonElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("I%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(I%v:I%v)", startDetailNonElectric+2, startDetailNonElectric+1+numberNonElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("J%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(J%v:J%v)", startDetailNonElectric+2, startDetailNonElectric+1+numberNonElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("K%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(K%v:K%v)", startDetailNonElectric+2, startDetailNonElectric+1+numberNonElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("L%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(L%v:L%v)", startDetailNonElectric+2, startDetailNonElectric+1+numberNonElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("M%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(M%v:M%v)", startDetailNonElectric+2, startDetailNonElectric+1+numberNonElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("N%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(N%v:N%v)", startDetailNonElectric+2, startDetailNonElectric+1+numberNonElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("O%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(O%v:O%v)", startDetailNonElectric+2, startDetailNonElectric+1+numberNonElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("P%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(P%v:P%v)", startDetailNonElectric+2, startDetailNonElectric+1+numberNonElectric))

	file.SetCellFormula(sheetName, fmt.Sprintf("Q%v", startDetailNonElectric+2+numberNonElectric), fmt.Sprintf("SUM(Q%v:Q%v)", startDetailNonElectric+2, startDetailNonElectric+1+numberNonElectric))

	var startCement int
	startCement = startDetailNonElectric + 5 + numberNonElectric

	file.SetCellValue(sheetName, fmt.Sprintf("A%v", startCement), "Data Realisasi Penjualan Batubara Untuk Memenuhi Kebutuhan Batubara Untuk Industri Semen")

	file.MergeCell(sheetName, fmt.Sprintf("A%v", startCement+1), fmt.Sprintf("A%v", startCement+2))
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startCement+1), fmt.Sprintf("D%v", startCement+2))

	file.MergeCell(sheetName, fmt.Sprintf("E%v", startCement+1), fmt.Sprintf("Q%v", startCement+1))

	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startCement+1), "Realisasi Penjualan Batubara (Pengguna Akhir Batubara) (Ton)")

	file.SetCellValue(sheetName, fmt.Sprintf("A%v", startCement+2), "No.")
	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startCement+2), "END USER")
	file.SetCellValue(sheetName, fmt.Sprintf("E%v", startCement+2), "January")
	file.SetCellValue(sheetName, fmt.Sprintf("F%v", startCement+2), "February")
	file.SetCellValue(sheetName, fmt.Sprintf("G%v", startCement+2), "March")
	file.SetCellValue(sheetName, fmt.Sprintf("H%v", startCement+2), "April")
	file.SetCellValue(sheetName, fmt.Sprintf("I%v", startCement+2), "May")
	file.SetCellValue(sheetName, fmt.Sprintf("J%v", startCement+2), "June")
	file.SetCellValue(sheetName, fmt.Sprintf("K%v", startCement+2), "July")
	file.SetCellValue(sheetName, fmt.Sprintf("L%v", startCement+2), "August")
	file.SetCellValue(sheetName, fmt.Sprintf("M%v", startCement+2), "September")
	file.SetCellValue(sheetName, fmt.Sprintf("N%v", startCement+2), "October")
	file.SetCellValue(sheetName, fmt.Sprintf("O%v", startCement+2), "November")
	file.SetCellValue(sheetName, fmt.Sprintf("P%v", startCement+2), "December")
	file.SetCellValue(sheetName, fmt.Sprintf("Q%v", startCement+2), "TOTAL")

	var numberCement int = 1

	var countCement int
	if countIndustry > 16 {
		countCement = countIndustry - 16
	}

	var baseCement int

	baseCement = 178 + 16 + countCement
	file.SetCellValue(chartSheetName, fmt.Sprintf("B%v", baseCement), "DETAIL REKAP TOTAL DMO BERDASARKAN JENIS INDUSTRI - Semen")

	for k, value := range reportSaleDetail.CompanyCement {

		file.SetCellValue(sheetName, fmt.Sprintf("B%v", startCement+(3+numberCement-1)), k)

		for _, v := range value {
			file.SetCellValue(sheetName, fmt.Sprintf("D%v", startCement+(3+numberCement-1)), v)
			file.SetCellValue(sheetName, fmt.Sprintf("A%v", startCement+(3+numberCement-1)), numberCement)

			if _, ok := reportSaleDetail.Cement.January[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("E%v", startCement+(3+numberCement-1)), reportSaleDetail.Cement.January[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("E%v", startCement+(3+numberCement-1)), 0)
			}

			if _, ok := reportSaleDetail.Cement.February[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("F%v", startCement+(3+numberCement-1)), reportSaleDetail.Cement.February[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("F%v", startCement+(3+numberCement-1)), 0)
			}

			if _, ok := reportSaleDetail.Cement.March[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("G%v", startCement+(3+numberCement-1)), reportSaleDetail.Cement.March[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("G%v", startCement+(3+numberCement-1)), 0)
			}

			if _, ok := reportSaleDetail.Cement.April[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("H%v", startCement+(3+numberCement-1)), reportSaleDetail.Cement.April[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("H%v", startCement+(3+numberCement-1)), 0)
			}

			if _, ok := reportSaleDetail.Cement.May[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("I%v", startCement+(3+numberCement-1)), reportSaleDetail.Cement.May[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("I%v", startCement+(3+numberCement-1)), 0)
			}

			if _, ok := reportSaleDetail.Cement.June[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("J%v", startCement+(3+numberCement-1)), reportSaleDetail.Cement.June[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("J%v", startCement+(3+numberCement-1)), 0)
			}

			if _, ok := reportSaleDetail.Cement.July[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("K%v", startCement+(3+numberCement-1)), reportSaleDetail.Cement.July[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("K%v", startCement+(3+numberCement-1)), 0)
			}

			if _, ok := reportSaleDetail.Cement.August[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("L%v", startCement+(3+numberCement-1)), reportSaleDetail.Cement.August[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("L%v", startCement+(3+numberCement-1)), 0)
			}

			if _, ok := reportSaleDetail.Cement.September[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("M%v", startCement+(3+numberCement-1)), reportSaleDetail.Cement.September[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("M%v", startCement+(3+numberCement-1)), 0)
			}

			if _, ok := reportSaleDetail.Cement.October[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("N%v", startCement+(3+numberCement-1)), reportSaleDetail.Cement.October[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("N%v", startCement+(3+numberCement-1)), 0)
			}

			if _, ok := reportSaleDetail.Cement.November[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("O%v", startCement+(3+numberCement-1)), reportSaleDetail.Cement.November[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("O%v", startCement+(3+numberCement-1)), 0)
			}

			if _, ok := reportSaleDetail.Cement.December[k]; ok {
				file.SetCellValue(sheetName, fmt.Sprintf("P%v", startCement+(3+numberCement-1)), reportSaleDetail.Cement.December[k][v])
			} else {
				file.SetCellValue(sheetName, fmt.Sprintf("P%v", startCement+(3+numberCement-1)), 0)
			}

			file.SetCellFormula(sheetName, fmt.Sprintf("Q%v", startCement+(3+numberCement-1)), fmt.Sprintf("SUM(E%v:P%v)", startCement+(3+numberCement-1), startCement+(3+numberCement-1)))

			file.SetCellFormula(chartSheetName, fmt.Sprintf("B%v", baseCement+1+countCement), fmt.Sprintf("Detail!D%v", startCement+(3+numberCement-1)))
			file.SetCellFormula(chartSheetName, fmt.Sprintf("C%v", baseCement+1+countCement), fmt.Sprintf("Detail!Q%v", startCement+(3+numberCement-1)))
			numberCement += 1
			countCement += 1
		}

		file.MergeCell(sheetName, fmt.Sprintf("B%v", startCement+(2+numberCement-1)), fmt.Sprintf("C%v", startCement+(2+numberCement-1)))
	}

	errChartSaleCement := file.SetCellStyle(chartSheetName, fmt.Sprintf("B%v", baseCement+1), fmt.Sprintf("C%v", baseCement+countCement), formatNumberStyle)

	if errChartSaleCement != nil {
		return file, errChartSaleCement
	}

	var seriesChartCompanyCement string

	seriesChartCompanyCement += fmt.Sprintf(`{
						"name": "Sales Semen",
						"categories": "%v!$B$%v:$B$%v",
						"values": "%v!$C$%v:$C$%v",
						"marker": {
									"symbol": "square"
								} 
				}`, chartSheetName, baseCement+1, baseCement+1+countCement, chartSheetName, baseCement+1, baseCement+1+countCement)

	valueChartCement := fmt.Sprintf(`{
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
	        "name": "Penjualan berdasarkan JENIS INDUSTRI (Semen)"
	    },
	    "show_blanks_as": "zero"
	}`, seriesChartCompanyCement)

	if err := file.AddChart(chartSheetName, fmt.Sprintf("G%v", baseCement), valueChartCement); err != nil {
		return file, err
	}

	errNoCementStyle := file.SetCellStyle(sheetName, fmt.Sprintf("A%v", startCement+3), fmt.Sprintf("A%v", startCement+3+numberCement-1), centerStyle)

	if errNoCementStyle != nil {
		return file, errNoCementStyle
	}

	errTitleTableCementStyle := file.SetCellStyle(sheetName, fmt.Sprintf("A%v", startCement+1), fmt.Sprintf("Q%v", startCement+2), boldTitleTableStyle)

	if errTitleTableCementStyle != nil {
		return file, errTitleTableCementStyle
	}

	errNumberCementStyle := file.SetCellStyle(sheetName, fmt.Sprintf("E%v", startCement+3), fmt.Sprintf("P%v", startCement+3+numberCement-1), formatNumberStyle)

	if errNumberCementStyle != nil {
		return file, errNumberCementStyle
	}

	errNumberTotalCementStyle := file.SetCellStyle(sheetName, fmt.Sprintf("Q%v", startCement+3), fmt.Sprintf("Q%v", startCement+3+numberCement-1), boldNumberStyle)

	if errNumberTotalCementStyle != nil {
		return file, errNumberTotalCementStyle
	}

	errCompanyCementStyle := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startCement+3), fmt.Sprintf("D%v", startCement+3+numberCement-1), boldTitleTableStyle)

	if errCompanyCementStyle != nil {
		return file, errCompanyCementStyle
	}

	file.SetCellValue(sheetName, fmt.Sprintf("B%v", startCement+2+numberCement), "TOTAL")
	file.MergeCell(sheetName, fmt.Sprintf("B%v", startCement+2+numberCement), fmt.Sprintf("D%v", startCement+2+numberCement))

	errSaleDetailCementTotalStyle := file.SetCellStyle(sheetName, fmt.Sprintf("B%v", startCement+2+numberCement), fmt.Sprintf("B%v", startCement+2+numberCement), boldNumberStyle)

	if errSaleDetailCementTotalStyle != nil {
		return file, errSaleDetailCementTotalStyle
	}

	errSaleDetailCementTotalNumberStyle := file.SetCellStyle(sheetName, fmt.Sprintf("E%v", startCement+2+numberCement), fmt.Sprintf("Q%v", startCement+2+numberCement), boldNumberStyle)

	if errSaleDetailCementTotalNumberStyle != nil {
		return file, errSaleDetailCementTotalNumberStyle
	}

	file.SetCellFormula(sheetName, fmt.Sprintf("E%v", startCement+2+numberCement), fmt.Sprintf("SUM(E%v:E%v)", startCement+3, startCement+numberCement+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("F%v", startCement+2+numberCement), fmt.Sprintf("SUM(F%v:F%v)", startCement+3, startCement+numberCement+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("G%v", startCement+2+numberCement), fmt.Sprintf("SUM(G%v:G%v)", startCement+3, startCement+numberCement+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("H%v", startCement+2+numberCement), fmt.Sprintf("SUM(H%v:H%v)", startCement+3, startCement+numberCement+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("I%v", startCement+2+numberCement), fmt.Sprintf("SUM(I%v:I%v)", startCement+3, startCement+numberCement+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("J%v", startCement+2+numberCement), fmt.Sprintf("SUM(J%v:J%v)", startCement+3, startCement+numberCement+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("K%v", startCement+2+numberCement), fmt.Sprintf("SUM(K%v:K%v)", startCement+3, startCement+numberCement+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("L%v", startCement+2+numberCement), fmt.Sprintf("SUM(L%v:L%v)", startCement+3, startCement+numberCement+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("M%v", startCement+2+numberCement), fmt.Sprintf("SUM(M%v:M%v)", startCement+3, startCement+numberCement+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("N%v", startCement+2+numberCement), fmt.Sprintf("SUM(N%v:N%v)", startCement+3, startCement+numberCement+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("O%v", startCement+2+numberCement), fmt.Sprintf("SUM(O%v:O%v)", startCement+3, startCement+numberCement+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("P%v", startCement+2+numberCement), fmt.Sprintf("SUM(P%v:P%v)", startCement+3, startCement+numberCement+1))

	file.SetCellFormula(sheetName, fmt.Sprintf("Q%v", startCement+2+numberCement), fmt.Sprintf("SUM(Q%v:Q%v)", startCement+3, startCement+numberCement+1))

	var startSaleExportImport int
	startSaleExportImport = startCement + 6 + numberCement

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

func (s *service) GetTransactionReport(iupopkId int, input TransactionReportInput, typeTransaction string) ([]TransactionReport, error) {
	transactionReport, err := s.repository.GetTransactionReport(iupopkId, input, typeTransaction)

	return transactionReport, err
}

func (s *service) CreateTransactionReport(file *excelize.File, sheetName string, iupopk iupopk.Iupopk, transactionData []TransactionReport) (*excelize.File, error) {
	file.SetCellValue(sheetName, "A1", iupopk.Name)

	custFmtQuantity := "#,##0.000"
	custFmtQuality := "#,##0.00"

	custFmtDate := "dd-mm-yyyy"

	for idx, value := range transactionData {
		file.SetCellValue(sheetName, "A1", iupopk.Name)

		file.SetCellValue(sheetName, fmt.Sprintf("A%v", 4+idx), value.TransactionType)
		file.SetCellValue(sheetName, fmt.Sprintf("B%v", 4+idx), idx+1)

		shippingDate := strings.Split(*value.ShippingDate, "T")
		file.SetCellValue(sheetName, fmt.Sprintf("C%v", 4+idx), shippingDate[0])
		file.SetCellValue(sheetName, fmt.Sprintf("D%v", 4+idx), value.Quantity)
		if value.Barge != nil && value.Tugboat != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("E%v", 4+idx), fmt.Sprintf("%s / %s", strings.ToUpper(value.Tugboat.Name), strings.ToUpper(value.Barge.Name)))
		}

		if value.Barge == nil && value.Tugboat != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("E%v", 4+idx), fmt.Sprintf("%s / -", strings.ToUpper(value.Tugboat.Name)))
		}

		if value.Barge != nil && value.Tugboat == nil {
			file.SetCellValue(sheetName, fmt.Sprintf("E%v", 4+idx), fmt.Sprintf("- / %s", strings.ToUpper(value.Barge.Name)))
		}

		if value.Barge == nil && value.Tugboat == nil {
			file.SetCellValue(sheetName, fmt.Sprintf("E%v", 4+idx), "- / -")
		}

		if value.SalesSystem != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("F%v", 4+idx), strings.ToUpper(value.SalesSystem.Name))
		}

		if value.Destination != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("G%v", 4+idx), strings.ToUpper(value.Destination.Name))
		}

		if value.Vessel != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("H%v", 4+idx), strings.ToUpper(value.Vessel.Name))
		}

		file.SetCellValue(sheetName, fmt.Sprintf("I%v", 4+idx), strings.ToUpper(iupopk.Name))

		if value.Customer != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("J%v", 4+idx), strings.ToUpper(value.Customer.CompanyName))
		}

		if value.LoadingPort != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("K%v", 4+idx), strings.ToUpper(value.LoadingPort.Name))

			if value.LoadingPort.PortLocation.Name != "" {
				file.SetCellValue(sheetName, fmt.Sprintf("L%v", 4+idx), strings.ToUpper(value.LoadingPort.PortLocation.Name))
			}
		}

		if value.UnloadingPort != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("M%v", 4+idx), strings.ToUpper(value.UnloadingPort.Name))

			if value.UnloadingPort.PortLocation.Name != "" {
				file.SetCellValue(sheetName, fmt.Sprintf("N%v", 4+idx), strings.ToUpper(value.UnloadingPort.PortLocation.Name))
			}
		}

		if value.TransactionType == "DN" {
			if value.DmoDestinationPort != nil {
				if value.SalesSystem != nil {
					if strings.Contains(value.SalesSystem.Name, "Barge") && value.Vessel != nil {
						file.SetCellValue(sheetName, fmt.Sprintf("O%v", 4+idx), fmt.Sprintf("%s @ %s", strings.ToUpper(value.Vessel.Name), strings.ToUpper(value.DmoDestinationPort.Name)))
					} else {
						file.SetCellValue(sheetName, fmt.Sprintf("O%v", 4+idx), strings.ToUpper(value.DmoDestinationPort.Name))
					}
				} else {
					file.SetCellValue(sheetName, fmt.Sprintf("O%v", 4+idx), strings.ToUpper(value.DmoDestinationPort.Name))
				}
			}
		}

		if value.TransactionType == "LN" {
			file.SetCellValue(sheetName, fmt.Sprintf("O%v", 4+idx), strings.ToUpper(value.DmoDestinationPortLnName))
		}

		if value.SkbDate != nil {
			skbDate := strings.Split(*value.SkbDate, "T")
			file.SetCellValue(sheetName, fmt.Sprintf("P%v", 4+idx), skbDate[0])
		}

		file.SetCellValue(sheetName, fmt.Sprintf("Q%v", 4+idx), strings.ToUpper(value.SkbNumber))

		if value.SkabDate != nil {
			skabDate := strings.Split(*value.SkabDate, "T")
			file.SetCellValue(sheetName, fmt.Sprintf("R%v", 4+idx), skabDate[0])
		}

		file.SetCellValue(sheetName, fmt.Sprintf("S%v", 4+idx), strings.ToUpper(value.SkabNumber))

		if value.BillOfLadingDate != nil {
			blDate := strings.Split(*value.BillOfLadingDate, "T")
			file.SetCellValue(sheetName, fmt.Sprintf("T%v", 4+idx), blDate[0])
		}

		file.SetCellValue(sheetName, fmt.Sprintf("U%v", 4+idx), strings.ToUpper(value.BillOfLadingNumber))

		file.SetCellValue(sheetName, fmt.Sprintf("V%v", 4+idx), value.RoyaltyRate/100)

		file.SetCellValue(sheetName, fmt.Sprintf("W%v", 4+idx), value.DpRoyaltyPrice)

		if value.DpRoyaltyDate != nil {
			dpRoyaltyDate := strings.Split(*value.DpRoyaltyDate, "T")
			file.SetCellValue(sheetName, fmt.Sprintf("X%v", 4+idx), dpRoyaltyDate[0])
		}

		if value.DpRoyaltyBillingCode != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("Y%v", 4+idx), strings.ToUpper(*value.DpRoyaltyBillingCode))
		}

		if value.DpRoyaltyNtpn != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("Z%v", 4+idx), strings.ToUpper(*value.DpRoyaltyNtpn))
		}

		file.SetCellValue(sheetName, fmt.Sprintf("AA%v", 4+idx), value.DpRoyaltyTotal)

		if value.PaymentDpRoyaltyDate != nil {
			paymentDpRoyaltyDate := strings.Split(*value.PaymentDpRoyaltyDate, "T")
			file.SetCellValue(sheetName, fmt.Sprintf("AD%v", 4+idx), paymentDpRoyaltyDate[0])
		}

		if value.PaymentDpRoyaltyBillingCode != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("AE%v", 4+idx), strings.ToUpper(*value.PaymentDpRoyaltyBillingCode))
		}

		if value.PaymentDpRoyaltyNtpn != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("AF%v", 4+idx), strings.ToUpper(*value.PaymentDpRoyaltyNtpn))
		}

		file.SetCellValue(sheetName, fmt.Sprintf("AG%v", 4+idx), value.PaymentDpRoyaltyTotal)

		if value.LhvDate != nil {
			lhvDate := strings.Split(*value.LhvDate, "T")
			file.SetCellValue(sheetName, fmt.Sprintf("AJ%v", 4+idx), lhvDate[0])
		}

		file.SetCellValue(sheetName, fmt.Sprintf("AK%v", 4+idx), strings.ToUpper(value.LhvNumber))

		if value.Surveyor != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("AL%v", 4+idx), strings.ToUpper(value.Surveyor.Name))
		}

		if value.CowDate != nil {
			cowDate := strings.Split(*value.CowDate, "T")
			file.SetCellValue(sheetName, fmt.Sprintf("AM%v", 4+idx), cowDate[0])
		}

		file.SetCellValue(sheetName, fmt.Sprintf("AN%v", 4+idx), strings.ToUpper(value.CowNumber))

		if value.CoaDate != nil {
			coaDate := strings.Split(*value.CoaDate, "T")
			file.SetCellValue(sheetName, fmt.Sprintf("AO%v", 4+idx), coaDate[0])
		}

		file.SetCellValue(sheetName, fmt.Sprintf("AP%v", 4+idx), strings.ToUpper(value.CoaNumber))

		file.SetCellValue(sheetName, fmt.Sprintf("AQ%v", 4+idx), value.QualityTmAr)

		file.SetCellValue(sheetName, fmt.Sprintf("AR%v", 4+idx), value.QualityImAdb)

		file.SetCellValue(sheetName, fmt.Sprintf("AS%v", 4+idx), value.QualityAshAdb)

		file.SetCellValue(sheetName, fmt.Sprintf("AT%v", 4+idx), value.QualityAshAr)

		file.SetCellValue(sheetName, fmt.Sprintf("AU%v", 4+idx), value.QualityVmAdb)

		file.SetCellValue(sheetName, fmt.Sprintf("AV%v", 4+idx), value.QualityFcAdb)

		file.SetCellValue(sheetName, fmt.Sprintf("AW%v", 4+idx), value.QualityTsAdb)

		file.SetCellValue(sheetName, fmt.Sprintf("AX%v", 4+idx), value.QualityTsAr)

		file.SetCellValue(sheetName, fmt.Sprintf("AY%v", 4+idx), value.QualityCaloriesAdb)

		file.SetCellValue(sheetName, fmt.Sprintf("AZ%v", 4+idx), value.QualityCaloriesAr)

		file.SetCellValue(sheetName, fmt.Sprintf("BA%v", 4+idx), value.BargingDistance)

		if value.InvoiceDate != nil {
			invoiceDate := strings.Split(*value.InvoiceDate, "T")
			file.SetCellValue(sheetName, fmt.Sprintf("BB%v", 4+idx), invoiceDate[0])
		}

		file.SetCellValue(sheetName, fmt.Sprintf("BC%v", 4+idx), strings.ToUpper(value.InvoiceNumber))

		file.SetCellValue(sheetName, fmt.Sprintf("BD%v", 4+idx), value.InvoicePriceUnit)

		file.SetCellValue(sheetName, fmt.Sprintf("BE%v", 4+idx), value.InvoicePriceTotal)

		file.SetCellValue(sheetName, fmt.Sprintf("BG%v", 4+idx), strings.ToUpper(value.ContractNumber))

		if value.ContractDate != nil {
			contractDate := strings.Split(*value.ContractDate, "T")
			file.SetCellValue(sheetName, fmt.Sprintf("BH%v", 4+idx), contractDate[0])
		}

		if value.DmoBuyer != nil {
			file.SetCellValue(sheetName, fmt.Sprintf("BI%v", 4+idx), strings.ToUpper(value.DmoBuyer.CompanyName))

			if value.DmoBuyer.IndustryType != nil {
				file.SetCellValue(sheetName, fmt.Sprintf("BJ%v", 4+idx), strings.ToUpper(value.DmoBuyer.IndustryType.SystemCategory))
			}
		}
	}

	border := []excelize.Border{
		{Type: "left", Color: "000000", Style: 1},
		{Type: "top", Color: "000000", Style: 1},
		{Type: "bottom", Color: "000000", Style: 1},
		{Type: "right", Color: "000000", Style: 1},
	}

	titleCenterStyle, _ := file.NewStyle(&excelize.Style{
		Border: border,
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})

	dateStyle, _ := file.NewStyle(&excelize.Style{
		CustomNumFmt: &custFmtDate,
	})

	quantityStyle, _ := file.NewStyle(&excelize.Style{
		CustomNumFmt: &custFmtQuantity,
	})

	qualityStyle, _ := file.NewStyle(&excelize.Style{
		CustomNumFmt: &custFmtQuality,
	})

	percentStyle, _ := file.NewStyle(&excelize.Style{
		NumFmt: 9,
	})

	errDate1 := file.SetColStyle(sheetName, "C", dateStyle)
	if errDate1 != nil {
		return file, errDate1
	}
	errDate2 := file.SetColStyle(sheetName, "P", dateStyle)
	if errDate2 != nil {
		return file, errDate2
	}
	errDate3 := file.SetColStyle(sheetName, "R", dateStyle)
	if errDate3 != nil {
		return file, errDate3
	}
	errDate4 := file.SetColStyle(sheetName, "T", dateStyle)
	if errDate4 != nil {
		return file, errDate4
	}
	errDate5 := file.SetColStyle(sheetName, "X", dateStyle)
	if errDate5 != nil {
		return file, errDate5
	}
	errDate6 := file.SetColStyle(sheetName, "AD", dateStyle)
	if errDate6 != nil {
		return file, errDate6
	}
	errDate7 := file.SetColStyle(sheetName, "AJ", dateStyle)
	if errDate7 != nil {
		return file, errDate7
	}
	errDate8 := file.SetColStyle(sheetName, "AM", dateStyle)
	if errDate8 != nil {
		return file, errDate8
	}
	errDate9 := file.SetColStyle(sheetName, "AO", dateStyle)
	if errDate9 != nil {
		return file, errDate9
	}
	errDate10 := file.SetColStyle(sheetName, "BB", dateStyle)
	if errDate10 != nil {
		return file, errDate10
	}
	errDate11 := file.SetColStyle(sheetName, "BH", dateStyle)
	if errDate11 != nil {
		return file, errDate11
	}

	errQuantity1 := file.SetColStyle(sheetName, "D", quantityStyle)
	if errQuantity1 != nil {
		return file, errQuantity1
	}
	errQuantity2 := file.SetColStyle(sheetName, "W", quantityStyle)
	if errQuantity2 != nil {
		return file, errQuantity2
	}
	errQuantity3 := file.SetColStyle(sheetName, "AA", quantityStyle)
	if errQuantity3 != nil {
		return file, errQuantity3
	}
	errQuantity4 := file.SetColStyle(sheetName, "AC", quantityStyle)
	if errQuantity4 != nil {
		return file, errQuantity4
	}
	errQuantity5 := file.SetColStyle(sheetName, "AG", quantityStyle)
	if errQuantity5 != nil {
		return file, errQuantity5
	}
	errQuantity6 := file.SetColStyle(sheetName, "AI", quantityStyle)
	if errQuantity6 != nil {
		return file, errQuantity6
	}
	errQuantity7 := file.SetColStyle(sheetName, "AY:AZ", quantityStyle)
	if errQuantity7 != nil {
		return file, errQuantity7
	}
	errQuantity8 := file.SetColStyle(sheetName, "BD:BE", quantityStyle)
	if errQuantity8 != nil {
		return file, errQuantity8
	}

	errQuality1 := file.SetColStyle(sheetName, "AQ:AX", qualityStyle)
	if errQuality1 != nil {
		return file, errQuality1
	}

	errPercent1 := file.SetColStyle(sheetName, "V", percentStyle)
	if errPercent1 != nil {
		return file, errPercent1
	}

	errTitleCenter := file.SetRowStyle(sheetName, 2, 3, titleCenterStyle)
	if errTitleCenter != nil {
		return file, errTitleCenter
	}

	file.SetColWidth(sheetName, "A", "A", 10)

	return file, nil
}
