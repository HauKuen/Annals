package v1

import (
	"github.com/HauKuen/Annals/model"
	"github.com/HauKuen/Annals/utils/respcode"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetUserInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var maps = make(map[string]interface{})
	data, code := model.GetUser(id)
	maps["username"] = data.Username
	maps["role"] = data.Role
	c.JSON(
		http.StatusOK, gin.H{
			"status":  code,
			"data":    maps,
			"message": respcode.GetErrMsg(code),
		},
	)
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

// CheckUser 用户是否存在
func CheckUser(c *gin.Context) {
	username := c.Query("username")
	code := model.CheckUser(username)
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": respcode.GetErrMsg(code),
	})
}

// AddUser 添加用户
func AddUser(c *gin.Context) {
	var data model.User
	_ = c.ShouldBindJSON(&data)
	code := model.CheckUser(data.Username)
	if code == respcode.SUCCESS {
		model.CreateUser(&data)
	} else {
		code = respcode.ErrorUsernameUsed
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  code,
		"message": respcode.GetErrMsg(code),
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
