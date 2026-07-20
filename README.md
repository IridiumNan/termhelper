# termhelper — 终端英语词汇辅助学习工具

[![Go Version](https://img.shields.io/badge/Go-1.26%2B-blue)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

**termhelper** 是一个终端 CLI 工具，帮助计算机专业学生和技术人员在日常工作中学习英文技术词汇。直接从报错信息、技术文档、终端输出等真实文本中提取词汇，借助 LLM 生成中文释义，并通过**间隔重复（Spaced Repetition）** 机制进行复习巩固。

---

## 功能特性

- 📖 **文本解析** — 从命令行参数、管道输入或文件读取英文技术文本
- 🤖 **多 LLM 支持** — 支持 Ollama（本地）、DeepSeek API、OpenAI API，按优先级自动切换
- 🎯 **智能提取** — 自动识别文本中的英文单词，跳过常见词和已掌握词汇
- 📚 **间隔复习** — 基于熟练度的间隔重复算法，到期自动提醒复习
- ✍️ **交互测验** — 四选一单词释义测验，答对/答错自动调整熟练度
- 📊 **状态统计** — 查看待复习单词数和各熟练度级别的单词分布
- 🔧 **可配置** — 支持通过 TOML 配置文件自定义 LLM 参数、复习策略等

---

## DEMO

[demo link](https://repo.waterman.xin/share/termhelper_demo.mp4)

---

## 安装

### 前置要求

- Go 1.26+
- 至少配置一种 LLM 后端（推荐使用 Ollama 本地运行，或 DeepSeek API）
- Linux

### Download from release link (recommanded)

[v0.0.1-github-link](https://github.com/IridiumNan/termhelper/releases/download/v0.0.1/termhelper-linux-amd64)
[v0.0.1-gitee-link](https://gitee.com/cai-zixiang_hainan/termhelper/releases/download/v0.0.1/termhelper-linux-amd64)

### 从源码编译

```bash
# 克隆仓库
git clone https://github.com/IridiumNan/termhelper.git
cd termhelper

# 编译当前平台
make build

# 或手动编译
go build -o bin/termhelper .
```

### 跨平台编译

```bash
# 编译所有支持平台
make build-all

# 编译指定平台
make build-linux
make build-darwin
make build-windows
```

编译产物位于 `bin/` 目录。
建议将 `bin/termhelper` 添加到 `PATH` 中以便全局使用。

---

## 快速开始

### 1. 初始化配置

首次运行会自动生成配置文件，或手动编辑：

```bash
termhelper config
```

配置文件位于 `~/.config/termhelper/config.toml`，主要需要配置：

- **Ollama**（推荐）：确保本地 Ollama 服务运行，默认地址 `http://localhost:11434`
- **DeepSeek**：设置 `api_key`（也可通过环境变量 `DEEPSEEK_API_KEY` 设置）
- **OpenAI**：设置 `api_key`（也可通过环境变量 `OPENAI_API_KEY` 设置）

### 2. 添加文本并提取词汇

```bash
# 直接传入文本
termhelper add "mv: cannot move '/etc/os-release' to './os-release': Permission denied"

# 从管道输入
echo "fatal: Not a valid object name" | termhelper add

# 从文件读取
termhelper add -f error.log
```

工具会自动将文本分块调用 LLM，提取其中的技术词汇并生成中文释义。

### 3. 查看单词状态

```bash
termhelper status
```

显示待复习单词数和各熟练度级别的分布：

- 🟢 **New** — 新词
- 🟡 **Learning** — 学习中
- 🔵 **Familiar** — 较熟悉
- 🟣 **Master** — 已掌握

### 4. 复习到期单词

```bash
# 列出最近到期的 30 个单词（默认）
termhelper list

# 指定数量
termhelper list 50
```

使用 `less` 分页查看，单词按熟练度着色显示。

### 5. 交互测验

```bash
# 从到期单词中随机出题
termhelper test

# 指定题目数量（至少 4 题）
termhelper test 20
```

四选一问答，输入 `A-D` 选择答案，`Q` 退出。每道题会显示正确答案和详细解释。

---

## 命令参考

| 命令 | 功能 |
| ------ | ------ |
| `termhelper add [文本]` | 解析文本中的英文词汇并获取中文解释 |
| `termhelper add -f <文件>` | 从文件读取文本并解析 |
| `termhelper list [数量]` | 列出待复习的到期单词 |
| `termhelper test [数量]` | 进行单词测验 |
| `termhelper status` | 查看单词学习状态统计 |
| `termhelper config` | 编辑配置文件 |
| `termhelper clear` | 清理临时文件 |

---

## 配置说明

配置文件路径：`~/.config/termhelper/config.toml`

详细配置项说明：

```toml
# LLM 提供商列表（按顺序尝试）
providers = ["ollama", "deepseek", "openai"]

# prompt_mode = "en2en"
prompt_mode = "zh2zh"
[deepseek]
api_key = "sk-xxx"           # DeepSeek API 密钥
model = "deepseek-v4-flash"  # 模型名称
api_base = "https://api.deepseek.com"

[ollama]
url = "http://localhost:11434"  # Ollama 服务地址
model = "qwen3.5:4b"           # 本地模型名称

[test]
batch_size = 30      # 每次加载的到期单词数
daily_limit = 100    # 单次测验最大题数

[test.proficiency]
correct_increment = 0.12   # 答对熟练度增量
wrong_decrement = 0.10     # 答错熟练度减量
initial_value = 0.15       # 新词初始熟练度
max_interval_days = 30     # 最大复习间隔（天）

[word_provider]
disabled_words = ["hello", "yes", "the", "no", "of"]  # 自动跳过的基础词汇
```

> **环境变量覆盖：** `DEEPSEEK_API_KEY`、`OPENAI_API_KEY`、`TERMHELPER_*` 可覆盖配置文件中对应项

---

## 数据存储

| 数据 | 路径 |
| ------ | ------ |
| 单词数据库 | `~/.local/share/termhelper/words.db` |
| 配置文件 | `~/.config/termhelper/config.toml` |
| 日志文件 | `~/.local/state/termhelper/log.txt` |
| 临时文件 | `/tmp/termhelper/` |

---

## 复习算法

基于**间隔重复（Spaced Repetition）** 原理：

- 每个新词初始熟练度为 `0.0`
- 答对增加熟练度（默认 `+0.12`），答错降低（默认 `-0.10`）
- 复习间隔 = `60分钟 × 2^(熟练度 × 10)`
- 熟练度越高，复习间隔越长，最高不超过 `max_interval_days` 天
- 每次 `list` 或 `test` 命令仅显示到期的单词

---

## 项目结构

```
termhelper/
├── main.go                  # 程序入口
├── cmd/                     # Cobra 命令定义
│   ├── root.go              # 根命令及初始化
│   ├── add.go               # 添加文本并提取词汇
│   ├── list.go              # 列出到期单词
│   ├── status.go            # 查看学习状态
│   ├── test.go              # 单词测验
│   ├── config.go            # 编辑配置
│   └── clear.go             # 清理临时文件
├── internal/
│   ├── config/              # 配置加载与管理
│   ├── extractor/           # 文本分块与词汇提取
│   ├── llm/                 # LLM 客户端（Ollama/DeepSeek/OpenAI）
│   ├── models/              # 数据模型与常量
│   ├── output/              # 日志与终端输出
│   ├── storage/             # BoltDB 数据持久化
│   ├── tester/              # 交互测验引擎
│   └── writer/              # 单词输出格式化
├── pkg/                     # 通用工具函数
├── demo/                    # 示例与演示代码
├── plan/                    # 开发计划文档
├── Makefile                 # 跨平台编译
└── go.mod                   # Go 模块定义
```

---

## 许可证

[MIT](LICENSE) License. Copyright © 2026 IridiumNan.

---

## 开发计划

参见 [todolist.md](todolist.md) 和 [plan/](plan/) 目录了解详细开发路线图。
