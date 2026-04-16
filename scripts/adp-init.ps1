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

# ─── Get latest version from GitHub ──────────────────────────────────────────
function Get-LatestVersion {
  try {
    $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest" -UseBasicParsing
    return $response.tag_name
  } catch {
    Err "Failed to fetch latest version from GitHub: $_"
    exit 1
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
  $version   = Get-LatestVersion
  $platform  = "win32"

  $archiveName  = "adp-${platform}-${arch}.zip"
  $downloadUrl  = "https://github.com/$REPO/releases/download/$version/$archiveName"
  $installDir   = Get-InstallDir
  $installPath  = Join-Path $installDir $BINARY_NAME

  Log "Platform : $platform"
  Log "Arch     : $arch"
  Log "Version  : $version"
  Log "Downloading $downloadUrl"

  # Create temp dir
  $tmpDir = Join-Path $env:TEMP "adp-install-$(Get-Random)"
  New-Item -ItemType Directory -Path $tmpDir -Force | Out-Null

  try {
    $archivePath = Join-Path $tmpDir $archiveName
    Invoke-WebRequest -Uri $downloadUrl -OutFile $archivePath -UseBasicParsing

    Log "Extracting..."
    Expand-Archive -Path $archivePath -DestinationPath $tmpDir -Force

    $extractedBin = Join-Path $tmpDir "adp.exe"
    if (-not (Test-Path $extractedBin)) {
      Err "Extracted binary not found at $extractedBin"
      exit 1
    }

    # Install
    New-Item -ItemType Directory -Path $installDir -Force | Out-Null
    Move-Item -Path $extractedBin -Destination $installPath -Force

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
