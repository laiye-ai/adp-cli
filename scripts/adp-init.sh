#!/usr/bin/env bash
set -euo pipefail

# ─── ADP CLI Installer for Linux / macOS ──────────────────────────────────────
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.sh | bash
# ─────────────────────────────────────────────────────────────────────────────

REPO="laiye-ai/adp-cli"
BINARY_NAME="adp"
INSTALL_DIR="/usr/local/bin"

# ─── Download mirrors (proxy first, GitHub fallback) ─────────────────────────
# Support custom version: ADP_VERSION=v1.2.3 bash adp-init.sh
ADP_VERSION="${ADP_VERSION:-latest}"
if [[ "$ADP_VERSION" == "latest" ]]; then
  GITHUB_URL="https://github.com/${REPO}/releases/latest/download"
else
  GITHUB_URL="https://github.com/${REPO}/releases/download/${ADP_VERSION}"
fi
MIRROR_URL="https://ghproxy.net/${GITHUB_URL}"

# ─── Color output ─────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

log()  { echo -e "${GREEN}[adp-cli]${NC} $*"; }
warn() { echo -e "${YELLOW}[adp-cli]${NC} $*"; }
err()  { echo -e "${RED}[adp-cli] ERROR:${NC} $*" >&2; }

# ─── Detect OS ────────────────────────────────────────────────────────────────
detect_os() {
  case "$(uname -s)" in
    Linux*)  echo "linux" ;;
    Darwin*) echo "darwin" ;;
    *)
      err "Unsupported OS: $(uname -s)"
      err "This script supports Linux and macOS only."
      err "Windows users: use scripts/adp-init.ps1"
      exit 1
      ;;
  esac
}

# ─── Detect Arch ──────────────────────────────────────────────────────────────
detect_arch() {
  case "$(uname -m)" in
    x86_64|amd64) echo "x64" ;;
    arm64|aarch64) echo "arm64" ;;
    *)
      err "Unsupported architecture: $(uname -m)"
      exit 1
      ;;
  esac
}


# ─── Add to PATH (persistent + current session) ─────────────────────────────
add_to_path() {
  local dir="$1"
  # Already in PATH — skip
  case ":$PATH:" in
    *":$dir:"*) return ;;
  esac

  local profile=""
  if [[ -n "${ZSH_VERSION:-}" ]] || [[ "$SHELL" == */zsh ]]; then
    profile="$HOME/.zshrc"
  elif [[ -f "$HOME/.bashrc" ]]; then
    profile="$HOME/.bashrc"
  elif [[ -f "$HOME/.profile" ]]; then
    profile="$HOME/.profile"
  fi

  if [[ -n "$profile" ]]; then
    if ! grep -qF "$dir" "$profile" 2>/dev/null; then
      echo "export PATH=\"$dir:\$PATH\"" >> "$profile"
      warn "Added $dir to $profile"
    fi
  fi

  # Make it available in the current session immediately
  export PATH="$dir:$PATH"
}

# ─── Main ─────────────────────────────────────────────────────────────────────
main() {
  # Check required commands
  for cmd in curl tar; do
    if ! command -v "$cmd" &>/dev/null; then
      err "'$cmd' is required but not found. Please install it first."
      exit 1
    fi
  done

  local os arch archive_name tmp_dir

  os=$(detect_os)
  arch=$(detect_arch)

  archive_name="adp-${os}-${arch}.tar.gz"

  log "Platform : ${os}"
  log "Arch     : ${arch}"

  tmp_dir=$(mktemp -d)
  trap 'rm -rf "$tmp_dir"' EXIT

  # Try mirror first (China-friendly), fallback to GitHub
  local downloaded=false
  for base_url in "$MIRROR_URL" "$GITHUB_URL"; do
    local url="${base_url}/${archive_name}"
    log "Downloading ${url}"
    if curl -fsSL --connect-timeout 10 --max-time 120 --progress-bar -o "${tmp_dir}/${archive_name}" "$url"; then
      downloaded=true
      break
    fi
    warn "Failed, trying next mirror..."
  done

  if [[ "$downloaded" != "true" ]]; then
    err "All download sources failed. Possible causes:"
    err "  - Network connectivity issue"
    err "  - No release found for platform '${os}-${arch}'"
    err "  - GitHub/mirror service unavailable"
    err ""
    err "Alternative: install via npm (recommended for China mainland):"
    err "  npm install -g @laiye-adp/agentic-doc-parse-and-extract-cli"
    exit 1
  fi

  log "Extracting..."
  tar -xzf "${tmp_dir}/${archive_name}" -C "$tmp_dir"

  # Determine install destination (fallback to ~/.local/bin if no write permission)
  local install_dest="${INSTALL_DIR}/${BINARY_NAME}"
  if [[ ! -w "$INSTALL_DIR" ]]; then
    warn "No write permission to ${INSTALL_DIR}, installing to ~/.local/bin instead."
    INSTALL_DIR="${HOME}/.local/bin"
    mkdir -p "$INSTALL_DIR"
    install_dest="${INSTALL_DIR}/${BINARY_NAME}"
  fi

  mv "${tmp_dir}/adp" "$install_dest"
  chmod +x "$install_dest"

  log "Installed adp to ${install_dest}"

  # Ensure INSTALL_DIR is in PATH
  add_to_path "$INSTALL_DIR"

  # Verify (use absolute path — PATH may not be refreshed in current shell)
  log "Verification: $("$install_dest" version)"
  echo ""
  echo -e "${GREEN}✓ ADP CLI installed successfully!${NC}"
  echo "  Run 'adp config set --api-key YOUR_API_KEY' to get started."

  if ! command -v adp &>/dev/null; then
    warn "adp is not on PATH in the current shell. Restart your terminal,"
    warn "or use the absolute path printed below."
  fi

  # Machine-readable install path (for agents / scripts to parse)
  echo "ADP_INSTALL_PATH=${install_dest}"
}

main
