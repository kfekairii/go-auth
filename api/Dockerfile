FROM golang:alpine as builder

WORKDIR /go/src/app

RUN go install github.com/cespare/reflex@latest

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY .. .

RUN go build -o ./run .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /go/src/app/run .

EXPOSE 3001

CMD ["./run"]

