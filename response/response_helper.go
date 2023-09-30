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

func Ok(w http.ResponseWriter, data interface{}) {
	JsonEncode(w, http.StatusOK, data)
}

func InteralServerError(w http.ResponseWriter, message string) {
	JsonEncode(w, http.StatusInternalServerError, message)
}
