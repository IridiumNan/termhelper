# TermHelper 开发计划（最终版）

## 一、项目全景

### 核心理念

- **所有词都值得学习**：不设停用词，全部进入词库
- **按熟练度调度**：SM-2 变体，弱词高频、强词低频
- **例句来自真实输入**：LLM 只负责释义，例句从用户文本提取
- **批量导入支持**：可快速构建个人领域词库
- **机制透明**：配置、数据、日志皆可触及，拒绝黑盒魔法

### 技术栈

| 组件 | 选型 | 理由 |
| ------ | ------ | ------ |
| 语言 | Go 1.21+ | 静态编译，无依赖，跨平台 |
| CLI | Cobra | 标准化子命令 |
| 配置 | Viper + TOML | 默认值、环境变量、flag 绑定；TOML 易读 |
| 模板嵌入 | `//go:embed` | 带注释的默认配置，独立文件维护 |
| 存储 | BoltDB | 嵌入式 KV，无外部依赖 |
| 日志 | `log/slog` | 标准库，结构化日志，双输出 |
| LLM | 统一接口 + 工厂模式 | 支持 Ollama/OpenAI/DeepSeek，可 fallback |

### 数据存储（XDG）

```
~/.config/termhelper/config.toml    # 配置文件
~/.local/share/termhelper/words.db  # BoltDB 词库
~/.local/share/termhelper/logs/termhelper.log  # 日志
```

### 目录结构

```
termhelper/
├── cmd/
│   ├── root.go          # 主命令，初始化配置/日志
│   ├── add.go           # 从文本提取新词
│   ├── test.go          # 复习测试
│   ├── list.go          # 列出词库
│   ├── stats.go         # 统计信息
│   ├── import.go        # 批量导入
│   └── config_cmd.go    # 配置管理（可选）
├── internal/
│   ├── config/
│   │   ├── config.go           # 加载、验证
│   │   └── config.toml.tmpl    # 带注释的默认配置（embed）
│   ├── models/
│   │   ├── word.go       # Word, Example
│   │   ├── extract.go    # ExtractedWord
│   │   └── llm.go        # WordRequest, WordResponse, ImportRequest, ImportResponse
│   ├── storage/
│   │   ├── db.go         # 打开/关闭 BoltDB
│   │   ├── word.go       # CRUD + 索引维护
│   │   └── schedule.go   # 查询到期词、更新调度
│   ├── llm/
│   │   ├── client.go         # LLMClient 接口
│   │   ├── factory.go        # NewClient（按 providers 顺序 fallback）
│   │   ├── openai_compatible.go  # 共享底层实现（OpenAI/DeepSeek）
│   │   ├── openai.go         # OpenAIClient（语义包装）
│   │   ├── deepseek.go       # DeepSeekClient（语义包装）
│   │   └── ollama.go         # OllamaClient
│   ├── extractor/
│   │   ├── extract.go    # 分词、去重、过滤
│   │   └── sentence.go   # 定位单词所在句子
│   ├── tester/
│   │   ├── quiz.go       # 出题、选项生成、交互
│   │   └── proficiency.go # 熟练度更新
│   ├── output/
│   │   └── render.go     # 高亮、颜色、格式化输出
│   └── logger/
│       └── logger.go     # 全局 Console/File Logger
├── pkg/
│   └── xdg/              # XDG 路径工具
├── main.go
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## 二、数据流总览

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           完整处理链路                                       │
└─────────────────────────────────────────────────────────────────────────────┘

[0] 输入采集
    │  stdin / -f file / 命令行参数
    ▼
    rawText: string
    │
    ▼
[1] 预处理（extractor）
    │  分词 → 去符号 → 小写 → 去重
    ▼
    allTokens: []string
    │
    ▼
[2] 词库比对（storage）
    │  查询词库 → 分出「已存在」和「新词」
    ▼
    newWords: []string (去重，无停用词过滤)
    │
    ▼
[3] 长度控制
    │  若 rawText 超长 → 按句子拆分为多个片段
    ▼
    tasks: []SegmentTask{Context, Words}
    │
    ▼
[4] LLM 批量解释（llm）
    │  构建 Prompt → 调用 ExplainBatch
    ▼
    responses: []WordResponse{Word, Definition, Examples}
    │
    ▼
[5] 存储（storage）
    │  写入 BoltDB（Proficiency=0.15，NextReviewDue=now）
    ▼
    (持久化)
    │
    ▼
[6] 输出展示（output）
    │  原文本高亮（按熟练度颜色）→ 下方列表（按熟练度升序）
    ▼
   终端
```

## 三、开发阶段（共 8 个阶段）

### 阶段 0：项目脚手架（0.5 天）

**目标**：可编译运行的空 CLI

| 任务 | 产出 |
| ------ | ------ |
| `go mod init` | go.mod |
| 安装 Cobra/Viper/BoltDB | 依赖 |
| `cobra-cli init` | cmd/root.go, main.go |
| 添加子命令骨架（add, test, list, stats, import） | cmd/*.go |
| 编写 Makefile | 编译/运行/测试 |

**验收**：`go build` 成功，`./termhelper --help` 显示帮助

### 阶段 1：XDG 目录 + 配置模块（1 天）

**目标**：首次运行自动创建目录和带注释的配置文件

| 任务 | 产出 |
| ------ | ------ |
| 实现 XDG 路径工具（pkg/xdg） | ~/.config, ~/.local/share 路径解析 |
| 创建 config.toml.tmpl（带详细注释） | internal/config/config.toml.tmpl |
| 实现 config.Load() | 读取/写入/默认值/环境变量 |
| 实现 config.Get() / Validate() | 全局配置访问 |
| 在 root.go 的 PersistentPreRun 中调用 | 启动时加载 |

**验收**：

- 首次运行自动生成 `~/.config/termhelper/config.toml`（带注释）
- 配置文件缺失字段时使用默认值
- 可通过 `TERMHELPER_PROVIDER=openai` 覆盖

### 阶段 2：日志模块（0.5 天）

**目标**：终端 + 文件双输出

| 任务 | 产出 |
| ------ | ------ |
| 实现 logger.Init() | 依赖 config.Logging |
| Console Logger（TextHandler，stderr） | 用户可见信息 |
| File Logger（JSONHandler，追加模式） | 详细调试日志 |
| 颜色支持（可选，后续可完善） | 终端输出彩色 |

**验收**：

- `logger.Console.Info("hello")` 输出到终端
- `logger.File.Debug("detail")` 写入日志文件
- 日志目录自动创建

### 阶段 3：数据模型 + 存储层（2 天）

**目标**：词库 CRUD + 熟练度调度

| 任务 | 产出 |
| ------ | ------ |
| 定义 models.Word, Example | internal/models/word.go |
| 实现 storage.Open() | BoltDB 初始化，创建 buckets |
| 实现 word CRUD（Put/Get/List/Delete） | 词条读写 |
| 实现索引维护（due_index） | 按 NextReviewDue 排序 |
| 实现 schedule.Queries() | 查询到期词，批量取 N 个 |
| 预置基础词（首次启动时插入） | file, command, error, etc. |

**验收**：

- 能写入/读取词条
- 查询到期词按时间排序
- 预置词熟练度 0.6

### 阶段 4：LLM 客户端（2 天）

**目标**：支持 Ollama/OpenAI/DeepSeek，可 fallback

| 任务 | 产出 |
| ------ | ------ |
| 定义 LLMClient 接口 | ExplainBatch, Ping |
| 实现 openAICompatibleClient（共享底层） | HTTP 请求/响应 |
| 实现 OpenAIClient（语义包装） | 调用底层 |
| 实现 DeepSeekClient（语义包装） | 调用底层 |
| 实现 OllamaClient | 调用 /api/generate |
| 实现 factory.NewClient() | 按 providers 顺序 fallback |
| 实现 Ping 健康检查 | 连接测试 |

**验收**：

- `providers = ["ollama", "openai"]` 时，Ollama 不可用自动切 OpenAI
- 每个 Client 的 Ping 能正确反映服务状态
- 环境变量可覆盖 API Key（不写入配置文件）

### 阶段 5：提取器（1.5 天）

**目标**：从文本中提取候选词 + 截取例句

| 任务 | 产出 |
| ------ | ------ |
| 实现分词（按空白+标点） | internal/extractor/extract.go |
| 实现去噪（数字/符号/单字母） | 纯数字过滤 |
| 实现句子截取 | 定位单词所在句子 |
| 实现文本拆分（超长处理） | 按句子边界拆分为片段 |

**验收**：

- `"Error: cannot find symbol"` → `["error", "cannot", "find", "symbol"]`
- 每个候选词能正确截取所在句子
- 超长文本（3000+ 字符）自动拆分

### 阶段 6：`add` 命令（2 天）

**目标**：完整串联：输入 → 提取 → LLM → 存储 → 输出

| 任务 | 产出 |
| ------ | ------ |
| 实现 cmd/add.go 主体流程 | 完整 pipeline |
| 集成 extractor | 提取候选词 |
| 集成 storage（查询已有词） | 过滤出真正新词 |
| 集成 llm（ExplainBatch） | 批量获取释义 |
| 集成 storage（写入） | 保存新词 |
| 集成 output 展示 | 高亮 + 列表 |

**验收**：

- `termhelper add "Error: cannot find symbol"` 能完整跑通
- 终端显示高亮文本 + 新词释义列表
- 词库中正确写入新词

### 阶段 7：`test` 命令（2 天）

**目标**：间隔重复测试

| 任务 | 产出 |
| ------ | ------ |
| 实现 proficiency 更新逻辑 | 答对 +0.12，答错 -0.10 |
| 实现 nextInterval 计算 | 1~30 天 |
| 实现 quiz 出题逻辑 | 正确 + 3 干扰项 |
| 实现交互循环 | 选择 → 反馈 → 下一题 |
| 集成 output 展示（彩色） | 正确绿色，错误红色 |

**验收**：

- `termhelper test` 能取出到期词
- 测试后熟练度正确更新
- 下次复习时间正确计算

### 阶段 8：`import` + 其他命令（2 天）

**目标**：批量导入 + 辅助命令

| 任务 | 产出 |
| ------ | ------ |
| 实现 import 解析器 | 自动检测格式（词 / 词+释义 / 词+释义+例句） |
| 实现 llm.ImportBatch | 转换用户输入为结构化词条 |
| 实现 import 存储 | 写入词库 |
| 实现 list 命令 | 列出词库（按熟练度排序） |
| 实现 stats 命令 | 统计总数/熟练度分布/到期数 |

**验收**：

- `termhelper import -f words.txt` 能批量导入
- `list` 显示所有词，`stats` 显示统计

### 阶段 9：打磨与发布（1 天）

**目标**：稳定可用，方便分发

| 任务 | 产出 |
| ------ | ------ |
| 错误处理完善 | 所有 error 有友好提示 |
| 帮助文档 | README, USAGE |
| 交叉编译 | Linux/macOS/Windows |
| GitHub Release | 二进制 + 安装脚本 |

**验收**：用户在 Linux 上 `curl ... | bash` 即可安装使用

## 四、里程碑

| 里程碑 | 时间 | 交付物 |
| -------- | ------ | -------- |
| M1 | 第 3 天 | 骨架 + 配置 + 日志，能生成配置文件和日志 |
| M2 | 第 7 天 | 存储 + LLM，能读写词库、调用大模型 |
| M3 | 第 10 天 | `add` 命令完整可用 |
| M4 | 第 13 天 | `test` 命令完整可用 |
| M5 | 第 15 天 | `import` + 辅助命令 |
| M6 | 第 16 天 | v0.1.0 发布 |

## 五、设计原则（贯穿全程）

1. **KISS**：能合并的不拆分，能显式的不隐式
2. **机制透明**：配置可读，日志详细，错误直接
3. **零魔法**：拒绝过度封装，拒绝黑盒抽象
4. **用户即维护者**：每个用户都能打开 config.toml 看懂并修改

## 六、建议的下一步

1. **立即开始阶段 0**：初始化项目，建立骨架
2. **阶段 1+2 可以并行**：配置和日志互不依赖
3. **阶段 3 是核心基石**：存储层稳定后再开发上层
