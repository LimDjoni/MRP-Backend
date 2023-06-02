package coareportln

import (
	"ajebackend/helper"
	"ajebackend/model/transaction"
	"bytes"
	"encoding/json"
	"net/http"
)

type Service interface {
	GetTransactionCoaReportLn(dateFrom string, dateTo string, iupopkId int) ([]transaction.Transaction, error)
	GetDetailTransactionCoaReportLn(id int, iupopkId int) (CoaReportLnDetail, error)
	RequestCreateExcelCoaReportLn(inputRequestCreateExcel InputRequestCreateUploadCoaReportLn) (map[string]interface{}, error)
	ListCoaReportLn(page int, sortFilter SortFilterCoaReportLn, iupopkId int) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetTransactionCoaReportLn(dateFrom string, dateTo string, iupopkId int) ([]transaction.Transaction, error) {
	listTransaction, listTransactionErr := s.repository.GetTransactionCoaReportLn(dateFrom, dateTo, iupopkId)

	return listTransaction, listTransactionErr
}

func (s *service) GetDetailTransactionCoaReportLn(id int, iupopkId int) (CoaReportLnDetail, error) {
	detailCoaReportLn, detailCoaReportLnErr := s.repository.GetDetailTransactionCoaReportLn(id, iupopkId)

	return detailCoaReportLn, detailCoaReportLnErr
}

func (s *service) RequestCreateExcelCoaReportLn(inputRequestCreateExcel InputRequestCreateUploadCoaReportLn) (map[string]interface{}, error) {
	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("BASE_JOB_URL")

	urlPost := baseURL + "/create/report/coa/ln"
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

func (s *service) ListCoaReportLn(page int, sortFilter SortFilterCoaReportLn, iupopkId int) (Pagination, error) {
	listCoaReportLn, listCoaReportLnErr := s.repository.ListCoaReportLn(page, sortFilter, iupopkId)

	return listCoaReportLn, listCoaReportLnErr
}
