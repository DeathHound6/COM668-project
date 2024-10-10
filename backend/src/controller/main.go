package controller

import (
	"com668-backend/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterControllers(engine *gin.Engine) {
	// Register authentication endpoints
	register(engine, http.MethodGet, "/authorise/slack", SlackRedirect())
	register(engine, http.MethodGet, "/authorise/slack/callback", AuthoriseSlack())

	// Register teams endpoints
	register(engine, http.MethodPost, "/teams", CreateTeam())
	register(engine, http.MethodDelete, "/teams/:team_id", DeleteTeam())

	// Register users endpoints
	register(engine, http.MethodPost, "/users", CreateUser())
	register(engine, http.MethodPost, "/users/login", LoginUser())

	// Register incident endpoints
	register(engine, http.MethodPost, "/incidents", CreateIncident())

	// Register settings endpoints
	register(engine, http.MethodGet, "/providers", GetProviders())
	register(engine, http.MethodPost, "/providers", CreateProvider())
	register(engine, http.MethodPut, "/providers/:provider_id", UpdateProvider())
	register(engine, http.MethodGet, "/providers/:provider_id/settings", GetSettings())
	register(engine, http.MethodPut, "/providers/:provider_id/settings", UpdateSettings())
}

func register(engine *gin.Engine, method string, endpoint string, handler gin.HandlerFunc) {
	var registerRouteFunc func(relativePath string, handlers ...gin.HandlerFunc) gin.IRoutes
	switch method {
	case http.MethodDelete:
		registerRouteFunc = engine.DELETE
	case http.MethodPatch:
		registerRouteFunc = engine.PATCH
	case http.MethodPost:
		registerRouteFunc = engine.POST
	case http.MethodPut:
		registerRouteFunc = engine.PUT
	case http.MethodGet:
		registerRouteFunc = engine.GET
	default:
		registerRouteFunc = engine.GET
	}

	// Register the endpoint
	registerRouteFunc(
		endpoint,
		middleware.TimingRequestMW(),
		middleware.UserAuthRequestMW(),
		middleware.TransactionRequestMW(),
		handler,
		middleware.TransactionResponseMW(),
		middleware.TimingResponseMW(),
		middleware.FormatResponseMW(),
	)
}
