package app

import "net/http"

type APIMessage struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

var (
	SussessDelete  = &APIMessage{Status: http.StatusNotFound, Message: "success deleted"}
	SuccessCreated = &APIMessage{Status: http.StatusCreated, Message: "success created"}
	AlreadExist    = &APIMessage{Status: http.StatusConflict, Message: "already exist"}
	SucceesSignIn  = &APIMessage{Status: http.StatusOK, Message: "authorization success"}
)
