package models

type User struct {
	ID        uint   `gorm:"primarykey"`
	UserName  string `gorm:"type:varchar(255);unique_index;not null"`
	PassWord  string `gorm:"type:varchar(255);not null"`
	Salt      string `gorm:"not null"`
	Phone     string `gorm:"type:varchar(20);unique_index"`
	Email     string `gorm:"type:varchar(100);unique_index"`
	CreatedAt int
	UpdatedAt int
	DeletedAt int
}
