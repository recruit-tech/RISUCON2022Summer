FROM node:17 AS builder
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json frontend/tsconfig.json ./
RUN npm ci
COPY frontend ./

RUN npm run build && npm run export

FROM nginx:1.19
WORKDIR /app
COPY --from=builder /app/out .
COPY nginx /etc/nginx
