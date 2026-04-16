# ADP CLI

ADP CLI 是来也科技 [ADP（Agentic Document Processing）](https://adp.laiye.com/) 产品的官方命令行工具，用于文档解析与智能抽取。

[English](#english) | [简体中文](#简体中文)

---

<a id="简体中文"></a>

## 功能概览

- **文档解析** — 将 PDF、图片、Office 文档转换为结构化数据
- **智能抽取** — 基于大模型提取发票、订单、证件等关键字段
- **自定义应用** — 创建和管理个性化抽取应用
- **批量处理** — 支持文件夹递归、URL 列表、并发处理
- **同步/异步** — 支持同步等待和异步任务查询两种模式
- **多语言** — 中英文界面自动切换（`--lang` 或 `ADP_LANG` 环境变量）
- **跨平台** — Windows / Linux / macOS，静态编译无依赖

### 支持的文件格式

`.jpg` `.jpeg` `.png` `.bmp` `.tiff` `.tif` `.pdf` `.doc` `.docx` `.xls` `.xlsx` `.ppt` `.pptx`（单文件最大 50MB）

## 快速开始

## 快速开始

### 安装

**方式一：npm 安装（推荐）**

```bash
npm install -g agentic-doc-parse-and-extract-cli
```

**方式二：Shell 脚本（Linux / macOS）**

```bash
curl -fsSL https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.sh | bash
```

**方式三：PowerShell（Windows）**

```powershell
irm https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.ps1 | iex
```

**方式四：手动下载**

从 [GitHub Releases](https://github.com/laiye-ai/adp-cli/releases) 下载对应平台的预编译二进制文件。

**方式五：源码构建**

```bash
git clone https://github.com/laiye-ai/adp-cli.git
cd adp-cli
go build -o adp .
```

### 配置

获取 API Key：访问 [ADP 中国区](https://adp.laiye.com/) 或 [ADP 全球区](https://adp-global.laiye.com/) 注册账户（新用户每月 100 免费积分）。

```bash
# 设置 API Key 和服务地址
adp config set --api-key <your-api-key>
adp config set --api-base-url https://adp.laiye.com

# 查看配置
adp config get
```

### 基本用法

```bash
# 查看可用应用
adp app-id list

# 解析本地文档
adp parse local ./invoice.pdf --app-id <app-id>

# 抽取关键字段
adp extract local ./invoice.pdf --app-id <app-id>

# 批量处理目录
adp parse local ./documents/ --app-id <app-id> --async

# 从 URL 处理
adp extract url https://example.com/file.pdf --app-id <app-id>

# 查询异步任务
adp parse query <task-id>

# 查看剩余积分
adp credit
```

## 命令参考

| 命令 | 说明 |
|------|------|
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

### 全局参数

| 参数 | 说明 |
|------|------|
| `--json` | 以 JSON 格式输出 |
| `--quiet` | 静默模式，仅输出结果 |
| `--lang <en\|zh>` | 指定界面语言 |

### 常用参数

| 参数 | 说明 |
|------|------|
| `--app-id` | 应用 ID（parse/extract 必填） |
| `--async` | 异步模式 |
| `--export <path>` | 导出结果到文件 |
| `--timeout <seconds>` | 超时时间（默认 900 秒） |
| `--concurrency <n>` | 并发数（免费用户最大 1，付费用户最大 2） |

## 构建

### 本地构建

```bash
go build -o adp .
```

### 跨平台构建

项目提供 Makefile 支持一键交叉编译 6 个平台：

```bash
make build-all VERSION=v1.0.0
```

输出到 `dist/` 目录：

| 平台 | 文件名 |
|------|--------|
| Windows x64 | `adp-win32-x64.exe` |
| Windows arm64 | `adp-win32-arm64.exe` |
| Linux x64 | `adp-linux-x64` |
| Linux arm64 | `adp-linux-arm64` |
| macOS x64 | `adp-darwin-x64` |
| macOS arm64 | `adp-darwin-arm64` |

版本号通过构建时注入：`-ldflags "-X github.com/laiye-ai/adp-cli/cmd.version=v1.0.0"`

所有构建均为静态编译（`CGO_ENABLED=0`），无外部依赖。

## 测试

### E2E 测试

```bash
# 离线测试（无需 API Key）
bash tests/test.sh

# 完整测试（需配置 API 凭据）
ADP_API_KEY=<key> ADP_API_BASE_URL=<url> bash tests/test.sh
```

测试报告输出到 `tests/test_report.txt`。

测试覆盖 40 个用例，包括：
- 版本和帮助信息
- 配置管理
- Schema 输出
- 应用列表与缓存
- 文档解析（本地/URL/目录/导出/并发）
- 文档抽取（本地/URL/目录/导出/并发）
- 自定义应用全生命周期（创建/查询/AI生成/更新/删除）
- 积分查询

## CI/CD

项目配置了 GitHub Actions：

- **CI**（`.github/workflows/ci.yml`）— push/PR 到 main 时触发，运行构建和 E2E 测试
- **Release**（`.github/workflows/release.yml`）— 推送 `v*` tag 时触发，交叉编译、创建 GitHub Release 并自动发布到 npm

## 项目结构

```
adp-cli/
├── main.go                  # 入口
├── cmd/                     # 命令定义（cobra）
│   ├── root.go              # 根命令、全局参数、i18n、版本检查
│   ├── config.go            # config 子命令
│   ├── appid.go             # app-id 子命令
│   ├── parse.go             # parse 子命令
│   ├── extract.go           # extract 子命令
│   ├── customapp.go         # custom-app 子命令
│   ├── credit.go            # credit 子命令
│   ├── schema.go            # schema 子命令
│   └── help.go              # 自定义 help
├── internal/
│   ├── api/client.go        # ADP API 客户端
│   ├── config/config.go     # 配置管理（AES-256-GCM 加密）
│   ├── formatter/formatter.go # 输出格式化
│   ├── i18n/i18n.go         # 国际化
│   ├── errors/errors.go     # 错误分类与退出码
│   ├── file/file_handler.go # 文件处理与校验
│   └── updater/updater.go   # 版本更新检查
├── scripts/
│   ├── postinstall.js       # npm 安装后自动下载二进制
│   ├── adp-init.sh          # Linux/macOS 一键安装脚本
│   └── adp-init.ps1         # Windows 一键安装脚本
├── tests/
│   ├── test.sh              # E2E 测试脚本
│   └── samples/             # 测试样本文件
├── package.json             # npm 包配置
├── Makefile                 # 跨平台构建
├── .github/workflows/       # CI/CD
└── go.mod
```

## 配置存储

- 配置目录：`~/.adp/`
- 配置文件：`~/.adp/config.json`
- API Key 加密存储（AES-256-GCM），密钥文件：`~/.adp/key.enc`
- 应用缓存：`~/.adp/app_cache.json`
- 版本检查缓存：`~/.adp/version_check.json`（每 24 小时更新一次）

## 环境变量

| 变量 | 说明 |
|------|------|
| `ADP_API_KEY` | API Key（优先于配置文件） |
| `ADP_API_BASE_URL` | 服务地址 |
| `ADP_LANG` | 界面语言（`en` / `zh`） |
| `ADP_LOG_LEVEL` | 日志级别（`debug` / `info` / `warn` / `error`） |

## 许可证

本项目采用商业许可协议（[license.md](license.md)）。非商业用途（个人学习、研究、教学、开源社区交流等）可免费使用、复制和分发。商业用途需获得来也科技书面授权。

ADP 服务按使用量计费，新用户每月 100 免费积分。

商业授权联系：global_product@laiye.com

---

<a id="english"></a>

## English

ADP CLI is the official command-line tool (Go edition) for [Laiye ADP (Agentic Document Processing)](https://adp-global.laiye.com/), providing document parsing and intelligent extraction capabilities.

### Quick Start

```bash
# Install (npm)
npm install -g agentic-doc-parse-and-extract-cli

# Install (Linux/macOS)
curl -fsSL https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.sh | bash

# Install (Windows PowerShell)
irm https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.ps1 | iex

# Configure
adp config set --api-key <your-api-key>
adp config set --api-base-url https://adp-global.laiye.com

# Parse a document
adp parse local ./invoice.pdf --app-id <app-id>

# Extract key fields
adp extract local ./invoice.pdf --app-id <app-id>

# Check credits
adp credit
```

### Features

- Document parsing (PDF, images, Office formats) to structured data
- Intelligent extraction of key fields (invoices, orders, certificates)
- Custom extraction applications with AI-powered field recommendation
- Batch processing with directory recursion and concurrent workers
- Sync and async processing modes
- English / Chinese interface (`--lang` flag or `ADP_LANG` env var)
- Cross-platform static binaries (Windows / Linux / macOS, x64 / arm64)
- AES-256-GCM encrypted API key storage
- Auto update notification (checks every 24 hours, non-blocking)

### Supported Formats

`.jpg` `.jpeg` `.png` `.bmp` `.tiff` `.tif` `.pdf` `.doc` `.docx` `.xls` `.xlsx` `.ppt` `.pptx` (max 50MB per file)

### Build

```bash
# Local build
go build -o adp .

# Cross-compile all platforms
make build-all VERSION=v1.0.0
```

### E2E Tests

```bash
# Offline tests (no API key needed)
bash tests/test.sh

# Full tests
ADP_API_KEY=<key> ADP_API_BASE_URL=<url> bash tests/test.sh
```

### License

This project is licensed under a Commercial License Agreement ([license.md](license.md)). Free for non-commercial use (personal learning, research, teaching, open-source community). Commercial use requires written authorization from Laiye Technology.

ADP service is billed by usage (100 free credits/month for new users).

Commercial licensing: global_product@laiye.com
