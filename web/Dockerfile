FROM oven/bun:1.3.13-alpine AS builder
WORKDIR /app

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .
RUN bun run build

FROM node:22-alpine
WORKDIR /app
RUN adduser -D -g '' appuser
COPY --from=builder /app/build ./build
USER appuser
EXPOSE 3000
ENV NODE_ENV=production
CMD ["node", "build/index.js"]
