package browser

import (
	database "ahripost_deploy/database"
	model_v1 "ahripost_deploy/models/v1"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Project(request *gin.Context) {
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

	var err error
	var project_id = request.Param("project_id")
	if project_id == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "no project id",
			"data": nil,
		})
		return
	}

	var id int64
	id, err = strconv.ParseInt(project_id, 10, 64)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "project id is not a number",
			"data": nil,
		})
		return
	}

	project := model_v1.Project{}
	result := database.DB.Where("_id = ? AND user_r_id = ?", id, user.RID).First(&project)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no project",
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

	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "find project success",
		"data": project,
	})
}

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

	projects := []model_v1.Project{}
	result := database.DB.Where("user_r_id = ?", user.RID).Find(&projects)
	if result.Error != nil {
		request.JSON(200, gin.H{
			"code": 50000,
			"msg":  "server error",
			"data": gin.H{
				"message": result.Error.Error(),
			},
		})
		return
	}

	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "find projects success",
		"data": projects,
	})
}
