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
		client_router.POST("/sync_api", client.Apis)
		client_router.POST("/delete_api", client.PutApi)
		client_router.POST("/sync_check", client.SyncCheck)
		client_router.POST("/sync_data", client.SyncData)
		client_router.GET("/project", client.Projects)
	}
	browser_router := r.Group("/browser/api")
	{
		browser_router.POST("/login", browser.Login)
		browser_router.GET("/project/public/:project_id", browser.PublicProject)
		browser_router.GET("/project/public", browser.PublicProjects)
		browser_router.GET("/api/public/:project_id", browser.PublicItems)
	}
	browser_router_auth := r.Group("/browser/api", middleware.AuthLogin())
	{
		browser_router_auth.GET("/project/:project_id", browser.Project)
		browser_router_auth.GET("/project", browser.Projects)
		browser_router_auth.POST("/project", browser.PostProject)
		browser_router_auth.PUT("/project/:project_id", browser.PutProject)
		browser_router_auth.GET("/api/:project_id", browser.Items)
		browser_router_auth.POST("/api/:project_id", browser.PostItem)
		browser_router_auth.GET("/member/:project_id/:member_id", browser.Member)
		browser_router_auth.GET("/member/:project_id", browser.Members)
		browser_router_auth.POST("/member/:project_id", browser.PostMember)
		browser_router_auth.DELETE("/member/:project_id/:member_id", browser.DeleteMember)
		browser_router_auth.GET("/userinfo", browser.UserInfo)
		browser_router_auth.GET("/token", browser.Tokens)
		browser_router_auth.POST("/token", browser.PostToken)
		browser_router_auth.DELETE("/token/:token_id", browser.DeleteToken)
	}
	browser_router_admin := r.Group("/browser/api", middleware.AdminMiddleware())
	{
		browser_router_admin.GET("/user/:user_id", browser.User)
		browser_router_admin.GET("/user", browser.Users)
		browser_router_admin.POST("/user", browser.PostUser)
		browser_router_admin.PUT("/user/:user_id", browser.PutUser)
		browser_router_admin.DELETE("/user/:user_id", browser.DeleteUser)
	}
}
