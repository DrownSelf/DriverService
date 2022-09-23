package swagger

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"DriverService/internal/app/middlewares"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes struct {
	routes  []Route
	handler *Handler
}

func NewRoutes(handler *Handler) *Routes {
	return &Routes{
		handler: handler,
		routes: []Route{
			Route{
				"Index",
				"GET",
				"/",
				Index,
			},

			Route{
				"ApiV1DriverGet",
				strings.ToUpper("Get"),
				"/api/v1/driver",
				handler.ApiV1DriverGet,
			},

			Route{
				"ApiV1DriverRatingGet",
				strings.ToUpper("Get"),
				"/api/v1/driver/rating",
				handler.ApiV1DriverRatingGet,
			},

			Route{
				"ApiV1DriverFeedbackRidePost",
				strings.ToUpper("Post"),
				"/api/v1/driver/feedbackRide",
				handler.ApiV1DriverFeedbackRidePost,
			},

			Route{
				"ApiV1DriverLoginPost",
				strings.ToUpper("Post"),
				"/api/v1/driver/login",
				handler.ApiV1DriverLoginPost,
			},

			Route{
				"ApiV1DriverRegisterPost",
				strings.ToUpper("Post"),
				"/api/v1/driver/register",
				handler.ApiV1DriverRegisterPost,
			},

			Route{
				"ApiV1DriverStatusPost",
				strings.ToUpper("Put"),
				"/api/v1/driver/status",
				handler.ApiV1DriverStatusPost,
			},
		},
	}
}

func NewRouter(h *Handler) *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range NewRoutes(h).routes {
		var handler http.Handler
		handler = route.HandlerFunc
		handler = middlewares.Logger(handler, route.Name)
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
