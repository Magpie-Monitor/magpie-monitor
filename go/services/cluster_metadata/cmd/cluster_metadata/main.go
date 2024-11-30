package main

import (
	"net/http"

	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/config"
	"github.com/Magpie-Monitor/magpie-monitor/services/cluster_metadata/pkg/services"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.AppModule,
		fx.Invoke(func(*http.Server, *services.MetadataService) {}),
	).Run()
}
