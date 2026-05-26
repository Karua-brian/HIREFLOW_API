package response

import (
	"encoding/json"
	"net/http"

)

// ErrorResponse represents the structure of an error response returned by the API.
type ErrorResponse struct {
	Message string `json:"message"` // A human-readable error message
}

// APIResponse represents the standard structure of a successful API response, which may include data and/or error information.
type APIResponse struct {
	Data interface{} `json:"data,omitempty"` // The actual response data (optional)
	Error *ErrorResponse `json:"error,omitempty"` // Error details if an error occurred (optional)
}

// JSON is a helper function to send a JSON response with the given status code and data.
func JSON(w http.ResponseWriter, statusCode int, message string) {

	// Set the status code and content type header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := APIResponse{
		Error: &ErrorResponse{
			Message: message,
		},
	}

	json.NewEncoder(w).Encode(resp)
}