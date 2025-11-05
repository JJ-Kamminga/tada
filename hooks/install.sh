#!/bin/bash
#
# Install Git hooks for the tada project
#

set -e

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}Installing Git hooks...${NC}"
echo ""

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo -e "${YELLOW}⚠ Not in a git repository. Are you in the project root?${NC}"
    exit 1
fi

# Create .git/hooks directory if it doesn't exist
mkdir -p .git/hooks

# Install pre-commit hook
echo -e "${BLUE}📝 Installing pre-commit hook...${NC}"
cp hooks/pre-commit .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
echo -e "${GREEN}✓ Pre-commit hook installed${NC}"
echo ""

# Check if golangci-lint is installed
if ! command -v golangci-lint &> /dev/null; then
    echo -e "${YELLOW}⚠ golangci-lint not found${NC}"
    echo ""
    echo "To install golangci-lint, run:"
    echo "  go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
    echo ""
    echo "The pre-commit hook will use 'go vet' as a fallback"
    echo ""
fi

echo -e "${GREEN}✅ Git hooks installation complete!${NC}"
echo ""
echo "The pre-commit hook will run:"
echo "  1. Code formatting checks (gofmt)"
echo "  2. Linting (golangci-lint or go vet)"
echo "  3. Tests (go test)"
echo ""
echo "To skip the hook when committing, use: git commit --no-verify"
echo ""
