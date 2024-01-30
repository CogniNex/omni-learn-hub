// Package app configures and runs application.
package app

import (
	"fmt"
	"omni-learn-hub/config"
	"omni-learn-hub/internal/controller/http/v1"
	"omni-learn-hub/internal/repository/pgsqlrepo"
	"omni-learn-hub/internal/service"
	"omni-learn-hub/pkg/hash"
	"omni-learn-hub/pkg/postgres"

	"omni-learn-hub/pkg/httpserver"
	"omni-learn-hub/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	// Repository
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PG.PoolMax))
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer pg.Close()

	// Services, Repos
	repos := pgsqlrepo.NewRepositories(pg)
	hasher := hash.NewBcryptPasswordHasher()
	services := service.NewServices(service.Deps{
		Repos:  repos,
		Hasher: hasher,
	})

	//// RabbitMQ RPC Server
	//rmqRouter := amqprpc.NewRouter(translationUseCase)
	//
	//rmqServer, err := server.New(cfg.RMQ.URL, cfg.RMQ.ServerExchange, rmqRouter, l)
	//if err != nil {
	//	l.Fatal(fmt.Errorf("app - Run - rmqServer - server.New: %w", err))
	//}

	// HTTP Server
	handler := gin.New()
	v1.NewRouter(handler, l, services)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())

	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}

	//err = rmqServer.Shutdown()
	//if err != nil {
	//	l.Error(fmt.Errorf("app - Run - rmqServer.Shutdown: %w", err))
	//}
}
