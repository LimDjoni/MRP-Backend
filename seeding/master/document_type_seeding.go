package seeding

import (
	"ajebackend/model/master/documenttype"
	"fmt"

	"gorm.io/gorm"
)

func SeedingDocumentType(db *gorm.DB) {

	tx := db.Begin()
	var checkDocumentType []documenttype.DocumentType
	tx.Find(&checkDocumentType)

	if len(checkDocumentType) > 0 {
		return
	}

	var createDocumentType []documenttype.DocumentType
	createDocumentType = append(createDocumentType,
		documenttype.DocumentType{Name: "BC 1.6 Pemberitahuan Pabean Pemasukan Barang Impor Untuk Ditimbun di Pusat Logistik Berikat", Code: "0207500"},
		documenttype.DocumentType{Name: "BC 2.0 Pemberitahuan Impor Barang (PIB)", Code: "0207501"},
		documenttype.DocumentType{Name: "BC 2.3 Pemberitahuan Impor Barang untuk Ditimbun di Tempat Penimbunan Berikat (TPB)", Code: "0207502"},
		documenttype.DocumentType{Name: "BC 2.5 Pemberitahuan Impor Barang dari Tempat Penimbunan Berikat (TPB)", Code: "0207503"},
		documenttype.DocumentType{Name: "BC 2.8 Pemberitahuan Impor Barang dari Pusat Logistik Berikat (PLB)", Code: "0207504"},
		documenttype.DocumentType{Name: "BC 3.3 Pemberitahuan Ekspor Barang (PEB) Dari Pusat Logistik Berikat (PLB)", Code: "0207002"},
		documenttype.DocumentType{Name: "BC 2.7  Pemberitahuan pengeluaran barang dari TPB ke TPB lainnya", Code: "0207518"},
		documenttype.DocumentType{Name: "P3BET", Code: "0207006"},
		documenttype.DocumentType{Name: "Pemberitahuan Ekspor Barang", Code: "0207001"},
		documenttype.DocumentType{Name: "PKBE (Pemberitahuan Konsolidasi Barang Ekspor)", Code: "0207005"},
		documenttype.DocumentType{Name: "PPFTZ 01 (Pemasukan asal LDP)", Code: "0207506"},
		documenttype.DocumentType{Name: "PPFTZ 01 (Pengeluaran dari Kawasan Bebas ke TLDDP)", Code: "0207505"},
		documenttype.DocumentType{Name: "PPFTZ (Pengeluaran dari Kawasan Bebas ke Luar DP)", Code: "0207003"},
		documenttype.DocumentType{Name: "PPKEK (Pemasukan asal LDP)", Code: "0207508"},
		documenttype.DocumentType{Name: "PPKEK (Pengeluaran dari KEK KE Luar DP)", Code: "0207004"},
		documenttype.DocumentType{Name: "PPKEK (Pengeluaran dari KEK ke TLDDP)", Code: "0207507"},
	)

	err := tx.Create(&createDocumentType).Error

	if err != nil {
		tx.Rollback()
		fmt.Println("Failed Seeding Document Type")
		return
	}

	tx.Commit()
}
