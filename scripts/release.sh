#!/usr/bin/env bash
# Release helper script
# Automates the release process

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
CHANGELOG="CHANGELOG.md"
README="README.md"

show_help() {
    echo "gx Release Helper"
    echo ""
    echo "Usage: ./scripts/release.sh [version]"
    echo ""
    echo "Arguments:"
    echo "  version    Version to release (e.g., v1.0.0)"
    echo ""
    echo "This script will:"
    echo "  1. Validate the version format"
    echo "  2. Check for uncommitted changes"
    echo "  3. Run tests"
    echo "  4. Update version references"
    echo "  5. Create a git tag"
    echo "  6. Build release packages"
    echo ""
    echo "Example:"
    echo "  ./scripts/release.sh v1.0.0"
}

validate_version() {
    local version=$1
    
    if [[ ! $version =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9]+)?$ ]]; then
        echo -e "${RED}Error: Invalid version format: $version${NC}"
        echo "Version must be in format: vX.Y.Z or vX.Y.Z-suffix"
        echo "Examples: v1.0.0, v1.2.3-beta, v2.0.0-rc1"
        exit 1
    fi
}

check_git_status() {
    if [[ -n $(git status -s) ]]; then
        echo -e "${RED}Error: Working directory is not clean${NC}"
        echo "Please commit or stash your changes before releasing"
        git status -s
        exit 1
    fi
}

check_branch() {
    local current_branch=$(git branch --show-current)
    
    if [[ "$current_branch" != "main" && "$current_branch" != "master" ]]; then
        echo -e "${YELLOW}Warning: You are not on main/master branch${NC}"
        echo "Current branch: $current_branch"
        read -p "Continue anyway? (y/N) " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            exit 1
        fi
    fi
}

run_tests() {
    echo -e "${BLUE}Running tests...${NC}"
    
    if ! make test; then
        echo -e "${RED}Error: Tests failed${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}Tests passed!${NC}"
}

check_tag_exists() {
    local version=$1
    
    if git rev-parse "$version" >/dev/null 2>&1; then
        echo -e "${RED}Error: Tag $version already exists${NC}"
        exit 1
    fi
}

update_changelog() {
    local version=$1
    local date=$(date +%Y-%m-%d)
    
    if [[ ! -f "$CHANGELOG" ]]; then
        echo -e "${YELLOW}Warning: $CHANGELOG not found, skipping...${NC}"
        return
    fi
    
    echo -e "${BLUE}Updating $CHANGELOG...${NC}"
    
    # Check if version already exists in changelog
    if grep -q "## \[$version\]" "$CHANGELOG"; then
        echo -e "${GREEN}Version $version already in changelog${NC}"
    else
        echo -e "${YELLOW}Please update $CHANGELOG manually${NC}"
        read -p "Press enter when ready to continue..."
    fi
}

create_tag() {
    local version=$1
    
    echo -e "${BLUE}Creating git tag $version...${NC}"
    
    git tag -a "$version" -m "Release $version"
    
    echo -e "${GREEN}Tag created: $version${NC}"
}

build_release() {
    local version=$1
    
    echo -e "${BLUE}Building release packages...${NC}"
    
    VERSION=$version make release
    
    echo -e "${GREEN}Release packages built successfully!${NC}"
    echo ""
    echo "Release artifacts:"
    ls -lh dist/
}

show_next_steps() {
    local version=$1
    
    echo ""
    echo -e "${GREEN}Release $version prepared successfully!${NC}"
    echo ""
    echo "Next steps:"
    echo "  1. Review the release packages in dist/"
    echo "  2. Push the tag: git push origin $version"
    echo "  3. GitHub Actions will automatically create a release"
    echo ""
    echo "Or manually create a GitHub release:"
    echo "  gh release create $version dist/* --title \"Release $version\" --notes-file release_notes.md"
}

# Main execution
if [[ $# -eq 0 ]] || [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
    show_help
    exit 0
fi

VERSION=$1

echo -e "${BLUE}Starting release process for $VERSION${NC}"
echo ""

# Pre-flight checks
echo "Running pre-flight checks..."
validate_version "$VERSION"
check_git_status
check_branch
check_tag_exists "$VERSION"

# Run tests
run_tests

# Update changelog
update_changelog "$VERSION"

# Confirm before proceeding
echo ""
echo -e "${YELLOW}Ready to create release $VERSION${NC}"
read -p "Continue? (y/N) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Release cancelled"
    exit 0
fi

# Create tag
create_tag "$VERSION"

# Build release
build_release "$VERSION"

# Show next steps
show_next_steps "$VERSION"
