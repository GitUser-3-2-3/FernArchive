package main

import (
	"context"
	"net/http"

	"FernArchive/internal/data"
)

type contextKey string

const userContextKey = contextKey("user")

func (bknd *backend) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (bknd *backend) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("user missing request context")
	}
	return user
}
