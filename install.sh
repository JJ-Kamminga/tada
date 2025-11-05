#!/bin/bash

set -e

echo "🎉 Installing tada..."
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed. Please install Go 1.21 or higher."
    echo "Visit: https://go.dev/doc/install"
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

# Ask user if they want to install to PATH
echo "Would you like to install tada to your PATH? (y/n)"
read -r response

if [[ "$response" =~ ^([yY][eE][sS]|[yY])$ ]]; then
    # Determine installation directory
    if [[ -w /usr/local/bin ]]; then
        INSTALL_DIR="/usr/local/bin"
        cp tada "$INSTALL_DIR/"
        echo "✓ Installed to $INSTALL_DIR/tada"
    elif [[ -d "$HOME/bin" ]] || mkdir -p "$HOME/bin" 2>/dev/null; then
        INSTALL_DIR="$HOME/bin"
        cp tada "$INSTALL_DIR/"
        echo "✓ Installed to $INSTALL_DIR/tada"

        # Check if ~/bin is in PATH
        if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
            echo ""
            echo "⚠️  Note: $HOME/bin is not in your PATH."
            echo "Add this line to your ~/.bashrc, ~/.zshrc, or shell config:"
            echo "  export PATH=\"\$HOME/bin:\$PATH\""
        fi
    else
        echo "❌ Could not find a suitable installation directory."
        echo "You can manually copy 'tada' to a directory in your PATH."
    fi
else
    echo "✓ Skipped PATH installation. You can run tada with ./tada"
fi

echo ""
echo "🎊 Installation complete!"
echo ""
echo "Next steps:"
echo "1. Create your todo file: mkdir -p ~/.tada && touch ~/.tada/todo.txt"
echo "2. Run: tada (or ./tada if not installed to PATH)"
echo ""
echo "Optional: Create a shell alias for 'td' by adding this to your shell config:"
echo "  alias td='tada'"
