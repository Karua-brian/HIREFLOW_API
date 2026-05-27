package response

import (
	"encoding/json"
	"net/http"

)

// ErrorResponse represents the structure of an error response returned by the API.
type ErrorResponse struct {
	Message string `json:"message"` // A human-readable error message
	Fields []ValidationError `json:"fields,omitempty"` // Optional list of field-specific validation errors
}

// APIResponse represents the standard structure of a successful API response, which may include data and/or error information.
type APIResponse struct {
	Data interface{} `json:"data,omitempty"` // The actual response data (optional)
	Error *ErrorResponse `json:"error,omitempty"` // Error details if an error occurred (optional)
}

// ValidationError represents a specific validation error for a particular field, used when validating input data.
type ValidationError struct {
	Field string `json:"field"` // The name of the field that failed validation
	Error string `json:"error"` // A human-readable error message describing the validation failure
}

// JSON is a helper function to send a JSON response with the given status code and data.
func JSON(w http.ResponseWriter, statusCode int, data interface{}) {

	// Set the status code and content type header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := APIResponse{
		Data: data,
	}

	json.NewEncoder(w).Encode(resp)
}

// Error is a helper function to send a JSON error response with the given status code and error message.
func Error(w http.ResponseWriter, statusCode int, message string, fields ...ValidationError) {
	// Set the status code and content type header
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := APIResponse{
		Error: &ErrorResponse{
			Message: message,
			Fields: fields,
		},
	}

	json.NewEncoder(w).Encode(resp)
}

// DecodeJSON is a helper function to decode a JSON request body into the provided destination struct.
func DecodeJSON(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Prevent unknown fields in the JSON body

	return decoder.Decode(dst)
}