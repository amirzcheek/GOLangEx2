package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
)

func main() {
	apiKey := "sk-proj-l9Bv3DyD55WiSaD0vMapT3BlbkFJsYtD4er0e7NJ4faZ29ep"
	apiEndpoint := "https://api.openai.com/v1/engines/gpt-3.5-turbo/completions"

	client := resty.New()

	for {
		// Get user input
		var userInput string
		fmt.Print("You: ")
		fmt.Scanln(&userInput)

		// Send user input to ChatGPT
		response, err := client.R().
			SetAuthToken(apiKey).
			SetHeader("Content-Type", "application/json").
			SetBody(map[string]interface{}{
				"model":      "gpt-3.5-turbo",
				"messages":   []interface{}{map[string]interface{}{"role": "user", "content": userInput}},
				"max_tokens": 512,
			}).
			Post(apiEndpoint)

		if err != nil {
			log.Fatalf("Error while sending the request: %v", err)
		}

		body := response.Body()

		var data map[string]interface{}
		err = json.Unmarshal(body, &data)
		if err != nil {
			fmt.Println("Error while decoding JSON response:", err)
			return
		}

		// Check if choices is not nil
		if choices, ok := data["choices"].([]interface{}); ok && len(choices) > 0 {
			// Extract the content from the JSON response
			if message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{}); ok {
				content := message["content"].(string)
				fmt.Println("ChatGPT:", content)
			}
		} else {
			fmt.Println("ChatGPT: Sorry, I couldn't understand that.")
		}

		// Exit the loop if the user says "exit"
		if strings.ToLower(userInput) == "exit" {
			break
		}
	}
}
