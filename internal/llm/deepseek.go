package llm

import (
	"context"

	"github.com/IridiumNan/termhelper/internal/config"
	"github.com/IridiumNan/termhelper/internal/models"
)

type DeepSeekClient struct {
	impl *openAICompatibleClient
}

func newDeepSeekClient() *DeepSeekClient {
	return &DeepSeekClient{
		impl: newOpenAICompatibleClient(
			config.GetGlobalConfig().DeepSeek.APIKey,
			models.ModelName(config.GetGlobalConfig().DeepSeek.Model),
			config.GetGlobalConfig().DeepSeek.APIBase,
		),
	}
}

func (client *DeepSeekClient) ExplainBatch(ctx context.Context, requests []*models.LLMRequest) (Responses []models.LLMExplainResults, err error) {
	return client.impl.explainBatch(ctx, requests)
}

func (client *DeepSeekClient) Ping(ctx context.Context) error {
	return client.impl.ping(ctx)
}
