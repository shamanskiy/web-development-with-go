package models

import (
	"errors"
	"fmt"

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
	ErrImageNotFound    = errors.New("models: image is not found")
)

type FileError struct {
	Issue string
}

func (fe FileError) Error() string {
	return fmt.Sprintf("invalid file: %v", fe.Issue)
}

func isSqlUniqueViolation(err error) bool {
	var pgError *pgconn.PgError
	if errors.As(err, &pgError) {
		if pgError.Code == pgerrcode.UniqueViolation {
			return true
		}
	}
	return false
}
