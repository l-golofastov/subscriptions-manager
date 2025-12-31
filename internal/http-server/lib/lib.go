package lib

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents error response
type ErrorResponse struct {
	ErrorMessage string `json:"error" example:"error message"`
}

// SuccessResponse represents success response
type SuccessResponse struct {
	Message string `json:"result" example:"success"`
}

func NewErrorResponse(error string) ErrorResponse {
	errorResponse := ErrorResponse{error}

	return errorResponse
}

func NewSuccessResponse(result string) SuccessResponse {
	successResponse := SuccessResponse{result}

	return successResponse
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload) //TODO: handle error
}

func RespondWithError(w http.ResponseWriter, statusCode int, errorMessage string) {
	RespondWithJSON(w, statusCode, NewErrorResponse(errorMessage))
}
