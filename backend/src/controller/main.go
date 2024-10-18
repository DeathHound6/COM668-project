package controller

import (
	"com668-backend/middleware"
	"com668-backend/utility"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type registerControllerOptions struct {
	useAuth bool
	useDB   bool
}

func RegisterControllers(engine *gin.Engine) {
	// Register authentication endpoints
	register(engine, http.MethodGet, "/authorise/slack", SlackRedirect(), registerControllerOptions{
		useAuth: true,
		useDB:   true,
	})
	register(engine, http.MethodGet, "/authorise/slack/callback", AuthoriseSlack(), registerControllerOptions{
		useAuth: true,
		useDB:   true,
	})

	// Register teams endpoints
	register(engine, http.MethodPost, "/teams", CreateTeam(), registerControllerOptions{
		useAuth: true,
		useDB:   true,
	})
	register(engine, http.MethodDelete, "/teams/:team_id", DeleteTeam(), registerControllerOptions{
		useAuth: true,
		useDB:   true,
	})

	// Register users endpoints
	register(engine, http.MethodPost, "/users", CreateUser(), registerControllerOptions{
		useAuth: false,
		useDB:   true,
	})
	register(engine, http.MethodPost, "/users/login", LoginUser(), registerControllerOptions{
		useAuth: false,
		useDB:   true,
	})

	// Register incident endpoints
	register(engine, http.MethodPost, "/incidents", CreateIncident(), registerControllerOptions{
		useAuth: true,
		useDB:   true,
	})

	// Register settings endpoints
	register(engine, http.MethodGet, "/providers", GetProviders(), registerControllerOptions{
		useAuth: true,
		useDB:   true,
	})
	register(engine, http.MethodPost, "/providers", CreateProvider(), registerControllerOptions{
		useAuth: true,
		useDB:   true,
	})
	register(engine, http.MethodPut, "/providers/:provider_id", UpdateProvider(), registerControllerOptions{
		useAuth: true,
		useDB:   true,
	})
}

func register(engine *gin.Engine, method string, endpoint string, handler gin.HandlerFunc, options registerControllerOptions) {
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

	handlers := []gin.HandlerFunc{}

	// Request handlers
	handlers = append(handlers, middleware.TimingRequestMW())
	if options.useDB {
		handlers = append(handlers, middleware.TransactionRequestMW())
	}
	if options.useAuth {
		handlers = append(handlers, middleware.UserAuthRequestMW())
	}

	// Endpoint handler
	handlers = append(handlers, func(ctx *gin.Context) {
		// Check if we have set an error response already
		// If we have, we skip the normal endpoint handler
		body, ok := ctx.Get("Body")
		if ok {
			_, ok = body.(*utility.ErrorResponseSchema)
		}
		if !ok {
			log.Default().Println("Running endpoint handler")
			handler(ctx)
		}
		ctx.Next()
	})

	// Response handlers
	if options.useDB {
		handlers = append(handlers, middleware.TransactionResponseMW())
	}
	handlers = append(handlers, middleware.TimingResponseMW())
	handlers = append(handlers, middleware.FormatResponseMW())

	// Register the endpoint
	registerRouteFunc(
		endpoint,
		handlers...,
	)
}
