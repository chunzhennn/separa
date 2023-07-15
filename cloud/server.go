package main

import (
	"separa/cloud/api"

	"github.com/gin-gonic/gin"
)

func Server(port string) {
	r := gin.Default()

	r.GET("/ping", api.Ping)
	// r.GET("/node/list", api.NodeList)
	r.GET("/node/add", api.NodeAdd)
	// r.GET("/node/update", api.NodeUpdate)
	// r.GET("/node/delete", api.NodeDelete)

	r.POST("/task/distribute", api.TaskDistribute)
	// r.GET("/task/status", api.TaskStatus)
	r.GET("/task/result", api.TaskResult)

	r.Run(":" + port)
}
