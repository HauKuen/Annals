package v1

import (
	"github.com/HauKuen/Annals/model"
	"github.com/HauKuen/Annals/utils/respcode"
	"github.com/gin-gonic/gin"
	"net/http"
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
			"status":  respcode.ERROR,
			"message": respcode.GetErrMsg(code),
		})
		return
	}

	// 创建新分类
	code = model.CreateCategory(&data)

	if code == respcode.SUCCESS {
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": respcode.GetErrMsg(code),
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  code,
			"message": respcode.GetErrMsg(code),
		})
	}
}
