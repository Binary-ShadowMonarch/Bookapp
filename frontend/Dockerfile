# syntax=docker/dockerfile:1.7

# ---- Dependencies Stage ----
FROM --platform=$BUILDPLATFORM node:20-alpine AS deps
WORKDIR /app
RUN corepack enable && corepack prepare pnpm@9.15.4 --activate
COPY package.json pnpm-lock.yaml .npmrc ./
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
	pnpm install --frozen-lockfile

# ---- Builder Stage ----
FROM --platform=$BUILDPLATFORM node:20-alpine AS builder
WORKDIR /app
RUN corepack enable && corepack prepare pnpm@9.15.4 --activate
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN pnpm build

# ---- Final Stage ----
FROM node:20-alpine AS final
WORKDIR /app
ENV NODE_ENV=production
RUN corepack enable && corepack prepare pnpm@9.15.4 --activate
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/pnpm-lock.yaml ./pnpm-lock.yaml
COPY --from=builder /app/.npmrc ./.npmrc
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
	pnpm install --prod --frozen-lockfile --ignore-scripts

# The Node adapter runs on port 3000 by default
EXPOSE 3000

CMD ["node", "build/index.js"]