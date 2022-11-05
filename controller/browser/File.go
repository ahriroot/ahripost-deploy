package browser

import (
	model_v1 "ahripost_deploy/models/v1"

	"github.com/gin-gonic/gin"
)

func File(request *gin.Context) {
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

	// TODO: 创建 gzip 压缩文件，保存到 /data

	request.JSON(200, gin.H{
		"code": 10000,
		"msg":  "login success",
		"data": gin.H{
			"token": user,
		},
	})
}
