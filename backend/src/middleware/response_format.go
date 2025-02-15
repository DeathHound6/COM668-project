package middleware

import (
	"com668-backend/utility"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func FormatResponseMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqID := ctx.GetString("ReqID")
		status, statusOk := ctx.Get("Status")
		if !statusOk {
			log.Default().Printf("[%s] Status not set in context. Assuming 200 OK\n", reqID)
			status = http.StatusOK
		}
		body, bodyOk := ctx.Get("Body")
		if strings.HasPrefix(strconv.Itoa(status.(int)), "2") {
			if http.MethodPost == ctx.Request.Method {
				log.Default().Printf("[%s] Body not set in context on POST. Assuming 201 Created\n", reqID)
				status = http.StatusCreated
			}
			if !bodyOk {
				log.Default().Printf("[%s] Body not set in context. Assuming 204 No Content\n", reqID)
				status = http.StatusNoContent
			}
		}
		if !bodyOk || status == http.StatusNoContent {
			log.Default().Printf("[%s] No body to return. Returning with status %d\n", reqID, status.(int))
			ctx.AbortWithStatus(status.(int))
			return
		}
		// for privacy/security reasons, we don't want to log the body unless we are in dev mode
		if gin.IsDebugging() {
			log.Default().Printf("[%s] Returning body with status %d %v\n", reqID, status.(int), body.(utility.ResponseSchema).String())
		}
		ctx.AbortWithStatusJSON(status.(int), body.(utility.ResponseSchema).JSON())
	}
}
