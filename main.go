package main

import (
	"testProject/test/config"
	"testProject/test/router"
)

func main() {
	config.Connect() //初始化MySQL Redis
	router.Router()
}
