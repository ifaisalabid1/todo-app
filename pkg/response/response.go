package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	Success   bool      `json:"success"`
	Message   string    `json:"message"`
	Data      any       `json:"data,omitzero"`
	Error     string    `json:"error,omitzero"`
	Timestamp time.Time `json:"timestamp"`
}

func Success(w http.ResponseWriter, status int, message string, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		fmt.Fprintf(w, "failed to respond with json: %v\n", err)
	}
}

func Error(w http.ResponseWriter, status int, message string, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := Response{
		Success:   false,
		Message:   message,
		Error:     err,
		Timestamp: time.Now().UTC(),
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		fmt.Fprintf(w, "failed to respond with json: %v\n", err)
	}
}
