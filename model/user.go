package model

import (
	"github.com/HauKuen/Annals/utils/errmsg"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique" json:"username"`
	Password string `gorm:"not null" json:"password"`
	Role     int    `gorm:"not null" json:"role"`
}

type APIUser struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Role      int    `json:"role"`
	CreatedAt string `json:"created_at"`
}

// GetUser 查询用户
func GetUser(id int) (User, int) {
	var user User
	err = db.Limit(1).Where("ID = ?", id).Find(&user).Error
	if err != nil {
		return user, errmsg.ERROR
	}
	return user, errmsg.SUCCESS
}

// GetUsers 查询用户列表
func GetUsers(pageSize int, pageNum int) ([]APIUser, int64) {
	var users []APIUser
	var total int64
	offset := (pageNum - 1) * pageSize
	db.Model(&User{}).Limit(pageSize).Offset(offset).Find(&users)
	db.Model(&User{}).Count(&total)
	return users, total
}

// CreateUser 添加用户
func CreateUser(data *User) (code int) {
	err = db.Create(&data).Error
	if err != nil {
		return errmsg.ERROR // 500
	}
	return errmsg.SUCCESS
}

// CheckUser 查询用户是否存在
func CheckUser(username string) (code int) {
	var user User
	db.Select("id").Where("username = ?", username).First(&user)
	if user.ID > 0 {
		return errmsg.ErrorUsernameUsed // 1001
	}
	return errmsg.SUCCESS
}
