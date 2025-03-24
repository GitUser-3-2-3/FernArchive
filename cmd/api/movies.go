package main

import (
	"errors"
	"fmt"
	"net/http"

	"FernArchive/internal/data"
	"FernArchive/internal/validator"
)

func (bknd *backend) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Runtime data.Runtime `json:"runtime"`
		Year    int32        `json:"year"`
		Genres  []string     `json:"genres"`
	}
	err := bknd.readJSON(w, r, &input)
	if err != nil {
		bknd.badRequestResponse(w, r, err)
		return
	}
	vldtr := validator.NewValidator()
	movie := &data.Movie{
		Title:   input.Title,
		Runtime: input.Runtime,
		Year:    input.Year,
		Genres:  input.Genres,
	}
	if data.ValidateMovie(vldtr, movie); !vldtr.Valid() {
		bknd.failedValidationResponse(w, r, vldtr.Errors)
		return
	}
	err = bknd.models.Movies.Insert(movie)
	if err != nil {
		bknd.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.Id))

	err = bknd.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		bknd.serverErrorResponse(w, r, err)
	}
}

func (bknd *backend) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r)
	if err != nil {
		bknd.notFoundResponse(w, r)
		return
	}
	movie, err := bknd.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			bknd.notFoundResponse(w, r)
		default:
			bknd.serverErrorResponse(w, r, err)
		}
		return
	}
	err = bknd.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		bknd.serverErrorResponse(w, r, err)
	}
}

func (bknd *backend) updateMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := bknd.readIdParam(r)
	if err != nil {
		bknd.notFoundResponse(w, r)
		return
	}
	movie, err := bknd.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			bknd.notFoundResponse(w, r)
		default:
			bknd.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	err = bknd.readJSON(w, r, &input)
	if err != nil {
		bknd.badRequestResponse(w, r, err)
		return
	}
	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	vldtr := validator.NewValidator()
	if data.ValidateMovie(vldtr, movie); !vldtr.Valid() {
		bknd.failedValidationResponse(w, r, vldtr.Errors)
		return
	}
	err = bknd.models.Movies.Update(movie)
	if err != nil {
		bknd.serverErrorResponse(w, r, err)
		return
	}
	err = bknd.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		bknd.serverErrorResponse(w, r, err)
	}
}
