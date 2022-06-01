package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/ihrk/rest-api-task/internal/app/apperrors"
	"github.com/ihrk/rest-api-task/internal/app/database"
)

var ErrInvalidRequest = errors.New("invalid request")

func ParseRequestJSON(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		return ErrInvalidRequest
	}

	return nil
}

func RespondJSON(w http.ResponseWriter, statusCode int, v interface{}) {
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		log.Println(err)
	}
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func RespondError(w http.ResponseWriter, err error) {
	if err == nil {
		return
	}

	statusCode := http.StatusInternalServerError

	if errors.Is(err, database.ErrNotFound) {
		statusCode = http.StatusNotFound
	} else if errors.Is(err, ErrInvalidRequest) {
		statusCode = http.StatusBadRequest
	} else if errors.Is(err, apperrors.ErrUnauthorized) {
		statusCode = http.StatusUnauthorized
	}

	RespondJSON(w, statusCode, &ErrorResponse{Message: err.Error()})
}

func GetIDVar(r *http.Request, name string) (int64, error) {
	strVal := mux.Vars(r)[name]

	id, err := strconv.Atoi(strVal)

	if err != nil {
		return 0, ErrInvalidRequest
	}

	return int64(id), nil
}
