# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Which versions are eligible for receiving such patches depends on the CVSS v3.0 Rating:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

**Please do not report security vulnerabilities through public GitHub issues.**

### Where to Report

If you discover a security vulnerability, please send an email to:

**security@fleeks.dev**

### What to Include

To help us triage and respond quickly, please include:

- Type of vulnerability (e.g., buffer overflow, SQL injection, cross-site scripting, etc.)
- Full paths of source file(s) related to the vulnerability
- Location of the affected source code (tag/branch/commit or direct URL)
- Step-by-step instructions to reproduce the issue
- Proof-of-concept or exploit code (if possible)
- Impact of the issue, including how an attacker might exploit it
- Any special configuration required to reproduce the issue

### What to Expect

- **Acknowledgment**: We will acknowledge receipt of your vulnerability report within 48 hours
- **Communication**: We will keep you informed about the progress of addressing the vulnerability
- **Timeline**: We aim to release a fix within:
  - **Critical vulnerabilities**: 7 days
  - **High vulnerabilities**: 14 days
  - **Medium vulnerabilities**: 30 days
  - **Low vulnerabilities**: 90 days
- **Credit**: We will credit you in the security advisory (unless you prefer to remain anonymous)

## Security Update Process

1. **Triage**: We assess the vulnerability and determine its severity
2. **Fix Development**: We develop a fix in a private repository
3. **Testing**: We thoroughly test the fix
4. **Release**: We release a new version with the security fix
5. **Disclosure**: We publish a security advisory detailing the vulnerability
6. **Notification**: We notify users through:
   - GitHub Security Advisories
   - Release notes
   - Email to registered users (if applicable)

## Security Best Practices

When using Fleeks CLI, we recommend:

### For Users

- **Keep Updated**: Always use the latest version of Fleeks CLI
- **Secure Credentials**: Never share your API keys or credentials
- **Environment Variables**: Use environment variables for sensitive data
- **Configuration Files**: Protect your `~/.fleeksconfig.yaml` file (chmod 600)
- **HTTPS Only**: Always use HTTPS endpoints
- **Review Permissions**: Regularly review workspace and container permissions

### For Developers

- **Code Review**: All code changes require review before merging
- **Dependency Scanning**: We scan dependencies for known vulnerabilities
- **Input Validation**: Always validate and sanitize user input
- **Secrets Management**: Never commit secrets to the repository
- **Secure Defaults**: Use secure defaults in all configurations
- **Least Privilege**: Request only necessary permissions

## Known Security Considerations

### API Keys and Tokens

- API keys are stored in `~/.fleeksconfig.yaml`
- This file should have restrictive permissions (600 on Unix-like systems)
- Tokens are transmitted over HTTPS only
- Tokens expire after a configurable period

### Container Security

- Containers run with least privilege by default
- Network isolation is enabled by default
- Volume mounts are restricted to user directories
- Container images are verified before use

### File Operations

- File operations respect system permissions
- Path traversal protection is implemented
- Symlinks are carefully handled to prevent attacks

### Network Security

- All API communication uses TLS 1.2 or higher
- Certificate validation is enabled by default
- WebSocket connections use WSS (secure WebSocket)

## Security Advisories

We publish security advisories through:

- [GitHub Security Advisories](https://github.com/fleeks-inc/fleeks-cli/security/advisories)
- Our website: https://fleeks.dev/security
- Email notifications to registered users

## Vulnerability Disclosure Policy

We follow a coordinated vulnerability disclosure process:

1. **Report Received**: Vulnerability reported to security@fleeks.dev
2. **Acknowledgment**: Within 48 hours
3. **Investigation**: 1-7 days
4. **Fix Development**: Varies by severity
5. **Private Testing**: With reporter if desired
6. **Public Release**: Security fix released
7. **Public Disclosure**: Advisory published (90 days maximum from report)

## Bug Bounty Program

We do not currently have a bug bounty program, but we:

- Acknowledge all valid security reports
- Credit researchers in advisories (with permission)
- Provide public recognition on our website
- May offer rewards on a case-by-case basis for exceptional findings

## Security Team

Our security team can be reached at:

- **Email**: security@fleeks.dev
- **PGP Key**: Available at https://fleeks.dev/security/pgp-key.txt
- **Response Time**: Within 48 hours

## Questions?

If you have questions about this security policy, please contact:

- **General Security Questions**: security@fleeks.dev
- **Policy Questions**: legal@fleeks.dev

---

**Last Updated**: November 11, 2025

**Version**: 1.0
