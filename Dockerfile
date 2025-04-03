FROM golang:1.23.5-alpine as builder

ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev

WORKDIR /

COPY ./src /src
COPY ./views /views
COPY ./public /public
COPY ./database /database

WORKDIR /src

RUN go mod tidy

RUN go build -o main ./main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder / .

EXPOSE 25000

WORKDIR /root/src

CMD ["./main"]

