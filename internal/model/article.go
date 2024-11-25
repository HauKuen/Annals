package model

import (
	"github.com/HauKuen/Annals/internal/utils"
	"github.com/HauKuen/Annals/internal/utils/respcode"
	"gorm.io/gorm"
)

type Article struct {
	gorm.Model
	Title      string `gorm:"type:varchar(100);not null" json:"title"`
	Content    string `gorm:"type:longtext;not null" json:"content"`
	Img        string `gorm:"type:varchar(200)" json:"img"`
	CategoryID uint   `gorm:"not null" json:"category_id"`
	UserID     uint   `gorm:"not null" json:"user_id"`

	// 关联
	Category Category `gorm:"foreignKey:CategoryID" json:"category"`
	User     User     `gorm:"foreignKey:UserID" json:"user"`
}

// GetArticles 获取文章列表
func GetArticles(pageSize int, pageNum int) ([]Article, int64) {
	var articles []Article
	var total int64
	offset := (pageNum - 1) * pageSize

	db.Model(&Article{}).Count(&total)
	db.Preload("Category").Preload("User").
		Limit(pageSize).
		Offset(offset).
		Find(&articles)

	return articles, total
}

// GetArticleByID 获取单个文章信息
func GetArticleByID(id int) (Article, int) {
	var article Article
	err := db.Preload("Category").Preload("User").First(&article, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return article, respcode.ErrorArtNotExist
		}
		return article, respcode.ERROR
	}
	return article, respcode.SUCCESS
}

// CreateArticle 创建文章
func CreateArticle(article *Article) int {
	// 检查标题是否为空
	if article.Title == "" {
		return respcode.ErrorArtTitleEmpty
	}

	// 检查内容是否为空
	if article.Content == "" {
		return respcode.ErrorArtContent
	}

	// 检查分类是否存在
	var category Category
	if err := db.First(&category, article.CategoryID).Error; err != nil {
		return respcode.ErrorCateNotExist
	}

	// 创建文章
	if err := db.Create(article).Error; err != nil {
		utils.Log.Error("创建文章失败:", err)
		return respcode.ERROR
	}

	// 加载关联的分类和用户信息
	if err := db.Preload("Category").Preload("User").First(article, article.ID).Error; err != nil {
		utils.Log.Error("加载文章关联信息失败:", err)
		return respcode.ERROR
	}

	return respcode.SUCCESS
}

// UpdateArticle 更新文章
func UpdateArticle(id int, article *Article) int {
	var existingArticle Article

	// 检查文章是否存在
	if err := db.First(&existingArticle, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return respcode.ErrorArtNotExist
		}
		return respcode.ERROR
	}

	// 检查分类是否存在
	if article.CategoryID != 0 {
		var category Category
		if err := db.First(&category, article.CategoryID).Error; err != nil {
			return respcode.ErrorCateNotExist
		}
	}

	// 只更新非空字段
	updates := make(map[string]interface{})
	if article.Title != "" {
		updates["title"] = article.Title
	}
	if article.Content != "" {
		updates["content"] = article.Content
	}
	if article.CategoryID != 0 {
		updates["category_id"] = article.CategoryID
	}

	if err := db.Model(&existingArticle).Updates(updates).Error; err != nil {
		return respcode.ERROR
	}
	return respcode.SUCCESS
}

// DeleteArticle 删除文章
func DeleteArticle(id int) int {
	var article Article
	if err := db.First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return respcode.ErrorArtNotExist
		}
		return respcode.ERROR
	}

	if err := db.Delete(&article).Error; err != nil {
		return respcode.ERROR
	}
	return respcode.SUCCESS
}

// GetArticlesByCategory 获取分类下的文章
func GetArticlesByCategory(categoryID int, pageSize int, pageNum int) ([]Article, int64, int) {
	var articles []Article
	var total int64
	offset := (pageNum - 1) * pageSize

	// 检查分类是否存在
	var category Category
	if err := db.First(&category, categoryID).Error; err != nil {
		return nil, 0, respcode.ErrorCateNotExist
	}

	db.Model(&Article{}).Where("category_id = ?", categoryID).Count(&total)
	if err := db.Preload("Category").Preload("User").
		Where("category_id = ?", categoryID).
		Limit(pageSize).
		Offset(offset).
		Find(&articles).Error; err != nil {
		return nil, 0, respcode.ERROR
	}

	return articles, total, respcode.SUCCESS
}

// GetArticlesByUser 获取用户的文章
func GetArticlesByUser(userID int, pageSize int, pageNum int) ([]Article, int64, int) {
	var articles []Article
	var total int64
	offset := (pageNum - 1) * pageSize

	// 检查用户是否存在
	var user User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, 0, respcode.ErrorUserNotExist
	}

	db.Model(&Article{}).Where("user_id = ?", userID).Count(&total)
	if err := db.Preload("Category").Preload("User").
		Where("user_id = ?", userID).
		Limit(pageSize).
		Offset(offset).
		Find(&articles).Error; err != nil {
		return nil, 0, respcode.ERROR
	}

	return articles, total, respcode.SUCCESS
}

// SearchArticles 搜索文章
func SearchArticles(keyword string, pageSize int, pageNum int) ([]Article, int64, int) {
	var articles []Article
	var total int64
	offset := (pageNum - 1) * pageSize

	if keyword == "" {
		return nil, 0, respcode.BadRequest
	}

	// 使用 LIKE 进行模糊搜索标题
	query := db.Model(&Article{}).Where("title LIKE ?", "%"+keyword+"%")

	// 获取总数
	query.Count(&total)

	// 获取分页数据
	err := query.Preload("Category").Preload("User").
		Limit(pageSize).
		Offset(offset).
		Find(&articles).Error

	if err != nil {
		return nil, 0, respcode.ERROR
	}

	return articles, total, respcode.SUCCESS
}
