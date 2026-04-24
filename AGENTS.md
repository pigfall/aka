Purpose
This file contains compact, high-signal notes for automated assistants (OpenCode sessions) working on this repository. Only include facts an agent would likely miss.

Quick facts
- Language / toolchain: Go. The module is declared in go.mod (module github.com/pigfall/aka).
- go.mod declares a toolchain: go1.23.6 (check with `go version`). Prefer using Go 1.23.x.
- Entrypoint: CLI implemented under cmd/aka. The binary's main is cmd/aka/main.go (Cobra-based).
- Command implementations: most subcommands live in pkg/cmd. Edit there for behavior changes.
- There is a prebuilt executable named `aka` at repository root. It is a built binary (do not blindly overwrite or execute unknown binaries when making changes — prefer building from source).

Setup (what an agent often misses)
- Ensure Go 1.23.x is installed and on PATH. go.mod specifies toolchain go1.23.6; mismatch may cause tool-specific issues.
- This repo uses mage for the convenience build target; mage is not required but is used by the included Magefile. Install via: `go install github.com/magefile/mage@v1.15.0` (or run the provided install_requirements.sh which contains this line).

Build & run (exact commands)
- Build (preferred reproducible): `go build -x -v -o out/aka ./cmd/aka/`
- Alternative (shortcut): `mage AKA` (requires mage installed). mage runs the go build above and writes out/aka.
- Run from source for quick tests: `go run ./cmd/aka <subcommand> [flags]` or execute the built binary `./out/aka <subcommand> [flags]`.

Tests
- Run all tests: `go test ./...`
- Run package tests only: `go test ./pkg/cmd -v`
- Run a single test: `go test ./pkg/cmd -run TestName -v`
- Tests in pkg/cmd use net/http/httptest (no external network calls). They do not require external services.

Notable commands / caveats to avoid mistakes
- Do NOT run `aka opencode install` (or `go run ./cmd/aka opencode install`) without explicit permission from the user. The opencode installer will:
  - Download a release archive from GitHub (network access).
  - Unpack to $HOME/tools/opencode and set executable bits.
  - Append `export PATH=$PATH:$HOME/tools/opencode` to shell rc files (.bashrc, .zshrc) it finds — this mutates user files.
  - It only supports these platform strings: darwin-arm64, darwin-amd64, linux-arm64, linux-amd64, windows-arm64, windows-amd64. It requires `tar` or `unzip`.

Where to edit
- Add or change CLI behavior: prefer pkg/cmd for subcommand logic and cmd/aka for wiring the Cobra commands.
- Small refactors affecting multiple commands should be done carefully; Cobra wiring is in cmd/aka and individual command code is in pkg/cmd.

Other files and scripts worth knowing
- install_requirements.sh: installs mage (go install github.com/magefile/mage@v1.15.0).
- install_golang.sh / install_golang.ps1: convenience scripts that download a Go toolchain. Note: the shell script downloads Go 1.22.1 — that is older than go.mod's declared toolchain. Prefer installing Go 1.23.x manually or with your system package manager.

Repository hygiene
- The large `aka` binary in the repo root is a built artifact. Do not accidentally modify or remove it unless instructed. When making code changes, rebuild with the source build commands above and update artifacts only if the user asks for a commit.

If something is unclear
- If docs conflict with behavior (for example install scripts vs go.mod toolchain), prefer the executable source of truth (go.mod and build commands in mage.go / cmd/aka) and ask a single clarifying question.
