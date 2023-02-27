package response

import (
	"encoding/json"
	"net/http"
)

func JsonEncode[T any](w http.ResponseWriter, statusCode int, data T) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	json.NewEncoder(w).Encode(data)
}