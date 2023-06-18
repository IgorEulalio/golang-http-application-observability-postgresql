package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/database"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/handler"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/logger"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/metrics"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/middleware"
	"github.com/IgorEulalio/golang-http-application-observability-postgresql/pkg/tracer"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

const SERVICE_NAME = "repository-service"

func main() {

	// Initialize components
	db := initializeDatabase()
	defer db.Close()

	shutdown := initializeTracer()

	// Create server
	r := setupRouter()

	// Setup repository route
	handler.GetAllRepositories(r, db, "/repositories")

	srv := startServer(r)

	// Wait for interrupt signal
	waitForShutdownSignal()

	// Shutdown
	shutdownServer(srv)
	shutdownTracer(shutdown)
}

func initializeDatabase() *sqlx.DB {
	db, err := database.ConnectToDatabase()
	if err != nil {
		logger.Log.Fatal(err)
	}
	logger.Log.Info("Connected to database.")
	return db
}

func initializeTracer() func(context.Context) error {
	shutdown, err := tracer.InitProvider()
	if err != nil {
		logger.Log.Fatal(err)
	}
	logger.Log.Info("Initialized tracer.")
	return shutdown
}

func InitAndSetGlobalMeterProvider() error {
	ctx := context.Background()
	meterProvider, err := metrics.InitMetricsProvider(ctx)
	if err != nil {
		return err
	}
	logger.Log.Info("Initialized metrics provider.")

	// Set as global MeterProvider
	otel.SetMeterProvider(meterProvider)

	return nil
}

func setupRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(otelmux.Middleware(SERVICE_NAME))
	r.Use(logger.LogRequestResponse)
	r.Use(middleware.HTTPRequestCounter)

	return r
}

func startServer(r *mux.Router) *http.Server {
	srv := &http.Server{
		Addr:    ":8080",
		Handler: otelhttp.NewHandler(r, SERVICE_NAME),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal(err)
		}
		logger.Log.Info("Listening on port 8080.")
	}()
	return srv
}

func waitForShutdownSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

func shutdownServer(srv *http.Server) {
	logger.Log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatalf("Server forced to shutdown: %v", err)
	}
	logger.Log.Info("Server successfully shutdown")
}

func shutdownTracer(shutdown func(context.Context) error) {
	logger.Log.Info("Shutting down TracerProvider...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := shutdown(ctx); err != nil {
		logger.Log.Fatalf("Failed to shutdown TracerProvider: %v", err)
	}
	logger.Log.Info("TracerProvider successfully shutdown")
}
