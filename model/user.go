package model

import (
	"errors"
	"fmt"
	"github.com/HauKuen/Annals/utils"
	"github.com/HauKuen/Annals/utils/respcode"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

type User struct {
	gorm.Model
	Username    string     `gorm:"unique;not null" json:"username" validate:"required"`
	Password    string     `gorm:"not null" json:"password" validate:"required"`
	Email       string     `gorm:"unique;not null" json:"email" validate:"required,email"`
	Role        int        `gorm:"not null" json:"role" validate:"required"`
	DisplayName string     `json:"display_name"`
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

// BeforeCreate 在创建用户前加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.Log.Error("Failed to hash password:", err)
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
func GetUser(id int) (APIUser, int) {
	var apiUser APIUser

	result := db.Model(&User{}).First(&apiUser, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			utils.Log.Error("User not found", err)
			return apiUser, respcode.ErrorUserNotExist
		}
		utils.Log.Error("Failed to get user:", err)
		return apiUser, respcode.ERROR
	}

	return apiUser, respcode.SUCCESS
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
		// 检查 MySQL 错误代码 1062（重复键错误）
		if strings.Contains(err.Error(), "Error 1062") {
			// 判断是哪个字段导致了重复错误
			if strings.Contains(err.Error(), "user.uni_user_username") {
				return respcode.ErrorUsernameUsed
			}
			if strings.Contains(err.Error(), "user.uni_user_email") {
				return respcode.ErrorEmailUsed
			}
		}
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
		message := fmt.Sprintf("User %s already exists", user.Username)
		utils.Log.Error(message)
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

// EditUser 编辑用户信息
func EditUser(id int, data *User) int {
	var user User
	var count int64

	// 查找要编辑的用户
	err := db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return respcode.ErrorUserNotExist
		}
		return respcode.ERROR
	}

	if data.Email != "" && data.Email != user.Email {
		db.Model(&User{}).Where("email = ? AND id != ?", data.Email, id).Count(&count)
		if count > 0 {
			return respcode.ErrorEmailUsed
		}
	}

	var updateFields []string
	if data.Email != "" {
		updateFields = append(updateFields, "email")
	}
	if data.Role != 0 {
		updateFields = append(updateFields, "role")
	}
	if data.DisplayName != "" {
		updateFields = append(updateFields, "display_name")
	}
	if data.Bio != "" {
		updateFields = append(updateFields, "bio")
	}
	if data.AvatarURL != "" {
		updateFields = append(updateFields, "avatar_url")
	}
	if data.IsActive != user.IsActive {
		updateFields = append(updateFields, "is_active")
	}

	// 执行更新
	err = db.Model(&user).Select(updateFields).Updates(data).Error
	if err != nil {
		return respcode.ERROR
	}

	return respcode.SUCCESS
}
