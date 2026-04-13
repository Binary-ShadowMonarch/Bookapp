# syntax=docker/dockerfile:1.7

# ---- Builder Stage ----
FROM --platform=$BUILDPLATFORM golang:1.24.5-bookworm AS builder
ARG TARGETOS
ARG TARGETARCH
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o main ./cmd/server/main.go

# ---- Final Stage ----
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
ENTRYPOINT ["/app/main"]