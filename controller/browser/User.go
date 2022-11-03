package browser

import (
	database "ahripost_deploy/database"
	model_v1 "ahripost_deploy/models/v1"
	"ahripost_deploy/tools"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UserInfo(request *gin.Context) {
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
	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": gin.H{
			"user": user,
		},
	})
}

func User(request *gin.Context) {
	var err error
	var user_id = request.Param("user_id")
	if user_id == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "no user id",
			"data": nil,
		})
		return
	}

	var id int64
	id, err = strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "user id is not a number",
			"data": nil,
		})
		return
	}

	u := model_v1.User{}
	result := database.DB.Where("_id = ?", id).First(&u)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no user",
				"data": nil,
			})
			return
		} else {
			request.JSON(200, gin.H{
				"code": 50000,
				"msg":  "database error",
				"data": nil,
			})
			return
		}
	}

	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": gin.H{
			"user": u,
		},
	})
}

func Users(request *gin.Context) {
	var err error
	var page = request.Query("page")
	var page_size = request.Query("size")

	var page_int int
	var page_size_int int

	if page == "" {
		page_int = 1
	} else {
		page_int, err = strconv.Atoi(page)
		if err != nil {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "page is not a number",
				"data": nil,
			})
			return
		}
	}

	if page_size == "" {
		page_size_int = 10
	} else {
		page_size_int, err = strconv.Atoi(page_size)
		if err != nil {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "page_size is not a number",
				"data": nil,
			})
			return
		}
	}

	var count int64 // 总数
	var users []model_v1.User
	result := database.DB.Model(&model_v1.User{}).Count(&count)
	if result.Error != nil {
		request.JSON(200, gin.H{
			"code": 50000,
			"msg":  "database error",
			"data": nil,
		})
		return
	}
	result = database.DB.Limit(page_size_int).Offset((page_int - 1) * page_size_int).Find(&users)
	if result.Error != nil {
		request.JSON(200, gin.H{
			"code": 50000,
			"msg":  "database error",
			"data": nil,
		})
		return
	}
	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": gin.H{
			"count": count,
			"users": users,
		},
	})
}

type FormUser struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func PostUser(request *gin.Context) {
	var err error
	var data FormUser
	err = request.ShouldBindJSON(&data)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "json error",
			"data": nil,
		})
		return
	}

	if data.Username == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "username is empty",
			"data": nil,
		})
		return
	}
	if data.Username == "admin" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "username con not be admin",
			"data": nil,
		})
		return
	}
	if data.Password == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "password is empty",
			"data": nil,
		})
		return
	}

	var user model_v1.User
	result := database.DB.Where("username = ?", data.Username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			user.Username = data.Username
			user.Password = tools.Sha256(data.Password)
			if database.DB.Create(&user).Error != nil {
				request.JSON(200, gin.H{
					"code": 50000,
					"msg":  "database error",
					"data": nil,
				})
				return
			}
			request.JSON(200, gin.H{
				"code": 10000,
				"msg":  "success",
				"data": gin.H{
					"user": user,
				},
			})
			return
		} else {
			request.JSON(200, gin.H{
				"code": 50000,
				"msg":  "database error",
				"data": gin.H{
					"message": result.Error.Error(),
				},
			})
			return
		}
	} else {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "username is arleady exist",
			"data": nil,
		})
		return
	}
}

func PutUser(request *gin.Context) {
	var err error
	var user_id = request.Param("user_id")
	if user_id == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "no user id",
			"data": nil,
		})
		return
	}

	var id int64
	id, err = strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "user id is not a number",
			"data": nil,
		})
		return
	}

	var data FormUser
	err = request.ShouldBindJSON(&data)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "json error",
			"data": nil,
		})
		return
	}

	if data.Username == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "username is empty",
			"data": nil,
		})
		return
	}

	if data.Password == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "password is empty",
			"data": nil,
		})
		return
	}

	var user model_v1.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			request.JSON(200, gin.H{
				"code": 50000,
				"msg":  "database error",
				"data": nil,
			})
			return
		} else {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "user is not exists",
				"data": nil,
			})
			return
		}
	}

	if user.Username == "admin" && data.Username != "admin" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "admin can not be modified",
			"data": nil,
		})
		return
	}

	user.Username = data.Username
	user.Password = tools.Sha256(data.Password)
	user.Token = ""
	result = database.DB.Save(&user)
	if result.Error != nil {
		request.JSON(200, gin.H{
			"code": 50000,
			"msg":  "database error",
			"data": nil,
		})
		return
	}

	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "success",
		"data": gin.H{
			"user": user,
		},
	})
}

func DeleteUser(request *gin.Context) {
	var err error
	var user_id = request.Param("user_id")
	if user_id == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "no user id",
			"data": nil,
		})
		return
	}

	var id int64
	id, err = strconv.ParseInt(user_id, 10, 64)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "user id is not a number",
			"data": nil,
		})
		return
	}

	var user model_v1.User
	result := database.DB.First(&user, id)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "user is not exists",
				"data": nil,
			})
			return
		} else {
			request.JSON(200, gin.H{
				"code": 50000,
				"msg":  "database error",
				"data": gin.H{
					"message": result.Error.Error(),
				},
			})
			return
		}
	} else {
		if user.Username == "admin" {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "admin can not delete",
				"data": nil,
			})
			return
		}
		result = database.DB.Delete(&user)
		if result.Error != nil {
			request.JSON(200, gin.H{
				"code": 50000,
				"msg":  "database error",
				"data": nil,
			})
			return
		}
		request.JSON(200, gin.H{
			"code": 10000,
			"msg":  "success",
			"data": nil,
		})
		return
	}
}
