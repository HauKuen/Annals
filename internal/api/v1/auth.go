package v1

import (
	"net/http"
	"strings"

	"github.com/HauKuen/Annals/internal/model"
	"github.com/HauKuen/Annals/internal/utils/jwt"
	"github.com/HauKuen/Annals/internal/utils/respcode"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  respcode.BadRequest,
			"message": respcode.GetErrMsg(respcode.BadRequest),
		})
		return
	}

	user, code := model.GetUserByUsername(req.Username)
	if code != respcode.SUCCESS {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  code,
			"message": respcode.GetErrMsg(code),
		})
		return
	}

	if !user.VerifyPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  respcode.ErrorPasswordWrong,
			"message": respcode.GetErrMsg(respcode.ErrorPasswordWrong),
		})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  respcode.ErrorUserInactive,
			"message": respcode.GetErrMsg(respcode.ErrorUserInactive),
		})
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  respcode.ERROR,
			"message": respcode.GetErrMsg(respcode.ERROR),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  respcode.SUCCESS,
		"message": respcode.GetErrMsg(respcode.SUCCESS),
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}

// ValidateToken 验证token的有效性
func ValidateToken(c *gin.Context) {
	auth := c.GetHeader("Authorization")
	if auth == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  respcode.Unauthorized,
			"message": respcode.GetErrMsg(respcode.Unauthorized),
			"valid":   false,
		})
		return
	}

	// 分割Bearer token
	parts := strings.SplitN(auth, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  respcode.ErrorTokenInvalid,
			"message": respcode.GetErrMsg(respcode.ErrorTokenInvalid),
			"valid":   false,
		})
		return
	}

	// 解析token
	claims, err := jwt.ParseToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  respcode.ErrorTokenInvalid,
			"message": respcode.GetErrMsg(respcode.ErrorTokenInvalid),
			"valid":   false,
		})
		return
	}

	// 检查用户是否存在且处于活动状态
	user, code := model.GetUserByUsername(claims.Username)
	if code != respcode.SUCCESS {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  code,
			"message": respcode.GetErrMsg(code),
			"valid":   false,
		})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  respcode.ErrorUserInactive,
			"message": respcode.GetErrMsg(respcode.ErrorUserInactive),
			"valid":   false,
		})
		return
	}

	// Token 有效，返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"status":  respcode.SUCCESS,
		"message": respcode.GetErrMsg(respcode.SUCCESS),
		"valid":   true,
		"data": gin.H{
			"id":       claims.ID,
			"username": claims.Username,
			"role":     claims.Role,
		},
	})
}
