#!/bin/bash
set -e

VERSION=${1:-"1.0.0"}
BINARY_NAME="hitpaw-mcp-server"
BUILD_DIR="build"
NPM_BIN_DIR="npm/bin"
LDFLAGS="-s -w -X main.serverVersion=${VERSION}"

echo "🔨 开始跨平台编译 v${VERSION}..."
echo ""

rm -rf ${BUILD_DIR} ${NPM_BIN_DIR}
mkdir -p ${BUILD_DIR} ${NPM_BIN_DIR}

PLATFORMS=(
  "darwin/arm64"
  "darwin/amd64"
  "linux/amd64"
  "linux/arm64"
  "windows/amd64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
  GOOS="${PLATFORM%/*}"
  GOARCH="${PLATFORM#*/}"

  OUTPUT="${BINARY_NAME}-${GOOS}-${GOARCH}"
  if [ "${GOOS}" = "windows" ]; then
    OUTPUT="${OUTPUT}.exe"
  fi

  echo "  编译 ${GOOS}/${GOARCH} → ${OUTPUT}"

  GOOS=${GOOS} GOARCH=${GOARCH} go build \
    -ldflags "${LDFLAGS}" \
    -o "${BUILD_DIR}/${OUTPUT}" \
    ./cmd/mcp-server/

  cp "${BUILD_DIR}/${OUTPUT}" "${NPM_BIN_DIR}/${OUTPUT}"
done

echo ""
echo "✅ 编译完成！产物目录："
echo ""
ls -lh ${BUILD_DIR}/
echo ""
echo "下一步："
echo "  1. 创建 GitHub Release: gh release create v${VERSION} ${BUILD_DIR}/*"
echo "  2. 发布 npm 包: cd npm && npm publish --access public"
