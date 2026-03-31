#!/usr/bin/env node

/**
 * HitPaw MCP Server - npm 入口
 *
 * 启动 Go 编译的 MCP Server 二进制文件。
 * 通过 npx @hitpaw/mcp-server 运行时执行此文件。
 */

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

function getBinaryName() {
  const platform = os.platform();
  const arch = os.arch();

  let suffix;
  switch (platform) {
    case 'darwin':
      suffix = arch === 'arm64' ? 'darwin-arm64' : 'darwin-amd64';
      break;
    case 'linux':
      suffix = arch === 'arm64' ? 'linux-arm64' : 'linux-amd64';
      break;
    case 'win32':
      suffix = 'windows-amd64.exe';
      break;
    default:
      console.error(`不支持的平台: ${platform}`);
      process.exit(1);
  }

  return `hitpaw-mcp-server-${suffix}`;
}

function findBinary() {
  const binaryName = getBinaryName();

  const localPath = path.join(__dirname, 'bin', binaryName);
  if (fs.existsSync(localPath)) return localPath;

  const globalPath = path.join(__dirname, binaryName);
  if (fs.existsSync(globalPath)) return globalPath;

  console.error(`找不到 MCP Server 二进制文件: ${binaryName}`);
  console.error('请运行: npm run postinstall 重新下载');
  process.exit(1);
}

const binaryPath = findBinary();

try {
  fs.chmodSync(binaryPath, 0o755);
} catch (e) {
  // Windows 上可能不需要
}

const child = spawn(binaryPath, process.argv.slice(2), {
  stdio: ['inherit', 'inherit', 'inherit'],
  env: process.env,
});

child.on('error', (err) => {
  console.error(`启动 MCP Server 失败: ${err.message}`);
  process.exit(1);
});

child.on('exit', (code) => {
  process.exit(code || 0);
});

process.on('SIGINT', () => child.kill('SIGINT'));
process.on('SIGTERM', () => child.kill('SIGTERM'));