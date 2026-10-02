package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/MalyginDanila/chat-service/internal/config"
	"github.com/MalyginDanila/chat-service/internal/logger"
	serverdebug "github.com/MalyginDanila/chat-service/internal/server-debug"
)

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func run() error {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.ParseAndValidate(*configPath)
	if err != nil {
		return fmt.Errorf("parse and validate config %q: %v", *configPath, err)
	}

	if err := logger.Init(logger.Options{
		Level:          cfg.Log.Level,
		ProductionMode: cfg.Global.IsProduction(),
	}); err != nil {
		return fmt.Errorf("init logger: %v", err)
	}
	defer logger.Sync()

	srvDebug, err := serverdebug.New(cfg.Servers.Debug.Addr)
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)

	// Run servers.
	eg.Go(func() error { return srvDebug.Run(ctx) })

	// Run services.
	// Ждём следующих модулей :)

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	zap.L().Info("graceful shutdown")
	return nil
}
