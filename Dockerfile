# 1. Build the frontend with Yarn for lower memory usage
FROM node:20-alpine AS frontend-build
WORKDIR /app/frontend
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL
COPY frontend/package*.json ./
RUN npm ci --only=production && npm cache clean --force
COPY frontend/ ./
RUN npm run build

# 2. Build the backend
FROM golang:1.24-alpine AS backend-build
WORKDIR /app

# Install security updates and required packages
RUN apk update && apk upgrade && \
    apk add --no-cache git ca-certificates tzdata && \
    update-ca-certificates

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build with security flags
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o main ./cmd/main.go

# 3. Final minimal security-hardened image
FROM alpine:latest

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Install security updates and minimal runtime dependencies
RUN apk update && apk upgrade && \
    apk add --no-cache ca-certificates tzdata && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*

WORKDIR /app

# Copy backend binary
COPY --from=backend-build /app/main .

# Copy frontend build output
COPY --from=frontend-build /app/frontend/dist ./frontend/dist

# Copy assets to dist folder for static serving
COPY frontend/src/assets ./frontend/dist/assets

# Set proper ownership
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/v1/teams || exit 1

# Start the backend (which should serve API and/or static files)
CMD ["./main"] 