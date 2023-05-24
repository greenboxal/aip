package main

import (
	"context"
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
	apimachinery2 "github.com/greenboxal/aip/aip-forddb/pkg/apimachinery"
	"github.com/greenboxal/aip/aip-forddb/pkg/apis"
	"github.com/greenboxal/aip/aip-forddb/pkg/forddb"
	forddbimpl "github.com/greenboxal/aip/aip-forddb/pkg/impl"
	"github.com/greenboxal/aip/aip-forddb/pkg/jobs"
	"github.com/greenboxal/aip/aip-forddb/pkg/objectstore/firestore"
	"github.com/greenboxal/aip/aip-forddb/pkg/tracing"
	"github.com/greenboxal/aip/aip-langchain/pkg/providers/openai"
	"github.com/greenboxal/aip/aip-langchain/pkg/vectorstore/milvus"
	"github.com/greenboxal/aip/aip-sdk/pkg/config"
	"github.com/greenboxal/aip/aip-sdk/pkg/network/ipfs"
	"github.com/greenboxal/aip/aip-sdk/pkg/network/p2p"
	"github.com/greenboxal/aip/aip-wiki/pkg/wiki"
)

func main() {
	app := fx.New(
		BuildLogging(),

		config.Module,
		apimachinery2.Module,
		apis.Module,
		p2p.Module,
		ipfs.Module,
		openai.Module,
		comms.Module,
		msn.Module,
		ford.Module,
		forddbimpl.Module,
		jobs.Module,
		daemon.Module,
		wiki.Module,
		tracing.Module,

		//inmemory.WithInIndexingMemoryDatabase(),
		//badger.WithBadgerStorage(),
		//memgraph.WithMemgraphStorage(),
		//ipfs.WithIpfsStorage(),
		firestore.WithObjectStore(),
		milvus.WithIndexStorage(),

		fx.Invoke(func(d *daemon.Daemon, db forddb.Database, _api *apimachinery2.Server) {
		}),

		fx.Provide(NewChatHandler),

		fx.Invoke(func(logger *zap.SugaredLogger, handler *ChatHandler) {
			logger.Info("I'm alive.")

			//	goprocess.Go(handler.Run)
		}),
	)

	go func() {
		signalCh := make(chan os.Signal, 1)

		signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		<-signalCh

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := app.Stop(shutdownCtx); err != nil {
			panic(err)
		}
	}()

	app.Run()
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
