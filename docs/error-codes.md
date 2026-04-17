# Error Codes / 错误码

[English](#english) | [简体中文](#简体中文)

---

<a id="简体中文"></a>

## 退出码

| 退出码 | 常量 | 含义 |
|--------|------|------|
| 0 | `ExitSuccess` | 全部成功 |
| 1 | `ExitGeneralError` | 网络错误 / API 错误 / 系统错误 / 批量全部失败 |
| 2 | `ExitParameterError` | 参数错误 / 请求格式错误 |
| 3 | `ExitResourceNotFound` | 资源未找到（远程 404 或本地文件不存在） |
| 4 | `ExitPermissionDenied` | 认证失败（401）/ 权限不足（403） |
| 5 | `ExitConflict` | 资源冲突（409）/ 重复创建 |
| 6 | `ExitPartialFailure` | 批量处理部分失败 |

## 错误类型

| 类型 | 值 | 退出码 | 可重试 | 说明 |
|------|----|--------|--------|------|
| `NETWORK_ERROR` | `ErrorTypeNetwork` | 1 | 是 | 传输层故障（DNS、TCP、TLS） |
| `API_ERROR` | `ErrorTypeAPI` | 1 | 是 | 服务端错误（429 限流、5xx、未知状态码） |
| `SYSTEM_ERROR` | `ErrorTypeSystem` | 1 | 否 | 未归类的兜底错误 |
| `PARAM_ERROR` | `ErrorTypeParam` | 2 | 否 | 客户端参数问题（400、JSON 格式、枚举值） |
| `RESOURCE_ERROR` | `ErrorTypeResource` | 3 | 否 | 资源不存在（404、本地文件缺失） |
| `AUTH_ERROR` | `ErrorTypeAuth` | 4 | 否 | 认证/授权失败（401、403） |
| `CONFLICT_ERROR` | `ErrorTypeConflict` | 5 | 否 | 资源冲突（409、重复） |

## 分类流程

错误分类按以下优先级从高到低匹配，命中即返回：

```
错误消息
  │
  ▼
① HTTP 状态码精确匹配
  │  "status code 400" → PARAM_ERROR     (exit 2, 不可重试)
  │  "status code 401" → AUTH_ERROR      (exit 4, 不可重试)
  │  "status code 403" → AUTH_ERROR      (exit 4, 不可重试)
  │  "status code 404" → RESOURCE_ERROR  (exit 3, 不可重试)
  │  "status code 409" → CONFLICT_ERROR  (exit 5, 不可重试)
  │  "status code 429" → API_ERROR       (exit 1, 可重试)
  │  "status code 5xx" → API_ERROR       (exit 1, 可重试)
  ▼
② 认证错误
  │  unauthorized / invalid api key / api key expired / authentication failed
  │  → AUTH_ERROR (exit 4)
  ▼
③ 权限错误
  │  forbidden / permission denied / access denied
  │  → AUTH_ERROR (exit 4)
  ▼
④ 资源未找到
  │  not found / does not exist / version_not_found / app not found
  │  → RESOURCE_ERROR (exit 3)
  ▼
⑤ 资源冲突
  │  conflict / already exists / duplicate
  │  → CONFLICT_ERROR (exit 5)
  ▼
⑥ 本地文件未找到
  │  file not found / no such file / enoent / path not found
  │  → RESOURCE_ERROR (exit 3)
  ▼
⑦ 参数错误
  │  failed to parse json / failed to decode / json decode / json unmarshal
  │  path traversal / invalid path / unsupported file type / invalid value / missing required
  │  → PARAM_ERROR (exit 2)
  ▼
⑧ 网络错误（传输层）
  │  dial tcp / connection refused / connection reset / no such host
  │  dns lookup / econnrefused / econnreset / i/o timeout
  │  tls handshake / certificate / network is unreachable
  │  → NETWORK_ERROR (exit 1, 可重试)
  ▼
⑨ 通用 API 错误
  │  兜底匹配 "status code"
  │  → API_ERROR (exit 1, 可重试)
  ▼
⑩ 系统错误
     以上均未命中
     → SYSTEM_ERROR (exit 1, 不可重试)
```

## 批量处理退出码

| 结果 | 退出码 |
|------|--------|
| 全部成功 | 0 |
| 部分失败 | 6 (`ExitPartialFailure`) |
| 全部失败 | 1 (`ExitGeneralError`) |

## 错误输出格式

### 终端 (TTY)

```
Error: Authentication error: status code 401: unauthorized
  Fix: Check your API key is correct and has not expired.
  Type: AUTH_ERROR
  Code: 4
  Retryable: false
```

### JSON / 非 TTY（Agent 消费）

```json
{
  "type": "AUTH_ERROR",
  "code": 4,
  "message": "Authentication error: status code 401: unauthorized",
  "fix": "Check your API key is correct and has not expired.",
  "retryable": false,
  "details": { "context": "parse local" }
}
```

## 重试机制

仅 `retryable: true` 的错误会被 `--retry` 触发重试：

- **NETWORK_ERROR** — 传输层故障，通常是瞬时的
- **API_ERROR** — 429 限流、5xx 服务端错误

重试策略：指数退避（1s → 2s → 4s → ...），最大重试次数由 `--retry N` 控制（默认 0，不重试）。

---

<a id="english"></a>

## English

## Exit Codes

| Code | Constant | Meaning |
|------|----------|---------|
| 0 | `ExitSuccess` | All succeeded |
| 1 | `ExitGeneralError` | Network / API / system error, or batch all failed |
| 2 | `ExitParameterError` | Invalid parameters or request format |
| 3 | `ExitResourceNotFound` | Resource not found (HTTP 404 or local file missing) |
| 4 | `ExitPermissionDenied` | Authentication (401) or authorization (403) failure |
| 5 | `ExitConflict` | Resource conflict (409) or duplicate |
| 6 | `ExitPartialFailure` | Batch processing partially failed |

## Error Types

| Type | Value | Exit Code | Retryable | Description |
|------|-------|-----------|-----------|-------------|
| `NETWORK_ERROR` | `ErrorTypeNetwork` | 1 | Yes | Transport-level failure (DNS, TCP, TLS) |
| `API_ERROR` | `ErrorTypeAPI` | 1 | Yes | Server error (429 rate limit, 5xx, unknown status) |
| `SYSTEM_ERROR` | `ErrorTypeSystem` | 1 | No | Unclassified fallback error |
| `PARAM_ERROR` | `ErrorTypeParam` | 2 | No | Client parameter issue (400, JSON format, enum) |
| `RESOURCE_ERROR` | `ErrorTypeResource` | 3 | No | Resource not found (404, local file missing) |
| `AUTH_ERROR` | `ErrorTypeAuth` | 4 | No | Authentication/authorization failure (401, 403) |
| `CONFLICT_ERROR` | `ErrorTypeConflict` | 5 | No | Resource conflict (409, duplicate) |

## Classification Flow

Errors are classified by priority (highest to lowest). First match wins:

```
Error message
  │
  ▼
① HTTP status code (exact match)
  │  "status code 400" → PARAM_ERROR     (exit 2, not retryable)
  │  "status code 401" → AUTH_ERROR      (exit 4, not retryable)
  │  "status code 403" → AUTH_ERROR      (exit 4, not retryable)
  │  "status code 404" → RESOURCE_ERROR  (exit 3, not retryable)
  │  "status code 409" → CONFLICT_ERROR  (exit 5, not retryable)
  │  "status code 429" → API_ERROR       (exit 1, retryable)
  │  "status code 5xx" → API_ERROR       (exit 1, retryable)
  ▼
② Auth errors
  │  unauthorized / invalid api key / api key expired / authentication failed
  │  → AUTH_ERROR (exit 4)
  ▼
③ Permission errors
  │  forbidden / permission denied / access denied
  │  → AUTH_ERROR (exit 4)
  ▼
④ Resource not found
  │  not found / does not exist / version_not_found / app not found
  │  → RESOURCE_ERROR (exit 3)
  ▼
⑤ Conflict errors
  │  conflict / already exists / duplicate
  │  → CONFLICT_ERROR (exit 5)
  ▼
⑥ Local file not found
  │  file not found / no such file / enoent / path not found
  │  → RESOURCE_ERROR (exit 3)
  ▼
⑦ Parameter errors
  │  failed to parse json / failed to decode / json decode / json unmarshal
  │  path traversal / invalid path / unsupported file type / invalid value / missing required
  │  → PARAM_ERROR (exit 2)
  ▼
⑧ Network errors (transport-level)
  │  dial tcp / connection refused / connection reset / no such host
  │  dns lookup / econnrefused / econnreset / i/o timeout
  │  tls handshake / certificate / network is unreachable
  │  → NETWORK_ERROR (exit 1, retryable)
  ▼
⑨ Generic API errors
  │  fallback match on "status code"
  │  → API_ERROR (exit 1, retryable)
  ▼
⑩ System error
     none of the above matched
     → SYSTEM_ERROR (exit 1, not retryable)
```

## Batch Processing Exit Codes

| Result | Exit Code |
|--------|-----------|
| All succeeded | 0 |
| Partial failure | 6 (`ExitPartialFailure`) |
| All failed | 1 (`ExitGeneralError`) |

## Error Output Format

### Terminal (TTY)

```
Error: Authentication error: status code 401: unauthorized
  Fix: Check your API key is correct and has not expired.
  Type: AUTH_ERROR
  Code: 4
  Retryable: false
```

### JSON / Non-TTY (for Agent consumption)

```json
{
  "type": "AUTH_ERROR",
  "code": 4,
  "message": "Authentication error: status code 401: unauthorized",
  "fix": "Check your API key is correct and has not expired.",
  "retryable": false,
  "details": { "context": "parse local" }
}
```

## Retry Mechanism

Only errors with `retryable: true` are retried when `--retry` is specified:

- **NETWORK_ERROR** — transient transport failures
- **API_ERROR** — 429 rate limit, 5xx server errors

Strategy: exponential backoff (1s → 2s → 4s → ...), max attempts controlled by `--retry N` (default 0, no retry).
