# Release Process

This document describes the step-by-step process for creating a new release of gx.

## Prerequisites

- [ ] All tests passing
- [ ] All features for the release are merged to main
- [ ] CHANGELOG.md is updated with release notes
- [ ] Version number decided (following semver)

## Release Steps

### 1. Prepare the Release

```bash
# Ensure you're on the main branch
git checkout main
git pull origin main

# Ensure working directory is clean
git status
```

### 2. Update Version Information

Update the CHANGELOG.md with the new version and date:

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- Feature 1
- Feature 2

### Fixed
- Bug fix 1
```

### 3. Run Tests

```bash
# Run all tests
make test

# Or using build scripts
.\build.ps1 test  # Windows
./build.sh test   # Linux/macOS
```

### 4. Create Release Using Helper Script (Recommended)

```bash
# Linux/macOS
./scripts/release.sh v1.0.0

# This will:
# - Validate version format
# - Check git status
# - Run tests
# - Create git tag
# - Build release packages
```

### 5. Manual Release (Alternative)

If not using the helper script:

```bash
# Create and push tag
git tag -a v1.0.0 -m "Release v1.0.0"

# Build release packages
VERSION=v1.0.0 make release  # Linux/macOS
.\build.ps1 release -Version v1.0.0  # Windows
```

### 6. Verify Release Packages

```bash
# Check dist/ directory
ls -lh dist/

# Verify checksums
cat dist/checksums.txt

# Test a binary
tar -xzf dist/gx-v1.0.0-linux-amd64.tar.gz
./gx --version
```

### 7. Push Tag to GitHub

```bash
git push origin v1.0.0
```

This will trigger the GitHub Actions workflow to:
- Run tests on all platforms
- Build release binaries
- Create a GitHub release
- Upload release artifacts

### 8. Create GitHub Release (Manual Alternative)

If not using GitHub Actions:

```bash
# Using GitHub CLI
gh release create v1.0.0 \
  dist/* \
  --title "Release v1.0.0" \
  --notes-file release_notes.md
```

Or manually through GitHub web interface:
1. Go to https://github.com/yourusername/gx/releases/new
2. Select the tag v1.0.0
3. Add release title: "Release v1.0.0"
4. Copy release notes from CHANGELOG.md
5. Upload files from dist/
6. Publish release

### 9. Verify Release

- [ ] GitHub release is created
- [ ] All platform binaries are attached
- [ ] Checksums file is attached
- [ ] Release notes are correct
- [ ] Download and test a binary

### 10. Announce Release

- [ ] Update project README if needed
- [ ] Announce on relevant channels
- [ ] Update documentation site (if applicable)

## Hotfix Release Process

For urgent bug fixes:

1. Create hotfix branch from the release tag:
   ```bash
   git checkout -b hotfix/v1.0.1 v1.0.0
   ```

2. Make the fix and commit

3. Follow steps 2-9 above with the new version

4. Merge hotfix back to main:
   ```bash
   git checkout main
   git merge hotfix/v1.0.1
   git push origin main
   ```

## Pre-release / Beta Release

For pre-releases, use version suffixes:

```bash
# Beta release
./scripts/release.sh v1.0.0-beta.1

# Release candidate
./scripts/release.sh v1.0.0-rc.1
```

Mark as pre-release in GitHub:
- Check "This is a pre-release" when creating the GitHub release

## Rollback

If a release has critical issues:

1. Delete the GitHub release
2. Delete the git tag:
   ```bash
   git tag -d v1.0.0
   git push origin :refs/tags/v1.0.0
   ```
3. Fix the issue
4. Create a new patch release (v1.0.1)

## Troubleshooting

### Build Fails

- Check Go version compatibility
- Ensure all dependencies are available
- Review error messages in build output

### GitHub Actions Fails

- Check workflow logs in GitHub Actions tab
- Verify secrets are configured (if needed)
- Ensure tag format is correct (v*.*.*)

### Release Already Exists

- Delete the existing release and tag
- Or increment version number

## Version Numbering

Follow [Semantic Versioning](https://semver.org/):

- **MAJOR** (X.0.0): Breaking changes
- **MINOR** (0.X.0): New features, backward compatible
- **PATCH** (0.0.X): Bug fixes, backward compatible

Examples:
- v1.0.0 - First stable release
- v1.1.0 - New features added
- v1.1.1 - Bug fixes
- v2.0.0 - Breaking changes
- v1.0.0-beta.1 - Beta release
- v1.0.0-rc.1 - Release candidate
