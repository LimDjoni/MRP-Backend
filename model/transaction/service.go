package transaction

import (
	"ajebackend/helper"
	"ajebackend/model/dmo"
	"ajebackend/model/trader"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
)

type Service interface {
	ListDataDN(page int, sortFilter SortAndFilter) (Pagination, error)
	DetailTransactionDN(id int) (Transaction, error)
	CheckDataUnique(inputTrans DataTransactionInput) (bool,bool,bool,bool)
	ListDataDNWithoutMinerba() ([]Transaction, error)
	CheckDataDnAndMinerba(listData []int)(bool, error)
	CheckDataDnAndMinerbaUpdate(listData []int, idMinerba int)([]Transaction, error)
	GetDetailMinerba(id int) (DetailMinerba, error)
	RequestCreateExcel(reqInput InputRequestCreateExcelMinerba) (map[string]interface{}, error)
	ListDataDNWithoutDmo() (ChooseTransactionDmo, error)
	CheckDataDnAndDmo(listData []int) ([]Transaction, error)
	GetDetailDmo(id int) (DetailDmo, error)
	RequestCreateDmo(reqInput InputRequestCreateUploadDmo) (map[string]interface{}, error)
	RequestCreateCustomDmo(dataDmo dmo.Dmo, traderEndUser trader.Trader ,bast *multipart.FileHeader, reconciliationLetter *multipart.FileHeader, statementLetter *multipart.FileHeader, authorization string ) (map[string]interface{}, error)
	GetReport(year int) (ReportRecapOutput, ReportDetailOutput, error)
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

func (s *service) CheckDataUnique(inputTrans DataTransactionInput) (bool,bool,bool,bool) {
	isDpRoyaltyNtpnUnique, isDpRoyaltyBillingCodeUnique, isPaymentDpRoyaltyNtpnUnique, isPaymentDpRoyaltyBillingCodeUnique := s.repository.CheckDataUnique(inputTrans)

	return isDpRoyaltyNtpnUnique, isDpRoyaltyBillingCodeUnique, isPaymentDpRoyaltyNtpnUnique, isPaymentDpRoyaltyBillingCodeUnique
}

func (s *service) ListDataDNWithoutMinerba() ([]Transaction, error) {
	listDataDNWithoutMinerba, listDataDNWithoutMinerbaErr := s.repository.ListDataDNWithoutMinerba()

	return listDataDNWithoutMinerba, listDataDNWithoutMinerbaErr
}

func (s *service) CheckDataDnAndMinerba(listData []int)(bool, error) {
	checkData, checkDataErr := s.repository.CheckDataDnAndMinerba(listData)

	return checkData, checkDataErr
}

func (s *service) CheckDataDnAndMinerbaUpdate(listData []int, idMinerba int)([]Transaction, error) {
	checkData, checkDataErr := s.repository.CheckDataDnAndMinerbaUpdate(listData, idMinerba)

	return checkData, checkDataErr
}

func (s *service) GetDetailMinerba(id int) (DetailMinerba, error) {
	detailMinerba, detailMinerbaErr := s.repository.GetDetailMinerba(id)

	return detailMinerba, detailMinerbaErr
}

func (s *service) RequestCreateExcel(reqInput InputRequestCreateExcelMinerba) (map[string]interface{}, error) {
	var res map[string]interface{}
	baseURL := helper.GetEnvWithKey("BASE_JOB_URL")

	urlPost := baseURL + "/create/minerba"
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

func (s *service) ListDataDNWithoutDmo() (ChooseTransactionDmo, error) {
	listDataDNWithoutDmo, listDataDNWithoutDmoErr := s.repository.ListDataDNWithoutDmo()

	return listDataDNWithoutDmo, listDataDNWithoutDmoErr
}

func (s *service) CheckDataDnAndDmo(listData []int) ([]Transaction, error) {
	checkData, checkDataErr := s.repository.CheckDataDnAndDmo(listData)

	return checkData, checkDataErr
}

func (s *service) GetDetailDmo(id int) (DetailDmo, error) {
	detailDmo, detailDmoErr := s.repository.GetDetailDmo(id)

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

func (s *service) RequestCreateCustomDmo(dataDmo dmo.Dmo, traderEndUser trader.Trader, bast *multipart.FileHeader, reconciliationLetter *multipart.FileHeader, statementLetter *multipart.FileHeader, authorization string ) (map[string]interface{}, error) {
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

	bastContent, _ := bast.Open()
	partBast, _ := w.CreateFormFile("bast", bast.Filename)
	bastByteContainer, _ := ioutil.ReadAll(bastContent)
	partBast.Write(bastByteContainer)

	reconciliationLetterContent, _ := reconciliationLetter.Open()
	partReconciliationLetter, _ := w.CreateFormFile("reconciliation_letter", reconciliationLetter.Filename)
	reconciliationLetterByteContainer, _ := ioutil.ReadAll(reconciliationLetterContent)
	partReconciliationLetter.Write(reconciliationLetterByteContainer)

	statementLetterContent, _ := bast.Open()
	partStatementLetter, _ := w.CreateFormFile("statement_letter", statementLetter.Filename)
	statementLetterByteContainer, _ := ioutil.ReadAll(statementLetterContent)
	partStatementLetter.Write(statementLetterByteContainer)

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

func (s *service) GetReport(year int) (ReportRecapOutput, ReportDetailOutput, error) {
	reportRecap, reportDetail, reportErr := s.repository.GetReport(year)

	return reportRecap, reportDetail, reportErr
}
