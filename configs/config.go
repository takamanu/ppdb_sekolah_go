package configs

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {

	dsn := "root:@tcp(127.0.0.1:3306)/my_fashion?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&models.User{})

}
