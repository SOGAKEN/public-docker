# フロントエンドのビルド
FROM node:14 AS frontend-builder
WORKDIR /app/frontend
COPY frontend/package.json ./
RUN npm install
COPY frontend ./
RUN npm run build

# バックエンドのビルド
FROM golang:1.16 AS backend-builder
WORKDIR /app/backend
COPY backend/go.mod ./
RUN go mod download
COPY backend ./
RUN CGO_ENABLED=0 GOOS=linux go build -o main

# 最終イメージ
FROM golang:1.16
WORKDIR /app
COPY --from=backend-builder /app/backend/main ./
COPY --from=frontend-builder /app/frontend/.next ./frontend/.next
COPY --from=frontend-builder /app/frontend/public ./frontend/public

# .envファイルをコピー
COPY .env ./

CMD ["./main"]
