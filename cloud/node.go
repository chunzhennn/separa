package main

import (
	"separa/cloud/api"

	"github.com/gin-gonic/gin"
)

func Node(port string) {
	r := gin.Default()

	r.GET("/ping", api.Ping)
	r.POST("/task/run", api.TaskRun)
	r.GET("/task/collect", api.Collect)

	r.Run(":" + port)
}
