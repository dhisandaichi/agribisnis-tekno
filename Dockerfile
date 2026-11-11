# ---------- build frontend ----------
FROM node:22-alpine AS web
WORKDIR /web
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ---------- build backend ----------
FROM golang:1.24-alpine AS api
WORKDIR /src
ENV GOTOOLCHAIN=local CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download
COPY backend/ ./backend/
RUN go build -trimpath -ldflags="-s -w" -o /out/server ./backend/cmd/server

# ---------- runtime ----------
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 app
ENV TZ=Asia/Jakarta PORT=8080 STATIC_DIR=/app/public
WORKDIR /app
COPY --from=api /out/server /app/server
COPY --from=web /web/dist /app/public
USER app
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s CMD wget -qO- http://127.0.0.1:8080/api/v1/health || exit 1
ENTRYPOINT ["/app/server"]
