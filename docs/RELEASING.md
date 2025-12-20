---
summary: 'Release checklist for ordercli (GitHub release)'
---

# Releasing ordercli

Use this checklist for each release (GitHub tag + Homebrew tap).

## Checklist
- Update version in `internal/version/version.go`.
- Update `CHANGELOG.md` with a new version header and date (move “Unreleased”).
- Run tests: `go test ./...`.
- Build binary: `go build ./cmd/ordercli`.
- Tag the release: `git tag -a v<version> -m "Release <version>"`.
- Push commits + tags: `git push origin main --tags`.
- Homebrew tap update (sibling `../homebrew-tap`):
  - Update/create `Formula/ordercli.rb` with `version`, `url`, `sha256`.
  - `url` format: `https://github.com/steipete/ordercli/archive/refs/tags/v<version>.tar.gz`
  - Hash: `curl -L -o /tmp/ordercli.tar.gz <url>` then `shasum -a 256 /tmp/ordercli.tar.gz`.
  - Commit + push tap: `git commit -am "ordercli v<version>" && git push origin main`.
- Create GitHub release `v<version>` using the matching changelog bullets (optional).
