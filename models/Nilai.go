package models

import "time"

type Nilai struct {
	Utama               uint64    `gorm:"column:Utama;primaryKey;autoIncrement" json:"id"`
	DataPokokID         uint64    `gorm:"column:datapokok_id" json:"datapokok_id"`
	Matematika          string    `gorm:"column:matematika" json:"matematika"`
	IlmuPengetahuanAlam string    `gorm:"column:ilmu_pengetahuan_alam" json:"ilmu_pengetahuan_alam"`
	BahasaIndonesia     string    `gorm:"column:bahasa_indonesia" json:"bahasa_indonesia"`
	TestMembacaAlQuran  string    `gorm:"column:test_membaca_al_quran" json:"test_membaca_al_quran"`
	Status              string    `gorm:"column:status" json:"status"`
	CreatedAt           time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Nilai) TableName() string {
	return "nilai"
}
