package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/IridiumNan/termhelper/internal/config"
	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/pkg"
	"github.com/ollama/ollama/api"
)

type OllamaClient struct {
	url        string
	model      string
	httpClient *http.Client
	opts       config.OllamaOptions
}

func newOllamaClient() *OllamaClient {
	cfg := config.GetGlobalConfig()

	return &OllamaClient{
		url:        cfg.Ollama.URL,
		model:      cfg.Ollama.Model,
		httpClient: &http.Client{Timeout: 60 * time.Second},
		opts:       cfg.Ollama.OllamaOpts,
	}
}

func (c *OllamaClient) Explain(ctx context.Context, chunk *models.Chunk) (responses *models.LLMExplainResults, err error) {
	url, err := url.Parse(c.url)
	if err != nil {
		return nil, err
	}
	client := api.NewClient(url, c.httpClient)

	prompt := models.NewPrompt(chunk.Words, chunk.Context)

	promptStr, err := prompt.Str()
	if err != nil {
		return nil, pkg.ReportErrorInputErr("convert prompt to str", err)
	}

	results := &models.LLMExplainResults{}

	req := &api.GenerateRequest{
		Model:  c.model,
		Prompt: promptStr,
		Stream: new(bool),
		Options: map[string]any{
			"temperature": c.opts.Temperature,
			"top_p":       c.opts.TopP,
			"num_predict": c.opts.NumPredict,
			"num_ctx":     c.opts.NumCtx,
		},
	}
	respFunc := func(resp api.GenerateResponse) error {
		err := json.Unmarshal([]byte(resp.Response), results)

		return err
	}

	err = client.Generate(ctx, req, respFunc)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (c *OllamaClient) LegacyExplain(ctx context.Context, chunk *models.Chunk) (responses *models.LLMExplainResults, err error) {
	if len(chunk.Words) == 0 {
		return nil, nil
	}

	prompt := models.NewPrompt(chunk.Words, chunk.Context)

	promptStr, err := prompt.Str()
	if err != nil {
		return nil, pkg.ReportErrorInputErr("convert prompt to str", err)
	}

	ollamaReq := map[string]interface{}{
		"model":  c.model,
		"prompt": promptStr,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": c.opts.Temperature,
			"top_p":       c.opts.TopP,
			"num_predict": c.opts.NumPredict,
			"num_ctx":     c.opts.NumCtx,
		},
	}

	jsonData, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, pkg.ReportErrorInputErr("marshal ollama request", err)
	}

	url := strings.TrimRight(c.url, "/") + "/api/generate"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
	if err != nil {
		return nil, pkg.ReportErrorInputErr("create request ", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, pkg.ReportErrorInputErr("http request", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var ollamaResp struct {
		Response   string `json:"response"`
		Done       bool   `json:"done"`
		DoneReason string `json:"done_reason"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, fmt.Errorf("decode response %w", err)
	}

	rawJSON := strings.TrimSpace(ollamaResp.Response)

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

func (c *OllamaClient) Ping(ctx context.Context) error {
	url := strings.TrimRight(c.url, "/") + "/api/tags"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama not ready: status: %d", resp.StatusCode)
	}

	return nil
}
