package models

// APIError is a generic error message returned by the API.
type APIError struct {
	Message string `json:"message"`
}
