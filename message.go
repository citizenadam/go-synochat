package synochat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

type ChatMessage struct {
	Text string `json:"text"`
	// FileUrl string `json:"file_url,omitempty"`
}

func (c *Client) SendMessage(msg *ChatMessage, token string) error {
	endpoint, err := url.JoinPath(c.BaseURL.String(), "/webapi/entry.cgi")
	if err != nil {
		return fmt.Errorf("failed to construct URL: %w", err)
	}

	queryParams := url.Values{
		"api":     {"SYNO.Chat.External"},
		"method":  {"incoming"},
		"version": {"2"},
		"token":   {token},
	}

	jsonData, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	formData := url.Values{
		"payload": {string(jsonData)},
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.URL.RawQuery = queryParams.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	log.Debugf("response status code: %d", resp.StatusCode)

	apiResponse, err := NewAPIResponseFromHTTPResponse(resp)
	if err != nil {
		return fmt.Errorf("failed to parse API response: %w", err)
	}

	if !apiResponse.Success {
		return fmt.Errorf("synochat API response error, name: %s - reason: %s",
			apiResponse.Error.Errors.Name, apiResponse.Error.Errors.Reason)
	}

	return nil
}
