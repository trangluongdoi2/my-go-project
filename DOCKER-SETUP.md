# Docker Network Setup Guide

Hướng dẫn setup PostgreSQL, RabbitMQ, Redis và Go application trong cùng một Docker network.

---

## Tổng quan

Tất cả services sẽ chạy trong Docker containers và kết nối với nhau thông qua một Docker network chung.

```
┌─────────────────── Docker Network: app-network ───────────────────┐
│                                                                    │
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐   │
│  │PostgreSQL│◄───┤   Redis  │◄───┤ RabbitMQ │◄───┤  Go App  │   │
│  │  :5432   │    │  :6379   │    │  :5672   │    │  :9512   │   │
│  └──────────┘    └──────────┘    └──────────┘    └──────────┘   │
│                                                                    │
└────────────────────────────────────────────────────────────────────┘
```

---

## 1. Tạo Docker Compose File

Tạo file `docker-compose.yml` ở root project:

```yaml
version: '3.8'

# Define custom network
networks:
  app-network:
    driver: bridge

# Define volumes for data persistence
volumes:
  postgres_data:
  rabbitmq_data:
  redis_data:

services:
  # PostgreSQL Database
  postgres:
    image: postgres:14-alpine
    container_name: pos-postgres
    networks:
      - app-network
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: nail-db
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

  # RabbitMQ Message Broker
  rabbitmq:
    image: rabbitmq:3.12-management-alpine
    container_name: pos-rabbitmq
    networks:
      - app-network
    ports:
      - "5672:5672"   # AMQP port
      - "15672:15672" # Management UI port
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
      RABBITMQ_DEFAULT_VHOST: /
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Redis Cache
  redis:
    image: redis:7-alpine
    container_name: pos-redis
    networks:
      - app-network
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5

  # Go Application (Optional - nếu muốn chạy trong Docker)
  app:
    build:
      context: .
      dockerfile: Dockerfile
    container_name: pos-api
    networks:
      - app-network
    ports:
      - "9512:9512"
    environment:
      CONFIG_PATH: /app/internal/config/docker.yaml
      ENV: docker
    depends_on:
      postgres:
        condition: service_healthy
      rabbitmq:
        condition: service_healthy
      redis:
        condition: service_healthy
    volumes:
      - ./internal/config:/app/internal/config
    restart: unless-stopped
```

---

## 2. Tạo Init SQL Script

Tạo folder và file `scripts/init.sql`:

```bash
mkdir -p scripts
```

File `scripts/init.sql`:

```sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create bookings table
CREATE TABLE IF NOT EXISTS bookings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer JSONB,
    service_name VARCHAR(255) NOT NULL,
    staff_id UUID,
    appointment_at TIMESTAMP NOT NULL,
    price DECIMAL(10,2),
    status SMALLINT DEFAULT 1,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create staff table
CREATE TABLE IF NOT EXISTS staff (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(50),
    email VARCHAR(255),
    position VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX idx_bookings_staff_id ON bookings(staff_id);
CREATE INDEX idx_bookings_appointment_at ON bookings(appointment_at);
CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_deleted_at ON bookings(deleted_at);

-- Insert sample data
INSERT INTO staff (name, phone, email, position) VALUES
    ('Alice Johnson', '0901234567', 'alice@example.com', 'Senior Technician'),
    ('Bob Smith', '0912345678', 'bob@example.com', 'Technician')
ON CONFLICT DO NOTHING;
```

---

## 3. Tạo Dockerfile cho Go App

Tạo file `Dockerfile` ở root project:

```dockerfile
# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/internal/config ./internal/config

EXPOSE 9512

CMD ["./main"]
```

---

## 4. Tạo Docker Config File

Tạo file `internal/config/docker.yaml`:

```yaml
env:
  mode: debug

server:
  port: 9512
  locale: 'en'
  timezone: UTC

database:
  host: postgres          # ← Tên service trong docker-compose
  port: 5432
  username: postgres
  password: postgres
  db_name: nail-db

rabbitmq:
  schema: amqp
  worker: 2
  host: rabbitmq          # ← Tên service trong docker-compose
  port: 5672
  vhost: /
  username: guest
  password: guest
  ssl: false
  ctag: nail-project
  retry: 3
  exchanges: ex.direct.delay
  queue: {}
  uri: amqp://guest:guest@rabbitmq:5672  # ← Dùng tên service

redis:
  host: redis             # ← Tên service trong docker-compose
  port: 6379
  password:
  db: 0
```

---

## 5. Tạo .dockerignore

Tạo file `.dockerignore`:

```
.git
.gitignore
README.md
SETUP.md
DOCKER-SETUP.md
*.md
.env
.vscode
.idea
bin/
tmp/
*.log
```

---

## 6. Commands để Chạy

### 6.1. Start tất cả services

```bash
# Start tất cả services (detached mode)
docker-compose up -d

# Hoặc start với logs
docker-compose up

# Chỉ start infrastructure (không start app)
docker-compose up -d postgres rabbitmq redis
```

### 6.2. Xem logs

```bash
# Xem logs tất cả services
docker-compose logs -f

# Xem logs của service cụ thể
docker-compose logs -f postgres
docker-compose logs -f rabbitmq
docker-compose logs -f redis
docker-compose logs -f app
```

### 6.3. Stop services

```bash
# Stop tất cả
docker-compose down

# Stop và xóa volumes (CẢNH BÁO: Mất data!)
docker-compose down -v
```

### 6.4. Restart services

```bash
# Restart tất cả
docker-compose restart

# Restart service cụ thể
docker-compose restart postgres
docker-compose restart app
```

---

## 7. Chạy Go App Local + Docker Services

Nếu bạn muốn chạy Go app ở local nhưng dùng Docker cho services:

### 7.1. Start chỉ infrastructure

```bash
docker-compose up -d postgres rabbitmq redis
```

### 7.2. Cập nhật config file

File `internal/config/development.yaml`:

```yaml
database:
  host: localhost        # ← localhost vì app chạy ngoài Docker
  port: 5432
  # ...

rabbitmq:
  host: localhost
  port: 5672
  uri: amqp://guest:guest@localhost:5672
  # ...

redis:
  host: localhost
  port: 6379
  # ...
```

### 7.3. Run Go app

```bash
go run cmd/api/main.go
```

---

## 8. Kiểm tra Services

### 8.1. Check containers đang chạy

```bash
docker-compose ps
```

Output mẫu:
```
NAME              IMAGE                              STATUS         PORTS
pos-api           go-backend-project-app             Up 2 minutes   0.0.0.0:9512->9512/tcp
pos-postgres      postgres:14-alpine                 Up 2 minutes   0.0.0.0:5432->5432/tcp
pos-rabbitmq      rabbitmq:3.12-management-alpine    Up 2 minutes   0.0.0.0:5672->5672/tcp, 0.0.0.0:15672->15672/tcp
pos-redis         redis:7-alpine                     Up 2 minutes   0.0.0.0:6379->6379/tcp
```

### 8.2. Check network

```bash
# List networks
docker network ls

# Inspect network
docker network inspect go-backend-project_app-network
```

### 8.3. Test connections

```bash
# Test PostgreSQL
docker exec -it pos-postgres psql -U postgres -d nail-db -c "SELECT version();"

# Test Redis
docker exec -it pos-redis redis-cli ping

# Test RabbitMQ
curl -u guest:guest http://localhost:15672/api/overview
```

---

## 9. Access Services

| Service           | URL/Command                                  |
|-------------------|----------------------------------------------|
| PostgreSQL        | `localhost:5432`                             |
| RabbitMQ AMQP     | `localhost:5672`                             |
| RabbitMQ UI       | http://localhost:15672 (guest/guest)         |
| Redis             | `localhost:6379`                             |
| API               | http://localhost:9512                        |

### Kết nối từ host machine

```bash
# PostgreSQL
psql -h localhost -p 5432 -U postgres -d nail-db

# Redis
redis-cli -h localhost -p 6379

# API
curl http://localhost:9512/bookings
```

---

## 10. Development Workflow

### 10.1. Initial Setup

```bash
# 1. Clone repo
git clone <repo-url>
cd go-backend-project

# 2. Start Docker services
docker-compose up -d

# 3. Wait for services to be healthy
docker-compose ps

# 4. Test API
curl http://localhost:9512/bookings
```

### 10.2. Daily Development

**Option 1: Tất cả trong Docker**
```bash
# Start
docker-compose up -d

# Rebuild app sau khi thay đổi code
docker-compose up -d --build app

# Stop
docker-compose down
```

**Option 2: App local + Services Docker**
```bash
# Start services
docker-compose up -d postgres rabbitmq redis

# Run app
go run cmd/api/main.go

# Stop services
docker-compose down
```

---

## 11. Troubleshooting

### Port đã được sử dụng

```bash
# Kiểm tra port
lsof -i :5432
lsof -i :5672
lsof -i :6379

# Đổi port trong docker-compose.yml
# Ví dụ: "5433:5432" thay vì "5432:5432"
```

### Container không start

```bash
# Xem logs
docker-compose logs <service-name>

# Restart service
docker-compose restart <service-name>

# Rebuild
docker-compose up -d --build <service-name>
```

### Database connection refused

```bash
# Kiểm tra health
docker-compose ps

# Restart postgres
docker-compose restart postgres

# Check logs
docker-compose logs postgres
```

### Xóa và tạo lại từ đầu

```bash
# Stop và xóa tất cả (bao gồm volumes)
docker-compose down -v

# Start lại
docker-compose up -d
```

---

## 12. Production Tips

### 12.1. Environment Variables

Tạo file `.env`:

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=secure_password_here
POSTGRES_DB=nail-db

RABBITMQ_USER=admin
RABBITMQ_PASS=secure_password_here

REDIS_PASSWORD=secure_password_here
```

Cập nhật `docker-compose.yml`:

```yaml
services:
  postgres:
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
```

### 12.2. Security

- Đổi default passwords
- Không expose ports không cần thiết
- Sử dụng secrets management
- Enable SSL/TLS cho production

---

## Quick Start Commands

```bash
# Start everything
docker-compose up -d

# View logs
docker-compose logs -f

# Stop everything
docker-compose down

# Rebuild app
docker-compose up -d --build app

# Clean everything (including data)
docker-compose down -v
```

Happy Dockerizing! 🐳