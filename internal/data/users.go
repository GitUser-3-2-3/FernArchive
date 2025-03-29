package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"FernArchive/internal/validator"

	"golang.org/x/crypto/bcrypt"
)

const constraintUniqueEmail = `pq: duplicate key value violates unique constraint "users_email_key"`

var ErrDuplicateEmail = errors.New("duplicate email")

type User struct {
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

type password struct {
	plainText *string
	hash      []byte
}

func (pass *password) SetPass(plainTxtPass string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTxtPass), 12)
	if err != nil {
		return err
	}
	pass.plainText = &plainTxtPass
	pass.hash = hash
	return nil
}

func (pass *password) CheckPass(plainTxtPass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(pass.hash, []byte(plainTxtPass))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(vldtr *validator.Validator, email string) {
	vldtr.Check(email != "", "email", "must be provided")
	vldtr.Check(validator.Matches(email, validator.EmailRX), "email", "invalid email address")
}

func ValidatePassPlainTxt(vldtr *validator.Validator, pass string) {
	vldtr.Check(pass != "", "password", "must be provided")
	vldtr.Check(len(pass) >= 8, "password", "must be greater than 8 chars")
	vldtr.Check(len(pass) <= 72, "password", "must be lesser than 72 chars")
}

func ValidateUser(vldtr *validator.Validator, user *User) {
	vldtr.Check(user.Name != "", "name", "must be provided")
	vldtr.Check(len(user.Name) <= 500, "name", "must be less than 500 chars")

	ValidateEmail(vldtr, user.Email)

	if user.Password.plainText != nil {
		ValidatePassPlainTxt(vldtr, *user.Password.plainText)
	}
	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}

type UserModel struct {
	Db *sql.DB
}

func (mdl *UserModel) InsertUser(user *User) error {
	query := `INSERT INTO users (name, email, password_hash, activated) VALUES ($1, $2, $3, $4)
                RETURNING id, created_at, version`

	args := []any{user.Name, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := mdl.Db.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == constraintUniqueEmail:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (mdl *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, created_at, name, email, password_hash, activated, version FROM users 
                WHERE email = $1`
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := mdl.Db.QueryRowContext(ctx, query, email).Scan(&user.Id,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (mdl *UserModel) UpdateUser(user *User) error {
	query := `UPDATE users SET name=$1, email=$2, password_hash=$3, activated=$4, version=version+1
                WHERE id = $5 AND version = $6 RETURNING version`

	args := []any{user.Name, user.Email, user.Password.hash, user.Activated, user.Id, user.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := mdl.Db.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		case err.Error() == constraintUniqueEmail:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (mdl *UserModel) GetForToken(scope, plainTxt string) (*User, error) {
	hash := sha256.Sum256([]byte(plainTxt))

	query := `SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
		    FROM users INNER JOIN tokens ON users.id = tokens.user_id
		    WHERE tokens.hash = $1 AND tokens.scope = $2 AND tokens.expiry > $3`

	args := []any{hash[:], scope, time.Now()}
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := mdl.Db.QueryRowContext(ctx, query, args...).Scan(&user.Id,
		&user.CreatedAt, &user.Name, &user.Email, &user.Password.hash, &user.Activated, &user.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
