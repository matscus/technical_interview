package handlers

import (
	"encoding/json"
	"net/http"

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
