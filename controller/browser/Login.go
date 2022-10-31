package browser

import (
	database "ahripost_deploy/database"
	model_v1 "ahripost_deploy/models/v1"
	"ahripost_deploy/tools"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(request *gin.Context) {
	var data LoginForm
	if err := request.ShouldBindJSON(&data); err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "data error",
			"data": nil,
		})
		return
	}

	password := tools.Sha256(data.Password)

	user := model_v1.User{}
	result := database.DB.Where("username = ? AND password = ?", data.Username, password).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			request.JSON(200, gin.H{
				"code": 0,
				"msg":  "no user",
				"data": nil,
			})
			return
		} else {
			request.JSON(200, gin.H{
				"code": 50000,
				"msg":  "server error",
				"data": gin.H{
					"message": result.Error.Error(),
				},
			})
			return
		}
	}

	token := tools.Sha256(tools.RandomString(8) + data.Username)
	user.Token = token
	database.DB.Save(&user)

	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "login success",
		"data": gin.H{
			"token": token,
		},
	})
}
