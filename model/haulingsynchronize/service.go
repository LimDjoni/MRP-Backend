package haulingsynchronize

type Service interface {
	SynchronizeTransactionIsp(syncData SynchronizeInputTransactionIsp) (bool, error)
	SynchronizeTransactionJetty(syncData SynchronizeInputTransactionJetty) (bool, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) SynchronizeTransactionIsp(syncData SynchronizeInputTransactionIsp) (bool, error) {
	isSync, isSyncErr := s.repository.SynchronizeTransactionIsp(syncData)

	return isSync, isSyncErr
}

func (s *service) SynchronizeTransactionJetty(syncData SynchronizeInputTransactionJetty) (bool, error) {
	isSync, isSyncErr := s.repository.SynchronizeTransactionJetty(syncData)

	return isSync, isSyncErr
}
