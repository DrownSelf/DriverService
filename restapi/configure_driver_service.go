// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"io"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/DrownSelf/DriverService/restapi/operations"
	"github.com/DrownSelf/DriverService/restapi/operations/actions_with_token"
	"github.com/DrownSelf/DriverService/restapi/operations/logging_actions"
)

//go:generate swagger generate server --target ../../DriverService --name DriverService --spec ../swagger.yml --principal interface{} --exclude-main

func configureFlags(api *operations.DriverServiceAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.DriverServiceAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.ApplicationsJSONProducer = runtime.ProducerFunc(func(w io.Writer, data interface{}) error {
		return errors.NotImplemented("applicationsJson producer has not yet been implemented")
	})
	api.JSONProducer = runtime.JSONProducer()

	if api.ActionsWithTokenGetDriverHandler == nil {
		api.ActionsWithTokenGetDriverHandler = actions_with_token.GetDriverHandlerFunc(func(params actions_with_token.GetDriverParams) middleware.Responder {
			return middleware.NotImplemented("operation actions_with_token.GetDriver has not yet been implemented")
		})
	}
	if api.ActionsWithTokenGetDriverLogoutHandler == nil {
		api.ActionsWithTokenGetDriverLogoutHandler = actions_with_token.GetDriverLogoutHandlerFunc(func(params actions_with_token.GetDriverLogoutParams) middleware.Responder {
			return middleware.NotImplemented("operation actions_with_token.GetDriverLogout has not yet been implemented")
		})
	}
	if api.ActionsWithTokenPostDriverEndRideHandler == nil {
		api.ActionsWithTokenPostDriverEndRideHandler = actions_with_token.PostDriverEndRideHandlerFunc(func(params actions_with_token.PostDriverEndRideParams) middleware.Responder {
			return middleware.NotImplemented("operation actions_with_token.PostDriverEndRide has not yet been implemented")
		})
	}
	if api.LoggingActionsPostDriverLoginHandler == nil {
		api.LoggingActionsPostDriverLoginHandler = logging_actions.PostDriverLoginHandlerFunc(func(params logging_actions.PostDriverLoginParams) middleware.Responder {
			return middleware.NotImplemented("operation logging_actions.PostDriverLogin has not yet been implemented")
		})
	}
	if api.ActionsWithTokenPostDriverRateUserHandler == nil {
		api.ActionsWithTokenPostDriverRateUserHandler = actions_with_token.PostDriverRateUserHandlerFunc(func(params actions_with_token.PostDriverRateUserParams) middleware.Responder {
			return middleware.NotImplemented("operation actions_with_token.PostDriverRateUser has not yet been implemented")
		})
	}
	if api.LoggingActionsPostDriverRegisterHandler == nil {
		api.LoggingActionsPostDriverRegisterHandler = logging_actions.PostDriverRegisterHandlerFunc(func(params logging_actions.PostDriverRegisterParams) middleware.Responder {
			return middleware.NotImplemented("operation logging_actions.PostDriverRegister has not yet been implemented")
		})
	}
	if api.ActionsWithTokenPostDriverTakeOrderHandler == nil {
		api.ActionsWithTokenPostDriverTakeOrderHandler = actions_with_token.PostDriverTakeOrderHandlerFunc(func(params actions_with_token.PostDriverTakeOrderParams) middleware.Responder {
			return middleware.NotImplemented("operation actions_with_token.PostDriverTakeOrder has not yet been implemented")
		})
	}
	if api.ActionsWithTokenPostOrderRateDriverHandler == nil {
		api.ActionsWithTokenPostOrderRateDriverHandler = actions_with_token.PostOrderRateDriverHandlerFunc(func(params actions_with_token.PostOrderRateDriverParams) middleware.Responder {
			return middleware.NotImplemented("operation actions_with_token.PostOrderRateDriver has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
