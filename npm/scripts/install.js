#!/usr/bin/env node

const childProcess = require("node:child_process");
const crypto = require("node:crypto");
const fs = require("node:fs");
const https = require("node:https");
const os = require("node:os");
const path = require("node:path");

const packageJson = require("../package.json");

const repo = "go-tapd/cli";
const version = process.env.TAPD_VERSION || packageJson.version;
const tag = version.startsWith("v") ? version : `v${version}`;

const platformNames = {
  darwin: "Darwin",
  linux: "Linux",
  win32: "Windows",
};

const archNames = {
  x64: "x86_64",
  arm64: "arm64",
};

function fail(message) {
  console.error(`tapd npm install failed: ${message}`);
  process.exit(1);
}

function assetName() {
  const platform = platformNames[process.platform];
  const arch = archNames[process.arch];

  if (!platform || !arch) {
    fail(`unsupported platform ${process.platform}/${process.arch}`);
  }

  const extension = process.platform === "win32" ? "zip" : "tar.gz";
  return `tapd_${version.replace(/^v/, "")}_${platform}_${arch}.${extension}`;
}

function request(url, redirects = 0) {
  return new Promise((resolve, reject) => {
    https
      .get(
        url,
        {
          headers: {
            "User-Agent": "@go-tapd/tapd npm installer",
          },
        },
        (response) => {
          const statusCode = response.statusCode || 0;
          const location = response.headers.location;

          if (statusCode >= 300 && statusCode < 400 && location) {
            response.resume();
            if (redirects > 5) {
              reject(new Error(`too many redirects while downloading ${url}`));
              return;
            }
            resolve(request(new URL(location, url).toString(), redirects + 1));
            return;
          }

          if (statusCode < 200 || statusCode >= 300) {
            response.resume();
            reject(new Error(`download failed with HTTP ${statusCode}: ${url}`));
            return;
          }

          resolve(response);
        },
      )
      .on("error", reject);
  });
}

async function downloadText(url) {
  const response = await request(url);
  return new Promise((resolve, reject) => {
    const chunks = [];
    response.on("data", (chunk) => chunks.push(chunk));
    response.on("end", () => resolve(Buffer.concat(chunks).toString("utf8")));
    response.on("error", reject);
  });
}

async function downloadFile(url, destination) {
  const response = await request(url);
  await new Promise((resolve, reject) => {
    const file = fs.createWriteStream(destination);
    response.pipe(file);
    file.on("finish", () => file.close(resolve));
    file.on("error", reject);
    response.on("error", reject);
  });
}

function checksumFor(checksums, name) {
  const line = checksums
    .split(/\r?\n/)
    .find((entry) => entry.includes(name));

  if (!line) {
    fail(`could not find ${name} in checksums.txt`);
  }

  const match = line.match(/[a-f0-9]{64}/i);
  if (!match) {
    fail(`could not parse checksum for ${name}`);
  }

  return match[0].toLowerCase();
}

function sha256(file) {
  const hash = crypto.createHash("sha256");
  hash.update(fs.readFileSync(file));
  return hash.digest("hex");
}

function extractArchive(archive, destination) {
  if (process.platform === "win32") {
    childProcess.execFileSync(
      "powershell",
      [
        "-NoProfile",
        "-ExecutionPolicy",
        "Bypass",
        "-Command",
        "Expand-Archive -LiteralPath $args[0] -DestinationPath $args[1] -Force",
        archive,
        destination,
      ],
      { stdio: "inherit" },
    );
    return;
  }

  childProcess.execFileSync("tar", ["-xzf", archive, "-C", destination], {
    stdio: "inherit",
  });
}

function findBinary(directory) {
  const binaryName = process.platform === "win32" ? "tapd.exe" : "tapd";
  const entries = fs.readdirSync(directory, { withFileTypes: true });

  for (const entry of entries) {
    const entryPath = path.join(directory, entry.name);
    if (entry.isFile() && entry.name === binaryName) {
      return entryPath;
    }
    if (entry.isDirectory()) {
      const found = findBinary(entryPath);
      if (found) {
        return found;
      }
    }
  }

  return "";
}

async function main() {
  if (process.env.TAPD_SKIP_DOWNLOAD) {
    console.log("Skipping tapd binary download because TAPD_SKIP_DOWNLOAD is set.");
    return;
  }

  const name = assetName();
  const baseUrl = `https://github.com/${repo}/releases/download/${tag}`;
  const archiveUrl = `${baseUrl}/${name}`;
  const checksumsUrl = `${baseUrl}/checksums.txt`;
  const tempDir = fs.mkdtempSync(path.join(os.tmpdir(), "tapd-npm-"));
  const archivePath = path.join(tempDir, name);
  const extractDir = path.join(tempDir, "extract");

  fs.mkdirSync(extractDir);

  try {
    console.log(`Downloading ${archiveUrl}`);
    await downloadFile(archiveUrl, archivePath);

    const expected = checksumFor(await downloadText(checksumsUrl), name);
    const actual = sha256(archivePath);

    if (actual !== expected) {
      fail(`checksum mismatch for ${name}: expected ${expected}, got ${actual}`);
    }

    extractArchive(archivePath, extractDir);

    const binary = findBinary(extractDir);
    if (!binary) {
      fail(`tapd binary not found in ${name}`);
    }

    const vendorDir = path.join(__dirname, "..", "vendor");
    const target = path.join(vendorDir, process.platform === "win32" ? "tapd.exe" : "tapd");
    fs.rmSync(vendorDir, { force: true, recursive: true });
    fs.mkdirSync(vendorDir, { recursive: true });
    fs.copyFileSync(binary, target);

    if (process.platform !== "win32") {
      fs.chmodSync(target, 0o755);
    }

    console.log(`Installed tapd ${version} to ${target}`);
  } finally {
    fs.rmSync(tempDir, { force: true, recursive: true });
  }
}

main().catch((error) => fail(error.message));
