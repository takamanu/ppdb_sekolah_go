package models

import "time"

type Datapokok struct {
	ID           uint64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID       uint64     `gorm:"column:user_id" json:"user_id"`
	Email        string     `gorm:"required;unique;column:email" json:"email"`
	PasFoto      string     `gorm:"required;column:pas_foto" json:"pas_foto"`
	NamaLengkap  string     `gorm:"required;column:nama_lengkap" json:"nama_lengkap"`
	NISN         string     `gorm:"required;unique;column:nisn" json:"nisn"`
	JenisKelamin string     `gorm:"required;column:jenis_kelamin" json:"jenis_kelamin"`
	TempatLahir  string     `gorm:"required;column:tempat_lahir" json:"tempat_lahir"`
	TanggalLahir *time.Time `gorm:"required;column:tanggal_lahir" json:"tanggal_lahir"`
	AsalSekolah  string     `gorm:"required;column:asal_sekolah" json:"asal_sekolah"`
	NamaAyah     string     `gorm:"required;column:nama_ayah" json:"nama_ayah"`
	NoWaAyah     string     `gorm:"required;column:no_wa_ayah" json:"no_wa_ayah"`
	NamaIbu      string     `gorm:"required;column:nama_ibu" json:"nama_ibu"`
	NoWaIbu      string     `gorm:"required;column:no_wa_ibu" json:"no_wa_ibu"`
	CreatedAt    *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at" json:"updated_at"`
	Nilai        []Nilai    `json:"nilai"  form:"nilai"`
}

func (Datapokok) TableName() string {
	return "datapokok"
}
