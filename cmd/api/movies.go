package main

import (
	"fmt"
	"net/http"
	"time"

	"FernArchive/internal/data"
)

func (bknd *backend) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int          `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genre   []string     `json:"genres"`
	}
	err := bknd.readJSON(w, r, &input)
	if err != nil {
		bknd.badRequestResponse(w, r, err)
		return
	}
	_, _ = fmt.Fprintf(w, "%+v\n", input)
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
