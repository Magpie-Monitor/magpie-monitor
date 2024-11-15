package logsstream

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type LogsStreamListener struct {
	applicationLogsReader ApplicationLogsStreamReader
	nodeLogsReader        NodeLogsStreamReader
	logger                *zap.Logger
}

func NewLogsStreamListener(
	lc fx.Lifecycle,
	logger *zap.Logger,
	applicationLogsReader ApplicationLogsStreamReader,
	nodeLogsReader NodeLogsStreamReader,
) *LogsStreamListener {

	listener := LogsStreamListener{
		logger:                logger,
		applicationLogsReader: applicationLogsReader,
		nodeLogsReader:        nodeLogsReader,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {

			go listener.nodeLogsReader.Listen()
			go listener.applicationLogsReader.Listen()
			return nil
		},
	})

	return &listener
}
