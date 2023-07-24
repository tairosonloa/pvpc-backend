package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthCheckHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	}
}
