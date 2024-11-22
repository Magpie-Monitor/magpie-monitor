package handlers_test

import (
	"context"
	"github.com/Magpie-Monitor/magpie-monitor/pkg/tests"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/brokers"
	"github.com/Magpie-Monitor/magpie-monitor/services/reports/internal/handlers"
	config "github.com/Magpie-Monitor/magpie-monitor/services/reports/pkg/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"testing"
)

func TestApplicationLogsStreamReader(t *testing.T) {

	type TestDependencies struct {
		fx.In
		Logger        *zap.Logger
		ReportHandler *handlers.ReportsHandler
	}

	testCases := []struct {
		request     *brokers.ReportRequest
		expectedErr *brokers.ReportRequestFailed
	}{
		{
			request:     &brokers.ReportRequest{},
			expectedErr: &brokers.ReportRequestFailed{},
		},
	}

	test := func(dependencies TestDependencies) {
		if dependencies.ReportHandler == nil {
			t.Fatal("Failed to load report handler")
		}

		response := dependencies.ReportHandler.ScheduleReport(context.Background(), "3213", &brokers.ReportRequest{})

	}

	tests.RunTest(test, t, config.AppModule)
}
