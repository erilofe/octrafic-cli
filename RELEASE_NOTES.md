# Release Notes v0.1.0

Initial release of Octrafic CLI - AI-powered API testing tool.

## Features

- Natural language API testing with AI
- Multi-provider support: Claude (Anthropic), OpenRouter, OpenAI
- Format support: OpenAPI/Swagger, Postman Collections, GraphQL, Markdown
- Auto-conversion of various API formats to OpenAPI
- Interactive TUI with project management
- Authentication support: Bearer, API Key, Basic Auth
- Cross-platform: Linux, macOS, Windows (amd64, arm64, armv7)

## Installation

**Quick install:**
```bash
# Linux & macOS
curl -fsSL https://octrafic.com/install.sh | bash

# Windows
iex (iwr -useb https://octrafic.com/install.ps1)

# macOS (Homebrew)
brew install octrafic/tap/octrafic

# Arch Linux
yay -S octrafic-bin
```

**Package managers:**
- DEB packages for Debian/Ubuntu
- RPM packages for Fedora/RHEL/CentOS
- Homebrew for macOS
- AUR for Arch Linux

## What's New

Initial release with core functionality:
- AI-powered test generation and execution
- Interactive project wizard with auth configuration
- Multi-architecture support
- Automated release pipeline with GoReleaser
