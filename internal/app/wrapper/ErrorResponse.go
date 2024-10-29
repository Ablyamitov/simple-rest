package wrapper

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Response struct {
	ErrorID   string `json:"error_id"`
	Message   string `json:"message"`
	Method    string `json:"method"`
	Timestamp string `json:"timestamp"`
}

func NewErrorResponse(err string, method string) *Response {
	return &Response{
		ErrorID:   uuid.New().String(),
		Message:   err,
		Method:    method,
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

func SendError(w http.ResponseWriter, statusCode int, err error, method string) {
	response := NewErrorResponse(err.Error(), method)
	log.Printf("Error ID: %s, Message: %s, Method: %s, Timestamp: %s",
		response.ErrorID, response.Message, response.Method, response.Timestamp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Printf("Failed to encode error wrapper response: %v", err)
	}
}

func LogError(err string, method string) {
	logger := setupLogger()
	response := NewErrorResponse(err, method)
	logger.Error("Error occurred",
		slog.String("error_id", response.ErrorID),
		slog.String("message", response.Message),
		slog.String("method", response.Method),
		slog.String("timestamp", response.Timestamp),
	)
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
