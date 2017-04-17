package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type apiErrorResponse struct {
	status int
	error  error
}

func apiError(rw http.ResponseWriter, statusCode int, err error) {
	if statusCode == http.StatusInternalServerError {
		log.Printf("ERROR: %v\n", err)
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)

	resp := apiErrorResponse{statusCode, err}
	if err := json.NewEncoder(rw).Encode(&resp); err != nil {
		log.Printf("Error encoding apiErrorResponse: %v\n", resp)
	}
}
