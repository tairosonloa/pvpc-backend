package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"go-pvpc/internal/listing"
	"go-pvpc/internal/platform/server/handler/health"
	"go-pvpc/internal/platform/server/handler/zones"
	"go-pvpc/internal/platform/storage/postgresql"

	"github.com/gin-gonic/gin"
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
		engine:          gin.Default(),
		httpAddr:        fmt.Sprintf("%s:%d", host, port),
		shutdownTimeout: shutdownTimeout,
		storage: storage{
			db:        db,
			dbTimeout: dbTimeout,
		},
	}

	srv.registerServices()
	srv.registerRoutes()

	return srv
}

func (s *Server) registerServices() {
	// Repositories
	pricesZoneRepository := postgresql.NewPricesZoneRepository(s.storage.db, s.storage.dbTimeout)

	// Services
	s.services.listingService = listing.NewListingService(pricesZoneRepository)
}

func (s *Server) registerRoutes() {
	// Health check
	s.engine.GET("/health", health.HealthCheckHandler(s.storage.db, s.storage.dbTimeout))

	// Zones
	s.engine.GET("/zones", zones.ListZonesHandler(s.services.listingService))
}

func (s *Server) Run() {
	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.engine,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Println("Server running on", s.httpAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Unexpected server shutdown", err)
		}
	}()

	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	fmt.Println("")
	log.Println("Shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
