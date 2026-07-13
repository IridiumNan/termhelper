# TermHelper 项目开发计划 (Plan) - 修订版

## 1. 项目概述

**TermHelper** 是一个面向计算机专业学生的 CLI 单词学习工具，专为技术英语（Linux 报错、编程语言文档、man 手册等）设计。它能够：

- 自动从英文报错/文档中提取生词，**仅借助 LLM 给出中文释义（1~2 个词）**。
- **例句直接取自用户输入的原始文本片段**（该词所在的句子或上下文），确保学习内容真实、有上下文。
- 将新词加入复习队列，基于改进的 SM-2 算法进行间隔重复测试。
- 根据熟练度动态调整复习间隔，弱词高频出现，强词低频出现。
- 数据完全本地存储，支持云端大模型（OpenAI）和本地模型（Ollama）。
- 编译为单一二进制文件，遵循 XDG 标准存放配置和数据，便于分发。

---

## 2. 技术选型（不变）

| 组件          | 选择                  | 理由                                                                 |
|---------------|-----------------------|----------------------------------------------------------------------|
| 编程语言      | Go 1.21+              | 编译为静态二进制，无依赖，跨平台，适合 CLI 工具。                    |
| CLI 框架      | Cobra + Viper         | 标准化子命令（add, test, list, stats, config），配置管理方便。       |
| 数据存储      | BoltDB（嵌入式 KV）   | 纯 Go 实现，无外部依赖，文件级数据库，适合本地轻量存储。             |
| LLM 客户端    | 自定义接口（OpenAI / Ollama） | 支持云端和本地模型，用户可配置。                                     |
| 构建工具      | Go 官方工具 + Makefile | 交叉编译多平台（Linux/macOS/Windows）。                              |
| 配置/数据路径 | XDG Base Directory     | `~/.config/termhelper/`（配置），`~/.local/share/termhelper/`（数据） |

---

## 3. 核心模块划分（略作调整）

```
termhelper/
├── cmd/                     # Cobra 命令
│   ├── root.go
│   ├── add.go
│   ├── test.go
│   ├── list.go
│   ├── stats.go
│   └── config.go
├── internal/
│   ├── storage/            # 数据库操作（BoltDB）
│   │   ├── db.go           # 初始化、bucket 创建
│   │   ├── word.go         # 单词 CRUD + 索引更新
│   │   └── schedule.go     # 查询到期单词、更新 next_review_due
│   ├── llm/                # LLM 客户端抽象
│   │   ├── client.go       # 接口定义
│   │   ├── openai.go
│   │   └── ollama.go
│   ├── extractor/          # 生词提取 + 例句提取（新增）
│   │   ├── extract.go      # 分词、过滤、词形还原
│   │   └── sentence.go     # 从原文中截取单词所在的句子作为例句
│   ├── tester/             # 测试逻辑
│   │   ├── quiz.go         # 出题、选项生成、交互
│   │   └── proficiency.go  # 熟练度更新算法
│   └── config/             # 配置管理（Viper）
│       └── config.go
├── pkg/                    # 可复用工具
│   └── xdg/                # XDG 路径工具
├── go.mod
├── main.go
└── Makefile
```

---

## 4. 数据存储设计（BoltDB）【变更：例句来源】

BoltDB 为单文件 `words.db`，包含以下 Bucket：

### 4.1 Bucket: `words`

存储每个单词的完整信息。**Key** = 单词小写（string），**Value** = JSON 对象。

```json
{
  "word": "compile",
  "definition": "编译",                    // LLM 生成的中文释义（1~2个词）
  "examples": [                           // 例句直接来自用户输入的原始文本
    {"en": "Please compile the source code before running.", "context": "…整个原始段落或文件名…"},
    {"en": "The compiler threw an error."}   // 可保留其他出现的句子
  ],
  "proficiency": 0.15,
  "last_reviewed": 0,
  "next_review_due": 0,
  "created_at": 1704067200
}
```

- **`examples` 数组**：存储该单词在用户历史输入中出现的原始句子（去重）。每个条目包含 `en`（英文原句）和可选的 `context`（整个段落或文件名，便于溯源）。**不需要中文翻译**，因为学习目标是理解原文。
- 当同一个单词再次出现在新的输入中，我们可以追加新的例句（去重后），丰富学习材料。

### 4.2 Bucket: `due_index` (二级索引)

不变，用于按 `next_review_due` 排序查询到期单词。  
**Key** = `fmt.Sprintf("%020d:%s", next_review_due, word)`  
**Value** = 空（仅用于占位）

### 4.3 Bucket: `metadata`

存储全局元数据（如版本、统计信息等），暂不扩展。

---

## 5. 核心处理链路【变更：例句提取】

1. **用户输入**原始英文文本（报错、手册片段等）。
2. **分词与过滤**：拆分成 token，去除数字、符号，进行简单的词形还原（小写、去复数等）。
3. **排除已存在词库**：只保留未记录的生词。
4. **例句提取**：对每个生词，在其原始输入文本中定位该词所在句子（按 `.` `!` `?` 等标点分割，若未找到则取前后若干词），作为该词的例句存入 `examples` 列表（若该句与已有例句重复则跳过）。
5. **调用 LLM**：仅针对该生词（带上上下文句子）请求**中文释义（1~2 个中文词）**，不要求生成例句。
6. **存储新词**：写入 `words`，初始 proficiency=0.15，`next_review_due=now`，并插入 `due_index` 索引。
7. 输出处理结果。

---

## 6. SM-2 调度算法实现（不变）

### 6.1 熟练度（Proficiency）

- 范围：`0.0 ~ 1.0`
- 初始值：`0.15`（新词）
- 更新规则：
  - 答对：`proficiency = min(1.0, proficiency + 0.12)`
  - 答错：`proficiency = max(0.0, proficiency - 0.10)`

### 6.2 复习间隔计算

```go
func nextInterval(proficiency float64) int64 {
    days := 1 + proficiency * 30
    return int64(math.Floor(days))
}
```

实际 `next_review_due = last_reviewed + interval * 86400`。

### 6.3 查询到期单词（不变）

`test` 命令执行时，从 `due_index` 中获取 `next_review_due <= now` 的单词，按顺序取前 `batch_size` 个。

### 6.4 测试流程【变更：例句展示】

1. 对每个到期单词，从 `examples` 列表中**随机选取一个例句**（原文句子）显示在屏幕上，目标单词用黄色高亮。
2. 生成 4 个中文释义选项：正确选项 + 3 个干扰项（从词库中随机抽取不同单词的 `definition`）。
3. 用户选择后：
   - 正确：绿色输出“正确”，更新 proficiency += 0.12
   - 错误：红色输出正确释义，更新 proficiency -= 0.10
4. 更新 `last_reviewed = now`，重新计算 `next_review_due`。
5. 更新数据库（单词记录 + 删除旧索引 Key + 插入新索引 Key）。
6. 按回车继续下一题。

### 6.5 干扰项生成（不变）

从词库中随机选取 3 个**不同单词**的 definition，确保不与正确答案重复。若词库不足，则用占位符。

---

## 7. CLI 命令详细设计【变更：add 流程】

### 7.1 全局 Flags（不变）

### 7.2 `add` 子命令

```bash
termhelper add [text]          # 从参数读取文本
termhelper add -f file.txt     # 从文件读取
echo "error message" | termhelper add   # 从 stdin 读取
```

**处理流程**：

1. 读取输入文本。
2. 提取所有候选 token，过滤已有词库中的词。
3. 对每个新词，**从输入文本中提取该词所在的句子**（作为例句）。
4. 调用 LLM（传递该词和上下文句子）获取中文释义。
5. 存储新词，包含释义和提取的例句（至少一个）。
6. 打印添加成功的单词数量及释义。

**注意**：若同一个单词在输入中出现多次，可提取多个不同句子作为例句，但存储时去重。

### 7.3 `test` 子命令（不变）

### 7.4 `list` 子命令（不变）

### 7.5 `stats` 子命令（不变）

### 7.6 `config` 子命令（不变）

---

## 8. LLM 集成【变更：Prompt 及返回格式】

### 8.1 接口定义（不变）

```go
type LLMClient interface {
    Explain(word, context string) (*WordEntry, error)  // 只返回释义，不含例句
}
```

### 8.2 Prompt 模板（修改）

```
You are a technical English tutor. Given the following context from a software development error message or manual:

Context sentence: "<sentence>"
Word to explain: "<word>"

Provide a short Chinese definition (1-2 Chinese words) for this word in this technical context.
Only return the definition, nothing else. Do not provide examples or extra text.
```

**返回格式**：仅字符串，如 `"编译"`。

### 8.3 OpenAI / Ollama 实现（不变，但解析简化为纯文本）

### 8.4 错误处理与重试（不变）

---

## 9. 配置管理（不变，但需调整配置项？不涉及）

---

## 10. 实现步骤与里程碑【调整优先级】

### 阶段 1：基础框架（第 1-2 天）- 不变

### 阶段 2：存储层（第 3-4 天）- 不变，但需适配新字段

### 阶段 3：LLM 客户端（第 5-6 天）- 调整为只返回释义字符串

### 阶段 4：生词提取 + 例句提取（第 7-8 天）【新增关键任务】

- 实现 `extractor/sentence.go`：按标点分割句子，并定位目标单词所在的句子（需处理大小写、标点包围的情况）。
- 实现例句去重（基于句子内容哈希）。

### 阶段 5：`add` 命令（第 9-10 天）- 集成上述提取和 LLM 调用

### 阶段 6：`test` 命令（第 11-14 天）- 展示例句（原文），而不是模型生成例句

### 阶段 7：其他命令（第 15 天）

### 阶段 8：打磨与测试（第 16-18 天）

### 阶段 9：构建与发布（第 19 天）

---

## 11. 扩展性考虑（新增）

- 后续可支持手动添加自定义例句（通过 `add --example "..."` 覆盖提取）。
- 可支持导入 Anki 卡片时保留原句字段。

---

## 12. 风险与缓解（新增风险）

| 风险 | 缓解措施 |
| ------ | ---------- |
| 例句提取可能截取不完整（如跨行） | 预处理时合并换行为空格，按 `.` `?` `!` 分割；若未找到句子，则截取前后 50 个字符作为上下文。 |
| 同一单词多次出现，例句重复 | 用句子内容的哈希去重，保留最多 5 个不同例句。 |
| LLM 仅返回释义可能不够精确 | 可允许用户手动编辑释义（通过 `edit` 命令后续实现）。 |

---

## 13. 测试策略（调整）

- 单元测试：新增 `sentence.go` 的例句提取逻辑。
- 集成测试：用真实报错文本测试 `add`，检查例句是否来自原文。

---

## 14. 文档（更新说明）

- 在 README 中强调“例句来自您的真实输入，不是 AI 生成”，突出产品优势。

---

## 15. 交付物（不变）

- 源代码(github仓库)
- 多平台二进制
- 一键安装和配置脚本

---
