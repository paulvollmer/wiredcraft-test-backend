// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"
	graceful "github.com/tylerb/graceful"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	"github.com/paulvollmer/wiredcraft-test-backend/models"
	"github.com/paulvollmer/wiredcraft-test-backend/restapi/operations"
	"github.com/paulvollmer/wiredcraft-test-backend/restapi/operations/user"
)

var (
	fileLogger lumberjack.Logger
	db         Database
	dbFile     = "server.db"
)

func configureFlags(api *operations.WiredcraftTestBackendAPI) {}

func configureAPI(api *operations.WiredcraftTestBackendAPI) http.Handler {
	api.ServeError = errors.ServeError

	// initialize the logger
	// a real api server need to log the incoming requests
	fileLogger = lumberjack.Logger{
		Filename:   "./server.log",
		MaxSize:    500,
		MaxBackups: 3,
		MaxAge:     28,
	}
	log.SetOutput(&fileLogger)
	api.Logger = log.Printf
	fmt.Printf("Logger initialized %q\n", fileLogger.Filename)

	// initialize the database
	// as database we use the embedded database boltDB
	db, err := NewDatabase(dbFile, 0644)
	if err != nil {
		fmt.Println("Database Error", err)
		os.Exit(127)
	}
	fmt.Printf("Database initialize %q\n", dbFile)

	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "x-token" header is set
	api.KeyAuth = func(token string) (interface{}, error) {
		// this set up the token only for demonstration
		if token == "secret" {
			return token, nil
		}
		return nil, errors.New(401, "incorrect api key auth")
	}

	api.UserUserIDDeleteHandler = user.UserIDDeleteHandlerFunc(func(params user.UserIDDeleteParams, principal interface{}) middleware.Responder {
		err := db.DeleteUser(params.Userid)
		if err != nil {
			return user.NewUserIDDeleteNotFound().WithPayload(&status404)
		}
		status := true
		return user.NewUserIDDeleteOK().WithPayload(user.UserIDDeleteOKBody{Status: &status})
	})

	api.UserUserIDGetHandler = user.UserIDGetHandlerFunc(func(params user.UserIDGetParams, principal interface{}) middleware.Responder {
		u, err := db.ReadUser(params.Userid)
		if err != nil {
			return user.NewUserIDGetNotFound().WithPayload(&status404)
		}
		return user.NewUserIDGetOK().WithPayload(u)
	})

	api.UserUserIDPutHandler = user.UserIDPutHandlerFunc(func(params user.UserIDPutParams, principal interface{}) middleware.Responder {
		d, err := db.UpdateUser(params.Userid, *params.Data)
		if err != nil {
			return user.NewUserIDPutNotFound().WithPayload(&status404)
		}
		return user.NewUserIDPutOK().WithPayload(d)
	})

	api.UserUserPostHandler = user.UserPostHandlerFunc(func(params user.UserPostParams, principal interface{}) middleware.Responder {
		created, err := db.CreateUser(*params.Data)
		if err != nil {
			return user.NewUserPostInternalServerError().WithPayload(&status500)
		}
		return user.NewUserPostCreated().WithPayload(created)
	})

	api.ServerShutdown = func() {
		// here we can graceful disconnect from other services
		log.Println("graceful shutdown the database")
		defer db.Close()
		log.Println("graceful shutdown the server")
		fileLogger.Rotate()
	}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *graceful.Server, scheme string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s\n", r.Method, r.URL.Path)
		fmt.Printf("%s %s\n", r.Method, r.URL.Path)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Server", "wiredcraft-user-api")
		handler.ServeHTTP(w, r)
	})
}

var (
	statuscode400 int64 = 400
	statusmsg400        = "Bad Request"
	status400           = models.ModelError{
		Statuscode: &statuscode400,
		Status:     &statusmsg400,
	}

	statuscode404 int64 = 404
	statusmsg404        = "Not Found"
	status404           = models.ModelError{
		Statuscode: &statuscode404,
		Status:     &statusmsg404,
	}

	statuscode500 int64 = 500
	statusmsg500        = "Internal Server Error"
	status500           = models.ModelError{
		Statuscode: &statuscode500,
		Status:     &statusmsg500,
	}
)
