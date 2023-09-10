package http

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"pvpc-backend/internal/platform/http/handlers/health"
	"pvpc-backend/internal/platform/http/handlers/prices"
	"pvpc-backend/internal/platform/http/handlers/zones"
	"pvpc-backend/internal/platform/http/middlewares"
	"pvpc-backend/internal/platform/providers/redataapi"
	"pvpc-backend/internal/platform/storage/postgresql"
	servicespkg "pvpc-backend/internal/services"
	"pvpc-backend/pkg/logger"
)

type HttpServer struct {
	address         string
	engine          *gin.Engine
	shutdownTimeout time.Duration
	storage         storage
	services        services
}

type storage struct {
	db        *sql.DB
	dbTimeout time.Duration
}

type services struct {
	pricesService servicespkg.PricesService
	zonesService  servicespkg.ZonesService
}

func NewHttpServer(host string, port uint, env string, shutdownTimeout time.Duration, db *sql.DB, dbTimeout time.Duration, reeApiUrl string) HttpServer {
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	srv := HttpServer{
		engine:          gin.New(),
		address:         fmt.Sprintf("%s:%d", host, port),
		shutdownTimeout: shutdownTimeout,
		storage: storage{
			db:        db,
			dbTimeout: dbTimeout,
		},
	}

	srv.registerMiddlewares()
	srv.registerServices(reeApiUrl)
	srv.registerRoutes()

	return srv
}

func (s *HttpServer) registerMiddlewares() {
	s.engine.Use(gin.Recovery())
	s.engine.Use(middlewares.Logger([]string{"/v1/health"}))
}

func (s *HttpServer) registerServices(reeApiUrl string) {
	// Providers
	pricesProvider := redataapi.NewREDataAPI(reeApiUrl)

	// Repositories
	pricesRepository := postgresql.NewPricesRepository(s.storage.db, s.storage.dbTimeout)
	zonesRepository := postgresql.NewZonesRepository(s.storage.db, s.storage.dbTimeout)

	// Services
	s.services.pricesService = servicespkg.NewPricesService(pricesProvider, pricesRepository, zonesRepository)
	s.services.zonesService = servicespkg.NewZonesService(zonesRepository)
}

func (s *HttpServer) registerRoutes() {
	// Health check
	s.engine.GET("/v1/health", health.HealthCheckHandlerV1(s.storage.db, s.storage.dbTimeout))

	// Prices
	s.engine.POST("/v1/prices", prices.CreatePricesV1(s.services.pricesService))

	// Zones
	s.engine.GET("/v1/zones", zones.ListZonesHandlerV1(s.services.zonesService))
}

func (s *HttpServer) Run() {
	srv := &http.Server{
		Addr:     s.address,
		Handler:  s.engine,
		ErrorLog: logger.ServerErrorLoggerFromDefault(),
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		logger.Info("Server running", "address", s.address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Unexpected server shutdown", "err", err)
		}
	}()

	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	fmt.Println() // Blank line for readability, so ^C is on its own line.
	logger.Info("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", "err", err)
	}

	logger.Info("Server exiting")
}
