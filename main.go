package main

import (
	router "ahripost_deploy/router"
	"ahripost_deploy/tools"
	"fmt"
	"strings"

	"ahripost_deploy/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	var cfg = tools.Cfg

	r := gin.Default()
	r.Use(middleware.Cors())
	router.RegisterRouter(r)

	var build strings.Builder
	build.WriteString(cfg.AppHost)
	build.WriteString(":")
	build.WriteString(cfg.AppPort)
	address := build.String()
	fmt.Println("Server run at http://" + address)
	// fmt.Println("cfg=%+v\n", cfg)

	if err := r.Run(address); err != nil {
		panic(err.Error())
	}
}
