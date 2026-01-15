FROM golang:alpine AS base
WORKDIR /app
RUN apk add --no-cache git make
COPY pkg/go.mod pkg/go.sum ./pkg/
COPY payment/go.mod payment/go.sum ./payment/
COPY order/go.mod order/go.sum ./order/
WORKDIR /app/payment
RUN go mod download

FROM base AS dev
RUN go install github.com/air-verse/air@latest
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY pkg /app/pkg
COPY payment /app/payment
COPY order/ /app/order
CMD ["sh", "-c", "swag init -g internal/server/server.go && air"]

FROM base AS test
RUN apk add --no-cache build-base
ENV CGO_ENABLED=1
COPY . .
CMD ["go", "test", "-v", "-race", "./tests/..."]

FROM base AS qa-builder
COPY pkg /app/pkg
COPY payment /app/payment
COPY order/ /app/order
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN	swag init -g internal/app/server.go
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/api/main.go

FROM gcr.io/distroless/static-debian12 AS qa
WORKDIR /root/
COPY --from=qa-builder /app/payment/server .
EXPOSE 8080
CMD ["./server"]

FROM base AS prod-builder
COPY pkg /app/pkg
COPY payment /app/payment
COPY order/ /app/order
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/api/main.go

FROM gcr.io/distroless/static-debian12 AS prod
WORKDIR /root/
COPY --from=prod-builder /app/payment/server .
EXPOSE 8080
CMD ["./server"]