# ADP CLI — Project Rules

## Project Overview

ADP CLI is a Go-based command-line tool for Laiye's Agentic Document Processing platform. It supports both human terminal use and AI Agent invocation. Module path: `github.com/laiye-ai/adp-cli`, Go 1.21, CGO disabled.

Key dependencies: cobra v1.8.0, viper v1.18.2, zerolog v1.31.0.

## Adding a New Command — Checklist

Every new command **must** update all of the following. Skip any = broken feature.

| # | File | What to do |
|---|------|------------|
| 1 | `cmd/<group>.go` | Command definition, flags, `init()` registration |
| 2 | `internal/api/client.go` | API methods (**strong-typed signatures**, see rules below) |
| 3 | `internal/i18n/i18n.go` | en + zh entries for descriptions, flag help, success messages |
| 4 | `cmd/root.go` | `reloadCommandTranslations()`: `Short`/`Long` + `updateFlagHelp` for every flag |
| 5 | `cmd/schema.go` | `GetFullSchema()`: command entry with all options (name, type, required, default, description, enum) |
| 6 | `tests/test.sh` | Offline tests (`--help`, `schema <group>`) + API tests if applicable |
| 7 | `README.md` + `README.zh.md` | Commands table row; Features table if new command group |
| 8 | `skills/.../SKILL.md` | Agent-facing usage guide, so Agents learn the new command without reading source |

---

## Agent-Friendly Design Principles

This CLI is built for both humans and AI Agents. Every new command **must** preserve these properties.

### 1. Machine-Readable Command Spec (`adp schema`)

`adp schema` outputs the full command tree as JSON (commands, subcommands, options, types, defaults, enums). Agents read this instead of parsing `--help`. **Missing schema = command invisible to Agents.**

### 2. Stable, Parseable Output

- **stdout** = data only (`PrintJSON`). Agents pipe stdout into JSON parsers.
- **stderr** = human messages (`PrintSuccess`, `PrintInfo`, `PrintWarning`, progress, errors).
- **Never mix.** `--json` forces JSON output globally; even progress becomes JSON Lines. `--quiet` suppresses all non-error messages.
- Never use `fmt.Println` — it bypasses the formatter and breaks `--json` mode.

### 3. Stable Exit Codes

| Code | Meaning | Constant |
|------|---------|----------|
| 0 | Success | — |
| 1 | General/system error | `ExitGeneralError` |
| 2 | Parameter error | `ExitParameterError` |
| 3 | Resource not found | `ExitResourceNotFound` |
| 4 | Permission denied | `ExitPermissionDenied` |
| 5 | Conflict | `ExitConflict` |
| 6 | Partial failure (batch) | `ExitPartialFailure` |

**Never invent new codes.** Code 6 is critical — Agents distinguish "all failed" from "some failed".

### 4. Self-Describing Errors

Every error is a structured `CLIError` with `Type`, `Code`, `Message`, `Fix` (actionable hint), `Retryable` (bool). In non-TTY / `--json` mode, errors serialize as JSON to stderr. **Always provide a `Fix` string** — silent errors waste Agent retries.

### 5. Resumable Async Workflows

Long jobs use **two-phase async** so Agents can crash and resume:

- **Phase 1**: `--async --no-wait --export tasks.json` → submits, writes task IDs to file.
- **Phase 2**: `query --watch --file tasks.json` → polls until done.

New async commands **must** support both phases.

### 6. Idempotent and Predictable

- Required flags use `MarkFlagRequired` — cobra rejects malformed calls before any side effect.
- `--retry <n>` retries only `Retryable=true` errors with exponential backoff.
- Destructive ops (delete) should be idempotent — no error if already gone.

### 7. Discoverable

- `adp config get` → `"configured": true/false`. Agents check before calling API.
- `adp app-id list` / `adp app-id cache` → runtime app ID discovery, never hardcode.
- `adp credit` → quota check before bulk runs.
- `adp version` → version probe.

### 8. Caller Attribution

`--source <name>` (e.g. `--source claude`) flows into telemetry. Agents/Skills should always pass it.

### 9. i18n Without Surprises

- `--lang en|zh` or `ADP_LANG` env var controls language.
- Error matching uses `code` field, not localized `message`.
- Schema output respects `--lang` for Agent-facing descriptions.

### Anti-Patterns

- Writing data to stderr or messages to stdout
- Exit code 1 for parameter errors (use 2)
- Inventing non-standard flag names
- Matching errors by human-readable message strings
- Adding a command without `schema.go` registration
- Swallowing errors silently
- Long-running ops without `--async` support

---

## Architecture

```
cmd/                  Cobra commands, one file per group
  root.go             Root command, PersistentPreRun, reloadCommandTranslations
  parse.go            parse + shared helpers (processLocalFiles, processURLs, queryTasks, initClientWithConfig)
  extract.go          extract (delegates to shared.go helpers with mode="extract")
  customapp.go        custom-app + helpers (parseJSONParam, parseStringList, loadConfigWithOverride)
  humanreview.go      human-review (delegates to shared.go helpers with mode="human-review")
  webhook.go          webhook
  shared.go           Shared helpers (ensureMaxConcurrency, initClientWithConfig, processLocalFiles, processURLs, processBase64, queryTasks, checkAPIResponse, batchProcess, batchSubmit, retryWithBackoff)
  schema.go           adp schema output (machine-readable command spec)
internal/api/         HTTP API client (single file client.go)
internal/config/      Config load/save, AES-256-GCM key encryption
internal/errors/      CLIError, ClassifyException, exit codes (0-6)
internal/formatter/   Output formatting (PrintSuccess/JSON/ExitWithError, TTY detection)
internal/i18n/        Bilingual translations (en/zh), T() lookup
internal/file/        File validation, URL list reading, JSON output
internal/telemetry/   Command usage tracking (Begin/End/SetError)
internal/updater/     Async version check via GitHub releases (24h cache)
tests/                E2E bash test suite + samples/
scripts/              Install scripts (adp-init.sh, adp-init.ps1)
skills/               Agent skills definition (SKILL.md)
npm/                  Platform-specific npm packages for distribution
```

---

## Command Pattern

### Template

```go
// Package-level var. Short == Long, both use i18n.T().
var myCmd = &cobra.Command{
    Use:   "my-command",
    Short: i18n.T("my_command_title"),
    Long:  i18n.T("my_command_title"),
    Run: func(cmd *cobra.Command, args []string) {
        // 1. Read flags
        appID, _ := cmd.Flags().GetString("app-id")

        // 2. Get client (validates config, exits if not configured)
        client, _ := initClientWithConfig("my-command")

        // 3. Call API (strong-typed method)
        result, err := client.DoSomething(appID, ...)
        if err != nil {
            formatterOut.ExitWithError(errors.ClassifyException(err, "my-command"))
        }
        checkAPIResponse(result, "my-command")  // catch HTTP 200 with error code

        // 4. Output
        formatterOut.PrintSuccess(i18n.T("my_command_success"))
        formatterOut.PrintJSON(result)
    },
}

func init() {
    parentCmd.AddCommand(myCmd)
    myCmd.Flags().String("app-id", "", i18n.T("my_option_app_id"))
    myCmd.MarkFlagRequired("app-id")
}
```

### Client Initialization — Two Patterns

| Pattern | When to use | Defined in |
|---------|-------------|-----------|
| `initClientWithConfig(mode)` | Default. Validates `IsConfigured`, exits with "run adp config set" hint. | `parse.go` |
| `loadConfigWithOverride(apiKey)` | Command has `--api-key` override flag. | `customapp.go` |

### Error Handling Chain

```
API client returns error
  → errors.ClassifyException(err, context) → *CLIError with Type/Code/Fix/Retryable
    → formatterOut.ExitWithError(cliErr) → stderr output + os.Exit(code)

API returns HTTP 200 but code != "success"
  → checkAPIResponse(result, context) → ExitWithError if error code found

Validation error (bad flag value, missing param)
  → errors.NewCLIError(msg, type, exitCode, retryable, fix, details) → ExitWithError
```

---

## API Client Rules (`internal/api/client.go`)

- **Strong-typed signatures.** Parameters are explicit Go types. Payload `map[string]interface{}` is built inside the client method, never passed from cmd layer.
- Optional params: use pointers (`*string`, `*bool`) or nil-able slices (`[]string`).
- Async methods return `(string, error)` for task ID.
- Sync/query methods return `(map[string]interface{}, error)`.
- All requests go through `c.request(method, endpoint, data)`.

```go
// Correct — strong-typed, payload in client layer
func (c *Client) CreateRule(appID, name string, status bool, rules []map[string]interface{}, logic int) (map[string]interface{}, error) {
    data := map[string]interface{}{"app_id": appID, "rule_name": name, ...}
    return c.request("POST", "...", data)
}

// Wrong — leaks wire format to cmd layer, typo-prone
func (c *Client) CreateRule(data map[string]interface{}) (map[string]interface{}, error) { ... }
```

---

## i18n Rules (`internal/i18n/i18n.go`)

### Key Naming Convention

```
{group}_{subcommand}_{element}
```

| Type | Example |
|------|---------|
| Group description | `human_review_description` |
| Command description | `human_review_rule_create_title` |
| Flag description | `human_review_option_rule_name` |
| Shared flag | `option_timeout`, `option_async`, `option_export` |
| Success message | `human_review_rule_created` |

### Adding Translations — 3 Steps

1. Add entry to **both** `en` and `zh` maps in `internal/i18n/i18n.go`
2. Add `cmd.Short` / `cmd.Long` reload in `reloadCommandTranslations()` in `cmd/root.go`
3. Add `updateFlagHelp(cmd, "flag-name", i18n.T("key"))` for **every** flag on that command

Missing step 2/3 = `--lang` switching silently broken for that command.

---

## Flag Conventions

- **kebab-case only**: `--app-id`, `--rule-name`, `--webhook-url` (never underscores)
- Reuse standard flags before inventing new ones:

| Flag | Type | Default | Scope |
|------|------|---------|-------|
| `--app-id` | string | required | Most commands |
| `--async` | bool | false | parse / extract / task-create |
| `--no-wait` | bool | false | With `--async` |
| `--export` | string | "" | Result output path |
| `--timeout` | int | 900 | All API commands |
| `--concurrency` | int | 1 | Batch operations |
| `--retry` | int | 0 | Batch operations |
| `--watch` | bool | false | Query commands |
| `--file` | string | "" | Query commands (read task IDs) |

---

## Helper Reuse

Before writing new helpers, check if one already exists:

| Need | Use | Defined in |
|------|-----|-----------|
| Parse JSON string or file path | `parseJSONParam(value)` | `customapp.go` |
| Parse comma-separated or JSON array | `parseStringList(value)` | `customapp.go` |
| Resolve URL or URL list file | `resolveURLInput(input)` | `parse.go` |
| Process local files (sync/async/batch) | `processLocalFiles(...)` | `parse.go` |
| Process URLs (sync/async/batch) | `processURLs(...)` | `parse.go` |
| Query async tasks with polling | `queryTasks(...)` | `parse.go` |
| Check API response code field | `checkAPIResponse(result, ctx)` | `helpers.go` |
| Validate concurrency vs user tier | `validateConcurrency(client, n)` | `parse.go` |

Cross-file shared helpers go in `cmd/helpers.go`. Group-specific helpers stay in their own file.

---

## Tests (`tests/test.sh`)

- E2E bash tests, no Go `_test.go` files.
- **Offline tests** (always run, no API key): `--help` for every command, `schema <group>` for every group.
- **API tests** (require `ADP_API_KEY` + `ADP_API_BASE_URL`): Full CRUD lifecycle per command group.
- Pattern: `run_test "test name" $ADP command args`
- CI (`ci.yml`) runs offline tests on every push/PR; API tests run when secrets are set.
- New command groups need: help test (offline), schema test (offline), API section (numbered).

---

## README Updates

Both `README.md` (English) and `README.zh.md` (Chinese) must be updated **together**, identical structure:

- **Commands table**: Add row for every new subcommand
- **Features table**: Add row if adding a new command group
- **Flags table**: Add row if introducing a new shared flag
- **Quick Examples**: Add example if workflow is non-obvious

---

## Telemetry (`internal/telemetry/`)

- `telemetry.Begin()` in `PersistentPreRun`, `telemetry.End()` in `PersistentPostRun` + exit hook.
- Skipped for: `config`, `version`, `help`, bare `adp`.
- Tracks: command, params (user flags excl. `--json`/`--quiet`/`--lang`/`--source`), duration, status, error, OS, arch, cli_version, source.
- Async send with 3s timeout, fire-and-forget.

---

## Build and Release

- `make build-all` cross-compiles for 6 targets (win/linux/darwin × x64/arm64), CGO_ENABLED=0.
- Version injected via ldflags: `-X $(MODULE)/cmd.version=$(VERSION)`.
- Release triggered by `v*` tag → GitHub Release + npm publish.
- Beta: versions containing `-` publish to npm beta tag.

---

## Environment Variables

| Variable | Description |
|----------|-------------|
| `ADP_API_KEY` | API key (overrides config file) |
| `ADP_API_BASE_URL` | Service URL |
| `ADP_LANG` | Interface language (`en` / `zh`) |
| `ADP_LOG_LEVEL` | Log level (`debug` / `info` / `warn` / `error`) |

---

## Code Style

- Go 1.21, use `interface{}` (not `any`) — matches existing codebase
- One file per command group in `cmd/`
- Don't duplicate helpers — check "Helper Reuse" table first
- Parent commands (groups) have no `Run` func; only leaf commands have `Run`
- `Short` and `Long` always set to the same `i18n.T()` call
