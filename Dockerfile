FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN go mod init uploader && go mod tidy && go build -o server .

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/server .
COPY upload.html .

EXPOSE 8080

CMD ["./server"]

