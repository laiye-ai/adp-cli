# ADP CLI Installer for Windows (PowerShell)
# Usage:
#   irm https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.ps1 | iex
# Or:
#   Invoke-WebRequest -Uri "https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.ps1" -OutFile adp-init.ps1
#   .\adp-init.ps1

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

$REPO        = "laiye-ai/adp-cli"
$BINARY_NAME = "adp.exe"

# ─── Download mirrors (proxy first, GitHub fallback) ─────────────────────────
# Support custom version: $env:ADP_VERSION="v1.2.3"; .\adp-init.ps1
$ADP_VERSION = if ($env:ADP_VERSION) { $env:ADP_VERSION } else { "latest" }
if ($ADP_VERSION -eq "latest") {
  $GITHUB_BASE = "https://github.com/$REPO/releases/latest/download"
} else {
  $GITHUB_BASE = "https://github.com/$REPO/releases/download/$ADP_VERSION"
}
$MIRROR_BASE = "https://ghproxy.net/$GITHUB_BASE"

# ─── Color helpers ────────────────────────────────────────────────────────────
function Log  { param($msg) Write-Host "[adp-cli] $msg" -ForegroundColor Green }
function Warn { param($msg) Write-Host "[adp-cli] $msg" -ForegroundColor Yellow }
function Err  { param($msg) Write-Host "[adp-cli] ERROR: $msg" -ForegroundColor Red }

# ─── Detect Arch ──────────────────────────────────────────────────────────────
function Get-Arch {
  switch ($env:PROCESSOR_ARCHITECTURE) {
    'AMD64' { return 'x64' }
    'ARM64' { return 'arm64' }
    default {
      # Fallback: check if running under WOW64 (32-bit on 64-bit OS)
      if ($env:PROCESSOR_ARCHITEW6432 -eq 'AMD64') { return 'x64' }
      Err "Unsupported architecture: $($env:PROCESSOR_ARCHITECTURE)"
      exit 1
    }
  }
}

# ─── Get install directory ────────────────────────────────────────────────────
function Get-InstallDir {
  # Prefer %LOCALAPPDATA%\Programs\adp (no admin required)
  $dir = Join-Path $env:LOCALAPPDATA "Programs\adp"
  return $dir
}

# ─── Add to PATH (user scope) ─────────────────────────────────────────────────
function Add-ToUserPath {
  param([string]$Dir)

  $currentPath = [Environment]::GetEnvironmentVariable('Path', 'User')
  if ($currentPath -split ';' -contains $Dir) {
    return # Already in PATH
  }

  [Environment]::SetEnvironmentVariable('Path', "$currentPath;$Dir", 'User')
  # Also update current session
  $env:Path += ";$Dir"
  Warn "Added $Dir to your user PATH."
  Warn "Restart your terminal for PATH changes to take effect."
}

# ─── Main ─────────────────────────────────────────────────────────────────────
function Main {
  $arch      = Get-Arch
  $platform  = "win32"

  $archiveName  = "adp-${platform}-${arch}.zip"
  $installDir   = Get-InstallDir
  $installPath  = Join-Path $installDir $BINARY_NAME

  Log "Platform : $platform"
  Log "Arch     : $arch"

  # Create temp dir
  $tmpDir = Join-Path $env:TEMP "adp-install-$(Get-Random)"
  New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

  try {
    $archivePath = Join-Path $tmpDir $archiveName

    # Try mirror first (China-friendly), fallback to GitHub
    $downloaded = $false
    foreach ($baseUrl in @($MIRROR_BASE, $GITHUB_BASE)) {
      $url = "$baseUrl/$archiveName"
      Log "Downloading $url"
      try {
        Invoke-WebRequest -Uri $url -OutFile $archivePath -UseBasicParsing -TimeoutSec 120
        $downloaded = $true
        break
      } catch {
        Warn "Failed, trying next mirror..."
      }
    }

    if (-not $downloaded) {
      Err "All download sources failed. Possible causes:"
      Err "  - Network connectivity issue"
      Err "  - No release found for platform '${platform}-${arch}'"
      Err "  - GitHub/mirror service unavailable"
      Err ""
      Err "Alternative: install via npm (recommended for China mainland):"
      Err "  npm install -g @laiye-adp/agentic-doc-parse-and-extract-cli"
      exit 1
    }

    Log "Extracting..."
    Expand-Archive -Path $archivePath -DestinationPath $tmpDir -Force

    $extractedBin = Get-ChildItem -Path $tmpDir -Recurse -Filter "adp.exe" | Select-Object -First 1
    if (-not $extractedBin) {
      Err "Extracted binary 'adp.exe' not found in archive"
      exit 1
    }

    # Install
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    Move-Item -Path $extractedBin.FullName -Destination $installPath -Force

    Log "Installed adp to $installPath"

    # Add to PATH
    Add-ToUserPath -Dir $installDir

    # Verify
    $adpVersion = & $installPath version 2>&1
    Log "Verification: $adpVersion"

    Write-Host ""
    Write-Host "[adp-cli] ✓ ADP CLI installed successfully!" -ForegroundColor Green
    Write-Host "  Run 'adp config set --api-key YOUR_API_KEY' to get started." -ForegroundColor Cyan

  } finally {
    Remove-Item -Path $tmpDir -Recurse -Force -ErrorAction SilentlyContinue
  }
}

Main
