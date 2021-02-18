package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gkatanacio/marvel-characters-api/internal/errs"
)

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	jsonEncoder := json.NewEncoder(w)
	jsonEncoder.Encode(data)
}

type errorResponseBody struct {
	Error string `json:"error"`
}

func errorResponse(w http.ResponseWriter, err error) {
	body := &errorResponseBody{}
	body.Error = err.Error()

	var status int
	switch e := err.(type) {
	case errs.HttpError:
		status = e.StatusCode()
	default:
		status = http.StatusInternalServerError
	}

	jsonResponse(w, body, status)
}
