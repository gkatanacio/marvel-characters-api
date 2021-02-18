package errs

import (
	"errors"
	"net/http"
)

type HttpError interface {
	Error() string
	StatusCode() int
}

type BadRequest struct {
	error
}

func NewBadRequest(err string) *BadRequest {
	return &BadRequest{errors.New(err)}
}

func (e *BadRequest) StatusCode() int {
	return http.StatusBadRequest
}

type NotFound struct {
	error
}

func NewNotFound(err string) *NotFound {
	return &NotFound{errors.New(err)}
}

func (e *NotFound) StatusCode() int {
	return http.StatusNotFound
}

type BadGateway struct {
	error
}

func NewBadGateway(err string) *BadGateway {
	return &BadGateway{errors.New(err)}
}

func (e *BadGateway) StatusCode() int {
	return http.StatusBadGateway
}
