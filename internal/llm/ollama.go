package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/IridiumNan/termhelper/internal/config"
	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/pkg"
)

const (
	defaultContentType = "application/json"
	ollamaEndPoint     = "/api/generate"

	jsonPrefix = "```json"
	jsonSuffix = "```"
)

type OllamaClient struct {
	URL        string
	Model      string
	Prompt     string
	ollamaOpts config.OllamaOptions

	http *http.Client
}

func newOllamaClient() *OllamaClient {
	return &OllamaClient{
		URL:        config.GetGlobalConfig().Ollama.URL,
		Model:      config.GetGlobalConfig().Ollama.Model,
		Prompt:     models.EmptyPrompt,
		ollamaOpts: config.GetGlobalConfig().Ollama.OllamaOpts,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (client *OllamaClient) ExplainBatch(ctx context.Context, requests []*models.LLMRequest) (results []models.LLMExplainResults, err error) {
	if len(requests) == 0 {
		return []models.LLMExplainResults{}, nil
	}

	results = make([]models.LLMExplainResults, 0)

	for i := range requests {
		context := requests[i].Context
		words := requests[i].Word

		prompt := NewPrompt(words, context)

		promptStr, _ := json.Marshal(prompt)

		ollamaReq := map[string]any{
			"model":  client.Model,
			"prompt": promptStr,
			"stream": false,
			"options": map[string]any{
				"temperature": client.ollamaOpts.Temperature,
				"top_p":       client.ollamaOpts.TopP,
				"num_predict": client.ollamaOpts.NumPredict,
				"num_ctx":     client.ollamaOpts.NumCtx,
			},
		}

		jsonData, err := json.Marshal(ollamaReq)
		if err != nil {
			return nil, pkg.ReportErrorInputErr("marshual ollama request", err)
		}

		url := client.URL + ollamaEndPoint

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonData))
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}

		req.Header.Set("Content-Type", defaultContentType)

		resp, err := client.http.Do(req)
		if err != nil {
			return nil, fmt.Errorf("http request: %w", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)

			return nil, fmt.Errorf("ollama API error status: %d, body: %s", resp.StatusCode, string(body))
		}

		var ollamaResp struct {
			Response   string `json:"response"`
			Done       bool   `json:"done"`
			DoneReason string `json:"dong_reason"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
			return nil, fmt.Errorf("decode response: %w", err)
		}

		if pkg.IsEmptyStr(ollamaResp.Response) {
			return nil, fmt.Errorf("empty response from ollama")
		}

		rawJson := strings.TrimSpace(ollamaResp.Response)

		if strings.HasPrefix(rawJson, jsonPrefix) {
			rawJson = strings.TrimPrefix(rawJson, jsonPrefix)
			rawJson = strings.TrimSuffix(rawJson, jsonSuffix)
		} else if strings.HasSuffix(rawJson, jsonSuffix) {
			rawJson = strings.TrimPrefix(rawJson, jsonSuffix)
			rawJson = strings.TrimSuffix(rawJson, jsonSuffix)
		}

		var responses []models.LLMExplainResponse
		if err := json.Unmarshal([]byte(rawJson), &responses); err != nil {
			var single models.LLMExplainResponse

			if errSingle := json.Unmarshal([]byte(rawJson), &single); errSingle == nil {
				responses = []models.LLMExplainResponse{single}
			} else {
				return nil, fmt.Errorf("parse JSON: %w, raw : %s", err, rawJson)
			}
		}

		results = append(results, models.LLMExplainResults{
			Results: responses,
		})
	}

	return results, nil
}

func (client *OllamaClient) Ping(ctx context.Context) error {
	url := client.URL + "/api/tags"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.http.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama not ready status %d", resp.StatusCode)
	}
	return nil
}
