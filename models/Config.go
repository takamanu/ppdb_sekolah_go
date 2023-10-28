package models

import "time"

type Config struct {
	ID         uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Pengumuman bool      `gorm:"column:pengumuman" json:"pengumuman"`
	RedirectWA string    `gorm:"column:redirect_wa" json:"redirect_wa"`
	CreatedAt  time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt  time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func (Config) TableName() string {
	return "config"
}
