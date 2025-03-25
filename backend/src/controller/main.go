package controller

import (
	"com668-backend/middleware"
	"com668-backend/utility"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type registerControllerOptions struct {
	useAdminAuth bool
	useAuth      bool
	useDB        bool
}

func RegisterControllers(engine *gin.Engine) {
	// Register authentication endpoints
	register(engine, http.MethodGet, "/authorise/slack", SlackRedirect(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: false,
	})
	register(engine, http.MethodGet, "/authorise/slack/callback", AuthoriseSlack(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: false,
	})

	// Register teams endpoints
	register(engine, http.MethodGet, "/teams", GetTeams(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: false,
	})

	// Register users endpoints
	register(engine, http.MethodPost, "/users", CreateUser(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: true,
	})
	register(engine, http.MethodPost, "/users/login", LoginUser(), registerControllerOptions{
		useAuth:      false,
		useDB:        true,
		useAdminAuth: false,
	})
	register(engine, http.MethodGet, "/me", GetUser(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: false,
	})

	// Register incident endpoints
	register(engine, http.MethodGet, "/incidents", GetIncidents(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: false,
	})
	register(engine, http.MethodGet, "/incidents/:incident_id", GetIncident(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: false,
	})
	register(engine, http.MethodPost, "/incidents", CreateIncident(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: true,
	})
	register(engine, http.MethodPut, "/incidents/:incident_id", UpdateIncident(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: false,
	})
	register(engine, http.MethodPost, "/incidents/:incident_id/comments", CreateIncidentComment(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: false,
	})
	register(engine, http.MethodDelete, "/incidents/:incident_id/comments/:comment_id", DeleteIncidentComment(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: false,
	})

	// Register settings endpoints
	register(engine, http.MethodGet, "/providers", GetProviders(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: true,
	})
	register(engine, http.MethodGet, "/providers/:provider_id", GetProvider(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: true,
	})
	register(engine, http.MethodPost, "/providers", CreateProvider(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: true,
	})
	register(engine, http.MethodPut, "/providers/:provider_id", UpdateProvider(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: true,
	})
	register(engine, http.MethodDelete, "/providers/:provider_id", DeleteProvider(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: true,
	})

	// Register hosts endpoints
	register(engine, http.MethodGet, "/hosts", GetHosts(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: false,
	})
	register(engine, http.MethodGet, "/hosts/:host_id", GetHost(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: false,
	})
	register(engine, http.MethodPost, "/hosts", CreateHost(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: true,
	})
	register(engine, http.MethodPut, "/hosts/:host_id", UpdateHost(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: true,
	})
	register(engine, http.MethodDelete, "/hosts/:host_id", DeleteHost(), registerControllerOptions{
		useAuth:      true,
		useDB:        true,
		useAdminAuth: true,
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
		if options.useAuth {
			handlers = append(handlers, middleware.UserAuthRequestMW(options.useAdminAuth))
		}
	}

	// Endpoint handler
	handlers = append(handlers, func(ctx *gin.Context) {
		reqID := ctx.GetString("ReqID")
		// Check if we have set an error response already
		// If we have, we skip the normal endpoint handler
		body, ok := ctx.Get("Body")
		if ok {
			_, ok = body.(*utility.ErrorResponseSchema)
		}
		if body == nil || !ok {
			log.Default().Printf("[%s] Running endpoint handler for %s %s\n", reqID, ctx.Request.Method, ctx.Request.URL.Path)
			handler(ctx)
		}
		log.Default().Printf("[%s] Endpoint handler for %s %s completed\n", reqID, ctx.Request.Method, ctx.Request.URL.Path)
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

func getCommonParams(ctx *gin.Context) (map[string]any, error) {
	params := make(map[string]any)

	pageStr := ctx.Query("page")
	if pageStr == "" {
		pageStr = "1"
	}
	pageInt, err := strconv.Atoi(pageStr)
	if err != nil {
		return nil, errors.New("page query parameter must be an integer")
	}
	params["page"] = pageInt

	pageSizeStr := ctx.Query("pageSize")
	if pageSizeStr == "" {
		pageSizeStr = "10"
	}
	pageSizeInt, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		return nil, errors.New("pageSize query parameter must be an integer")
	}
	params["pageSize"] = pageSizeInt

	return params, nil
}
