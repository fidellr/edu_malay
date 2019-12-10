package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fidellr/edu_malay/model"
	"github.com/sirupsen/logrus"
)

var (
	// ErrContextNil is thrown if the context passed is nil to any function that require to use context.
	ErrContextNil = errors.New("Context is Nil")

	// ErrNotFound is thrown if the item requested does not exist in the system
	ErrNotFound = errors.New("Your requested item does not exists")

	// ErrNotModified is thrown to the client when the cached copy of a particular file is up to date with the server.
	ErrNotModified = errors.New("")
)

// ConstraintError represents a custom error for a contstraint things.
type ConstraintError string

func (e ConstraintError) Error() string {
	return string(e)
}

// ConstraintErrorf constructs ConstraintError with formatted message.
func ConstraintErrorf(format string, a ...interface{}) ConstraintError {
	return ConstraintError(fmt.Sprintf(format, a...))
}

// ErrorFromResponseStatusCode generates error based on the status code from *http.Response.
// For example, it will generate ErrNotFound when given status code of 404.
func ErrorFromResponseStatusCode(code int, message string) (err error) {
	switch code {
	case http.StatusNotFound:
		err = ErrNotFound
	case http.StatusBadRequest:
		err = ConstraintErrorf(message)
	case http.StatusNotModified:
		err = ErrNotModified
	default:
		err = fmt.Errorf(message)
	}

	return
}

// GetStatusCode get status code for the given error
func GetStatusCode(err error) int {

	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)

	switch err {
	case model.INTERNAL_SERVER_ERROR:
		return http.StatusInternalServerError
	case model.NOT_FOUND_ERROR:
		return http.StatusNotFound
	case model.CONFLICT_ERROR:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
