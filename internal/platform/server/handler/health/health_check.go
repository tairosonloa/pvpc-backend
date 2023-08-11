package health

import (
	"database/sql"
	"time"

	"github.com/alexliesenfeld/health"
	"github.com/gin-gonic/gin"
)

func HealthCheckHandler(db *sql.DB, dbTimeout time.Duration) gin.HandlerFunc {
	checker := health.NewChecker(
		health.WithCheck(health.Check{
			Name:    "database",
			Timeout: dbTimeout,
			Check:   db.PingContext,
		}),
	)

	return gin.WrapF(health.NewHandler(checker))
}
