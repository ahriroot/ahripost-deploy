package client

import (
	database "ahripost_deploy/database"
	model_v1 "ahripost_deploy/models/v1"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FormApi struct {
	Key        string `json:"key"`
	Project    int64  `json:"project"`
	LastSync   int64  `json:"last_sync"`
	LastUpdate int64  `json:"last_update"`
}

type FormProject struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type FormRequest struct {
	Apis    []FormApi   `json:"apis"`
	Project FormProject `json:"project"`
}

func SyncCheck(request *gin.Context) {
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
	var data FormRequest
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

	project := model_v1.Project{}
	result := database.DB.Where("key = ? AND user_r_id = ?", data.Project.Key, user.RID).First(&project)
	if result.Error != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "项目不存在！",
			"data": nil,
		})
		return
	}

	items := []model_v1.Item{}
	database.DB.Where("user_r_id = ? AND project_r_id = ?", user.RID, project.RID).Find(&items)
	map_remote := map[string]model_v1.Item{}
	for _, item := range items {
		map_remote[item.Key] = item
	}

	items_download := make([]string, 0)
	items_upload := make([]string, 0)
	for _, api := range data.Apis {
		if item, exist := map_remote[api.Key]; exist {
			if item.LastUpdate > api.LastUpdate {
				items_download = append(items_download, item.Key)
			} else if item.LastUpdate < api.LastUpdate {
				items_upload = append(items_upload, item.Type)
			}
			delete(map_remote, api.Key)
		} else {
			items_upload = append(items_upload, api.Key)
		}
	}
	for _, item := range map_remote {
		items_download = append(items_download, item.Key)
	}

	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "sync check success",
		"data": gin.H{
			"items_download": items_download,
			"items_upload":   items_upload,
		},
	})
}

func SyncData(request *gin.Context) {
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

	items_upload := data["items_upload"].([]interface{})
	items_download := data["items_download"].([]interface{})
	items_project := data["project"].(map[string]interface{})
	project_key := items_project["key"]
	project := model_v1.Project{}
	result := database.DB.Where("key = ? AND user_r_id = ?", project_key, user.RID).First(&project)
	if result.Error != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "项目不存在！",
			"data": nil,
		})
		return
	}

	items := []model_v1.Item{}
	database.DB.Where("user_r_id = ? AND key IN ?", user.RID, items_download).Find(&items)

	utc_timestame := time.Now().UnixMilli()
	count := 0
	for _, api := range items_upload {
		data_item := api.(map[string]interface{})
		item := model_v1.Item{}
		result := database.DB.Where("user_r_id = ? AND key = ?", user.RID, data_item["key"]).First(&item)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				item.ID = int64(data_item["id"].(float64))
				item.Key = data_item["key"].(string)
				item.Label = data_item["label"].(string)
				item.Type = data_item["type"].(string)
				item.ProjectRID = project.RID
				item.UserRID = user.RID
				item.Parent = int64(data_item["parent"].(float64))
				item.LastSync = utc_timestame
				item.LastUpdate = int64(data_item["last_update"].(float64))
				item.Request = data_item["request"].(string)
				item.Response = data_item["response"].(string)
				database.DB.Create(&item)

				count++
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
			item.Key = data_item["key"].(string)
			item.Label = data_item["label"].(string)
			item.Type = data_item["type"].(string)
			item.Parent = int64(data_item["parent"].(float64))
			item.UserRID = user.RID
			item.LastSync = utc_timestame
			item.LastUpdate = int64(data_item["last_update"].(float64))
			item.Request = data_item["request"].(string)
			item.Response = data_item["response"].(string)
			database.DB.Save(&item)

			count++
		}
	}

	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "sync data success",
		"data": items,
	})
}
