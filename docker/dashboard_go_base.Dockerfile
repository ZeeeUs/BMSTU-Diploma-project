FROM golang:1.16.4

WORKDIR /app

COPY . .

RUN make build-go