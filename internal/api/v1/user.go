package v1

import (
	"net/http"
	"strconv"

	"github.com/HauKuen/Annals/internal/model"
	"github.com/HauKuen/Annals/internal/utils/respcode"
	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.BadRequest,
			"message": respcode.GetErrMsg(respcode.BadRequest),
		})
		return
	}

	data, code := model.GetUser(id)
	response := gin.H{
		"status":  code,
		"data":    data,
		"message": respcode.GetErrMsg(code),
	}

	if code == respcode.SUCCESS {
	}

	c.JSON(http.StatusOK, response)
}

func GetUsers(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.Query("pageSize"))
	pageNum, _ := strconv.Atoi(c.Query("pageNum"))
	if pageSize <= 0 {
		pageSize = 10 // 默认页面大小
	}
	if pageNum <= 0 {
		pageNum = 1 // 默认页码为第一页
	}

	data, total := model.GetUsers(pageSize, pageNum)
	code := respcode.SUCCESS
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"data":    data,
		"total":   total,
		"message": respcode.GetErrMsg(code),
	})
}

// AddUser 添加用户
func AddUser(c *gin.Context) {
	var data model.User

	// 绑定 JSON 数据并检查错误
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.BadRequest,
			"message": respcode.GetErrMsg(respcode.BadRequest),
		})
		return
	}

	// 创建用户
	if code := model.CreateUser(&data); code != respcode.SUCCESS {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  code,
			"message": respcode.GetErrMsg(code),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  respcode.SUCCESS,
		"message": respcode.GetErrMsg(respcode.SUCCESS),
		"data":    data,
	})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	code := model.DeleteUser(id)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": respcode.GetErrMsg(code),
	})
}

// EditUser 更新用户信息
func EditUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.BadRequest,
			"message": respcode.GetErrMsg(respcode.BadRequest),
		})
		return
	}

	var data model.User
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.ERROR,
			"message": "Invalid input data",
		})
		return
	}

	// 检查用户名是否已被其他用户使用
	if data.Username != "" {
		code := model.CheckUser(data.Username)
		if code == respcode.ErrorUsernameUsed {
			c.JSON(http.StatusConflict, gin.H{
				"status":  respcode.ErrorUsernameUsed,
				"message": respcode.GetErrMsg(code),
			})
			return
		}
	}

	code := model.EditUser(id, &data)

	switch code {
	case respcode.SUCCESS:
		c.JSON(http.StatusOK, gin.H{
			"status":  code,
			"message": respcode.GetErrMsg(code),
			"data":    data,
		})
	case respcode.ErrorUserNotExist:
		c.JSON(http.StatusNotFound, gin.H{
			"status":  code,
			"message": respcode.GetErrMsg(code),
		})
	case respcode.ErrorUsernameUsed, respcode.ERROR:
		c.JSON(http.StatusConflict, gin.H{
			"status":  code,
			"message": respcode.GetErrMsg(code),
		})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  code,
			"message": respcode.GetErrMsg(code),
		})
	}
}
