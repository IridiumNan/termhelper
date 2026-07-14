package llm

import (
	"context"

	"github.com/IridiumNan/termhelper/internal/config"
	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/pkg"
)

type LLMClient interface {
	ExplainBatch(ctx context.Context, requests []models.LLMRequest) (Responses []models.LLMResponse, err error)

	Ping(ctx context.Context) error
}

func NextProvider() (*LLMClient, error) {
	modelName, err := config.NextProviderStr()
	if err != nil {
		return nil, pkg.ReportErrorInputErr("NextProvider", err)
	}
	switch modelName {
	case models.OllamaClientStr:
	}
}
