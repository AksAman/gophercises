package routing

import (
	"net/http"

	"github.com/AksAman/gophercises/quietHN/controllers"
	"github.com/gorilla/mux"
)

type URLRoute struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
	Methods []string
}

var Routes = []URLRoute{
	{Pattern: "/", Handler: controllers.Home, Methods: []string{"GET"}},
	{Pattern: "/stories/", Handler: controllers.Stories, Methods: []string{"GET"}},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = controllers.NotFoundHandler{}

	for _, route := range Routes {
		router.HandleFunc(route.Pattern, route.Handler).Methods(route.Methods...)
	}

	return router
}
