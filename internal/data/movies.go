package data

import (
	"database/sql"
	"errors"
	"time"

	"FernArchive/internal/validator"

	"github.com/lib/pq"
)

type Movie struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Runtime   Runtime   `json:"runtime,omitempty"`
	Year      int32     `json:"year,omitempty"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"`
}

type MovieModel struct {
	DB *sql.DB
}

func (mdl *MovieModel) Insert(movie *Movie) error {
	query := `INSERT INTO movies (title, year, runtime, genres)
		    VALUES ($1, $2, $3, $4)
		    RETURNING id, created_at, version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	return mdl.DB.QueryRow(query, args...).Scan(&movie.Id, &movie.CreatedAt, &movie.Version)
}

func (mdl *MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `SELECT id, created_at, title, year, runtime, genres, version FROM movies
                WHERE id = $1`
	movie := new(Movie)

	err := mdl.DB.QueryRow(query, id).Scan(&movie.Id, &movie.CreatedAt,
		&movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return movie, nil
}

func (mdl *MovieModel) Update(movie *Movie) error {
	query := `UPDATE movies SET title=$1, year=$2, runtime=$3, genres=$4, version=version+1
                WHERE id = $5 RETURNING version`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.Id}
	return mdl.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (mdl *MovieModel) Delete(id int64) error {
	return nil
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
