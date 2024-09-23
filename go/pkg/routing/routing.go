package routing

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route interface {
	http.Handler
	Pattern() string
}

func NewRootRouter() *mux.Router {
	mux := mux.NewRouter().PathPrefix("/v1").Subrouter()

	return mux
}
