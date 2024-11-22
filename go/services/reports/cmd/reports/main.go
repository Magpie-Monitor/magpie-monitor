package main

import (
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/handlers"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/services"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/config"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/openai"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.AppModule,
		fx.Invoke(func(
			reportsService *services.ReportsService,
			batchPoller *openai.BatchPoller,
			reportsHandler *handlers.ReportsHandler) {

			// Listen for ReportRequested messages
			go reportsHandler.ListenForReportRequests()

			// Poll for reports pending generation
			go reportsHandler.PollReports()

			// Poll for pending OpenAi batches
			go batchPoller.Start()
		},
		),
	).Run()
}
