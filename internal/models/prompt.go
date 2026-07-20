package models

import (
	_ "embed"
	"encoding/json"
)

// PromptMode controls which language the LLM uses to explain words.
type PromptMode string

const (
	PromptModeZh PromptMode = "zh2zh" // Chinese explanation (original behavior)
	PromptModeEn PromptMode = "en2en" // Simple English explanation
)

type Prompt struct {
	Role    string   `json:"role"`
	Task    string   `json:"task"`
	Format  string   `json:"format"`
	Example string   `json:"example"`
	Words   []string `json:"words"`
	Context string   `json:"context"`
}

func NewPrompt(words []string, context string, mode PromptMode) *Prompt {
	role := promptRole
	task := promptTask
	format := promptFormat
	example := promptExample

	if mode == PromptModeEn {
		role = promptEnRole
		task = promptEnTask
		format = promptEnFormat
		example = promptEnExample
	}

	return &Prompt{
		Role:    role,
		Task:    task,
		Format:  format,
		Example: example,
		Words:   words,
		Context: context,
	}
}

func (prompt *Prompt) Str() (string, error) {
	promptStr, err := json.Marshal(prompt)
	if err != nil {
		return EmptyStr, err
	}

	return string(promptStr), nil
}

// TODO: extrat this Prompt to config.toml
const (
	promptRole = `你是一位技术英语导师，专门帮助计算机专业学生理解英文技术文档中的词汇。`

	promptTask = `我会给你一段英文技术文本（可能是报错信息、手册片段或文档），以及一个需要解释的词汇列表。
请针对列表中的每个词，在该技术上下文语境下给出：
1. 简单释义：1-2 个中文词，如果存在多个常见释义，用中文逗号「，」分隔。
2. 详细解释：一段中文说明，解释该词在技术语境中的具体含义和作用。`

	promptFormat = `请严格按照以下 JSON 格式返回，不要添加任何额外文本或说明：
{
  "results": [
    {
      "word": "词汇原文",
      "simple_definition": "中文释义，多个用逗号分隔",
      "detailed_explanation": "详细解释"
    }
  ]
}
要求：
- results 数组中的顺序必须与输入词汇列表的顺序完全一致。
- 只解释指定的词汇，不要添加其他词汇。
- 确保 JSON 是有效的，不包含注释或多余逗号。`

	promptExample = `输入示例：
上下文文本：
mv: cannot move '/etc/os-release' to './os-release': Permission denied
词汇列表：["mv", "release", "permission", "denied"]

输出示例：
{
  "results": [
    {"word":"mv","simple_definition":"移动","detailed_explanation":"用于移动或重命名文件或目录的 Linux 命令"},
    {"word":"release","simple_definition":"发行版本","detailed_explanation":"指软件或系统的特定发布版本，如 /etc/os-release 文件包含当前系统版本信息"},
    {"word":"permission","simple_definition":"权限","detailed_explanation":"文件或目录的访问权限，决定用户能否读取、写入或执行"},
    {"word":"denied","simple_definition":"被拒绝","detailed_explanation":"操作未被允许，通常由于当前用户缺乏足够的文件系统权限"}
  ]
}`
)
