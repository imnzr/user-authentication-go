package handler

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	logger *zap.Logger
}

func NewBaseHandler(logger *zap.Logger) *BaseHandler {
	return &BaseHandler{logger: logger}
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    int    `json:"code"`
}

// SendJSON sends a JSON response
func (h *BaseHandler) SendJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", zap.Error(err))
	}
}

// SendSuccess sends a success response
func (h *BaseHandler) SendSuccess(w http.ResponseWriter, data interface{}, message string) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	h.SendJSON(w, http.StatusOK, response)
}

// SendCreated sends a created response
func (h *BaseHandler) SendCreated(w http.ResponseWriter, data interface{}, message string) {
	response := Response{
		Success: true,
		Message: message,
		Data:    data,
	}
	h.SendJSON(w, http.StatusCreated, response)
}

// SendError sends an error response
func (h *BaseHandler) SendError(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponse{
		Success: false,
		Error:   message,
		Code:    statusCode,
	}
	h.SendJSON(w, statusCode, response)
}
