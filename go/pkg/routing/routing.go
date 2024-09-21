package routing

import (
	"go.uber.org/fx"
	"net/http"
)

type Route interface {
	http.Handler
	Pattern() string
}

func NewServeMux(routes []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, route := range routes {
		mux.Handle(route.Pattern(), route)
	}
	return mux
}

func ProvideAsRootServeMux(f any) any {
	return fx.Annotate(
		NewServeMux,
		fx.ParamTags(`group:"routes"`),
	)
}

func ProvideAsRoute(f any) any {

	return fx.Annotate(
		f,
		fx.As(new(Route)),
		fx.ResultTags(`group:"routes"`),
	)
}
