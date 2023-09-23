package models

import (
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
)

var (
	// users
	ErrEmailTaken    = errors.New("models: email address is already in use")
	ErrEmailNotFound = errors.New("models: email address is not found")
	ErrPasswordWrong = errors.New("models: password is wrong")
	ErrInvalidToken  = errors.New("models: invalid reset password token")

	// galleries
	ErrResourceNotFound = errors.New("models: resource not found")
)

func isSqlUniqueViolation(err error) bool {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		if pgError.Code == pgerrcode.UniqueViolation {
			return true
		}
	}
	return false
}
