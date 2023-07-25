package royaltyreport

import (
	"ajebackend/helper"

	"bytes"
	"encoding/json"
	"net/http"
)

type Service interface {
	GetTransactionRoyaltyReport(dateFrom string, dateTo string, iupopkId int) ([]RoyaltyReportData, error)
	GetDetailTransactionRoyaltyReport(id int, iupopkId int) (RoyaltyReportDetail, error)
	RequestCreateExcelRoyaltyReport(inputRequestCreateExcel InputRequestCreateUploadRoyaltyReport) (map[string]interface{}, error)
	ListRoyaltyReport(page int, sortFilter SortFilterRoyaltyReport, iupopkId int) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetTransactionRoyaltyReport(dateFrom string, dateTo string, iupopkId int) ([]RoyaltyReportData, error) {
	listTransaction, listTransactionErr := s.repository.GetTransactionRoyaltyReport(dateFrom, dateTo, iupopkId)

	return listTransaction, listTransactionErr
}

func (s *service) GetDetailTransactionRoyaltyReport(id int, iupopkId int) (RoyaltyReportDetail, error) {
	detailRoyaltyReport, detailRoyaltyReportErr := s.repository.GetDetailTransactionRoyaltyReport(id, iupopkId)

	return detailRoyaltyReport, detailRoyaltyReportErr
}

func (s *service) RequestCreateExcelRoyaltyReport(inputRequestCreateExcel InputRequestCreateUploadRoyaltyReport) (map[string]interface{}, error) {
	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("BASE_JOB_URL")

	urlPost := baseURL + "/create/report/royalty/report"
	body, bodyErr := json.Marshal(inputRequestCreateExcel)

	if bodyErr != nil {
		return res, bodyErr
	}
	var payload = bytes.NewBufferString(string(body))

	req, doReqErr := http.NewRequest("POST", urlPost, payload)

	if req != nil {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
	}
	client := &http.Client{}
	resp, doReqErr := client.Do(req)

	if doReqErr != nil {
		return res, doReqErr
	}

	json.NewDecoder(resp.Body).Decode(&res)

	return res, doReqErr
}

func (s *service) ListRoyaltyReport(page int, sortFilter SortFilterRoyaltyReport, iupopkId int) (Pagination, error) {
	listRoyaltyReport, listRoyaltyReportErr := s.repository.ListRoyaltyReport(page, sortFilter, iupopkId)

	return listRoyaltyReport, listRoyaltyReportErr
}
