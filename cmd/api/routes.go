package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	post  = http.MethodPost
	put   = http.MethodPut
	patch = http.MethodPatch
	get   = http.MethodGet
	dlete = http.MethodDelete
)

func (bknd *backend) routes() http.Handler {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(bknd.methodNotAllowedResponse)
	router.NotFound = http.HandlerFunc(bknd.notFoundResponse)

	router.HandlerFunc(get, "/v1/healthcheck", bknd.healthcheckHandler)

	router.HandlerFunc(post, "/v1/movies", bknd.requirePermission("movies:write", bknd.createMovieHandler))
	router.HandlerFunc(patch, "/v1/movies/:id", bknd.requirePermission("movies:write", bknd.updateMovieHandler))
	router.HandlerFunc(get, "/v1/movies", bknd.requirePermission("movies:read", bknd.listMovieHandler))
	router.HandlerFunc(get, "/v1/movies/:id", bknd.requirePermission("movies:read", bknd.showMovieHandler))
	router.HandlerFunc(dlete, "/v1/movies/:id", bknd.requirePermission("movies:write", bknd.deleteMovieHandler))

	router.HandlerFunc(post, "/v1/users", bknd.registerUserHandler)
	router.HandlerFunc(put, "/v1/users/activated", bknd.activateUserHandler)

	router.HandlerFunc(post, "/v1/tokens/authentication", bknd.createAuthTokenHandler)

	return bknd.recoverPanic(bknd.rateLimiter(bknd.authenticate(router)))
}
