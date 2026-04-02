#!/bin/bash
# 构建 macOS DMG
# 使用 hdiutil 创建 DMG 镜像
set -e

VERSION=${1:-"1.0.0"}
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"
BUILD_DIR="$PROJECT_DIR/build/output"
DMG_NAME="cinaseek-client-${VERSION}-macos"
DMG_STAGING="/tmp/cinaseek-dmg-staging"

echo "Building macOS DMG v${VERSION}..."

rm -rf "$DMG_STAGING"
mkdir -p "$DMG_STAGING/CinaSeek.app/Contents/MacOS"
mkdir -p "$DMG_STAGING/CinaSeek.app/Contents/Resources"

if [ -f "$BUILD_DIR/cinaseek-client-darwin-arm64" ]; then
    cp "$BUILD_DIR/cinaseek-client-darwin-arm64" "$DMG_STAGING/CinaSeek.app/Contents/MacOS/cinaseek-client"
elif [ -f "$BUILD_DIR/cinaseek-client-darwin-amd64" ]; then
    cp "$BUILD_DIR/cinaseek-client-darwin-amd64" "$DMG_STAGING/CinaSeek.app/Contents/MacOS/cinaseek-client"
else
    echo "ERROR: No macOS binary found in $BUILD_DIR"
    exit 1
fi
chmod +x "$DMG_STAGING/CinaSeek.app/Contents/MacOS/cinaseek-client"

cp "$SCRIPT_DIR/Info.plist" "$DMG_STAGING/CinaSeek.app/Contents/Info.plist"

hdiutil create -volname "CinaSeek" \
    -srcfolder "$DMG_STAGING" \
    -ov -format UDZO \
    "$BUILD_DIR/${DMG_NAME}.dmg"

rm -rf "$DMG_STAGING"
echo "DMG created: $BUILD_DIR/${DMG_NAME}.dmg"
