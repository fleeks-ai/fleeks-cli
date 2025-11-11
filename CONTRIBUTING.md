# Contributing to Fleeks CLI

First off, thank you for considering contributing to Fleeks CLI! 🎉 It's people like you that make Fleeks CLI such a great tool.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [How to Contribute](#how-to-contribute)
- [Development Setup](#development-setup)
- [Submitting Changes](#submitting-changes)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Documentation](#documentation)
- [Community](#community)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to conduct@fleeks.dev.

## Getting Started

### Types of Contributions We're Looking For

We welcome many types of contributions:

- 🐛 **Bug fixes** - Found a bug? Submit a fix!
- ✨ **New features** - Have an idea? Let's discuss it!
- 📝 **Documentation** - Improvements to docs are always welcome
- 🎨 **UI/UX improvements** - Make the CLI more beautiful
- 🧪 **Tests** - More tests = more stability
- 🌍 **Translations** - Help make Fleeks accessible globally
- 💡 **Ideas** - Share your thoughts in discussions

### What We're NOT Looking For

- Breaking changes without prior discussion
- Features that don't align with project goals
- Malicious code or security vulnerabilities
- Plagiarized code

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates.

When you create a bug report, include as many details as possible:

**Use the bug report template:**
- Clear, descriptive title
- Steps to reproduce
- Expected vs actual behavior
- CLI version (`fleeks version`)
- OS and version
- Relevant logs/screenshots

### Suggesting Features

We love feature suggestions! Before suggesting a feature:

1. Check if it's already been suggested
2. Make sure it aligns with project goals
3. Be clear and detailed about the use case

**Use the feature request template:**
- Problem you're trying to solve
- Proposed solution
- Alternative solutions considered
- Additional context

### Pull Requests

1. Fork the repo and create your branch from `main`
2. Make your changes
3. Add tests if applicable
4. Update documentation
5. Ensure tests pass
6. Submit a pull request

## Development Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Docker (optional, for backend testing)

### Setup Steps

```bash
# 1. Fork and clone the repository
git clone https://github.com/YOUR_USERNAME/fleeks-cli.git
cd fleeks-cli

# 2. Add upstream remote
git remote add upstream https://github.com/fleeks-inc/fleeks-cli.git

# 3. Install dependencies
go mod download

# 4. Build the CLI
go build -o fleeks main.go

# 5. Run tests
go test ./...

# 6. Test the CLI
./fleeks --version
```

### Project Structure

```
fleeks-cli/
├── cmd/              # Command definitions
├── internal/         # Internal packages
│   ├── client/      # API client
│   └── config/      # Configuration
├── .github/         # GitHub workflows
├── docs/            # Documentation
└── main.go          # Entry point
```

## Submitting Changes

### Branch Naming

Use descriptive branch names:
- `feature/add-workspace-templates`
- `fix/authentication-bug`
- `docs/update-readme`
- `test/add-workspace-tests`

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Formatting, missing semicolons, etc.
- `refactor`: Code restructuring
- `test`: Adding tests
- `chore`: Maintenance tasks

**Examples:**
```bash
feat(workspace): add template support for Rust projects

fix(auth): resolve token expiration issue

docs(readme): update installation instructions

test(agent): add unit tests for agent controller
```

### Pull Request Process

1. **Update your branch:**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Run tests:**
   ```bash
   go test ./...
   go vet ./...
   ```

3. **Push changes:**
   ```bash
   git push origin your-branch-name
   ```

4. **Create PR:**
   - Use a clear, descriptive title
   - Fill out the PR template completely
   - Link related issues
   - Add screenshots/GIFs for UI changes
   - Request review from maintainers

5. **Address feedback:**
   - Respond to all review comments
   - Make requested changes
   - Re-request review when ready

6. **Merge:**
   - Maintainers will merge once approved
   - Squash commits if requested

## Coding Standards

### Go Style Guide

Follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments):

- Use `gofmt` for formatting
- Use meaningful variable names
- Keep functions small and focused
- Add comments for exported functions
- Handle errors explicitly
- Use context for cancelation

### Code Example

```go
// CreateWorkspace creates a new workspace with the given name and template
func CreateWorkspace(ctx context.Context, name, template string) (*Workspace, error) {
    if name == "" {
        return nil, fmt.Errorf("workspace name cannot be empty")
    }

    // Create workspace request
    req := &CreateWorkspaceRequest{
        Name:     name,
        Template: template,
    }

    // Send request to API
    resp, err := client.Post(ctx, "/workspaces", req)
    if err != nil {
        return nil, fmt.Errorf("failed to create workspace: %w", err)
    }

    return resp.Workspace, nil
}
```

### File Organization

- One command per file in `cmd/`
- Group related functionality
- Keep files under 500 lines
- Use internal packages for shared code

## Testing Guidelines

### Writing Tests

```go
func TestCreateWorkspace(t *testing.T) {
    tests := []struct {
        name        string
        workspaceName string
        template    string
        wantErr     bool
    }{
        {
            name:        "valid workspace",
            workspaceName: "my-project",
            template:    "python",
            wantErr:     false,
        },
        {
            name:        "empty name",
            workspaceName: "",
            template:    "python",
            wantErr:     true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := CreateWorkspace(context.Background(), tt.workspaceName, tt.template)
            if (err != nil) != tt.wantErr {
                t.Errorf("CreateWorkspace() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -run TestCreateWorkspace ./cmd

# Run with verbose output
go test -v ./...
```

### Test Coverage

- Aim for >80% coverage for new code
- Test edge cases and error conditions
- Mock external dependencies
- Use table-driven tests

## Documentation

### Code Documentation

- Add GoDoc comments to all exported functions
- Include examples in documentation
- Document parameters and return values
- Explain complex logic with inline comments

### User Documentation

When adding features:
- Update README.md
- Add command examples
- Update relevant guides
- Add troubleshooting tips

### Documentation Style

- Use clear, concise language
- Include code examples
- Add screenshots/GIFs when helpful
- Keep it up-to-date

## Community

### Where to Get Help

- 💬 **Discord:** https://discord.gg/fleeks
- 📧 **Email:** support@fleeks.dev
- 💭 **Discussions:** GitHub Discussions tab
- 🐛 **Issues:** GitHub Issues tab

### Communication Guidelines

- Be respectful and inclusive
- Stay on topic
- Help others when you can
- Follow the Code of Conduct

### Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Credited in relevant documentation
- Given GitHub contributor badge

## Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):
- `MAJOR.MINOR.PATCH`
- Major: Breaking changes
- Minor: New features
- Patch: Bug fixes

### Release Schedule

- **Patch releases:** As needed for critical bugs
- **Minor releases:** Monthly
- **Major releases:** Quarterly

## Questions?

Don't hesitate to ask! You can:
- Open a discussion on GitHub
- Join our Discord server
- Email us at dev@fleeks.ai

## Thank You! 🙏

Your contributions make Fleeks CLI better for everyone. We appreciate your time and effort!

---

**Happy Contributing! 🚀**

*Last updated: November 11, 2025*
