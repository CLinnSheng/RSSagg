package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// a function that convert the payload to json
// and send it back to the client
func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	// convert the payload to json
	data, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to marshal JSON response")
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5xx error: " + msg)
	}

	type errResponse struct {
		// specifying the key object will be error in the json format
		Error string `json:"error"`
	}

	respondWithJson(w, code, errResponse{
		Error: msg,
	})
}
