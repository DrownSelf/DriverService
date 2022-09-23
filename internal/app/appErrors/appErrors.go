package appErrors

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Res struct {
	Message string `json:"message"`
	Code    int    `json:"status"`
}

var (
	ErrInvalidSession    = errors.New("Session is expired. please log in system again")
	ErrDriverExists      = errors.New("Driver with this phone exists")
	ErrMethodNotAllowed  = errors.New("Method not allowed")
	ErrDriverDoesntExist = errors.New("Driver doensn't exists")
	ErrWrongPassword     = errors.New("Wrong password")
	ErrInvalidData       = errors.New("Invalid data input")
)

func HandleErr(w http.ResponseWriter, err error) {
	var res Res
	switch err {
	case ErrInvalidSession:
		res = Res{ErrInvalidSession.Error(), http.StatusUnauthorized}
	case ErrDriverExists:
		res = Res{ErrDriverExists.Error(), http.StatusConflict}
	case ErrMethodNotAllowed:
		res = Res{ErrMethodNotAllowed.Error(), http.StatusMethodNotAllowed}
	case ErrWrongPassword:
		res = Res{ErrWrongPassword.Error(), http.StatusBadRequest}
	case ErrDriverDoesntExist:
		res = Res{ErrDriverDoesntExist.Error(), http.StatusBadRequest}
	case ErrInvalidData:
		res = Res{ErrInvalidData.Error(), http.StatusBadRequest}
	default:
		res = Res{"", http.StatusInternalServerError}
	}
	log.Println(res)
	w.WriteHeader(res.Code)
	_ = json.NewEncoder(w).Encode(res.Message)
}
