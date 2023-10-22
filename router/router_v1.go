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
		v1.POST("/register", hander.RegisterUser)        //注册
		v1.POST("/login", hander.Login)                  //登录
		v1.POST("/reset_password", hander.ResetPassword) //登录
	}
	router.Run(":8080")
}
