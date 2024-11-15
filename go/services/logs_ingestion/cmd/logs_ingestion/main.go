package main

import (
	"github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/config"
	logsstream "github.com/Magpie-Monitor/magpie-monitor/services/logs_ingestion/pkg/logs_stream"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.AppModule,
		fx.Invoke(func(*logsstream.LogsStreamListener) {}),
	).Run()
}
