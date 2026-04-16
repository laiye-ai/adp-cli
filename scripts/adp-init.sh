#!/usr/bin/env bash
set -euo pipefail

# ─── ADP CLI Installer for Linux / macOS ──────────────────────────────────────
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.sh | bash
# ─────────────────────────────────────────────────────────────────────────────

REPO="laiye-ai/adp-cli"
BINARY_NAME="adp"
INSTALL_DIR="/usr/local/bin"

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

# ─── Get latest version from GitHub ──────────────────────────────────────────
get_latest_version() {
  local version
  version=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep '"tag_name"' \
    | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
  if [[ -z "$version" ]]; then
    err "Failed to fetch latest version from GitHub."
    exit 1
  fi
  echo "$version"
}

# ─── Main ─────────────────────────────────────────────────────────────────────
main() {
  local os arch version archive_name download_url tmp_dir

  os=$(detect_os)
  arch=$(detect_arch)
  version=$(get_latest_version)

  archive_name="adp-${os}-${arch}.tar.gz"
  download_url="https://github.com/${REPO}/releases/download/${version}/${archive_name}"

  log "Platform : ${os}"
  log "Arch     : ${arch}"
  log "Version  : ${version}"
  log "Downloading ${download_url}"

  tmp_dir=$(mktemp -d)
  trap 'rm -rf "$tmp_dir"' EXIT

  curl -fsSL --progress-bar -o "${tmp_dir}/${archive_name}" "$download_url"

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

  # Verify
  if command -v adp &>/dev/null; then
    log "Verification: $(adp version)"
    echo ""
    echo -e "${GREEN}✓ ADP CLI installed successfully!${NC}"
    echo "  Run 'adp config set --api-key YOUR_API_KEY' to get started."
  else
    warn "adp installed to ${install_dest} but is not in PATH."
    warn "Add the following to your shell profile:"
    warn "  export PATH=\"${INSTALL_DIR}:\$PATH\""
  fi
}

main
