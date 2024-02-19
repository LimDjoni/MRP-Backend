package haulingsynchronize

import (
	"encoding/base64"
	"time"
)

type Service interface {
	SynchronizeTransactionIsp(syncData SynchronizeInputTransactionIsp) (bool, error)
	SynchronizeTransactionJetty(syncData SynchronizeInputTransactionJetty) (bool, error)
	UpdateSyncMasterIsp(iupopkId uint, dateTime time.Time) (bool, error)
	UpdateSyncMasterJetty(iupopkId uint, dateTime time.Time) (bool, error)
	GetSyncMasterDataIsp(iupopkId uint) (MasterDataIsp, error)
	GetSyncMasterDataJetty(iupopkId uint) (MasterDataJetty, error)
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

func (s *service) UpdateSyncMasterIsp(iupopkId uint, dateTime time.Time) (bool, error) {
	isUpdated, isUpdatedErr := s.repository.UpdateSyncMasterIsp(iupopkId, dateTime)

	return isUpdated, isUpdatedErr
}

func (s *service) UpdateSyncMasterJetty(iupopkId uint, dateTime time.Time) (bool, error) {
	isUpdated, isUpdatedErr := s.repository.UpdateSyncMasterJetty(iupopkId, dateTime)

	return isUpdated, isUpdatedErr
}

func (s *service) GetSyncMasterDataIsp(iupopkId uint) (MasterDataIsp, error) {
	getSyncMaster, getSyncMasterErr := s.repository.GetSyncMasterDataIsp(iupopkId)

	return getSyncMaster, getSyncMasterErr
}

func (s *service) GetSyncMasterDataJetty(iupopkId uint) (MasterDataJetty, error) {
	getSyncMaster, getSyncMasterErr := s.repository.GetSyncMasterDataJetty(iupopkId)

	return getSyncMaster, getSyncMasterErr
}
