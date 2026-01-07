FROM golang:1.25-alpine AS base

WORKDIR /app

COPY go.mod go.sum ./

FROM base AS build

RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest AS runtine

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 2112

CMD ["./app"]
