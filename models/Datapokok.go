package models

import "time"

type Datapokok struct {
	ID                     uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID                 uint64     `gorm:"column:user_id" json:"user_id"`
	Email                  string     `gorm:"column:email" json:"email"`
	PasFoto                string     `gorm:"column:pas_foto" json:"pas_foto"`
	NamaLengkap            string     `gorm:"column:nama_lengkap" json:"nama_lengkap"`
	NISN                   string     `gorm:"column:nisn" json:"nisn"`
	JenisKelamin           string     `gorm:"column:jenis_kelamin" json:"jenis_kelamin"`
	TempatLahir            string     `gorm:"column:tempat_lahir" json:"tempat_lahir"`
	TanggalLahir           *time.Time `gorm:"column:tanggal_lahir" json:"tanggal_lahir"`
	Agama                  string     `gorm:"column:agama" json:"agama"`
	AsalSekolah            string     `gorm:"column:asal_sekolah" json:"asal_sekolah"`
	AlamatSekolah          string     `gorm:"column:alamat_sekolah" json:"alamat_sekolah"`
	JumlahHafalan          string     `gorm:"column:jumlah_hafalan" json:"jumlah_hafalan"`
	NamaAyah               string     `gorm:"column:nama_ayah" json:"nama_ayah"`
	PekerjaanAyah          string     `gorm:"column:pekerjaan_ayah" json:"pekerjaan_ayah"`
	PenghasilanAyah        string     `gorm:"column:penghasilan_ayah" json:"penghasilan_ayah"`
	PendidikanTerakhirAyah string     `gorm:"column:pendidikan_terakir_ayah" json:"pendidikan_terakir_ayah"`
	NoWaAyah               string     `gorm:"column:no_wa_ayah" json:"no_wa_ayah"`
	NamaIbu                string     `gorm:"column:nama_ibu" json:"nama_ibu"`
	PekerjaanIbu           string     `gorm:"column:pekerjaan_ibu" json:"pekerjaan_ibu"`
	PenghasilanIbu         string     `gorm:"column:penghasilan_ibu" json:"penghasilan_ibu"`
	PendidikanTerakhirIbu  string     `gorm:"column:pendidikan_terakir_ibu" json:"pendidikan_terakir_ibu"`
	NoWaIbu                string     `gorm:"column:no_wa_ibu" json:"no_wa_ibu"`
	NamaWaliSiswa          *string    `gorm:"column:nama_wali_siswa" json:"nama_wali_siswa"`
	HubunganDenganSiswa    *string    `gorm:"column:hubungan_dengan_siswa" json:"hubungan_dengan_siswa"`
	AlamatWaliSiswa        *string    `gorm:"column:alamat_wali_siswa" json:"alamat_wali_siswa"`
	PekerjaanWali          *string    `gorm:"column:pekerjaan_wali" json:"pekerjaan_wali"`
	PenghasilanWali        *string    `gorm:"column:penghasilan_wali" json:"penghasilan_wali"`
	PendidikanTerakhirWali *string    `gorm:"column:pendidikan_terakir_wali" json:"pendidikan_terakir_wali"`
	NoWaWaliSiswa          *string    `gorm:"column:no_wa_wali_siswa" json:"no_wa_wali_siswa"`
	Motivasi               string     `gorm:"column:motivasi" json:"motivasi"`
	DaftarSekolahLain      bool       `gorm:"column:daftar_sekolah_lain" json:"daftar_sekolah_lain"`
	NamaSekolahJikaDaftar  string     `gorm:"column:nama_sekolahnya_jika_daftar" json:"nama_sekolahnya_jika_daftar"`
	InformasiDapatkanDari  string     `gorm:"column:informasi_didapatkan_dari" json:"informasi_didapatkan_dari"`
	CreatedAt              *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt              *time.Time `gorm:"column:updated_at" json:"updated_at"`
	Nilai                  []Nilai    `json:"nilai"  form:"nilai"`
}

func (Datapokok) TableName() string {
	return "datapokok"
}
