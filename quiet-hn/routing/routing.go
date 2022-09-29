package routing

import (
	"net/http"

	"github.com/AksAman/gophercises/quietHN/views"
	"github.com/gorilla/mux"
)

type URLRoute struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
	Methods []string
}

var Routes = []URLRoute{
	{Pattern: "/", Handler: views.Home, Methods: []string{"GET"}},
	{Pattern: "/stories/", Handler: views.Stories, Methods: []string{"GET"}},
}

func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.NotFoundHandler = views.NotFoundHandler{}

	for _, route := range Routes {
		router.HandleFunc(route.Pattern, route.Handler).Methods(route.Methods...)
	}

	return router
}
