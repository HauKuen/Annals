package model

import "github.com/HauKuen/Annals/utils/respcode"

type Category struct {
	Name string `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
}

// CreateCategory 创建分类
func CreateCategory(data *Category) int {
	if err := db.Create(data).Error; err != nil {
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
