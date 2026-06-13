FROM golang:1.20 AS builder

WORKDIR /app

COPY . .

RUN go mod init uploader && go mod tidy && go build -o server .

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends gosu \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/server .
COPY upload.html .
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
CMD ["./server"]
