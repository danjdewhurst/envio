#!/bin/sh
set -e

REPO="danjdewhurst/envio"
INSTALL_DIR="/usr/local/bin"

# Get latest version tag
VERSION=$(curl -sI "https://github.com/${REPO}/releases/latest" | grep -i '^location:' | sed 's|.*/tag/v||;s/[[:space:]]//g')

if [ -z "$VERSION" ]; then
  echo "Error: could not determine latest version" >&2
  exit 1
fi

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64|amd64) ARCH="amd64" ;;
  aarch64|arm64) ARCH="arm64" ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

case "$OS" in
  darwin|linux) ;;
  *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

TARBALL="envio_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/v${VERSION}/${TARBALL}"

echo "Installing envio v${VERSION} (${OS}/${ARCH})..."

TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

curl -sL "$URL" -o "$TMP/$TARBALL"
tar -xzf "$TMP/$TARBALL" -C "$TMP"

if [ -w "$INSTALL_DIR" ]; then
  mv "$TMP/envio" "$INSTALL_DIR/envio"
else
  sudo mv "$TMP/envio" "$INSTALL_DIR/envio"
fi

echo "envio v${VERSION} installed to ${INSTALL_DIR}/envio"
