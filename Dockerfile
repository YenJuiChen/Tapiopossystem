# syntax=docker/dockerfile:1.4

# -------- frontend builder --------
    FROM node:20 AS frontend-builder
    WORKDIR /app/frontend
    
    COPY frontend/package.json frontend/yarn.lock ./
    RUN yarn install
    
    COPY frontend ./
    RUN yarn build
    
    # -------- backend builder --------
    FROM golang:1.22.3-alpine AS backend-builder
    WORKDIR /app
    
    # ✅ 設定交叉編譯平台（linux/amd64）
    ENV GOOS=linux
    ENV GOARCH=amd64
    ENV CGO_ENABLED=0
    
    COPY backend/go.mod backend/go.sum ./
    RUN go mod download
    
    COPY backend ./
    RUN go build -o /server .
    
    # -------- final stage --------
    FROM alpine:latest
    WORKDIR /root/
    
    # ✅ 複製後端 binary
    COPY --from=backend-builder /server .
    
    # ✅ 複製前端靜態檔到 /root/dist（Go 程式從 ./dist serve）
    COPY --from=frontend-builder /app/frontend/dist ./dist
    
    # ✅ 若有 template 也一起複製
    COPY backend/templates ./templates
    
    # ✅ 安裝時區（UTC+8）
    RUN apk add --no-cache tzdata && \
        cp /usr/share/zoneinfo/Asia/Taipei /etc/localtime && \
        echo "Asia/Taipei" > /etc/timezone
    
    EXPOSE 80
    CMD ["./server"]
    