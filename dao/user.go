package dao

import (
	"testProject/test/models"

	"gorm.io/gorm"
)

func GetUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	user := &models.User{}
	if err := db.Where("user_name = ?", username).First(&user).Error; err != nil {
		return nil, err // 查询出错
	}
	return user, nil
}

func CreateUser(db *gorm.DB, user *models.User) error {
	// 在这里执行插入操作
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func UpdatePassword(db *gorm.DB, newPasswrod int) error {
	if err := db.Update("password", newPasswrod).Error; err != nil {
		return err
	}
	return nil
}
