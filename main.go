package main

import (
	router "ahripost_deploy/router"

	"ahripost_deploy/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(middleware.Cors())
	router.RegisterRouter(r)

	println("ahripost deploy server start at 0.0.0.0:8080")
	r.Run() // listen and serve on 0.0.0.0:8080
}
