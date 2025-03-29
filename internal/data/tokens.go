package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"FernArchive/internal/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserId    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func generateToken(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{UserId: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]
	return token, nil
}

func ValidateTokenPlainText(vldtr *validator.Validator, plainTxt string) {
	vldtr.Check(plainTxt != "", "plainTxt", "must be provided")
	vldtr.Check(len(plainTxt) == 26, "plainTxt", "must be 26 characters long")
}

type TokenModel struct {
	Db *sql.DB
}

func (mdl *TokenModel) NewToken(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = mdl.Insert(token)
	return token, err
}

func (mdl *TokenModel) Insert(token *Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) VALUES ($1, $2, $3, $4)`
	args := []any{token.Hash, token.UserId, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := mdl.Db.ExecContext(ctx, query, args...)
	return err
}

func (mdl *TokenModel) DeleteAllForUser(scope string, userId int64) error {
	query := `DELETE FROM tokens WHERE scope=$1 AND user_id=$2`
	args := []any{scope, userId}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := mdl.Db.ExecContext(ctx, query, args...)
	return err
}
