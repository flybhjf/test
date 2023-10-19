package router

import (
	"testProject/test/hander"

	"github.com/gin-gonic/gin"
)

func Router() {
	router := gin.Default()
	// 简单的路由组: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/register", hander.RegisterUser)
		v1.POST("/login", hander.Login)
	}
	router.Run(":8080")
}
