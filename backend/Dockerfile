FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/stockify ./cmd/server/main.go

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /go/bin/stockify_datasync ./cmd/datasync/main.go

FROM alpine:latest

RUN apk add --no-cache curl

RUN addgroup -S stockifygroup && adduser -S stockifyuser -G stockifygroup

COPY --from=builder /go/bin/stockify /app/stockify

COPY --from=builder /go/bin/stockify_datasync /app/stockify_datasync

COPY ./entrypoint.sh /app/entrypoint.sh

RUN chmod +x /app/entrypoint.sh

USER stockifyuser

EXPOSE 80

ENTRYPOINT ["/app/entrypoint.sh"]

CMD ["/app/stockify"]
