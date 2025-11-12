#!/bin/bash

set -e

echo "🎉 Installing tada..."
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed. Please install Go 1.21 or higher."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✓ Found Go version: $GO_VERSION"
echo ""

# Install dependencies
echo "📦 Installing dependencies..."
go mod download
echo "✓ Dependencies installed"
echo ""

# Build the application
echo "🔨 Building tada..."
go build -o tada
echo "✓ Build complete"
echo ""

INSTALL_DIR="$HOME/go/bin"
sudo cp -i tada "$INSTALL_DIR/"
echo "✓ Installed to $INSTALL_DIR/tada"

echo ""
echo "🎊 Installation complete!"
echo ""
echo "Next steps:"
echo "1. Create your todo file: mkdir -p ~/.tada && touch ~/.tada/todo.txt"
echo "2. Run: tada (or ./tada if not installed to PATH)"
echo ""
echo "Optional: Create a shell alias for 'td' by adding this to your shell config:"
echo "  alias td='tada'"
