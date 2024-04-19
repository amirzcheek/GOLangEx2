package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

func main() {
	apiKey := "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"
	client := openai.NewClient(apiKey)

	// Variables for rate limiting
	retryDelay := 20 * time.Second
	maxRetries := 3
	retryCount := 0
	var prevInput string

	// Open a log file for writing
	logFile, err := os.OpenFile("request_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logFile.Close()
	logger := log.New(logFile, "", log.LstdFlags)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			// Get user input from the form
			userInput := r.FormValue("question")

			// Check if the input contains keywords
			filterWords := []string{"alcohol", "18+", "drugs"}
			for _, word := range filterWords {
				if strings.Contains(userInput, word) {
					// Decline the request if it contains any filter word
					fmt.Fprint(w, "Declined: Your request was declined because your question is not suitable for children.")
					return
				}
			}

			logger.Printf("Request: %s\n", userInput)
			// Concatenate with previous input if it doesn't contain keywords
			userInput = prevInput + " " + userInput

			// Send user input to ChatGPT
			resp, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model: openai.GPT3Dot5Turbo,
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleUser,
							Content: userInput,
						},
					},
				},
			)

			if err != nil {
				if strings.Contains(err.Error(), "Rate limit reached") {
					// Retry logic
					retryCount++
					if retryCount > maxRetries {
						log.Fatalf("Max retry limit reached. Exiting.")
						return
					}
					fmt.Printf("Rate limit reached. Waiting for %v before retrying.\n", retryDelay)
					time.Sleep(retryDelay)
					return
				}
				log.Fatalf("ChatCompletion error: %v", err)
				return
			}

			// Reset retry count on successful response
			retryCount = 0

			// Print the response
			var response string
			if len(resp.Choices) > 0 {
				response = resp.Choices[0].Message.Content
			} else {
				response = "Sorry, I couldn't understand that."
			}

			// Display the response on the web page
			fmt.Fprint(w, response)

		} else {
			http.ServeFile(w, r, "index.html")
		}
	})

	log.Println("Starting server on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
