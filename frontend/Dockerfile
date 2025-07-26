# ---- Dependencies Stage ----
FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN npm i -g pnpm && pnpm install --frozen-lockfile

# ---- Builder Stage ----
FROM node:20-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm i -g pnpm && pnpm build

# ---- Final Stage ----
FROM node:20-alpine AS final
WORKDIR /app
ENV NODE_ENV=production
COPY --from=builder /app/build ./build
COPY --from=builder /app/package.json ./package.json
COPY --from=builder /app/pnpm-lock.yaml ./pnpm-lock.yaml

RUN npm i -g pnpm && pnpm install --prod --frozen-lockfile

# The Node adapter runs on port 3000 by default
EXPOSE 3000

CMD ["node", "build/index.js"]