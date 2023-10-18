package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserName string `gorm:"type:varchar(255);unique_index;not null"`
	PassWord string `gorm:"type:varchar(255);not null"`
	Salt     int    `gorm:"not null"`
	Phone    string `gorm:"type:varchar(20);unique_index"`
	Email    string `gorm:"type:varchar(100);unique_index"`
}
