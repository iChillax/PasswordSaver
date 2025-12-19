# PasswordSaver

Hi everyone, my name is Derek, and I'm a DevOps Engineer. I met some troubles and issues when I was looking for a secure place to save my accounts and passwords. So that is the main reason I made this project, so everyone can host themselves password management portal.

## 🚀 Quick Start

Get started in 5 minutes: [QUICKSTART.md](QUICKSTART.md)

## 📚 Documentation

### Getting Started
- [Quick Start Guide](QUICKSTART.md) - Setup in 5 minutes
- [Troubleshooting](TROUBLESHOOTING.md) - Common issues and solutions

### Features
- [Ollama Setup](docs/OLLAMA_SETUP.md) - AI-powered chat setup
- [Tag Colors](docs/features/tag-colors.md) - Unique colors for tags and categories
- [Secret Notes](docs/features/secret-notes.md) - Add descriptions to secrets
- [Account Type](docs/features/account-secret-type.md) - Store username + password pairs
- [Modal Expansion](docs/features/modal-expansion.md) - Smooth secret detail view
- [Phase 1 Tasks](docs/tasks/phase-1-tasks.md) - Core features
- [Phase 2 Tasks](docs/tasks/phase-2-tasks.md) - AI chatbot integration

### Development
- [Release Process](docs/RELEASE_PROCESS.md) - How to create releases
- [Release Workflow](.github/RELEASE_QUICKSTART.md) - Quick release guide
- [GitHub Workflows](.github/workflows/README.md) - CI/CD documentation

### Project Planning
- [Full Plan](plan.md) - Complete project roadmap
- [Branching Convention](docs/BRANCHING_CONVENTION.md) - Git workflow

## ✨ Features

### Core Security
- 🔐 Secure password storage with AES-256 encryption
- 🔑 JWT-based authentication
- 🤖 AI-powered natural language search (Ollama)

### Secret Management
- 🏷️ Organize secrets with categories and tags
- 🎨 Unique colors for each tag and category (automatic assignment)
- 📝 Add notes and descriptions to secrets
- 👤 Account type for storing username + password pairs
- 🔍 Advanced filtering by category and tags

### User Experience
- 📱 Responsive web interface with modal expansion
- ✨ Smooth animations and transitions
- 👁️ Show/hide toggle for password visibility
- 📋 Copy-to-clipboard functionality
- 🏷️ Inline tag editing with color-coded tags
- 🐳 Docker-based deployment

## 🛠️ Tech Stack

**Backend:**
- Go (Gin framework)
- MongoDB
- AES-256 encryption
- JWT authentication

**Frontend:**
- Vue.js 3
- Vite
- Tailwind CSS
- Pinia (state management)

**AI/ML:**
- Ollama (local LLM)
- Mistral/LLaMA models

**Infrastructure:**
- Docker & Docker Compose
- MongoDB
- Ollama service

## 📦 Installation

See [QUICKSTART.md](QUICKSTART.md) for detailed setup instructions.

```bash
# 1. Clone repository
git clone https://github.com/yourusername/PasswordSaver.git
cd PasswordSaver

# 2. Setup environment
cp backend/.env.example backend/.env
# Edit backend/.env with your keys

# 3. Start services
cd dockerfiles
docker-compose up -d

# 4. Setup Ollama (optional)
./setup-ollama.sh
```

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:

1. Label your PRs appropriately (see [Release Process](docs/RELEASE_PROCESS.md))
2. Follow the [branching convention](docs/BRANCHING_CONVENTION.md)
3. Write tests for new features
4. Update documentation

## 📝 License

See [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Built with ❤️ by Derek and contributors.
