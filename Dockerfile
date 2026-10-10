# Stage 1: Build Frontend
FROM node:22-alpine AS web-builder

WORKDIR /web

COPY web/package*.json ./
RUN npm ci

COPY web/ ./
RUN npm run build

# Stage 2: Build Go backend with embedded frontend
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download -x

COPY . .
# Copy frontend build output from web-builder
COPY --from=web-builder /web/dist /app/web/dist

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/simpleschedule ./cmd/main.go


FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/simpleschedule .

RUN mkdir -p /app/storage

EXPOSE 2811

ENTRYPOINT ["./simpleschedule"]