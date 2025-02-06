package test_test

import (
	"bytes"
	"com668-backend/controller"
	"com668-backend/database"
	"com668-backend/middleware"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

const (
	TestAdminEmail    string = "test@example.com"
	TestAdminPassword string = "system_user"
	TestUserEmail     string = "user1@example.com"
	TestUserPassword  string = "test_user"
)

func setup() *gin.Engine {
	if err := database.Connect(); err != nil {
		panic(err)
	}
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	controller.RegisterControllers(engine)
	engine.Use(gin.CustomRecoveryWithWriter(gin.DefaultErrorWriter, middleware.RecoveryHandler()))
	engine.HandleMethodNotAllowed = true
	return engine
}

func getJWT(engine *gin.Engine, userEmail string, userPassword string) (string, error) {
	body, err := getJSONBodyAsReader(map[string]any{
		"email":    userEmail,
		"password": userPassword,
	})
	if err != nil {
		return "", err
	}
	req, err := http.NewRequest(http.MethodPost, "/users/login", body)
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

func getJSONBodyAsReader(body map[string]any) (io.Reader, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(jsonBody), nil
}
