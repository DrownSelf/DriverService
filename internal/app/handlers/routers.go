package swagger

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "DriverService/cmd/docs"
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
				"ApiV1DriverLogOut",
				strings.ToUpper("Get"),
				"/api/v1/driver/logout",
				handler.ApiV1DriverLogOut,
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
	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:80/swagger/doc.json"), //The url pointing to API definition
	)).Methods(http.MethodGet)

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
