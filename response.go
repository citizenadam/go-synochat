package synochat

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Success bool `json:"success"`
	Error   struct {
		Code   int `json:"code"`
		Errors struct {
			Name   string `json:"name"`
			Reason string `json:"reason"`
		} `json:"errors,omitempty"`
	} `json:"error,omitempty"`
}

// NewAPIResponseFromHTTPResponse creates an APIResponse from an http.Response object
func NewAPIResponseFromHTTPResponse(resp *http.Response) (*APIResponse, error) {
	defer resp.Body.Close()

	var response APIResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
