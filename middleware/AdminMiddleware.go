package middleware

import (
	database "ahripost_deploy/database"
	model_v1 "ahripost_deploy/models/v1"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token_string := c.GetHeader("Authorization")

		user := model_v1.User{}
		result := database.DB.Where("token = ?", token_string).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(200, gin.H{
					"code": 0,
					"msg":  "error token",
					"data": nil,
				})
				c.Abort()
				return
			} else {
				c.JSON(200, gin.H{
					"code": 50000,
					"msg":  "server error",
					"data": gin.H{
						"message": result.Error.Error(),
					},
				})
				c.Abort()
				return
			}
		}

		if user.Username != "admin" {
			c.JSON(200, gin.H{
				"code": 0,
				"msg":  "not admin",
				"data": nil,
			})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
