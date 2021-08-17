// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/ehabterra/flash_api/internal/externals"

	"github.com/go-openapi/swag"

	"github.com/ehabterra/flash_api/internal/services"

	"github.com/ehabterra/flash_api/internal/database"

	"github.com/ehabterra/flash_api/internal/handlers"

	"github.com/ehabterra/flash_api/internal/middlewares"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"

	"github.com/ehabterra/flash_api/api/restapi/operations"
)

//go:generate swagger generate server --target ../../api --name Flash --spec ../../docs/swagger.yml --principal interface{} --exclude-main

var flashConfig = struct {
	Datasource string `long:"datasource" description:"database host:port" default:"localhost:21212" env:"DATASOURCE"`
}{}

func configureFlags(api *operations.FlashAPI) {
	api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{{
		ShortDescription: "Config",
		LongDescription:  "Flash config",
		Options:          &flashConfig,
	}}
}

func configureAPI(api *operations.FlashAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.UseSwaggerUI()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	db := database.NewVoltDB(flashConfig.Datasource)
	bank := externals.NewBank()

	users := services.NewUsers(bank, db)
	rates := services.NewRates()

	// Applies when the "Authorization" header is set
	api.BearerAuth = middlewares.ValidateHeader

	api.HomeGetHandler = handlers.NewHomeHandler()

	api.UsersLoginHandler = handlers.NewUsersLoginHandler(users)
	api.UsersConnectHandler = handlers.NewUsersConnectHandler(users)
	api.UsersGetBalanceHandler = handlers.NewUsersGetBalanceHandler(users)
	api.UsersSendHandler = handlers.NewUsersSendHandler(users)
	api.UsersUploadHandler = handlers.NewUsersUploadHandler(users)

	api.RatesGetRatesHandler = handlers.NewRatesGetRatesHandler(rates)

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
