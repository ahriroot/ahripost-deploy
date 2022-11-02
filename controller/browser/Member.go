package browser

import (
	database "ahripost_deploy/database"
	model_v1 "ahripost_deploy/models/v1"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Member(request *gin.Context) {
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
	var member_id = request.Param("member_id")
	if member_id == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "no member id",
			"data": nil,
		})
		return
	}

	var id int64
	id, err = strconv.ParseInt(member_id, 10, 64)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "member id is not a number",
			"data": nil,
		})
		return
	}

	member := model_v1.Member{}
	result := database.DB.Where("_id = ? AND user_r_id = ?", id, user.RID).First(&member)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no member",
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
		"msg":  "find member success",
		"data": member,
	})
}

func Members(request *gin.Context) {
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

	members := []model_v1.Member{}
	result := database.DB.Preload("Member").Where("user_r_id = ? AND project_r_id = ?", user.RID, project_id).Find(&members)
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
		"msg":  "find members success",
		"data": members,
	})
}

type PostMemberRequest struct {
	Username string `json:"username"`
}

func PostMember(request *gin.Context) {
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

	var data PostMemberRequest
	err = request.BindJSON(&data)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "bad request",
			"data": nil,
		})
		return
	}

	member_user := model_v1.User{}
	result := database.DB.Where("username = ?", data.Username).First(&member_user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no member",
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

	member := model_v1.Member{}
	result = database.DB.Where("user_r_id = ? AND member_r_id = ? AND project_r_id = ?", user.RID, member_user.RID, project_id).First(&member)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			member = model_v1.Member{
				UserRID:    user.RID,
				MemberRID:  member_user.RID,
				ProjectRID: project_id,
				Status:     0,
			}
			result = database.DB.Create(&member)
			if result.Error != nil {
				request.JSON(200, gin.H{
					"code": 50000,
					"msg":  "server error",
					"data": gin.H{
						"message": result.Error.Error(),
					},
				})
				return
			} else {
				request.JSON(200, gin.H{
					"code": 10000,
					"msg":  "add member success",
					"data": nil,
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
	} else {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "member already exist",
			"data": nil,
		})
		return
	}
}

func PutMember(request *gin.Context) {
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

	var member_id = request.Param("member_id")
	if member_id == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "no member id",
			"data": nil,
		})
		return
	}

	var id int64
	var err error
	id, err = strconv.ParseInt(member_id, 10, 64)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "member id is not a number",
			"data": nil,
		})
		return
	}

	var member_form model_v1.Member
	err = request.BindJSON(&member_form)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "bad request",
			"data": nil,
		})
		return
	}

	var member model_v1.Member
	result := database.DB.Where("_id = ? AND user_r_id = ?", id, user.RID).First(&member)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no member",
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

	member.RID = id
	member.Status = member_form.Status
	member.UserRID = user.RID
	result = database.DB.Save(&member)
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
		"msg":  "update member success",
		"data": member,
	})
}

func DeleteMember(request *gin.Context) {
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

	var member_id = request.Param("member_id")
	if member_id == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "no member id",
			"data": nil,
		})
		return
	}

	var id int64
	var err error
	id, err = strconv.ParseInt(member_id, 10, 64)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "member id is not a number",
			"data": nil,
		})
		return
	}

	var member model_v1.Member
	result := database.DB.Where("_id = ? AND user_r_id = ?", id, user.RID).First(&member)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			request.JSON(200, gin.H{
				"code": 40000,
				"msg":  "no member",
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

	result = database.DB.Delete(&member)
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
		"msg":  "delete member success",
		"data": nil,
	})
}
