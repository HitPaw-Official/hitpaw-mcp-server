#!/usr/bin/env node

/**
 * HitPaw MCP Server - 安装脚本
 *
 * npm install 后自动执行，下载对应平台的 Go 二进制文件。
 * 二进制文件托管在 GitHub Releases 上。
 */

const https = require('https');
const http = require('http');
const fs = require('fs');
const path = require('path');
const os = require('os');

// ==================== 配置 ====================
const REPO = 'HitPaw-Official/hitpaw-mcp-server';
const VERSION = '1.0.4';
const BASE_URL = `https://github.com/${REPO}/releases/download/v${VERSION}`;
// =============================================

function getPlatformInfo() {
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
      console.error(`❌ 不支持的平台: ${platform}/${arch}`);
      process.exit(1);
  }

  return {
    binaryName: `hitpaw-mcp-server-${suffix}`,
    platform,
    arch,
  };
}

function downloadFile(url, destPath) {
  return new Promise((resolve, reject) => {
    console.log(`📥 正在下载: ${url}`);

    const makeRequest = (requestUrl) => {
      const cli = requestUrl.startsWith('https') ? https : http;
      cli.get(requestUrl, (response) => {
        if (response.statusCode === 301 || response.statusCode === 302) {
          makeRequest(response.headers.location);
          return;
        }

        if (response.statusCode !== 200) {
          reject(new Error(`下载失败: HTTP ${response.statusCode}`));
          return;
        }

        const fileStream = fs.createWriteStream(destPath);
        const totalSize = parseInt(response.headers['content-length'] || '0', 10);
        let downloadedSize = 0;

        response.on('data', (chunk) => {
          downloadedSize += chunk.length;
          if (totalSize > 0) {
            const percent = ((downloadedSize / totalSize) * 100).toFixed(1);
            process.stdout.write(`\r  进度: ${percent}% (${formatSize(downloadedSize)}/${formatSize(totalSize)})`);
          }
        });

        response.pipe(fileStream);
        fileStream.on('finish', () => {
          fileStream.close();
          console.log('\n✅ 下载完成');
          resolve();
        });
        fileStream.on('error', (err) => {
          fs.unlinkSync(destPath);
          reject(err);
        });
      }).on('error', reject);
    };

    makeRequest(url);
  });
}

function formatSize(bytes) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / 1024 / 1024).toFixed(1)} MB`;
}

async function main() {
  const info = getPlatformInfo();
  console.log(`🖥️  平台: ${info.platform}/${info.arch}`);
  console.log(`📦 二进制: ${info.binaryName}`);

  const binDir = path.join(__dirname, 'bin');
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  const destPath = path.join(binDir, info.binaryName);

  if (fs.existsSync(destPath)) {
    console.log('✅ 二进制文件已存在，跳过下载');
    return;
  }

  const downloadUrl = `${BASE_URL}/${info.binaryName}`;

  try {
    await downloadFile(downloadUrl, destPath);
    if (info.platform !== 'win32') {
      fs.chmodSync(destPath, 0o755);
    }
    console.log(`✅ 安装成功: ${destPath}`);
  } catch (err) {
    console.error(`❌ 安装失败: ${err.message}`);
    console.error('');
    console.error('你也可以手动下载二进制文件:');
    console.error(`  ${downloadUrl}`);
    console.error(`然后放到: ${binDir}/`);
    console.error('');
    console.error('或者从源码编译:');
    console.error('  git clone https://github.com/' + REPO);
    console.error('  cd hitpaw-mcp-server && make build');
    process.exit(1);
  }
}

main();