package transaction

import (
	"ajebackend/helper"
	"bytes"
	"encoding/json"
	"net/http"
)

type Service interface {
	ListDataDN(page int, sortFilter SortAndFilter) (Pagination, error)
	DetailTransactionDN(id int) (Transaction, error)
	ListDataDNWithoutMinerba() ([]Transaction, error)
	CheckDataDnAndMinerba(listData []int)(bool, error)
	GetDetailMinerba(id int) (DetailMinerba, error)
	RequestCreateExcel(reqInput InputRequestCreateExcelMinerba) (map[string]interface{}, error)
	ListDataDNWithoutDmo() ([]Transaction, error)
	CheckDataDnAndDmo(listData []int) (bool, error)
	GetDetailDmo(id int) (DetailDmo, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListDataDN(page int, sortFilter SortAndFilter) (Pagination, error) {
	listDN, listDNErr := s.repository.ListDataDN(page, sortFilter)

	return listDN, listDNErr
}

func (s *service) DetailTransactionDN(id int) (Transaction, error) {
	detailTransactionDN, detailTransactionDNErr := s.repository.DetailTransactionDN(id)

	return detailTransactionDN, detailTransactionDNErr
}

func (s *service) ListDataDNWithoutMinerba() ([]Transaction, error) {
	listDataDNWithoutMinerba, listDataDNWithoutMinerbaErr := s.repository.ListDataDNWithoutMinerba()

	return listDataDNWithoutMinerba, listDataDNWithoutMinerbaErr
}

func (s *service) CheckDataDnAndMinerba(listData []int)(bool, error) {
	checkData, checkDataErr := s.repository.CheckDataDnAndMinerba(listData)

	return checkData, checkDataErr
}

func (s *service) GetDetailMinerba(id int) (DetailMinerba, error) {
	detailMinerba, detailMinerbaErr := s.repository.GetDetailMinerba(id)

	return detailMinerba, detailMinerbaErr
}

func (s *service) RequestCreateExcel(reqInput InputRequestCreateExcelMinerba) (map[string]interface{}, error) {
	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("BASE_JOB_URL")

	urlPost := baseURL + "/createexcel"
	body, bodyErr := json.Marshal(reqInput)

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

func (s *service) ListDataDNWithoutDmo() ([]Transaction, error) {
	listDataDNWithoutDmo, listDataDNWithoutDmoErr := s.repository.ListDataDNWithoutDmo()

	return listDataDNWithoutDmo, listDataDNWithoutDmoErr
}

func (s *service) CheckDataDnAndDmo(listData []int) (bool, error) {
	checkData, checkDataErr := s.repository.CheckDataDnAndDmo(listData)

	return checkData, checkDataErr
}

func (s *service) GetDetailDmo(id int) (DetailDmo, error) {
	detailDmo, detailDmoErr := s.repository.GetDetailDmo(id)

	return detailDmo, detailDmoErr
}
