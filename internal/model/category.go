package model

import (
	"fmt"

	"github.com/HauKuen/Annals/internal/utils/respcode"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);uniqueIndex;not null" json:"name" validate:"required"`
}

// CreateCategory 创建分类
func CreateCategory(data *Category) int {
	if err := db.Create(data).Error; err != nil {
		fmt.Print(err)
		return respcode.ERROR
	}
	return respcode.SUCCESS
}

// CheckCategory 查询分类是否存在
func CheckCategory(name string) int {
	var category Category
	if err := db.Select("id").Where("name = ?", name).First(&category).Error; err != nil {
		return respcode.ErrorCateNotExist
	}
	return respcode.SUCCESS
}

// DeleteCategory 删除分类
func DeleteCategory(id int) int {
	var category Category
	// 查询分类是否存在
	if err := db.Select("id").Where("id = ?", id).First(&category).Error; err != nil {
		return respcode.ErrorCateNotExist
	}

	if err := db.Where("id = ?", id).Delete(&Category{}).Error; err != nil {
		return respcode.ERROR
	}
	return respcode.SUCCESS
}

// GetCategory 查询分类
func GetCategory(id int) (Category, int) {
	var category Category
	if err := db.Where("id = ?", id).First(&category).Error; err != nil {
		return category, respcode.ErrorCateNotExist
	}
	return category, respcode.SUCCESS
}

// GetCategories 获取所有分类
func GetCategories() ([]Category, int) {
	var categories []Category

	if err := db.Find(&categories).Error; err != nil {
		return nil, respcode.ERROR
	}

	return categories, respcode.SUCCESS
}
