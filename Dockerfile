FROM golang:1.16

WORKDIR /app

COPY . ./

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER buser WITH SUPERUSER PASSWORD 'bpassword';" &&\
    createdb -O buser bdb &&\
    /etc/init.d/postgresql stop


RUN go mod download

RUN go build -o dashboard ./cmd

CMD [ "./dashboard" ]
