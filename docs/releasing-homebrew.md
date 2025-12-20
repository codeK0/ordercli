# ordercli Homebrew Release Playbook

Automates the Homebrew tap update for new releases.

## 0) Prereqs
- Clean git tree on `main`.
- Tagged release pushed (e.g., `v0.1.0`).
- Tap repo at `../homebrew-tap`.

## 1) Generate formula fields
Run:
```sh
scripts/release-homebrew.sh 0.1.0
```
Copy the printed `version`, `url`, and `sha256`.

## 2) Update the tap
Edit `../homebrew-tap/Formula/ordercli.rb`:
- Set `version "X.Y.Z"`.
- Set `url "https://github.com/steipete/ordercli/archive/refs/tags/vX.Y.Z.tar.gz"`.
- Paste `sha256`.

Commit + push in the tap repo:
```sh
git -C ../homebrew-tap commit -am "ordercli vX.Y.Z"
git -C ../homebrew-tap push origin main
```

## 3) Sanity-check install
```sh
brew uninstall ordercli || true
brew untap steipete/tap || true
brew tap steipete/tap
brew install steipete/tap/ordercli
brew test steipete/tap/ordercli
ordercli --help
```
