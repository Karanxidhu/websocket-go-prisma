package helper

import (
	"encoding/json"
	"net/http"

	"github.com/karanxidhu/go-websocket/data/response"
)

func ReadRequest(r *http.Request, result interface{}) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(result)
	if err != nil {
		panic(err)
	}
}

func WriteResponse(w http.ResponseWriter, result interface{}) {
	w.Header().Set("Content-Type", "application/json")

	var statusCode int

	// Check if the result is a WebResponse struct
	if res, ok := result.(response.WebResponse); ok {
		statusCode = res.Code // Extract the status code from the struct
	} else {
		// Default to 500 if the type assertion fails
		statusCode = http.StatusInternalServerError
	}

	// Set the HTTP status code
	w.WriteHeader(statusCode)

	// Encode the response as JSON
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		// Replace panic with proper error handling in production
		panic(err)
	}
}
