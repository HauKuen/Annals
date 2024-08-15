package model

import (
	"errors"
	"fmt"
	"github.com/HauKuen/Annals/utils/respcode"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"time"
)

type User struct {
	gorm.Model
	Username    string     `gorm:"unique" json:"username"`
	Password    string     `gorm:"not null" json:"password"`
	Email       string     `gorm:"unique;not null" json:"email"`
	Role        int        `gorm:"not null" json:"role"`
	DisplayName string     `gorm:"not null" json:"display_name"`
	Bio         string     `json:"bio"`
	AvatarURL   string     `json:"avatar_url"`
	LastLogin   *time.Time `json:"last_login"`
	IsActive    bool       `gorm:"default:true" json:"is_active"`
}

type APIUser struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Role        int    `json:"role"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
	AvatarURL   string `json:"avatar_url"`
	CreatedAt   string `json:"created_at"`
	LastLogin   string `json:"last_login"`
	IsActive    bool   `json:"is_active"`
}

// BeforeCreate GORM 的钩子，在创建用户前加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	u.Password = string(hashedPassword)
	return nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// GetUser 查询用户
func GetUser(id int) (User, int) {
	var user User
	err := db.Limit(1).Where("ID = ?", id).Find(&user).Error
	if err != nil {
		return user, respcode.ERROR
	}
	return user, respcode.SUCCESS
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
func CreateUser(data *User) int {
	err := db.Create(data).Error
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		return respcode.ERROR // 500
	}
	return respcode.SUCCESS
}

// CheckUser 查询用户是否存在
func CheckUser(username string) int {
	var user User
	db.Select("id").Where("username = ?", username).First(&user)
	if user.ID > 0 {
		return respcode.ErrorUsernameUsed // 1001
	}
	return respcode.SUCCESS
}

// DeleteUser 删除用户
func DeleteUser(id int) int {
	var user User
	// 先查询用户是否存在
	err := db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return respcode.ErrorUserNotExist // 1003
		}
		return respcode.ERROR
	}

	// 用户存在，进行删除操作
	err = db.Delete(&user).Error
	if err != nil {
		return respcode.ERROR
	}
	return respcode.SUCCESS
}
