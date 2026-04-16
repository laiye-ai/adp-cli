# ADP CLI

来也科技 [ADP（Agentic Document Processing）](https://adp.laiye.com/) 产品的官方命令行工具 —— 文档解析与智能字段抽取。

[English](README.md) | [简体中文](README.zh.md)

## 功能

- **文档解析** — 将文档转换为结构化文本（Markdown 或 JSON）
- **发票/订单抽取** — 从发票、订单等文档中抽取关键字段或表格
- **自定义文档抽取** — 从任意类型的文档中抽取自定义字段或表格
- **批量处理** — 支持批量处理文件夹或 URL 列表中的多个文档
- **同步/异步** — 支持同步和异步两种模式
- **跨平台** — 支持 Windows / Linux / macOS

## 支持的文件格式

`.jpg` `.jpeg` `.png` `.bmp` `.tiff` `.tif` `.pdf` `.doc` `.docx` `.xls` `.xlsx` `.ppt` `.pptx`（单文件最大 50MB）

## 安装

```bash
# npm（推荐）
npm install -g agentic-doc-parse-and-extract-cli

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
| `adp parse query <task-id>` | 查询异步解析任务 |
| `adp extract local <path>` | 抽取本地文件/目录 |
| `adp extract url <url>` | 抽取远程文件 |
| `adp extract base64 <data>` | 抽取 Base64 编码内容 |
| `adp extract query <task-id>` | 查询异步抽取任务 |
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
| `--export <path>` | 导出结果到文件 |
| `--timeout <seconds>` | 超时时间（默认 900 秒） |
| `--concurrency <n>` | 并发数（免费用户最大 1，付费用户最大 2） |

## 异步工作流

处理大文件或批量任务时，使用 `--async` 提交任务，CLI 返回 `task-id`，再用 `parse query` / `extract query` 轮询结果：

```bash
adp parse local ./big.pdf --app-id <app-id> --async
# 返回一个 task-id

adp parse query <task-id>
```

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
