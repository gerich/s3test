package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gerich/s3test/internal/adapters/filesystem"
	"github.com/gerich/s3test/internal/app"
	"github.com/gerich/s3test/internal/config"
	"github.com/gerich/s3test/internal/ports/rest"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func logger(cfg *config.Config) *zap.Logger {
	var zlevel zapcore.Level

	_ = zlevel.Set(cfg.LogLevel)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			NameKey:        "logger",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			CallerKey:      "caller",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}),
		zapcore.AddSync(os.Stdout),
		zap.NewAtomicLevelAt(zlevel),
	)

	return zap.New(
		core,
		zap.AddCaller(),
	)
}

func main() {
	cfg := config.New()
	log := logger(cfg)
	repository := filesystem.New(cfg.FileSystemConfig, log)
	if err := repository.Init(); err != nil {
		log.Fatal("cant prepare filesystem", zap.Error(err))
	}
	service := app.NewService(repository, log)
	server := rest.NewServer(service, cfg.HTTP, log)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	server.Run()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	<-sigCh
	server.Stop(ctx)
}
