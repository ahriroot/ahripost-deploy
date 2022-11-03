package middleware

import (
	database "ahripost_deploy/database"
	model_v1 "ahripost_deploy/models/v1"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TokenLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		token_string := c.GetHeader("Authorization")

		token := model_v1.Token{}
		result := database.DB.Where("token = ?", token_string).First(&token)
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

		user := model_v1.User{}
		result = database.DB.Where("_id = ?", token.UserRID).First(&user)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(200, gin.H{
					"code": 0,
					"msg":  "no user",
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

		c.Set("user", user)
		c.Next()
	}
}
