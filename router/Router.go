package router

import (
	client "ahripost_deploy/controller/client"
	"ahripost_deploy/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	client_router := r.Group("/client/api", middleware.AuthLogin())
	{
		client_router.POST("/sync", client.Sync)
	}
}
