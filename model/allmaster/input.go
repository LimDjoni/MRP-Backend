package allmaster

type RegisterUserRoleInput struct {
	UserId uint `json:"user_id" validate:"required"`
	RoleId uint `json:"role_id" validate:"required"`
}

type RegisterBrandInput struct {
	BrandName string `json:"brand_name" validate:"required"`
}

type RegisterHeavyEquipmentInput struct {
	BrandId            uint   `json:"brand_id" validate:"required"`
	HeavyEquipmentName string `json:"heavy_equipment_name" validate:"required"`
}

type RegisterSeriesInput struct {
	BrandId          uint   `json:"brand_id" validate:"required"`
	HeavyEquipmentId uint   `json:"heavy_equipment_id" validate:"required"`
	SeriesName       string `json:"series_name" validate:"required"`
}

type RegisterKartuKeluargaInput struct {
	NomorKartuKeluarga    string `json:"nomor_kartu_keluarga"`
	NamaIbuKandung        string `json:"nama_ibu_kandung"`
	KontakDarurat         string `json:"kontak_darurat"`
	NamaKontakDarurat     string `json:"nama_kontak_darurat"`
	HubunganKontakDarurat string `json:"hubungan_kontak_darurat"`
}

type RegisterKTPInput struct {
	NamaSesuaiKTP string `json:"nama_sesuai_ktp"`
	NomorKTP      string `json:"nomor_ktp"`
	TempatLahir   string `json:"tempat_lahir"`
	TanggalLahir  string `json:"tanggal_lahir"`
	Gender        string `json:"gender"`
	Alamat        string `json:"alamat"`
	RT            string `json:"rt"`
	RW            string `json:"rw"`
	Kel           string `json:"kel"`
	Kec           string `json:"kec"`
	Kota          string `json:"kota"`
	Prov          string `json:"prov"`
	KodePos       string `json:"kode_pos"`
	GolonganDarah string `json:"golongan_darah"`
	Agama         string `json:"agama"`
	RingKTP       string `json:"ringktp"`
}

type RegisterPendidikanInput struct {
	PendidikanLabel    string `json:"pendidikan_label"`
	PendidikanTerakhir string `json:"pendidikan_terakhir"`
	Jurusan            string `json:"jurusan"`
}

type RegisterDOHInput struct {
	EmployeeId    uint   `json:"employee_id"`
	TanggalDoh    string `json:"tanggal_doh"`
	TanggalEndDoh string `json:"tanggal_end_doh"`
	PT            string `json:"pt"`
	Penempatan    string `json:"penempatan"`
	StatusKontrak string `json:"status_kontrak"`
}

type RegisterJabatanInput struct {
	EmployeeId uint   `json:"employee_id"`
	DateMove   string `json:"date_move"`
	PositionId uint   `json:"position_id"`
}

type RegisterSertifikatInput struct {
	EmployeeId    uint   `json:"employee_id"`
	DateEffective string `json:"date_effective"`
	Sertifikat    string `json:"sertifikat"`
	Remark        string `json:"remark"`
}

type RegisterMCUInput struct {
	EmployeeId uint   `json:"employee_id"`
	DateMCU    string `json:"date_mcu"`
	DateEndMCU string `json:"date_end_mcu"`
	HasilMCU   string `json:"hasil_mcu"`
	MCU        string `json:"mcu"`
}

type RegisterHistoryInput struct {
	EmployeeId     uint   `json:"employee_id"`
	StatusTerakhir string `json:"status_terakhir"`
	Tanggal        string `json:"tanggal"`
	Keterangan     string `json:"keterangan"`
}

type EmployeeDOHExpired struct {
	ID             uint   `json:"id"`
	EmployeeID     uint   `json:"employee_id"`
	TanggalDOH     string `json:"tanggal_doh"`
	TanggalEndDOH  string `json:"tanggal_end_doh"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	DepartmentName string `json:"department_name"`
	PositionName   string `json:"position_name"`
}

type EmployeeMCUBerkala struct {
	ID             uint   `json:"id"`
	EmployeeID     uint   `json:"employee_id"`
	DateMCU        string `json:"date_mcu"`
	DateEndMCU     string `json:"date_end_mcu"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	DepartmentName string `json:"department_name"`
	PositionName   string `json:"position_name"`
}

type SortFilterDohKontrak struct {
	CodeEmp string `json:"code_emp"`
	PT      string `json:"pt"`
	Year    string `json:"year"`
	Field   string
	Sort    string
}
