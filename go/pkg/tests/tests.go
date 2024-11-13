package tests

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"testing"
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

	app.RequireStart()
	app.RequireStop()
}
