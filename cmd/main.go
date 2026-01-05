package main

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/jllovet/go-server-template/config"
	"github.com/jllovet/go-server-template/internal/server"
	"github.com/jllovet/go-server-template/internal/todo"
	"github.com/jllovet/go-server-template/internal/todo/memory"
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

	logger := logger.New(stdout, getenv("SERVICE_NAME", "todo-service"))

	repo := memory.New()
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
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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
