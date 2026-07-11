# ---- build ----
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/api

# ---- run ----
FROM alpine:3.20
RUN apk add --no-cache wget && adduser -D -u 10001 appuser
WORKDIR /app
COPY --from=builder /out/server .
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s CMD wget -qO- http://localhost:8080/health || exit 1
ENTRYPOINT ["./server"]
