# ---- Builder Stage ----
FROM golang:1.24.5-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

# ---- Final Stage ----
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
ENTRYPOINT ["/app/main"]