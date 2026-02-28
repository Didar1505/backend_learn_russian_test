package main

import (
	"github.com/Didar1505/project_test.git/internal/app"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	app := app.NewApplication(router)
	app.InitApp()
	// db init
	// config load environments
	// load logger pkg
	// cors middleware | middlewares

	// init all instances (modules)
	
	// register all routes from handler

	// defer log sync
	router.Run(":8080")
}