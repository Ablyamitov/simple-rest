package wrapper

import (
	"log/slog"
	"os"
	"time"

	"github.com/google/uuid"
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
