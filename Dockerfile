# Stage 1: Build the dashboard (frontend)
FROM node:22-alpine AS dashboard-builder

WORKDIR /workspace

# Install pnpm
RUN corepack enable && corepack prepare pnpm@10.17.0 --activate

# Copy package files
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/dash/package.json ./apps/dash/

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy dashboard source
COPY apps/dash ./apps/dash
COPY tsconfig.json ./

# Build dashboard
RUN pnpm run dash:build

# Stage 2: Build the Go backend
FROM golang:1.25 AS builder

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY migrations ./migrations
COPY core ./core
COPY internal ./internal
COPY store ./store
COPY stremio ./stremio
COPY *.go ./

# Copy built dashboard from previous stage
COPY --from=dashboard-builder /workspace/apps/dash/.output/public/ ./internal/dash/fs/

RUN CGO_ENABLED=1 GOOS=linux go build --tags 'fts5' -o ./stremthru -a -ldflags '-linkmode external -extldflags "-static"'

# Stage 3: Final runtime image
FROM alpine

RUN apk add --no-cache git

WORKDIR /app

COPY --from=builder /workspace/stremthru ./stremthru

# Create data directory
RUN mkdir -p /app/data

VOLUME ["/app/data"]

ENV STREMTHRU_ENV=prod

EXPOSE 8080

ENTRYPOINT []
CMD ["./stremthru"]
