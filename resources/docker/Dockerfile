# Development stage with hot reload
FROM golang:1.25rc1-alpine

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Install reflex for hot reload
RUN go install github.com/cespare/reflex@latest

# Copy go mod files
COPY api/go.mod api/go.sum ./
RUN go mod download

# Install oapi-codegen
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest

# Expose port
EXPOSE 8080

# Create directory for generated code
RUN mkdir -p /app/gen

# Use reflex to watch for file changes and rebuild
# Watch .go files and regenerate OpenAPI code when openapi.yaml changes
CMD ["reflex", "-c", "/app/reflex.conf"]
