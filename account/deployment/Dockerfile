FROM golang:alpine AS base
WORKDIR /app
RUN apk add --no-cache git make
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY account/go.mod account/go.sum ./account/
WORKDIR /app/account
RUN go mod download

FROM base AS dev
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY pkg /app/pkg
COPY account /app/account
CMD ["sh", "-c", "swag init -g internal/server/server.go && air"]

FROM base AS test
RUN apk add --no-cache build-base
ENV CGO_ENABLED=1
COPY . .
CMD ["go", "test", "-v", "-race", "./tests/..."]

FROM base AS qa-builder
COPY pkg /app/pkg
COPY account /app/account
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN	swag init -g internal/app/server.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/api/main.go

FROM gcr.io/distroless/static-debian12 AS qa
WORKDIR /root/
COPY --from=qa-builder /app/account/server .
EXPOSE 8080
CMD ["./server"]

FROM base AS prod-builder
COPY pkg /app/pkg
COPY account /app/account
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/api/main.go

FROM gcr.io/distroless/static-debian12 AS prod
WORKDIR /root/
COPY --from=prod-builder /app/account/server .
EXPOSE 8080
CMD ["./server"]