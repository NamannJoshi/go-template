package utils

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetIdFromRequest(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, err
	}

	return id, nil
}
