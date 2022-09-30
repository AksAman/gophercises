package routing

import (
	"net/http"

	"github.com/AksAman/gophercises/quietHN/controllers"
	"github.com/AksAman/gophercises/quietHN/middlewares"
	"github.com/gorilla/mux"
)

type URLRoute struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
	Methods []string
}

var Routes = []URLRoute{
	{Pattern: "/", Handler: controllers.Home, Methods: []string{"GET"}},
	{Pattern: "/stories", Handler: controllers.Stories, Methods: []string{"GET"}},
	{Pattern: "/panic", Handler: controllers.FakeError, Methods: []string{"GET"}},
	{Pattern: "/panic-after", Handler: controllers.FakeErrorAfter, Methods: []string{"GET"}},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = controllers.NotFoundHandler{}

	for _, route := range Routes {
		router.HandleFunc(route.Pattern, route.Handler).Methods(route.Methods...)
	}

	router.Use(middlewares.RecoveryMiddleware)
	return router
}
