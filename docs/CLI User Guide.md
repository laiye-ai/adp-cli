# ADP CLI User Guide

ADP CLI is the official command-line tool for [Laiye ADP (Agentic Document Processing)](https://github.com/laiye-ai/adp-cli). It provides AI-powered document parsing and intelligent field extraction, converting PDFs, images, and Office documents into structured data.

---

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Document Parsing (parse)](#document-parsing-parse)
- [Field Extraction (extract)](#field-extraction-extract)
- [Async Workflow](#async-workflow)
- [Batch Processing & Concurrency](#batch-processing--concurrency)
- [Application Management (app-id)](#application-management-app-id)
- [Custom Applications (custom-app)](#custom-applications-custom-app)
- [Credit Query (credit)](#credit-query-credit)
- [Command Schema (schema)](#command-schema-schema)
- [Global Options](#global-options)
- [Environment Variables](#environment-variables)
- [Exit Codes](#exit-codes)
- [Config Storage](#config-storage)
- [Supported File Formats](#supported-file-formats)
- [Getting Help](#getting-help)

---

## Installation

### Via npm (Recommended)

```bash
npm install -g @laiye-adp/agentic-doc-parse-and-extract-cli
```

### Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.sh | bash
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.ps1 | iex
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

### parse local / url / base64 Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--app-id` | string | — | **Required.** Application ID |
| `--async` | bool | false | Async mode |
| `--no-wait` | bool | false | Submit tasks only, do not wait for results (use with `--async`) |
| `--export` | string | — | Export path (file or directory) |
| `--timeout` | int | 900 | Timeout in seconds |
| `--concurrency` | int | 1 | Concurrent workers (free: max 1, paid: max 2) |
| `--retry` | int | 0 | Retry count with exponential backoff |
| `--file-name` | string | document | File name (base64 subcommand only) |

### parse query Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--watch` | bool | false | Poll until task completion |
| `--file` | string | — | Read task IDs from JSON file (output of `--no-wait`) |
| `--export` | string | — | Export results to file or directory |
| `--timeout` | int | 900 | Timeout in seconds |
| `--concurrency` | int | 1 | Concurrent query workers |

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

### extract local / url / base64 Flags

Identical to `parse local / url / base64` flags — see the table above.

### extract query Flags

Identical to `parse query` flags — see the table above.

---

## Async Workflow

For large files or batch jobs, submit with `--async` and the CLI returns a `task-id`. Poll for results with `parse query` / `extract query`:

```bash
adp parse local ./big.pdf --app-id APP_ID --async
# returns a task-id

adp parse query TASK_ID
```

### Two-Phase Async (`--no-wait`)

By default, `--async` submits and polls until completion — ideal for AI agents. For resumable workflows, use two-phase mode:

**Phase 1: Submit tasks**

```bash
adp extract local ./documents/ --app-id APP_ID --async --no-wait --export tasks.json
```

Output is a JSON array with task IDs:

```json
[
  {"path": "invoice.pdf", "task_id": "task_abc123"},
  {"path": "contract.pdf", "task_id": "task_def456"}
]
```

**Phase 2: Query results**

```bash
adp extract query --watch --file tasks.json
adp extract query --watch --file tasks.json --export ./results/
```

Even if the CLI crashes mid-way, task IDs in `tasks.json` are preserved — resume anytime with `query --file`.

---

## Batch Processing & Concurrency

When processing directories or multiple URLs, the CLI automatically enters batch mode:

```bash
# Process all files concurrently with automatic retry
adp parse local ./documents/ --app-id APP_ID --concurrency 2 --retry 3 --export ./results/
```

### Output Structure

When batch processing, the CLI writes each result to a separate file:

```
adp_results_20250417_153020/
├── _summary.json              # Summary (total, success, failed, per-file status)
├── invoice_01.pdf.json        # Successful result
├── contract_02.docx.json
└── report_03.pdf.error.json   # Error details
```

- `--export <dir>` — specify output directory
- Without `--export` — auto-creates `adp_results_<timestamp>/`
- Single file — outputs to stdout or the `--export` file path

### Batch Processing Features

- Individual JSON result files per document, plus a `_summary.json` overview
- Colored progress bar on TTY terminals; JSON-lines progress on non-TTY
- Exit code 6 on partial failure, 1 if all tasks fail

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

# Limit number of results
adp app-id list --limit 50
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--app-label` | string | — | Filter by label |
| `--app-type` | int | 0 | Application type (0=all, 1=custom) |
| `--limit` | int | 120 | Maximum number of results |

### View Local Cache

```bash
adp app-id cache
```

---

## Custom Applications (custom-app)

> All `custom-app` subcommands support the `--api-key` flag to specify an API Key (overrides the config file value).

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

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--app-name` | string | — | **Required.** Application name |
| `--extract-fields` | string | — | **Required.** Field definitions (JSON string or file path) |
| `--parse-mode` | string | — | **Required.** Parse mode (`standard` / `accurate` / `fast`) |
| `--app-label` | string | — | Application label |
| `--enable-long-doc` | string | — | Enable long document mode |
| `--long-doc-config` | string | — | Long document configuration |
| `--api-key` | string | — | Specify API Key |

### Update a Custom Application

```bash
adp custom-app update \
  --app-id APP_ID \
  --extract-fields ./updated_fields.json \
  --parse-mode accurate \
  --enable-long-doc true
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--app-id` | string | — | **Required.** Application ID |
| `--extract-fields` | string | — | **Required.** Field definitions |
| `--parse-mode` | string | — | **Required.** Parse mode |
| `--enable-long-doc` | string | — | **Required.** Enable long document mode |
| `--app-name` | string | — | Application name |
| `--app-label` | string | — | Application label |
| `--long-doc-config` | string | — | Long document configuration |
| `--api-key` | string | — | Specify API Key |

### View Application Config

```bash
adp custom-app get-config --app-id APP_ID

# View a specific version
adp custom-app get-config --app-id APP_ID --config-version 2
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--app-id` | string | — | **Required.** Application ID |
| `--config-version` | string | — | Specific config version |
| `--api-key` | string | — | Specify API Key |

### Delete Application / Version

```bash
adp custom-app delete --app-id APP_ID
adp custom-app delete-version --app-id APP_ID --config-version 2
```

| Command | Required Flags |
|---------|----------------|
| `delete` | `--app-id` |
| `delete-version` | `--app-id`, `--config-version` |

### AI-Generated Field Recommendations

Upload a sample document and let AI recommend extraction fields:

```bash
# From a local file
adp custom-app ai-generate --app-id APP_ID --file-local ./sample.pdf

# From a URL
adp custom-app ai-generate --app-id APP_ID --file-url https://example.com/sample.pdf

# From Base64
adp custom-app ai-generate --app-id APP_ID --base64 BASE64_DATA
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--app-id` | string | — | **Required.** Application ID |
| `--file-url` | string | — | Sample file URL (one of three input methods) |
| `--file-local` | string | — | Local sample file path (one of three input methods) |
| `--base64` | string | — | Base64-encoded sample data (one of three input methods) |
| `--api-key` | string | — | Specify API Key |

---

## Credit Query (credit)

```bash
adp credit

# Query with a specific API Key
adp credit --api-key YOUR_API_KEY
```

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--api-key` | string | — | Specify API Key (overrides config file) |

---

## Command Schema (schema)

Output the full command structure in JSON format, suitable for AI agent integration:

```bash
adp schema
```

AI agents should call `adp schema` for the machine-readable, authoritative command spec rather than relying on documentation.

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

## Environment Variables

| Variable | Description | Priority |
|----------|-------------|----------|
| `ADP_API_KEY` | API Key (overrides config file) | Higher than config |
| `ADP_API_BASE_URL` | API base URL | Higher than config |
| `ADP_LANG` | Language (`en` / `zh`) | Higher than system locale |
| `ADP_LOG_LEVEL` | Log level (`debug` / `info` / `warn` / `error`) | — |

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

## Config Storage

| Path | Description |
|------|-------------|
| `~/.adp/` | Config directory |
| `~/.adp/config.json` | Config file |
| `~/.adp/key.enc` | Encrypted API Key (AES-256-GCM) |
| `~/.adp/app_cache.json` | Application list cache |
| `~/.adp/version_check.json` | Version check cache (refreshed every 24h) |

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
