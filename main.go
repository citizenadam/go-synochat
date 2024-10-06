package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"./api"
	"./models"
	"github.com/0ghny/go-synochat"
)

var (
	synoClient   *synochat.Client
	claudeClient *api.Client
)

func init() {
	var err error

	// Initialize Synology Chat client
	synoBaseURL := os.Getenv("SYNOLOGY_BASE_URL")
	if synoBaseURL == "" {
		log.Fatal("SYNOLOGY_BASE_URL environment variable is not set")
	}
	synoClient, err = synochat.NewClient(synoBaseURL)
	if err != nil {
		log.Fatalf("Failed to create Synology Chat client: %v", err)
	}

	// Initialize Claude API client
	claudeClient, err = api.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Claude API client: %v", err)
	}
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := os.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var synoMessage synochat.ChatMessage
	err = json.Unmarshal(body, &synoMessage)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusBadRequest)
		return
	}

	// Process the message with Claude API
	claudeResponse, err := processWithClaude(synoMessage.Text)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing with Claude: %v", err), http.StatusInternalServerError)
		return
	}

	// Send the response back to Synology Chat
	err = synoClient.SendMessage(&synochat.ChatMessage{Text: claudeResponse}, os.Getenv("SYNOLOGY_WEBHOOK_TOKEN"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error sending message to Synology Chat: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func processWithClaude(message string) (string, error) {
	request := models.MessageRequest{
		Model: "claude-3-opus-20240229", // Update this to the latest model version
		Messages: []models.Message{
			{Role: "user", Content: message},
		},
		MaxTokens: 1000,
	}

	response, err := claudeClient.SendMessage(request)
	if err != nil {
		return "", err
	}

	// Extract the text from the response
	var responseText string
	for _, content := range response.Content {
		if content.Type == "text" {
			responseText += content.Text
		}
	}

	return responseText, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/webhook", handleWebhook)

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
