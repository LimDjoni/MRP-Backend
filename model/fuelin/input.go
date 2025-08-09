package fuelin

type RegisterFuelInInput struct {
	Date            string  `json:"date"`
	Vendor          string  `json:"vendor"`
	Code            string  `json:"code"`
	NomorSuratJalan string  `json:"nomor_surat_jalan"`
	NomorPlatMobil  string  `json:"nomor_plat_mobil"`
	Qty             float64 `json:"qty"`
	QtyNow          float64 `json:"qty_now"`
	Driver          string  `json:"driver"`
	TujuanAwal      string  `json:"tujuan_awal"`
}

type SortFilterFuelIn struct {
	Field           string
	Sort            string
	Date            string `json:"date"`
	Vendor          string `json:"vendor"`
	Code            string `json:"code"`
	NomorSuratJalan string `json:"nomor_surat_jalan"`
	NomorPlatMobil  string `json:"nomor_plat_mobil"`
	Qty             string `json:"qty"`
	QtyNow          string `json:"qty_now"`
	Driver          string `json:"driver"`
	TujuanAwal      string `json:"tujuan_awal"`
}
