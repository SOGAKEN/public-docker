FROM golang:1.22.0-alpine as backend

WORKDIR /app/backend

COPY backend/go.mod .
COPY backend/go.sum .
RUN go mod download
COPY backend/ .
RUN go build -o main

FROM node:18-alpine as frontend

WORKDIR /app/frontend

COPY frontend/package.json .
COPY frontend/package-lock.json .
RUN npm install
COPY frontend/ .
RUN npm run build

FROM node:18-alpine

WORKDIR /app

COPY --from=backend /app/backend/main .
COPY --from=frontend /app/frontend/.next ./frontend/.next
COPY --from=frontend /app/frontend/node_modules ./frontend/node_modules
COPY --from=frontend /app/frontend/package.json ./frontend/package.json
COPY --from=frontend /app/frontend/public ./frontend/public
COPY --from=frontend /app/frontend/ ./frontend/

EXPOSE 8080

CMD ["sh", "-c", "./main & cd frontend && npm run start"]
