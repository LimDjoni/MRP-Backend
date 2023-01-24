package insw

import (
	"ajebackend/helper"
	"bytes"
	"encoding/json"
	"net/http"
)

type Service interface {
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RequestCreateExcel(reqInput InputUpdateDocumentInsw) (map[string]interface{}, error) {
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
