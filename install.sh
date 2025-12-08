#!/usr/bin/env bash
set -e

APP_NAME="nano-tunnel"
BINARY_URL="https://nano-tunnel.onrender.com/nano-tunnel-linux-amd64"
INSTALL_DIR="/usr/local/bin"
TMP_FILE="/tmp/${APP_NAME}"

echo "Installing ${APP_NAME}..."

# Make sure curl exists
if ! command -v curl >/dev/null 2>&1; then
    echo "curl is required. Install it with: sudo apt install curl"
    exit 1
fi

echo "→ Downloading binary..."
curl -fsSL "$BINARY_URL" -o "$TMP_FILE"

echo "→ Making it executable..."
chmod +x "$TMP_FILE"

echo "→ Moving to ${INSTALL_DIR} (requires sudo)..."
sudo mv "$TMP_FILE" "${INSTALL_DIR}/${APP_NAME}"

echo "✔ Installation complete!"
echo "Run: ${APP_NAME} --help"
