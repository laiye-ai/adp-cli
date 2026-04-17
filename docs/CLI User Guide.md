# ADP CLI User Guide

ADP CLI is the official command-line tool for [Laiye ADP (Agentic Document Processing)](https://github.com/laiye-ai/adp-cli). It provides AI-powered document parsing and intelligent field extraction, converting PDFs, images, and Office documents into structured data.

---

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Document Parsing (parse)](#document-parsing-parse)
- [Field Extraction (extract)](#field-extraction-extract)
- [Application Management (app-id)](#application-management-app-id)
- [Custom Applications (custom-app)](#custom-applications-custom-app)
- [Credit Query (credit)](#credit-query-credit)
- [Global Options](#global-options)
- [Batch Processing & Concurrency](#batch-processing--concurrency)
- [Environment Variables](#environment-variables)
- [Exit Codes](#exit-codes)

---

## Installation

### Via npm

```bash
npm install -g agentic-doc-parse-and-extract-cli
```

### Via GitHub Releases

Download the binary for your platform from [Releases](https://github.com/laiye-ai/adp-cli/releases).

Supported platforms: Windows / Linux / macOS (x64 and arm64).

Verify installation:

```bash
adp version
```

---

## Configuration

An API Key is required before using the CLI. Credentials are stored with AES-256-GCM encryption in `~/.adp/config.json`.

### Set Configuration

```bash
# Set API Key
adp config set --api-key YOUR_API_KEY

# Set API base URL (optional, for private deployments)
adp config set --api-base-url https://your-server.com
```

### View Configuration

```bash
adp config get
```

The API Key is displayed in masked form.

### Clear Configuration

```bash
adp config clear        # Requires confirmation
adp config clear -y     # Skip confirmation
```

---

## Document Parsing (parse)

Convert documents into structured data (Markdown, tables, etc.).

### Parse Local Files

```bash
# Parse a single file
adp parse local ./invoice.pdf --app-id APP_ID

# Parse an entire directory (recursively scans supported files)
adp parse local ./documents/ --app-id APP_ID

# Export results to a file
adp parse local ./invoice.pdf --app-id APP_ID --export ./result.json
```

### Parse from URL

```bash
# Parse a single URL
adp parse url https://example.com/doc.pdf --app-id APP_ID

# Batch parse URLs from a file (one URL per line)
adp parse url ./url_list.txt --app-id APP_ID
```

### Parse Base64 Data

```bash
adp parse base64 BASE64_DATA --app-id APP_ID --file-name invoice.pdf
```

### Async Parsing + Task Query

```bash
# Submit an async task
adp parse local ./doc.pdf --app-id APP_ID --async

# Query task result
adp parse query TASK_ID

# Watch until completion
adp parse query TASK_ID --watch
```

### Common Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--app-id` | string | — | **Required.** Application ID |
| `--async` | bool | false | Async mode |
| `--export` | string | — | Export path (file or directory) |
| `--timeout` | int | 900 | Timeout in seconds |
| `--concurrency` | int | 1 | Concurrent workers (free: max 1, paid: max 2) |
| `--retry` | int | 0 | Retry count with exponential backoff |

---

## Field Extraction (extract)

Extract specific fields (invoice number, amount, date, etc.) from documents using LLM. Subcommands and flags are identical to `parse`.

```bash
# Extract fields from a local file
adp extract local ./invoice.pdf --app-id APP_ID

# Extract fields from a URL
adp extract url https://example.com/doc.pdf --app-id APP_ID

# Batch extract from a directory
adp extract local ./documents/ --app-id APP_ID --concurrency 2 --export ./results/
```

---

## Application Management (app-id)

### List Available Applications

```bash
# List all applications
adp app-id list

# Filter by label
adp app-id list --app-label "invoice"

# List custom applications only
adp app-id list --app-type 1
```

### View Local Cache

```bash
adp app-id cache
```

---

## Custom Applications (custom-app)

### Create a Custom Application

```bash
adp custom-app create \
  --app-name "Invoice Extractor" \
  --parse-mode standard \
  --extract-fields '[{"name":"invoice_no","type":"string","description":"Invoice number"}]'
```

`--extract-fields` accepts an inline JSON string or a path to a JSON file:

```bash
adp custom-app create \
  --app-name "Invoice Extractor" \
  --parse-mode accurate \
  --extract-fields ./fields.json
```

### Update a Custom Application

```bash
adp custom-app update \
  --app-id APP_ID \
  --extract-fields ./updated_fields.json \
  --parse-mode accurate \
  --enable-long-doc true
```

### View Application Config

```bash
adp custom-app get-config --app-id APP_ID

# View a specific version
adp custom-app get-config --app-id APP_ID --config-version 2
```

### Delete Application / Version

```bash
adp custom-app delete --app-id APP_ID
adp custom-app delete-version --app-id APP_ID --config-version 2
```

### AI-Generated Field Recommendations

Upload a sample document and let AI recommend extraction fields:

```bash
# From a local file
adp custom-app ai-generate --app-id APP_ID --file-local ./sample.pdf

# From a URL
adp custom-app ai-generate --app-id APP_ID --file-url https://example.com/sample.pdf
```

---

## Credit Query (credit)

```bash
adp credit
```

---

## Global Options

Available on all commands:

| Option | Description |
|--------|-------------|
| `--json` | Output in JSON format |
| `--quiet` | Suppress all output except errors |
| `--lang en` | Set language to English |
| `--lang zh` | Set language to Chinese |

Examples:

```bash
# JSON output for piping
adp parse local ./doc.pdf --app-id APP_ID --json | jq '.markdown'

# Quiet mode — errors only
adp parse local ./docs/ --app-id APP_ID --quiet --export ./out/
```

---

## Batch Processing & Concurrency

When processing directories or multiple URLs, the CLI automatically enters batch mode:

```bash
# Process all files concurrently with automatic retry
adp parse local ./documents/ --app-id APP_ID --concurrency 2 --retry 3 --export ./results/
```

Batch processing features:
- Individual JSON result files per document, plus a `_summary.json` overview
- Colored progress bar on TTY terminals; JSON-lines progress on non-TTY
- Exit code 6 on partial failure, 1 if all tasks fail

---

## Environment Variables

| Variable | Description | Priority |
|----------|-------------|----------|
| `ADP_API_KEY` | API Key (overrides config file) | Higher than config |
| `ADP_API_BASE_URL` | API base URL | Higher than config |
| `ADP_LANG` | Language (`en` / `zh`) | Higher than system locale |
| `ADP_LOG_LEVEL` | Log level | — |

---

## Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General / network / API error |
| 2 | Parameter error |
| 3 | Resource not found |
| 4 | Permission denied |
| 5 | Conflict error |
| 6 | Batch processing partial failure |

---

## Supported File Formats

PDF, JPG, JPEG, PNG, BMP, TIFF, TIF, DOC, DOCX, XLS, XLSX, PPT, PPTX

Maximum file size: **50 MB**.

---

## Getting Help

```bash
# General help
adp --help

# Subcommand help
adp parse --help
adp parse local --help
adp custom-app create --help

# Command schema as JSON (for AI agent integration)
adp schema
```
