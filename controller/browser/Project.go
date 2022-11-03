package browser

import (
	database "ahripost_deploy/database"
	model_v1 "ahripost_deploy/models/v1"

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

	var project_id = request.Param("project_id")
	if project_id == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "no project id",
			"data": nil,
		})
		return
	}

	project := model_v1.Project{}
	result := database.DB.Where("key = ?", project_id).First(&project)
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

	if project.UserRID != user.RID {
		member := model_v1.Member{}
		result = database.DB.Where("project_r_id = ? AND member_r_id = ?", project.Key, user.RID).First(&member)
		if result.Error != nil {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no project",
				"data": nil,
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

	members := []model_v1.Member{}
	result = database.DB.Where("member_r_id = ?", user.RID).Find(&members)
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
	project_ids := []string{}
	for _, member := range members {
		project_ids = append(project_ids, member.ProjectRID)
	}

	projects_by_member := []model_v1.Project{}
	result = database.DB.Where("key IN ?", project_ids).Find(&projects_by_member)
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
		"data": gin.H{
			"projects":           projects,
			"projects_by_member": projects_by_member,
		},
	})
}

type FormProject struct {
	Name   string `json:"name"`
	Public bool   `json:"public"`
}

func PutProject(request *gin.Context) {
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

	var project_id = request.Param("project_id")
	if project_id == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "no project id",
			"data": nil,
		})
		return
	}

	var data FormProject
	if err := request.ShouldBindJSON(&data); err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "no project id",
			"data": nil,
		})
		return
	}

	project, err := CheckPermission(user, project_id, []int{1})
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "no permission",
			"data": nil,
		})
		return
	}
	if data.Name != "" {
		project.Name = data.Name
	}
	project.Public = data.Public
	result := database.DB.Save(&project)
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
		"msg":  "update project success",
		"data": project,
	})
}
