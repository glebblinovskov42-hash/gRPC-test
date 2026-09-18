FROM golang:1.27.1-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/exe ./server

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/exe /app/exe

CMD [ "/app/exe" ]