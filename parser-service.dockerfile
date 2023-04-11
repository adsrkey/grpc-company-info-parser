FROM golang:1.20.1 as builder

WORKDIR /go-parser-service/

COPY . .

RUN CGO_ENABLED=0 go build -o service cmd/api/main.go

FROM alpine:latest

WORKDIR /go-parser-service

COPY --from=builder /go-parser-service/ /go-parser-service

EXPOSE 8080
EXPOSE 9090

CMD ./service