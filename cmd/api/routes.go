package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (bknd *backend) routes() http.Handler {
	router := httprouter.New()

	router.MethodNotAllowed = http.HandlerFunc(bknd.methodNotAllowedResponse)
	router.NotFound = http.HandlerFunc(bknd.notFoundResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", bknd.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", bknd.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", bknd.showMovieHandler)

	return bknd.recoverPanic(router)
}
