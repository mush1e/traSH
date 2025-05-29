package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mush1e/traSH/config"
)

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []interface{} `json:"messages"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

func chat(prompt, apiKey string) (string, error) {
	reqBody := chatRequest{
		Model: "gpt-3.5-turbo",
		Messages: []interface{}{
			message{Role: "system", Content: "you are a DERANGED AND UNHINGED SHELL CHAT COMPANION"},
			message{Role: "user", Content: prompt},
		},
	}
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(reqBody); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", buf)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("OpenAI error: %s", string(body))
	}

	var cr chatResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return "", err
	}
	if len(cr.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}
	return cr.Choices[0].Message.Content, nil
}

func HandleAI(cmd *Command) error {
	prompt := strings.Join(cmd.args, " ")
	key := config.GetConfig().GetAPIKey()
	if key == "" {
		return fmt.Errorf("no OpenAI key configured; set `openai_key` in your ~/.trashrc`")
	}

	resp, err := chat(prompt, key)
	if err != nil {
		return fmt.Errorf("AI error: %v", err)
	}
	fmt.Println(resp)
	return nil
}

func HandleExplain(cmd *Command) error {
	toExplain := strings.Join(cmd.args, " ")
	prompt := fmt.Sprintf("Explain what this shell command does: %s", toExplain)
	key := config.GetConfig().GetAPIKey()
	if key == "" {
		return fmt.Errorf("no OpenAI key configured; set `openai_key` in your ~/.trashrc`")
	}

	resp, err := chat(prompt, key)
	if err != nil {
		return fmt.Errorf("AI error: %v", err)
	}
	fmt.Println(resp)
	return nil
}
