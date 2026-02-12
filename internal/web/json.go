package web

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, err error, msg string) {
	type errorResponse struct {
		Message string `json:"error"`
	}

	log.Printf("error: %v\n", err)

	respondWithJSON(w, code, errorResponse{
		Message: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to encode response"))
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}
