#!/usr/bin/env node

const { spawnSync } = require("node:child_process");
const path = require("node:path");

const binaryName = process.platform === "win32" ? "tapd.exe" : "tapd";
const binaryPath = path.join(__dirname, "..", "vendor", binaryName);

const result = spawnSync(binaryPath, process.argv.slice(2), {
  stdio: "inherit",
});

if (result.error) {
  console.error(`Failed to run tapd binary at ${binaryPath}: ${result.error.message}`);
  process.exit(1);
}

if (result.signal) {
  console.error(`tapd exited because of signal ${result.signal}`);
  process.exit(1);
}

process.exit(result.status ?? 0);
