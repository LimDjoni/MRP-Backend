package allmaster

import (
	"mrpbackend/model/master/brand"
	"mrpbackend/model/master/department"
	"mrpbackend/model/master/doh"
	"mrpbackend/model/master/form"
	"mrpbackend/model/master/heavyequipment"
	"mrpbackend/model/master/history"
	"mrpbackend/model/master/jabatan"
	"mrpbackend/model/master/kartukeluarga"
	"mrpbackend/model/master/ktp"
	"mrpbackend/model/master/mcu"
	"mrpbackend/model/master/pendidikan"
	"mrpbackend/model/master/position"
	"mrpbackend/model/master/role"
	"mrpbackend/model/master/series"
	"mrpbackend/model/master/sertifikat"
	"mrpbackend/model/master/userrole"
)

type Service interface {
	CreateUserRole(userRoleInput RegisterUserRoleInput) (userrole.UserRole, error)
	CreateBrand(brandInput RegisterBrandInput) (brand.Brand, error)
	CreateHeavyEquipment(heavyEquipmentInput RegisterHeavyEquipmentInput) (heavyequipment.HeavyEquipment, error)
	CreateSeries(seriesInput RegisterSeriesInput) (series.Series, error)
	CreateKartuKeluarga(kartukeluargaInput RegisterKartuKeluargaInput) (kartukeluarga.KartuKeluarga, error)
	CreateKTP(ktpInput RegisterKTPInput) (ktp.KTP, error)
	CreatePendidikan(pendidikanInput RegisterPendidikanInput) (pendidikan.Pendidikan, error)
	CreateDOH(dohInput RegisterDOHInput) (doh.DOH, error)
	CreateJabatan(dohInput RegisterJabatanInput) (jabatan.Jabatan, error)
	CreateSertifikat(sertifikatInput RegisterSertifikatInput) (sertifikat.Sertifikat, error)
	CreateMCU(mcuInput RegisterMCUInput) (mcu.MCU, error)
	CreateHistory(historyInput RegisterHistoryInput) (history.History, error)
	FindUserRole() ([]userrole.UserRole, error)
	FindUserRoleById(id uint) (userrole.UserRole, error)
	FindBrand() ([]brand.Brand, error)
	FindBrandById(id uint) (brand.Brand, error)
	FindHeavyEquipment() ([]heavyequipment.HeavyEquipment, error)
	FindHeavyEquipmentById(id uint) (heavyequipment.HeavyEquipment, error)
	FindHeavyEquipmentByBrandID(brandId uint) ([]heavyequipment.HeavyEquipment, error)
	FindSeries() ([]series.Series, error)
	FindSeriesById(id uint) (series.Series, error)
	FindSeriesByBrandAndEquipmentdID(brandId uint, heavyequipmentId uint) ([]series.Series, error)
	FindDepartment() ([]department.Department, error)
	FindRole() ([]role.Role, error)
	FindPosition() ([]position.Position, error)

	FindDOHById(id uint) (doh.DOH, error)
	UpdateDOH(inputDOH RegisterDOHInput, id int) (doh.DOH, error)
	DeleteDOH(id uint) (bool, error)
	FindDohKontrak(page int, sortFilter SortFilterDohKontrak) (Pagination, error)

	FindJabatanById(id uint) (jabatan.Jabatan, error)
	UpdateJabatan(inputJabatan RegisterJabatanInput, id int) (jabatan.Jabatan, error)
	DeleteJabatan(id uint) (bool, error)

	FindSertifikatById(id uint) (sertifikat.Sertifikat, error)
	UpdateSertifikat(inputSertifikat RegisterSertifikatInput, id int) (sertifikat.Sertifikat, error)
	DeleteSertifikat(id uint) (bool, error)

	FindMCUById(id uint) (mcu.MCU, error)
	UpdateMCU(inputMCU RegisterMCUInput, id int) (mcu.MCU, error)
	DeleteMCU(id uint) (bool, error)
	FindMCUBerkala(page int, sortFilter SortFilterDohKontrak) (Pagination, error)

	FindHistoryById(id uint) (history.History, error)
	UpdateHistory(inputHistory RegisterHistoryInput, id int) (history.History, error)
	DeleteHistory(id uint) (bool, error)

	GenerateSideBar(userID uint) ([]form.Form, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateUserRole(userRoleInput RegisterUserRoleInput) (userrole.UserRole, error) {
	newUserRole, err := s.repository.CreateUserRole(userRoleInput)

	return newUserRole, err
}

func (s *service) CreateBrand(brandInput RegisterBrandInput) (brand.Brand, error) {
	newBrand, err := s.repository.CreateBrand(brandInput)

	return newBrand, err
}

func (s *service) CreateSeries(seriesInput RegisterSeriesInput) (series.Series, error) {
	newSeries, err := s.repository.CreateSeries(seriesInput)

	return newSeries, err
}

func (s *service) CreateHeavyEquipment(heavyEquipmentInput RegisterHeavyEquipmentInput) (heavyequipment.HeavyEquipment, error) {
	newHeavyEquipment, err := s.repository.CreateHeavyEquipment(heavyEquipmentInput)

	return newHeavyEquipment, err
}

func (s *service) CreateKartuKeluarga(kartukeluargaInput RegisterKartuKeluargaInput) (kartukeluarga.KartuKeluarga, error) {
	newKartuKeluarga, err := s.repository.CreateKartuKeluarga(kartukeluargaInput)

	return newKartuKeluarga, err
}

func (s *service) CreateKTP(ktpInput RegisterKTPInput) (ktp.KTP, error) {
	newKTP, err := s.repository.CreateKTP(ktpInput)

	return newKTP, err
}

func (s *service) CreatePendidikan(pendidikanInput RegisterPendidikanInput) (pendidikan.Pendidikan, error) {
	newPendidikan, err := s.repository.CreatePendidikan(pendidikanInput)

	return newPendidikan, err
}

func (s *service) CreateDOH(dohInput RegisterDOHInput) (doh.DOH, error) {
	newDOH, err := s.repository.CreateDOH(dohInput)

	return newDOH, err
}

func (s *service) CreateJabatan(jabatanInput RegisterJabatanInput) (jabatan.Jabatan, error) {
	newJabatan, err := s.repository.CreateJabatan(jabatanInput)

	return newJabatan, err
}

func (s *service) CreateSertifikat(sertifikatInput RegisterSertifikatInput) (sertifikat.Sertifikat, error) {
	newSertifikat, err := s.repository.CreateSertifikat(sertifikatInput)

	return newSertifikat, err
}

func (s *service) CreateMCU(mcuInput RegisterMCUInput) (mcu.MCU, error) {
	newMCU, err := s.repository.CreateMCU(mcuInput)

	return newMCU, err
}

func (s *service) CreateHistory(historyInput RegisterHistoryInput) (history.History, error) {
	newHistory, err := s.repository.CreateHistory(historyInput)

	return newHistory, err
}

func (s *service) FindUserRole() ([]userrole.UserRole, error) {
	userRole, err := s.repository.FindUserRole()

	return userRole, err
}

func (s *service) FindUserRoleById(id uint) (userrole.UserRole, error) {
	userRole, err := s.repository.FindUserRoleById(id)

	return userRole, err
}

func (s *service) FindBrand() ([]brand.Brand, error) {
	brand, err := s.repository.FindBrand()

	return brand, err
}

func (s *service) FindBrandById(id uint) (brand.Brand, error) {
	brand, err := s.repository.FindBrandById(id)

	return brand, err
}

func (s *service) FindHeavyEquipment() ([]heavyequipment.HeavyEquipment, error) {
	heavyEquipment, err := s.repository.FindHeavyEquipment()

	return heavyEquipment, err
}

func (s *service) FindHeavyEquipmentById(id uint) (heavyequipment.HeavyEquipment, error) {
	heavyEquipment, err := s.repository.FindHeavyEquipmentById(id)

	return heavyEquipment, err
}

func (s *service) FindHeavyEquipmentByBrandID(brandId uint) ([]heavyequipment.HeavyEquipment, error) {
	heavyEquipment, err := s.repository.FindHeavyEquipmentByBrandID(brandId)

	return heavyEquipment, err
}

func (s *service) FindSeries() ([]series.Series, error) {
	series, err := s.repository.FindSeries()

	return series, err
}

func (s *service) FindSeriesById(id uint) (series.Series, error) {
	series, err := s.repository.FindSeriesById(id)

	return series, err
}

func (s *service) FindSeriesByBrandAndEquipmentdID(brandId uint, heavyequipmentId uint) ([]series.Series, error) {
	series, err := s.repository.FindSeriesByBrandAndEquipmentdID(brandId, heavyequipmentId)

	return series, err
}

func (s *service) FindDepartment() ([]department.Department, error) {
	department, err := s.repository.FindDepartment()

	return department, err
}

func (s *service) FindRole() ([]role.Role, error) {
	roles, err := s.repository.FindRole()

	return roles, err
}

func (s *service) FindPosition() ([]position.Position, error) {
	positions, err := s.repository.FindPosition()

	return positions, err
}

func (s *service) FindDOHById(id uint) (doh.DOH, error) {
	doh, err := s.repository.FindDOHById(id)

	return doh, err
}

func (s *service) UpdateDOH(inputDOH RegisterDOHInput, id int) (doh.DOH, error) {
	updateDOH, updateDOHErr := s.repository.UpdateDOH(inputDOH, id)

	return updateDOH, updateDOHErr
}

func (s *service) DeleteDOH(id uint) (bool, error) {
	isDeletedDOH, isDeletedDOHErr := s.repository.DeleteDOH(id)

	return isDeletedDOH, isDeletedDOHErr
}

func (s *service) FindDohKontrak(page int, sortFilter SortFilterDohKontrak) (Pagination, error) {
	doh, err := s.repository.FindDohKontrak(page, sortFilter)

	return doh, err
}

func (s *service) FindJabatanById(id uint) (jabatan.Jabatan, error) {
	jabatan, err := s.repository.FindJabatanById(id)

	return jabatan, err
}

func (s *service) UpdateJabatan(inputJabatan RegisterJabatanInput, id int) (jabatan.Jabatan, error) {
	updateJabatan, updateJabatanErr := s.repository.UpdateJabatan(inputJabatan, id)

	return updateJabatan, updateJabatanErr
}

func (s *service) DeleteJabatan(id uint) (bool, error) {
	isDeletedJabatan, isDeletedJabatanErr := s.repository.DeleteJabatan(id)

	return isDeletedJabatan, isDeletedJabatanErr
}

func (s *service) FindSertifikatById(id uint) (sertifikat.Sertifikat, error) {
	sertifikat, err := s.repository.FindSertifikatById(id)

	return sertifikat, err
}

func (s *service) UpdateSertifikat(inputSertifikat RegisterSertifikatInput, id int) (sertifikat.Sertifikat, error) {
	updateSertifikat, updateSertifikatErr := s.repository.UpdateSertifikat(inputSertifikat, id)

	return updateSertifikat, updateSertifikatErr
}

func (s *service) DeleteSertifikat(id uint) (bool, error) {
	isDeletedSertifikat, isDeletedSertifikatErr := s.repository.DeleteSertifikat(id)

	return isDeletedSertifikat, isDeletedSertifikatErr
}

func (s *service) FindMCUById(id uint) (mcu.MCU, error) {
	sertifikat, err := s.repository.FindMCUById(id)

	return sertifikat, err
}

func (s *service) UpdateMCU(inputMCU RegisterMCUInput, id int) (mcu.MCU, error) {
	updateMCU, updateMCUErr := s.repository.UpdateMCU(inputMCU, id)

	return updateMCU, updateMCUErr
}

func (s *service) DeleteMCU(id uint) (bool, error) {
	isDeletedMCU, isDeletedMCUErr := s.repository.DeleteMCU(id)

	return isDeletedMCU, isDeletedMCUErr
}

func (s *service) FindMCUBerkala(page int, sortFilter SortFilterDohKontrak) (Pagination, error) {
	mcu, err := s.repository.FindMCUBerkala(page, sortFilter)

	return mcu, err
}

func (s *service) FindHistoryById(id uint) (history.History, error) {
	sertifikat, err := s.repository.FindHistoryById(id)

	return sertifikat, err
}

func (s *service) UpdateHistory(inputHistory RegisterHistoryInput, id int) (history.History, error) {
	updateHistory, updateHistoryErr := s.repository.UpdateHistory(inputHistory, id)

	return updateHistory, updateHistoryErr
}

func (s *service) DeleteHistory(id uint) (bool, error) {
	isDeletedHistory, isDeletedHistoryErr := s.repository.DeleteHistory(id)

	return isDeletedHistory, isDeletedHistoryErr
}

func (s *service) GenerateSideBar(userID uint) ([]form.Form, error) {
	form, err := s.repository.GenerateSideBar(userID)

	return form, err
}
