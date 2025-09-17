package employee

import (
	"fmt"
	"mrpbackend/model/master/apd"
	"mrpbackend/model/master/bank"
	"mrpbackend/model/master/bpjskesehatan"
	"mrpbackend/model/master/bpjsketenagakerjaan"
	"mrpbackend/model/master/kartukeluarga"
	"mrpbackend/model/master/ktp"
	"mrpbackend/model/master/laporan"
	"mrpbackend/model/master/npwp"
	"mrpbackend/model/master/pendidikan"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	CreateEmployee(employees RegisterEmployeeInput) (Employee, error)
	FindEmployee(empCode uint) ([]Employee, error)
	FindEmployeeById(id uint) (Employee, error)
	FindEmployeeByDepartmentId(departmentId uint) ([]Employee, error)
	ListEmployee(page int, sortFilter SortFilterEmployee) (Pagination, error)
	UpdateEmployee(inputEmployee UpdateEmployeeInput, id int) (Employee, error)
	DeleteEmployee(id uint) (bool, error)
	ListDashboard(empCode uint, dashboardSort SortFilterDashboardEmployee) (DashboardEmployee, error)
	ListDashboardTurnover(empCode uint, dashboardSort SortFilterDashboardEmployeeTurnOver) (DashboardEmployeeTurnOver, error)
	ListDashboardKontrak(empCode uint, dashboardSort SortFilterDashboardEmployeeTurnOver) (DashboardEmployeeKontrak, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) CreateEmployee(EmployeeInput RegisterEmployeeInput) (Employee, error) {
	var newEmployee Employee

	var newKartuKeluarga kartukeluarga.KartuKeluarga = EmployeeInput.KartuKeluarga
	if newKartuKeluarga.NomorKartuKeluarga != "" {
		if err := r.db.Create(&newKartuKeluarga).Error; err != nil {
			return Employee{}, err
		}
	}

	kartuKeluargaId := newKartuKeluarga.ID

	var newKTP ktp.KTP = EmployeeInput.KTP
	if newKTP.NomorKTP != "" {
		if err := r.db.Create(&newKTP).Error; err != nil {
			return Employee{}, err
		}
	}

	ktpId := newKTP.ID

	var newPendidikan pendidikan.Pendidikan = EmployeeInput.Pendidikan
	if newPendidikan.PendidikanTerakhir != "" {
		if err := r.db.Create(&newPendidikan).Error; err != nil {
			return Employee{}, err
		}
	}

	pendidikanId := newPendidikan.ID

	var newLaporan laporan.Laporan = EmployeeInput.Laporan
	if newLaporan.RingSerapan != "" {
		if err := r.db.Create(&newLaporan).Error; err != nil {
			return Employee{}, err
		}
	}

	laporanId := newLaporan.ID

	var newAPD apd.APD = *EmployeeInput.APD

	if newAPD.UkuranBaju != "" {
		if err := r.db.Create(&newAPD).Error; err != nil {
			return Employee{}, err
		}
	}

	newEmployee.APDId = &newAPD.ID

	var newNPWP npwp.NPWP = *EmployeeInput.NPWP
	if newNPWP.StatusPajak != "" {
		if err := r.db.Create(&newNPWP).Error; err != nil {
			return Employee{}, err
		}
	}

	npwpId := newNPWP.ID

	var newBank bank.Bank = EmployeeInput.Bank
	if newBank.NamaBank != "" {
		if err := r.db.Create(&newBank).Error; err != nil {
			return Employee{}, err
		}
	}

	bankId := newBank.ID

	var newBPJSKesehatan bpjskesehatan.BPJSKesehatan = *EmployeeInput.BPJSKesehatan
	if newBPJSKesehatan.NomorKesehatan != "" {
		if err := r.db.Create(&newBPJSKesehatan).Error; err != nil {
			return Employee{}, err
		}
		newEmployee.BPJSKesehatanId = &newBPJSKesehatan.ID
	}

	var newBPJSKetenagakerjaan bpjsketenagakerjaan.BPJSKetenagakerjaan = *EmployeeInput.BPJSKetenagakerjaan
	if newBPJSKetenagakerjaan.NomorKetenagakerjaan != "" {
		if err := r.db.Create(&newBPJSKetenagakerjaan).Error; err != nil {
			return Employee{}, err
		}
		newEmployee.BPJSKetenagakerjaanId = &newBPJSKetenagakerjaan.ID
	}

	newEmployee.NomorKaryawan = EmployeeInput.NomorKaryawan
	newEmployee.DepartmentId = EmployeeInput.DepartmentId
	newEmployee.Firstname = EmployeeInput.Firstname
	newEmployee.Lastname = EmployeeInput.Lastname
	newEmployee.PhoneNumber = EmployeeInput.PhoneNumber
	newEmployee.Email = EmployeeInput.Email
	newEmployee.Level = EmployeeInput.Level
	newEmployee.Status = EmployeeInput.Status
	newEmployee.HiredBy = EmployeeInput.HiredBy
	newEmployee.RoleId = EmployeeInput.RoleId
	newEmployee.KartuKeluargaId = kartuKeluargaId
	newEmployee.KTPId = ktpId
	newEmployee.PendidikanId = pendidikanId
	newEmployee.LaporanId = laporanId
	newEmployee.BankId = bankId
	newEmployee.NPWPId = &npwpId
	newEmployee.PositionId = EmployeeInput.PositionId
	newEmployee.DateOfHire = EmployeeInput.DateOfHire

	err := r.db.Create(&newEmployee).Error
	if err != nil {
		return newEmployee, err
	}

	employeeId := newEmployee.ID

	for _, doh := range EmployeeInput.DOH {
		doh.EmployeeId = employeeId
		r.db.Create(&doh)
	}

	for _, jab := range EmployeeInput.Jabatan {
		jab.EmployeeId = employeeId
		r.db.Create(&jab)
	}

	for _, sert := range EmployeeInput.Sertifikat {
		sert.EmployeeId = employeeId
		r.db.Create(&sert)
	}

	for _, mcu := range EmployeeInput.MCU {
		mcu.EmployeeId = employeeId
		r.db.Create(&mcu)
	}

	if EmployeeInput.History != nil {
		for _, h := range *EmployeeInput.History {
			h.EmployeeId = employeeId
			r.db.Create(&h)
		}
	}

	return newEmployee, nil
}

func (r *repository) FindEmployee(empCode uint) ([]Employee, error) {
	var employees []Employee

	queryFilter := "employees.id > 0"

	// Exclude PTs depending on empCode
	if empCode == 1 {
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%MJSU%'"
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%IBS%'"
	} else if empCode == 2 {
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%MRP%'"
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%TRIOP%'"
	}

	errFind := r.db.
		Preload("Department").
		Preload("Role").
		Preload("Position").
		Preload("KartuKeluarga").
		Preload("KTP").
		Preload("Pendidikan").
		Preload("Laporan").
		Preload("APD").
		Preload("NPWP").
		Preload("Bank").
		Preload("BPJSKesehatan").
		Preload("BPJSKetenagakerjaan").
		Preload("DOH").
		Preload("Jabatan").
		Preload("Jabatan.Position").
		Preload("Sertifikat").
		Preload("MCU").
		Preload("History").
		Where(queryFilter).
		Order("nomor_karyawan ASC").Find(&employees).Error

	return employees, errFind
}

func (r *repository) FindEmployeeById(id uint) (Employee, error) {
	var employees Employee

	orderByID := func(db *gorm.DB) *gorm.DB {
		return db.Order("id ASC")
	}
	errFind := r.db.
		Preload("Department", orderByID).
		Preload("Role", orderByID).
		Preload("Position", orderByID).
		Preload("KartuKeluarga", orderByID).
		Preload("KTP", orderByID).
		Preload("Pendidikan", orderByID).
		Preload("Laporan", orderByID).
		Preload("APD", orderByID).
		Preload("NPWP", orderByID).
		Preload("Bank", orderByID).
		Preload("BPJSKesehatan", orderByID).
		Preload("BPJSKetenagakerjaan", orderByID).
		Preload("DOH", orderByID).
		Preload("Jabatan", orderByID).
		Preload("Jabatan.Position", orderByID).
		Preload("Sertifikat", orderByID).
		Preload("MCU", orderByID).
		Preload("History", orderByID).
		Where("id = ?", id).First(&employees).Error
	return employees, errFind
}

func (r *repository) FindEmployeeByDepartmentId(departmentId uint) ([]Employee, error) {
	var employees []Employee

	errFind := r.db.Preload("Department").Where("department_id = ?", departmentId).Find(&employees).Error

	return employees, errFind
}

func (r *repository) ListEmployee(page int, sortFilter SortFilterEmployee) (Pagination, error) {
	var listEmployee []Employee
	var pagination Pagination

	pagination.Limit = 10
	pagination.Page = page
	queryFilter := "employees.id > 0"
	querySort := "employees.id desc"

	if sortFilter.Field != "" && sortFilter.Sort != "" {
		querySort = sortFilter.Field + " " + sortFilter.Sort
	}
	//CodeEmp
	//1 -> MRP/TRIOP
	//2 -> MJSU/IBS
	if sortFilter.CodeEmp == "1" {
		queryFilter += " AND (cast(nomor_karyawan AS TEXT) ILIKE '%MRP%' OR cast(nomor_karyawan AS TEXT) ILIKE '%TRIOP%')"
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%MJSU%'"
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%IBS%'"

		if sortFilter.NomorKaryawan != "" {
			queryFilter += " AND cast(nomor_karyawan AS TEXT) ILIKE '%" + sortFilter.NomorKaryawan + "%'"
		}
	}

	if sortFilter.CodeEmp == "2" {
		queryFilter += " AND (cast(nomor_karyawan AS TEXT) ILIKE '%MJSU%' OR cast(nomor_karyawan AS TEXT) ILIKE '%IBS%')"
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%MRP%'"
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%TRIOP%'"

		if sortFilter.NomorKaryawan != "" {
			queryFilter += " AND cast(nomor_karyawan AS TEXT) ILIKE '%" + sortFilter.NomorKaryawan + "%'"
		}
	}

	if sortFilter.DepartmentId != "" {
		queryFilter = queryFilter + " AND department_id = " + sortFilter.DepartmentId
	}

	if sortFilter.Firstname != "" {
		queryFilter = queryFilter + " AND (CAST(firstname AS TEXT) || ' ' || CAST(lastname AS TEXT)) ILIKE '%" + sortFilter.Firstname + "%'"
	}

	if sortFilter.HireBy != "" {
		queryFilter = queryFilter + " AND cast(hired_by AS TEXT) ILIKE '%" + sortFilter.HireBy + "%'"
	}

	if sortFilter.Agama != "" {
		queryFilter = queryFilter + " AND cast(ktps.agama AS TEXT) LIKE '%" + sortFilter.Agama + "%'"
	}

	if sortFilter.Level != "" {
		queryFilter = queryFilter + " AND cast(level AS TEXT) LIKE '%" + sortFilter.Level + "%'"
	}

	if sortFilter.Gender != "" {
		queryFilter = queryFilter + " AND cast(gender AS TEXT) ILIKE '" + sortFilter.Gender + "%'"
	}

	if sortFilter.KategoriLokalNonLokal != "" {
		queryFilter = queryFilter + " AND cast(laporans.kategori_lokal_non_lokal AS TEXT) LIKE '" + sortFilter.KategoriLokalNonLokal + "'"
	}

	if sortFilter.KategoriTriwulan != "" {
		queryFilter = queryFilter + " AND cast(laporans.kategori_laporan_twiwulan AS TEXT) LIKE '%" + sortFilter.KategoriTriwulan + "%'"
	}

	if sortFilter.Status != "" {
		queryFilter = queryFilter + " AND cast(status AS TEXT) LIKE '%" + sortFilter.Status + "%'"
	}

	if sortFilter.Kontrak != "" {
		queryFilter = queryFilter + " AND cast(latest_doh.status_kontrak AS TEXT) LIKE '%" + sortFilter.Kontrak + "%'"
	}

	if sortFilter.RoleId != "" {
		queryFilter = queryFilter + " AND role_id = " + sortFilter.RoleId
	}

	if sortFilter.PositionId != "" {
		queryFilter = queryFilter + " AND position_id = " + sortFilter.PositionId
	}

	query := r.db.
		Model(&Employee{}).
		Joins("LEFT JOIN ktps ON ktps.id = employees.ktp_id").
		Joins("LEFT JOIN laporans ON laporans.id = employees.laporan_id").
		Joins(`
			LEFT JOIN (
				SELECT DISTINCT ON (employee_id) *
				FROM dohs
				ORDER BY employee_id, tanggal_doh DESC
			) latest_doh ON latest_doh.employee_id = employees.id
		`).
		Where(queryFilter)

	paginationScope := paginateData(&Employee{}, &pagination, query)

	errFind := query.
		Preload("Department").
		Preload("Role").
		Preload("Position").
		Preload("KartuKeluarga").
		Preload("KTP").
		Preload("Pendidikan").
		Preload("Laporan").
		Preload("APD").
		Preload("NPWP").
		Preload("Bank").
		Preload("BPJSKesehatan").
		Preload("BPJSKetenagakerjaan").
		Preload("DOH").
		Preload("Jabatan").
		Preload("Sertifikat").
		Preload("MCU").
		Preload("History").
		Preload(clause.Associations).
		Order(querySort).Scopes(paginationScope).Find(&listEmployee).Error

	if errFind != nil {
		return pagination, errFind
	}

	pagination.Data = listEmployee

	return pagination, nil
}

func UintToPtr(u uint) *uint {
	if u == 0 {
		return nil
	}
	return &u
}

func (r *repository) UpdateEmployee(inputEmployee UpdateEmployeeInput, id int) (Employee, error) {

	var updatedEmployee Employee

	var kartuKeluargaId uint

	if inputEmployee.KartuKeluarga.ID > 0 {
		// Update existing
		if err := r.db.Model(&kartukeluarga.KartuKeluarga{}).
			Where("id = ?", inputEmployee.KartuKeluarga.ID).
			Updates(inputEmployee.KartuKeluarga).Error; err != nil {
			return updatedEmployee, err
		}
		kartuKeluargaId = inputEmployee.KartuKeluarga.ID
	} else {
		// Create new

		var newKartuKeluarga kartukeluarga.KartuKeluarga = inputEmployee.KartuKeluarga

		if newKartuKeluarga.NomorKartuKeluarga != "" {
			if err := r.db.Create(&newKartuKeluarga).Error; err != nil {
				return Employee{}, err
			}
		}
		kartuKeluargaId = newKartuKeluarga.ID
	}

	var ktpId uint

	if inputEmployee.KTP.ID > 0 {
		// Update existing
		if err := r.db.Model(&ktp.KTP{}).
			Where("id = ?", inputEmployee.KTP.ID).
			Updates(inputEmployee.KTP).Error; err != nil {
			return updatedEmployee, err
		}
		ktpId = inputEmployee.KTP.ID
	} else {
		// Create new
		var newKTP ktp.KTP = inputEmployee.KTP

		if newKTP.NomorKTP != "" {
			if err := r.db.Create(&newKTP).Error; err != nil {
				return Employee{}, err
			}
		}
		ktpId = newKTP.ID
	}

	var pendidikanId uint

	if inputEmployee.Pendidikan.ID > 0 {
		// Update existing
		if err := r.db.Model(&pendidikan.Pendidikan{}).
			Where("id = ?", inputEmployee.Pendidikan.ID).
			Updates(inputEmployee.Pendidikan).Error; err != nil {
			return updatedEmployee, err
		}
		pendidikanId = inputEmployee.Pendidikan.ID
	} else {
		// Create new
		var newPendidikan pendidikan.Pendidikan = inputEmployee.Pendidikan

		if newPendidikan.PendidikanTerakhir != "" {
			if err := r.db.Create(&newPendidikan).Error; err != nil {
				return Employee{}, err
			}
		}
		pendidikanId = newPendidikan.ID
	}

	var laporanId uint

	if inputEmployee.Laporan.ID > 0 {
		// Update existing
		if err := r.db.Model(&laporan.Laporan{}).
			Where("id = ?", inputEmployee.Laporan.ID).
			Updates(inputEmployee.Laporan).Error; err != nil {
			return updatedEmployee, err
		}
		laporanId = inputEmployee.Laporan.ID
	} else {
		// Create new
		var newLaporan laporan.Laporan = inputEmployee.Laporan

		if newLaporan.RingSerapan != "" {
			if err := r.db.Create(&newLaporan).Error; err != nil {
				return Employee{}, err
			}
		}
		laporanId = newLaporan.ID
	}

	var apdId uint

	if inputEmployee.APD != nil {
		if inputEmployee.APD.ID > 0 {
			if err := r.db.Model(&apd.APD{}).
				Where("id = ?", inputEmployee.APD.ID).
				Updates(inputEmployee.APD).Error; err != nil {
				return updatedEmployee, err
			}
			apdId = inputEmployee.APD.ID

		} else if inputEmployee.APD.UkuranBaju != "" {
			newAPD := *inputEmployee.APD
			if err := r.db.Create(&newAPD).Error; err != nil {
				return Employee{}, err
			}
			apdId = newAPD.ID
		}
	}

	var npwpId uint

	if inputEmployee.NPWP != nil {
		if inputEmployee.NPWP.ID > 0 {
			if err := r.db.Model(&npwp.NPWP{}).
				Where("id = ?", inputEmployee.NPWP.ID).
				Updates(inputEmployee.NPWP).Error; err != nil {
				return updatedEmployee, err
			}
			npwpId = inputEmployee.NPWP.ID

		} else if *inputEmployee.NPWP.NomorNPWP != "" {
			newNPWP := *inputEmployee.NPWP
			if err := r.db.Create(&newNPWP).Error; err != nil {
				return Employee{}, err
			}
			npwpId = newNPWP.ID
		}
	}

	var bankId uint

	if inputEmployee.Bank.ID > 0 {
		// Update existing
		if err := r.db.Model(&bank.Bank{}).
			Where("id = ?", inputEmployee.Bank.ID).
			Updates(inputEmployee.Bank).Error; err != nil {
			return updatedEmployee, err
		}
		bankId = inputEmployee.Bank.ID
	} else {
		// Create new
		var newBank bank.Bank = inputEmployee.Bank

		if newBank.NamaBank != "" {
			if err := r.db.Create(&newBank).Error; err != nil {
				return Employee{}, err
			}
		}
		bankId = newBank.ID
	}

	var bpjskesehatanId uint

	if inputEmployee.BPJSKesehatan != nil {
		if inputEmployee.BPJSKesehatan.ID > 0 {
			// Update existing
			if err := r.db.Model(&bpjskesehatan.BPJSKesehatan{}).
				Where("id = ?", inputEmployee.BPJSKesehatan.ID).
				Updates(inputEmployee.BPJSKesehatan).Error; err != nil {
				return updatedEmployee, err
			}
			bpjskesehatanId = inputEmployee.BPJSKesehatan.ID

		} else if inputEmployee.BPJSKesehatan.NomorKesehatan != "" {
			// Create new
			newBPJSKesehatan := *inputEmployee.BPJSKesehatan
			if err := r.db.Create(&newBPJSKesehatan).Error; err != nil {
				return Employee{}, err
			}
			bpjskesehatanId = newBPJSKesehatan.ID
		}
	}

	var bpjsketenagakerjaanId uint

	if inputEmployee.BPJSKetenagakerjaan != nil {
		if inputEmployee.BPJSKetenagakerjaan.ID > 0 {
			// Update existing
			if err := r.db.Model(&bpjsketenagakerjaan.BPJSKetenagakerjaan{}).
				Where("id = ?", inputEmployee.BPJSKetenagakerjaan.ID).
				Updates(inputEmployee.BPJSKetenagakerjaan).Error; err != nil {
				return updatedEmployee, err
			}
			bpjsketenagakerjaanId = inputEmployee.BPJSKetenagakerjaan.ID

		} else if inputEmployee.BPJSKetenagakerjaan.NomorKetenagakerjaan != "" {
			// Create new
			newBPJSKetenagakerjaan := *inputEmployee.BPJSKetenagakerjaan
			if err := r.db.Create(&newBPJSKetenagakerjaan).Error; err != nil {
				return Employee{}, err
			}
			bpjsketenagakerjaanId = newBPJSKetenagakerjaan.ID
		}
	}

	// 1. Find existing employee
	if err := r.db.Where("id = ?", id).First(&updatedEmployee).Error; err != nil {
		return updatedEmployee, err
	}

	// 2. Update main fields
	updatedEmployee.NomorKaryawan = inputEmployee.NomorKaryawan
	updatedEmployee.DepartmentId = inputEmployee.DepartmentId
	updatedEmployee.Firstname = inputEmployee.Firstname
	updatedEmployee.Lastname = inputEmployee.Lastname
	updatedEmployee.PhoneNumber = inputEmployee.PhoneNumber
	updatedEmployee.Email = inputEmployee.Email
	updatedEmployee.Level = inputEmployee.Level
	updatedEmployee.Status = inputEmployee.Status
	updatedEmployee.RoleId = inputEmployee.RoleId
	updatedEmployee.HiredBy = inputEmployee.HiredBy
	updatedEmployee.KartuKeluargaId = kartuKeluargaId
	updatedEmployee.KTPId = ktpId
	updatedEmployee.PendidikanId = pendidikanId
	updatedEmployee.LaporanId = laporanId
	updatedEmployee.APDId = UintToPtr(apdId)
	updatedEmployee.NPWPId = UintToPtr(npwpId)
	updatedEmployee.BankId = bankId
	updatedEmployee.BPJSKesehatanId = UintToPtr(bpjskesehatanId)
	updatedEmployee.BPJSKetenagakerjaanId = UintToPtr(bpjsketenagakerjaanId)
	updatedEmployee.PositionId = inputEmployee.PositionId
	updatedEmployee.DateOfHire = inputEmployee.DateOfHire

	if err := r.db.Save(&updatedEmployee).Error; err != nil {
		return updatedEmployee, err
	}

	// 4. Return
	return updatedEmployee, nil
}

func (r *repository) DeleteEmployee(id uint) (bool, error) {
	tx := r.db.Begin()
	var employees Employee

	// Check existence (this automatically ignores soft-deleted entries)
	if err := tx.Where("id = ?", id).First(&employees).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// ✅ Soft delete (do NOT use Unscoped)
	if err := tx.Delete(&employees).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	tx.Commit()
	return true, nil
}

func (r *repository) ListDashboard(empCode uint, dashboardSort SortFilterDashboardEmployee) (DashboardEmployee, error) {
	var dashboardEmployees DashboardEmployee
	var totalEmployeeCount, maleCount, femaleCount, totalHOCount, totalSiteCount int64

	queryFilter := "employees.id > 0"
	ptConditions := []string{}

	// === Handle PT Logic ===
	ptList := []string{}
	if strings.TrimSpace(dashboardSort.PT) != "" {
		ptList = strings.Split(dashboardSort.PT, ",")
	} else {
		// Default to PTs based on empCode
		if empCode == 1 {
			ptList = []string{"PT. MRP", "PT. TRIOP"}
		} else if empCode == 2 {
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

	// Exclude PTs depending on empCode
	if empCode == 1 {
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%MJSU%'"
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%IBS%'"
		queryFilter += " AND cast(status AS TEXT) ILIKE 'AKTIF'"
	} else if empCode == 2 {
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%MRP%'"
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%TRIOP%'"
		queryFilter += " AND cast(status AS TEXT) ILIKE 'AKTIF'"
	}

	// === Handle Department Filter (Optional) ===
	if strings.TrimSpace(dashboardSort.DepartmentId) != "" {
		deptIDs := strings.Split(dashboardSort.DepartmentId, ",")
		cleaned := []string{}
		for _, id := range deptIDs {
			id = strings.TrimSpace(id)
			if id != "" {
				cleaned = append(cleaned, id)
			}
		}
		if len(cleaned) > 0 {
			queryFilter += fmt.Sprintf(" AND department_id IN (%s)", strings.Join(cleaned, ","))
		}
	}

	// === Queries ===
	r.db.Model(&Employee{}).Where(queryFilter).Count(&totalEmployeeCount)
	dashboardEmployees.TotalEmployee = uint(totalEmployeeCount)

	queryMale := "cast(ktps.gender AS TEXT) ILIKE 'MALE' AND " + queryFilter
	queryFemale := "cast(ktps.gender AS TEXT) ILIKE 'FEMALE' AND " + queryFilter

	r.db.Model(&Employee{}).
		Joins("JOIN ktps ON employees.ktp_id = ktps.id").
		Where(queryMale).Count(&maleCount)
	dashboardEmployees.TotalMale = uint(maleCount)

	r.db.Model(&Employee{}).
		Joins("JOIN ktps ON employees.ktp_id = ktps.id").
		Where(queryFemale).Count(&femaleCount)
	dashboardEmployees.TotalFemale = uint(femaleCount)

	querySite := "cast(employees.hired_by AS TEXT) ILIKE '%SITE%' AND " + queryFilter
	queryHO := "cast(employees.hired_by AS TEXT) ILIKE '%HEAD OFFICE%' AND " + queryFilter

	r.db.Model(&Employee{}).Where(querySite).Count(&totalSiteCount)
	dashboardEmployees.HireSite = uint(totalSiteCount)

	r.db.Model(&Employee{}).Where(queryHO).Count(&totalHOCount)
	dashboardEmployees.HireHO = uint(totalHOCount)

	// BASED ON AGE
	var ageDist BasedOnAge
	err := r.db.Model(&Employee{}).
		Joins("JOIN ktps ON ktps.id = employees.ktp_id").
		Select(`
			COUNT(CASE WHEN DATE_PART('year', AGE(CURRENT_DATE, ktps.tanggal_lahir::DATE)) BETWEEN 18 AND 27 THEN 1 END) AS stage1,
			COUNT(CASE WHEN DATE_PART('year', AGE(CURRENT_DATE, ktps.tanggal_lahir::DATE)) BETWEEN 28 AND 37 THEN 1 END) AS stage2,
			COUNT(CASE WHEN DATE_PART('year', AGE(CURRENT_DATE, ktps.tanggal_lahir::DATE)) BETWEEN 38 AND 47 THEN 1 END) AS stage3,
			COUNT(CASE WHEN DATE_PART('year', AGE(CURRENT_DATE, ktps.tanggal_lahir::DATE)) BETWEEN 48 AND 57 THEN 1 END) AS stage4,
			COUNT(CASE WHEN DATE_PART('year', AGE(CURRENT_DATE, ktps.tanggal_lahir::DATE)) BETWEEN 58 AND 67 THEN 1 END) AS stage5
		`).
		Where(queryFilter). // ✅ use your dynamic filter
		Scan(&ageDist).Error

	if err != nil {
		return dashboardEmployees, err
	}
	dashboardEmployees.BasedOnAge = ageDist

	// BASED ON Year Service
	var yearService BasedOnYear
	errYear := r.db.Model(&Employee{}).
		Select(`
		COUNT(CASE 
			WHEN AGE(CURRENT_DATE, date_of_hire::DATE) < INTERVAL '6 months' 
			THEN 1 END) AS year1,
		COUNT(CASE 
			WHEN AGE(CURRENT_DATE, date_of_hire::DATE) >= INTERVAL '6 months' AND AGE(CURRENT_DATE, date_of_hire::DATE) < INTERVAL '1 year' 
			THEN 1 END) AS year2,
		COUNT(CASE 
			WHEN AGE(CURRENT_DATE, date_of_hire::DATE) >= INTERVAL '1 year' AND AGE(CURRENT_DATE, date_of_hire::DATE) < INTERVAL '2 years' 
			THEN 1 END) AS year3,
		COUNT(CASE 
			WHEN AGE(CURRENT_DATE, date_of_hire::DATE) >= INTERVAL '2 years' 
			THEN 1 END) AS year4
	`).
		Where(queryFilter).
		Scan(&yearService).Error

	if errYear != nil {
		return dashboardEmployees, errYear
	}
	dashboardEmployees.BasedOnYear = yearService

	// BASED ON EDUCATION
	var education BasedOnEducation

	errEducation := r.db.Model(&Employee{}).
		Joins("JOIN pendidikans ON employees.pendidikan_id = pendidikans.id").
		Select(`
			COUNT(CASE WHEN LOWER(pendidikans.pendidikan_label) LIKE 'sd%' THEN 1 END) AS Edu1,
			COUNT(CASE WHEN LOWER(pendidikans.pendidikan_label) LIKE 'smp%' THEN 1 END) AS Edu2,
			COUNT(CASE WHEN LOWER(pendidikans.pendidikan_label) LIKE 'sma%' THEN 1 END) AS Edu3,
			COUNT(CASE WHEN LOWER(pendidikans.pendidikan_label) LIKE 'diploma%' THEN 1 END) AS Edu4,
			COUNT(CASE WHEN LOWER(pendidikans.pendidikan_label) LIKE 'sarjana%' THEN 1 END) AS Edu5
		`).
		Where(queryFilter).
		Scan(&education).Error

	if errEducation != nil {
		return dashboardEmployees, errEducation
	}
	dashboardEmployees.BasedOnEducation = education

	// BASED ON DEPARTMENT
	var department BasedOnDepartment

	errDepartment := r.db.Model(&Employee{}).
		Joins("JOIN departments ON employees.department_id = departments.id").
		Select(`
			COUNT(CASE WHEN LOWER(departments.department_name) LIKE 'engineering' THEN 1 END) AS Engineering,
			COUNT(CASE WHEN LOWER(departments.department_name) LIKE 'finance' THEN 1 END) AS Finance,
			COUNT(CASE WHEN LOWER(departments.department_name) LIKE 'hrga' THEN 1 END) AS HRGA,
			COUNT(CASE WHEN LOWER(departments.department_name) LIKE 'operation' THEN 1 END) AS Operation,
			COUNT(CASE WHEN LOWER(departments.department_name) LIKE 'plant' THEN 1 END) AS Plant,
			COUNT(CASE WHEN LOWER(departments.department_name) LIKE 'she' THEN 1 END) AS SHE,
			COUNT(CASE WHEN LOWER(departments.department_name) LIKE 'coal loading' THEN 1 END) AS coal_loading,
			COUNT(CASE WHEN LOWER(departments.department_name) LIKE 'stockpile' THEN 1 END) AS Stockpile,
			COUNT(CASE WHEN LOWER(departments.department_name) LIKE 'shipping' THEN 1 END) AS Shipping,
			COUNT(CASE WHEN LOWER(departments.department_name) LIKE 'plant & logistic' THEN 1 END) AS plant_logistic,
			COUNT(CASE WHEN LOWER(departments.department_name) LIKE 'keamanan & eksternal' THEN 1 END) AS keamanan_eksternal,
			COUNT(CASE WHEN LOWER(departments.department_name) LIKE 'oshe (operation & she)' THEN 1 END) AS Oshe,
			COUNT(CASE WHEN LOWER(departments.department_name) LIKE 'management' THEN 1 END) AS Management
		`).
		Where(queryFilter).
		Scan(&department).Error

	if errDepartment != nil {
		return dashboardEmployees, errDepartment
	}
	dashboardEmployees.BasedOnDepartment = department

	// BASED ON RING
	var ring BasedOnRing

	errRing := r.db.Model(&Employee{}).
		Joins("JOIN laporans ON employees.laporan_id = laporans.id").
		Select(` 
			COUNT(CASE WHEN LOWER(laporans.ring_r_ip_pm) LIKE 'ring i' THEN 1 END) AS Ring1,
			COUNT(CASE WHEN LOWER(laporans.ring_r_ip_pm) LIKE 'ring ii' THEN 1 END) AS Ring2,
			COUNT(CASE WHEN LOWER(laporans.ring_r_ip_pm) LIKE 'ring iii' THEN 1 END) AS Ring3,
			COUNT(CASE WHEN LOWER(laporans.ring_r_ip_pm) LIKE 'luar ring' THEN 1 END) AS luar_ring
		`).
		Where(queryFilter).
		Scan(&ring).Error

	if errRing != nil {
		return dashboardEmployees, errRing
	}
	dashboardEmployees.BasedOnRing = ring

	// BASED ON LOKAL NON LOKAL
	var lokalNonLokal BasedOnLokal

	errLokal := r.db.Model(&Employee{}).
		Joins("JOIN laporans ON employees.laporan_id = laporans.id").
		Select(`
			COUNT(CASE WHEN LOWER(laporans.kategori_lokal_non_lokal) LIKE 'lokal' THEN 1 END) AS lokal,
			COUNT(CASE WHEN LOWER(laporans.kategori_lokal_non_lokal) LIKE 'non lokal' THEN 1 END) AS non_lokal
		`).
		Where(queryFilter).
		Scan(&lokalNonLokal).Error

	if errLokal != nil {
		return dashboardEmployees, errLokal
	}
	dashboardEmployees.BasedOnLokal = lokalNonLokal

	return dashboardEmployees, nil
}

func (r *repository) ListDashboardTurnover(empCode uint, dashboardSort SortFilterDashboardEmployeeTurnOver) (DashboardEmployeeTurnOver, error) {
	var dashboardEmployees DashboardEmployeeTurnOver
	var totalNewHireCount, TotalResignCount, TotalBerakhirKontrakCount, TotalPHKCount int64

	queryFilter := "employees.id > 0"
	ptConditions := []string{}

	// === Handle PT Logic ===
	ptList := []string{}
	if strings.TrimSpace(dashboardSort.PT) != "" {
		ptList = strings.Split(dashboardSort.PT, ",")
	} else {
		// Default to PTs based on empCode
		if empCode == 1 {
			ptList = []string{"PT. MRP", "PT. TRIOP"}
		} else if empCode == 2 {
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

	// Exclude PTs depending on empCode
	if empCode == 1 {
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%MJSU%'"
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%IBS%'"
	} else if empCode == 2 {
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%MRP%'"
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%TRIOP%'"
	}

	if strings.TrimSpace(dashboardSort.Year) == "" {
		// Set to current year
		dashboardSort.Year = fmt.Sprintf("%d", time.Now().Year())
	}

	queryNewHire := "date_of_hire IS NOT NULL AND cast(status AS TEXT) ILIKE 'AKTIF' AND EXTRACT(YEAR FROM date_of_hire::DATE) = " + dashboardSort.Year + " AND " + queryFilter

	r.db.Model(&Employee{}).Where(queryNewHire).Count(&totalNewHireCount)
	dashboardEmployees.TotalHire = uint(totalNewHireCount)

	queryResign := "EXTRACT(YEAR FROM histories.tanggal::DATE) = " + dashboardSort.Year + "  AND cast(status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'resign' AND " + queryFilter

	r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Where(queryResign).Count(&TotalResignCount)
	dashboardEmployees.TotalResign = uint(TotalResignCount)

	queryBerakhirPKWT := "EXTRACT(YEAR FROM histories.tanggal::DATE) = " + dashboardSort.Year + " AND cast(status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'berakhir kontrak' AND " + queryFilter

	r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Where(queryBerakhirPKWT).Count(&TotalBerakhirKontrakCount)
	dashboardEmployees.TotalBerakhirPkwt = uint(TotalBerakhirKontrakCount)

	queryPHK := "EXTRACT(YEAR FROM histories.tanggal::DATE) = " + dashboardSort.Year + " AND cast(status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'phk' AND " + queryFilter

	r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Where(queryPHK).Count(&TotalPHKCount)
	dashboardEmployees.TotalPhk = uint(TotalPHKCount)

	// BASED ON MONTH

	queryFilterNewHire := "EXTRACT(YEAR FROM date_of_hire::DATE) = " + dashboardSort.Year + " AND " + queryFilter
	queryFilterAll := "EXTRACT(YEAR FROM histories.tanggal::DATE) = " + dashboardSort.Year + " AND " + queryFilter
	var januari DataStatus
	var newHireCountJanuari int64

	errNewHire := r.db.Model(&Employee{}).
		Select("COUNT(*)").
		Where(`
			date_of_hire IS NOT NULL AND 
			cast(status AS TEXT) ILIKE 'AKTIF' AND 
			EXTRACT(MONTH FROM date_of_hire::DATE) = 1`).
		Where(queryFilterNewHire).
		Count(&newHireCountJanuari).Error

	if errNewHire != nil {
		return dashboardEmployees, errNewHire
	}

	errJanuari := r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'berakhir kontrak' THEN 1 END) AS berakhir_pkwt,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'resign' THEN 1 END) AS Resign,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'phk' THEN 1 END) AS PHK
		`).
		Where(queryFilterAll).
		Scan(&januari).Error

	if errJanuari != nil {
		return dashboardEmployees, errJanuari
	}
	dashboardEmployees.Januari = DataStatus{
		NewHire:      uint(newHireCountJanuari),
		BerakhirPkwt: januari.BerakhirPkwt,
		Resign:       januari.Resign,
		PHK:          januari.PHK,
	}

	var februari DataStatus
	var newHireCountFebruari int64

	errNewHire = r.db.Model(&Employee{}).
		Select("COUNT(*)").
		Where(`
				date_of_hire IS NOT NULL AND 
				cast(status AS TEXT) ILIKE 'AKTIF' AND 
				EXTRACT(MONTH FROM date_of_hire::DATE) = 2`).
		Where(queryFilterNewHire).
		Count(&newHireCountFebruari).Error

	if errNewHire != nil {
		return dashboardEmployees, errNewHire
	}

	errFebruari := r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'berakhir kontrak' THEN 1 END) AS berakhir_pkwt,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'resign' THEN 1 END) AS Resign,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'phk' THEN 1 END) AS PHK
		`).
		Where(queryFilterAll).
		Scan(&februari).Error

	if errFebruari != nil {
		return dashboardEmployees, errFebruari
	}
	dashboardEmployees.Februari = DataStatus{
		NewHire:      uint(newHireCountFebruari),
		BerakhirPkwt: februari.BerakhirPkwt,
		Resign:       februari.Resign,
		PHK:          februari.PHK,
	}

	var maret DataStatus
	var newHireCountMaret int64

	errNewHire = r.db.Model(&Employee{}).
		Select("COUNT(*)").
		Where(`
				date_of_hire IS NOT NULL AND 
				cast(status AS TEXT) ILIKE 'AKTIF' AND 
				EXTRACT(MONTH FROM date_of_hire::DATE) = 3`).
		Where(queryFilterNewHire).
		Count(&newHireCountMaret).Error

	if errNewHire != nil {
		return dashboardEmployees, errNewHire
	}

	errMaret := r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'berakhir kontrak' THEN 1 END) AS berakhir_pkwt,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'resign' THEN 1 END) AS Resign,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'phk' THEN 1 END) AS PHK
		`).
		Where(queryFilterAll).
		Scan(&maret).Error

	if errMaret != nil {
		return dashboardEmployees, errMaret
	}
	dashboardEmployees.Maret = DataStatus{
		NewHire:      uint(newHireCountMaret),
		BerakhirPkwt: maret.BerakhirPkwt,
		Resign:       maret.Resign,
		PHK:          maret.PHK,
	}

	var april DataStatus
	var newHireCountApril int64

	errNewHire = r.db.Model(&Employee{}).
		Select("COUNT(*)").
		Where(`
				date_of_hire IS NOT NULL AND 
				cast(status AS TEXT) ILIKE 'AKTIF' AND 
				EXTRACT(MONTH FROM date_of_hire::DATE) = 4`).
		Where(queryFilterNewHire).
		Count(&newHireCountApril).Error

	if errNewHire != nil {
		return dashboardEmployees, errNewHire
	}
	errApril := r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'berakhir kontrak' THEN 1 END) AS berakhir_pkwt,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'resign' THEN 1 END) AS Resign,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'phk' THEN 1 END) AS PHK
		`).
		Where(queryFilterAll).
		Scan(&april).Error

	if errApril != nil {
		return dashboardEmployees, errApril
	}
	dashboardEmployees.April = DataStatus{
		NewHire:      uint(newHireCountApril),
		BerakhirPkwt: april.BerakhirPkwt,
		Resign:       april.Resign,
		PHK:          april.PHK,
	}

	var mei DataStatus
	var newHireCountMei int64

	errNewHire = r.db.Model(&Employee{}).
		Select("COUNT(*)").
		Where(`
				date_of_hire IS NOT NULL AND 
				cast(status AS TEXT) ILIKE 'AKTIF' AND 
				EXTRACT(MONTH FROM date_of_hire::DATE) = 5`).
		Where(queryFilterNewHire).
		Count(&newHireCountMei).Error

	if errNewHire != nil {
		return dashboardEmployees, errNewHire
	}

	errMei := r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'berakhir kontrak' THEN 1 END) AS berakhir_pkwt,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'resign' THEN 1 END) AS Resign,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'phk' THEN 1 END) AS PHK
		`).
		Where(queryFilterAll).
		Scan(&mei).Error

	if errMei != nil {
		return dashboardEmployees, errMei
	}
	dashboardEmployees.Mei = DataStatus{
		NewHire:      uint(newHireCountMei),
		BerakhirPkwt: mei.BerakhirPkwt,
		Resign:       mei.Resign,
		PHK:          mei.PHK,
	}

	var juni DataStatus
	var newHireCountJuni int64

	errNewHire = r.db.Model(&Employee{}).
		Select("COUNT(*)").
		Where(`
				date_of_hire IS NOT NULL AND 
				cast(status AS TEXT) ILIKE 'AKTIF' AND 
				EXTRACT(MONTH FROM date_of_hire::DATE) = 6`).
		Where(queryFilterNewHire).
		Count(&newHireCountJuni).Error

	if errNewHire != nil {
		return dashboardEmployees, errNewHire
	}

	errJuni := r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'berakhir kontrak' THEN 1 END) AS berakhir_pkwt,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'resign' THEN 1 END) AS Resign,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'phk' THEN 1 END) AS PHK
		`).
		Where(queryFilterAll).
		Scan(&juni).Error

	if errJuni != nil {
		return dashboardEmployees, errJuni
	}
	dashboardEmployees.Juni = DataStatus{
		NewHire:      uint(newHireCountJuni),
		BerakhirPkwt: juni.BerakhirPkwt,
		Resign:       juni.Resign,
		PHK:          juni.PHK,
	}

	var juli DataStatus
	var newHireCountJuli int64

	errNewHire = r.db.Model(&Employee{}).
		Select("COUNT(*)").
		Where(`
				date_of_hire IS NOT NULL AND 
				cast(status AS TEXT) ILIKE 'AKTIF' AND 
				EXTRACT(MONTH FROM date_of_hire::DATE) = 7`).
		Where(queryFilterNewHire).
		Count(&newHireCountJuli).Error

	if errNewHire != nil {
		return dashboardEmployees, errNewHire
	}

	errJuli := r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'berakhir kontrak' THEN 1 END) AS berakhir_pkwt,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'resign' THEN 1 END) AS Resign,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'phk' THEN 1 END) AS PHK
		`).
		Where(queryFilterAll).
		Scan(&juli).Error

	if errJuli != nil {
		return dashboardEmployees, errJuli
	}
	dashboardEmployees.Juli = DataStatus{
		NewHire:      uint(newHireCountJuli),
		BerakhirPkwt: juli.BerakhirPkwt,
		Resign:       juli.Resign,
		PHK:          juli.PHK,
	}

	var agustus DataStatus
	var newHireCountAgustus int64

	errNewHire = r.db.Model(&Employee{}).
		Select("COUNT(*)").
		Where(`
				date_of_hire IS NOT NULL AND 
				cast(status AS TEXT) ILIKE 'AKTIF' AND 
				EXTRACT(MONTH FROM date_of_hire::DATE) = 8`).
		Where(queryFilterNewHire).
		Count(&newHireCountAgustus).Error

	if errNewHire != nil {
		return dashboardEmployees, errNewHire
	}

	errAgustus := r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'berakhir kontrak' THEN 1 END) AS berakhir_pkwt,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'resign' THEN 1 END) AS Resign,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'phk' THEN 1 END) AS PHK
		`).
		Where(queryFilterAll).
		Scan(&agustus).Error

	if errAgustus != nil {
		return dashboardEmployees, errAgustus
	}
	dashboardEmployees.Agustus = DataStatus{
		NewHire:      uint(newHireCountAgustus),
		BerakhirPkwt: agustus.BerakhirPkwt,
		Resign:       agustus.Resign,
		PHK:          agustus.PHK,
	}

	var september DataStatus
	var newHireCountSeptember int64

	errNewHire = r.db.Model(&Employee{}).
		Select("COUNT(*)").
		Where(`
				date_of_hire IS NOT NULL AND 
				cast(status AS TEXT) ILIKE 'AKTIF' AND 
				EXTRACT(MONTH FROM date_of_hire::DATE) = 9`).
		Where(queryFilterNewHire).
		Count(&newHireCountSeptember).Error

	if errNewHire != nil {
		return dashboardEmployees, errNewHire
	}

	errSeptember := r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'berakhir kontrak' THEN 1 END) AS berakhir_pkwt,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'resign' THEN 1 END) AS Resign,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'phk' THEN 1 END) AS PHK
		`).
		Where(queryFilterAll).
		Scan(&september).Error

	if errSeptember != nil {
		return dashboardEmployees, errSeptember
	}
	dashboardEmployees.September = DataStatus{
		NewHire:      uint(newHireCountSeptember),
		BerakhirPkwt: september.BerakhirPkwt,
		Resign:       september.Resign,
		PHK:          september.PHK,
	}

	var oktober DataStatus
	var newHireCountOktober int64

	errNewHire = r.db.Model(&Employee{}).
		Select("COUNT(*)").
		Where(`
				date_of_hire IS NOT NULL AND 
				cast(status AS TEXT) ILIKE 'AKTIF' AND 
				EXTRACT(MONTH FROM date_of_hire::DATE) = 10`).
		Where(queryFilterNewHire).
		Count(&newHireCountOktober).Error

	if errNewHire != nil {
		return dashboardEmployees, errNewHire
	}

	errOktober := r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'berakhir kontrak' THEN 1 END) AS berakhir_pkwt,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'resign' THEN 1 END) AS Resign,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'phk' THEN 1 END) AS PHK
		`).
		Where(queryFilterAll).
		Scan(&oktober).Error

	if errOktober != nil {
		return dashboardEmployees, errOktober
	}
	dashboardEmployees.Oktober = DataStatus{
		NewHire:      uint(newHireCountOktober),
		BerakhirPkwt: oktober.BerakhirPkwt,
		Resign:       oktober.Resign,
		PHK:          oktober.PHK,
	}

	var november DataStatus
	var newHireCountNovember int64

	errNewHire = r.db.Model(&Employee{}).
		Select("COUNT(*)").
		Where(`
				date_of_hire IS NOT NULL AND 
				cast(status AS TEXT) ILIKE 'AKTIF' AND 
				EXTRACT(MONTH FROM date_of_hire::DATE) = 11`).
		Where(queryFilterNewHire).
		Count(&newHireCountNovember).Error

	if errNewHire != nil {
		return dashboardEmployees, errNewHire
	}

	errNovember := r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'berakhir kontrak' THEN 1 END) AS berakhir_pkwt,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'resign' THEN 1 END) AS Resign,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'phk' THEN 1 END) AS PHK
		`).
		Where(queryFilterAll).
		Scan(&november).Error

	if errNovember != nil {
		return dashboardEmployees, errNovember
	}
	dashboardEmployees.November = DataStatus{
		NewHire:      uint(newHireCountNovember),
		BerakhirPkwt: november.BerakhirPkwt,
		Resign:       november.Resign,
		PHK:          november.PHK,
	}

	var desember DataStatus
	var newHireCountDesember int64

	errNewHire = r.db.Model(&Employee{}).
		Select("COUNT(*)").
		Where(`
				date_of_hire IS NOT NULL AND 
				cast(status AS TEXT) ILIKE 'AKTIF' AND 
				EXTRACT(MONTH FROM date_of_hire::DATE) = 12`).
		Where(queryFilterNewHire).
		Count(&newHireCountDesember).Error

	if errNewHire != nil {
		return dashboardEmployees, errNewHire
	}

	errDesember := r.db.Model(&Employee{}).
		Joins("JOIN histories ON employees.id = histories.employee_id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 12  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'berakhir kontrak' THEN 1 END) AS berakhir_pkwt,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 12  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'resign' THEN 1 END) AS Resign,
			COUNT(CASE WHEN EXTRACT(MONTH FROM histories.tanggal::DATE) = 12  AND cast(employees.status AS TEXT) ILIKE 'TIDAK AKTIF' AND LOWER(histories.status_terakhir) LIKE 'phk' THEN 1 END) AS PHK
		`).
		Where(queryFilterAll).
		Scan(&desember).Error

	if errDesember != nil {
		return dashboardEmployees, errDesember
	}
	dashboardEmployees.Desember = DataStatus{
		NewHire:      uint(newHireCountDesember),
		BerakhirPkwt: desember.BerakhirPkwt,
		Resign:       desember.Resign,
		PHK:          desember.PHK,
	}

	return dashboardEmployees, nil
}

func (r *repository) ListDashboardKontrak(empCode uint, dashboardSort SortFilterDashboardEmployeeTurnOver) (DashboardEmployeeKontrak, error) {
	var dashboardEmployees DashboardEmployeeKontrak

	queryFilter := "employees.id > 0"
	ptConditions := []string{}

	// === Handle PT Logic ===
	ptList := []string{}
	if strings.TrimSpace(dashboardSort.PT) != "" {
		ptList = strings.Split(dashboardSort.PT, ",")
	} else {
		// Default to PTs based on empCode
		if empCode == 1 {
			ptList = []string{"PT. MRP", "PT. TRIOP"}
		} else if empCode == 2 {
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

	// Exclude PTs depending on empCode
	if empCode == 1 {
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%MJSU%'"
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%IBS%'"
	} else if empCode == 2 {
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%MRP%'"
		queryFilter += " AND cast(nomor_karyawan AS TEXT) NOT ILIKE '%TRIOP%'"
	}

	if strings.TrimSpace(dashboardSort.Year) == "" {
		// Set to current year
		dashboardSort.Year = fmt.Sprintf("%d", time.Now().Year())
	}

	// BASED ON MONTH
	queryFilterAll := "EXTRACT(YEAR FROM dohs.tanggal_end_doh::DATE) = " + dashboardSort.Year + " AND " + queryFilter

	var januari DepartmentName
	errJanuari := r.db.Model(&Employee{}).
		Joins("JOIN dohs ON employees.id = dohs.employee_id").
		Joins("JOIN departments ON employees.department_id = departments.id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'operation' THEN 1 END) AS Operation, 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant' THEN 1 END) AS Plant,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'hrga' THEN 1 END) AS Hrga,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'she' THEN 1 END) AS She,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'finance' THEN 1 END) AS Finance,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'engineering' THEN 1 END) AS Engineering,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'coal loading' THEN 1 END) AS CoalLoading,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'stockpile' THEN 1 END) AS Stockpile,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'shippping' THEN 1 END) AS Shipping,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant & logistic' THEN 1 END) AS PlantLogistic,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'keamanan & eksternal' THEN 1 END) AS keamanan_eksternal,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'oshe (operation & she)' THEN 1 END) AS Oshe,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 1  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'management' THEN 1 END) AS Management
		`).
		Where(queryFilterAll).
		Scan(&januari).Error

	if errJanuari != nil {
		return dashboardEmployees, errJanuari
	}

	dashboardEmployees.Januari = januari

	var februari DepartmentName

	errFebruari := r.db.Model(&Employee{}).
		Joins("JOIN dohs ON employees.id = dohs.employee_id").
		Joins("JOIN departments ON employees.department_id = departments.id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'operation' THEN 1 END) AS Operation, 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant' THEN 1 END) AS Plant,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'hrga' THEN 1 END) AS Hrga,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'she' THEN 1 END) AS She,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'finance' THEN 1 END) AS Finance,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'engineering' THEN 1 END) AS Engineering,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'coal loading' THEN 1 END) AS CoalLoading,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'stockpile' THEN 1 END) AS Stockpile,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'shippping' THEN 1 END) AS Shipping,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant & logistic' THEN 1 END) AS PlantLogistic,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'keamanan & eksternal' THEN 1 END) AS keamanan_eksternal,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'oshe (operation & she)' THEN 1 END) AS Oshe,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 2  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'management' THEN 1 END) AS Management
		`).
		Where(queryFilterAll).
		Scan(&februari).Error

	if errFebruari != nil {
		return dashboardEmployees, errFebruari
	}

	dashboardEmployees.Februari = februari

	var maret DepartmentName

	errMaret := r.db.Model(&Employee{}).
		Joins("JOIN dohs ON employees.id = dohs.employee_id").
		Joins("JOIN departments ON employees.department_id = departments.id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'operation' THEN 1 END) AS Operation, 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant' THEN 1 END) AS Plant,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'hrga' THEN 1 END) AS Hrga,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'she' THEN 1 END) AS She,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'finance' THEN 1 END) AS Finance,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'engineering' THEN 1 END) AS Engineering,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'coal loading' THEN 1 END) AS CoalLoading,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'stockpile' THEN 1 END) AS Stockpile,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'shippping' THEN 1 END) AS Shipping,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant & logistic' THEN 1 END) AS PlantLogistic,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'keamanan & eksternal' THEN 1 END) AS keamanan_eksternal,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'oshe (operation & she)' THEN 1 END) AS Oshe,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 3  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'management' THEN 1 END) AS Management
		`).
		Where(queryFilterAll).
		Scan(&maret).Error

	if errMaret != nil {
		return dashboardEmployees, errMaret
	}

	dashboardEmployees.Maret = maret

	var april DepartmentName

	errApril := r.db.Model(&Employee{}).
		Joins("JOIN dohs ON employees.id = dohs.employee_id").
		Joins("JOIN departments ON employees.department_id = departments.id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'operation' THEN 1 END) AS Operation, 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant' THEN 1 END) AS Plant,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'hrga' THEN 1 END) AS Hrga,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'she' THEN 1 END) AS She,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'finance' THEN 1 END) AS Finance,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'engineering' THEN 1 END) AS Engineering,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'coal loading' THEN 1 END) AS CoalLoading,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'stockpile' THEN 1 END) AS Stockpile,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'shippping' THEN 1 END) AS Shipping,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant & logistic' THEN 1 END) AS PlantLogistic,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'keamanan & eksternal' THEN 1 END) AS keamanan_eksternal,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'oshe (operation & she)' THEN 1 END) AS Oshe,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 4  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'management' THEN 1 END) AS Management
		`).
		Where(queryFilterAll).
		Scan(&april).Error

	if errApril != nil {
		return dashboardEmployees, errApril
	}
	dashboardEmployees.April = april

	var mei DepartmentName

	errMei := r.db.Model(&Employee{}).
		Joins("JOIN dohs ON employees.id = dohs.employee_id").
		Joins("JOIN departments ON employees.department_id = departments.id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'operation' THEN 1 END) AS Operation, 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant' THEN 1 END) AS Plant,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'hrga' THEN 1 END) AS Hrga,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'she' THEN 1 END) AS She,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'finance' THEN 1 END) AS Finance,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'engineering' THEN 1 END) AS Engineering,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'coal loading' THEN 1 END) AS CoalLoading,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'stockpile' THEN 1 END) AS Stockpile,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'shippping' THEN 1 END) AS Shipping,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant & logistic' THEN 1 END) AS PlantLogistic,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'keamanan & eksternal' THEN 1 END) AS keamanan_eksternal,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'oshe (operation & she)' THEN 1 END) AS Oshe,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 5  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'management' THEN 1 END) AS Management
		`).
		Where(queryFilterAll).
		Scan(&mei).Error

	if errMei != nil {
		return dashboardEmployees, errMei
	}
	dashboardEmployees.Mei = mei

	var juni DepartmentName

	errJuni := r.db.Model(&Employee{}).
		Joins("JOIN dohs ON employees.id = dohs.employee_id").
		Joins("JOIN departments ON employees.department_id = departments.id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'operation' THEN 1 END) AS Operation, 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant' THEN 1 END) AS Plant,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'hrga' THEN 1 END) AS Hrga,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'she' THEN 1 END) AS She,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'finance' THEN 1 END) AS Finance,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'engineering' THEN 1 END) AS Engineering,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'coal loading' THEN 1 END) AS CoalLoading,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'stockpile' THEN 1 END) AS Stockpile,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'shippping' THEN 1 END) AS Shipping,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant & logistic' THEN 1 END) AS PlantLogistic,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'keamanan & eksternal' THEN 1 END) AS keamanan_eksternal,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'oshe (operation & she)' THEN 1 END) AS Oshe,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 6  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'management' THEN 1 END) AS Management
		`).
		Where(queryFilterAll).
		Scan(&juni).Error

	if errJuni != nil {
		return dashboardEmployees, errJuni
	}
	dashboardEmployees.Juni = juni

	var juli DepartmentName
	errJuli := r.db.Model(&Employee{}).
		Joins("JOIN dohs ON employees.id = dohs.employee_id").
		Joins("JOIN departments ON employees.department_id = departments.id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'operation' THEN 1 END) AS Operation, 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant' THEN 1 END) AS Plant,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'hrga' THEN 1 END) AS Hrga,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'she' THEN 1 END) AS She,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'finance' THEN 1 END) AS Finance,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'engineering' THEN 1 END) AS Engineering,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'coal loading' THEN 1 END) AS CoalLoading,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'stockpile' THEN 1 END) AS Stockpile,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'shippping' THEN 1 END) AS Shipping,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant & logistic' THEN 1 END) AS PlantLogistic,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'keamanan & eksternal' THEN 1 END) AS keamanan_eksternal,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'oshe (operation & she)' THEN 1 END) AS Oshe,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 7  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'management' THEN 1 END) AS Management
		`).
		Where(queryFilterAll).
		Scan(&juli).Error

	if errJuli != nil {
		return dashboardEmployees, errJuli
	}
	dashboardEmployees.Juli = juli

	var agustus DepartmentName
	errAgustus := r.db.Model(&Employee{}).
		Joins("JOIN dohs ON employees.id = dohs.employee_id").
		Joins("JOIN departments ON employees.department_id = departments.id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'operation' THEN 1 END) AS Operation, 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant' THEN 1 END) AS Plant,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'hrga' THEN 1 END) AS Hrga,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'she' THEN 1 END) AS She,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'finance' THEN 1 END) AS Finance,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'engineering' THEN 1 END) AS Engineering,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'coal loading' THEN 1 END) AS CoalLoading,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'stockpile' THEN 1 END) AS Stockpile,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'shippping' THEN 1 END) AS Shipping,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant & logistic' THEN 1 END) AS PlantLogistic,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'keamanan & eksternal' THEN 1 END) AS keamanan_eksternal,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'oshe (operation & she)' THEN 1 END) AS Oshe,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 8  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'management' THEN 1 END) AS Management
		`).
		Where(queryFilterAll).
		Scan(&agustus).Error

	if errAgustus != nil {
		return dashboardEmployees, errAgustus
	}
	dashboardEmployees.Agustus = agustus

	var september DepartmentName

	errSeptember := r.db.Model(&Employee{}).
		Joins("JOIN dohs ON employees.id = dohs.employee_id").
		Joins("JOIN departments ON employees.department_id = departments.id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'operation' THEN 1 END) AS Operation, 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant' THEN 1 END) AS Plant,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'hrga' THEN 1 END) AS Hrga,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'she' THEN 1 END) AS She,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'finance' THEN 1 END) AS Finance,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'engineering' THEN 1 END) AS Engineering,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'coal loading' THEN 1 END) AS CoalLoading,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'stockpile' THEN 1 END) AS Stockpile,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'shippping' THEN 1 END) AS Shipping,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant & logistic' THEN 1 END) AS PlantLogistic,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'keamanan & eksternal' THEN 1 END) AS keamanan_eksternal,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'oshe (operation & she)' THEN 1 END) AS Oshe,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 9  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'management' THEN 1 END) AS Management
		`).
		Where(queryFilterAll).
		Scan(&september).Error

	if errSeptember != nil {
		return dashboardEmployees, errSeptember
	}
	dashboardEmployees.September = september

	var oktober DepartmentName
	errOktober := r.db.Model(&Employee{}).
		Joins("JOIN dohs ON employees.id = dohs.employee_id").
		Joins("JOIN departments ON employees.department_id = departments.id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'operation' THEN 1 END) AS Operation, 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant' THEN 1 END) AS Plant,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'hrga' THEN 1 END) AS Hrga,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'she' THEN 1 END) AS She,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'finance' THEN 1 END) AS Finance,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'engineering' THEN 1 END) AS Engineering,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'coal loading' THEN 1 END) AS CoalLoading,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'stockpile' THEN 1 END) AS Stockpile,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'shippping' THEN 1 END) AS Shipping,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant & logistic' THEN 1 END) AS PlantLogistic,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'keamanan & eksternal' THEN 1 END) AS keamanan_eksternal,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'oshe (operation & she)' THEN 1 END) AS Oshe,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 10  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'management' THEN 1 END) AS Management
		`).
		Where(queryFilterAll).
		Scan(&oktober).Error

	if errOktober != nil {
		return dashboardEmployees, errOktober
	}
	dashboardEmployees.Oktober = oktober

	var november DepartmentName

	errNovember := r.db.Model(&Employee{}).
		Joins("JOIN dohs ON employees.id = dohs.employee_id").
		Joins("JOIN departments ON employees.department_id = departments.id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'operation' THEN 1 END) AS Operation, 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant' THEN 1 END) AS Plant,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'hrga' THEN 1 END) AS Hrga,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'she' THEN 1 END) AS She,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'finance' THEN 1 END) AS Finance,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'engineering' THEN 1 END) AS Engineering,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'coal loading' THEN 1 END) AS CoalLoading,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'stockpile' THEN 1 END) AS Stockpile,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'shippping' THEN 1 END) AS Shipping,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant & logistic' THEN 1 END) AS PlantLogistic,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'keamanan & eksternal' THEN 1 END) AS keamanan_eksternal,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'oshe (operation & she)' THEN 1 END) AS Oshe,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 11  AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'management' THEN 1 END) AS Management
		`).
		Where(queryFilterAll).
		Scan(&november).Error

	if errNovember != nil {
		return dashboardEmployees, errNovember
	}
	dashboardEmployees.November = november

	var desember DepartmentName
	errDesember := r.db.Model(&Employee{}).
		Joins("JOIN dohs ON employees.id = dohs.employee_id").
		Joins("JOIN departments ON employees.department_id = departments.id").
		Select(` 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 12 AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'operation' THEN 1 END) AS Operation, 
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 12 AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant' THEN 1 END) AS Plant,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 12 AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'hrga' THEN 1 END) AS Hrga,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 12 AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'she' THEN 1 END) AS She,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 12 AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'finance' THEN 1 END) AS Finance,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 12 AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'engineering' THEN 1 END) AS Engineering,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 12 AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'coal loading' THEN 1 END) AS CoalLoading,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 12 AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'stockpile' THEN 1 END) AS Stockpile,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 12 AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'shippping' THEN 1 END) AS Shipping,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 12 AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'plant & logistic' THEN 1 END) AS PlantLogistic,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 12 AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'keamanan & eksternal' THEN 1 END) AS keamanan_eksternal,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 12 AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'oshe (operation & she)' THEN 1 END) AS Oshe,
			COUNT(CASE WHEN EXTRACT(MONTH FROM dohs.tanggal_end_doh::DATE) = 12 AND cast(employees.status AS TEXT) ILIKE 'AKTIF' AND LOWER(departments.department_name) LIKE 'management' THEN 1 END) AS Management
		`).
		Where(queryFilterAll).
		Scan(&desember).Error

	if errDesember != nil {
		return dashboardEmployees, errDesember
	}
	dashboardEmployees.Desember = desember

	return dashboardEmployees, nil
}
