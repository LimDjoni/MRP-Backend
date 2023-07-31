package royaltyrecon

import (
	"ajebackend/helper"

	"bytes"
	"encoding/json"
	"net/http"
)

type Service interface {
	GetTransactionRoyaltyRecon(dateFrom string, dateTo string, iupopkId int) ([]RoyaltyReconData, error)
	GetDetailTransactionRoyaltyRecon(id int, iupopkId int) (RoyaltyReconDetail, error)
	RequestCreateExcelRoyaltyRecon(inputRequestCreateExcel InputRequestCreateUploadRoyaltyRecon) (map[string]interface{}, error)
	ListRoyaltyRecon(page int, sortFilter SortFilterRoyaltyRecon, iupopkId int) (Pagination, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetTransactionRoyaltyRecon(dateFrom string, dateTo string, iupopkId int) ([]RoyaltyReconData, error) {
	listTransaction, listTransactionErr := s.repository.GetTransactionRoyaltyRecon(dateFrom, dateTo, iupopkId)

	return listTransaction, listTransactionErr
}

func (s *service) GetDetailTransactionRoyaltyRecon(id int, iupopkId int) (RoyaltyReconDetail, error) {
	detailRoyaltyRecon, detailRoyaltyReconErr := s.repository.GetDetailTransactionRoyaltyRecon(id, iupopkId)

	return detailRoyaltyRecon, detailRoyaltyReconErr
}

func (s *service) RequestCreateExcelRoyaltyRecon(inputRequestCreateExcel InputRequestCreateUploadRoyaltyRecon) (map[string]interface{}, error) {
	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("BASE_JOB_URL")

	urlPost := baseURL + "/create/report/royalty/recon"
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

func (s *service) ListRoyaltyRecon(page int, sortFilter SortFilterRoyaltyRecon, iupopkId int) (Pagination, error) {
	listRoyaltyRecon, listRoyaltyReconErr := s.repository.ListRoyaltyRecon(page, sortFilter, iupopkId)

	return listRoyaltyRecon, listRoyaltyReconErr
}
