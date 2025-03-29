package main

import (
	"errors"
	"net/http"
	"time"

	"FernArchive/internal/data"
	"FernArchive/internal/validator"
)

func (bknd *backend) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := bknd.readJSON(w, r, &input)
	if err != nil {
		bknd.badRequestResponse(w, r, err)
		return
	}
	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}
	err = user.Password.SetPass(input.Password)
	if err != nil {
		bknd.serverErrorResponse(w, r, err)
		return
	}
	vldtr := validator.NewValidator()

	if data.ValidateUser(vldtr, user); !vldtr.Valid() {
		bknd.failedValidationResponse(w, r, vldtr.Errors)
		return
	}
	err = bknd.models.Users.InsertUser(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			vldtr.AddError("email", "email already taken")
			bknd.failedValidationResponse(w, r, vldtr.Errors)
		default:
			bknd.serverErrorResponse(w, r, err)
		}
		return
	}
	token, err := bknd.models.Tokens.NewToken(user.Id, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		bknd.serverErrorResponse(w, r, err)
		return
	}
	bknd.background(func() {
		userData := map[string]any{"activationToken": token.PlainText,
			"userId":   user.Id,
			"username": user.Name,
		}
		err = bknd.mailer.SendEmail(user.Email, "user_welcome.gohtml", userData)
		if err != nil {
			bknd.logger.Error(err.Error())
		}
	})
	err = bknd.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		bknd.serverErrorResponse(w, r, err)
	}
}
