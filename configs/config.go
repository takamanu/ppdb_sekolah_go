package configs

import (
	"ppdb_sekolah_go/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	dsn := "root:@tcp(127.0.0.1:3306)/ppdb_smp?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Nilai{})
	DB.AutoMigrate(&models.Datapokok{})
	DB.AutoMigrate(&models.Config{})

}
