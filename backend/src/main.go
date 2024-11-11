package main

import (
	"com668-backend/controller"
	"com668-backend/database"
	_ "com668-backend/docs" // import docs to register the swagger definition
	"com668-backend/middleware"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	swaggerGin "github.com/swaggo/gin-swagger"
)

//	@title						A.I.M.S Swagger
//	@version					1.0
//	@host						localhost:5000
//	@BasePath					/
//	@schemes					https
//	@accept						json
//	@produce					json
//	@securitydefinitions.apikey	JWT
//	@in							header
//	@name						Authorization
//	@description				The JWT
func main() {
	if err := database.Connect(); err != nil {
		panic(err)
	}

	// Setup HTTP webserver
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(gin.CustomRecoveryWithWriter(gin.DefaultErrorWriter, middleware.RecoveryHandler()))
	engine.GET("/swagger/*any", swaggerGin.WrapHandler(swaggerFiles.Handler))
	engine.HandleMethodNotAllowed = true
	controller.RegisterControllers(engine)

	// Run the webserver in a goroutine (non blocking call)
	go (func() {
		addr := fmt.Sprintf(":%d", 5000)
		tlsCertFile := os.Getenv("TLS_CERT_FILE")
		tlsKeyFile := os.Getenv("TLS_KEY_FILE")
		if err := engine.RunTLS(addr, tlsCertFile, tlsKeyFile); err != nil {
			panic(err)
		}
	})()

	// Graceful exit handler
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, syscall.SIGINT, syscall.SIGTERM)
	<-exitSignal
	close(exitSignal)
}
