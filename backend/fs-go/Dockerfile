FROM golang:1.22 as builder

ENV CGO_ENABLED=0

WORKDIR /app
COPY ./main.go .
RUN go mod init example.com/myapp
RUN go get -d -v
RUN go build -o main .

FROM gcr.io/distroless/static AS distroless

WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
