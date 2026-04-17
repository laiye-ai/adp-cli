# ADP CLI 使用指南

ADP CLI 是 [Laiye ADP（Agentic Document Processing）](https://github.com/laiye-ai/adp-cli) 的官方命令行工具，用于 AI 驱动的文档解析与智能字段提取。支持 PDF、图片、Office 文档等格式，可将非结构化文档转为结构化数据。

---

## 目录

- [安装](#安装)
- [配置](#配置)
- [文档解析（parse）](#文档解析parse)
- [字段提取（extract）](#字段提取extract)
- [应用管理（app-id）](#应用管理app-id)
- [自定义应用（custom-app）](#自定义应用custom-app)
- [额度查询（credit）](#额度查询credit)
- [全局选项](#全局选项)
- [批量处理与并发](#批量处理与并发)
- [环境变量](#环境变量)
- [退出码](#退出码)

---

## 安装

### npm 安装

```bash
npm install -g agentic-doc-parse-and-extract-cli
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

### 常用参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--app-id` | string | — | **必填**，应用 ID |
| `--async` | bool | false | 异步模式 |
| `--export` | string | — | 导出路径（文件或目录） |
| `--timeout` | int | 900 | 超时时间（秒） |
| `--concurrency` | int | 1 | 并发数（免费用户最大 1，付费用户最大 2） |
| `--retry` | int | 0 | 失败重试次数（指数退避） |

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
```

### 查看本地缓存

```bash
adp app-id cache
```

---

## 自定义应用（custom-app）

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
  --parse-mode accurate \
  --extract-fields ./fields.json
```

### 更新自定义应用

```bash
adp custom-app update \
  --app-id APP_ID \
  --extract-fields ./updated_fields.json \
  --parse-mode accurate \
  --enable-long-doc true
```

### 查看应用配置

```bash
adp custom-app get-config --app-id APP_ID

# 查看指定版本
adp custom-app get-config --app-id APP_ID --config-version 2
```

### 删除应用 / 版本

```bash
adp custom-app delete --app-id APP_ID
adp custom-app delete-version --app-id APP_ID --config-version 2
```

### AI 推荐提取字段

上传样本文档，让 AI 自动推荐提取字段：

```bash
# 从本地文件
adp custom-app ai-generate --app-id APP_ID --file-local ./sample.pdf

# 从 URL
adp custom-app ai-generate --app-id APP_ID --file-url https://example.com/sample.pdf
```

---

## 额度查询（credit）

```bash
adp credit
```

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

## 批量处理与并发

处理目录或多个 URL 时，CLI 自动进入批量模式：

```bash
# 并发处理目录中所有文件，失败自动重试 3 次
adp parse local ./documents/ --app-id APP_ID --concurrency 2 --retry 3 --export ./results/
```

批量处理特性：
- 结果按文件输出为独立 JSON 文件，附带 `_summary.json` 汇总
- TTY 终端显示彩色进度条，非 TTY 输出 JSON 行格式进度
- 部分失败时退出码为 6，全部失败退出码为 1

---

## 环境变量

| 环境变量 | 说明 | 优先级 |
|----------|------|--------|
| `ADP_API_KEY` | API Key（覆盖配置文件） | 高于配置文件 |
| `ADP_API_BASE_URL` | API 地址 | 高于配置文件 |
| `ADP_LANG` | 语言（`en` / `zh`） | 高于系统语言 |
| `ADP_LOG_LEVEL` | 日志级别 | — |

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
