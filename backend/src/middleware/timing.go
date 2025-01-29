package middleware

import (
	"com668-backend/utility"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	timingField string = "reqStart"
)

func TimingRequestMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now().UnixNano()
		ctx.Set(timingField, now)
		ctx.Next()
	}
}

func TimingResponseMW() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		reqStart := ctx.GetInt64(timingField)
		now := time.Now().UnixNano()
		duration, err := time.ParseDuration(fmt.Sprintf("%dns", now-reqStart))
		if err != nil {
			ctx.Set("Status", http.StatusInternalServerError)
			ctx.Set("Body", &utility.ErrorResponseSchema{
				Error: err.Error(),
			})
			ctx.Next()
			return
		}
		ctx.Header("X-Timing", duration.String())
		ctx.Next()
	}
}
