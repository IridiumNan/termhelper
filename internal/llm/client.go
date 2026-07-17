package llm

import (
	"context"
	"fmt"

	"github.com/IridiumNan/termhelper/internal/config"
	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/pkg"
)

type LLMClient interface {
	Explain(ctx context.Context, chunk *models.Chunk) (response *models.LLMExplainResults, err error)
	Ping(ctx context.Context) error
}

func NextClient() (LLMClient, error) {
	if !config.HasNextProvider() {
		return nil, fmt.Errorf("there is no available provider")
	}

	nextProviderStr, _ := config.NextProviderStr()

	switch nextProviderStr {
	case models.OllamaClientStr:
		return newOllamaClient(), nil
		// case models.DeepSeekClientStr:
	}

	return nil, pkg.ReportErrorInputStr("NextClient", "fail to get provider")
}
