FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download -x

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/simpleschedule ./cmd/main.go

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/simpleschedule .

RUN mkdir -p /app/storage

EXPOSE 2811

ENTRYPOINT ["./simpleschedule"]