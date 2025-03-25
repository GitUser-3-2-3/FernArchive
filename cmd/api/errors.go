package main

import (
	"fmt"
	"net/http"
)

func (bknd *backend) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)
	bknd.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (bknd *backend) errorResponseJSON(w http.ResponseWriter, r *http.Request, status int, msg any) {
	env := envelope{"error": msg}
	err := bknd.writeJSON(w, status, env, nil)
	if err != nil {
		bknd.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (bknd *backend) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("%s method not supported for this request", r.Method)
	bknd.errorResponseJSON(w, r, http.StatusMethodNotAllowed, msg)
}

func (bknd *backend) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "requested resource not found"
	bknd.errorResponseJSON(w, r, http.StatusNotFound, msg)
}

func (bknd *backend) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	bknd.errorResponseJSON(w, r, http.StatusBadRequest, err.Error())
}

func (bknd *backend) failedValidationResponse(w http.ResponseWriter, r *http.Request,
	errs map[string]string,
) {
	bknd.errorResponseJSON(w, r, http.StatusUnprocessableEntity, errs)
}

func (bknd *backend) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	msg := "unable to update record due an edit conflict, please try again"
	bknd.errorResponseJSON(w, r, http.StatusConflict, msg)
}

func (bknd *backend) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	bknd.logError(r, err)
	msg := "server encountered a problem and could not process your request"
	bknd.errorResponseJSON(w, r, http.StatusInternalServerError, msg)
}
