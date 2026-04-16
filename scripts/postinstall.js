#!/usr/bin/env node

'use strict';

const https = require('https');
const fs = require('fs');
const path = require('path');
const os = require('os');
const { execSync } = require('child_process');

// ─── Platform detection ───────────────────────────────────────────────────────

const PLATFORM = process.platform; // 'win32' | 'linux' | 'darwin'
const ARCH = process.arch;         // 'x64' | 'arm64'

const SUPPORTED_PLATFORMS = ['win32', 'linux', 'darwin'];
const SUPPORTED_ARCHS = ['x64', 'arm64'];

if (!SUPPORTED_PLATFORMS.includes(PLATFORM)) {
  console.error(`[adp-cli] Unsupported platform: ${PLATFORM}`);
  console.error(`[adp-cli] Supported platforms: ${SUPPORTED_PLATFORMS.join(', ')}`);
  process.exit(1);
}

if (!SUPPORTED_ARCHS.includes(ARCH)) {
  console.error(`[adp-cli] Unsupported architecture: ${ARCH}`);
  console.error(`[adp-cli] Supported architectures: ${SUPPORTED_ARCHS.join(', ')}`);
  process.exit(1);
}

// ─── Paths ────────────────────────────────────────────────────────────────────

const PKG_ROOT = path.resolve(__dirname, '..');
const BIN_DIR = path.join(PKG_ROOT, 'bin');
const BIN_NAME = PLATFORM === 'win32' ? 'adp.exe' : 'adp';
const BIN_PATH = path.join(BIN_DIR, BIN_NAME);

// ─── Version ──────────────────────────────────────────────────────────────────

const pkg = JSON.parse(fs.readFileSync(path.join(PKG_ROOT, 'package.json'), 'utf8'));
// pkg.version is kept in sync with the GitHub Release tag by release.yml
// (npm version <x.y.z> --no-git-tag-version runs automatically on tag push)
const VERSION = `v${pkg.version}`;

const REPO = 'laiye-ai/adp-cli';

// ─── Download URL ─────────────────────────────────────────────────────────────

const EXT = PLATFORM === 'win32' ? '.zip' : '.tar.gz';
const ARCHIVE_NAME = `adp-${PLATFORM}-${ARCH}${EXT}`;

function buildDownloadUrl(version) {
  return `https://github.com/${REPO}/releases/download/${version}/${ARCHIVE_NAME}`;
}

// Fetch the latest release tag from GitHub API as fallback
function fetchLatestVersion() {
  return new Promise((resolve, reject) => {
    https.get(
      `https://api.github.com/repos/${REPO}/releases/latest`,
      { headers: { 'User-Agent': 'adp-cli-postinstall' } },
      (res) => {
        let data = '';
        res.on('data', (chunk) => { data += chunk; });
        res.on('end', () => {
          try {
            const json = JSON.parse(data);
            if (!json.tag_name) throw new Error('tag_name not found in response');
            resolve(json.tag_name);
          } catch (e) {
            reject(new Error(`Failed to parse latest version: ${e.message}`));
          }
        });
      }
    ).on('error', reject);
  });
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

function log(msg) {
  if (process.env.npm_config_loglevel !== 'silent') {
    console.log(`[adp-cli] ${msg}`);
  }
}

function mkdirSafe(dir) {
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true });
  }
}

function download(url, dest) {
  return new Promise((resolve, reject) => {
    const file = fs.createWriteStream(dest);

    const request = (targetUrl) => {
      https.get(targetUrl, (res) => {
        // Follow redirects (GitHub releases use 302)
        if (res.statusCode === 301 || res.statusCode === 302) {
          file.destroy();
          request(res.headers.location);
          return;
        }

        if (res.statusCode !== 200) {
          file.destroy();
          fs.unlinkSync(dest);
          reject(new Error(`Download failed: HTTP ${res.statusCode} for ${targetUrl}`));
          return;
        }

        const total = parseInt(res.headers['content-length'] || '0', 10);
        let received = 0;
        let lastPct = -1;

        res.on('data', (chunk) => {
          received += chunk.length;
          if (total > 0) {
            const pct = Math.floor((received / total) * 100);
            if (pct !== lastPct && pct % 10 === 0) {
              log(`Downloading... ${pct}%`);
              lastPct = pct;
            }
          }
        });

        res.pipe(file);
        file.on('finish', () => file.close(resolve));
        file.on('error', (err) => {
          fs.unlinkSync(dest);
          reject(err);
        });
      }).on('error', (err) => {
        fs.unlinkSync(dest);
        reject(err);
      });
    };

    request(url);
  });
}

function extractTarGz(archivePath, destDir) {
  execSync(`tar -xzf "${archivePath}" -C "${destDir}"`, { stdio: 'pipe' });
}

function extractZip(archivePath, destDir) {
  // Use PowerShell on Windows (available on Win8+)
  execSync(
    `powershell -NoProfile -Command "Expand-Archive -Path '${archivePath}' -DestinationPath '${destDir}' -Force"`,
    { stdio: 'pipe' }
  );
}

// ─── Main ─────────────────────────────────────────────────────────────────────

async function resolveDownloadUrl() {
  const primaryUrl = buildDownloadUrl(VERSION);

  // Probe the primary URL (HEAD request); fall back to latest if 404
  return new Promise((resolve) => {
    const req = https.request(primaryUrl, { method: 'HEAD' }, (res) => {
      if (res.statusCode === 200 || res.statusCode === 302 || res.statusCode === 301) {
        resolve({ url: primaryUrl, version: VERSION });
      } else if (res.statusCode === 404) {
        log(`Release ${VERSION} not found (404), falling back to latest...`);
        fetchLatestVersion()
          .then((latestVersion) => resolve({ url: buildDownloadUrl(latestVersion), version: latestVersion }))
          .catch(() => resolve({ url: primaryUrl, version: VERSION })); // last resort: try primary anyway
      } else {
        resolve({ url: primaryUrl, version: VERSION });
      }
    });
    req.on('error', () => resolve({ url: primaryUrl, version: VERSION }));
    req.end();
  });
}

async function main() {
  // Skip if binary already exists (e.g. re-install)
  if (fs.existsSync(BIN_PATH)) {
    log(`Binary already exists at ${BIN_PATH}, skipping download.`);
    return;
  }

  log(`Platform : ${PLATFORM}`);
  log(`Arch     : ${ARCH}`);

  const { url: downloadUrl, version: resolvedVersion } = await resolveDownloadUrl();
  log(`Version  : ${resolvedVersion}`);
  log(`Downloading ${downloadUrl}`);

  mkdirSafe(BIN_DIR);

  const tmpDir = fs.mkdtempSync(path.join(os.tmpdir(), 'adp-'));
  const archivePath = path.join(tmpDir, ARCHIVE_NAME);

  try {
    await download(downloadUrl, archivePath);
    log(`Download complete, extracting...`);

    if (PLATFORM === 'win32') {
      extractZip(archivePath, tmpDir);
      // Binary inside zip: adp-win32-x64.exe
      const extracted = path.join(tmpDir, `adp-${PLATFORM}-${ARCH}.exe`);
      fs.renameSync(extracted, BIN_PATH);
    } else {
      extractTarGz(archivePath, tmpDir);
      // Binary inside tar.gz: adp-linux-x64 or adp-darwin-arm64 etc.
      const extracted = path.join(tmpDir, `adp-${PLATFORM}-${ARCH}`);
      fs.renameSync(extracted, BIN_PATH);
      fs.chmodSync(BIN_PATH, 0o755);
    }

    log(`Installed adp binary to ${BIN_PATH}`);
    log(`Run "adp version" to verify the installation.`);
  } finally {
    // Cleanup temp files
    try {
      fs.rmSync(tmpDir, { recursive: true, force: true });
    } catch (_) {}
  }
}

main().catch((err) => {
  console.error(`[adp-cli] Installation failed: ${err.message}`);
  console.error(`[adp-cli] You can install manually:`);
  console.error(`[adp-cli]   Linux/macOS : curl -fsSL https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.sh | bash`);
  console.error(`[adp-cli]   Windows     : irm https://raw.githubusercontent.com/laiye-ai/adp-cli/main/scripts/adp-init.ps1 | iex`);
  process.exit(0); // Exit 0 to not break npm install
});
