package main

import (
	"fmt"
	"net/http"
	"time"

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
	movie := data.Movie{
		Id:        id,
		CreatedAt: time.Now(),
		Title:     "Infinity War",
		Runtime:   102,
		Genres:    []string{"action", "sci-fi", "war"},
		Version:   1,
	}
	err = bknd.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		bknd.serverErrorResponse(w, r, err)
	}
}

func (bknd *backend) updateMovieHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "update a movie")
}
