package v1

import (
	"net/http"
	"strconv"

	"github.com/HauKuen/Annals/internal/model"
	"github.com/HauKuen/Annals/internal/utils/respcode"
	"github.com/gin-gonic/gin"
)

// GetArticles 获取文章列表
func GetArticles(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))

	data, total := model.GetArticles(pageSize, pageNum)
	code := respcode.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": respcode.GetErrMsg(code),
	})
}

// GetArticle 获取单个文章
func GetArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.BadRequest,
			"message": respcode.GetErrMsg(respcode.BadRequest),
		})
		return
	}

	data, code := model.GetArticleByID(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"message": respcode.GetErrMsg(code),
	})
}

// AddArticle 添加文章
func AddArticle(c *gin.Context) {
	var article model.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.BadRequest,
			"message": respcode.GetErrMsg(respcode.BadRequest),
		})
		return
	}

	// 设置当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  respcode.Unauthorized,
			"message": respcode.GetErrMsg(respcode.Unauthorized),
		})
		return
	}
	article.UserID = userID.(uint)

	code := model.CreateArticle(&article)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": respcode.GetErrMsg(code),
		"data":    article,
	})
}

// EditArticle 更新文章
func EditArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.BadRequest,
			"message": respcode.GetErrMsg(respcode.BadRequest),
		})
		return
	}

	var article model.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.BadRequest,
			"message": respcode.GetErrMsg(respcode.BadRequest),
		})
		return
	}

	// 检查是否是文章作者
	userID := c.GetUint("user_id")
	existingArticle, code := model.GetArticleByID(id)
	if code != respcode.SUCCESS {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  code,
			"message": respcode.GetErrMsg(code),
		})
		return
	}

	if existingArticle.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  respcode.ErrorNoPermission,
			"message": respcode.GetErrMsg(respcode.ErrorNoPermission),
		})
		return
	}

	code = model.UpdateArticle(id, &article)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": respcode.GetErrMsg(code),
	})
}

// DeleteArticle 删除文章
func DeleteArticle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.BadRequest,
			"message": respcode.GetErrMsg(respcode.BadRequest),
		})
		return
	}

	// 检查是否是文章作者
	userID := c.GetUint("user_id")
	article, code := model.GetArticleByID(id)
	if code != respcode.SUCCESS {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  code,
			"message": respcode.GetErrMsg(code),
		})
		return
	}

	if article.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  respcode.ErrorNoPermission,
			"message": respcode.GetErrMsg(respcode.ErrorNoPermission),
		})
		return
	}

	code = model.DeleteArticle(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": respcode.GetErrMsg(code),
	})
}

// GetCategoryArticles 获取分类下的文章
func GetCategoryArticles(c *gin.Context) {
	categoryID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.BadRequest,
			"message": respcode.GetErrMsg(respcode.BadRequest),
		})
		return
	}

	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))

	data, total, code := model.GetArticlesByCategory(categoryID, pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": respcode.GetErrMsg(code),
	})
}

// GetUserArticles 获取用户的文章
func GetUserArticles(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.BadRequest,
			"message": respcode.GetErrMsg(respcode.BadRequest),
		})
		return
	}

	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))

	data, total, code := model.GetArticlesByUser(userID, pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": respcode.GetErrMsg(code),
	})
}

// SearchArticles 搜索文章
func SearchArticles(c *gin.Context) {
	keyword := c.Query("keyword")
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))

	if pageSize <= 0 {
		pageSize = 10
	}
	if pageNum <= 0 {
		pageNum = 1
	}

	data, total, code := model.SearchArticles(keyword, pageSize, pageNum)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": respcode.GetErrMsg(code),
	})
}
