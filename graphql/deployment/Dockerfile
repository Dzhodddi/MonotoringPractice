FROM golang:alpine AS base
WORKDIR /app
RUN apk add --no-cache git make
COPY go.work go.work
COPY account/go.mod account/go.sum ./account/
COPY product/go.mod  product/go.sum ./product/
COPY order/go.mod  order/go.sum ./order/
#COPY recommender/go.mod  recommender/go.sum ./recommender/
COPY payment/go.mod  payment/go.sum ./payment/
COPY graphql/go.mod  graphql/go.sum ./graphql/
COPY pkg/go.mod pkg/go.sum ./pkg/

WORKDIR /app/graphql
RUN go mod download

FROM base AS dev
RUN go install github.com/air-verse/air@latest
COPY pkg /app/pkg
COPY account /app/account
COPY product /app/product
COPY order /app/order
#COPY recommender /app/recommender
COPY payment /app/payment
COPY graphql /app/graphql
CMD ["sh", "-c", "air"]

FROM base AS test
RUN apk add --no-cache build-base
ENV CGO_ENABLED=1
COPY . .
CMD ["go", "test", "-v", "-race", "./tests/..."]

FROM base AS qa-builder
COPY pkg /app/pkg
COPY account /app/account
COPY product /app/product
COPY order /app/order
#COPY recommender /app/recommender
COPY payment /app/payment
COPY graphql /app/graphql
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
COPY product /app/product
COPY order /app/order
#COPY recommender /app/recommender
COPY payment /app/payment
COPY graphql /app/graphql
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server cmd/api/main.go

FROM gcr.io/distroless/static-debian12 AS prod
WORKDIR /root/
COPY --from=prod-builder /app/graphql/server .
EXPOSE 8080
CMD ["./server"]