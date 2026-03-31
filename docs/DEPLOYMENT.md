# Go Backend Deployment Guide

## Prerequisites

- Go 1.23+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose (optional)

## Quick Start

### Option 1: Docker Compose (Recommended)

```bash
cd go-backend

# Start all services
docker-compose up -d

# Check logs
docker-compose logs -f go-backend

# Stop services
docker-compose down
```

### Option 2: Manual Deployment

#### 1. Setup Database

```bash
# PostgreSQL
createdb digitalhuman
psql -d digitalhuman -f migrations/001_init.sql

# Redis
redis-server
```

#### 2. Configure Environment

```bash
export DB_HOST=localhost
export DB_NAME=digitalhuman
export DB_USER=postgres
export DB_PASSWORD=your_password
export REDIS_HOST=localhost
export LOG_LEVEL=info
```

#### 3. Build and Run

```bash
# Build
go build -o bin/server ./cmd/server

# Run
./bin/server
```

## Configuration

Edit `configs/config.yaml`:

```yaml
common:
  name: "Awesome-Digital-Human"
  version: "v3.0.0"
  log_level: "debug"

server:
  ip: "0.0.0.0"
  port: 8881

database:
  host: "${DB_HOST:localhost}"
  port: 5432
  name: "${DB_NAME:digitalhuman}"
  user: "${DB_USER:postgres}"
  password: "${DB_PASSWORD:}"

redis:
  host: "${REDIS_HOST:localhost}"
  port: 6379

agents:
  support_list:
    - name: "RepeaterAgent"
      type: "AGENT"
      desc: "Repeat user input"
    - name: "OpenAIAgent"
      type: "AGENT"
      desc: "OpenAI GPT Agent"
      config:
        api_key: "${OPENAI_API_KEY}"
        model: "gpt-3.5-turbo"
  default: "RepeaterAgent"

engines:
  tts:
    support_list:
      - name: "EdgeTTS"
        type: "TTS"
        desc: "Microsoft Edge TTS"
    default: "EdgeTTS"
  asr:
    support_list: []
    default: ""
```

## Nginx Configuration

```nginx
# Proxy Go backend
location /adh/agent/ {
    proxy_pass http://127.0.0.1:8881;
    proxy_http_version 1.1;
    proxy_buffering off;
}

location /adh/tts/ {
    proxy_pass http://127.0.0.1:8881;
}

location /adh/asr/ {
    proxy_pass http://127.0.0.1:8881;
}
```

## API Endpoints

### Agent API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/adh/agent/v0/engine` | GET | List agents |
| `/adh/agent/v0/engine/default` | GET | Get default agent |
| `/adh/agent/v0/engine/:engine` | GET | Get agent params |
| `/adh/agent/v0/engine/:engine` | POST | Create conversation |
| `/adh/agent/v0/engine` | POST | Agent inference (SSE) |

### TTS API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/adh/tts/v0/engine` | GET | List TTS engines |
| `/adh/tts/v0/engine/default` | GET | Get default engine |
| `/adh/tts/v0/engine/:engine` | GET | Get engine params |
| `/adh/tts/v0/engine/:engine/voice` | GET | Get voice list |
| `/adh/tts/v0/engine` | POST | TTS inference |

### ASR API

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/adh/asr/v0/engine` | GET | List ASR engines |
| `/adh/asr/v0/engine/default` | GET | Get default engine |
| `/adh/asr/v0/engine/:engine` | GET | Get engine params |
| `/adh/asr/v0/engine` | POST | ASR inference |
| `/adh/asr/v0/engine/file` | POST | ASR file inference |
| `/adh/asr/v0/engine/stream` | WS | Streaming ASR |

## Health Check

```bash
curl http://localhost:8881/health
```

## Monitoring

### Logs

```bash
# Docker
docker-compose logs -f go-backend

# Manual
tail -f /var/log/adh-go-backend.log
```

### Metrics

- Port: 8881
- Health: `/health`

## Troubleshooting

### Database Connection Failed

```bash
# Check PostgreSQL
psql -h localhost -U postgres -d digitalhuman

# Check Redis
redis-cli ping
```

### Port Already in Use

```bash
# Find process
lsof -i :8881

# Kill process
kill -9 <PID>
```

### Configuration Not Loading

```bash
# Check environment
echo $CONFIG_PATH

# Verify config file
cat configs/config.yaml
```