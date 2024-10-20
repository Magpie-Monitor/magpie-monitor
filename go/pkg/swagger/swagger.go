package swagger

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type SwaggerRouter struct {
	logger *zap.Logger
	router *mux.Router
}

func NewSwaggerRouter(swaggerHandler *SwaggerHandler, rootRouter *mux.Router) *SwaggerRouter {
	router := rootRouter.PathPrefix("/swagger").Subrouter()
	router.Methods(http.MethodGet).Path("/openapi.yaml").HandlerFunc(swaggerHandler.GetOpenApiSpecification)
	router.Methods(http.MethodGet).HandlerFunc(swaggerHandler.Get)

	return &SwaggerRouter{
		router: router,
	}
}

func (rtr *SwaggerRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rtr.router.ServeHTTP(w, r)
}

type SwaggerHandler struct {
	logger *zap.Logger
	config *SwaggerConfig
	ui     string
}

func NewSwaggerHandler(logger *zap.Logger, config *SwaggerConfig) *SwaggerHandler {

	swaggerUi, err := os.ReadFile("/srv/swagger-ui.html")
	if err != nil {
		panic(fmt.Sprintf("Failed to open swagger-ui file %s", err.Error()))
	}

	swaggerUiWithCustomHost := strings.ReplaceAll(string(swaggerUi), "SWAGGER_HOST", config.Host)

	return &SwaggerHandler{
		logger: logger,
		ui:     swaggerUiWithCustomHost,
		config: config,
	}
}

func (h *SwaggerHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(h.ui))
}
func (h *SwaggerHandler) GetOpenApiSpecification(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "/srv/openapi.yaml")
}

type SwaggerConfig struct {
	Host string
}

func ProvideSwaggerConfig() func() *SwaggerConfig {
	host, ok := os.LookupEnv("SWAGGER_HOST")
	if !ok {
		panic("SWAGGER_HOST is not set")
	}

	return func() *SwaggerConfig {
		return &SwaggerConfig{
			Host: host,
		}
	}
}
