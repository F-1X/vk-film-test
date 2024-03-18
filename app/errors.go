package app

import (
	"net/http"
)

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return e.Message
}

func (e APIError) APIError() (int, string) {
	return e.Status, e.Message
}

var (
	ErrAuth       = &APIError{Status: http.StatusUnauthorized, Message: "invalid token"}
	ErrNotFound   = &APIError{Status: http.StatusNotFound, Message: "not found"}
	ErrBadRequest = &APIError{Status: http.StatusBadRequest, Message: "bad request"}
	ErrDuplicate  = &APIError{Status: http.StatusBadRequest, Message: "duplicate"}

	ErrActorNotFound = &APIError{Status: http.StatusNotFound, Message: "actor not found"}
	ErrFilmNotFound  = &APIError{Status: http.StatusNotFound, Message: "film not found"}

	ErrActorConflict = &APIError{Status: http.StatusConflict, Message: "actor already exist"}
	ErrActorDelete   = &APIError{Status: http.StatusNoContent, Message: "delete success"}

	ErrFilmConflict = &APIError{Status: http.StatusConflict, Message: "film already exist"}
	ErrUserConflict = &APIError{Status: http.StatusConflict, Message: "user already exist"}

	ErrCreated = &APIError{Status: http.StatusCreated, Message: "success created"}

	ErrSignUp = &APIError{Status: http.StatusCreated, Message: "failed signup"}

	ErrUnauthorized    = &APIError{Status: http.StatusUnauthorized, Message: "failed singin"}
	ErrSuccessUpdate   = &APIError{Status: http.StatusCreated, Message: "update success"}
	ErrNotConcretActor = &APIError{Status: http.StatusBadRequest, Message: "not concrect actor"}

	ErrBadPassword = &APIError{Status: http.StatusUnauthorized, Message: "bad password sorry"}

	ErrCookies = &APIError{Status: http.StatusBadRequest, Message: "bad cookies"}
)
