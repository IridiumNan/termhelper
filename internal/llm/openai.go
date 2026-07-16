package llm

import (
	"context"

	"github.com/IridiumNan/termhelper/internal/config"
	"github.com/IridiumNan/termhelper/internal/models"
)

type OpenAIClient struct {
	impl *openAICompatibleClient
}

func newOpenAIClient() *OpenAIClient {
	return &OpenAIClient{
		impl: newOpenAICompatibleClient(
			config.GetGlobalConfig().OpenAI.APIKey,
			models.ModelName(config.GetGlobalConfig().OpenAI.Model),
			config.GetGlobalConfig().OpenAI.APIBase,
		),
	}
}

func (client *OpenAIClient) ExplainBatch(ctx context.Context, requests []*models.LLMRequest) (Responses []models.LLMExplainResults, err error) {
	return client.impl.explainBatch(ctx, requests)
}

func (client *OpenAIClient) Ping(ctx context.Context) error {
	return client.impl.ping(ctx)
}
