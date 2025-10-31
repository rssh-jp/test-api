# Build stage
FROM golang:1.25.3-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY api/go.mod api/go.sum ./
RUN go mod download

# Copy source code
COPY api/ ./
COPY resources/ /app/resources/

# Generate OpenAPI code
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest
RUN mkdir -p gen
RUN oapi-codegen -package gen -generate types,server,spec /app/resources/openapi/openapi.yaml > gen/openapi.gen.go

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
