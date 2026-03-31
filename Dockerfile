FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /server ./cmd/server

FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /server /app/server
COPY configs/config.yaml /app/configs/config.yaml

ENV CONFIG_PATH=/app/configs/config.yaml
ENV TZ=Asia/Shanghai

EXPOSE 8881

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8881/health || exit 1

CMD ["/app/server"]