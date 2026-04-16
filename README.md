# ADP CLI

Official command-line tool for [Laiye ADP (Agentic Document Processing)](https://adp-global.laiye.com/) — document parsing and intelligent field extraction.

[English](README.md) | [简体中文](README.zh.md)

## Features

- **Document Parsing** — Turn any type of document into structured data (Markdown or JSON)
- **Invoice/Receipt/Purchase Order Extraction** — Extract key information and line items from invoices, receipts, purchase orders and more
- **Custom Document Extraction** — Extract custom fields from any type of document
- **Batch Processing** — Concurrently process multiple documents in a folder or from URLs
- **Sync/Async** — Support both sync and async processing modes
- **Cross Platform** — Support Windows / Linux / macOS

## Supported Formats

`.jpg` `.jpeg` `.png` `.bmp` `.tiff` `.tif` `.pdf` `.doc` `.docx` `.xls` `.xlsx` `.ppt` `.pptx` (max 50MB per file)

## Install

```bash
# npm (recommended)
npm install -g agentic-doc-parse-and-extract-cli

# Linux / macOS
curl -fsSL https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.sh | bash

# Windows (PowerShell)
irm https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.ps1 | iex
```

Or download a prebuilt binary from [GitHub Releases](https://github.com/laiye-ai/adp-cli/releases).

## Configure

Get an API key at [https://adp-global.laiye.com/](https://adp-global.laiye.com/) (new users get 100 free credits per month).

```bash
adp config set --api-key <your-api-key>
adp config set --api-base-url https://adp-global.laiye.com
adp config get
```

## Quick Examples

```bash
# List available apps
adp app-id list

# Parse a local document
adp parse local ./invoice.pdf --app-id <app-id>

# Extract key fields
adp extract local ./invoice.pdf --app-id <app-id>

# Parse a directory in async mode
adp parse local ./documents/ --app-id <app-id> --async

# Process a remote URL
adp extract url https://example.com/file.pdf --app-id <app-id>

# Query an async task
adp parse query <task-id>

# Check remaining credits
adp credit
```

## Commands

> AI agents should call `adp schema` for the machine-readable, authoritative command spec. The table below is a human-friendly summary.

| Command | Description |
|---|---|
| `adp version` | Print version |
| `adp config set` | Set API key / base URL |
| `adp config get` | Show current config |
| `adp config clear` | Clear config |
| `adp app-id list` | List available apps |
| `adp app-id cache` | Read app list from local cache |
| `adp parse local <path>` | Parse local file/directory |
| `adp parse url <url>` | Parse remote file (URL list file supported) |
| `adp parse base64 <data>` | Parse Base64-encoded content |
| `adp parse query <task-id>` | Query an async parse task |
| `adp extract local <path>` | Extract from local file/directory |
| `adp extract url <url>` | Extract from remote file |
| `adp extract base64 <data>` | Extract from Base64-encoded content |
| `adp extract query <task-id>` | Query an async extract task |
| `adp custom-app create` | Create a custom extraction app |
| `adp custom-app update` | Update custom app config |
| `adp custom-app get-config` | Show app config |
| `adp custom-app delete` | Delete a custom app |
| `adp custom-app delete-version` | Delete a specific config version |
| `adp custom-app ai-generate` | AI-recommend extraction fields |
| `adp credit` | Show remaining credits |
| `adp schema` | Output command schema (for AI agents) |

## Flags

| Flag | Description |
|---|---|
| `--json` | Output JSON |
| `--quiet` | Quiet mode, output result only |
| `--lang <en\|zh>` | Interface language |
| `--app-id` | App ID (required for parse / extract) |
| `--async` | Async mode |
| `--export <path>` | Export result to file |
| `--timeout <seconds>` | Timeout (default 900s) |
| `--concurrency <n>` | Concurrent workers (free: max 1, paid: max 2) |

## Async Workflow

For large files or batch jobs, submit with `--async` and the CLI returns a `task-id`. Poll for results with `parse query` / `extract query`:

```bash
adp parse local ./big.pdf --app-id <app-id> --async
# returns a task-id

adp parse query <task-id>
```

## Environment Variables

| Variable | Description |
|---|---|
| `ADP_API_KEY` | API key (overrides config file) |
| `ADP_API_BASE_URL` | Service URL |
| `ADP_LANG` | Interface language (`en` / `zh`) |
| `ADP_LOG_LEVEL` | Log level (`debug` / `info` / `warn` / `error`) |

## Config Storage

- Config dir: `~/.adp/`
- Config file: `~/.adp/config.json`
- Encrypted API key: `~/.adp/key.enc` (AES-256-GCM)
- App cache: `~/.adp/app_cache.json`
- Version check cache: `~/.adp/version_check.json` (refreshed every 24h)

## License

Commercial license — see [license.md](license.md). Free for non-commercial use (personal learning, research, teaching, open-source community). Commercial use requires written authorization from Laiye Technology. Contact: global_product@laiye.com

## Contributing

Build all platforms with `make build-all VERSION=v1.0.0`. Run E2E tests with `bash tests/test.sh`.
