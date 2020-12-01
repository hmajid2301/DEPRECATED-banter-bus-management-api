package serverModels

// APIError is a generic error message returned by the API.
type APIError struct {
	Message string `json:"message" description:"The error message returned to the client." example:"Game quibly not found."`
}
