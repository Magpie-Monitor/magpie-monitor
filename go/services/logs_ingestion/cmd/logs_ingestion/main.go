package main

import (
	// "github.com/Magpie-Monitor/magpie-monitor/pkg/tests"
	"github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/config"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.AppModule,
		fx.Invoke(func(*config.LogsStreamListener) {}),
	).Run()
}
