# ============================
# Base Builder Stage
# ============================
FROM golang:1.21-alpine AS base-builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev sqlite ca-certificates

WORKDIR /app

# Copy Go mod files first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# ============================
# Service Build Stages
# ============================
FROM base-builder AS red-builder
RUN go build -o /app/bin/redteam ./redteam/main.go

FROM base-builder AS blue-builder
RUN go build -o /app/bin/blueteam ./blueteam/main.go

FROM base-builder AS dashboard-builder
RUN go build -o /app/bin/dashboard ./dashboard/main.go

# ============================
# Runtime Base
# ============================
FROM alpine:3.18 AS runtime-base
RUN apk --no-cache add ca-certificates && \
    mkdir -p /data && chmod a+rw /data  # Critical fix for SQLite permissions
WORKDIR /app

# ============================
# Service Runtime Stages
# ============================
FROM runtime-base AS red-runtime
COPY --from=red-builder /app/bin/redteam .
EXPOSE 8082
CMD ["./redteam"]

FROM runtime-base AS blue-runtime
COPY --from=blue-builder /app/bin/blueteam .
EXPOSE 8081
CMD ["./blueteam"]

FROM runtime-base AS dashboard-runtime
COPY --from=dashboard-builder /app/bin/dashboard .
EXPOSE 8080
CMD ["./dashboard"]

# ============================
# Frontend Build & Runtime
# ============================
FROM node:18-alpine AS frontend-builder
WORKDIR /app
COPY frontend/package*.json ./
COPY frontend/.env .env
RUN npm install
COPY frontend/ .
RUN npm run build

FROM nginx:alpine AS frontend-runtime
COPY --from=frontend-builder /app/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]