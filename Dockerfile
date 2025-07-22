# 1. Build the frontend with Yarn for lower memory usage
FROM node:20 AS frontend-build
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm install -g yarn
RUN yarn install
COPY frontend/ ./
RUN yarn build

# 2. Build the backend
FROM golang:1.21 AS backend-build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go

# 3. Final image
FROM debian:bookworm-slim
WORKDIR /app

# Copy backend binary
COPY --from=backend-build /app/main .

# Copy frontend build output
COPY --from=frontend-build /app/frontend/dist ./frontend/dist

# (Optional) Install runtime dependencies
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

EXPOSE 8080

# Start the backend (which should serve API and/or static files)
CMD ["./main"] 