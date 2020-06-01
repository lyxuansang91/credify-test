# Builder
FROM golang:1.14.2-alpine3.11 as builder

RUN apk update && apk upgrade && \
    apk --update add git make

WORKDIR /app

COPY go.mod go.sum /app/



COPY . .

RUN make engine

# Distribution
FROM alpine:latest

RUN apk update && apk add bash
RUN apk add --no-cache ca-certificates openssl tzdata



WORKDIR /app 

COPY config.json .
COPY entry-point.sh .

EXPOSE 9090

COPY --from=builder /app/credify-test .

CMD ["./entry-point.sh"]
