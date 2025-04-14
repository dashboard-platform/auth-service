// Package handler provides HTTP handlers for the authentication service.
// These handlers process incoming requests, interact with the service layer, and return responses.
package handler

// Response represents the structure of an HTTP response.
type Response struct {
	Error bool        `json:"error"` // Indicates whether the response contains an error.
	Data  interface{} `json:"data"`  // The response data or error message.
}
