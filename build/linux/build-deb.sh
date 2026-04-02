#!/bin/bash
# 构建 .deb 包
set -e

VERSION=${1:-"1.0.0"}
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"
BUILD_DIR="$PROJECT_DIR/build/output"
DEB_STAGING="/tmp/cinaseek-deb-staging"
DEB_PKG="cinaseek_${VERSION}_amd64.deb"

echo "Building .deb package v${VERSION}..."

rm -rf "$DEB_STAGING"
mkdir -p "$DEB_STAGING/DEBIAN"
mkdir -p "$DEB_STAGING/usr/bin"
mkdir -p "$DEB_STAGING/etc/systemd/system"

# Copy binary
cp "$BUILD_DIR/cinaseek-client-linux-amd64" "$DEB_STAGING/usr/bin/cinaseek-client" 2>/dev/null || {
    echo "ERROR: cinaseek-client-linux-amd64 not found"
    exit 1
}
chmod 755 "$DEB_STAGING/usr/bin/cinaseek-client"

# Copy systemd service
cp "$SCRIPT_DIR/cinaseek.service" "$DEB_STAGING/etc/systemd/system/cinaseek.service"

# Create control file
cat > "$DEB_STAGING/DEBIAN/control" << EOF
Package: cinaseek
Version: $VERSION
Section: devel
Priority: optional
Architecture: amd64
Depends: libc6 (>= 2.31)
Maintainer: CinaGroup <dev@cinagroup.com>
Description: CinaSeek - Lightweight VM Web Management Platform
 A lightweight VM remote management platform based on Multipass,
 supporting OpenClaw AI development environment one-click deployment.
EOF

# Build deb
dpkg-deb --build "$DEB_STAGING" "$BUILD_DIR/$DEB_PKG"

rm -rf "$DEB_STAGING"
echo "DEB created: $BUILD_DIR/$DEB_PKG"
