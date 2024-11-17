package helper

import (
	"encoding/json"
	"net/http"
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
	encoder := json.NewEncoder(w)
	err := encoder.Encode(result)
	if err != nil {
		panic(err)
	}
}
