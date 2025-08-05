# Single-stage build for Railway compatibility
FROM golang:1.24-alpine

WORKDIR /app

# Install minimal dependencies with specific Node.js version
RUN apk add --no-cache nodejs npm git ca-certificates && \
    npm config set registry https://registry.npmjs.org/ && \
    npm cache clean --force && \
    node --version && \
    npm --version

# Copy frontend files and build
COPY frontend/ ./frontend/
WORKDIR /app/frontend
RUN npm install --production=false --no-optional
RUN ls -la && echo "Starting build..." && npm run build

# Back to app root and build backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

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