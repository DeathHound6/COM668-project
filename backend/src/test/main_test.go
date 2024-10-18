package test_test

import (
	"bytes"
	"com668-backend/controller"
	"com668-backend/database"
	"com668-backend/middleware"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

const (
	TestUserEmail    string = "test@example.com"
	TestUserPassword string = "system_user"
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

func getJWT(engine *gin.Engine, userEmail string, userPassword string) (string, error) {
	body := []byte(fmt.Sprintf("{\"email\":\"%s\",\"password\":\"%s\"}", userEmail, userPassword))
	req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	writer := makeRequest(engine, req)
	if writer.Code != http.StatusNoContent {
		return "", fmt.Errorf("status code %d != %d", writer.Code, http.StatusNoContent)
	}
	jwt := writer.Result().Header.Get(middleware.AuthHeaderNameString)
	return jwt, nil
}

func makeRequest(engine *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	writer := httptest.NewRecorder()
	engine.ServeHTTP(writer, req)
	return writer
}
