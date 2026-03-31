# Awesome Digital Human - Go Backend

Go implementation of the Python FastAPI backend for the Awesome Digital Human Live2D project.

## Features

- **Agent API**: OpenAI, Dify, Coze, FastGPT, Repeater agents with SSE streaming
- **TTS API**: Edge TTS with 30+ voices
- **ASR API**: HTTP and WebSocket streaming ASR
- **PostgreSQL + Redis**: Persistent storage with caching

## Tech Stack

| Component | Technology |
|-----------|------------|
| HTTP Framework | Gin |
| WebSocket | gorilla/websocket |
| Configuration | Viper |
| Logging | Zap |
| Database | PostgreSQL |
| Cache | Redis |

## Project Structure

```
go-backend/
├── cmd/server/main.go          # Entry point
├── configs/config.yaml         # Configuration
├── internal/
│   ├── config/                 # Config loader
│   ├── protocol/               # API protocol definitions
│   ├── agent/                  # Agent implementations
│   ├── engine/                 # TTS/ASR engines
│   ├── storage/                # Database layer
│   ├── server/                 # HTTP handlers
│   └── pkg/logger/             # Logging
├── migrations/                 # SQL migrations
├── deploy/nginx/               # Nginx config
├── Dockerfile
├── docker-compose.yml
└── scripts/test_api.sh         # API tests
```

## Quick Start

### Docker (Recommended)

```bash
docker-compose up -d
```

### Manual

```bash
# Install dependencies
go mod download

# Setup database
psql -d digitalhuman -f migrations/001_init.sql

# Run
go run ./cmd/server
```

## API Endpoints

### Agent

```
GET  /adh/agent/v0/engine           # List agents
GET  /adh/agent/v0/engine/default   # Default agent
GET  /adh/agent/v0/engine/:engine   # Agent params
POST /adh/agent/v0/engine/:engine   # Create conversation
POST /adh/agent/v0/engine           # Agent inference (SSE)
```

### TTS

```
GET  /adh/tts/v0/engine             # List TTS engines
GET  /adh/tts/v0/engine/default     # Default engine
GET  /adh/tts/v0/engine/:engine     # Engine params
GET  /adh/tts/v0/engine/:engine/voice  # Voice list
POST /adh/tts/v0/engine             # TTS inference
```

### ASR

```
GET  /adh/asr/v0/engine             # List ASR engines
GET  /adh/asr/v0/engine/default     # Default engine
GET  /adh/asr/v0/engine/:engine     # Engine params
POST /adh/asr/v0/engine             # ASR inference
POST /adh/asr/v0/engine/file        # File inference
WS   /adh/asr/v0/engine/stream      # Streaming ASR
```

## Test

```bash
./scripts/test_api.sh
```

## Documentation

- [Deployment Guide](docs/DEPLOYMENT.md)
- [Design Spec](../docs/superpowers/specs/2026-03-30-go-backend-migration-design.md)

## License

MIT