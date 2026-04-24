#!/bin/bash
# Install golang - architecture aware
echo "Installing Go with arch detection..."

set -e

# Allow overriding Go version (default 1.22.1)
GO_VERSION="${GO_VERSION:-1.22.1}"

OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$OS" in
  linux) OS_NAME="linux" ;;
  darwin) OS_NAME="darwin" ;;
  *) echo "Unsupported OS: $OS"; exit 1 ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH_NAME="amd64" ;;
  arm64|aarch64) ARCH_NAME="arm64" ;;
  *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

GO_TAR="go${GO_VERSION}.${OS_NAME}-${ARCH_NAME}.tar.gz"
GO_URL="https://go.dev/dl/${GO_TAR}"

mkdir -p "$HOME/tools"
cd "$HOME/tools"
echo "Downloading $GO_URL ..."
curl -fL -o go.tar "$GO_URL"
tar -xf go.tar

mkdir -p "$HOME/gopath"
grep -q 'export GOPATH=' "$HOME/.bashrc" || echo 'export GOPATH=$HOME/gopath' >> "$HOME/.bashrc"
grep -q 'export PATH=.*$HOME/tools/go/bin' "$HOME/.bashrc" || echo 'export PATH=$PATH:$HOME/tools/go/bin:$GOPATH/bin' >> "$HOME/.bashrc"
export GOPATH="$HOME/gopath"
export PATH="$PATH:$HOME/tools/go/bin:$GOPATH/bin"

echo "Installing gopls..."
go install golang.org/x/tools/gopls@latest

echo "Go installation complete!"
