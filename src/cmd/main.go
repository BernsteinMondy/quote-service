package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/BernsteinMondy/quote-service/src/internal/httpserver"
	"github.com/BernsteinMondy/quote-service/src/internal/impl"
	"github.com/BernsteinMondy/quote-service/src/internal/service"
	"github.com/BernsteinMondy/quote-service/src/pkg/database"
	"github.com/gorilla/mux"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	err := run()
	if err != nil {
		slog.Error("run() returned error", slog.String("error", err.Error()))
	}
}

func run() (err error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP)
	defer stop()

	slog.Info("Loading config")
	cfg, err := loadConfigFromEnv()
	if err != nil {
		return fmt.Errorf("load config from env: %w", err)
	}
	slog.Info("Config loaded")

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	slog.Info("Connecting to database")
	db, err := sqlDB(cfg.DB)
	if err != nil {
		return fmt.Errorf("new sql database: %w", err)
	}
	slog.Info("Connected to database")
	defer func() {
		slog.Info("Closing database connection")
		err = errors.Join(err, db.Close())
		slog.Info("Database connection closed")
	}()

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	quoteRepo := impl.NewQuoteRepository(db)
	quoteService := service.New(quoteRepo)

	router := mux.NewRouter()
	server := httpserver.New(quoteService, router, cfg.Server.Port)
	stopWg := sync.WaitGroup{}

	stopWg.Add(1)
	go func(ctx context.Context) {
		defer stopWg.Done()
		httpSrvErr := launchHTTPServer(ctx, server)
		if httpSrvErr != nil {
			slog.Error("launchHTTPServer() returned error", slog.String("error", httpSrvErr.Error()))
		}
	}(ctx)

	<-ctx.Done()
	stopWg.Wait()
	return nil
}

func sqlDB(cfg DB) (*sql.DB, error) {
	dbConfig := database.Config{
		User:     cfg.User,
		Password: cfg.Password,
		Name:     cfg.Name,
		Host:     cfg.Host,
		Port:     cfg.Port,
		SSLMode:  cfg.SSLMode,
	}

	return database.NewSQLDatabase(dbConfig)
}

func launchHTTPServer(ctx context.Context, server *http.Server) (err error) {
	var httpServerShutDownError error
	defer func() {
		err = errors.Join(err, httpServerShutDownError)
	}()

	shutDownDone := make(chan struct{})
	go func(ctx context.Context) {
		<-ctx.Done()

		slog.Info("Shutting down http server")
		httpServerShutDownError = server.Shutdown(ctx)
		slog.Info("Http server shut down")

		close(shutDownDone)
	}(ctx)

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	slog.Info("Starting http server", slog.String("addr", server.Addr))
	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("listen on %s: %w", server.Addr, err)
	}

	<-shutDownDone
	return nil
}
