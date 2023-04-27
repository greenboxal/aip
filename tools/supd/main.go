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

	"github.com/greenboxal/aip/pkg/api"
	"github.com/greenboxal/aip/pkg/collective/comms"
	"github.com/greenboxal/aip/pkg/collective/transports/pubsub"
	"github.com/greenboxal/aip/pkg/collective/transports/slack"
	"github.com/greenboxal/aip/pkg/daemon"
	"github.com/greenboxal/aip/pkg/ford"
	"github.com/greenboxal/aip/pkg/network/p2p"
	"github.com/greenboxal/aip/pkg/supervisor"
)

func main() {
	app := fx.New(
		BuildLogging(),

		api.Module,
		daemon.Module,
		p2p.Module,
		ford.Module,

		fx.Provide(supervisor.NewManager),
		fx.Provide(comms.NewManager),
		fx.Provide(slack.NewTransport),
		fx.Provide(pubsub.NewTransport),

		fx.Invoke(func(d *daemon.Daemon, _api *api.API) error {
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
	return fx.Module(
		"Logging",

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
