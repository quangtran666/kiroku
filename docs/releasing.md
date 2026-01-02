# Release Guide

This document documents the release process for Kiroku.

## Prerequisites

- Access to the GitHub repository.
- Git installed locally.

## Release Process

The release process is fully automated using **GitHub Actions** and **GoReleaser**. You do not need to build binaries manually.

### 1. Update Version and Documentation

Before releasing, ensure all changes are committed and documentation is up to date.

```bash
# Check status
git status

# Commit any pending changes
git add .
git commit -m "chore: prepare for release v1.0.1"
git push
```

### 2. Trigger a Release

To trigger a release, simply create and push a semantic version tag (e.g., `v1.0.0`, `v1.1.0`, `v2.0.0`).

```bash
# 1. Create a tag
git tag v1.0.1

# 2. Push the tag to GitHub
git push origin v1.0.1
```

### 3. Verification

Once the tag is pushed:

1.  Go to the **Actions** tab in the GitHub repository.
2.  You will see a workflow named `release` running.
3.  Wait for it to complete (usually 2-3 minutes).

### 4. Result

GitHub Actions will automatically:

1.  Build binaries for:
    - Linux (AMD64, ARM64)
    - macOS (Intel, Apple Silicon)
    - Windows (AMD64, ARM64)
2.  Create Linux packages (`.deb`, `.rpm`).
3.  Draft a new Release on GitHub.
4.  Upload all artifacts (`.tar.gz`, `.deb`, `.rpm`, checksums) to the release page.

### 5. Accessing the Release

Users can download the latest version from the [Releases Page](https://github.com/quangtran666/kiroku/releases).

**Installation for Ubuntu/Debian users:**
Download the `.deb` file and run:

```bash
sudo dpkg -i kiroku_1.0.1_linux_amd64.deb
```
