package llm

import (
	"context"

	"github.com/IridiumNan/termhelper/internal/config"
	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/pkg"
)

type LLMClient interface {
	ExplainBatch(ctx context.Context, requests []*models.LLMRequest) (Responses []models.LLMExplainResults, err error)

	Ping(ctx context.Context) error
}

func NextProvider() (LLMClient, error) {
	modelName, err := config.NextProviderStr()
	if err != nil {
		return nil, pkg.ReportErrorInputErr("NextProvider", err)
	}
	switch modelName {
	case models.OllamaClientStr:
		return newOllamaClient(), nil

	case models.DeepSeekClientStr:
		return newDeepSeekClient(), nil

	case models.OpenAIClientStr:
		return newOpenAIClient(), nil
	}

	return nil, pkg.ReportErrorInputStr("NextProvider", " unknown option")
}

func NewPrompt(words []string, context string) *models.Prompt {
	return &models.Prompt{
		Role:    models.PromptRole,
		Task:    models.PromptTask,
		Format:  models.PromptFormat,
		Example: models.PromptExample,
		Words:   words,
		Context: context,
	}
}
