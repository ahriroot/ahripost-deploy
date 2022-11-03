package browser

import (
	database "ahripost_deploy/database"
	model_v1 "ahripost_deploy/models/v1"
	"ahripost_deploy/tools"
	"strconv"

	"github.com/gin-gonic/gin"
)

func Tokens(request *gin.Context) {
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

	var tokens []model_v1.Token
	result := database.DB.Where("user_r_id = ?", user.RID).Find(&tokens)
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
		"msg":  "success",
		"data": tokens,
	})
}

func PostToken(request *gin.Context) {
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

	var token string = tools.Sha256(tools.RandomString(32))
	var token_model = model_v1.Token{
		Token:   token,
		UserRID: user.RID,
	}

	result := database.DB.Create(&token_model)
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
		"msg":  "success",
		"data": token,
	})
}

func DeleteToken(request *gin.Context) {
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

	var token_id string = request.Param("token_id")
	if token_id == "" {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "no token id",
			"data": nil,
		})
		return
	}

	var id int64
	id, err := strconv.ParseInt(token_id, 10, 64)
	if err != nil {
		request.JSON(200, gin.H{
			"code": 40000,
			"msg":  "token id is not a number",
			"data": nil,
		})
		return
	}

	var token_model model_v1.Token
	result := database.DB.Where("_id = ? AND user_r_id = ?", id, user.RID).First(&token_model)
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

	result = database.DB.Delete(&token_model)
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
		"msg":  "success",
		"data": nil,
	})
}
