package middleware

import (
	"strings"

	"github.com/HauKuen/Annals/internal/utils/jwt"
	"github.com/HauKuen/Annals/internal/utils/respcode"
	"github.com/gin-gonic/gin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(401, gin.H{
				"status":  respcode.Unauthorized,
				"message": respcode.GetErrMsg(respcode.Unauthorized),
			})
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(401, gin.H{
				"status":  respcode.ErrorTokenInvalid,
				"message": respcode.GetErrMsg(respcode.ErrorTokenInvalid),
			})
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(401, gin.H{
				"status":  respcode.ErrorTokenInvalid,
				"message": respcode.GetErrMsg(respcode.ErrorTokenInvalid),
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.ID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}
