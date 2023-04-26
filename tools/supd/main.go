package main

import (
	"context"
	"os"
	"os/signal"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"

	"github.com/greenboxal/aip/pkg/collective/transports/slack"
	"github.com/greenboxal/aip/pkg/supervisor"
)

func main() {
	app := fx.New(
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

		fx.Provide(supervisor.NewManager),
		fx.Provide(slack.NewTransport),

		fx.Provide(NewRouting),
		fx.Provide(NewDaemon),

		fx.Invoke(func(d *Daemon) error {
			return d.Run()
		}),
	)

	go func() {
		signalCh := make(chan os.Signal, 1)

		signal.Notify(signalCh, os.Interrupt)

		<-signalCh

		if err := app.Stop(context.Background()); err != nil {
			panic(err)
		}
	}()

	app.Run()
}
