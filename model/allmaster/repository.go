package allmaster

import (
	"encoding/json"
	"fmt"
	"math"
	"mrpbackend/model/master/brand"
	"mrpbackend/model/master/department"
	"mrpbackend/model/master/departmentform"
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
	"mrpbackend/model/master/roleform"
	"mrpbackend/model/master/series"
	"mrpbackend/model/master/sertifikat"
	"mrpbackend/model/master/userrole"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	CreateUserRole(userRoleInput RegisterUserRoleInput) (userrole.UserRole, error)
	CreateBrand(brandInput RegisterBrandInput) (brand.Brand, error)
	CreateHeavyEquipment(heavyEquipmentInput RegisterHeavyEquipmentInput) (heavyequipment.HeavyEquipment, error)
	CreateSeries(seriesInput RegisterSeriesInput) (series.Series, error)
	CreateKartuKeluarga(kartukeluargaInput RegisterKartuKeluargaInput) (kartukeluarga.KartuKeluarga, error)
	CreateKTP(ktpInput RegisterKTPInput) (ktp.KTP, error)
	CreatePendidikan(ktpInput RegisterPendidikanInput) (pendidikan.Pendidikan, error)
	CreateDOH(dohInput RegisterDOHInput) (doh.DOH, error)
	CreateJabatan(jabatanInput RegisterJabatanInput) (jabatan.Jabatan, error)
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

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateUserRole(userRoleInput RegisterUserRoleInput) (userrole.UserRole, error) {
	var newUserRole userrole.UserRole

	newUserRole.UserId = userRoleInput.UserId
	newUserRole.RoleId = userRoleInput.RoleId

	err := r.db.Create(&newUserRole).Error
	if err != nil {
		return newUserRole, err
	}

	return newUserRole, nil
}

func (r *repository) CreateBrand(brandInput RegisterBrandInput) (brand.Brand, error) {
	var newBrand brand.Brand

	newBrand.BrandName = brandInput.BrandName

	err := r.db.Create(&newBrand).Error
	if err != nil {
		return newBrand, err
	}

	return newBrand, nil
}

func (r *repository) CreateHeavyEquipment(heavyEquipmentInput RegisterHeavyEquipmentInput) (heavyequipment.HeavyEquipment, error) {
	var newHeavyEquipment heavyequipment.HeavyEquipment

	newHeavyEquipment.BrandId = heavyEquipmentInput.BrandId
	newHeavyEquipment.HeavyEquipmentName = heavyEquipmentInput.HeavyEquipmentName

	err := r.db.Create(&newHeavyEquipment).Error

	if err != nil {
		return newHeavyEquipment, err
	}

	return newHeavyEquipment, nil
}

func (r *repository) CreateSeries(brandInput RegisterSeriesInput) (series.Series, error) {
	var newSeries series.Series

	newSeries.BrandId = brandInput.BrandId
	newSeries.HeavyEquipmentId = brandInput.HeavyEquipmentId
	newSeries.SeriesName = brandInput.SeriesName

	err := r.db.Create(&newSeries).Error
	if err != nil {
		return newSeries, err
	}

	return newSeries, nil
}

func (r *repository) CreateKartuKeluarga(kartukeluargaInput RegisterKartuKeluargaInput) (kartukeluarga.KartuKeluarga, error) {
	var newKartuKeluarga kartukeluarga.KartuKeluarga

	newKartuKeluarga.NomorKartuKeluarga = kartukeluargaInput.NomorKartuKeluarga
	newKartuKeluarga.NamaIbuKandung = kartukeluargaInput.NamaIbuKandung
	newKartuKeluarga.KontakDarurat = kartukeluargaInput.KontakDarurat
	newKartuKeluarga.NamaKontakDarurat = kartukeluargaInput.NamaKontakDarurat
	newKartuKeluarga.HubunganKontakDarurat = kartukeluargaInput.HubunganKontakDarurat

	err := r.db.Create(&newKartuKeluarga).Error
	if err != nil {
		return newKartuKeluarga, err
	}

	return newKartuKeluarga, nil
}

func (r *repository) CreateKTP(ktpInput RegisterKTPInput) (ktp.KTP, error) {
	var newKTP ktp.KTP

	newKTP.NamaSesuaiKTP = ktpInput.NamaSesuaiKTP
	newKTP.NomorKTP = ktpInput.NomorKTP
	newKTP.TempatLahir = ktpInput.TempatLahir
	newKTP.TanggalLahir = ktpInput.TanggalLahir
	newKTP.Gender = ktpInput.Gender
	newKTP.Alamat = ktpInput.Alamat
	newKTP.RT = ktpInput.RT
	newKTP.RW = ktpInput.RW
	newKTP.Kel = ktpInput.Kel
	newKTP.Kec = ktpInput.Kec
	newKTP.Kota = ktpInput.Kota
	newKTP.Prov = ktpInput.Prov
	newKTP.KodePos = ktpInput.KodePos
	newKTP.GolonganDarah = ktpInput.GolonganDarah
	newKTP.Agama = ktpInput.Agama
	newKTP.RingKTP = ktpInput.RingKTP

	err := r.db.Create(&newKTP).Error
	if err != nil {
		return newKTP, err
	}

	return newKTP, nil
}

func (r *repository) CreatePendidikan(pendidikanInput RegisterPendidikanInput) (pendidikan.Pendidikan, error) {
	var newPendidikan pendidikan.Pendidikan

	newPendidikan.PendidikanLabel = pendidikanInput.PendidikanLabel
	newPendidikan.PendidikanTerakhir = pendidikanInput.PendidikanTerakhir
	newPendidikan.Jurusan = pendidikanInput.Jurusan

	err := r.db.Create(&newPendidikan).Error
	if err != nil {
		return newPendidikan, err
	}

	return newPendidikan, nil
}

func (r *repository) CreateDOH(dohInput RegisterDOHInput) (doh.DOH, error) {
	var newDOH doh.DOH

	newDOH.EmployeeId = dohInput.EmployeeId
	newDOH.TanggalDoh = dohInput.TanggalDoh
	newDOH.TanggalEndDoh = dohInput.TanggalEndDoh
	newDOH.PT = dohInput.PT
	newDOH.Penempatan = dohInput.Penempatan
	newDOH.StatusKontrak = dohInput.StatusKontrak

	err := r.db.Create(&newDOH).Error
	if err != nil {
		return newDOH, err
	}

	return newDOH, nil
}

func (r *repository) CreateJabatan(jabatanInput RegisterJabatanInput) (jabatan.Jabatan, error) {
	var newJabatan jabatan.Jabatan

	newJabatan.EmployeeId = jabatanInput.EmployeeId
	newJabatan.DateMove = jabatanInput.DateMove
	newJabatan.PositionId = jabatanInput.PositionId

	err := r.db.Create(&newJabatan).Error
	if err != nil {
		return newJabatan, err
	}

	// 🔁 Reload the record with associations
	err = r.db.
		Preload("Position").
		First(&newJabatan, newJabatan.ID).Error
	if err != nil {
		return newJabatan, err
	}

	return newJabatan, nil
}

func (r *repository) CreateSertifikat(sertifikatInput RegisterSertifikatInput) (sertifikat.Sertifikat, error) {
	var newSertifikat sertifikat.Sertifikat

	newSertifikat.EmployeeId = sertifikatInput.EmployeeId
	newSertifikat.DateEffective = sertifikatInput.DateEffective
	newSertifikat.Sertifikat = sertifikatInput.Sertifikat
	newSertifikat.Remark = sertifikatInput.Remark

	err := r.db.Create(&newSertifikat).Error
	if err != nil {
		return newSertifikat, err
	}

	return newSertifikat, nil
}

func (r *repository) CreateMCU(mcuInput RegisterMCUInput) (mcu.MCU, error) {
	var newMCU mcu.MCU

	newMCU.EmployeeId = mcuInput.EmployeeId
	newMCU.DateMCU = mcuInput.DateMCU
	newMCU.DateEndMCU = mcuInput.DateEndMCU
	newMCU.HasilMCU = mcuInput.HasilMCU
	newMCU.MCU = mcuInput.MCU

	err := r.db.Create(&newMCU).Error
	if err != nil {
		return newMCU, err
	}

	return newMCU, nil
}

func (r *repository) CreateHistory(historyInput RegisterHistoryInput) (history.History, error) {
	var newHistory history.History

	newHistory.EmployeeId = historyInput.EmployeeId
	newHistory.StatusTerakhir = historyInput.StatusTerakhir
	newHistory.Tanggal = historyInput.Tanggal
	newHistory.Keterangan = historyInput.Keterangan

	err := r.db.Create(&newHistory).Error
	if err != nil {
		return newHistory, err
	}

	return newHistory, nil
}

func (r *repository) FindUserRole() ([]userrole.UserRole, error) {
	var userRole []userrole.UserRole

	errFind := r.db.Preload("Role").Find(&userRole).Error

	return userRole, errFind
}

func (r *repository) FindUserRoleById(id uint) (userrole.UserRole, error) {
	var userRole userrole.UserRole

	errFind := r.db.Preload("Role").Where("id = ?", id).First(&userRole).Error

	return userRole, errFind
}

func (r *repository) FindBrand() ([]brand.Brand, error) {
	var brand []brand.Brand

	errFind := r.db.Find(&brand).Error

	return brand, errFind
}

func (r *repository) FindBrandById(id uint) (brand.Brand, error) {
	var brand brand.Brand

	errFind := r.db.Where("id = ?", id).First(&brand).Error

	return brand, errFind
}

func (r *repository) FindHeavyEquipment() ([]heavyequipment.HeavyEquipment, error) {
	var heavyEquipment []heavyequipment.HeavyEquipment

	query := `SELECT DISTINCT ON (heavy_equipment_name) * FROM heavy_equipments ORDER BY heavy_equipment_name, id`

	err := r.db.Raw(query).Scan(&heavyEquipment).Error
	return heavyEquipment, err
}

func (r *repository) FindHeavyEquipmentById(id uint) (heavyequipment.HeavyEquipment, error) {
	var heavyEquipment heavyequipment.HeavyEquipment

	errFind := r.db.Preload("Brand").Where("id = ?", id).First(&heavyEquipment).Error
	return heavyEquipment, errFind
}

func (r *repository) FindHeavyEquipmentByBrandID(brandId uint) ([]heavyequipment.HeavyEquipment, error) {
	var heavyEquipment []heavyequipment.HeavyEquipment

	err := r.db.
		Preload("Brand").
		Where("brand_id = ?", brandId).
		Find(&heavyEquipment).Error

	return heavyEquipment, err
}

func (r *repository) FindSeries() ([]series.Series, error) {
	var series []series.Series

	query := `SELECT DISTINCT ON (series_name) * FROM series ORDER BY series_name, id`

	err := r.db.Raw(query).Scan(&series).Error
	return series, err
}

func (r *repository) FindSeriesById(id uint) (series.Series, error) {
	var series series.Series

	errFind := r.db.
		Preload("Brand").
		Preload("HeavyEquipment").
		Where("id = ?", id).First(&series).Error
	return series, errFind
}

func (r *repository) FindSeriesByBrandAndEquipmentdID(brandId uint, heavyequipmentId uint) ([]series.Series, error) {
	var series []series.Series

	err := r.db.
		Preload("Brand").
		Preload("HeavyEquipment").
		Where("brand_id = ? AND heavy_equipment_id = ?", brandId, heavyequipmentId).
		Find(&series).Error

	return series, err
}

func (r *repository) FindDepartment() ([]department.Department, error) {
	var department []department.Department

	errFind := r.db.Find(&department).Error

	return department, errFind
}

func (r *repository) FindRole() ([]role.Role, error) {
	var series []role.Role

	query := `SELECT DISTINCT ON (name) * FROM roles ORDER BY name, id`

	err := r.db.Raw(query).Scan(&series).Error
	return series, err
}

func (r *repository) FindPosition() ([]position.Position, error) {
	var positions []position.Position

	query := `SELECT DISTINCT ON (position_name) * FROM positions ORDER BY position_name, id`

	err := r.db.Raw(query).Scan(&positions).Error
	return positions, err
}

func findFormIDByTitle(forms []form.Form, title string) uint {
	for _, f := range forms {
		if f.FormName == title {
			return f.ID
		}
	}
	return 0
}

func (r *repository) FindDOHById(id uint) (doh.DOH, error) {
	var doh doh.DOH

	errFind := r.db.
		Where("id = ?", id).First(&doh).Error
	return doh, errFind
}

func (r *repository) UpdateDOH(inputDOH RegisterDOHInput, id int) (doh.DOH, error) {

	var updatedDOH doh.DOH
	errFind := r.db.Where("id = ?", id).First(&updatedDOH).Error

	if errFind != nil {
		return updatedDOH, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputDOH)

	if errorMarshal != nil {
		return updatedDOH, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedDOH, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedDOH).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedDOH, updateErr
	}

	return updatedDOH, nil
}

func (r *repository) DeleteDOH(id uint) (bool, error) {
	tx := r.db.Begin()
	var doh doh.DOH

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&doh).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&doh).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}

func (r *repository) FindDohKontrak(page int, sortFilter SortFilterDohKontrak) (Pagination, error) {
	var results []EmployeeDOHExpired
	var pagination Pagination

	// Validate and compute pagination
	pagination.Page = page
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	pagination.Limit = 10
	offset := pagination.GetOffset()

	queryFilter := "dohs.id > 0"
	ptConditions := []string{}

	// === Handle PT Logic ===
	ptList := []string{}
	if strings.TrimSpace(sortFilter.PT) != "" {
		ptList = strings.Split(sortFilter.PT, ",")
	} else {
		// Default to PTs based on empCode
		if sortFilter.CodeEmp == "1" {
			ptList = []string{"PT. MRP", "PT. TRIOP"}
		} else if sortFilter.CodeEmp == "2" {
			ptList = []string{"PT. MJSU", "PT. IBS"}
		}
	}

	for _, pt := range ptList {
		pt = strings.TrimSpace(pt)
		switch pt {
		case "PT. MRP":
			ptConditions = append(ptConditions, "cast(nomor_karyawan AS TEXT) ILIKE '%MRP%'")
		case "PT. TRIOP":
			ptConditions = append(ptConditions, "cast(nomor_karyawan AS TEXT) ILIKE '%TRIOP%'")
		case "PT. MJSU":
			ptConditions = append(ptConditions, "cast(nomor_karyawan AS TEXT) ILIKE '%MJSU%'")
		case "PT. IBS":
			ptConditions = append(ptConditions, "cast(nomor_karyawan AS TEXT) ILIKE '%IBS%'")
		}
	}

	if len(ptConditions) > 0 {
		queryFilter += " AND (" + strings.Join(ptConditions, " OR ") + ")"
	} else {
		queryFilter += " AND (1=0)" // Fallback: No valid PT
	}

	if strings.TrimSpace(sortFilter.Year) == "" {
		// Set to current year
		sortFilter.Year = fmt.Sprintf("%d", time.Now().Year())
	}

	// BASED ON MONTH
	queryFilterAll := "EXTRACT(YEAR FROM dohs.tanggal_end_doh::DATE) = " + sortFilter.Year + " AND " + queryFilter

	// Default sort
	querySort := "employee_id"
	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	//CodeEmp
	//1 -> MRP/TRIOP
	//2 -> MJSU/IBS
	if sortFilter.CodeEmp == "1" {
		queryFilter += " AND (cast(e.nomor_karyawan AS TEXT) ILIKE '%MRP%' OR cast(e.nomor_karyawan AS TEXT) ILIKE '%TRIOP%')"
		queryFilter += " AND cast(e.nomor_karyawan AS TEXT) NOT ILIKE '%MJSU%'"
		queryFilter += " AND cast(e.nomor_karyawan AS TEXT) NOT ILIKE '%IBS%'"
		queryFilter += " AND cast(status AS TEXT) ILIKE 'AKTIF'"
	}

	if sortFilter.CodeEmp == "2" {
		queryFilter += " AND (cast(e.nomor_karyawan AS TEXT) ILIKE '%MJSU%' OR cast(e.nomor_karyawan AS TEXT) ILIKE '%IBS%')"
		queryFilter += " AND cast(e.nomor_karyawan AS TEXT) NOT ILIKE '%MRP%'"
		queryFilter += " AND cast(e.nomor_karyawan AS TEXT) NOT ILIKE '%TRIOP%'"
		queryFilter += " AND cast(status AS TEXT) ILIKE 'AKTIF'"
	}

	// Paginated query
	query := `
		WITH latest_dohs AS (
			SELECT DISTINCT ON (dohs.employee_id)
				dohs.*, 
				e.firstname AS firstname,
				e.lastname AS lastname, 
				d.department_name, 
				p.position_name
			FROM dohs
			JOIN employees e ON dohs.employee_id = e.id
			JOIN departments d ON e.department_id = d.id
			JOIN positions p ON e.position_id = p.id
			WHERE dohs.deleted_at IS NULL  AND ` + queryFilterAll + `
			ORDER BY dohs.employee_id, TO_DATE(dohs.tanggal_doh, 'YYYY-MM-DD') DESC
		)
		SELECT * FROM latest_dohs
		WHERE CURRENT_DATE >= TO_DATE(tanggal_end_doh, 'YYYY-MM-DD') - INTERVAL '1 month'
		ORDER BY ` + querySort + `
		LIMIT ? OFFSET ?
	`
	if err := r.db.Raw(query, pagination.Limit, offset).Scan(&results).Error; err != nil {
		return pagination, err
	}

	// Count query
	countQuery := `
		WITH latest_dohs AS (
			SELECT DISTINCT ON (dohs.employee_id)
				dohs.*, 
				e.firstname AS firstname,
				e.lastname AS lastname, 
				d.department_name, 
				p.position_name
			FROM dohs
			JOIN employees e ON dohs.employee_id = e.id
			JOIN departments d ON e.department_id = d.id
			JOIN positions p ON e.position_id = p.id
			WHERE dohs.deleted_at IS NULL AND ` + queryFilterAll + `
			ORDER BY dohs.employee_id, TO_DATE(dohs.tanggal_doh, 'YYYY-MM-DD') DESC
		)
		SELECT count(*) FROM latest_dohs
		WHERE CURRENT_DATE >= TO_DATE(tanggal_end_doh, 'YYYY-MM-DD') - INTERVAL '1 month'
	`
	var totalRows int64
	if err := r.db.Raw(countQuery).Scan(&totalRows).Error; err != nil {
		return pagination, err
	}

	// Fill pagination struct
	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.Data = results

	return pagination, nil
}

func (r *repository) FindJabatanById(id uint) (jabatan.Jabatan, error) {
	var jabatan jabatan.Jabatan

	errFind := r.db.
		Where("id = ?", id).First(&jabatan).Error
	return jabatan, errFind
}

func (r *repository) UpdateJabatan(inputJabatan RegisterJabatanInput, id int) (jabatan.Jabatan, error) {

	var updatedJabatan jabatan.Jabatan
	errFind := r.db.Where("id = ?", id).First(&updatedJabatan).Error

	if errFind != nil {
		return updatedJabatan, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputJabatan)

	if errorMarshal != nil {
		return updatedJabatan, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedJabatan, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedJabatan).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedJabatan, updateErr
	}
	// 🔁 Reload the record with associations
	err := r.db.
		Preload("Position").
		First(&updatedJabatan, updatedJabatan.ID).Error

	if err != nil {
		return updatedJabatan, err
	}

	return updatedJabatan, nil
}

func (r *repository) DeleteJabatan(id uint) (bool, error) {
	tx := r.db.Begin()
	var jabatan jabatan.Jabatan

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&jabatan).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&jabatan).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}

func (r *repository) FindSertifikatById(id uint) (sertifikat.Sertifikat, error) {
	var sertifikat sertifikat.Sertifikat

	errFind := r.db.
		Where("id = ?", id).First(&sertifikat).Error
	return sertifikat, errFind
}

func (r *repository) UpdateSertifikat(inputSertifikat RegisterSertifikatInput, id int) (sertifikat.Sertifikat, error) {

	var updatedSertifikat sertifikat.Sertifikat
	errFind := r.db.Where("id = ?", id).First(&updatedSertifikat).Error

	if errFind != nil {
		return updatedSertifikat, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputSertifikat)

	if errorMarshal != nil {
		return updatedSertifikat, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedSertifikat, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedSertifikat).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedSertifikat, updateErr
	}

	return updatedSertifikat, nil
}

func (r *repository) DeleteSertifikat(id uint) (bool, error) {
	tx := r.db.Begin()
	var sertifikat sertifikat.Sertifikat

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&sertifikat).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&sertifikat).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}

func (r *repository) FindMCUById(id uint) (mcu.MCU, error) {
	var mcu mcu.MCU

	errFind := r.db.
		Where("id = ?", id).First(&mcu).Error
	return mcu, errFind
}

func (r *repository) UpdateMCU(inputMCU RegisterMCUInput, id int) (mcu.MCU, error) {

	var updatedMCU mcu.MCU
	errFind := r.db.Where("id = ?", id).First(&updatedMCU).Error

	if errFind != nil {
		return updatedMCU, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputMCU)

	if errorMarshal != nil {
		return updatedMCU, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedMCU, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedMCU).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedMCU, updateErr
	}

	return updatedMCU, nil
}

func (r *repository) DeleteMCU(id uint) (bool, error) {
	tx := r.db.Begin()
	var mcu mcu.MCU

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&mcu).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&mcu).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}

func (r *repository) FindMCUBerkala(page int, sortFilter SortFilterDohKontrak) (Pagination, error) {
	var results []EmployeeMCUBerkala
	var pagination Pagination

	// Validate and compute pagination
	pagination.Page = page
	if pagination.Page < 1 {
		pagination.Page = 1
	}
	pagination.Limit = 10
	offset := pagination.GetOffset()

	queryFilter := "mcus.id > 0"
	ptConditions := []string{}

	// === Handle PT Logic ===
	ptList := []string{}
	if strings.TrimSpace(sortFilter.PT) != "" {
		ptList = strings.Split(sortFilter.PT, ",")
	} else {
		// Default to PTs based on empCode
		if sortFilter.CodeEmp == "1" {
			ptList = []string{"PT. MRP", "PT. TRIOP"}
		} else if sortFilter.CodeEmp == "2" {
			ptList = []string{"PT. MJSU", "PT. IBS"}
		}
	}

	for _, pt := range ptList {
		pt = strings.TrimSpace(pt)
		switch pt {
		case "PT. MRP":
			ptConditions = append(ptConditions, "cast(nomor_karyawan AS TEXT) ILIKE '%MRP%'")
		case "PT. TRIOP":
			ptConditions = append(ptConditions, "cast(nomor_karyawan AS TEXT) ILIKE '%TRIOP%'")
		case "PT. MJSU":
			ptConditions = append(ptConditions, "cast(nomor_karyawan AS TEXT) ILIKE '%MJSU%'")
		case "PT. IBS":
			ptConditions = append(ptConditions, "cast(nomor_karyawan AS TEXT) ILIKE '%IBS%'")
		}
	}

	if len(ptConditions) > 0 {
		queryFilter += " AND (" + strings.Join(ptConditions, " OR ") + ")"
	} else {
		queryFilter += " AND (1=0)" // Fallback: No valid PT
	}

	if strings.TrimSpace(sortFilter.Year) == "" {
		// Set to current year
		sortFilter.Year = fmt.Sprintf("%d", time.Now().Year())
	}

	// BASED ON MONTH
	queryFilterAll := "EXTRACT(YEAR FROM mcus.date_end_mcu::DATE) = " + sortFilter.Year + " AND " + queryFilter

	// Default sort
	querySort := "employee_id"
	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}

	//CodeEmp
	//1 -> MRP/TRIOP
	//2 -> MJSU/IBS
	if sortFilter.CodeEmp == "1" {
		queryFilter += " AND (cast(e.nomor_karyawan AS TEXT) ILIKE '%MRP%' OR cast(e.nomor_karyawan AS TEXT) ILIKE '%TRIOP%')"
		queryFilter += " AND cast(e.nomor_karyawan AS TEXT) NOT ILIKE '%MJSU%'"
		queryFilter += " AND cast(e.nomor_karyawan AS TEXT) NOT ILIKE '%IBS%'"
		queryFilter += " AND cast(status AS TEXT) ILIKE 'AKTIF'"
	}

	if sortFilter.CodeEmp == "2" {
		queryFilter += " AND (cast(e.nomor_karyawan AS TEXT) ILIKE '%MJSU%' OR cast(e.nomor_karyawan AS TEXT) ILIKE '%IBS%')"
		queryFilter += " AND cast(e.nomor_karyawan AS TEXT) NOT ILIKE '%MRP%'"
		queryFilter += " AND cast(e.nomor_karyawan AS TEXT) NOT ILIKE '%TRIOP%'"
		queryFilter += " AND cast(status AS TEXT) ILIKE 'AKTIF'"
	}

	// Main query
	query := `
		WITH latest_mcus AS (
			SELECT DISTINCT ON (mcus.employee_id)
				mcus.*,
				e.firstname AS firstname,
				e.lastname AS lastname,
				d.department_name,
				p.position_name
			FROM mcus
			JOIN employees e ON mcus.employee_id = e.id
			JOIN departments d ON e.department_id = d.id
			JOIN positions p ON e.position_id = p.id
			WHERE mcus.deleted_at IS NULL AND ` + queryFilterAll + `
			ORDER BY mcus.employee_id, TO_DATE(mcus.date_end_mcu, 'YYYY-MM-DD') DESC
		)
		SELECT * FROM latest_mcus
		WHERE CURRENT_DATE >= TO_DATE(date_end_mcu, 'YYYY-MM-DD') - INTERVAL '1 month'
		ORDER BY ` + querySort + `
		LIMIT ? OFFSET ?
	`
	if err := r.db.Raw(query, pagination.Limit, offset).Scan(&results).Error; err != nil {
		return pagination, err
	}

	// Count query
	countQuery := `
		WITH latest_mcus AS (
			SELECT DISTINCT ON (mcus.employee_id)
				mcus.*,
				e.firstname AS firstname,
				e.lastname AS lastname,
				d.department_name,
				p.position_name
			FROM mcus
			JOIN employees e ON mcus.employee_id = e.id
			JOIN departments d ON e.department_id = d.id
			JOIN positions p ON e.position_id = p.id
			WHERE mcus.deleted_at IS NULL AND ` + queryFilterAll + `
			ORDER BY mcus.employee_id, TO_DATE(mcus.date_end_mcu, 'YYYY-MM-DD') DESC
		)
		SELECT count(*) FROM latest_mcus
		WHERE CURRENT_DATE >= TO_DATE(date_end_mcu, 'YYYY-MM-DD') - INTERVAL '1 month' 
	`
	var totalRows int64
	if err := r.db.Raw(countQuery).Scan(&totalRows).Error; err != nil {
		return pagination, err
	}

	// Finalize pagination
	pagination.TotalRows = totalRows
	pagination.TotalPages = int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.Data = results

	return pagination, nil
}

func (r *repository) FindHistoryById(id uint) (history.History, error) {
	var history history.History

	errFind := r.db.
		Where("id = ?", id).First(&history).Error
	return history, errFind
}

func (r *repository) UpdateHistory(inputHistory RegisterHistoryInput, id int) (history.History, error) {

	var updatedHistory history.History
	errFind := r.db.Where("id = ?", id).First(&updatedHistory).Error

	if errFind != nil {
		return updatedHistory, errFind
	}

	dataInput, errorMarshal := json.Marshal(inputHistory)

	if errorMarshal != nil {
		return updatedHistory, errorMarshal
	}

	var dataInputMapString map[string]interface{}

	errorUnmarshal := json.Unmarshal(dataInput, &dataInputMapString)

	if errorUnmarshal != nil {
		return updatedHistory, errorUnmarshal
	}

	updateErr := r.db.Model(&updatedHistory).Updates(dataInputMapString).Error

	if updateErr != nil {
		return updatedHistory, updateErr
	}

	return updatedHistory, nil
}

func (r *repository) DeleteHistory(id uint) (bool, error) {
	tx := r.db.Begin()
	var history history.History

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&history).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&history).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}

func (r *repository) GenerateSideBar(userID uint) ([]form.Form, error) {
	// Step 1: Get department and role for user
	var userRoles []userrole.UserRole
	if err := r.db.Where("user_id = ?", userID).Find(&userRoles).Error; err != nil {
		return nil, err
	}

	var deptIDs, roleIDs []uint
	for _, ur := range userRoles {
		deptIDs = append(deptIDs, ur.DepartmentId)
		roleIDs = append(roleIDs, ur.RoleId)
	}

	// Step 2: Get allowed department_form
	var deptForms []departmentform.DepartmentForm
	if err := r.db.Where("department_id IN ? AND role_id IN ?", deptIDs, roleIDs).Find(&deptForms).Error; err != nil {
		return nil, err
	}

	var deptFormIDs []uint
	for _, df := range deptForms {
		deptFormIDs = append(deptFormIDs, df.ID)
	}

	// Step 3: Get allowed role_form
	var roleForms []roleform.RoleForm
	if err := r.db.Where("department_form_id IN ?", deptFormIDs).Find(&roleForms).Error; err != nil {
		return nil, err
	}

	var formIDs []uint
	for _, rf := range roleForms {
		formIDs = append(formIDs, rf.FormId)
	}

	flagMap := make(map[uint]roleform.RoleForm)
	for _, rf := range roleForms {
		flagMap[rf.FormId] = rf
	}

	// Step 4: Get allowed forms
	var forms []form.Form
	if err := r.db.Where("id IN ?", formIDs).Order("sequence ASC").Find(&forms).Error; err != nil {
		return nil, err
	}

	// Step 5: Build form tree with flags
	formMap := make(map[uint][]form.Form)

	for _, f := range forms {
		roleFlags := flagMap[f.ID] // assume this map contains permissions

		formItem := form.Form{
			ID:         f.ID,
			FormName:   f.FormName,
			Path:       f.Path,
			Sequence:   f.Sequence,
			CreateFlag: roleFlags.CreateFlag,
			UpdateFlag: roleFlags.UpdateFlag,
			ReadFlag:   roleFlags.ReadFlag,
			DeleteFlag: roleFlags.DeleteFlag,
		}

		parentID := uint(0)
		if f.ParentID != nil {
			parentID = *f.ParentID
		}
		formMap[parentID] = append(formMap[parentID], formItem)
	}

	// Step 6: Recursive tree assembly
	var buildTree func(parentID uint) []form.Form
	buildTree = func(parentID uint) []form.Form {
		items := formMap[parentID]
		for i := range items {
			items[i].Children = buildTree(items[i].ID)
		}
		return items
	}

	tree := buildTree(0)
	return tree, nil
}
