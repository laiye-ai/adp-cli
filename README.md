# ADP CLI

Official command-line tool for [Laiye ADP (Agentic Document Processing)](https://adp-global.laiye.com/) — document parsing and intelligent field extraction.

[English](README.md) | [简体中文](README.zh.md)

## Features

- **Document Parsing** — Turn any type of document into structured data (Markdown or JSON)
- **Invoice/Receipt/Purchase Order Extraction** — Extract key information and line items from invoices, receipts, purchase orders and more
- **Custom Document Extraction** — Extract custom fields from any type of document
- **Batch Processing** — Concurrently process multiple documents in a folder or from URLs, with per-file result output
- **Sync/Async** — Support both sync and async processing modes
- **Two-Phase Async** — `--async --no-wait` submits tasks and outputs task-id list; `query --file` resumes from where you left off
- **Reliability** — Automatic retry with exponential backoff (`--retry`), fine-grained exit codes
- **Cross Platform** — Windows / Linux / macOS, static binaries with no dependencies

## Supported Formats

`.jpg` `.jpeg` `.png` `.bmp` `.tiff` `.tif` `.pdf` `.doc` `.docx` `.xls` `.xlsx` (max 50MB per file)

## Agent Integration

If you are an AI agent, install the ADP skills:

```bash
npx skills add laiye-ai/adp-cli -y -g
```

The skills package will guide you through CLI installation, authentication, and usage automatically.

## Install (Manual)

```bash
# npm (recommended)
npm install -g @laiye-adp/agentic-doc-parse-and-extract-cli

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

# Two-phase async (submit + query separately, resumable)
adp extract local ./documents/ --app-id <app-id> --async --no-wait --export tasks.json
adp extract query --watch --file tasks.json

# Auto retry on failure (up to 2 retries)
adp parse local ./documents/ --app-id <app-id> --retry 2

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
| `adp parse query <task-id...>` | Query async parse tasks (supports multiple IDs or `--file`) |
| `adp extract local <path>` | Extract from local file/directory |
| `adp extract url <url>` | Extract from remote file |
| `adp extract base64 <data>` | Extract from Base64-encoded content |
| `adp extract query <task-id...>` | Query async extract tasks (supports multiple IDs or `--file`) |
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
| `--no-wait` | Submit tasks only, do not wait for results (use with `--async`) |
| `--export <path>` | Export result to file (single file) or directory (batch) |
| `--timeout <seconds>` | Timeout (default 900s) |
| `--concurrency <n>` | Concurrent workers (free: max 1, paid: max 2) |
| `--retry <n>` | Retries for retryable errors (default 0) |
| `--file <path>` | Read task IDs from JSON file (output of `--no-wait`, query only) |

## Async Workflow

For large files or batch jobs, submit with `--async` and the CLI returns a `task-id`. Poll for results with `parse query` / `extract query`:

```bash
adp parse local ./big.pdf --app-id <app-id> --async
# returns a task-id

adp parse query <task-id>
```

### Two-Phase Async (`--no-wait`)

By default, `--async` submits and polls until completion — ideal for AI agents. For resumable workflows, use two-phase mode:

**Phase 1: Submit tasks**

```bash
adp extract local ./documents/ --app-id <app-id> --async --no-wait --export tasks.json
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

## Batch Processing

When processing multiple files/URLs, the CLI writes each result to a separate file:

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

## Exit Codes

| Code | Meaning |
|------|---------|
| `0` | All success |
| `1` | All failed / system error |
| `2` | Parameter error |
| `3` | Resource not found |
| `4` | Permission denied |
| `5` | Conflict |
| `6` | Partial failure (some tasks failed in batch) |

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
