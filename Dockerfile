FROM node:20-alpine AS frontend
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN mkdir -p /app/internal/interfaces/http/static
RUN npm run build

FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY --from=frontend /app/internal/interfaces/http/static ./internal/interfaces/http/static
COPY internal/ ./internal/
COPY cmd/ ./cmd/
RUN go build -o server ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/server ./
EXPOSE 8080
CMD ["./server"]
