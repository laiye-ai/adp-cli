# ADP CLI 使用指南

ADP CLI 是 [Laiye ADP（Agentic Document Processing）](https://github.com/laiye-ai/adp-cli) 的官方命令行工具，用于 AI 驱动的文档解析与智能字段提取。支持 PDF、图片、Office 文档等格式，可将非结构化文档转为结构化数据。

---

## 目录

- [安装](#安装)
- [配置](#配置)
- [文档解析（parse）](#文档解析parse)
- [字段提取（extract）](#字段提取extract)
- [异步工作流](#异步工作流)
- [批量处理与并发](#批量处理与并发)
- [应用管理（app-id）](#应用管理app-id)
- [自定义应用（custom-app）](#自定义应用custom-app)
- [额度查询（credit）](#额度查询credit)
- [命令 Schema（schema）](#命令-schemaschema)
- [全局选项](#全局选项)
- [环境变量](#环境变量)
- [退出码](#退出码)
- [配置存储](#配置存储)
- [支持的文件格式](#支持的文件格式)
- [获取帮助](#获取帮助)

---

## 安装

### npm 安装（推荐）

```bash
npm install -g @laiye-adp/agentic-doc-parse-and-extract-cli
```

### Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.sh | bash
```

### Windows（PowerShell）

```powershell
irm https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.ps1 | iex
```

### GitHub Release 下载

从 [Releases](https://github.com/laiye-ai/adp-cli/releases) 页面下载对应平台的二进制文件。

支持平台：Windows / Linux / macOS（x64 和 arm64）。

安装完成后验证：

```bash
adp version
```

---

## 配置

使用前需配置 API Key 和服务地址。配置信息以 AES-256-GCM 加密存储在 `~/.adp/config.json` 中。

### 设置配置

```bash
# 设置 API Key
adp config set --api-key YOUR_API_KEY

# 设置 API 地址（可选，如使用私有部署）
adp config set --api-base-url https://your-server.com
```

### 查看配置

```bash
adp config get
```

API Key 将以掩码形式显示。

### 清除配置

```bash
adp config clear        # 需确认
adp config clear -y     # 跳过确认
```

---

## 文档解析（parse）

将文档转换为结构化数据（Markdown、表格等）。

### 解析本地文件

```bash
# 解析单个文件
adp parse local ./invoice.pdf --app-id APP_ID

# 解析整个目录（递归扫描所有支持的文件）
adp parse local ./documents/ --app-id APP_ID

# 导出结果到文件
adp parse local ./invoice.pdf --app-id APP_ID --export ./result.json
```

### 解析 URL

```bash
# 解析单个 URL
adp parse url https://example.com/doc.pdf --app-id APP_ID

# 从文件中批量读取 URL（每行一个 URL）
adp parse url ./url_list.txt --app-id APP_ID
```

### 解析 Base64 数据

```bash
adp parse base64 BASE64_DATA --app-id APP_ID --file-name invoice.pdf
```

### 异步解析 + 查询任务

```bash
# 提交异步任务
adp parse local ./doc.pdf --app-id APP_ID --async

# 查询任务结果
adp parse query TASK_ID

# 等待任务完成
adp parse query TASK_ID --watch
```

### parse local / url / base64 参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--app-id` | string | — | **必填**，应用 ID |
| `--async` | bool | false | 异步模式 |
| `--no-wait` | bool | false | 仅提交任务，不等待结果（与 `--async` 配合使用） |
| `--export` | string | — | 导出路径（文件或目录） |
| `--timeout` | int | 900 | 超时时间（秒） |
| `--concurrency` | int | 1 | 并发数（免费用户最大 1，付费用户最大 2） |
| `--retry` | int | 0 | 失败重试次数（指数退避） |
| `--file-name` | string | document | 文件名（仅 base64 子命令可用） |

### parse query 参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--watch` | bool | false | 轮询等待任务完成 |
| `--file` | string | — | 从 JSON 文件读取任务 ID（`--no-wait` 的输出文件） |
| `--export` | string | — | 导出结果到文件或目录 |
| `--timeout` | int | 900 | 超时时间（秒） |
| `--concurrency` | int | 1 | 并发查询数 |

---

## 字段提取（extract）

基于 LLM 从文档中提取指定字段（如发票号、金额、日期等）。子命令和参数与 `parse` 完全一致。

```bash
# 提取本地文件字段
adp extract local ./invoice.pdf --app-id APP_ID

# 提取 URL 文件字段
adp extract url https://example.com/doc.pdf --app-id APP_ID

# 批量提取目录中的文件
adp extract local ./documents/ --app-id APP_ID --concurrency 2 --export ./results/
```

### extract local / url / base64 参数

与 `parse local / url / base64` 参数完全一致，参见上方表格。

### extract query 参数

与 `parse query` 参数完全一致，参见上方表格。

---

## 异步工作流

处理大文件或批量任务时，使用 `--async` 提交任务，CLI 返回 `task-id`，再用 `parse query` / `extract query` 轮询结果：

```bash
adp parse local ./big.pdf --app-id APP_ID --async
# 返回一个 task-id

adp parse query TASK_ID
```

### 两阶段异步（`--no-wait`）

默认情况下，`--async` 会提交并轮询直到完成——适合 AI Agent 使用。对于可恢复的工作流，使用两阶段模式：

**第一阶段：提交任务**

```bash
adp extract local ./documents/ --app-id APP_ID --async --no-wait --export tasks.json
```

输出为包含任务 ID 的 JSON 数组：

```json
[
  {"path": "invoice.pdf", "task_id": "task_abc123"},
  {"path": "contract.pdf", "task_id": "task_def456"}
]
```

**第二阶段：查询结果**

```bash
adp extract query --watch --file tasks.json
adp extract query --watch --file tasks.json --export ./results/
```

即使 CLI 中途崩溃，`tasks.json` 中的任务 ID 也会被保留——随时可用 `query --file` 恢复查询。

---

## 批量处理与并发

处理目录或多个 URL 时，CLI 自动进入批量模式：

```bash
# 并发处理目录中所有文件，失败自动重试 3 次
adp parse local ./documents/ --app-id APP_ID --concurrency 2 --retry 3 --export ./results/
```

### 输出结构

批量处理时，CLI 将每个结果写入单独的文件：

```
adp_results_20250417_153020/
├── _summary.json              # 汇总（总数、成功、失败、每文件状态）
├── invoice_01.pdf.json        # 成功结果
├── contract_02.docx.json
└── report_03.pdf.error.json   # 错误详情
```

- `--export <dir>` — 指定输出目录
- 不加 `--export` — 自动创建 `adp_results_<timestamp>/`
- 单文件 — 输出到 stdout 或 `--export` 指定的文件路径

### 批量处理特性

- 结果按文件输出为独立 JSON 文件，附带 `_summary.json` 汇总
- TTY 终端显示彩色进度条，非 TTY 输出 JSON 行格式进度
- 部分失败时退出码为 6，全部失败退出码为 1

---

## 应用管理（app-id）

### 列出可用应用

```bash
# 列出所有应用
adp app-id list

# 按标签过滤
adp app-id list --app-label "invoice"

# 仅列出自定义应用
adp app-id list --app-type 1

# 仅列出系统预设应用（开箱即用）
adp app-id list --app-type 0

# 限制返回数量
adp app-id list --limit 50
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--app-label` | string | — | 按标签过滤 |
| `--app-type` | int | 不传 | 应用类型（不传=全部，0=系统预设，1=自定义） |
| `--limit` | int | 120 | 返回数量限制 |

### 查看本地缓存

```bash
adp app-id cache
```

---

## 自定义应用（custom-app）

> 以下所有 `custom-app` 子命令均支持 `--api-key` 参数，用于指定 API Key（覆盖配置文件中的值）。

### 创建自定义应用

```bash
adp custom-app create \
  --app-name "发票提取" \
  --parse-mode standard \
  --extract-fields '[{"name":"invoice_no","type":"string","description":"发票号码"}]'
```

`--extract-fields` 支持 JSON 字符串或 JSON 文件路径：

```bash
adp custom-app create \
  --app-name "发票提取" \
  --parse-mode agentic \
  --extract-fields ./fields.json
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--app-name` | string | — | **必填**，应用名称 |
| `--extract-fields` | string | — | **必填**，提取字段定义（JSON 字符串或文件路径） |
| `--parse-mode` | string | — | **必填**，解析模式（`advance` / `standard` / `agentic`） |
| `--app-label` | string | — | 应用标签 |
| `--enable-long-doc` | string | — | 是否启用长文档模式 |
| `--long-doc-config` | string | — | 长文档配置 |
| `--api-key` | string | — | 指定 API Key |

### 更新自定义应用

```bash
adp custom-app update \
  --app-id APP_ID \
  --extract-fields ./updated_fields.json \
  --parse-mode agentic \
  --enable-long-doc true
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--app-id` | string | — | **必填**，应用 ID |
| `--extract-fields` | string | — | **必填**，提取字段定义 |
| `--parse-mode` | string | — | **必填**，解析模式 |
| `--enable-long-doc` | string | — | **必填**，是否启用长文档模式 |
| `--app-name` | string | — | 应用名称 |
| `--app-label` | string | — | 应用标签 |
| `--long-doc-config` | string | — | 长文档配置 |
| `--api-key` | string | — | 指定 API Key |

### 查看应用配置

```bash
adp custom-app get-config --app-id APP_ID

# 查看指定版本
adp custom-app get-config --app-id APP_ID --config-version 2
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--app-id` | string | — | **必填**，应用 ID |
| `--config-version` | string | — | 指定配置版本 |
| `--api-key` | string | — | 指定 API Key |

### 删除应用 / 版本

```bash
adp custom-app delete --app-id APP_ID
adp custom-app delete-version --app-id APP_ID --config-version 2
```

| 命令 | 必填参数 |
|------|----------|
| `delete` | `--app-id` |
| `delete-version` | `--app-id`、`--config-version` |

### AI 推荐提取字段

上传样本文档，让 AI 自动推荐提取字段：

```bash
# 从本地文件
adp custom-app ai-generate --app-id APP_ID --file-local ./sample.pdf

# 从 URL
adp custom-app ai-generate --app-id APP_ID --file-url https://example.com/sample.pdf

# 从 Base64
adp custom-app ai-generate --app-id APP_ID --base64 BASE64_DATA
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--app-id` | string | — | **必填**，应用 ID |
| `--file-url` | string | — | 样本文件 URL（三选一） |
| `--file-local` | string | — | 本地样本文件路径（三选一） |
| `--base64` | string | — | Base64 编码的样本数据（三选一） |
| `--api-key` | string | — | 指定 API Key |

---

## 额度查询（credit）

```bash
adp credit

# 使用指定 API Key 查询
adp credit --api-key YOUR_API_KEY
```

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--api-key` | string | — | 指定 API Key（覆盖配置文件） |

---

## 命令 Schema（schema）

输出完整的命令结构（JSON 格式），适合 AI Agent 集成：

```bash
adp schema
```

AI Agent 应调用 `adp schema` 获取机器可读的权威命令规格，而非依赖文档。

---

## 全局选项

所有命令均支持以下选项：

| 选项 | 说明 |
|------|------|
| `--json` | 以 JSON 格式输出 |
| `--quiet` | 静默模式，仅输出错误 |
| `--lang zh` | 设置中文界面 |
| `--lang en` | 设置英文界面 |

示例：

```bash
# JSON 输出，方便管道处理
adp parse local ./doc.pdf --app-id APP_ID --json | jq '.markdown'

# 静默模式，仅关注错误
adp parse local ./docs/ --app-id APP_ID --quiet --export ./out/
```

---

## 环境变量

| 环境变量 | 说明 | 优先级 |
|----------|------|--------|
| `ADP_API_KEY` | API Key（覆盖配置文件） | 高于配置文件 |
| `ADP_API_BASE_URL` | API 地址 | 高于配置文件 |
| `ADP_LANG` | 语言（`en` / `zh`） | 高于系统语言 |
| `ADP_LOG_LEVEL` | 日志级别（`debug` / `info` / `warn` / `error`） | — |

---

## 退出码

| 退出码 | 含义 |
|--------|------|
| 0 | 成功 |
| 1 | 一般错误 / 网络错误 / API 错误 |
| 2 | 参数错误 |
| 3 | 资源未找到 |
| 4 | 权限不足 |
| 5 | 冲突错误 |
| 6 | 批量处理部分失败 |

---

## 配置存储

| 路径 | 说明 |
|------|------|
| `~/.adp/` | 配置目录 |
| `~/.adp/config.json` | 配置文件 |
| `~/.adp/key.enc` | 加密的 API Key（AES-256-GCM） |
| `~/.adp/app_cache.json` | 应用列表缓存 |
| `~/.adp/version_check.json` | 版本检查缓存（每 24 小时刷新） |

---

## 支持的文件格式

PDF、JPG、JPEG、PNG、BMP、TIFF、TIF、DOC、DOCX、XLS、XLSX、PPT、PPTX

单文件大小限制：**50 MB**。

---

## 获取帮助

```bash
# 查看总帮助
adp --help

# 查看子命令帮助
adp parse --help
adp parse local --help
adp custom-app create --help

# 查看命令结构（JSON，适合 AI Agent 集成）
adp schema
```
