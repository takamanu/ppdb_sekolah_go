package models

import "time"

type Nilai struct {
	Utama               uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	DataPokokID         uint64    `gorm:"column:datapokok_id" json:"datapokok_id"`
	Matematika          uint      `gorm:"column:matematika" json:"matematika"`
	IlmuPengetahuanAlam uint      `gorm:"column:ilmu_pengetahuan_alam" json:"ilmu_pengetahuan_alam"`
	BahasaIndonesia     uint      `gorm:"column:bahasa_indonesia" json:"bahasa_indonesia"`
	TestMembacaAlQuran  uint      `gorm:"column:test_membaca_al_quran" json:"test_membaca_al_quran"`
	Status              string    `gorm:"column:status" json:"status"`
	CreatedAt           time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Nilai) TableName() string {
	return "nilai"
}
