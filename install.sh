#!/usr/bin/env bash
# gx Quick Installation Script for Linux/macOS
# This script builds gx and runs init-install

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║     gx - Go Version Manager            ║${NC}"
echo -e "${BLUE}║     Quick Installation Script          ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed${NC}"
    echo "Please install Go first: https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✓${NC} Found Go: $GO_VERSION"
echo ""

# Build gx
echo -e "${YELLOW}Building gx...${NC}"
if [ -f "build.sh" ]; then
    ./build.sh build
else
    go build -o build/gx ./cmd/gx
fi

if [ ! -f "build/gx" ]; then
    echo -e "${RED}Error: Build failed${NC}"
    exit 1
fi

echo -e "${GREEN}✓${NC} Build complete"
echo ""

# Run init-install
echo -e "${YELLOW}Running installation...${NC}"
echo ""
./build/gx init-install

echo ""
echo -e "${GREEN}╔════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║     Installation Complete!             ║${NC}"
echo -e "${GREEN}╚════════════════════════════════════════╝${NC}"
