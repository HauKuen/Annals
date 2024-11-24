package v1

import (
	"net/http"

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
