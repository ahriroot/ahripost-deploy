package browser

import (
	database "ahripost_deploy/database"
	model_v1 "ahripost_deploy/models/v1"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Item(request *gin.Context) {
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
	var item_id = request.Param("item_id")
	if item_id == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "no item id",
			"data": nil,
		})
		return
	}

	var id int64
	id, err = strconv.ParseInt(item_id, 10, 64)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "item id is not a number",
			"data": nil,
		})
		return
	}

	item := model_v1.Item{}
	result := database.DB.Where("_id = ? AND user_r_id = ?", id, user.RID).First(&item)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no item",
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
		"msg":  "find item success",
		"data": item,
	})
}

func Items(request *gin.Context) {
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

	items := []model_v1.Item{}
	result = database.DB.Where("project_r_id = ?", project.Key).Find(&items)
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
		"msg":  "find items success",
		"data": items,
	})
}

func PostItem(request *gin.Context) {
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
		if member.Status != 1 { // 该成员没有上传权限
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no permission",
				"data": nil,
			})
			return
		}
	}

	var data map[string]interface{}
	err = request.BindJSON(&data)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "json error",
			"data": nil,
		})
		return
	}
	utc_timestame := time.Now().UnixMilli()

	var item model_v1.Item
	item.ID = int64(data["id"].(float64))
	item.Key = data["key"].(string)
	item.Label = data["label"].(string)
	item.Type = data["type"].(string)
	item.From = data["from"].(string)
	item.ProjectRID = project.Key
	item.UserRID = user.RID
	item.Parent = data["parent"].(string)
	item.LastSync = utc_timestame
	item.LastUpdate = int64(data["last_update"].(float64))
	item.Template = data["template"].(string)
	if data["request"] == nil {
		item.Request = ""
	} else {
		item.Request = data["request"].(string)
	}
	if data["response"] == nil {
		item.Response = ""
	} else {
		item.Response = data["response"].(string)
	}

	result = database.DB.Create(&item)
	if result.Error != nil {
		request.JSON(200, gin.H{
			"code": 50000,
			"msg":  "server error",
			"data": gin.H{
				"message": result,
			},
		})
		return
	}

	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "create item success",
		"data": item,
	})
}
