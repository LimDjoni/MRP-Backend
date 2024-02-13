package haulingsynchronize

import (
	"ajebackend/helper"

	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
)

type Service interface {
	SynchronizeTransactionIsp(syncData SynchronizeInputTransactionIsp) (bool, error)
	SynchronizeTransactionJetty(syncData SynchronizeInputTransactionJetty) (bool, error)
	GetSynchronizeMasterData(iupopkId int) (SynchronizeInputMaster, error)
	UpdateSynchronizeMaster(iupopkId int) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (s *service) SynchronizeTransactionIsp(syncData SynchronizeInputTransactionIsp) (bool, error) {
	isSync, isSyncErr := s.repository.SynchronizeTransactionIsp(syncData)

	return isSync, isSyncErr
}

func (s *service) SynchronizeTransactionJetty(syncData SynchronizeInputTransactionJetty) (bool, error) {
	isSync, isSyncErr := s.repository.SynchronizeTransactionJetty(syncData)

	return isSync, isSyncErr
}

func (s *service) UpdateSynchronizeMaster(iupopkId int) (bool, error) {
	upd, updErr := s.repository.UpdateSynchronizeMaster(iupopkId)

	return upd, updErr
}

func (s *service) GetSynchronizeMasterData(iupopkId int) (SynchronizeInputMaster, error) {
	syncData, syncDataErr := s.repository.GetSynchronizeMasterData(iupopkId)

	return syncData, syncDataErr
}

func (s *service) SynchronizeToIsp(syncData SynchronizeInputMaster) (map[string]interface{}, error) {

	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("ISP_URL")

	userBasic := helper.GetEnvWithKey("USERNAME_BASIC")
	passBasic := helper.GetEnvWithKey("PASSWORD_BASIC")

	urlPost := baseURL

	body, bodyErr := json.Marshal(syncData)

	if bodyErr != nil {
		return res, bodyErr
	}
	var payload = bytes.NewBufferString(string(body))

	req, doReqErr := http.NewRequest("POST", urlPost, payload)

	if req != nil {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", "Basic "+basicAuth(userBasic, passBasic))
	}

	client := &http.Client{}
	resp, doReqErr := client.Do(req)

	if doReqErr != nil {
		return res, doReqErr
	}

	json.NewDecoder(resp.Body).Decode(&res)

	return res, doReqErr
}
