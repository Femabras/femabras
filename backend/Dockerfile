# ----------------------------------
# Stage 1: The Builder
# ----------------------------------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Install certificates so Go can fetch packages securely
RUN apk add --no-cache ca-certificates git

# Copy module files first to cache dependencies (saves build time)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of your Go code
COPY . .

# Compile the app. CGO_ENABLED=0 makes it a standalone static binary!
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o femabras-server ./cmd/server

# ----------------------------------
# Stage 2: The Lean Runner (Production)
# ----------------------------------
FROM alpine:latest

WORKDIR /app

# We need certificates here too, otherwise Google Auth will fail!
RUN apk add --no-cache ca-certificates

# Copy ONLY the compiled binary from the builder stage
COPY --from=builder /app/femabras-server .

EXPOSE 8080

CMD ["./femabras-server"]