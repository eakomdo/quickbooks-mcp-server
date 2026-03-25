FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source and build
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o qbo-mcp ./cmd/qbo-mcp


FROM alpine:3.20

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/qbo-mcp .

# The MCP server communicates over stdio by default.
# Set PORT or MCP_TRANSPORT=http for HTTP mode.
ENTRYPOINT ["./qbo-mcp"]

