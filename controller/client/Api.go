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
	result := database.DB.Where("key = ?", data_project["key"].(string)).Find(&project)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no project",
				"data": gin.H{
					"message": result.Error.Error(),
				},
			})
			return
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

	if project.UserRID != user.RID { // 该用户不拥有该项目
		member := model_v1.Member{}
		result = database.DB.Where("project_r_id = ? AND member_r_id = ?", project.Key, user.RID).First(&member)
		if result.Error != nil { // 该用户不是该项目的成员
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no project",
				"data": nil,
			})
			return
		}
		if member.Status != 1 { // 该成员没有上传权限
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no permission",
				"data": nil,
			})
			return
		}
	}

	data_apis := data["apis"].([]interface{})
	keys := []string{}
	for _, api := range data_apis {
		api := api.(map[string]interface{})
		keys = append(keys, api["key"].(string))
	}

	var db_apis []model_v1.Item
	result = database.DB.Where("project_r_id = ? AND key IN ?", project.Key, keys).Find(&db_apis)
	if result.Error != nil {
		request.JSON(200, gin.H{
			"code": 50000,
			"msg":  "sync apis error",
			"data": gin.H{
				"message": result.Error.Error(),
			},
		})
		return
	}

	data_remote := make(map[string]interface{})
	for _, api := range db_apis {
		data_remote[api.Key] = api
	}

	items_upload := make([]model_v1.Item, 0) // 上传到远程的数据

	items_update := make([]model_v1.Item, 0) // 更新远程
	items_sync := make([]model_v1.Item, 0)   // 更新本地

	// for data_local
	utc_timestame := time.Now().UnixMilli()
	for _, api := range data_apis {
		api := api.(map[string]interface{})
		key := api["key"].(string)
		if _, exist := data_remote[key]; exist {
			item := data_remote[key].(model_v1.Item)
			client := api["client"].(string)

			last_update_local := int64(api["last_update"].(float64))
			if user.RID == item.UserRID && client == item.Client { // 上次是本人在同一台设备上修改的
				item.Key = api["key"].(string)
				item.Label = api["label"].(string)
				item.Type = api["type"].(string)
				item.Parent = api["parent"].(string)
				item.UserRID = user.RID
				item.LastSync = utc_timestame
				item.LastUpdate = last_update_local
				item.Request = api["request"].(string)
				item.Response = api["response"].(string)
				item.Template = api["template"].(string)
				item.Client = client
				items_update = append(items_update, item)
			} else {
				// last_sync_remote := data_remote[key].(model_v1.Item).LastSync
				// last_sync_local := int64(api["last_sync"].(float64))
				last_update_remote := data_remote[key].(model_v1.Item).LastUpdate
				if last_update_local > last_update_remote { // 本地数据更新, 上传到远程
					item.Key = api["key"].(string)
					item.Label = api["label"].(string)
					item.Type = api["type"].(string)
					item.Parent = api["parent"].(string)
					item.UserRID = user.RID
					item.LastSync = utc_timestame
					item.LastUpdate = last_update_local
					item.Request = api["request"].(string)
					item.Response = api["response"].(string)
					item.Template = api["template"].(string)
					item.Client = client
					items_update = append(items_update, item)
				} else if last_update_local < last_update_remote {
					items_sync = append(items_sync, item)
				}
			}
		} else {
			items_upload = append(items_upload, model_v1.Item{
				ID:         int64(api["id"].(float64)),
				Key:        api["key"].(string),
				Label:      api["label"].(string),
				Type:       api["type"].(string),
				From:       "client",
				ProjectRID: project.Key,
				UserRID:    user.RID,
				Parent:     api["parent"].(string),
				LastSync:   utc_timestame,
				LastUpdate: int64(api["last_update"].(float64)),
				Request:    api["request"].(string),
				Response:   api["response"].(string),
				Template:   api["template"].(string),
				Client:     api["client"].(string),
			})
		}
	}

	for _, api := range items_upload {
		// create
		database.DB.Create(&api)
	}

	for _, api := range items_update {
		// update
		database.DB.Save(&api)
	}

	request.JSON(200, gin.H{
		"code": 10001,
		"msg":  "sync success",
		"data": gin.H{
			"items_sync":   items_sync,
			"items_update": items_update,
		},
	})
}

func PutApi(request *gin.Context) {
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

	data_project := data["project"].(string)
	project := model_v1.Project{}
	result := database.DB.Where("key = ?", data_project).Find(&project)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no project",
				"data": gin.H{
					"message": result.Error.Error(),
				},
			})
			return
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

	if project.UserRID != user.RID { // 该用户不拥有该项目
		member := model_v1.Member{}
		result = database.DB.Where("project_r_id = ? AND member_r_id = ?", project.Key, user.RID).First(&member)
		if result.Error != nil { // 该用户不是该项目的成员
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no project",
				"data": nil,
			})
			return
		}
		if member.Status != 1 { // 该成员没有上传权限
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no permission",
				"data": nil,
			})
			return
		}
	}

	data_keys := data["keys"].([]interface{})

	result = database.DB.Model(&model_v1.Item{}).Where("project_r_id = ? AND key IN (?)", project.Key, data_keys).Update("tag", true)
	if result.Error != nil {
		request.JSON(200, gin.H{
			"code": 50000,
			"msg":  "delete error",
			"data": gin.H{
				"message": result.Error.Error(),
			},
		})
		return
	}

	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "delete success",
		"data": data_keys,
	})
}
