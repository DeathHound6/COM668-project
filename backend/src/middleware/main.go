package middleware

import "github.com/gin-gonic/gin"

var (
	RequestMW []gin.HandlerFunc = []gin.HandlerFunc{
		TransactionRequestMW(),
	}

	ResponseMW []gin.HandlerFunc = []gin.HandlerFunc{
		TransactionResponseMW(),
	}
)
