package llm

import (
	"context"
	"fmt"
	"time"

	"github.com/IridiumNan/termhelper/internal/config"
	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/internal/output"
)

type LLMClient interface {
	// Explain function explain words on the chunk and return Results which can be parased and push to database
	Explain(ctx context.Context, chunk *models.Chunk) (response *models.LLMExplainResults, err error)

	// Ping function test if this LLM Provider is available
	// Just check the network
	Ping(ctx context.Context) error
}

func GetClient() (client LLMClient, err error) {
	if !config.HasNextProvider() {
		return nil, fmt.Errorf("there is no available provider")
	}

	for {
		nextProviderStr, _ := config.NextProviderStr()

		switch nextProviderStr {
		case models.OllamaClientStr:
			client = newOllamaClient()
		case models.DeepSeekClientStr:
			client = NewDeepSeekClient()
			// TODO: OpenAI client (Optional)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		defer cancel()

		err = client.Ping(ctx)
		if err == nil {
			return
		} else {
			output.Console.Info("fail to ping model provider, try next", "provider", nextProviderStr, "err", err)
			output.File.Info("fail to ping model provider, try next", "provider", nextProviderStr, "err", err)
		}

	}
}
