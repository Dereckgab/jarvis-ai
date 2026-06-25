# JARVIS Full IA - Your Personal AI Assistant

> **Enterprise-Grade AI Assistant with Real-Time System Monitoring, Semantic Memory, and Multi-Modal Interactions**

![Status](https://img.shields.io/badge/status-production--ready-brightgreen)
![License](https://img.shields.io/badge/license-MIT-blue)
![Go](https://img.shields.io/badge/Go-1.22-00ADD8)
![Node.js](https://img.shields.io/badge/Node.js-22-339933)
![React](https://img.shields.io/badge/React-19-61DAFB)

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose (latest version)
- 4GB RAM minimum
- 10GB disk space

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/jarvis-fullia.git
   cd jarvis-fullia
   ```

2. **Configure environment variables**
   ```bash
   cp backend/.env.example .env
   # Edit .env with your API keys
   nano .env
   ```

3. **Start the entire stack**
   ```bash
   ./start.sh
   # or using make
   make dev
   ```

4. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - Qdrant Dashboard: http://localhost:6333/dashboard

## 📋 Architecture Overview

### System Architecture (Clean/Hexagonal + Vertical Slice)

```
┌─────────────────────────────────────────────────────────┐
│                    Frontend (Next.js)                    │
│         React 19 + TypeScript + Tailwind + Framer        │
└────────────────────────┬────────────────────────────────┘
                         │
                    HTTP/REST
                         │
┌────────────────────────▼────────────────────────────────┐
│                 Backend (Go + Fiber)                     │
│                                                           │
│  ┌──────────────────────────────────────────────────┐   │
│  │  Interfaces Layer (HTTP Handlers & DTOs)         │   │
│  └──────────────────────────────────────────────────┘   │
│                         │                                 │
│  ┌──────────────────────▼──────────────────────────┐   │
│  │  Application Layer (Services & Use Cases)       │   │
│  │  • AI Service (DeepSeek/OpenAI)                │   │
│  │  • Memory Service (Semantic Search)            │   │
│  │  • System Info Service (Telemetry)             │   │
│  │  • TTS Service (Text-to-Speech)                │   │
│  └──────────────────────────────────────────────────┘   │
│                         │                                 │
│  ┌──────────────────────▼──────────────────────────┐   │
│  │  Domain Layer (Entities & Repositories)         │   │
│  │  • User, SystemInfo, Game, Memory               │   │
│  └──────────────────────────────────────────────────┘   │
│                         │                                 │
│  ┌──────────────────────▼──────────────────────────┐   │
│  │  Infrastructure Layer (Implementations)         │   │
│  │  • MySQL (GORM)                                 │   │
│  │  • Redis (Caching)                              │   │
│  │  • Qdrant (Vector DB)                           │   │
│  │  • External APIs (DeepSeek, OpenAI)             │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
```

### Technology Stack

| Component | Technology | Purpose |
|-----------|-----------|---------|
| **Frontend** | Next.js 16 + React 19 + TypeScript | Modern UI with SSR |
| **Backend** | Go 1.22 + Fiber | High-performance REST API |
| **Database** | MySQL 8.0 | Relational data storage |
| **Cache** | Redis 7.0 | Session & query caching |
| **Vector DB** | Qdrant | Semantic memory & embeddings |
| **AI Engine** | DeepSeek R1/V3 + OpenAI GPT-4o | Cognitive processing |
| **TTS** | Piper/Coqui (Local) | Text-to-speech synthesis |
| **Telemetry** | OpenTelemetry | Distributed tracing |
| **Monitoring** | gopsutil | Cross-platform system metrics |

## 🏗️ Project Structure

```
jarvis-fullia/
├── backend/                          # Go backend
│   ├── cmd/main.go                   # Application entry point
│   ├── config/                       # Configuration management
│   ├── internal/
│   │   ├── domain/                   # Business logic & entities
│   │   │   ├── entity/               # Domain models
│   │   │   └── repository/           # Repository interfaces
│   │   ├── application/              # Use cases & services
│   │   │   └── service/              # Business services
│   │   ├── infrastructure/           # External integrations
│   │   │   ├── repository/           # Repository implementations
│   │   │   ├── database/             # DB connections
│   │   │   ├── ai/                   # AI service implementations
│   │   │   ├── cache/                # Redis client
│   │   │   └── qdrant/               # Qdrant client
│   │   └── interfaces/               # HTTP handlers & middleware
│   │       └── http/                 # REST API handlers
│   ├── pkg/                          # Reusable packages
│   │   ├── errors/                   # Error types
│   │   ├── logger/                   # Logging utility
│   │   ├── security/                 # JWT & security
│   │   ├── telemetry/                # OpenTelemetry setup
│   │   ├── sysmon/                   # System monitoring
│   │   ├── tts/                      # TTS implementation
│   │   └── deepseek-api-go/          # DeepSeek API client
│   ├── Dockerfile                    # Multi-stage build
│   ├── go.mod & go.sum               # Go dependencies
│   └── .env.example                  # Environment template
│
├── frontend/                         # Next.js frontend
│   ├── src/
│   │   ├── app/                      # Next.js app router
│   │   │   ├── page.tsx              # Root page (redirect)
│   │   │   ├── login/                # Login page
│   │   │   ├── register/             # Registration page
│   │   │   └── dashboard/            # Protected routes
│   │   │       ├── page.tsx          # Dashboard home
│   │   │       ├── chat/             # Chat interface
│   │   │       ├── system/           # System monitor
│   │   │       └── games/            # Games search
│   │   ├── components/               # Reusable components
│   │   │   ├── Button.tsx            # Button with animations
│   │   │   ├── Input.tsx             # Input with validation
│   │   │   ├── Card.tsx              # Card component
│   │   │   ├── Skeleton.tsx          # Loading skeleton
│   │   │   └── DashboardLayout.tsx   # Layout wrapper
│   │   ├── context/                  # React context
│   │   │   └── AuthContext.tsx       # Auth state management
│   │   ├── lib/                      # Utilities
│   │   │   └── api-client.ts         # API client
│   │   └── app/globals.css           # Global styles
│   ├── Dockerfile                    # Multi-stage build
│   ├── package.json                  # Dependencies
│   └── .env.local                    # Frontend config
│
├── docker-compose.yml                # Service orchestration
├── .env                              # Environment variables
├── .env.example                      # Environment template
├── Makefile                          # Build automation
├── start.sh                          # Quick start script
└── README.md                         # This file
```

## 🔐 Security Features

### Authentication & Authorization
- **JWT-based authentication** with short-lived access tokens
- **Refresh token rotation** for enhanced security
- **Role-Based Access Control (RBAC)** for resource protection
- **Multi-Factor Authentication (MFA)** ready infrastructure

### Data Protection
- **Password hashing** with Bcrypt + Argon2
- **HTTPS/TLS 1.2+** enforcement for all communications
- **Data encryption at rest** for sensitive information
- **CORS configuration** to prevent cross-origin attacks

### API Security
- **Rate limiting** to prevent brute force attacks
- **SQL injection prevention** via parameterized queries (GORM)
- **XSS protection** with Content Security Policy (CSP)
- **CSRF tokens** for state-changing operations
- **Input validation** on all endpoints

### Secrets Management
- **Environment variables** for API keys and credentials
- **No hardcoded secrets** in source code
- **Automatic key rotation** support
- **Secure secret storage** in production (AWS Secrets Manager, Vault)

## 🚀 Features

### Core Capabilities

#### 1. **AI-Powered Chat Interface**
- Real-time conversation with DeepSeek R1/V3 or GPT-4o
- Context-aware responses with semantic memory
- Multi-turn conversation support
- Fallback mechanism for API failures

#### 2. **Semantic Memory System**
- Long-term memory storage using Qdrant vector database
- Semantic search across conversation history
- Automatic memory consolidation
- Privacy-preserving embeddings

#### 3. **Real-Time System Monitoring**
- CPU, Memory, Disk, GPU metrics
- Network I/O statistics
- Cross-platform support (Windows, Linux, macOS)
- Historical data tracking

#### 4. **Game Compatibility Checker**
- Search game database
- Check system requirements vs. actual hardware
- Compatibility recommendations
- Performance predictions

#### 5. **Text-to-Speech (TTS)**
- Local TTS synthesis (Piper/Coqui)
- Low-latency audio generation
- Multiple voice options
- ElevenLabs integration (optional)

## 📊 API Endpoints

### Authentication
```
POST   /api/auth/login              # User login
POST   /api/auth/register           # User registration
POST   /api/auth/refresh            # Refresh access token
```

### AI & Chat
```
POST   /api/ai/chat                 # Send message to AI
```

### System Information
```
GET    /api/system-info/latest      # Get current system metrics
GET    /api/system-info/history     # Get historical metrics
```

### Text-to-Speech
```
POST   /api/tts/generate            # Generate speech from text
```

### Memory Management
```
POST   /api/memory                  # Save memory
GET    /api/memory                  # Get all memories
POST   /api/memory/search           # Search memories
```

### Games
```
GET    /api/games/search            # Search games
GET    /api/games/:id               # Get game details
```

## 🛠️ Development

### Backend Development

```bash
# Install dependencies
cd backend
go mod download

# Run tests
go test ./... -v

# Build locally
go build -o jarvis ./cmd/main.go

# Run with hot reload (requires air)
air
```

### Frontend Development

```bash
# Install dependencies
cd frontend
npm install

# Development server
npm run dev

# Build for production
npm run build

# Run production build
npm start

# Run tests
npm test
```

## 📦 Deployment

### Docker Compose (Development/Staging)
```bash
# Start all services
make dev

# View logs
make logs

# Stop services
make down
```

### Production Deployment

1. **Set production environment variables**
   ```bash
   export APP_ENV=production
   export SECURITY_JWT_SECRET=<strong-random-secret>
   export DATABASE_PASSWORD=<strong-password>
   export REDIS_PASSWORD=<strong-password>
   ```

2. **Deploy with Docker Compose**
   ```bash
   docker compose -f docker-compose.yml up -d
   ```

3. **Configure reverse proxy (Nginx/Traefik)**
   ```nginx
   server {
       listen 443 ssl http2;
       server_name api.jarvis.example.com;
       
       ssl_certificate /etc/letsencrypt/live/jarvis.example.com/fullchain.pem;
       ssl_certificate_key /etc/letsencrypt/live/jarvis.example.com/privkey.pem;
       
       location / {
           proxy_pass http://localhost:8080;
           proxy_set_header Host $host;
           proxy_set_header X-Real-IP $remote_addr;
       }
   }
   ```

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./... -v -cover
```

### Frontend Tests
```bash
cd frontend
npm test
```

### Integration Tests
```bash
# Run full stack tests
make test
```

## 📈 Performance Optimization

### Backend
- **Connection pooling** for database and cache
- **Query optimization** with proper indexing
- **Caching layer** with Redis for frequently accessed data
- **Async processing** for long-running tasks
- **Compression** for API responses

### Frontend
- **Code splitting** for faster initial load
- **Image optimization** with Next.js Image component
- **Lazy loading** for routes and components
- **Service workers** for offline support
- **CSS-in-JS** with Tailwind for minimal bundle size

## 🐛 Troubleshooting

### Services won't start
```bash
# Check Docker is running
docker ps

# View logs
docker compose logs

# Rebuild images
docker compose build --no-cache
```

### Database connection errors
```bash
# Check MySQL is healthy
docker compose ps mysql

# Reset database
docker compose down -v
docker compose up -d
```

### API connection issues
```bash
# Check backend is running
curl http://localhost:8080/health

# Check CORS configuration
# Verify SECURITY_CORS_ORIGINS in .env
```

## 📚 Documentation

- [API Documentation](./docs/API.md)
- [Architecture Guide](./docs/ARCHITECTURE.md)
- [Deployment Guide](./docs/DEPLOYMENT.md)
- [Contributing Guide](./CONTRIBUTING.md)

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- DeepSeek for the powerful R1/V3 models
- OpenAI for GPT-4o fallback support
- Qdrant for vector database excellence
- The Go and React communities

## 📞 Support

For issues, questions, or suggestions:
- Open an issue on GitHub
- Check existing documentation
- Review the troubleshooting guide

---

**Made with ❤️ by the JARVIS Team**

*Last Updated: 2026-06-19*
