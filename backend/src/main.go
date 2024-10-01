package main

import (
	"com668-backend/controller"
	"com668-backend/database"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := database.Connect(); err != nil {
		panic(err)
	}

	// Setup HTTP webserver
	gin.SetMode(gin.DebugMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
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
}
