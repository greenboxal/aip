package main

import (
	"context"
	"encoding/base64"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/aip-controller/pkg/collective/comms"
	"github.com/greenboxal/aip/aip-controller/pkg/collective/msn"
	"github.com/greenboxal/aip/aip-controller/pkg/daemon"
	"github.com/greenboxal/aip/aip-controller/pkg/ford"
	"github.com/greenboxal/aip/aip-forddb/pkg/apimachinery"
	"github.com/greenboxal/aip/aip-forddb/pkg/apis"
	forddbimpl "github.com/greenboxal/aip/aip-forddb/pkg/impl"
	"github.com/greenboxal/aip/aip-forddb/pkg/jobs"
	"github.com/greenboxal/aip/aip-forddb/pkg/objectstore/firestore"
	"github.com/greenboxal/aip/aip-forddb/pkg/tracing"
	"github.com/greenboxal/aip/aip-langchain/pkg/providers/openai"
	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore/milvus"
	"github.com/greenboxal/aip/aip-sdk/pkg/cli"
	"github.com/greenboxal/aip/aip-sdk/pkg/config"
	"github.com/greenboxal/aip/aip-sdk/pkg/network/ipfs"
	"github.com/greenboxal/aip/aip-sdk/pkg/network/p2p"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki"
)

func main() {
	var cliManager *cli.Manager

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	defer close(signalCh)
	defer cancel()

	SetupGoogleCredentials()

	app := fx.New(
		BuildLogging(),

		config.Module,
		cli.Module,
		apimachinery.Module,
		apis.Module,
		p2p.Module,
		ipfs.Module,
		openai.Module,
		daemon.Module,
		jobs.Module,
		comms.Module,
		msn.Module,
		ford.Module,
		forddbimpl.Module,
		wiki.Module,
		tracing.Module,

		milvus.WithIndexStorage(),
		firestore.WithObjectStore(),

		fx.Invoke(func(m *cli.Manager) {
			cliManager = m
		}),
	)

	go func() {
		_, _ = <-signalCh

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := app.Stop(shutdownCtx); err != nil {
			panic(err)
		}

		os.Exit(0)
	}()

	if err := app.Start(ctx); err != nil {
		panic(err)
	}

	if err := cliManager.Run(); err != nil {
		panic(err)
	}
}

func BuildLogging() fx.Option {
	return fx.Options(
		fx.Provide(func() (*zap.Logger, error) {
			return zap.NewDevelopment()
		}),

		fx.Provide(func(l *zap.Logger) *zap.SugaredLogger {
			return l.Sugar()
		}),

		fx.WithLogger(func(l *zap.Logger) fxevent.Logger {
			zl := &fxevent.ZapLogger{Logger: l}
			zl.UseLogLevel(-2)
			zl.UseErrorLevel(zap.ErrorLevel)
			return zl
		}),
	)
}

func SetupGoogleCredentials() {
	if creds := os.Getenv("GOOGLE_CREDENTIALS"); creds != "" {
		tmpFile, err := os.CreateTemp("", "google-credentials-*.json")

		if err != nil {
			panic(err)
		}

		data, err := base64.StdEncoding.DecodeString(creds)

		if err != nil {
			panic(err)
		}

		if err := os.WriteFile(tmpFile.Name(), data, 0600); err != nil {
			panic(err)
		}

		if err := os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", tmpFile.Name()); err != nil {
			panic(err)
		}
	}
}
