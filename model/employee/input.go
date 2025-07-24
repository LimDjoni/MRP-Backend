package employee

import (
	"mrpbackend/model/master/apd"
	"mrpbackend/model/master/bank"
	"mrpbackend/model/master/bpjskesehatan"
	"mrpbackend/model/master/bpjsketenagakerjaan"
	"mrpbackend/model/master/doh"
	"mrpbackend/model/master/history"
	"mrpbackend/model/master/jabatan"
	"mrpbackend/model/master/kartukeluarga"
	"mrpbackend/model/master/ktp"
	"mrpbackend/model/master/laporan"
	"mrpbackend/model/master/mcu"
	"mrpbackend/model/master/npwp"
	"mrpbackend/model/master/pendidikan"
	"mrpbackend/model/master/position"
	"mrpbackend/model/master/sertifikat"
)

type RegisterEmployeeInput struct {
	NomorKaryawan string `json:"nomor_karyawan" validate:"required"`
	DepartmentId  uint   `json:"department_id" validate:"required"`
	Firstname     string `json:"firstname" validate:"required"`
	Lastname      string `json:"lastname"`
	PhoneNumber   string `json:"phone_number"`
	Email         string `json:"email"`
	Level         string `json:"level" validate:"required"`
	Status        string `json:"status" validate:"required"`
	RoleId        uint   `json:"role_id" validate:"required"`
	HiredBy       string `json:"hire_by" validate:"required"`
	PositionId    uint   `json:"position_id" validate:"required"`
	DateOfHire    string `json:"date_of_hire"`

	// Related Data
	KartuKeluarga       kartukeluarga.KartuKeluarga
	KTP                 ktp.KTP
	Pendidikan          pendidikan.Pendidikan
	Laporan             laporan.Laporan
	APD                 *apd.APD
	NPWP                *npwp.NPWP
	Bank                bank.Bank
	BPJSKesehatan       *bpjskesehatan.BPJSKesehatan
	BPJSKetenagakerjaan *bpjsketenagakerjaan.BPJSKetenagakerjaan
	DOH                 []doh.DOH
	Jabatan             []jabatan.Jabatan
	Sertifikat          []sertifikat.Sertifikat
	MCU                 []mcu.MCU
	History             *[]history.History
	Position            position.Position
}

type UpdateEmployeeInput struct {
	NomorKaryawan string `json:"nomor_karyawan" validate:"required"`
	DepartmentId  uint   `json:"department_id" validate:"required"`
	Firstname     string `json:"firstname" validate:"required"`
	Lastname      string `json:"lastname"`
	PhoneNumber   string `json:"phone_number"`
	Email         string `json:"email"`
	Level         string `json:"level" validate:"required"`
	Status        string `json:"status" validate:"required"`
	RoleId        uint   `json:"role_id" validate:"required"`
	HiredBy       string `json:"hire_by" validate:"required"`
	PositionId    uint   `json:"position_id" validate:"required"`
	DateOfHire    string `json:"date_of_hire"`

	// Related Data
	KartuKeluarga       kartukeluarga.KartuKeluarga
	KTP                 ktp.KTP
	Pendidikan          pendidikan.Pendidikan
	Laporan             laporan.Laporan
	APD                 *apd.APD
	NPWP                *npwp.NPWP
	Bank                bank.Bank
	BPJSKesehatan       *bpjskesehatan.BPJSKesehatan
	BPJSKetenagakerjaan *bpjsketenagakerjaan.BPJSKetenagakerjaan
	Position            position.Position
}

type SortFilterEmployee struct {
	Field                 string
	Sort                  string
	CodeEmp               string `json:"code_emp" validate:"required"`
	NomorKaryawan         string `json:"nomor_karyawan" validate:"required"`
	DepartmentId          string `json:"department_id" validate:"required"`
	Firstname             string `json:"firstname" validate:"required"`
	HireBy                string `json:"hire_by"`
	Agama                 string `json:"agama"`
	Level                 string `json:"level" validate:"required"`
	Gender                string `json:"gender"`
	KategoriLokalNonLokal string `json:"kategori_lokal_non_lokal"`
	KategoriTriwulan      string `json:"kategori_triwulan"`
	Status                string `json:"status"`
	Kontrak               string `json:"kontrak"`
	RoleId                string `json:"role_id"`
	PositionId            string `json:"position_id"`
}

type BasedOnRing struct {
	Ring1    uint `json:"ring_1"`
	Ring2    uint `json:"ring_2"`
	Ring3    uint `json:"ring_3"`
	LuarRing uint `json:"luar_ring"`
}

type BasedOnDepartment struct {
	Engineering       uint `json:"engineering"`
	Finance           uint `json:"finance"`
	HRGA              uint `json:"hrga"`
	Operation         uint `json:"operation"`
	Plant             uint `json:"plant"`
	SHE               uint `json:"she"`
	CoalLoading       uint `json:"coal_loading"`
	Stockpile         uint `json:"stockpile"`
	Shipping          uint `json:"shipping"`
	PlantLogistic     uint `json:"plant_logistic"`
	KeamananEksternal uint `json:"keamanan_eksternal"`
	Oshe              uint `json:"oshe"`
	Management        uint `json:"management"`
}

type BasedOnEducation struct {
	Edu1 uint `json:"education_1"`
	Edu2 uint `json:"education_2"`
	Edu3 uint `json:"education_3"`
	Edu4 uint `json:"education_4"`
	Edu5 uint `json:"education_5"`
}

type BasedOnYear struct {
	Year1 uint `json:"year_1"`
	Year2 uint `json:"year_2"`
	Year3 uint `json:"year_3"`
	Year4 uint `json:"year_4"`
}

type BasedOnAge struct {
	Stage1 uint `json:"stage_1"`
	Stage2 uint `json:"stage_2"`
	Stage3 uint `json:"stage_3"`
	Stage4 uint `json:"stage_4"`
	Stage5 uint `json:"stage_5"`
}

type BasedOnLokal struct {
	Lokal    uint `json:"lokal"`
	NonLokal uint `json:"non_lokal"`
}

type DashboardEmployee struct {
	TotalEmployee     uint              `json:"total_employee"`
	TotalMale         uint              `json:"total_male"`
	TotalFemale       uint              `json:"total_female"`
	HireHO            uint              `json:"hired_ho"`
	HireSite          uint              `json:"hired_site"`
	BasedOnAge        BasedOnAge        `json:"based_on_age"`
	BasedOnYear       BasedOnYear       `json:"based_on_year"`
	BasedOnEducation  BasedOnEducation  `json:"based_on_education"`
	BasedOnDepartment BasedOnDepartment `json:"based_on_department"`
	BasedOnRing       BasedOnRing       `json:"based_on_ring"`
	BasedOnLokal      BasedOnLokal      `json:"based_on_lokal"`
}

type SortFilterDashboardEmployee struct {
	PT           string `json:"pt"`
	DepartmentId string `json:"department_id"`
}

type DataStatus struct {
	NewHire      uint `json:"new_hire"`
	BerakhirPkwt uint `json:"berakhir_pkwt"`
	Resign       uint `json:"resign"`
	PHK          uint `json:"phk"`
}

type SortFilterDashboardEmployeeTurnOver struct {
	PT   string `json:"pt"`
	Year string `json:"year"`
}

type DashboardEmployeeTurnOver struct {
	TotalHire         uint       `json:"total_hire"`
	TotalResign       uint       `json:"total_resign"`
	TotalBerakhirPkwt uint       `json:"total_berakhir_pkwt"`
	TotalPhk          uint       `json:"total_phk"`
	Januari           DataStatus `json:"januari"`
	Februari          DataStatus `json:"februari"`
	Maret             DataStatus `json:"maret"`
	April             DataStatus `json:"april"`
	Mei               DataStatus `json:"mei"`
	Juni              DataStatus `json:"juni"`
	Juli              DataStatus `json:"juli"`
	Agustus           DataStatus `json:"agustus"`
	September         DataStatus `json:"september"`
	Oktober           DataStatus `json:"oktober"`
	November          DataStatus `json:"november"`
	Desember          DataStatus `json:"desember"`
}

type DepartmentName struct {
	Operation         uint `json:"operation"`
	Plant             uint `json:"plant"`
	Hrga              uint `json:"hrga"`
	She               uint `json:"she"`
	Finance           uint `json:"finance"`
	Engineering       uint `json:"engineering"`
	CoalLoading       uint `json:"coal_loading"`
	Stockpile         uint `json:"stockpile"`
	Shipping          uint `json:"shipping"`
	PlantLogistic     uint `json:"plant_logistic"`
	KeamananEksternal uint `json:"keamanan_eksternal"`
	Oshe              uint `json:"oshe"`
	Management        uint `json:"management"`
}

type DashboardEmployeeKontrak struct {
	Januari   DepartmentName `json:"januari"`
	Februari  DepartmentName `json:"februari"`
	Maret     DepartmentName `json:"maret"`
	April     DepartmentName `json:"april"`
	Mei       DepartmentName `json:"mei"`
	Juni      DepartmentName `json:"juni"`
	Juli      DepartmentName `json:"juli"`
	Agustus   DepartmentName `json:"agustus"`
	September DepartmentName `json:"september"`
	Oktober   DepartmentName `json:"oktober"`
	November  DepartmentName `json:"november"`
	Desember  DepartmentName `json:"desember"`
}
