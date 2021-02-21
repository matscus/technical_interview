package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/matscus/technical_interview/users"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := users.GetAllUsers()
	if err != nil {
		WriteHTTPError(w, http.StatusInternalServerError, err)
	}
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		WriteHTTPError(w, http.StatusInternalServerError, err)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, ok := params["id"]
	if ok {
		uuid, err := uuid.Parse(id)
		if err != nil {
			WriteHTTPError(w, http.StatusOK, errors.New("uuid format is not valide"))
			return
		}
		user, err := users.GetUser(uuid)
		if err != nil {
			WriteHTTPError(w, http.StatusInternalServerError, err)
			return
		}
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			WriteHTTPError(w, http.StatusInternalServerError, err)
			return
		}
		return
	}
	WriteHTTPError(w, http.StatusOK, errors.New("params id not found"))
}
