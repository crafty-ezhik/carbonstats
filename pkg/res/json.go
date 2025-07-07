package res

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, data any, statusCode int) *http.ResponseWriter {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return nil
	}
	return &w
}
