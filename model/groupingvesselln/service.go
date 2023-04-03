package groupingvesselln

import (
	"ajebackend/helper"
	"bytes"
	"encoding/json"
	"net/http"
)

type Service interface {
	ListGroupingVesselLn(page int, sortFilter SortFilterGroupingVesselLn, iupopkId int) (Pagination, error)
	ListGroupingVesselLnWithPeriod(month string, year int, iupopkId int) ([]GroupingVesselLn, error)
	DetailInsw(id int, iupopkId int) (DetailInsw, error)
	RequestCreateExcelInsw(reqInput InputRequestCreateUploadInsw) (map[string]interface{}, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListGroupingVesselLn(page int, sortFilter SortFilterGroupingVesselLn, iupopkId int) (Pagination, error) {
	listGroupingVesselLn, listGroupingVesselLnErr := s.repository.ListGroupingVesselLn(page, sortFilter, iupopkId)

	return listGroupingVesselLn, listGroupingVesselLnErr
}

func (s *service) ListGroupingVesselLnWithPeriod(month string, year int, iupopkId int) ([]GroupingVesselLn, error) {
	listGroupingVesselLnWithoutInsw, listGroupingVesselLnWithoutInswErr := s.repository.ListGroupingVesselLnWithPeriod(month, year, iupopkId)

	return listGroupingVesselLnWithoutInsw, listGroupingVesselLnWithoutInswErr
}

func (s *service) DetailInsw(id int, iupopkId int) (DetailInsw, error) {
	detailInsw, detailInswErr := s.repository.DetailInsw(id, iupopkId)

	return detailInsw, detailInswErr
}

func (s *service) RequestCreateExcelInsw(reqInput InputRequestCreateUploadInsw) (map[string]interface{}, error) {
	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("BASE_JOB_URL")

	urlPost := baseURL + "/create/insw"
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
