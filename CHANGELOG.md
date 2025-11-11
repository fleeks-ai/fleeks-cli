# Changelog

All notable changes to Fleeks CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial open source release
- Complete documentation suite
- GitHub Actions CI/CD pipeline
- Community contribution templates

## [1.0.0] - 2025-11-11

### Added

#### Core Features
- 🎨 Beautiful gradient ASCII logo with smooth color transitions
- 🔐 Complete authentication system (login, logout, register, refresh)
- 🤖 AI agent management with real-time streaming
- 📦 Workspace management (create, list, delete, file operations)
- 🐳 Container management (create, start, stop, logs)
- 📁 File operations (create, update, delete, list)
- 💻 Terminal operations with command execution
- 🌍 Environment management (development, staging, production)

#### Configuration
- Environment-based configuration system
- User config stored in `~/.fleeksconfig.yaml`
- Support for `.env.development`, `.env.staging`, `.env.production`
- Automatic config migration from deprecated fields

#### Documentation
- Comprehensive README with installation and usage instructions
- BUILD_AND_TEST.md for developers
- BACKEND_INTEGRATION.md for backend team (25+ API endpoints)
- INFRASTRUCTURE_DEPLOYMENT.md for DevOps team
- PRODUCTION_READINESS.md with launch checklist
- CONTRIBUTING.md for contributors
- CODE_OF_CONDUCT.md for community guidelines
- SECURITY.md for security policies

#### Developer Experience
- Cobra CLI framework for robust command structure
- Viper for configuration management
- Color support with gookit/color and fatih/color
- Error handling and validation
- Progress indicators and status messages

### Technical Details

#### Commands

**Authentication**
```bash
fleeks auth login              # Log in to Fleeks
fleeks auth logout             # Log out from Fleeks
fleeks auth register           # Register new account
fleeks auth refresh            # Refresh authentication token
fleeks auth whoami             # Display current user info
```

**Workspaces**
```bash
fleeks workspace create        # Create new workspace
fleeks workspace list          # List all workspaces
fleeks workspace delete        # Delete workspace
fleeks workspace get           # Get workspace details
```

**Agents**
```bash
fleeks agent start             # Start AI agent
fleeks agent stop              # Stop AI agent
fleeks agent status            # Get agent status
fleeks agent stream            # Stream agent output
```

**Containers**
```bash
fleeks container create        # Create container
fleeks container start         # Start container
fleeks container stop          # Stop container
fleeks container logs          # View container logs
fleeks container list          # List containers
```

**Files**
```bash
fleeks files create            # Create file
fleeks files update            # Update file
fleeks files delete            # Delete file
fleeks files get               # Get file content
fleeks files list              # List files
```

**Terminal**
```bash
fleeks terminal run            # Execute terminal command
fleeks terminal start          # Start terminal session
fleeks terminal stop           # Stop terminal session
```

**Environment**
```bash
fleeks env list                # List environments
fleeks env set                 # Set environment variables
fleeks env get                 # Get environment variable
fleeks env delete              # Delete environment variable
```

#### API Integration
- REST API client with authentication
- WebSocket support for real-time streaming
- Automatic token refresh
- Request/response logging
- Error handling and retries

#### Build System
- Multi-platform support (Windows, macOS, Linux)
- Architecture support (amd64, arm64)
- Version injection at build time
- Automated releases via GitHub Actions

### Infrastructure
- Docker support for local development
- Environment-based API endpoints
- TLS/SSL support
- Configuration validation

### Security
- Secure credential storage
- HTTPS-only API communication
- Token-based authentication
- Environment variable support for secrets

## Platform Support

### Officially Supported
- ✅ Windows 10/11 (amd64)
- ✅ macOS 12+ Intel (amd64)
- ✅ macOS 12+ Apple Silicon (arm64)
- ✅ Linux (amd64)
- ✅ Linux (arm64)

### Requirements
- Go 1.21 or higher (for building from source)
- Internet connection (for API access)
- Docker (optional, for container features)

## Known Issues

None at release.

## Upgrade Instructions

### From Pre-1.0 Versions

This is the first stable release. If you were using a pre-release version:

1. Download the latest binary for your platform
2. Replace your old binary
3. Run `fleeks auth login` to re-authenticate
4. Your configuration will be automatically migrated

## Contributors

Special thanks to all contributors who made this release possible!

- [@fleeks-team](https://github.com/fleeks-team) - Core development

## Links

- [Documentation](https://github.com/fleeks-inc/fleeks-cli/blob/main/README.md)
- [Report Issues](https://github.com/fleeks-inc/fleeks-cli/issues)
- [Discussions](https://github.com/fleeks-inc/fleeks-cli/discussions)

---

[Unreleased]: https://github.com/fleeks-inc/fleeks-cli/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/fleeks-inc/fleeks-cli/releases/tag/v1.0.0
