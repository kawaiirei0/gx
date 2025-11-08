#!/usr/bin/env bash
# gx Build Script for Linux/macOS
# Bash build script for users without Make

set -e

# Configuration
BINARY_NAME="gx"
MAIN_PACKAGE="./cmd/gx"
BUILD_DIR="build"
DIST_DIR="dist"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Get version information
if [ -z "$VERSION" ]; then
    VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
fi

COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# Build flags
LDFLAGS="-X main.Version=$VERSION -X main.Commit=$COMMIT -X main.BuildDate=$BUILD_DATE -s -w"

# Platform configurations
declare -a PLATFORMS=(
    "windows/amd64"
    "windows/386"
    "linux/amd64"
    "linux/386"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
)

show_help() {
    echo "gx Build Script for Linux/macOS"
    echo ""
    echo "Usage: ./build.sh [target]"
    echo ""
    echo "Available targets:"
    echo "  build          - Build for current platform (default)"
    echo "  build-all      - Build for all supported platforms"
    echo "  release        - Create release packages for all platforms"
    echo "  install        - Install gx to GOPATH/bin"
    echo "  clean          - Remove build artifacts"
    echo "  test           - Run tests"
    echo "  version        - Show version information"
    echo ""
    echo "Environment variables:"
    echo "  VERSION        - Override version (default: git describe)"
    echo ""
    echo "Examples:"
    echo "  ./build.sh build"
    echo "  VERSION=v1.0.0 ./build.sh build-all"
    echo "  ./build.sh release"
}

build_current() {
    echo -e "${GREEN}Building $BINARY_NAME for current platform...${NC}"
    
    mkdir -p "$BUILD_DIR"
    
    # Detect current platform
    local current_os
    local current_arch
    
    case "$OSTYPE" in
        linux*)   current_os="linux" ;;
        darwin*)  current_os="darwin" ;;
        msys*|win32*|cygwin*) current_os="windows" ;;
        *)        current_os=$(uname -s | tr '[:upper:]' '[:lower:]') ;;
    esac
    
    case $(uname -m) in
        x86_64|amd64) current_arch="amd64" ;;
        i386|i686)    current_arch="386" ;;
        aarch64|arm64) current_arch="arm64" ;;
        *)            current_arch=$(uname -m) ;;
    esac
    
    local output_path="$BUILD_DIR/$BINARY_NAME"
    if [ "$current_os" = "windows" ]; then
        output_path="$output_path.exe"
    fi
    
    # Set GOOS and GOARCH explicitly for current platform
    GOOS=$current_os GOARCH=$current_arch go build -ldflags "$LDFLAGS" -o "$output_path" "$MAIN_PACKAGE"
    
    echo -e "${GREEN}Build complete: $output_path${NC}"
}

build_all() {
    echo -e "${GREEN}Building for all platforms...${NC}"
    
    clean_build
    
    for platform in "${PLATFORMS[@]}"; do
        IFS='/' read -r os arch <<< "$platform"
        platform_dir="${os}_${arch}"
        output_dir="$BUILD_DIR/$platform_dir"
        
        echo -e "${GREEN}Building for $os/$arch...${NC}"
        
        mkdir -p "$output_dir"
        
        local binary_ext=""
        if [ "$os" = "windows" ]; then
            binary_ext=".exe"
        fi
        
        output_path="$output_dir/$BINARY_NAME$binary_ext"
        
        GOOS=$os GOARCH=$arch go build -ldflags "$LDFLAGS" -o "$output_path" "$MAIN_PACKAGE"
        
        if [ $? -ne 0 ]; then
            echo -e "${RED}Build failed for $os/$arch!${NC}"
            exit 1
        fi
    done
    
    echo -e "${GREEN}All builds complete!${NC}"
}

create_release() {
    echo -e "${GREEN}Creating release packages...${NC}"
    
    build_all
    
    mkdir -p "$DIST_DIR"
    
    for platform in "${PLATFORMS[@]}"; do
        IFS='/' read -r os arch <<< "$platform"
        platform_dir="${os}_${arch}"
        archive_name="$BINARY_NAME-$VERSION-$os-$arch"
        
        echo -e "${GREEN}Packaging $archive_name...${NC}"
        
        source_dir="$BUILD_DIR/$platform_dir"
        
        if [ "$os" = "windows" ]; then
            archive_path="$DIST_DIR/$archive_name.zip"
            (cd "$source_dir" && zip -q "../../../$archive_path" "$BINARY_NAME.exe")
        else
            archive_path="$DIST_DIR/$archive_name.tar.gz"
            tar -czf "$archive_path" -C "$source_dir" "$BINARY_NAME"
        fi
        
        if [ -f "$archive_path" ]; then
            size=$(du -h "$archive_path" | cut -f1)
            echo -e "  ${GREEN}Created: $archive_path ($size)${NC}"
        fi
    done
    
    echo ""
    echo -e "${GREEN}Release packages created in $DIST_DIR/${NC}"
    ls -lh "$DIST_DIR/"
    
    # Generate checksums
    echo ""
    echo -e "${GREEN}Generating checksums...${NC}"
    (cd "$DIST_DIR" && shasum -a 256 * > checksums.txt)
    echo -e "${GREEN}Checksums saved to $DIST_DIR/checksums.txt${NC}"
}

install_binary() {
    echo -e "${GREEN}Installing $BINARY_NAME...${NC}"
    
    go install -ldflags "$LDFLAGS" "$MAIN_PACKAGE"
    
    gopath=$(go env GOPATH)
    echo -e "${GREEN}Installed to $gopath/bin/$BINARY_NAME${NC}"
}

clean_build() {
    echo -e "${GREEN}Cleaning build artifacts...${NC}"
    
    rm -rf "$BUILD_DIR" "$DIST_DIR"
    
    echo -e "${GREEN}Clean complete${NC}"
}

run_tests() {
    echo -e "${GREEN}Running tests...${NC}"
    
    go test -v -race -coverprofile=coverage.out ./...
    
    echo -e "${GREEN}Tests passed!${NC}"
}

show_version() {
    echo "Version: $VERSION"
    echo "Commit: $COMMIT"
    echo "Build Date: $BUILD_DATE"
}

# Main execution
TARGET="${1:-build}"

case "$TARGET" in
    build)
        build_current
        ;;
    build-all)
        build_all
        ;;
    release)
        create_release
        ;;
    install)
        install_binary
        ;;
    clean)
        clean_build
        ;;
    test)
        run_tests
        ;;
    version)
        show_version
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        echo -e "${RED}Unknown target: $TARGET${NC}"
        echo ""
        show_help
        exit 1
        ;;
esac
