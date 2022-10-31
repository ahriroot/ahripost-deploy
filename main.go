package main

import (
	router "ahripost_deploy/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	router.RegisterRouter(r)

	println("ahripost deploy server start at 0.0.0.0:8080")
	r.Run() // listen and serve on 0.0.0.0:8080
}
