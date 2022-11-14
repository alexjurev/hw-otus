package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/server/http"
	"github.com/alexjurev/hw-otus/hw12_13_14_15_calendar/internal/storagebuilder"
	log "github.com/sirupsen/logrus"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		log.Errorf("failed to start %v", err)
		return
	}
	err = logger.PrepareLogger(config.Logger)
	if err != nil {
		log.Errorf("failed to start %v", err)
		return
	}
	storage, err := storagebuilder.New(config.Storage)
	if err != nil {
		log.Errorf("failed to start %v", err)
		return
	}

	calendar := app.New(storage)
	server := internalhttp.NewServer(config.Server, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			log.Errorf("failed to stop http server: %s", err.Error())
		}
	}()

	log.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		log.Errorf("failed to start http server: %s", err.Error())
		cancel()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err := storage.Close(ctx)
		if err != nil {
			log.Errorf("failed to close storage: %s", err.Error())
		}
		os.Exit(1) //nolint:gocritic
	}
	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err = storage.Close(ctx)
	if err != nil {
		log.Errorf("failed to close storage: %v", err)
	}
}
