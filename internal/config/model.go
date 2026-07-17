package config

import (
	"log/slog"

	"github.com/IridiumNan/termhelper/internal/models"
)

type Config struct {
	Providers    []models.ModelName `mapstructure:"providers"`
	DeepSeek     DeepSeekConfig     `mapstructure:"deepseek"`
	OpenAI       OpenAIConfig       `mapstructure:"openai"`
	Ollama       OllamaConfig       `mapstructure:"ollama"`
	Test         TestConfig         `mapstructure:"test"`
	Logging      LoggingConfig      `mapstructure:"logging"`
	WordProvider WordProvider       `mapstructure:"word_provider"`
}

type WordProvider struct {
	DisabledWords []string `mapstructure:"disabled_words"`
}

type DeepSeekConfig struct {
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
	APIBase string `mapstructure:"api_base"`
}

type OpenAIConfig struct {
	APIKey  string `mapstructure:"api_key"`
	Model   string `mapstructure:"model"`
	APIBase string `mapstructure:"api_base"`
}

type OllamaConfig struct {
	URL        string        `mapstructure:"url"`
	Model      string        `mapstructure:"model"`
	OllamaOpts OllamaOptions `mapstructure:"options"`
}

type OllamaOptions struct {
	Temperature float32 `mapstructure:"temperature"`
	TopP        float32 `mapstructure:"top_p"`
	NumPredict  int     `mapstructure:"num_predict"`
	NumCtx      int     `mapstructure:"num_ctx"`
}

type TestConfig struct {
	BatchSize   float32         `mapstructure:"batch_size"`
	DailyLimit  int             `mapstructure:"daily_limit"`
	Proficiency TestProficiency `mapstructure:"proficiency"`
}

type TestProficiency struct {
	CorrectIncrement float32 `mapstructure:"correct_increment"`
	WrongDecrement   float32 `mapstructure:"wrong_decrement"`
	InitialValue     float32 `mapstructure:"initial_value"`
	MaxIntervalDays  int     `mapstructure:"max_interval_days"`
}

type LoggingConfig struct {
	ConsoleLevel slog.Level `mapstructure:"console_level"`
	FileLevel    slog.Level `mapstructure:"file_level"`
}

type LLMConfig struct {
	MaxWordsPerChunk uint16 `mapstructure:"max_words_per_chunk"`
	MaxChunkChars    uint16 `mapstructure:"max_chunk_chars"`
}
