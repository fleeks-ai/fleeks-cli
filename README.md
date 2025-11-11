# 🚀 Fleeks CLI - Revolutionary AI-Powered Development Platform

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org/dl/)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey.svg)](https://github.com/fleeks-inc/fleeks-cli/releases)

```
▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄▄
██ ███████ ██      ██████████ ██████ ██    ██  ██████ ██
██ ██      ██      ██      ██ ██     ██  ██   ██      ██
██ ███████ ██      ████████   ██████ ████     ██████  ██
██ ██      ██      ██         ██     ██  ██         ██ ██
██ ██      ███████ ██         ██████ ██    ██ ██████  ██
▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀▀
```

## 🌟 The World's First Universal AI Software Engineer CLI

Fleeks CLI revolutionizes software development with capabilities **no competitor has**:

✅ **Single universal AI software engineer** (handles ALL project types and tasks)  
✅ **Dynamic expertise adaptation** (web, mobile, blockchain, games, AI/ML, IoT)  
✅ **Hybrid local-cloud workspace management**  
✅ **Persistent project memory across sessions**  
✅ **Real-time streaming collaboration**  
✅ **Full container orchestration integrated with AI**

---

## 🚀 Quick Start

### Installation

```bash
# Download the latest release
curl -sSL https://releases.fleeks.dev/cli/install.sh | bash

# Or build from source
git clone https://github.com/fleeks-inc/fleeks-cli.git
cd fleeks-cli
go build -o fleeks .
```

### Initial Setup

```bash
# 1. Authenticate with Fleeks
fleeks auth login

# 2. Create your first workspace
fleeks workspace create my-api --template microservices

# 3. Start your AI software engineer
fleeks agent start --task "Design and implement user authentication service"

# 4. Watch the magic happen!
fleeks agent watch my-api
```

---

## 🏗️ Environment Configuration

Fleeks CLI supports three environments for seamless development-to-production workflows:

### Development Environment
```bash
# Use local backend services
fleeks --environment development workspace create my-app

# Automatically connects to:
# - Main API: http://localhost:8000
# - LSP Service: http://localhost:8001  
# - MCP Service: http://localhost:8002
# - Neo4j: bolt://localhost:7687
# - Qdrant: http://localhost:6333
```

### Staging Environment
```bash
# Use staging services for testing
fleeks --environment staging agent start --task "Implement user authentication"
```

### Production Environment
```bash
# Use production services
fleeks --environment production workspace list
```

### Environment Management
```bash
# View current environment
fleeks env info

# Test connectivity
fleeks env test

# List all settings
fleeks env list
```

---

## 📚 Core Commands

### 🏗️ Workspace Management
```bash
# Create workspace with template
fleeks workspace create my-project --template python
fleeks workspace create api-service --template microservices
fleeks workspace create frontend --template react

# List workspaces
fleeks workspace list

# Get workspace info
fleeks workspace info my-project

# Sync local changes to cloud
fleeks workspace sync my-project --watch

# Delete workspace
fleeks workspace delete my-project
```

### 🤖 AI Software Engineer
```bash
# Start the universal AI software engineer
fleeks agent start --task "Design and implement user authentication"
fleeks agent start --task "Build React Native mobile app"
fleeks agent start --task "Create Solidity smart contracts"
fleeks agent start --task "Setup CI/CD pipeline with testing"

# Agent automatically adapts expertise based on task:
# - "React Native" → Mobile development skills loaded
# - "Solidity" → Blockchain development skills loaded
# - "CI/CD" → DevOps skills loaded

# Monitor agent progress
fleeks agent watch my-project
fleeks agent status my-project

# Chat with your software engineer
fleeks chat my-project

# Stop agent
fleeks agent stop my-project
```

### 🐳 Container Orchestration
```bash
# Get container info
fleeks container info my-project

# View resource usage
fleeks container stats my-project

# Execute commands in container
fleeks container exec my-project "npm test"
fleeks container exec my-project "python manage.py migrate"

# View logs
fleeks container logs my-project --follow

# Scale resources
fleeks container scale my-project --cpu 4 --memory 8Gi
```

### 📁 Smart File Operations
```bash
# List files in workspace
fleeks files list my-project --recursive

# Upload files with smart sync
fleeks files upload my-project ./src /workspace/src --recursive
fleeks files upload my-project ./package.json /workspace/package.json

# Download files
fleeks files download my-project /workspace/dist ./dist --recursive

# Create files remotely
fleeks files create my-project /workspace/README.md "# My Project"

# Watch for file changes
fleeks files watch my-project

# Delete files
fleeks files delete my-project /workspace/old-file.txt
```

### 🖥️ Terminal Operations
```bash
# Execute commands
fleeks terminal exec my-project "npm run build"
fleeks terminal exec my-project "python manage.py test"

# Run background jobs
fleeks terminal run my-project "python server.py" --background
fleeks terminal run my-project "npm run dev" --background

# Manage background jobs
fleeks terminal jobs my-project
fleeks terminal stop my-project job-123
fleeks terminal logs my-project job-123

# Interactive shell (coming soon)
fleeks terminal shell my-project
```

### 🔐 Authentication
```bash
# Login to Fleeks
fleeks auth login

# Check authentication status
fleeks auth status

# View user information
fleeks auth whoami

# Logout
fleeks auth logout
```

---

## 🔥 Revolutionary Features

### 1. Universal AI Software Engineer with Dynamic Expertise
Unlike any competitor, Fleeks provides **one intelligent agent that adapts to ANY project type**:

```bash
# Single agent handles complete workflow
fleeks agent start --task "Build a DeFi staking protocol with mobile app"

# Agent automatically detects project types and loads expertise:
# ✅ Detects "DeFi" + "staking" → Blockchain skills loaded
# ✅ Detects "mobile app" → Mobile development skills loaded
# ✅ Dynamically switches between Solidity, React Native, testing, deployment

# Agent works on multiple project types simultaneously:
fleeks chat my-project
> "Create smart contracts for token staking"
# → Blockchain expertise active

> "Now build React Native app to interact with contracts"  
# → Mobile + Blockchain expertise both active

> "Add Unity game integration for rewards"
# → Game + Blockchain + Mobile expertise all active
```

**No context switching required** - the agent is a true polyglot software engineer!

### 2. Hybrid Local-Cloud Workspaces
**Revolutionary workspace management** that competitors can't match:

```bash
# Create local workspace, sync to cloud instantly
fleeks workspace create my-app --template python
# ✅ Local files created
# ✅ Cloud container ready in <100ms
# ✅ Smart sync active

# Real-time bidirectional sync
fleeks workspace sync my-app --watch --bidirectional
# ✅ Local changes → Cloud instantly
# ✅ Agent changes → Local instantly
# ✅ Conflict resolution built-in
```

### 3. Real-Time Streaming Collaboration
**Live monitoring** of AI software engineer activities:

```bash
fleeks agent watch my-project
# 🔴 LIVE: Analyzing requirements for DeFi protocol...
# 🔴 LIVE: [BLOCKCHAIN] Writing Solidity smart contracts...
# 🔴 LIVE: [BLOCKCHAIN] Implementing staking logic with rewards...
# 🔴 LIVE: [MOBILE] Creating React Native app structure...
# 🔴 LIVE: [MOBILE] Building Web3 wallet integration...
# 🔴 LIVE: [TESTING] Writing contract tests with Hardhat...
# 🔴 LIVE: [DEVOPS] Configuring deployment pipeline...
```

### 4. Container Orchestration + AI
**AI-integrated container management** no competitor offers:

```bash
# AI automatically manages resources based on workload
fleeks container stats my-project
# CPU: 45% (AI will auto-scale at 80%)
# Memory: 2.1Gi/4Gi (AI optimizing allocation)
# Network: 150 Mbps (AI-optimized routing)
```

---

## 🛠️ Advanced Usage

### Environment Variables
```bash
# Override default settings
export FLEEKS_API_BASE_URL="https://custom-api.example.com"
export FLEEKS_WORKSPACE_BASE_PATH="/custom/workspaces"
export FLEEKS_AGENT_MAX_CONCURRENT="5"

# Development mode
export FLEEKS_DEV_MODE=true
export FLEEKS_DEV_VERBOSE=true
```

### Configuration File
Create `~/.fleeksconfig.yaml`:

```yaml
api:
  base_url: "https://api.fleeks.dev"
  timeout: "30s"
  retry_count: 3

workspace:
  default_template: "python"
  sync_enabled: true
  local_path: "./workspaces"

agent:
  max_iterations: 10
  streaming_enabled: true
  preserve_context: true

streaming:
  enabled: true
  buffer_size: 1024
```

### Workspace Templates

Available templates for instant setup:

- **`python`** - Python application with virtual environment
- **`node`** - Node.js application with npm/yarn
- **`go`** - Go application with modules
- **`rust`** - Rust application with Cargo
- **`microservices`** - Multi-service architecture
- **`react`** - React frontend application
- **`vue`** - Vue.js frontend application
- **`django`** - Django web application
- **`fastapi`** - FastAPI microservice
- **`nextjs`** - Next.js full-stack application

---

## 🆚 Competitive Advantage

### vs GitHub Copilot CLI
| Feature | Fleeks CLI | GitHub Copilot CLI |
|---------|------------|-------------------|
| Universal software engineer | ✅ **Handles ALL project types** | ❌ Limited to code completion |
| Dynamic expertise | ✅ **Web, Mobile, Blockchain, Games, AI/ML, IoT** | ❌ Generic coding assistance |
| Hybrid workspaces | ✅ **Local-cloud sync** | ❌ Local only |
| Real-time streaming | ✅ **Live collaboration** | ❌ No streaming |
| Container integration | ✅ **AI-managed containers** | ❌ No container support |
| Multi-project support | ✅ **Switch between types seamlessly** | ❌ Single context only |

### vs Claude Code
| Feature | Fleeks CLI | Claude Code |
|---------|------------|-------------|
| Persistent memory | ✅ **Cross-session context** | ❌ Session-based only |
| Background execution | ✅ **Long-running tasks** | ❌ Interactive only |
| Production deployment | ✅ **Full DevOps pipeline** | ❌ Development only |
| Project type detection | ✅ **Automatic from conversation** | ❌ Manual specification |

### vs Gemini CLI
| Feature | Fleeks CLI | Gemini CLI |
|---------|------------|------------|
| Local-cloud hybrid | ✅ **Revolutionary architecture** | ❌ Cloud-only |
| File synchronization | ✅ **Smart bidirectional sync** | ❌ Manual operations |
| Container orchestration | ✅ **Integrated management** | ❌ No container support |
| Polyglot expertise | ✅ **11+ project types** | ❌ General purpose |

---

## 🧪 Full Development Testing

### Start Backend Services
```bash
# In fleeks-backend-services directory
# Start development environment
npm run start:dev
# or
docker-compose -f docker/docker-compose.yml up

# Verify services are running
fleeks --environment development env test
```

### Test Complete Workflow
```bash
# 1. Create workspace
fleeks --environment development workspace create test-api --template microservices

# 2. Start AI software engineer
fleeks --environment development agent start --task "Design and implement user authentication API with JWT"

# Agent automatically detects:
# - "authentication API" → Web development skills loaded
# - "JWT" → Security best practices loaded

# 3. Watch live progress
fleeks --environment development agent watch test-api

# 4. Upload initial files
fleeks --environment development files upload test-api ./requirements.txt /workspace/requirements.txt

# 5. Execute setup commands
fleeks --environment development terminal exec test-api "pip install -r requirements.txt"

# 6. Run tests
fleeks --environment development terminal run test-api "python -m pytest" --background

# 7. Monitor container resources
fleeks --environment development container stats test-api
```

---

## 📋 Troubleshooting

### Common Issues

#### Authentication Errors
```bash
# Check authentication status
fleeks auth status

# Re-authenticate if needed
fleeks auth logout
fleeks auth login
```

#### Connection Issues
```bash
# Test environment connectivity
fleeks env test

# Check specific environment
fleeks --environment development env test
fleeks --environment staging env test
```

#### Workspace Sync Problems
```bash
# Force sync workspace
fleeks workspace sync my-project --force

# Check workspace status
fleeks workspace info my-project
```

#### Agent Not Responding
```bash
# Check agent status
fleeks agent status my-project

# Restart agent (preserves conversation history)
fleeks agent stop my-project
fleeks agent start --task "Continue previous work"
```

### Debug Mode
```bash
# Enable verbose output
fleeks --verbose workspace create debug-project

# Check environment configuration
fleeks env list

# View detailed logs
fleeks --environment development --verbose agent watch my-project
```

### Environment File Issues
Ensure environment files exist:
- `.env.development` - For local development
- `.env.staging` - For staging environment  
- `.env.production` - For production environment

### Service Dependencies
Required services for full functionality:
- **Main API** (port 8000)
- **LSP Service** (port 8001) 
- **MCP Service** (port 8002)
- **Neo4j** (port 7687)
- **Qdrant** (port 6333)
- **Redis** (port 6379)
- **PostgreSQL** (port 5432)

---

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

## 🌐 Links

- **Website**: [https://fleeks.dev](https://fleeks.dev)
- **Documentation**: [https://docs.fleeks.dev](https://docs.fleeks.dev)
- **API Reference**: [https://api.fleeks.dev/docs](https://api.fleeks.dev/docs)
- **Discord Community**: [https://discord.gg/fleeks](https://discord.gg/fleeks)
- **GitHub**: [https://github.com/fleeks-inc/fleeks-cli](https://github.com/fleeks-inc/fleeks-cli)

---

## 🎯 What's Next?

Fleeks CLI is just the beginning. Coming soon:

- 🧠 **Expanded Expertise** - More domain-specific skills (IoT, gaming, blockchain)
- 🔄 **CI/CD Integration** - GitHub Actions, GitLab CI, Jenkins
- 🎨 **Visual Workspace** - Web-based workspace management
- 📱 **Mobile App** - Monitor your AI engineer on the go
- 🌍 **Multi-cloud** - AWS, GCP, Azure deployment
- 🔐 **Enterprise SSO** - SAML, OIDC, Active Directory
- 🎮 **Enhanced Project Detection** - Even smarter automatic expertise loading

**Join the revolution in AI-powered development!** 🚀