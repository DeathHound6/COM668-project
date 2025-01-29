package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func FormatResponseMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		status, statusOk := ctx.Get("Status")
		if !statusOk {
			log.Default().Println("Status not set in context. Assuming 200 OK")
			status = http.StatusOK
		}
		if http.MethodPost == ctx.Request.Method {
			log.Default().Println("Body not set in context on POST. Assuming 201 Created")
			status = http.StatusCreated
		}
		body, bodyOk := ctx.Get("Body")
		if !bodyOk {
			log.Default().Printf("Body not set in context on %s. Assuming 204 No Content\n", ctx.Request.Method)
			status = http.StatusNoContent
		}
		if !bodyOk || status.(int) == http.StatusNoContent {
			log.Default().Printf("No body to return. Returning with status %d\n", status.(int))
			ctx.AbortWithStatus(status.(int))
			return
		}
		if gin.IsDebugging() {
			log.Default().Printf("Returning body with status %d\n%v\n", status.(int), body)
		}
		ctx.AbortWithStatusJSON(status.(int), body)
	}
}
