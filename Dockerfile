# Single-stage build for Railway compatibility
FROM node:20-alpine

WORKDIR /app

# Install Go 1.24.5 from official binary
RUN apk add --no-cache git ca-certificates wget && \
    wget -O go.tar.gz https://go.dev/dl/go1.24.5.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm go.tar.gz && \
    ln -s /usr/local/go/bin/go /usr/local/bin/go && \
    npm config set registry https://registry.npmjs.org/ && \
    npm cache clean --force && \
    node --version && \
    npm --version && \
    go version

# Copy frontend files and build (SvelteKit)
COPY frontend/ ./frontend/
WORKDIR /app/frontend
RUN npm install --production=false
RUN npm rebuild
RUN npm install @rollup/rollup-linux-x64-musl
RUN ls -la && echo "Starting SvelteKit build..." && npm run build

# Back to app root and build backend
WORKDIR /app
COPY go.mod go.sum ./

# Copy backend files (excluding frontend to avoid overwriting built assets)
COPY cmd/ ./cmd/
COPY internal/ ./internal/
COPY railway.json ./

RUN go mod download

# Verify SvelteKit build output exists
RUN ls -la frontend/dist/ && echo "SvelteKit build successful"

# Build Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main ./cmd/main.go

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/teams || exit 1

CMD ["./main"] 