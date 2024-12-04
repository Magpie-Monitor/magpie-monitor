package tests

import (
	"context"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

const (
	PRODUCTION_ENVIRONMENT = "prod"
	TEST_ENVIRONMENT       = "test"
)

func RunTest[T any](test func(dependencies T), t *testing.T, appModule fx.Option) {

	app := fxtest.New(
		t,
		appModule,
		fx.Invoke(test),
	)

	app.Start(context.Background())
	app.Stop(context.Background())
}
