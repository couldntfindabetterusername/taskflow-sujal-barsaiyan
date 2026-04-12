# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates postgresql-client

# Install golang-migrate
RUN wget -O /tmp/migrate.tar.gz https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz && \
    tar -xzf /tmp/migrate.tar.gz -C /usr/local/bin && \
    rm /tmp/migrate.tar.gz

WORKDIR /root/
COPY --from=builder /app/server .
COPY migrations ./migrations

# Create entrypoint script to run migrations before starting server
RUN echo '#!/bin/sh' > /root/entrypoint.sh && \
    echo 'set -e' >> /root/entrypoint.sh && \
    echo '' >> /root/entrypoint.sh && \
    echo '# Wait for PostgreSQL to be ready' >> /root/entrypoint.sh && \
    echo 'echo "Waiting for PostgreSQL to be ready..."' >> /root/entrypoint.sh && \
    echo 'until pg_isready -h $(echo $DATABASE_URL | sed "s/.*@\(.*\):.*/\1/") -U $(echo $DATABASE_URL | sed "s/.*:\/\/\(.*\):.*/\1/"); do' >> /root/entrypoint.sh && \
    echo '  sleep 1' >> /root/entrypoint.sh && \
    echo 'done' >> /root/entrypoint.sh && \
    echo '' >> /root/entrypoint.sh && \
    echo 'echo "PostgreSQL is ready!"' >> /root/entrypoint.sh && \
    echo '' >> /root/entrypoint.sh && \
    echo '# Run database migrations' >> /root/entrypoint.sh && \
    echo 'echo "Running database migrations..."' >> /root/entrypoint.sh && \
    echo 'migrate -path=/root/migrations -database "$DATABASE_URL" up' >> /root/entrypoint.sh && \
    echo '' >> /root/entrypoint.sh && \
    echo 'echo "Migrations completed successfully!"' >> /root/entrypoint.sh && \
    echo '' >> /root/entrypoint.sh && \
    echo '# Start the server' >> /root/entrypoint.sh && \
    echo 'echo "Starting TaskFlow server..."' >> /root/entrypoint.sh && \
    echo 'exec ./server' >> /root/entrypoint.sh && \
    chmod +x /root/entrypoint.sh

EXPOSE 8080
CMD ["/root/entrypoint.sh"]
