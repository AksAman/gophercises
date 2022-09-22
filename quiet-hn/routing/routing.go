package routing

import (
	"net/http"

	"github.com/AksAman/gophercises/quietHN/views"
)

type URLRoute struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
}

var Routes = []URLRoute{
	{Pattern: "/", Handler: views.Home},
	{Pattern: "/stories/", Handler: views.Stories},
}
