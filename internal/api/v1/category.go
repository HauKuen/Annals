package v1

import (
	"net/http"
	"strconv"

	"github.com/HauKuen/Annals/internal/model"
	"github.com/HauKuen/Annals/internal/utils/respcode"
	"github.com/gin-gonic/gin"
)

// AddCategory 添加分类
func AddCategory(c *gin.Context) {
	var data model.Category

	// 绑定 JSON 数据到结构体
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.ERROR,
			"message": "无效的数据格式",
		})
		return
	}

	// 检查分类名是否为空
	if data.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.ErrorEmptyCatename,
			"message": respcode.GetErrMsg(respcode.ErrorEmptyCatename),
		})
		return
	}

	// 检查分类是否已存在
	code := model.CheckCategory(data.Name)
	if code == respcode.SUCCESS {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.ErrorCatenameUsed,
			"message": respcode.GetErrMsg(respcode.ErrorCatenameUsed),
		})
		return
	}

	// 创建新分类
	code = model.CreateCategory(&data)

	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": respcode.GetErrMsg(code),
	})
}

// DeleteCategory 删除分类
func DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DeleteCategory(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": respcode.GetErrMsg(code),
	})
}

// GetCategory 获取分类
func GetCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	data, code := model.GetCategory(id)
	response := gin.H{
		"status":  code,
		"message": respcode.GetErrMsg(code),
	}
	if code == respcode.SUCCESS {
		response["data"] = data
	}
	c.JSON(http.StatusOK, response)
}
