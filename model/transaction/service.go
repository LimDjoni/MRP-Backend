package transaction

import (
	"ajebackend/helper"
	"ajebackend/model/dmo"
	"ajebackend/model/dmovessel"
	"ajebackend/model/master/trader"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type Service interface {
	ListData(page int, sortFilter SortAndFilter, transactionType string, iupopkId int) (Pagination, error)
	DetailTransaction(id int, transactionType string, iupopkId int) (Transaction, error)
	CheckDataUnique(inputTrans DataTransactionInput) (bool, bool, bool, bool)
	ListDataDNWithoutMinerba(iupopkId int) ([]Transaction, error)
	CheckDataDnAndMinerba(listData []int, iupopkId int) (bool, error)
	CheckDataDnAndMinerbaUpdate(listData []int, idMinerba int, iupopkId int) ([]Transaction, error)
	GetDetailMinerba(id int, iupopkId int) (DetailMinerba, error)
	RequestCreateExcel(reqInput InputRequestCreateExcelMinerba) (map[string]interface{}, error)
	RequestCreateExcelLn(reqInput InputRequestCreateExcelMinerbaLn) (map[string]interface{}, error)
	ListDataDNBargeWithoutVessel(iupopkId int) ([]Transaction, error)
	ListDataDNBargeWithVessel(iupopkId int) ([]Transaction, error)
	ListDataDNVessel(iupopkId int) ([]Transaction, error)
	CheckDataDnAndDmo(listData []int, iupopkId int) ([]Transaction, error)
	CheckGroupingVesselAndDmo(listData []int, iupopkId int) ([]dmovessel.DmoVessel, error)
	GetDetailDmo(id int, iupopkId int) (DetailDmo, error)
	RequestCreateDmo(reqInput InputRequestCreateUploadDmo) (map[string]interface{}, error)
	RequestCreateCustomDmo(dataDmo dmo.Dmo, traderEndUser trader.Trader, reconciliationLetter *multipart.FileHeader, authorization string, reqInputCreateUploadDmo InputRequestCreateUploadDmo) (map[string]interface{}, error)
	GetReport(year int, iupopkId int) (ReportRecapOutput, ReportDetailOutput, error)
	GetListForReport(iupopkId int) (ListForCreatingReportDmoOutput, error)
	GetDetailGroupingVesselDn(id int, iupopkId int) (DetailGroupingVesselDn, error)
	ListDataDnWithoutGroup(iupopkId int) (ListTransactionNotHaveGroupingVessel, error)
	GetDetailGroupingVesselLn(id int, iupopkId int) (DetailGroupingVesselLn, error)
	ListDataLnWithoutGroup(iupopkId int) ([]Transaction, error)
	GetDetailMinerbaLn(id int, iupopkId int) (DetailMinerbaLn, error)
	ListDataLNWithoutMinerba(iupopkId int) ([]Transaction, error)
	CheckDataLnAndMinerbaLnUpdate(listData []int, idMinerba int, iupopkId int) ([]Transaction, error)
	CheckDataLnAndMinerbaLn(listData []int, iupopkId int) (bool, error)
	GetDataDmo(id uint, iupopkId int) (ListTransactionDmoBackgroundJob, error)
	RequestCreateReportDmo(input InputRequestCreateReportDmo) (map[string]interface{}, error)
	GetDetailReportDmo(id int, iupopkId int) (DetailReportDmo, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) ListData(page int, sortFilter SortAndFilter, transactionType string, iupopkId int) (Pagination, error) {
	listDN, listDNErr := s.repository.ListData(page, sortFilter, transactionType, iupopkId)

	return listDN, listDNErr
}

func (s *service) DetailTransaction(id int, transactionType string, iupopkId int) (Transaction, error) {
	detailTransactionDN, detailTransactionDNErr := s.repository.DetailTransaction(id, transactionType, iupopkId)

	return detailTransactionDN, detailTransactionDNErr
}

func (s *service) CheckDataUnique(inputTrans DataTransactionInput) (bool, bool, bool, bool) {
	isDpRoyaltyNtpnUnique, isDpRoyaltyBillingCodeUnique, isPaymentDpRoyaltyNtpnUnique, isPaymentDpRoyaltyBillingCodeUnique := s.repository.CheckDataUnique(inputTrans)

	return isDpRoyaltyNtpnUnique, isDpRoyaltyBillingCodeUnique, isPaymentDpRoyaltyNtpnUnique, isPaymentDpRoyaltyBillingCodeUnique
}

func (s *service) ListDataDNWithoutMinerba(iupopkId int) ([]Transaction, error) {
	listDataDNWithoutMinerba, listDataDNWithoutMinerbaErr := s.repository.ListDataDNWithoutMinerba(iupopkId)

	return listDataDNWithoutMinerba, listDataDNWithoutMinerbaErr
}

func (s *service) CheckDataDnAndMinerba(listData []int, iupopkId int) (bool, error) {
	checkData, checkDataErr := s.repository.CheckDataDnAndMinerba(listData, iupopkId)

	return checkData, checkDataErr
}

func (s *service) CheckDataDnAndMinerbaUpdate(listData []int, idMinerba int, iupopkId int) ([]Transaction, error) {
	checkData, checkDataErr := s.repository.CheckDataDnAndMinerbaUpdate(listData, idMinerba, iupopkId)

	return checkData, checkDataErr
}

func (s *service) GetDetailMinerba(id int, iupopkId int) (DetailMinerba, error) {
	detailMinerba, detailMinerbaErr := s.repository.GetDetailMinerba(id, iupopkId)

	return detailMinerba, detailMinerbaErr
}

func (s *service) RequestCreateExcel(reqInput InputRequestCreateExcelMinerba) (map[string]interface{}, error) {
	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("BASE_JOB_URL")

	urlPost := baseURL + "/create/minerba/dn"
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

func (s *service) RequestCreateExcelLn(reqInput InputRequestCreateExcelMinerbaLn) (map[string]interface{}, error) {
	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("BASE_JOB_URL")

	urlPost := baseURL + "/create/minerba/ln"
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

func (s *service) ListDataDNBargeWithoutVessel(iupopkId int) ([]Transaction, error) {
	listDataDNBargeWithoutVessel, listDataDNBargeWithoutVesselErr := s.repository.ListDataDNBargeWithoutVessel(iupopkId)

	return listDataDNBargeWithoutVessel, listDataDNBargeWithoutVesselErr
}

func (s *service) ListDataDNBargeWithVessel(iupopkId int) ([]Transaction, error) {
	listDataDNBargeWithVessel, listDataDNBargeWithVesselErr := s.repository.ListDataDNBargeWithVessel(iupopkId)

	return listDataDNBargeWithVessel, listDataDNBargeWithVesselErr
}

func (s *service) ListDataDNVessel(iupopkId int) ([]Transaction, error) {
	listDataDNVessel, listDataDNVesselErr := s.repository.ListDataDNVessel(iupopkId)

	return listDataDNVessel, listDataDNVesselErr
}

func (s *service) CheckDataDnAndDmo(listData []int, iupopkId int) ([]Transaction, error) {
	checkData, checkDataErr := s.repository.CheckDataDnAndDmo(listData, iupopkId)

	return checkData, checkDataErr
}

func (s *service) CheckGroupingVesselAndDmo(listData []int, iupopkId int) ([]dmovessel.DmoVessel, error) {
	checkGrouping, checkGroupingErr := s.repository.CheckGroupingVesselAndDmo(listData, iupopkId)

	return checkGrouping, checkGroupingErr
}

func (s *service) GetDetailDmo(id int, iupopkId int) (DetailDmo, error) {
	detailDmo, detailDmoErr := s.repository.GetDetailDmo(id, iupopkId)

	return detailDmo, detailDmoErr
}

func (s *service) RequestCreateDmo(reqInput InputRequestCreateUploadDmo) (map[string]interface{}, error) {
	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("BASE_JOB_URL")

	urlPost := baseURL + "/create/dmo"
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

func (s *service) RequestCreateCustomDmo(dataDmo dmo.Dmo, traderEndUser trader.Trader, reconciliationLetter *multipart.FileHeader, authorization string, reqInputCreateUploadDmo InputRequestCreateUploadDmo) (map[string]interface{}, error) {
	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("BASE_JOB_URL")
	var (
		buf = new(bytes.Buffer)
		w   = multipart.NewWriter(buf)
	)

	dataDmoMarshal, _ := json.Marshal(dataDmo)
	partDmo, _ := w.CreateFormField("data_dmo")
	partDmo.Write(dataDmoMarshal)

	traderEndUserMarshal, _ := json.Marshal(traderEndUser)
	partTraderEndUser, _ := w.CreateFormField("trader_end_user")
	partTraderEndUser.Write(traderEndUserMarshal)

	authorizationMarshal, _ := json.Marshal(authorization)
	partAuthorization, _ := w.CreateFormField("authorization")
	partAuthorization.Write(authorizationMarshal)

	reconciliationLetterContent, _ := reconciliationLetter.Open()
	partReconciliationLetter, _ := w.CreateFormFile("reconciliation_letter", reconciliationLetter.Filename)
	reconciliationLetterByteContainer, _ := ioutil.ReadAll(reconciliationLetterContent)
	partReconciliationLetter.Write(reconciliationLetterByteContainer)

	reqInputCreateUploadDmoMarshal, _ := json.Marshal(reqInputCreateUploadDmo)
	partDataTransaction, _ := w.CreateFormField("data_transaction")
	partDataTransaction.Write(reqInputCreateUploadDmoMarshal)

	w.Close()
	urlPost := baseURL + "/upload/dmo/custom"

	req, doReqErr := http.NewRequest("POST", urlPost, buf)

	if req != nil {
		req.Header.Add("Content-Type", w.FormDataContentType())
		//req.Header.Add("Accept", "multipart/form-data")
	}
	client := &http.Client{}
	resp, doReqErr := client.Do(req)

	if doReqErr != nil {
		return res, doReqErr
	}

	json.NewDecoder(resp.Body).Decode(&res)

	return res, doReqErr
}

func (s *service) GetReport(year int, iupopkId int) (ReportRecapOutput, ReportDetailOutput, error) {
	reportRecap, reportDetail, reportErr := s.repository.GetReport(year, iupopkId)

	return reportRecap, reportDetail, reportErr
}

func (s *service) GetListForReport(iupopkId int) (ListForCreatingReportDmoOutput, error) {
	listForReport, listForReportErr := s.repository.GetListForReport(iupopkId)

	return listForReport, listForReportErr
}

func (s *service) GetDetailGroupingVesselDn(id int, iupopkId int) (DetailGroupingVesselDn, error) {
	detailGroupingVesselDn, detailGroupingVesselDnErr := s.repository.GetDetailGroupingVesselDn(id, iupopkId)

	return detailGroupingVesselDn, detailGroupingVesselDnErr
}

func (s *service) ListDataDnWithoutGroup(iupopkId int) (ListTransactionNotHaveGroupingVessel, error) {
	listWithoutGroup, listWithoutGroupErr := s.repository.ListDataDnWithoutGroup(iupopkId)

	return listWithoutGroup, listWithoutGroupErr
}

func (s *service) GetDetailGroupingVesselLn(id int, iupopkId int) (DetailGroupingVesselLn, error) {
	detailGroupingVesselLn, detailGroupingVesselLnErr := s.repository.GetDetailGroupingVesselLn(id, iupopkId)

	return detailGroupingVesselLn, detailGroupingVesselLnErr
}

func (s *service) ListDataLnWithoutGroup(iupopkId int) ([]Transaction, error) {
	listWithoutGroup, listWithoutGroupErr := s.repository.ListDataLnWithoutGroup(iupopkId)

	return listWithoutGroup, listWithoutGroupErr
}

func (s *service) GetDetailMinerbaLn(id int, iupopkId int) (DetailMinerbaLn, error) {
	detailMinerbaLn, detailMinerbaLnErr := s.repository.GetDetailMinerbaLn(id, iupopkId)

	return detailMinerbaLn, detailMinerbaLnErr
}

func (s *service) ListDataLNWithoutMinerba(iupopkId int) ([]Transaction, error) {
	listDataLNWithoutMinerba, listDataLNWithoutMinerbaErr := s.repository.ListDataLNWithoutMinerba(iupopkId)

	return listDataLNWithoutMinerba, listDataLNWithoutMinerbaErr
}

func (s *service) CheckDataLnAndMinerbaLnUpdate(listData []int, idMinerba int, iupopkId int) ([]Transaction, error) {
	checkData, checkDataErr := s.repository.CheckDataLnAndMinerbaLnUpdate(listData, idMinerba, iupopkId)

	return checkData, checkDataErr
}

func (s *service) CheckDataLnAndMinerbaLn(listData []int, iupopkId int) (bool, error) {
	checkData, checkDataErr := s.repository.CheckDataLnAndMinerbaLn(listData, iupopkId)

	return checkData, checkDataErr
}

func (s *service) GetDataDmo(id uint, iupopkId int) (ListTransactionDmoBackgroundJob, error) {
	getDataReportDmo, getDataReportDmoErr := s.repository.GetDataDmo(id, iupopkId)

	return getDataReportDmo, getDataReportDmoErr
}

func (s *service) RequestCreateReportDmo(input InputRequestCreateReportDmo) (map[string]interface{}, error) {
	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("BASE_JOB_URL")

	urlPost := baseURL + "/create/dmo/report"
	body, bodyErr := json.Marshal(input)

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

func (s *service) GetDetailReportDmo(id int, iupopkId int) (DetailReportDmo, error) {
	detailReportDmo, detailReportDmoErr := s.repository.GetDetailReportDmo(id, iupopkId)

	return detailReportDmo, detailReportDmoErr
}
