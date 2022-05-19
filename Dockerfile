FROM golang:1.16

WORKDIR /app

COPY . ./
RUN go mod download

RUN go build -o dashboard ./cmd

CMD [ "./dashboard" ]
