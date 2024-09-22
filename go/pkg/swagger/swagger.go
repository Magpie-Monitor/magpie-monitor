package swagger

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"strings"
)

type SwaggerRouter struct {
	logger *zap.Logger
	mux    *http.ServeMux
}

func NewSwaggerRouter(swaggerHandler *SwaggerHandler) *SwaggerRouter {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /swagger/openapi.yaml", swaggerHandler.GetOpenApiSpecification)
	mux.HandleFunc("GET /", swaggerHandler.Get)

	return &SwaggerRouter{
		mux: mux,
	}
}

func (r *SwaggerRouter) Pattern() string {
	return "/swagger/"
}

func (router *SwaggerRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
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
