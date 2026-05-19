#!/usr/bin/env node

const fs = require("node:fs");
const path = require("node:path");

const rawVersion = process.argv[2] || process.env.GITHUB_REF_NAME || process.env.RELEASE_VERSION || "";
const version = rawVersion.replace(/^refs\/tags\//, "").replace(/^v/, "");

if (!/^\d+\.\d+\.\d+(?:[-+][0-9A-Za-z.-]+)?$/.test(version)) {
  console.error(`Invalid release version: ${rawVersion}`);
  process.exit(1);
}

const packagePath = path.join(__dirname, "..", "package.json");
const packageJson = JSON.parse(fs.readFileSync(packagePath, "utf8"));

packageJson.version = version;
fs.writeFileSync(packagePath, `${JSON.stringify(packageJson, null, 2)}\n`);

console.log(`Set npm package version to ${version}`);
