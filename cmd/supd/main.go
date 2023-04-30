package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sashabaranov/go-openai"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/pkg/apimachinery"
	"github.com/greenboxal/aip/pkg/apis"
	"github.com/greenboxal/aip/pkg/collective/comms"
	"github.com/greenboxal/aip/pkg/config"
	"github.com/greenboxal/aip/pkg/daemon"
	"github.com/greenboxal/aip/pkg/ford"
	"github.com/greenboxal/aip/pkg/ford/forddb"
	"github.com/greenboxal/aip/pkg/network/ipfs"
	"github.com/greenboxal/aip/pkg/network/p2p"
	"github.com/greenboxal/aip/pkg/storage/memgraph"
	"github.com/greenboxal/aip/pkg/storage/milvus"
)

func main() {
	app := fx.New(
		BuildLogging(),

		config.Module,
		apimachinery.Module,
		apis.Module,
		p2p.Module,
		ipfs.Module,
		comms.Module,
		ford.Module,
		daemon.Module,

		fx.Provide(func() *openai.Client {
			return openai.NewClient(os.Getenv("OPENAI_API_KEY"))
		}),

		//badger.WithBadgerStorage(),
		memgraph.WithMemgraphStorage(),
		milvus.WithMilvusStorage(),

		fx.Invoke(func(d *daemon.Daemon, db forddb.Database, _api *apimachinery.API) error {
			if err := forddb.ImportPath(db, "./data"); err != nil {
				return err
			}

			return d.Run()
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
