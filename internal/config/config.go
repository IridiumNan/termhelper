package config

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/IridiumNan/termhelper/internal/models"
	"github.com/IridiumNan/termhelper/pkg"
	"github.com/spf13/viper"
)

//go:embed config.toml
var defaultConfigTemplate string

const (
	maxProviderIdx = 2
)

var (
	globalConfig         Config
	currentProviderIdx   int = 0
	defaultDisabledWords     = []string{
		"a",
		"about",
		"after",
		"again",
		"all",
		"also",
		"an",
		"and",
		"any",
		"are",
		"as",
		"at",
		"back",
		"be",
		"because",
		"been",
		"before",
		"being",
		"between",
		"both",
		"but",
		"by",
		"can",
		"could",
		"did",
		"do",
		"does",
		"doing",
		"down",
		"during",
		"each",
		"even",
		"every",
		"few",
		"first",
		"for",
		"from",
		"had",
		"has",
		"have",
		"having",
		"he",
		"her",
		"here",
		"him",
		"his",
		"how",
		"if",
		"in",
		"into",
		"is",
		"it",
		"its",
		"just",
		"last",
		"like",
		"many",
		"may",
		"me",
		"might",
		"more",
		"most",
		"much",
		"must",
		"my",
		"next",
		"no",
		"not",
		"now",
		"of",
		"off",
		"on",
		"once",
		"one",
		"only",
		"or",
		"other",
		"our",
		"out",
		"over",
		"quite",
		"rather",
		"said",
		"same",
		"shall",
		"she",
		"should",
		"so",
		"some",
		"still",
		"such",
		"than",
		"that",
		"the",
		"their",
		"them",
		"then",
		"there",
		"these",
		"they",
		"this",
		"those",
		"three",
		"through",
		"to",
		"two",
		"under",
		"unless",
		"until",
		"up",
		"us",
		"very",
		"was",
		"we",
		"well",
		"were",
		"what",
		"when",
		"where",
		"whether",
		"which",
		"while",
		"who",
		"whom",
		"whose",
		"why",
		"will",
		"with",
		"without",
		"would",
		"yes",
		"yet",
		"you",
		"your",
	}
)

func LoadConfig() error {
	funcName := "LoadConfig"
	if err := pkg.EnsureDir(models.GetConfigDirPath()); err != nil {
		return pkg.ReportErrorInputErr(funcName+" ensure config dir", err)
	}

	if _, err := os.Stat(models.GetConfigFilePath()); os.IsNotExist(err) {
		if err := os.WriteFile(models.GetConfigFilePath(), []byte(defaultConfigTemplate), 0o644); err != nil {
			return pkg.ReportErrorInputErr(funcName+" wirte default config template", err)
		}
		fmt.Println("config file " + models.GetConfigFilePath() + " not found, use the default one")
		fmt.Println("you can edit it -> " + models.GetConfigFilePath())

	}

	viper.SetConfigName(models.GetConfigRawName())
	viper.SetConfigType(models.GetConfigType())
	viper.AddConfigPath(models.GetConfigDirPath())

	viper.AutomaticEnv()
	viper.SetEnvPrefix("TERMHELPER")

	setDefaults()

	if err := viper.ReadInConfig(); err != nil {
		return pkg.ReportErrorInputErr(funcName+" read config", err)
	}

	if err := viper.Unmarshal(&globalConfig); err != nil {
		return pkg.ReportErrorInputErr(funcName+"Unmarshal file to globalConfig", err)
	}

	return nil
}

func setDefaults() {
	viper.SetDefault("providers", []string{"ollama", "deepseek", "openai"})

	viper.SetDefault("deepseek.model", "deepseek-v4-flash")
	viper.SetDefault("deepseek.api_base", "https://api.deepseek.com")

	viper.SetDefault("openai.model", "gpt-3.5-turbo")
	viper.SetDefault("openai.api_base", "https://api.openai.com/v1")

	viper.SetDefault("ollama.url", "http://localhost:11434")
	viper.SetDefault("ollama.model", "frob/qwen3.5-instruct:4b")
	viper.SetDefault("ollama.options.temperature", 0.1)
	viper.SetDefault("ollama.options.top_p", 0.9)
	viper.SetDefault("ollama.options.num_predict", 4096)
	viper.SetDefault("ollama.options.num_ctx", 8192)

	viper.SetDefault("test.batch_size", 30)
	viper.SetDefault("test.daily_limit", 100)

	viper.SetDefault("test.proficiency.correct_increment", 0.12)
	viper.SetDefault("test.proficiency.wrong_decrement", 0.10)
	viper.SetDefault("initial_value", 0.15)
	viper.SetDefault("max_interval_days", 30)

	viper.SetDefault("prompt_mode", string(models.PromptModeZh))

	viper.SetDefault("logging.console_level", 0)
	viper.SetDefault("logging.file_level", -4)

	viper.SetDefault("word_provider", defaultDisabledWords)
}

func PrintConfig() {
	fmt.Println(globalConfig)
}

func GetGlobalConfig() Config {
	return globalConfig
}

func NextProviderStr() (models.ModelName, error) {
	if currentProviderIdx > maxProviderIdx {
		return "", fmt.Errorf("all provider has been tried yet fail")
	}
	currentProviderIdx = currentProviderIdx + 1

	return globalConfig.Providers[currentProviderIdx-1], nil
}

func HasNextProvider() bool {
	return currentProviderIdx <= maxProviderIdx
}
