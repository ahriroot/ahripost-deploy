package client

import (
	"time"

	database "ahripost_deploy/database"
	model_v1 "ahripost_deploy/models/v1"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Apis(request *gin.Context) {
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
	var data map[string]interface{}
	if err = request.ShouldBind(&data); err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "参数错误！",
			"data": gin.H{
				"message": err.Error(),
			},
		})
		return
	}
	data_project := data["project"].(map[string]interface{})
	project := model_v1.Project{}
	result := database.DB.First(&project, int64(data_project["_id"].(float64)))
	utc_timestame := time.Now().UnixMilli()
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			project.ID = int64(data_project["id"].(float64))
			project.UserRID = user.RID
			project.Key = data_project["key"].(string)
			project.Name = data_project["name"].(string)
			project.CreateAt = utc_timestame
			database.DB.Create(&project)
		} else {
			request.JSON(200, gin.H{
				"code": 50000,
				"msg":  "sync project error",
				"data": gin.H{
					"message": result.Error.Error(),
				},
			})
			return
		}
	}

	data_item := data["item"].(map[string]interface{})
	item := model_v1.Item{}
	result = database.DB.Where("key = ? AND project_r_id = ?", data_item["key"].(string), project.RID).First(&item)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			item.ID = int64(data_item["id"].(float64))
			item.Key = data_item["key"].(string)
			item.Label = data_item["label"].(string)
			item.Type = data_item["type"].(string)
			item.ProjectRID = project.RID
			item.UserRID = user.RID
			item.Parent = data_item["parent"].(string)
			item.LastSync = utc_timestame
			item.LastUpdate = int64(data_item["last_update"].(float64))
			item.Request = data_item["request"].(string)
			item.Response = data_item["response"].(string)
			database.DB.Create(&item)

			request.JSON(200, gin.H{
				"code": 10000,
				"msg":  "sync success",
				"data": gin.H{
					"project": project,
					"item":    item,
				},
			})
			return
		} else {
			request.JSON(200, gin.H{
				"code": 50000,
				"msg":  "sync api error",
				"data": gin.H{
					"message": result.Error.Error(),
				},
			})
			return
		}
	} else {
		local_last_sync := int64(data_item["last_sync"].(float64))
		remote_last_sync := item.LastSync
		// local_last_update := int64(data_item["last_update"].(float64))
		// remote_last_update := item.LastUpdate

		// 上次同步时间早于上次更新时间，说明有其他人更新了数据，需要同步到本地
		if local_last_sync < remote_last_sync {
			request.JSON(200, gin.H{
				"code": 10002,
				"msg":  "sync conflict",
				"data": gin.H{
					"project": project,
					"item":    item,
				},
			})
			return
		} else {
			item.Key = data_item["key"].(string)
			item.Label = data_item["label"].(string)
			item.Type = data_item["type"].(string)
			item.Parent = data_item["parent"].(string)
			item.UserRID = user.RID
			item.LastSync = utc_timestame
			item.LastUpdate = int64(data_item["last_update"].(float64))
			item.Request = data_item["request"].(string)
			item.Response = data_item["response"].(string)
			database.DB.Save(&item)

			request.JSON(200, gin.H{
				"code": 10001,
				"msg":  "sync success",
				"data": gin.H{
					"project": project,
					"item":    item,
				},
			})
			return
		}
	}
}
