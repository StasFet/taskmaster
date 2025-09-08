FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o app .

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/app .

EXPOSE 3001

CMD ["./app"]