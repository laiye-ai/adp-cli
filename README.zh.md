# ADP CLI

来也科技 [ADP（Agentic Document Processing）](https://adp.laiye.com/) 产品的官方命令行工具 —— 文档解析与智能字段抽取。

[English](README.md) | [简体中文](README.zh.md)

## 功能

- **文档解析** — 将文档转换为结构化文本（Markdown 或 JSON）
- **发票/订单抽取** — 从发票、收据、采购订单等文档中抽取关键字段和表格
- **自定义文档抽取** — 从任意类型的文档中抽取自定义字段或表格
- **批量处理** — 并发处理文件夹或 URL 列表中的多个文档，每个文件单独输出结果
- **同步/异步** — 支持同步和异步两种模式
- **两阶段异步** — `--async --no-wait` 仅提交任务并输出 task-id 列表；`query --file` 从中断处恢复查询
- **可靠性** — 自动重试与指数退避（`--retry`），细粒度退出码
- **跨平台** — Windows / Linux / macOS，静态二进制无依赖

## 支持的文件格式

`.jpg` `.jpeg` `.png` `.bmp` `.tiff` `.tif` `.pdf` `.doc` `.docx` `.xls` `.xlsx` `.ppt` `.pptx`（单文件最大 50MB）

## 安装

```bash
# npm（推荐）
npm install -g @laiye-adp/agentic-doc-parse-and-extract-cli

# Linux / macOS
curl -fsSL https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.sh | bash

# Windows（PowerShell）
irm https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.ps1 | iex
```

或从 [GitHub Releases](https://github.com/laiye-ai/adp-cli/releases) 下载预编译二进制。

## 配置

访问 [https://adp.laiye.com/](https://adp.laiye.com/) 注册并获取 API Key（新用户每月 100 免费积分）。

```bash
adp config set --api-key <your-api-key>
adp config set --api-base-url https://adp.laiye.com
adp config get
```

## 快速示例

```bash
# 查看可用应用
adp app-id list

# 解析本地文档
adp parse local ./invoice.pdf --app-id <app-id>

# 抽取关键字段
adp extract local ./invoice.pdf --app-id <app-id>

# 异步解析目录
adp parse local ./documents/ --app-id <app-id> --async

# 处理远程 URL
adp extract url https://example.com/file.pdf --app-id <app-id>

# 查询异步任务
adp parse query <task-id>

# 两阶段异步（分开提交和查询，支持断点续传）
adp extract local ./documents/ --app-id <app-id> --async --no-wait --export tasks.json
adp extract query --watch --file tasks.json

# 失败自动重试（最多 2 次）
adp parse local ./documents/ --app-id <app-id> --retry 2

# 查看剩余积分
adp credit
```

## 命令

> AI Agent 应调用 `adp schema` 获取机器可读的权威命令规格。下表仅供人类速查。

| 命令 | 说明 |
|---|---|
| `adp version` | 显示版本号 |
| `adp config set` | 设置 API Key / 服务地址 |
| `adp config get` | 查看当前配置 |
| `adp config clear` | 清除配置 |
| `adp app-id list` | 列出可用应用 |
| `adp app-id cache` | 从本地缓存读取应用列表 |
| `adp parse local <path>` | 解析本地文件/目录 |
| `adp parse url <url>` | 解析远程文件（支持 URL 列表文件） |
| `adp parse base64 <data>` | 解析 Base64 编码内容 |
| `adp parse query <task-id...>` | 查询异步解析任务（支持多个 ID 或 `--file`） |
| `adp extract local <path>` | 抽取本地文件/目录 |
| `adp extract url <url>` | 抽取远程文件 |
| `adp extract base64 <data>` | 抽取 Base64 编码内容 |
| `adp extract query <task-id...>` | 查询异步抽取任务（支持多个 ID 或 `--file`） |
| `adp custom-app create` | 创建自定义抽取应用 |
| `adp custom-app update` | 更新自定义应用配置 |
| `adp custom-app get-config` | 查看应用配置 |
| `adp custom-app delete` | 删除自定义应用 |
| `adp custom-app delete-version` | 删除指定配置版本 |
| `adp custom-app ai-generate` | AI 推荐抽取字段 |
| `adp credit` | 查看剩余积分 |
| `adp schema` | 输出命令 Schema（供 AI Agent 使用） |

## 参数

| 参数 | 说明 |
|---|---|
| `--json` | 以 JSON 格式输出 |
| `--quiet` | 静默模式，仅输出结果 |
| `--lang <en\|zh>` | 指定界面语言 |
| `--app-id` | 应用 ID（parse / extract 必填） |
| `--async` | 异步模式 |
| `--no-wait` | 仅提交任务，不等待结果（与 `--async` 配合使用） |
| `--export <path>` | 导出结果到文件（单文件）或目录（批量） |
| `--timeout <seconds>` | 超时时间（默认 900 秒） |
| `--concurrency <n>` | 并发数（免费用户最大 1，付费用户最大 2） |
| `--retry <n>` | 可重试错误的重试次数（默认 0） |
| `--file <path>` | 从 JSON 文件读取任务 ID（`--no-wait` 的输出文件，仅 query 可用） |

## 异步工作流

处理大文件或批量任务时，使用 `--async` 提交任务，CLI 返回 `task-id`，再用 `parse query` / `extract query` 轮询结果：

```bash
adp parse local ./big.pdf --app-id <app-id> --async
# 返回一个 task-id

adp parse query <task-id>
```

### 两阶段异步（`--no-wait`）

默认情况下，`--async` 会提交并轮询直到完成——适合 AI Agent 使用。对于可恢复的工作流，使用两阶段模式：

**第一阶段：提交任务**

```bash
adp extract local ./documents/ --app-id <app-id> --async --no-wait --export tasks.json
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

## 批量处理

处理多个文件/URL 时，CLI 会将每个结果写入单独的文件：

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

## 退出码

| 退出码 | 含义 |
|------|---------|
| `0` | 全部成功 |
| `1` | 全部失败 / 系统错误 |
| `2` | 参数错误 |
| `3` | 资源未找到 |
| `4` | 权限不足 |
| `5` | 冲突 |
| `6` | 部分失败（批量中部分任务失败） |

## 环境变量

| 变量 | 说明 |
|---|---|
| `ADP_API_KEY` | API Key（优先于配置文件） |
| `ADP_API_BASE_URL` | 服务地址 |
| `ADP_LANG` | 界面语言（`en` / `zh`） |
| `ADP_LOG_LEVEL` | 日志级别（`debug` / `info` / `warn` / `error`） |

## 配置存储

- 配置目录：`~/.adp/`
- 配置文件：`~/.adp/config.json`
- 加密的 API Key：`~/.adp/key.enc`（AES-256-GCM）
- 应用缓存：`~/.adp/app_cache.json`
- 版本检查缓存：`~/.adp/version_check.json`（每 24 小时刷新）

## 许可证

商业许可协议 —— 详见 [license.md](license.md)。非商业用途（个人学习、研究、教学、开源社区交流等）可免费使用。商业用途需获得来也科技书面授权。联系：global_product@laiye.com

## 参与贡献

跨平台构建：`make build-all VERSION=v1.0.0`。运行 E2E 测试：`bash tests/test.sh`。
