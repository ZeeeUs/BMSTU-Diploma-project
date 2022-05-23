FROM golang:1.16 AS build

ADD . /app
WORKDIR /app
RUN go build ./cmd/main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata
ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get install -y python3
RUN apt-get install -y python3-pip
RUN pip install Faker
RUN pip install psycopg2-binary

RUN apt-get install -y redis-server

RUN apt-get -y update && apt-get install -y postgresql-12
USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER buser WITH SUPERUSER PASSWORD 'bpassword';" &&\
    createdb -O buser bdb &&\
    /etc/init.d/postgresql stop

EXPOSE 5432
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=build /app/main/ .

EXPOSE 5000
ENV PGPASSWORD bpassword
CMD service postgresql start && psql -h localhost -d bdb -U buser -p 5432 -a -q -f ./db/dump.sql && redis-server --port 7500 --daemonize yes && python3 ./py_script/main.py && ./main