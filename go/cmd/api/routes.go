package main

import (
	"expvar"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/api/v1/participant", app.listParticipantsHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/participant", app.createParticipantHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/participant/:uuid", app.showParticipantHandler)
	//router.HandlerFunc(http.MethodGet, "/api/v1/participant/:id", app.requirePermission("participants:read", app.showParticipantHandler))
	router.HandlerFunc(http.MethodPatch, "/api/v1/participant/:uuid", app.updateParticipantHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/participant/:uuid", app.requestPaymentParticipantHandler)
	router.HandlerFunc(http.MethodDelete, "/api/v1/participant/:id", app.requirePermission("participants:write", app.deleteParticipantHandler))

	router.HandlerFunc(http.MethodPost, "/api/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/api/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(SpaFileSystem(http.Dir(app.config.webroot))))
	mux.Handle("/api/", app.metrics(app.authenticate(router)))

	return mux
	//return app.metrics(app.authenticate(router))
}

type spaFileSystem struct {
	root http.FileSystem
}

func (fs *spaFileSystem) Open(name string) (http.File, error) {
	f, err := fs.root.Open(name)
	if os.IsNotExist(err) {
		return fs.root.Open("index.html")
	}
	return f, err
}
func SpaFileSystem(fs http.FileSystem) *spaFileSystem {
	return &spaFileSystem{root: fs}
}
