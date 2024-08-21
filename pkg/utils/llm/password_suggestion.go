package llm

import (
	"cli-project/internal/config"
	"cli-project/pkg/utils/math"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	generalPrompt = "This request is automated please respond carefully or it might break the system. Only generate output in the specified format. Give answers to the best of your ability. "
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

}

func PasswordSuggestion() string {
	apiKey := os.Getenv("OPENAI_API")
	client := resty.New()

	var (
		passwordPrompt = fmt.Sprintf("Task: Generate a very strong password of minimum 12 character length containing at least 1  uppercase, 1 lowercase, 1 digit  and 1 special character. Also you must use these digits %d in your password generation randomly. Output : Just return the password itself", math.RandomInt())
	)

	finalPrompt := generalPrompt + passwordPrompt

	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model":      config.GPT_MODEL,
			"messages":   []interface{}{map[string]interface{}{"role": "system", "content": finalPrompt}},
			"max_tokens": 50,
		}).
		Post(config.GPT_API_ENDPOINT)

	if err != nil {
		log.Fatalf("Error while sending the request: %v", err)
	}

	body := response.Body()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error while decoding JSON response:", err)
		return ""
	}

	content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	return content
}
