# JARVIS — AI Assistant

> Assistente de IA pessoal com monitoramento de hardware em tempo real, chat inteligente em português e verificação de compatibilidade de jogos.

![Status](https://img.shields.io/badge/status-ativo-brightgreen)
![License](https://img.shields.io/badge/license-MIT-blue)
![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)
![Next.js](https://img.shields.io/badge/Next.js-16-black?logo=next.js)
![React](https://img.shields.io/badge/React-19-61DAFB?logo=react)
![TypeScript](https://img.shields.io/badge/TypeScript-5-3178C6?logo=typescript)

---

## Sobre o projeto

JARVIS é um assistente de IA que conhece o seu PC. Ele lê suas especificações reais de hardware e usa esse contexto em todas as respostas — seja para dizer se um jogo roda na sua máquina, analisar gargalos de performance ou simplesmente conversar em português.

### Telas

**Dashboard — monitoramento em tempo real**

![Dashboard](https://raw.githubusercontent.com/Dereckgab/jarvis-ai/main/docs/screenshots/dashboard.png)

> Uso de CPU, RAM, Disco e GPU com anéis animados e specs detalhadas do hardware.

**Games — verificação de compatibilidade**

![Games](https://raw.githubusercontent.com/Dereckgab/jarvis-ai/main/docs/screenshots/games.png)

> Pesquise qualquer jogo: JARVIS mostra a capa, um veredicto (RODA / NÃO RODA / DEPENDE), tabela comparativa das suas specs vs requisitos mínimos e recomendados, e dicas de otimização.

**Chat — IA em português**

![Chat](https://raw.githubusercontent.com/Dereckgab/jarvis-ai/main/docs/screenshots/chat.png)

> Converse com a IA sobre jogos, hardware ou qualquer assunto. Quando a pergunta é sobre compatibilidade, JARVIS dá um veredicto direto e oferece análise detalhada com um clique.

---

## Funcionalidades

- **Chat com IA** — integrado com Groq (llama-3.3-70b-versatile), responde em português, usa suas specs como contexto
- **Verificação de jogos** — capa real via RAWG API, veredicto animado, tabela de requisitos com indicadores coloridos
- **Em alta esta semana** — 5 jogos mais jogados da semana carregados automaticamente na página de Games
- **Monitoramento de hardware** — CPU, RAM, Disco, GPU, Rede, temperatura — atualizado a cada 5 segundos
- **Autenticação segura** — JWT com refresh token, bcrypt nas senhas

---

## Stack

| Camada | Tecnologia |
|--------|-----------|
| Frontend | Next.js 16 · React 19 · TypeScript · Tailwind v4 · Framer Motion |
| Backend | Go 1.25 · Fiber v2 · GORM · Clean Architecture |
| Banco de dados | MySQL 8.0 |
| Cache | Redis 7 |
| IA | Groq API (llama-3.3-70b-versatile) |
| Imagens de jogos | RAWG API |
| Monitoramento | gopsutil (CPU, RAM, Disco, GPU, Rede) |
| Infraestrutura | Docker · Docker Compose |

---

## Arquitetura do backend

```
backend/
├── cmd/                  # Entry point
├── config/               # Configuração via .env
└── internal/
    ├── domain/           # Entidades e interfaces (puro Go, sem dependências)
    ├── application/      # Serviços e regras de negócio
    ├── infrastructure/   # MySQL, Redis, Groq, gopsutil
    └── interfaces/       # Handlers HTTP, middlewares, DTOs
```

Segue os princípios de **Clean Architecture**: a camada de domínio não depende de nada externo. Repositórios são interfaces — a implementação fica na infraestrutura.

---

## Como rodar localmente

### Pré-requisitos

- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Go 1.21+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)

### 1. Clone o repositório

```bash
git clone https://github.com/Dereckgab/jarvis-ai.git
cd jarvis-ai
```

### 2. Configure as variáveis de ambiente

```bash
cp backend/.env.example backend/.env
```

Edite `backend/.env` e preencha:

```env
AI_GROQ_API_KEY=sua_chave_aqui        # groq.com (gratuito)
SECURITY_JWT_SECRET=uma_chave_segura
```

Crie `frontend/.env.local`:

```env
NEXT_PUBLIC_RAWG_API_KEY=sua_chave_aqui   # rawg.io (gratuito)
```

### 3. Suba o banco de dados e o Redis

```bash
docker compose up -d
```

### 4. Rode o backend

```bash
cd backend
go run cmd/server/main.go
```

### 5. Rode o frontend

```bash
cd frontend
npm install
npm run dev
```

Acesse **http://localhost:3000**

---

## Endpoints da API

```
POST  /api/auth/register         Cadastro
POST  /api/auth/login            Login
POST  /api/auth/refresh          Refresh token

POST  /api/ai/chat               Chat com a IA

GET   /api/system-info/latest    Specs atuais do hardware
GET   /api/system-info/history   Histórico de métricas

GET   /api/games/search          Busca de jogos
```

---

## Variáveis de ambiente

Veja [`backend/.env.example`](backend/.env.example) para a lista completa.

As principais:

| Variável | Descrição |
|----------|-----------|
| `AI_GROQ_API_KEY` | Chave da Groq API (obrigatório) |
| `AI_MODEL` | Modelo de IA (padrão: `llama-3.3-70b-versatile`) |
| `DATABASE_*` | Configurações do MySQL |
| `SECURITY_JWT_SECRET` | Chave secreta para JWT |
| `NEXT_PUBLIC_RAWG_API_KEY` | Chave da RAWG API para capas de jogos |

---

## Autor

Desenvolvido por **Gabriel** — estudante de desenvolvimento focado em Go e React.

[![GitHub](https://img.shields.io/badge/GitHub-Dereckgab-181717?logo=github)](https://github.com/Dereckgab)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-gabriel--dereck-0A66C2?logo=linkedin)](https://linkedin.com/in/gabriel-dereck)

---

## Licença

MIT — veja [LICENSE](LICENSE) para detalhes.
