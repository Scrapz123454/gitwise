#!/usr/bin/env node

const { execSync } = require("child_process");
const fs = require("fs");
const path = require("path");
const os = require("os");
const https = require("https");

const VERSION = require("./package.json").version;
const REPO = "aymenm/gitwise";

function getPlatform() {
  const platform = os.platform();
  const arch = os.arch();

  const platformMap = {
    linux: "linux",
    darwin: "darwin",
    win32: "windows",
  };

  const archMap = {
    x64: "amd64",
    arm64: "arm64",
  };

  const p = platformMap[platform];
  const a = archMap[arch];

  if (!p || !a) {
    console.error(`Unsupported platform: ${platform}/${arch}`);
    process.exit(1);
  }

  return { os: p, arch: a };
}

function download(url) {
  return new Promise((resolve, reject) => {
    https
      .get(url, (res) => {
        if (res.statusCode >= 300 && res.statusCode < 400 && res.headers.location) {
          return download(res.headers.location).then(resolve).catch(reject);
        }
        const chunks = [];
        res.on("data", (chunk) => chunks.push(chunk));
        res.on("end", () => resolve(Buffer.concat(chunks)));
        res.on("error", reject);
      })
      .on("error", reject);
  });
}

async function main() {
  const { os: osName, arch } = getPlatform();
  const ext = osName === "windows" ? "zip" : "tar.gz";
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/gitwise_${VERSION}_${osName}_${arch}.${ext}`;

  console.log(`Downloading gitwise v${VERSION} for ${osName}/${arch}...`);

  const binDir = path.join(__dirname, "bin");
  fs.mkdirSync(binDir, { recursive: true });

  try {
    const data = await download(url);
    const tmpFile = path.join(os.tmpdir(), `gitwise.${ext}`);
    fs.writeFileSync(tmpFile, data);

    if (ext === "tar.gz") {
      execSync(`tar -xzf "${tmpFile}" -C "${binDir}" gitwise`, { stdio: "inherit" });
    } else {
      // For Windows, use PowerShell to extract
      execSync(
        `powershell -command "Expand-Archive -Path '${tmpFile}' -DestinationPath '${binDir}' -Force"`,
        { stdio: "inherit" }
      );
    }

    // Make executable
    const binaryPath = path.join(binDir, osName === "windows" ? "gitwise.exe" : "gitwise");
    if (osName !== "windows") {
      fs.chmodSync(binaryPath, 0o755);
    }

    fs.unlinkSync(tmpFile);
    console.log("gitwise installed successfully!");
  } catch (err) {
    console.error(`Failed to download gitwise: ${err.message}`);
    console.error("You can install it manually from: https://github.com/" + REPO + "/releases");
    process.exit(1);
  }
}

main();
