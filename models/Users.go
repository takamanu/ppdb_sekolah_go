package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string      `gorm:"required" json:"name" form:"name"`
	Email     string      `gorm:"required, unique" json:"email"  form:"email"`
	Password  string      `gorm:"required" json:"password" form:"password"`
	Role      uint        `gorm:"default:0" json:"role"  form:"role"`
	Datapokok []Datapokok `json:"datapokok"  form:"datapokok"`
}

type UserResponse struct {
	ID    uint   `json:"id" form:"name"`
	Name  string `json:"name" form:"name"`
	Email string `json:"email"  form:"email"`
	Role  uint   `json:"role"  form:"role"`
	Token string `json:"token" form:"token"`
}
