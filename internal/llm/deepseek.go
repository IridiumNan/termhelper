package llm

import (
	"context"

	"github.com/IridiumNan/termhelper/internal/models"
)

type DeepSeekClient struct {
	impl *openAICompatibleClient
}

func (client DeepSeekClient) ExplainBatch(ctx context.Context, request []models.LLMRequest) (Responses []models.LLMResponse) {
}
