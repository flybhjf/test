package router

import "github.com/gin-gonic/gin"

func Router() {
	router := gin.Default()
	// 简单的路由组: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login")
		v1.POST("/submit")
		v1.POST("/read")
	}
	router.Run(":8080")
}
