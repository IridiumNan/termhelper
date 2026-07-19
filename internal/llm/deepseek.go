package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/IridiumNan/termhelper/internal/config"
	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/internal/output"
	"github.com/cohesion-org/deepseek-go"
)

//
// type LLMClient interface {
// 	Explain(ctx context.Context, chunk *models.Chunk) (response *models.LLMExplainResults, err error)
// 	Ping(ctx context.Context) error
// }

type DeepSeekClient struct {
	client *deepseek.Client

	model string
}

func NewDeepSeekClient() *DeepSeekClient {
	apiKey := config.GetGlobalConfig().DeepSeek.APIKey

	model := config.GetGlobalConfig().DeepSeek.Model

	return &DeepSeekClient{
		client: deepseek.NewClient(apiKey),
		model:  model,
	}
}

func (d *DeepSeekClient) Explain(ctx context.Context, chunk *models.Chunk) (response *models.LLMExplainResults, err error) {
	prompt := models.NewPrompt(chunk.Words, chunk.Context)

	promptStr, err := prompt.Str()
	if err != nil {
		output.Console.Error("error when generate prompt", "error", err)
		output.File.Error("error when generate prompt", "error", err)
		return
	}

	req := &deepseek.ChatCompletionRequest{
		Model: d.model,
		Messages: []deepseek.ChatCompletionMessage{
			{Role: deepseek.ChatMessageRoleSystem, Content: promptStr},
		},
	}

	resp, err := d.client.CreateChatCompletion(ctx, req)
	if err != nil {
		output.Console.Error("error when create chat with deepseek", "error", err)
		output.File.Error("error when create chat with deepseek", "error", err)
		return
	}

	fmt.Println("raw choice[0].message.content => ", resp.Choices[0].Message.Content)

	rawJSON := resp.Choices[0].Message.Content

	if strings.HasPrefix(rawJSON, "```json") {
		rawJSON = strings.TrimPrefix(rawJSON, "```json")
		rawJSON = strings.TrimSuffix(rawJSON, "```")
	} else if strings.HasSuffix(rawJSON, "```") {
		rawJSON = strings.TrimPrefix(rawJSON, "```")
		rawJSON = strings.TrimSuffix(rawJSON, "```")
	}

	rawJSON = strings.TrimSpace(rawJSON)

	var results models.LLMExplainResults

	if err = json.Unmarshal([]byte(rawJSON), &results); err != nil {
		return nil, fmt.Errorf("parse JSON: %w, raw: %s", err, rawJSON)
	}

	return &results, nil
}

func (d *DeepSeekClient) Ping(ctx context.Context) error {
	url := d.client.BaseURL
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)

	if resp.StatusCode == http.StatusUnauthorized {
		return nil
	}

	defer resp.Body.Close()

	return err
}
