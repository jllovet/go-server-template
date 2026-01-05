package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jllovet/go-server-template/config"
	"github.com/jllovet/go-server-template/internal/server"
	"github.com/jllovet/go-server-template/internal/todo"
	"github.com/jllovet/go-server-template/internal/todo/memory"
	"github.com/jllovet/go-server-template/internal/todo/postgres"
	"github.com/jllovet/go-server-template/logger"
)

// See: https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
func main() {
	ctx := context.Background()
	if err := run(
		ctx,
		os.Args,
		config.GetEnv,
		os.Stdin,
		os.Stdout,
		os.Stderr,
	); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	args []string,
	getenv func(key string, defaultValue string) string,
	stdin io.Reader,
	stdout, stderr io.Writer,
) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	config := &config.InitializedConfig

	logOutput := stdout
	if getenv("DISABLE_LOGGING", "") == "true" {
		logOutput = io.Discard
	}
	logger := logger.New(logOutput, getenv("SERVICE_NAME", "todo-service"))

	var repo todo.Repository
	if config.DatabaseURL != "" {
		db, err := sql.Open("pgx", config.DatabaseURL)
		if err != nil {
			return fmt.Errorf("open db: %w", err)
		}
		defer db.Close()
		if err := db.Ping(); err != nil {
			return fmt.Errorf("ping db: %w", err)
		}
		repo = postgres.New(db)
	} else {
		repo = memory.New()
	}
	service := todo.NewService(repo)

	srv := server.NewServer(
		service,
		config,
		logger,
	)
	httpServer := &http.Server{
		Addr:    net.JoinHostPort(config.Host, config.Port),
		Handler: srv,
	}
	go func() {
		logger.Info("listening", "address", httpServer.Addr)
		var err error
		if config.CertFile != "" && config.KeyFile != "" {
			err = httpServer.ListenAndServeTLS(config.CertFile, config.KeyFile)
		} else {
			err = httpServer.ListenAndServe()
		}
		if err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()
	wg.Wait()
	return nil
}
