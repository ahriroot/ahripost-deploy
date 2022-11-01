package client

import (
	database "ahripost_deploy/database"
	model_v1 "ahripost_deploy/models/v1"

	"github.com/gin-gonic/gin"
)

func Projects(request *gin.Context) {
	var user model_v1.User
	if u, exist := request.Get("user"); exist {
		user = u.(model_v1.User)
	} else {
		request.JSON(200, gin.H{
			"code": 0,
			"msg":  "not login",
			"data": nil,
		})
		return
	}

	var projects []model_v1.Project
	result := database.DB.Where("user_r_id = ?", user.RID).Find(&projects)
	if result.Error != nil {
		request.JSON(200, gin.H{
			"code": 50000,
			"msg":  "get projects error",
			"data": gin.H{
				"message": result.Error.Error(),
			},
		})
		return
	}

	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": gin.H{
			"projects": projects,
		},
	})
}
