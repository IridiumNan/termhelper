package llm

import (
	"context"
	"net/http"
	"time"

	"github.com/IridiumNan/termhelper/internal/models"
)

type openAICompatibleClient struct {
	apiKey    string
	modelName models.ModelName
	baseURL   string
	http      *http.Client
}

func newOpenAICompatibleClient(apiKey string, model models.ModelName, baseURL string) *openAICompatibleClient {
	return &openAICompatibleClient{
		apiKey:    apiKey,
		modelName: model,
		baseURL:   baseURL,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (client *openAICompatibleClient) explainBatch(ctx context.Context, request []*models.LLMRequest) ([]models.LLMExplainResults, error) {
	// TODO:
	//
	return nil, nil
}

func (client *openAICompatibleClient) ping(ctx context.Context) error {
	return nil
}
