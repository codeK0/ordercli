---
summary: 'Release checklist for ordercli (GitHub release + Homebrew tap)'
---

# Releasing ordercli

Follow these steps for each release. Title GitHub releases as `ordercli <version>`.

## Checklist
- Update version in `internal/version/version.go`.
- Update `CHANGELOG.md` with the new version section.
- Tag the release: `git tag -a v<version> -m "Release <version>"` and push tags after commits.
- Build source archive for Homebrew: `git archive --format=tar.gz --output /tmp/ordercli-<version>.tar.gz v<version>`.
- Compute checksum: `shasum -a 256 /tmp/ordercli-<version>.tar.gz`.
- Update Homebrew tap (`../homebrew-tap/Formula/ordercli.rb`):
  - Set `version "X.Y.Z"`.
  - Set `url "https://github.com/steipete/ordercli/archive/refs/tags/vX.Y.Z.tar.gz"`.
  - Paste `sha256` from the archive.
- Commit and push changes in ordercli and the tap; push tags: `git push origin main --tags` then `git push` in `../homebrew-tap`.
- Create GitHub release for `v<version>`:
  - Title: `ordercli <version>`
  - Body: bullets from `CHANGELOG.md` for that version.
  - Assets: attach `/tmp/ordercli-<version>.tar.gz` with SHA256 sum in the body.
