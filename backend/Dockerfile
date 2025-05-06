FROM golang:1.24.2-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o backend-server .

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash

WORKDIR /root/
COPY --from=builder /app/backend-server .
COPY wait-for-it.sh /usr/local/bin/wait-for-it
RUN chmod +x /usr/local/bin/wait-for-it

CMD ["wait-for-it", "db:3306", "-t", "60", "--", "./backend-server"]