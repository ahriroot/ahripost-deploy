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
	Project    string `json:"project"`
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
	result := database.DB.Where("key = ?", data.Project.Key).First(&project)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			project.ID = 0
			project.UserRID = user.RID
			project.Key = data.Project.Key
			project.Name = data.Project.Name
			project.CreateAt = time.Now().UnixMilli()
			result = database.DB.Create(&project)
			if result.Error != nil {
				request.JSON(200, gin.H{
					"code": 50000,
					"msg":  "sync project error",
					"data": gin.H{
						"message": result.Error.Error(),
					},
				})
				return
			}
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

	items := []model_v1.Item{}
	database.DB.Where("project_r_id = ?", project.Key).Find(&items)
	map_remote := map[string]model_v1.Item{}
	for _, item := range items {
		map_remote[item.Key] = item
	}

	keys_delete := make([]string, 0)
	items_download := make([]string, 0)
	items_upload := make([]string, 0)

	for _, api := range data.Apis {
		if item, exist := map_remote[api.Key]; exist {
			if map_remote[api.Key].Tag {
				keys_delete = append(keys_delete, api.Key)
			} else if item.LastUpdate > api.LastUpdate {
				items_download = append(items_download, item.Key)
			} else if item.LastUpdate < api.LastUpdate {
				items_upload = append(items_upload, item.Key)
			}
			delete(map_remote, api.Key)
		} else {
			items_upload = append(items_upload, api.Key)
		}
	}
	for _, item := range map_remote {
		if !item.Tag {
			items_download = append(items_download, item.Key)
		}
	}

	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "sync check success",
		"data": gin.H{
			"items_download": items_download,
			"items_upload":   items_upload,
			"keys_delete":    keys_delete,
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
	result := database.DB.Where("key = ?", project_key).First(&project)
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

	items := []model_v1.Item{}
	database.DB.Where("project_r_id = ? AND key IN ?", project.Key, items_download).Find(&items)

	utc_timestame := time.Now().UnixMilli()
	count := 0

	for _, api := range items_upload {
		data_item := api.(map[string]interface{})
		item := model_v1.Item{}
		result := database.DB.Where("key = ?", data_item["key"]).First(&item)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				item.ID = int64(data_item["id"].(float64))
				item.Key = data_item["key"].(string)
				item.Label = data_item["label"].(string)
				item.Type = data_item["type"].(string)
				item.ProjectRID = project.Key
				item.UserRID = user.RID
				item.Parent = data_item["parent"].(string)
				item.LastSync = utc_timestame
				item.LastUpdate = int64(data_item["last_update"].(float64))
				item.Request = data_item["request"].(string)
				item.Response = data_item["response"].(string)
				item.Template = data_item["template"].(string)
				item.Client = data_item["client"].(string)
				item.Tag = data_item["tag"].(bool)
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
			item.Parent = data_item["parent"].(string)
			item.UserRID = user.RID
			item.LastSync = utc_timestame
			item.LastUpdate = int64(data_item["last_update"].(float64))
			item.Request = data_item["request"].(string)
			item.Response = data_item["response"].(string)
			item.Template = data_item["template"].(string)
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
