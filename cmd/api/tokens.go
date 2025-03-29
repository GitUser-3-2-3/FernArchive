package main

import (
	"errors"
	"net/http"
	"time"

	"FernArchive/internal/data"
	"FernArchive/internal/validator"
)

func (bknd *backend) createAuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := bknd.readJSON(w, r, &input)
	if err != nil {
		bknd.badRequestResponse(w, r, err)
		return
	}
	vldtr := validator.NewValidator()

	data.ValidatePasswordPlainTxt(vldtr, input.Password)
	data.ValidateEmail(vldtr, input.Email)
	if !vldtr.Valid() {
		bknd.failedValidationResponse(w, r, vldtr.Errors)
		return
	}
	user, err := bknd.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			bknd.invalidCredentialsResponse(w, r)
		default:
			bknd.serverErrorResponse(w, r, err)
		}
		return
	}
	correct, err := user.Password.CheckPass(input.Password)
	if err != nil {
		bknd.serverErrorResponse(w, r, err)
		return
	}
	if !correct {
		bknd.invalidCredentialsResponse(w, r)
		return
	}
	token, err := bknd.models.Tokens.NewToken(user.Id, 360*time.Hour, data.ScopeActivation)
	if err != nil {
		bknd.serverErrorResponse(w, r, err)
		return
	}
	err = bknd.writeJSON(w, http.StatusCreated, envelope{"authentication_token": token}, nil)
	if err != nil {
		bknd.serverErrorResponse(w, r, err)
	}
}
