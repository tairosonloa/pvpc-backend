package server

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-gonic/gin"

	"pvpc-backend/internal/listing"
	"pvpc-backend/internal/platform/server/handler/health"
	"pvpc-backend/internal/platform/server/handler/zones"
	"pvpc-backend/internal/platform/server/middleware"
	"pvpc-backend/internal/platform/storage/postgresql"
)

type Server struct {
	httpAddr        string
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
	listingService listing.ListingService
}

func New(host string, port uint, env string, shutdownTimeout time.Duration, db *sql.DB, dbTimeout time.Duration) Server {
	if env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}
	srv := Server{
		engine:          gin.New(),
		httpAddr:        fmt.Sprintf("%s:%d", host, port),
		shutdownTimeout: shutdownTimeout,
		storage: storage{
			db:        db,
			dbTimeout: dbTimeout,
		},
	}

	srv.registerMiddlewares()
	srv.registerServices()
	srv.registerRoutes()

	return srv
}

func (s *Server) registerMiddlewares() {
	s.engine.Use(gin.Recovery())
	s.engine.Use(middleware.Logger([]string{"/health"}))
}

func (s *Server) registerServices() {
	// Repositories
	pricesZonesRepository := postgresql.NewPricesZonesRepository(s.storage.db, s.storage.dbTimeout)

	// Services
	s.services.listingService = listing.NewListingService(pricesZonesRepository)
}

func (s *Server) registerRoutes() {
	// Health check
	s.engine.GET("/v1/health", health.HealthCheckHandlerV1(s.storage.db, s.storage.dbTimeout))

	// Zones
	s.engine.GET("/v1/zones", zones.ListZonesHandlerV1(s.services.listingService))
}

func (s *Server) Run() {
	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.engine,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info("Server running", "addr", s.httpAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Unexpected server shutdown: %v", err)
		}
	}()

	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	fmt.Println() // Blank line for readability, so ^C is on its own line.
	log.Infof("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Info("Server exiting")
}
