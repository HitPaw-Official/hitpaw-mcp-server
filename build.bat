@echo off
chcp 65001 >nul
setlocal

set VERSION=1.0.4
set BINARY_NAME=hitpaw-mcp-server
set BUILD_DIR=build
set LDFLAGS=-s -w -X main.serverVersion=%VERSION%

echo ========================================
echo  HitPaw MCP Server Cross-Platform Build
echo  Version: %VERSION%
echo ========================================
echo.

if exist %BUILD_DIR% rmdir /s /q %BUILD_DIR%
mkdir %BUILD_DIR%

echo [1/5] Building darwin/arm64 (macOS Apple Silicon)...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags "%LDFLAGS%" -o %BUILD_DIR%\%BINARY_NAME%-darwin-arm64 ./cmd/mcp-server/
if %errorlevel% neq 0 goto :error

echo [2/5] Building darwin/amd64 (macOS Intel)...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o %BUILD_DIR%\%BINARY_NAME%-darwin-amd64 ./cmd/mcp-server/
if %errorlevel% neq 0 goto :error

echo [3/5] Building linux/amd64...
set GOOS=linux
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o %BUILD_DIR%\%BINARY_NAME%-linux-amd64 ./cmd/mcp-server/
if %errorlevel% neq 0 goto :error

echo [4/5] Building linux/arm64...
set GOOS=linux
set GOARCH=arm64
go build -ldflags "%LDFLAGS%" -o %BUILD_DIR%\%BINARY_NAME%-linux-arm64 ./cmd/mcp-server/
if %errorlevel% neq 0 goto :error

echo [5/5] Building windows/amd64...
set GOOS=windows
set GOARCH=amd64
go build -ldflags "%LDFLAGS%" -o %BUILD_DIR%\%BINARY_NAME%-windows-amd64.exe ./cmd/mcp-server/
if %errorlevel% neq 0 goto :error

echo.
echo ========================================
echo  Build complete! Output:
echo ========================================
dir %BUILD_DIR%
echo.
echo Next steps:
echo   1. gh release create v%VERSION% %BUILD_DIR%\*
echo   2. cd npm ^&^& npm publish --access public
goto :eof

:error
echo.
echo BUILD FAILED!
exit /b 1