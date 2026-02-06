# Release Process

## Quick Release

1. Create and push a new tag:
```bash
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

2. GitHub Actions will automatically:
   - Run tests
   - Build binaries for all platforms
   - Create GitHub Release with:
     - Binary archives (.tar.gz, .zip)
     - DEB packages (Debian/Ubuntu)
     - RPM packages (Fedora/RHEL)
     - Checksums
     - Install scripts (install.sh, install.ps1)
     - Changelog
   - Update Homebrew tap (octrafic/homebrew-tap)
   - Update AUR package (if configured)

## First-Time Setup

### Required Repositories

1. **Homebrew Tap** (required for Homebrew installation):
   ```bash
   # Create on GitHub:
   # Repository: octrafic/homebrew-tap
   # Description: Homebrew formulae for Octrafic
   # Public repository
   ```

2. **AUR Package** (optional, for Arch Linux):
   - Create account on https://aur.archlinux.org
   - Generate SSH key for AUR
   - Add to GitHub Secrets as `AUR_KEY`

### GitHub Secrets

Required secrets (Settings → Secrets and variables → Actions):

- `GITHUB_TOKEN` - Automatically provided by GitHub ✓
- `AUR_KEY` - SSH private key for AUR (optional)

## Release Checklist

Before creating a release:

- [ ] Update version numbers if needed
- [ ] Update CHANGELOG.md
- [ ] Run tests locally: `go test ./...`
- [ ] Test local build: `goreleaser build --snapshot --clean`
- [ ] Commit all changes
- [ ] Create and push tag

## Tag Naming

- Format: `vMAJOR.MINOR.PATCH`
- Examples: `v1.0.0`, `v1.2.3`, `v2.0.0-beta.1`

## Testing Release Locally

Test release process without publishing:

```bash
goreleaser release --snapshot --clean --skip=publish
```

Check generated files in `dist/`:
- Binaries: `dist/octrafic_*/`
- Archives: `dist/*.tar.gz`, `dist/*.zip`
- Packages: `dist/*.deb`, `dist/*.rpm`
- Homebrew: `dist/homebrew/octrafic.rb`
- AUR: `dist/aur/octrafic-bin.pkgbuild`

## Rollback

If release fails or has issues:

1. Delete the tag locally and remotely:
```bash
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0
```

2. Delete the GitHub Release (if created)

3. Fix issues and create new tag

## Post-Release

After successful release:

1. Verify installation methods:
   ```bash
   # Homebrew (macOS/Linux)
   brew install octrafic/tap/octrafic

   # Debian/Ubuntu
   wget <deb-url>
   sudo apt install ./octrafic_VERSION_linux_amd64.deb

   # Install script
   curl -fsSL https://octrafic.com/install.sh | bash
   ```

2. Update documentation if needed

3. Announce release (social media, etc.)

## Troubleshooting

### Homebrew tap not updating

- Check if `octrafic/homebrew-tap` repository exists
- Verify `GITHUB_TOKEN` has write permissions
- Check GitHub Actions logs

### AUR not updating

- Verify `AUR_KEY` secret is set correctly
- Check if AUR package name is reserved
- Verify SSH key permissions on AUR

### Release workflow fails

1. Check GitHub Actions logs
2. Test locally: `goreleaser release --snapshot --clean --skip=publish`
3. Verify all required files exist (LICENSE, README.md, etc.)
4. Check go.mod and dependencies
