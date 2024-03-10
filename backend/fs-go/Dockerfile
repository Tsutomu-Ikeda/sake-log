FROM golang:alpine AS builder

ENV CGO_ENABLED=1

RUN apk add --no-cache git gcc musl-dev

WORKDIR /app
COPY . .
RUN go mod init example.com/myapp
RUN go get -d -v
RUN go build -o main .

FROM alpine

WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
