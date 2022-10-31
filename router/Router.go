package router

import (
	browser "ahripost_deploy/controller/browser"
	client "ahripost_deploy/controller/client"
	"ahripost_deploy/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	client_router := r.Group("/client/api", middleware.TokenLogin())
	{
		client_router.POST("/sync", client.Sync)
	}
	browser_router := r.Group("/browser/api")
	{
		browser_router.POST("/login", browser.Login)
	}
	browser_router_auth := r.Group("/browser/api", middleware.AuthLogin())
	{
		browser_router_auth.GET("/project/:project_id", browser.Project)
		browser_router_auth.GET("/project", browser.Projects)
		browser_router_auth.GET("/api/:project_id", browser.Items)
	}
}
