package test_test

import (
	"com668-backend/controller"
	"com668-backend/database"

	"github.com/gin-gonic/gin"
)

func setup() *gin.Engine {
	if err := database.Connect(); err != nil {
		panic(err)
	}
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	controller.RegisterControllers(engine)
	engine.Use(gin.Recovery())
	engine.HandleMethodNotAllowed = true
	return engine
}
